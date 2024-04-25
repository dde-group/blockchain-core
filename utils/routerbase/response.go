package routerbase

/**
 * @Author: lee
 * @Description:
 * @File: response
 * @Date: 2022/4/1 6:49 下午
 */
import (
	"github.com/gin-gonic/gin"

	"gitlab.xbit.trade/blockchain/blockchain-core/utils/logutils"
	"go.uber.org/zap"
	"net/http"
	"reflect"
)

// ResponseOK
/* @Description:
 * @param ack interface{}
 * @param c *gin.Context
 */
func ResponseOK(ack interface{}, c *gin.Context) {

	ackType := reflect.TypeOf(ack)
	ackValue := reflect.ValueOf(nil)
	if reflect.Ptr == ackType.Kind() {
		ackValue = reflect.ValueOf(ack).Elem()
	}

	kind := ackValue.Kind()
	if reflect.Struct == kind {
		resultValue := ackValue.FieldByName("Code")

		if resultValue.CanSet() && resultValue.IsValid() && reflect.Int == resultValue.Kind() {
			resultValue.SetInt(int64(Success))
		}
		msgValue := ackValue.FieldByName("Msg")
		if msgValue.CanSet() && msgValue.IsValid() && reflect.String == msgValue.Kind() {
			msgValue.SetString(ErrorMessage[Success])
		}
	}

	logutils.WithContext(c).Debug("success", zap.Any("ack", ack))

	c.JSON(http.StatusOK, ack)
}

func ResponseFailed(result int, c *gin.Context) {
	ack := AckBase{
		Code: result,
		Msg:  ErrorMessage[result],
	}
	logutils.WithContext(c).Warn("failed", zap.Any("ack", ack))
	c.JSON(http.StatusOK, ack)
}

func ResponseFailedWithRemark(result int, c *gin.Context, remark string) {
	msg := ErrorMessage[result]
	if "" == msg {
		msg = remark
	} else {
		msg += ", " + remark
	}
	ack := AckBase{
		Code: result,
		Msg:  msg,
	}
	logutils.WithContext(c).Warn("failed", zap.Any("ack", ack))
	c.JSON(http.StatusOK, ack)
}

func ResponseErrorWithAck(result int, c *gin.Context, ack interface{}, remark string) {
	ackType := reflect.TypeOf(ack)
	ackValue := reflect.ValueOf(nil)
	if reflect.Ptr == ackType.Kind() {
		ackValue = reflect.ValueOf(ack).Elem()
	}

	msg := ErrorMessage[result]
	if "" == msg {
		msg = remark
	} else {
		msg += ", " + remark
	}

	kind := ackValue.Kind()
	if reflect.Struct == kind {
		resultValue := ackValue.FieldByName("Code")

		if resultValue.CanSet() && resultValue.IsValid() && reflect.Int == resultValue.Kind() {
			resultValue.SetInt(int64(result))
		}
		msgValue := ackValue.FieldByName("Msg")
		if msgValue.CanSet() && msgValue.IsValid() && reflect.String == msgValue.Kind() {
			msgValue.SetString(msg)
		}
	}

	logutils.WithContext(c).Warn("failed", zap.Any("ack", ack))

	c.JSON(http.StatusOK, ack)
}
