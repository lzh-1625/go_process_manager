package api

import (
	"msm/consts/ctxflag"
	"msm/consts/role"
	"msm/dao"
	"msm/model"
	"msm/service/process"
	"strconv"

	"github.com/gin-gonic/gin"
)

type procApi struct{}

var ProcApi = new(procApi)

func (p *procApi) CreateNewProcess(ctx *gin.Context) {
	req := model.Process{}
	ctx.ShouldBindJSON(&req)
	index, err := dao.ProcessDao.AddProcessConfig(req)
	errCheck(ctx, err != nil, err)
	req.Uuid = index
	proc, err := process.RunNewProcess(req)
	errCheck(ctx, err != nil, err)
	process.ProcessCtlService.AddProcess(req.Uuid, proc)
	rOk(ctx, "创建成功", gin.H{
		"id": req.Uuid,
	})
}

func (p *procApi) DeleteNewProcess(ctx *gin.Context) {
	uuid, err := strconv.Atoi(ctx.Query("uuid"))
	errCheck(ctx, err != nil, "参数有误")
	process.ProcessCtlService.KillProcess(uuid)
	process.ProcessCtlService.DeleteProcess(uuid)
	err = dao.ProcessDao.DeleteProcessConfig(uuid)
	errCheck(ctx, err != nil, err)
	rOk(ctx, "删除成功", nil)
}

func (p *procApi) KillProcess(ctx *gin.Context) {
	uuid, err := strconv.Atoi(ctx.Query("uuid"))
	errCheck(ctx, err != nil, "参数有误")
	err = process.ProcessCtlService.KillProcess(uuid)
	errCheck(ctx, err != nil, err)
	rOk(ctx, "成功", nil)
}

func (p *procApi) StartProcess(ctx *gin.Context) {
	uuid, err := strconv.Atoi(ctx.Query("uuid"))
	errCheck(ctx, err != nil, "参数有误")
	prod, err := process.ProcessCtlService.GetProcess(uuid)
	if err != nil { // 进程不存在则创建
		proc, err := process.RunNewProcess(dao.ProcessDao.GetProcessConfigById(uuid))
		errCheck(ctx, err != nil, err)
		process.ProcessCtlService.AddProcess(uuid, proc)
		rOk(ctx, "成功", nil)
		return
	}
	errCheck(ctx, prod.GetStateState() == 1, "进程还在运行中")
	prod.ResetRestartTimes()
	prod.ReStart()
	// dao.UpdateServerAutoStart(uuid, true)
	rOk(ctx, "成功", nil)
}

func (p *procApi) GetProcessList(ctx *gin.Context) {
	if ctx.GetInt(ctxflag.ROLE) < int(role.USER) {
		rOk(ctx, "进程列表获取成功", process.ProcessCtlService.GetProcessList())
	} else {
		rOk(ctx, "进程列表获取成功", process.ProcessCtlService.GetProcessListByUser(ctx.GetString(ctxflag.USER_NAME)))
	}
}

func (p *procApi) UpdateProcessConfig(ctx *gin.Context) {
	req := model.Process{}
	ctx.ShouldBindJSON(&req)
	process.ProcessCtlService.UpdateProcessConfig(req)
	err := dao.ProcessDao.UpdateProcessConfig(req)
	errCheck(ctx, err != nil, err)
	rOk(ctx, "更改配置成功", nil)
}

func (p *procApi) GetProcessConfig(ctx *gin.Context) {
	uuid, err := strconv.Atoi(ctx.Query("uuid"))
	errCheck(ctx, err != nil, "参数有误")
	data := dao.ProcessDao.GetProcessConfigById(uuid)
	errCheck(ctx, data.Uuid == 0, "未查询到信息")
	rOk(ctx, "success", data)
}

func (p *procApi) ProcessControl(ctx *gin.Context) {
	user := ctx.GetString(ctxflag.USER_NAME)
	uuid, err := strconv.Atoi(ctx.Query("uuid"))
	errCheck(ctx, err != nil, "参数有误")
	proc, err := process.ProcessCtlService.GetProcess(uuid)
	errCheck(ctx, err != nil, "进程控制权获取失败")
	proc.ProcessControl(user)
	rOk(ctx, "获取进程控权成功", nil)
}
