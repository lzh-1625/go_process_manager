package api

import (
	"github.com/lzh-1625/go_process_manager/internal/app/logic"
	"github.com/lzh-1625/go_process_manager/internal/app/model"
	"github.com/lzh-1625/go_process_manager/internal/app/repository"

	"github.com/gin-gonic/gin"
)

type taskApi struct{}

var TaskApi = new(taskApi)

func (t *taskApi) CreateTask(ctx *gin.Context, req model.Task) {
	err := logic.TaskLogic.CreateTask(req)
	errCheck(ctx, err != nil, err)
	rOk(ctx, "Operation successful!", nil)
}

func (t *taskApi) GetTaskById(ctx *gin.Context) {
	result, err := repository.TaskRepository.GetTaskById(getQueryInt(ctx, "id"))
	errCheck(ctx, err != nil, "Query failed!")
	rOk(ctx, "Operation successful!", result)
}

func (t *taskApi) GetTaskList(ctx *gin.Context) {
	result := logic.TaskLogic.GetAllTaskJob()
	rOk(ctx, "Operation successful!", result)
}

func (t *taskApi) DeleteTaskById(ctx *gin.Context) {
	err := logic.TaskLogic.DeleteTask(getQueryInt(ctx, "id"))
	errCheck(ctx, err != nil, err)
	rOk(ctx, "Operation successful!", nil)
}

func (t *taskApi) StartTask(ctx *gin.Context) {
	go logic.TaskLogic.RunTaskById(getQueryInt(ctx, "id"))
	rOk(ctx, "Operation successful!", nil)
}

func (t *taskApi) StopTask(ctx *gin.Context) {
	errCheck(ctx, logic.TaskLogic.StopTaskJob(getQueryInt(ctx, "id")) != nil, "Operation failed!")
	rOk(ctx, "Operation successful!", nil)
}

func (t *taskApi) EditTask(ctx *gin.Context, req model.Task) {
	err := logic.TaskLogic.EditTask(req)
	errCheck(ctx, err != nil, err)
	rOk(ctx, "Operation successful!", nil)
}

func (t *taskApi) EditTaskEnable(ctx *gin.Context, req model.Task) {
	err := logic.TaskLogic.EditTaskEnable(req.Id, req.Enable)
	errCheck(ctx, err != nil, err)
	rOk(ctx, "Operation successful!", nil)
}

func (t *taskApi) RunTaskByKey(ctx *gin.Context) {
	err := logic.TaskLogic.RunTaskByKey(ctx.Param("key"))
	errCheck(ctx, err != nil, err)
	rOk(ctx, "Operation successful!", nil)
}

func (t *taskApi) CreateTaskApiKey(ctx *gin.Context) {
	err := logic.TaskLogic.CreateApiKey(getQueryInt(ctx, "id"))
	errCheck(ctx, err != nil, err)
	rOk(ctx, "Operation successful!", nil)
}
