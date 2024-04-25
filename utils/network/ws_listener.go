package network

/**
 * @Author: lee
 * @Description:
 * @File: websocket_client
 * @Date: 2021/9/9 11:24 上午
 */
import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/juju/ratelimit"
	"github.com/spf13/cast"
	"gitlab.xbit.trade/blockchain/blockchain-core/utils/dumputils"
	"gitlab.xbit.trade/blockchain/blockchain-core/utils/logutils"
	"go.uber.org/zap"
	"io"
	"net/http"
	"time"
)

const (
	OptionLimitDuration = "limit-duration"
	OptionMaxLimit      = "max-limit"
	OptionMaxFrequent   = "max-frequent"
	OptionPongWait      = "pong-wait"
	OptionWriteWait     = "write-wait"
)

type OptionFunc func(ls *WebsocketListener)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 51200,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type WebsocketListener struct {
	// The websocket connection.
	conn   *websocket.Conn
	bucket *ratelimit.Bucket

	limitDuration time.Duration
	maxLimit      int
	maxFreqCount  int
	pongWait      time.Duration
	writeWait     time.Duration
	Ctx           context.Context

	send chan []byte

	OnMessage            func(ctx context.Context, conn *websocket.Conn, msg []byte)        //读接口回掉
	OnSend               func(ctx context.Context, conn *websocket.Conn, msg []byte) []byte //写接口回掉
	OnCloseReadCallback  func()
	OnCloseWriteCallback func()

	ip      string
	wsid    string
	isClose bool
	//isZip   bool
}

func NewWebsocketListener(w http.ResponseWriter, r *http.Request, options ...OptionFunc) (*WebsocketListener, error) {
	maxLimit := 50
	limitDuration := 2 * time.Second
	maxFrequent := 20
	pongWait := 30 * time.Second
	writeWait := 5 * time.Second

	ip := GetHttpRequestIP(r)

	requestId := r.Header.Get(HeaderXRequestID)
	fields := []zap.Field{
		zap.String("c-id", requestId),
	}
	ctx := context.WithValue(r.Context(), logutils.LogCtxKey, fields)

	conn, err := upgrader.Upgrade(w, r, nil)
	if nil != err {
		return nil, fmt.Errorf("upgrade err: %s", err.Error())
	}

	ret := &WebsocketListener{
		conn:          conn,
		limitDuration: limitDuration,
		maxLimit:      maxLimit,
		maxFreqCount:  maxFrequent,
		pongWait:      pongWait,
		writeWait:     writeWait,

		send: make(chan []byte, 1024),

		ip:   ip,
		wsid: requestId,
		Ctx:  ctx,
	}

	for _, option := range options {
		if nil != option {
			option(ret)
		}
	}

	ret.bucket = ratelimit.NewBucketWithQuantum(limitDuration, int64(maxLimit), int64(maxLimit))

	return ret, nil
}

func (c *WebsocketListener) DoJob() {
	go c.readPump()
	go c.writePump()
	logutils.WithContext(c.Ctx).Info("client connected")
}

func (c *WebsocketListener) IsClosed() bool {
	return c.isClose
}

func (c *WebsocketListener) WriteMessage(msg []byte) error {
	//newMsg := msg
	//var err error
	//if c.isZip {
	//	newMsg, err = GzipCompress(msg)
	//	if nil != err {
	//		return err
	//	}
	//}

	c.send <- msg
	return nil
}

func (c *WebsocketListener) readPump() {
	defer dumputils.SkipPanic(func() {
		logutils.WithContext(c.Ctx).Debug("readPump defer")
		if nil != c.OnCloseReadCallback {
			c.OnCloseReadCallback()
		}
		_ = c.conn.WriteMessage(websocket.CloseMessage, []byte("read exit"))
		_ = c.conn.Close()
		c.isClose = true
	})

	cnFrequent := 0
	pongWaitTime := time.Now().Add(c.pongWait)

	//c.conn.SetReadLimit(maxMessageSize)
	_ = c.conn.SetReadDeadline(pongWaitTime)

	//客户端主动发ping
	c.conn.SetPingHandler(func(string) error {
		pongWaitTime = time.Now().Add(c.pongWait)
		_ = c.conn.WriteMessage(websocket.PongMessage, []byte("pong"))
		return c.conn.SetReadDeadline(pongWaitTime)
	})
	for {
		if c.isClose {
			return
		}
		cnFrequent = 0
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				logutils.WithContext(c.Ctx).Debug("readPump IsUnexpectedCloseError", zap.Error(err))
				break
			}
			logutils.WithContext(c.Ctx).Debug("readPump ReadMessage failed", zap.Error(err))
			break
		}

		if c.bucket.TakeAvailable(1) <= 0 {
			cnFrequent++
			if cnFrequent > c.maxFreqCount {
				logutils.WithContext(c.Ctx).Warn("readPump|reach max_frequent_count")
				break
			}

			c.responseTooMuch()
			continue
		}

		if nil != c.OnMessage {
			c.OnMessage(c.Ctx, c.conn, message)
		}
	}
}

func (c *WebsocketListener) writePump() {
	defer dumputils.SkipPanic(func() {
		logutils.WithContext(c.Ctx).Warn("writePump defer")
		if nil != c.OnCloseWriteCallback {
			c.OnCloseWriteCallback()
		}
		_ = c.conn.Close()
		c.isClose = true
	})

	for {
		select {
		case message, ok := <-c.send:
			if c.isClose {
				return
			}
			if !ok {
				logutils.WithContext(c.Ctx).Debug("writePump read send buffer failed")
				_ = c.conn.WriteMessage(websocket.CloseMessage, []byte("no message to send"))
				return
			}
			now := time.Now()
			_ = c.conn.SetWriteDeadline(now.Add(c.writeWait))

			w, err := c.nextWriter()
			if err != nil {
				logutils.WithContext(c.Ctx).Warn("nextWriter err", zap.Error(err))
				_ = c.conn.WriteMessage(websocket.CloseMessage, []byte("no writer"))
				return
			}

			if nil != c.OnSend {
				message = c.OnSend(c.Ctx, c.conn, message)
			}

			n, err := c.responseMsg(w, message)
			if nil != err {
				logutils.WithContext(c.Ctx).Error("writePump responseMsg err", zap.Error(err))
				_ = c.conn.WriteMessage(websocket.CloseMessage, []byte("write failed"))
				return
			}

			since := time.Since(now)
			if since > 1*time.Second {
				logutils.WithContext(c.Ctx).Warn("writePump slow speed", zap.Duration("d", since))
			}
			logutils.WithContext(c.Ctx).Debug("writePump finish", zap.Int("len", n), zap.Duration("d", since))

		}
	}
}

func (c *WebsocketListener) responseMsg(w io.WriteCloser, message []byte) (int, error) {
	n, err := w.Write(message)
	if nil != err {
		logutils.WithContext(c.Ctx).Debug("writePump Write failed", zap.Error(err))
	}

	if err = w.Close(); nil != err {
		return n, fmt.Errorf("writePump Close err: %s", err.Error())
	}

	return n, nil
}

func (c *WebsocketListener) responseTooMuch() {
	ack := &AckWsServerBase{
		Event: EventError,
		Code:  ErrorTooMuch,
		Msg:   ErrMsgTooMuch,
	}

	respBody, _ := json.Marshal(ack)

	_ = c.WriteMessage(respBody)
}

func GzipCompress(src []byte) ([]byte, error) {
	var in bytes.Buffer
	w := gzip.NewWriter(&in)
	_, err := w.Write(src)
	if err != nil {
		return nil, err
	}
	_ = w.Close()
	return in.Bytes(), nil
}

func WithOption(name string, value interface{}) OptionFunc {
	switch name {
	case OptionLimitDuration:
		return func(c *WebsocketListener) {
			limitDuration, err := cast.ToDurationE(value)
			if nil == err {
				c.limitDuration = limitDuration
			}
		}

	case OptionMaxLimit:
		return func(c *WebsocketListener) {
			maxLimit, err := cast.ToIntE(value)
			if nil == err {
				c.maxLimit = maxLimit
			}
		}

	case OptionMaxFrequent:
		return func(c *WebsocketListener) {
			maxFrequent, err := cast.ToIntE(value)
			if nil == err {
				c.maxFreqCount = maxFrequent
			}
		}

	case OptionPongWait:
		return func(c *WebsocketListener) {
			pongWait, err := cast.ToDurationE(value)
			if nil == err {
				c.pongWait = pongWait
			}
		}

	case OptionWriteWait:
		return func(c *WebsocketListener) {
			writeWait, err := cast.ToDurationE(value)
			if nil == err {
				c.pongWait = writeWait
			}
		}
	}

	return nil
}

func (c *WebsocketListener) nextWriter() (io.WriteCloser, error) {
	w, err := c.conn.NextWriter(c.getMessageType())
	if err != nil {
		return nil, err
	}
	return w, err
}

func (c *WebsocketListener) getMessageType() int {
	return websocket.TextMessage
}
