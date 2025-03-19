package api

import (
	"github.com/lzh-1625/go_process_manager/internal/app/model"
	"github.com/lzh-1625/go_process_manager/internal/app/repository"

	"github.com/gin-gonic/gin"
)

var PermissionApi = new(permissionApi)

type permissionApi struct{}

func (p *permissionApi) EditPermssion(ctx *gin.Context) {
	req := bind[model.Permission](ctx)
	err := repository.PermissionRepository.EditPermssion(req)
	errCheck(ctx, err != nil, err)
	rOk(ctx, "Operation successful!", nil)
}

func (p *permissionApi) GetPermissionList(ctx *gin.Context) {
	result := repository.PermissionRepository.GetPermssionList(getQueryString(ctx, "account"))
	rOk(ctx, "Query successful!", result)
}
