package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Response struct {
	Code    ReCode
	Message string
	Data    interface{}
}

func RespError(ctx *gin.Context, code ReCode) {
	re := &Response{
		Code:    code,
		Message: code.Msg(),
		Data:    nil,
	}
	ctx.JSON(http.StatusBadRequest, re)
}

func RespErrorWithMsg(ctx *gin.Context, code ReCode, data interface{}) {
	re := &Response{
		Code:    code,
		Message: code.Msg(),
		Data:    data,
	}
	ctx.JSON(http.StatusBadRequest, re)
}

func RespSuccess(ctx *gin.Context, data interface{}) {
	re := &Response{
		Code:    CodeSuccess,
		Message: CodeSuccess.Msg(),
		Data:    data,
	}
	ctx.JSON(http.StatusOK, re)
}
