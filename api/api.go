package api

import (
	"errors"
	"msm/consts/ctxflag"
	"net/http"

	"github.com/gin-gonic/gin"
)

func rOk(ctx *gin.Context, message string, data any) {
	jsonData := map[string]any{
		"code": 0,
		"msg":  message,
	}
	if data != nil {
		jsonData["data"] = data
	}
	ctx.JSON(http.StatusOK, jsonData)
}

func errCheck(ctx *gin.Context, isErr bool, errData any) {
	if !isErr {
		return
	}
	if err, ok := errData.(error); ok {
		ctx.Set(ctxflag.ERR, err)
	}
	if err, ok := errData.(string); ok {
		ctx.Set(ctxflag.ERR, errors.New(err))
	}
	panic(0)
}
