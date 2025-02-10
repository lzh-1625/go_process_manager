package api

import (
	"github.com/lzh-1625/go_process_manager/internal/app/logic"
	"github.com/lzh-1625/go_process_manager/internal/app/model"
	"github.com/lzh-1625/go_process_manager/internal/app/repository"

	"github.com/gin-gonic/gin"
)

type procApi struct{}

var ProcApi = new(procApi)

func (p *procApi) CreateNewProcess(ctx *gin.Context, req model.Process) {
	index, err := repository.ProcessRepository.AddProcessConfig(req)
	errCheck(ctx, err != nil, err)
	req.Uuid = index
	proc, err := logic.ProcessCtlLogic.RunNewProcess(req)
	errCheck(ctx, err != nil, err)
	logic.ProcessCtlLogic.AddProcess(req.Uuid, proc)
	rOk(ctx, "Operation successful!", gin.H{
		"id": req.Uuid,
	})
}

func (p *procApi) DeleteNewProcess(ctx *gin.Context) {
	uuid := getQueryInt(ctx, "uuid")
	logic.ProcessCtlLogic.KillProcess(uuid)
	logic.ProcessCtlLogic.DeleteProcess(uuid)
	err := repository.ProcessRepository.DeleteProcessConfig(uuid)
	errCheck(ctx, err != nil, err)
	rOk(ctx, "Operation successful!", nil)
}

func (p *procApi) KillProcess(ctx *gin.Context) {
	uuid := getQueryInt(ctx, "uuid")
	err := logic.ProcessCtlLogic.KillProcess(uuid)
	errCheck(ctx, err != nil, err)
	rOk(ctx, "Operation successful!", nil)
}

func (p *procApi) StartProcess(ctx *gin.Context) {
	uuid := getQueryInt(ctx, "uuid")
	prod, err := logic.ProcessCtlLogic.GetProcess(uuid)
	if err != nil { // 进程不存在则创建
		proc, err := logic.ProcessCtlLogic.RunNewProcess(repository.ProcessRepository.GetProcessConfigById(uuid))
		errCheck(ctx, err != nil, err)
		logic.ProcessCtlLogic.AddProcess(uuid, proc)
		rOk(ctx, "Operation successful!", nil)
		return
	}
	errCheck(ctx, prod.State.State == 1, "The process is currently running.")
	prod.ResetRestartTimes()
	err = prod.Start()
	errCheck(ctx, err != nil, err)
	rOk(ctx, "Operation successful!", nil)
}

func (p *procApi) StartAllProcess(ctx *gin.Context) {
	if isAdmin(ctx) {
		logic.ProcessCtlLogic.ProcessStartAll()
	} else {
		logic.ProcessCtlLogic.ProcesStartAllByUsername(getUserName(ctx))
	}
	rOk(ctx, "Operation successful!", nil)
}

func (p *procApi) KillAllProcess(ctx *gin.Context) {
	if isAdmin(ctx) {
		logic.ProcessCtlLogic.KillAllProcess()
	} else {
		logic.ProcessCtlLogic.KillAllProcessByUserName(getUserName(ctx))
	}
	rOk(ctx, "Operation successful!", nil)
}

func (p *procApi) GetProcessList(ctx *gin.Context) {
	if isAdmin(ctx) {
		rOk(ctx, "Query successful!", logic.ProcessCtlLogic.GetProcessList())
	} else {
		rOk(ctx, "Query successful!", logic.ProcessCtlLogic.GetProcessListByUser(getUserName(ctx)))
	}
}

func (p *procApi) UpdateProcessConfig(ctx *gin.Context, req model.Process) {
	logic.ProcessCtlLogic.UpdateProcessConfig(req)
	err := repository.ProcessRepository.UpdateProcessConfig(req)
	errCheck(ctx, err != nil, err)
	rOk(ctx, "Operation successful!", nil)
}

func (p *procApi) GetProcessConfig(ctx *gin.Context) {
	uuid := getQueryInt(ctx, "uuid")
	data := repository.ProcessRepository.GetProcessConfigById(uuid)
	errCheck(ctx, data.Uuid == 0, "No information found!")
	rOk(ctx, "Query successful!", data)
}

func (p *procApi) ProcessControl(ctx *gin.Context) {
	user := getUserName(ctx)
	uuid := getQueryInt(ctx, "uuid")
	proc, err := logic.ProcessCtlLogic.GetProcess(uuid)
	errCheck(ctx, err != nil, err)
	proc.ProcessControl(user)
	rOk(ctx, "Operation successful!", nil)
}
