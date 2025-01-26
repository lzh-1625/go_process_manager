package api

import (
	"github.com/lzh-1625/go_process_manager/internal/app/service"

	"github.com/gin-gonic/gin"
)

type configApi struct{}

var ConfigApi = new(configApi)

func (c *configApi) GetSystemConfiguration(ctx *gin.Context) {
	result := service.ConfigService.GetSystemConfiguration()
	rOk(ctx, "Operation successful!", result)
}

func (c *configApi) SetSystemConfiguration(ctx *gin.Context, req map[string]string) {
	errCheck(ctx, service.ConfigService.SetSystemConfiguration(req) != nil, "Set config fail!")
	rOk(ctx, "Operation successful!", nil)
}

func (c *configApi) EsConfigReload(ctx *gin.Context) {
	errCheck(ctx, !service.EsService.InitEs(), "Incorrect username or password!")
	rOk(ctx, "Operation successful!", nil)
}
