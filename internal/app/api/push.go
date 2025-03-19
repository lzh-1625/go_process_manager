package api

import (
	"github.com/lzh-1625/go_process_manager/internal/app/model"
	"github.com/lzh-1625/go_process_manager/internal/app/repository"

	"github.com/gin-gonic/gin"
)

type pushApi struct{}

var PushApi = new(pushApi)

func (p *pushApi) GetPushList(ctx *gin.Context) {
	rOk(ctx, "Query successful!", repository.PushRepository.GetPushList())
}

func (p *pushApi) GetPushById(ctx *gin.Context) {
	id := getQueryInt(ctx, "id")
	rOk(ctx, "Query successful!", repository.PushRepository.GetPushConfigById(id))
}

func (p *pushApi) AddPushConfig(ctx *gin.Context) {
	req := bind[model.Push](ctx)
	err := repository.PushRepository.AddPushConfig(req)
	errCheck(ctx, err != nil, err)
	rOk(ctx, "Operation successful!", nil)
}

func (p *pushApi) UpdatePushConfig(ctx *gin.Context) {
	req := bind[model.Push](ctx)
	err := repository.PushRepository.UpdatePushConfig(req)
	errCheck(ctx, err != nil, err)
	rOk(ctx, "Operation successful!", nil)
}

func (p *pushApi) DeletePushConfig(ctx *gin.Context) {
	id := getQueryInt(ctx, "id")
	err := repository.PushRepository.DeletePushConfig(id)
	errCheck(ctx, err != nil, err)
	rOk(ctx, "Operation successful!", nil)
}
