package api

import (
	"msm/model"

	"msm/dao"

	"github.com/gin-gonic/gin"
)

var PermissionApi = new(permissionApi)

type permissionApi struct{}

func (p *permissionApi) EditPermssion(ctx *gin.Context) {
	per := model.Permission{}
	err := ctx.ShouldBindJSON(&per)
	errCheck(ctx, err != nil, err)
	err = dao.PermissionDao.EditPermssion(per)
	errCheck(ctx, err != nil, err)
	rOk(ctx, "权限修改成功", nil)
}

func (p *permissionApi) GetPermissionList(ctx *gin.Context) {
	result := dao.PermissionDao.GetPermssionList(ctx.Query("account"))
	rOk(ctx, "查询成功", result)
}
