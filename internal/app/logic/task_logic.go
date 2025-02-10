package logic

import (
	"context"
	"errors"
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

var TaskLogic taskLogic

func (t *taskLogic) InitTaskJob() {
	for _, v := range repository.TaskRepository.GetAllTask() {
		tj := &model.TaskJob{
			Task:      &v,
			StartTime: time.Now(),
		}
		if tj.Task.Cron != nil {
			c := cron.New()
			_, err := c.AddFunc(*v.Cron, t.cronHandle(tj))
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
	c, ok := t.taskJobMap.Load(id)
	if !ok {
		return errors.New("id不存在")
	}
	taskJob := c.(*model.TaskJob)
	if taskJob.Running {
		taskJob.Cancel()
	}
	return nil
}

func (t *taskLogic) StartTaskJob(id int) error {
	c, ok := t.taskJobMap.Load(id)
	if !ok {
		return errors.New("id不存在")
	}
	TaskJob := c.(*model.TaskJob)
	TaskJob.Cron.Run()
	return nil
}

func (t *taskLogic) GetAllTaskJob() []model.TaskVo {
	result := repository.TaskRepository.GetAllTaskWithProcessName()
	for i, v := range result {
		item, ok := t.taskJobMap.Load(v.Id)
		if !ok {
			continue
		}
		task := item.(*model.TaskJob)
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
	if data.Cron != nil {
		if _, err := cron.ParseStandard(*data.Cron); err != nil { // cron表达式校验
			log.Logger.Errorw("cron解析失败", "cron", *data.Cron, "err", err)
			return err
		} else {
			c := cron.New()
			c.AddFunc(*data.Cron, t.cronHandle(tj))
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
	if data.Cron != nil {
		if _, err := cron.ParseStandard(*data.Cron); err != nil {
			return err
		}
	}
	v, ok := t.taskJobMap.Load(data.Id)
	if !ok {
		return errors.New("don't exist this task id")
	}
	tj := v.(*model.TaskJob)
	tj.Task = &data
	return repository.TaskRepository.EditTask(data)
}

func (t *taskLogic) EditTaskEnable(id int, status bool) error {
	v, ok := t.taskJobMap.Load(id)
	if !ok {
		return errors.New("don't exist this task id")
	}
	tj := v.(*model.TaskJob)
	tj.Task.Enable = status
	repository.TaskRepository.EditTaskEnable(id, status)
	if tj.Cron != nil {
		if status {
			tj.Cron.Start()
		} else {
			tj.Cron.Stop()
		}
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
