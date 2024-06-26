package api

import (
	"msm/config"
	"msm/model"
	"msm/service/es"

	"github.com/gin-gonic/gin"
)

type logApi struct{}

var LogApi = new(logApi)

func (a *logApi) GetLog(ctx *gin.Context) {
	req := model.GetLogReq{}
	errCheck(ctx, !config.CF.EsEnable, "elasticsearch未启用或账号密码错误")
	errCheck(ctx, ctx.ShouldBindJSON(&req) != nil, "请求体格式错误")
	rOk(ctx, "查询成功", es.EsService.Search(req))
}
