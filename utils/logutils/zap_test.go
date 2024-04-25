package logutils

import (
	"context"
	"go.uber.org/zap"
	"testing"
)

func Test_Log(t *testing.T) {
	def := DefaultZapConfig
	def.WarnName = "warn"
	def.ErrorName = "error"
	InitLogger(def)

	Info("info")
	Error("err")
	Panic("panic")
}

func Test_Ctx(t *testing.T) {

	reqId := zap.Int("reqId", 123)
	reqReq := zap.String("req", "123")
	fields := make([]zap.Field, 0, 2)
	fields = append(fields, reqReq)
	fields = append(fields, reqId)
	ctx := context.WithValue(context.TODO(), LogCtxKey, fields)
	InitLogger(DefaultZapConfig)
	logger := WithContext(ctx)
	logger.Info("info", zap.String("symbol", "btc"))
	Info("success")

}
