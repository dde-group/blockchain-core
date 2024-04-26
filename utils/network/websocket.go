package network

/**
 * @Author: lee
 * @Description:
 * @File: websocket_client
 * @Date: 2021/9/9 11:24 上午
 */

import (
	"context"
	"fmt"
	"github.com/dde-group/blockchain-core/utils/logutils"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	MessagePrefix = "WsPrefix:"
)

type WebsocketAgent struct {
	NetAgentBase
	client      *websocket.Conn
	reqChan     chan string
	OnPing      func(string) error
	OnPong      func(string) error
	OnMessage   func(*WebsocketAgent, string)      //收到消息回调
	OnSend      func(*WebsocketAgent, int, string) //发送消息回调
	OnClose     func(*WebsocketAgent)
	OnConnected func() //连接被断开回调
	errConn     error
	sendElapse  int //发送消息时间间隔 单位ms 用于限频
	sendCache   []string
	header      http.Header
	Ctx         context.Context
}

type WsOptionFunc func(agent *WebsocketAgent)

func NewWebsocketAgent(host string, port uint, path string, isSecure bool, elapse int, options ...WsOptionFunc) *WebsocketAgent {
	hostUrl := ""
	trimHost := strings.TrimLeft(host, " ")

	if strings.HasPrefix(trimHost, "ws") && strings.Contains(trimHost, "://") {
		hostUrl = trimHost
	} else {
		if isSecure {
			hostUrl += "wss://" + trimHost
		} else {
			hostUrl += "ws://" + trimHost
		}
	}

	if 0 != port {
		hostUrl += ":" + strconv.FormatUint(uint64(port), 10)
	}

	hostUrl += path

	rawUrl, err := url.Parse(hostUrl)
	if nil != err {
		panic(err.Error())
	}

	ret := &WebsocketAgent{
		NetAgentBase: NetAgentBase{
			URL:      rawUrl,
			isAlive:  false,
			timeout:  5000,
			isClosed: false,
		},
		reqChan:    make(chan string, 128),
		sendCache:  make([]string, 0, 16),
		sendElapse: elapse,
		Ctx:        context.WithValue(context.TODO(), logutils.LogCtxKey, nil),
	}

	for _, option := range options {
		if nil != option {
			option(ret)
		}
	}

	return ret
}

func WithWsHeader(header http.Header) WsOptionFunc {
	return func(agent *WebsocketAgent) {
		agent.header = header
	}
}

func (ws *WebsocketAgent) SetPingHandler(handler func(string) error) {
	ws.client.SetPingHandler(handler)
}

func (ws *WebsocketAgent) SetPongHandler(handler func(string) error) {
	ws.client.SetPongHandler(handler)
}

func (ws *WebsocketAgent) SetCloseHandler(handler func(code int, text string) error) {
	ws.client.SetCloseHandler(handler)
}

func (ws *WebsocketAgent) Connect() {
	go func() {
		for {
			if ws.isClosed {
				break
			}

			if !ws.isAlive && !ws.isClosed {
				if err := ws.dial(); nil != err {
					logutils.WithContext(ws.Ctx).Warn("WebsocketAgent dial fatal", zap.Error(err), zap.String("url", ws.URL.String()))
				}
			}

			time.Sleep(time.Duration(ws.timeout) * time.Millisecond)
		}
	}()
	ws.doSendThread()
	ws.doReceiveThread()
}

func (ws *WebsocketAgent) Reconnect() {
	ws.isAlive = false
}

func (ws *WebsocketAgent) Close() error {
	ws.isClosed = true
	if nil != ws.client {
		return ws.client.Close()
	}
	return nil
}

func (ws *WebsocketAgent) Send(msg string) {
	//断线了就不发了减少sendMsg阻塞
	if !ws.isAlive {
		return
	}
	messageType := fmt.Sprintf("%02d", websocket.TextMessage)

	ws.reqChan <- MessagePrefix + messageType + msg
}

func (ws *WebsocketAgent) SendPongMsg(data []byte) {
	//断线了就不发了减少sendMsg阻塞
	if !ws.isAlive {
		return
	}
	messageType := fmt.Sprintf("%02d", websocket.PongMessage)
	ws.reqChan <- MessagePrefix + messageType + string(data)
}
func (ws *WebsocketAgent) SendPingMsg(data []byte) {
	//断线了就不发了减少sendMsg阻塞
	if !ws.isAlive {
		return
	}
	messageType := fmt.Sprintf("%02d", websocket.PingMessage)
	ws.reqChan <- MessagePrefix + messageType + string(data)
}

func (ws *WebsocketAgent) WaitForConnected() <-chan error {
	ret := make(chan error, 1)
	tick := time.Tick(100 * time.Millisecond)
	timer := time.NewTimer(10 * time.Second)
	go func() {

		for {
			select {
			case <-tick:
				{
					if ws.isAlive {
						ret <- nil
						return
					}
				}
			case <-timer.C:
				{
					errMsg := ""
					if nil != ws.errConn {
						errMsg = ws.errConn.Error()
					}
					err := fmt.Errorf("wait for websocket connect time out 30s, url: %s, err: %s ", ws.URL.String(), errMsg)
					ret <- err
					return
				}
			}
		}
	}()

	return ret
}

func (ws *WebsocketAgent) dial() error {
	var err error
	var client *websocket.Conn
	urlStr := ws.URL.String()

	header := ws.header
	if nil == header {
		header = http.Header{}
	}
	requestId := uuid.New().String()
	header.Add(HeaderXRequestID, requestId)

	logutils.Warn("dial websocket", zap.String("url", urlStr), zap.String("req-id", requestId))
	client, _, err = websocket.DefaultDialer.Dial(urlStr, header)

	if nil != err {
		ws.errConn = err
		logutils.WithContext(ws.Ctx).Error("dial websocket failed", zap.String("url", urlStr))
		return err
	}
	fields := []zap.Field{
		zap.String("c-id", requestId),
	}
	ws.Ctx = context.WithValue(context.TODO(), logutils.LogCtxKey, fields)
	ws.client = client

	if nil != ws.OnPing {
		ws.client.SetPingHandler(ws.OnPing)
	}

	if nil != ws.OnPong {
		ws.client.SetPongHandler(ws.OnPong)
	}

	ws.client.SetCloseHandler(func(code int, text string) error {
		ws.isAlive = false
		if nil != ws.OnClose {
			ws.OnClose(ws)
		}
		message := websocket.FormatCloseMessage(code, "")
		ws.client.WriteControl(websocket.CloseMessage, message, time.Now().Add(time.Second))
		return nil
	})

	//将alive设置提前，不能放在ws.OnConnected()后面，里面可能会发送消息，如果管道满了导致阻塞alive将不被设置，发送协程因alive未设置不发送消息了导致死锁了
	ws.isAlive = true

	if nil != ws.OnConnected {
		ws.OnConnected()
	}

	ws.errConn = nil

	return nil
}

func (ws *WebsocketAgent) doSendThread() {
	//logutils.Warn("doSendThread", zap.String("url", ws.URL.String()))
	go func() {
		for {
			if ws.isClosed {
				break
			}

			if !ws.isAlive {
				time.Sleep(100 * time.Millisecond)
				continue
			}

			msg := <-ws.reqChan

			ws.sendCache = append(ws.sendCache, msg)

			cnSuccess := 0
			for index, rawMsg := range ws.sendCache {
				if "" == rawMsg {
					cnSuccess++
					continue
				}
				messageType, sendMsg := ParseMessage(rawMsg)
				var err error
				switch messageType {
				case websocket.TextMessage, websocket.BinaryMessage:
					err = ws.client.WriteMessage(messageType, []byte(sendMsg))
				case websocket.PongMessage, websocket.PingMessage, websocket.CloseMessage:
					err = ws.client.WriteControl(messageType, []byte(sendMsg), time.Now().Add(time.Second))
				}

				if nil != err {
					ws.isAlive = false
					logutils.WithContext(ws.Ctx).Warn("doSendThread fatal", zap.String("url", ws.URL.String()), zap.Error(err))
					time.Sleep(100 * time.Millisecond)
					//控制消息不用重发了
					if messageType != websocket.TextMessage && messageType != websocket.BinaryMessage {
						ws.sendCache[index] = ""
					}
					break
				}

				//发送成功清空字符串
				ws.sendCache[index] = ""

				if nil != ws.OnSend {
					ws.OnSend(ws, messageType, sendMsg)
				}
				if ws.sendElapse > 0 {
					elapse := time.Duration(ws.sendElapse) * time.Millisecond
					time.Sleep(elapse)
				}
				cnSuccess++
			}

			//全部发送成功后清空缓存
			if cnSuccess == len(ws.sendCache) {
				ws.sendCache = ws.sendCache[0:0]
			}

		}
	}()
}

func (ws *WebsocketAgent) doReceiveThread() {
	//logutils.Warn("doReceiveThread", zap.String("url", ws.URL.String()))
	go func() {
		for {
			if ws.isClosed {
				break
			}
			if !ws.isAlive {
				time.Sleep(100 * time.Millisecond)
				continue
			}

			_, msg, err := ws.client.ReadMessage()
			if nil != err {
				ws.isAlive = false
				logutils.WithContext(ws.Ctx).Warn("doReceiveThread fatal", zap.String("url", ws.URL.String()), zap.Error(err))
				continue
			}

			ws.OnMessage(ws, string(msg))
		}
	}()
}

func ParseMessage(msg string) (int, string) {
	prefix := msg[0:len(MessagePrefix)]
	if prefix != MessagePrefix {
		return 0, ""
	}

	remain := msg[len(MessagePrefix):]
	msgType := remain[0:2]
	messageType, _ := strconv.ParseInt(msgType, 10, 64)
	remain = remain[2:]

	return int(messageType), remain
}
