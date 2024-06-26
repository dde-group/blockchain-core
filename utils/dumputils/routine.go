package dumputils

import (
	"bytes"
	"errors"
	"github.com/dde-group/blockchain-core/utils"
	"github.com/dde-group/blockchain-core/utils/logutils"
	"go.uber.org/zap"
	"runtime"
)

/**
 * @Author: lee
 * @Description:
 * @File: routine
 * @Date: 2021/10/14 2:53 下午
 */

// HandlePanic
/* @Description: 第一个参数是函数，后面是参数
 * @param params ...interface{}
 */
func HandlePanic(params ...interface{}) {
	var err error
	var stack string
	if r := recover(); nil != r {
		stack = string(PanicTrace(4))
		switch r.(type) {
		case error:
			err = r.(error)
			break
		case string:
			err = errors.New(r.(string))
			break
		default:
			err = errors.New("Unknown panic")
		}

		pc := make([]uintptr, 1)
		numFrames := runtime.Callers(5, pc)
		if numFrames < 1 {
			return
		}

		frame, _ := runtime.CallersFrames(pc).Next()
		//log.Println("rame function, file, line", frame.Function, frame.File, frame.Line)

		//log.Println("panic stack:\n "+stack+"\n", err.Error())
		//params = append(params, err)
		utils.Invoke0(params...)

		logutils.Error("frame function, file, line", zap.String("func", frame.Function), zap.String("file", frame.File),
			zap.Int("line", frame.Line))
		logutils.Panic("panic stack:", zap.Error(err), zap.String("stack", stack))

	}
}

func SkipPanic(params ...interface{}) {
	var err error
	var stack string
	if r := recover(); nil != r {
		stack = string(PanicTrace(4))
		switch r.(type) {
		case error:
			err = r.(error)
			break
		case string:
			err = errors.New(r.(string))
			break
		default:
			err = errors.New("Unknown panic")
		}

		pc := make([]uintptr, 1)
		numFrames := runtime.Callers(5, pc)
		if numFrames < 1 {
			return
		}

		frame, _ := runtime.CallersFrames(pc).Next()
		//log.Println("rame function, file, line", frame.Function, frame.File, frame.Line)

		//log.Println("panic stack:\n "+stack+"\n", err.Error())
		utils.Invoke0(params)

		logutils.Error("frame function, file, line", zap.String("func", frame.Function), zap.String("file", frame.File),
			zap.Int("line", frame.Line), zap.Error(err), zap.String("stack", stack))
	}
}

func PanicTrace(kb int) []byte {
	s := []byte("/src/runtime/panic.go")
	e := []byte("\ngoroutine ")
	line := []byte("\n")
	stack := make([]byte, kb<<10) //4KB
	length := runtime.Stack(stack, true)
	start := bytes.Index(stack, s)
	stack = stack[start:length]
	start = bytes.Index(stack, line) + 1
	stack = stack[start:]
	end := bytes.LastIndex(stack, line)
	if end != -1 {
		stack = stack[:end]
	}
	end = bytes.Index(stack, e)
	if end != -1 {
		stack = stack[:end]
	}
	stack = bytes.TrimRight(stack, "\n")
	return stack
}
