package logic

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/lzh-1625/go_process_manager/internal/app/constants"
	"github.com/lzh-1625/go_process_manager/internal/app/model"
	"github.com/lzh-1625/go_process_manager/internal/app/repository"
	"github.com/lzh-1625/go_process_manager/log"
	"github.com/lzh-1625/go_process_manager/utils"

	"github.com/robfig/cron/v3"
)

type taskLogic struct {
	taskJobMap sync.Map
}

var TaskLogic = new(taskLogic)

func (t *taskLogic) getTaskJob(id int) (*model.TaskJob, error) {
	c, ok := t.taskJobMap.Load(id)
	if !ok {
		return nil, errors.New("don't exist this task id")
	}
	return c.(*model.TaskJob), nil
}

func (t *taskLogic) InitTaskJob() {
	for _, v := range repository.TaskRepository.GetAllTask() {
		tj := &model.TaskJob{
			Task:      &v,
			StartTime: time.Now(),
		}
		if tj.Task.Cron != "" {
			c := cron.New()
			_, err := c.AddFunc(v.Cron, t.cronHandle(tj))
			if err != nil {
				log.Logger.Errorw("定时任务创建失败", "err", err, "id", v.Id)
				continue
			}
			if v.Enable {
				c.Start()
			}
			tj.Cron = c
		}
		t.taskJobMap.Store(v.Id, tj)
	}
}

func (t *taskLogic) cronHandle(data *model.TaskJob) func() {
	return func() {
		log.Logger.AddAdditionalInfo("id", data.Task.Id)
		defer log.Logger.DeleteAdditionalInfo(1)
		log.Logger.Infow("定时任务启动")
		if data.Running {
			log.Logger.Infow("任务已在运行，跳过当前任务")
			return
		}
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		data.Cancel = cancel
		t.run(ctx, data)
		log.Logger.Infow("定时任务结束")
	}
}

func (t *taskLogic) StopTaskJob(id int) error {
	taskJob, err := t.getTaskJob(id)
	if err != nil {
		return errors.New("don't exist this task id")
	}
	if taskJob.Running {
		taskJob.Cancel()
	}
	return nil
}

func (t *taskLogic) StartTaskJob(id int) error {
	taskJob, err := t.getTaskJob(id)
	if err != nil {
		return errors.New("don't exist this task id")
	}
	taskJob.Cron.Run()
	return nil
}

func (t *taskLogic) GetAllTaskJob() []model.TaskVo {
	result := repository.TaskRepository.GetAllTaskWithProcessName()
	for i, v := range result {
		task, err := t.getTaskJob(v.Id)
		if err != nil {
			continue
		}
		result[i].Id = task.Task.Id
		result[i].Running = task.Running
		result[i].Enable = task.Task.Enable
	}
	return result
}

func (t *taskLogic) DeleteTask(id int) (err error) {
	t.StopTaskJob(id)
	t.EditTaskEnable(id, false)
	t.taskJobMap.Delete(id)
	err = repository.TaskRepository.DeleteTask(id)
	if err != nil {
		return
	}
	return
}

func (t *taskLogic) CreateTask(data model.Task) error {
	tj := &model.TaskJob{
		Task:      &data,
		StartTime: time.Now(),
	}
	if data.Cron != "" {
		if _, err := cron.ParseStandard(data.Cron); err != nil { // cron表达式校验
			log.Logger.Errorw("cron解析失败", "cron", data.Cron, "err", err)
			return err
		} else {
			c := cron.New()
			c.AddFunc(data.Cron, t.cronHandle(tj))
			tj.Cron = c
		}
	}
	taskId, err := repository.TaskRepository.AddTask(data)
	if err != nil {
		return err
	}
	data.Id = taskId
	t.taskJobMap.Store(data.Id, tj)
	return nil
}

func (t *taskLogic) EditTask(data model.Task) error {
	tj, err := t.getTaskJob(data.Id)
	if err != nil {
		return fmt.Errorf("task with id %v does not exist", data.Id)
	}

	if tj.Running {
		return errors.New("can't edit task while it is running")
	}

	// 如果 Cron 已经存在，停止并清理
	if tj.Cron != nil {
		tj.Cron.Stop()
		tj.Cron = nil
	}

	// 更新任务
	tj.Task = &data

	// 如果 Cron 字段为空，直接禁用任务并返回
	if tj.Task.Cron == "" {
		tj.Task.Enable = false
		return repository.TaskRepository.EditTask(data)
	}

	// 校验 Cron 表达式
	if _, err := cron.ParseStandard(tj.Task.Cron); err != nil {
		tj.Task.Enable = false
		return fmt.Errorf("invalid cron expression: %v", err)
	}

	// 创建 Cron 调度器
	c := cron.New()
	_, err = c.AddFunc(data.Cron, t.cronHandle(tj))
	if err != nil {
		log.Logger.Errorw("failed to create cron job", "err", err, "id", data.Id)
		tj.Task.Enable = false
		return fmt.Errorf("failed to create cron job: %v", err)
	}

	// 启动 Cron 调度器
	if data.Enable {
		c.Start()
	}

	tj.Cron = c

	// 更新任务到数据库
	return repository.TaskRepository.EditTask(data)
}

func (t *taskLogic) EditTaskEnable(id int, status bool) error {
	tj, err := t.getTaskJob(id)
	if err != nil {
		return errors.New("don't exist this task id")
	}
	if tj.Cron != nil {
		if status {
			tj.Cron.Start()
		} else {
			tj.Cron.Stop()
		}
	} else if status {
		return errors.New("cron job create failed")
	}

	if err := repository.TaskRepository.EditTaskEnable(id, status); err != nil {
		return err
	}
	return nil
}

func (t *taskLogic) CreateApiKey(id int) error {
	data, err := repository.TaskRepository.GetTaskById(id)
	if err != nil {
		return err
	}
	key := utils.RandString(10)
	data.Key = &key
	repository.TaskRepository.EditTask(data)
	return nil
}

func (t *taskLogic) RunTaskByKey(key string) error {
	data, err := repository.TaskRepository.GetTaskByKey(key)
	if err != nil {
		return errors.New("don't exist key")
	}
	go t.RunTaskById(data.Id)
	return nil
}

func (t *taskLogic) RunTaskByTriggerEvent(processName string, event constants.ProcessState) {
	taskList := repository.TaskRepository.GetTriggerTask(processName, event)
	if len(taskList) == 0 {
		return
	}
	log.Logger.Infow("获取触发任务", "count", len(taskList), "prcess", processName, "触发事件", event)
	for _, v := range taskList {
		log.Logger.Infow("执行触发任务", "taskId", v.Id)
		t.RunTaskById(v.Id)
	}
}
