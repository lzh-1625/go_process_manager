package api

import (
	"github.com/lzh-1625/go_process_manager/internal/app/model"
	"github.com/lzh-1625/go_process_manager/internal/app/repository"
	"github.com/lzh-1625/go_process_manager/internal/app/service"

	"github.com/gin-gonic/gin"
)

type taskApi struct{}

var TaskApi = new(taskApi)

func (t *taskApi) CreateTask(ctx *gin.Context, req model.Task) {
	err := service.TaskService.CreateTask(req)
	errCheck(ctx, err != nil, err)
	rOk(ctx, "Operation successful!", nil)
}

func (t *taskApi) GetTaskById(ctx *gin.Context) {
	result, err := repository.TaskRepository.GetTaskById(getQueryInt(ctx, "id"))
	errCheck(ctx, err != nil, "Query failed!")
	rOk(ctx, "Operation successful!", result)
}

func (t *taskApi) GetTaskList(ctx *gin.Context) {
	result := service.TaskService.GetAllTaskJob()
	rOk(ctx, "Operation successful!", result)
}

func (t *taskApi) DeleteTaskById(ctx *gin.Context) {
	err := service.TaskService.DeleteTask(getQueryInt(ctx, "id"))
	errCheck(ctx, err != nil, err)
	rOk(ctx, "Operation successful!", nil)
}

func (t *taskApi) StartTask(ctx *gin.Context) {
	go service.TaskService.RunTaskById(getQueryInt(ctx, "id"))
	rOk(ctx, "Operation successful!", nil)
}

func (t *taskApi) StopTask(ctx *gin.Context) {
	errCheck(ctx, service.TaskService.StopTaskJob(getQueryInt(ctx, "id")) != nil, "Operation failed!")
	rOk(ctx, "Operation successful!", nil)
}

func (t *taskApi) EditTask(ctx *gin.Context, req model.Task) {
	err := service.TaskService.EditTask(req)
	errCheck(ctx, err != nil, err)
	rOk(ctx, "Operation successful!", nil)
}

func (t *taskApi) EditTaskEnable(ctx *gin.Context, req model.Task) {
	err := service.TaskService.EditTaskEnable(req.Id, req.Enable)
	errCheck(ctx, err != nil, err)
	rOk(ctx, "Operation successful!", nil)
}

func (t *taskApi) RunTaskByKey(ctx *gin.Context) {
	err := service.TaskService.RunTaskByKey(ctx.Param("key"))
	errCheck(ctx, err != nil, err)
	rOk(ctx, "Operation successful!", nil)
}

func (t *taskApi) CreateTaskApiKey(ctx *gin.Context) {
	err := service.TaskService.CreateApiKey(getQueryInt(ctx, "id"))
	errCheck(ctx, err != nil, err)
	rOk(ctx, "Operation successful!", nil)
}
