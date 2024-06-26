package api

import (
	"msm/model"
	"strconv"

	"msm/dao"

	"github.com/gin-gonic/gin"
)

type pushApi struct{}

var PushApi = new(pushApi)

func (p *pushApi) GetPushList(ctx *gin.Context) {
	rOk(ctx, "查询成功", dao.PushDao.GetPushList())
}

func (p *pushApi) GetPushById(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Query("id"))
	errCheck(ctx, err != nil, err)
	rOk(ctx, "查询成功", dao.PushDao.GetPushConfigById(id))
}

func (p *pushApi) AddPushConfig(ctx *gin.Context) {
	data := model.Push{}
	err := ctx.ShouldBindJSON(&data)
	errCheck(ctx, err != nil, err)
	err = dao.PushDao.AddPushConfig(data)
	errCheck(ctx, err != nil, err)
	rOk(ctx, "添加成功", nil)
}

func (p *pushApi) UpdatePushConfig(ctx *gin.Context) {
	data := model.Push{}
	err := ctx.ShouldBindJSON(&data)
	errCheck(ctx, err != nil, err)
	err = dao.PushDao.UpdatePushConfig(data)
	errCheck(ctx, err != nil, err)
	rOk(ctx, "更新成功", nil)
}

func (p *pushApi) DeletePushConfig(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Query("id"))
	errCheck(ctx, err != nil, err)
	err = dao.PushDao.DeletePushConfig(id)
	errCheck(ctx, err != nil, err)
	rOk(ctx, "删除成功", nil)
}
