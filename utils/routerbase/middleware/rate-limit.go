package middleware

/**
 * @Author: lee
 * @Description:
 * @File: rate-limit
 * @Date: 2022-04-14 7:01 下午
 */

import (
	"github.com/gin-gonic/gin"
	"github.com/juju/ratelimit"
	"gitlab.xbit.trade/blockchain/blockchain-core/utils/dumputils"
	"gitlab.xbit.trade/blockchain/blockchain-core/utils/logutils"
	"gitlab.xbit.trade/blockchain/blockchain-core/utils/routerbase"
	"go.uber.org/zap"
	"runtime/debug"
	"time"
)

func RateLimit() func(c *gin.Context) {
	bucket := ratelimit.NewBucketWithQuantum(2*time.Second, 50, 50)
	return func(c *gin.Context) {
		// 如果取不到令牌就中断本次请求返回 rate limit
		if bucket.TakeAvailable(1) < 1 {

			routerbase.ResponseFailed(routerbase.ErrorRateLimit, c)
			c.Abort()
			return
		}
		c.Next()
	}
}

func Recover(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			//打印错误堆栈信息
			logutils.Error("panic", zap.Any("err", r))
			debug.PrintStack()
			logutils.Error("Stack \n" + string(dumputils.PanicTrace(4)))

			//封装通用json返回
			c.JSON(200, gin.H{
				"code": "4444",
				"msg":  "服务器内部错误",
			})
		}
	}()
	//加载完 defer recover，继续后续接口调用
	c.Next()
}

func LogContext(c *gin.Context) {
	requestId := c.GetHeader("X-Request-ID")
	path := c.Request.URL.Path
	fields := []zap.Field{
		zap.String("c-req-id", requestId),
		zap.String("c-path", path),
	}
	c.Set(logutils.LogCtxKey, fields)

	c.Next()
}
