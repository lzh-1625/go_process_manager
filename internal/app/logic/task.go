package logic

import (
	"context"
	"errors"
	"time"

	"github.com/lzh-1625/go_process_manager/config"
	"github.com/lzh-1625/go_process_manager/internal/app/constants"
	"github.com/lzh-1625/go_process_manager/internal/app/middle"
	"github.com/lzh-1625/go_process_manager/internal/app/model"
	"github.com/lzh-1625/go_process_manager/log"
)

func (t *taskLogic) RunTaskById(id int) error {
	v, ok := t.taskJobMap.Load(id)
	if !ok {
		return errors.New("don't exist task id")
	}
	task := v.(*model.TaskJob)
	if task.Running {
		return errors.New("task is running")
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	task.Cancel = cancel
	t.run(ctx, task)
	return nil
}

func (t *taskLogic) run(ctx context.Context, data *model.TaskJob) {
	data.Running = true
	middle.TaskWaitCond.Trigger()
	defer func() {
		data.Running = false
		middle.TaskWaitCond.Trigger()
	}()
	log.Logger.AddAdditionalInfo("taskId", data.Task.Id)
	defer log.Logger.DeleteAdditionalInfo(1)
	var ok bool
	// 判断条件是否满足
	if data.Task.Condition == constants.PASS {
		ok = true
	} else {
		proc, err := ProcessCtlLogic.GetProcess(data.Task.OperationTarget)
		if err != nil {
			return
		}
		ok = conditionHandle[data.Task.Condition](data.Task, proc)
	}
	log.Logger.Debugw("任务条件判断", "pass", ok)
	if !ok {
		return
	}

	proc, err := ProcessCtlLogic.GetProcess(data.Task.OperationTarget)
	if err != nil {
		log.Logger.Debugw("不存在该进程，结束任务")
		return
	}

	// 执行操作
	log.Logger.Infow("任务开始执行")
	if !OperationHandle[data.Task.Operation](data.Task, proc) {
		log.Logger.Errorw("任务执行失败")
		return
	}
	log.Logger.Infow("任务执行成功", "target", data.Task.OperationTarget)

	if data.Task.NextId != nil {
		v, ok := t.taskJobMap.Load(*data.Task.NextId)
		nextTask := v.(*model.TaskJob)
		if !ok {
			log.Logger.Errorw("无法获取到下一个节点,结束任务", "nextId", data.Task.NextId)
			return
		}
		select {
		case <-ctx.Done():
			log.Logger.Infow("任务流被手动结束")
		default:
			log.Logger.Debugw("执行下一个节点", "nextId", *data.Task.NextId)
			if nextTask.Running {
				log.Logger.Errorw("下一个节点已在运行，结束任务", "nextId", data.Task.NextId)
				return
			}
			t.run(ctx, nextTask)
		}
	} else {
		log.Logger.Infow("任务流结束")
	}
}

type conditionFunc func(data *model.Task, proc *ProcessBase) bool

var conditionHandle = map[constants.Condition]conditionFunc{
	constants.RUNNING: func(data *model.Task, proc *ProcessBase) bool {
		return proc.State.State == 1
	},
	constants.NOT_RUNNING: func(data *model.Task, proc *ProcessBase) bool {
		return proc.State.State != 1
	},
	constants.EXCEPTION: func(data *model.Task, proc *ProcessBase) bool {
		return proc.State.State == 2
	},
}

// 执行操作，返回结果是否成功
type operationFunc func(data *model.Task, proc *ProcessBase) bool

var OperationHandle = map[constants.TaskOperation]operationFunc{
	constants.TASK_START: func(data *model.Task, proc *ProcessBase) bool {
		if proc.State.State == 1 {
			log.Logger.Debugw("进程已在运行")
			return false
		}
		return proc.Start() == nil
	},

	constants.TASK_START_WAIT_DONE: func(data *model.Task, proc *ProcessBase) bool {
		if proc.State.State == 1 {
			log.Logger.Debugw("进程已在运行")
			return false
		}
		if err := proc.Start(); err != nil {
			log.Logger.Debugw("进程启动失败")
			return false
		}
		select {
		case <-proc.StopChan:
			log.Logger.Debugw("进程停止，任务完成")
			return true
		case <-time.After(time.Second * time.Duration(config.CF.TaskTimeout)):
			log.Logger.Errorw("任务超时")
			return false
		}
	},

	constants.TASK_STOP: func(data *model.Task, proc *ProcessBase) bool {
		if proc.State.State != 1 {
			log.Logger.Debugw("进程未在运行")
			return false
		}
		log.Logger.Debugw("异步停止任务")
		go proc.Kill()
		return true
	},

	constants.TASK_STOP_WAIT_DONE: func(data *model.Task, proc *ProcessBase) bool {
		if proc.State.State != 1 {
			log.Logger.Debugw("进程未在运行")
			return false
		}
		log.Logger.Debugw("停止任务并等待结束")
		return proc.Kill() == nil
	},
}
