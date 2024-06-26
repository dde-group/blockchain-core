package logutils

import (
	"context"
	"fmt"
	"go.uber.org/zap"
)

/**
 * @Author: lee
 * @Description:
 * @File: init
 * @Date: 2021/9/14 3:41 下午
 */

var (
	loggerModule ILogger
	logInit      = false

	errorNotInit = fmt.Errorf("log module not inited")
)

type ILogger interface {
	Info(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
	Warn(msg string, fields ...zap.Field)
	Debug(msg string, fields ...zap.Field)
	Fatal(msg string, fields ...zap.Field)
	DPanic(msg string, fields ...zap.Field)
	Panic(msg string, fields ...zap.Field)
	WithContext(ctx context.Context) ILogger
}

func InitLogger(config interface{}) {
	var err error
	if v, ok := config.(ZapConfig); ok {
		loggerModule, err = NewZapLogModule(v)
		if nil != err {
			panic(fmt.Errorf("zap log init fault, err: %s", err.Error()))
		}

		logInit = true
	}
}

func Logger() ILogger {
	return loggerModule.(*ZapLogModule)
}

func WithContext(ctx context.Context) ILogger {
	if !logInit {
		panic(errorNotInit)
	}
	return loggerModule.WithContext(ctx)
}

func Info(msg string, fields ...zap.Field) {
	if !logInit {
		panic(errorNotInit)
	}
	loggerModule.Info(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	if !logInit {
		panic(errorNotInit)
	}
	loggerModule.Warn(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	if !logInit {
		panic(errorNotInit)
	}
	loggerModule.Error(msg, fields...)
}

func Debug(msg string, fields ...zap.Field) {
	if !logInit {
		panic(errorNotInit)
	}
	loggerModule.Debug(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	if !logInit {
		panic(errorNotInit)
	}
	loggerModule.Fatal(msg, fields...)
}

func DPanic(msg string, fields ...zap.Field) {
	if !logInit {
		panic(errorNotInit)
	}
	loggerModule.DPanic(msg, fields...)
}

func Panic(msg string, fields ...zap.Field) {
	if !logInit {
		panic(errorNotInit)
	}
	loggerModule.Panic(msg, fields...)
}
