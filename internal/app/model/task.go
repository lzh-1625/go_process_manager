package model

import (
	"context"
	"time"

	"github.com/lzh-1625/go_process_manager/internal/app/constants"

	"github.com/robfig/cron/v3"
)

type Task struct {
	Id              int                     `gorm:"column:id;NOT NULL;primaryKey;autoIncrement;" json:"id" `
	ProcessId       int                     `gorm:"column:process_id;NOT NULL" json:"processId" `
	Condition       constants.Condition     `gorm:"column:condition;NOT NULL" json:"condition" `
	NextId          *int                    `gorm:"column:next_id;" json:"nextId" `
	Operation       constants.TaskOperation `gorm:"column:operation;NOT NULL" json:"operation" `
	TriggerEvent    *constants.ProcessState `gorm:"column:trigger_event;" json:"triggerEvent" `
	TriggerTarget   *int                    `gorm:"column:trigger_target;" json:"triggerTarget" `
	OperationTarget int                     `gorm:"column:operation_target;NOT NULL" json:"operationTarget" `
	Cron            string                  `gorm:"column:cron;" json:"cron" `
	Enable          bool                    `gorm:"column:enable;" json:"enable" `
	ApiEnable       bool                    `gorm:"column:api_enable;" json:"apiEnable" `
	Key             *string                 `gorm:"column:key;" json:"key" `
}

func (*Task) TableName() string {
	return "task"
}

type TaskJob struct {
	Cron      *cron.Cron         `json:"-"`
	Task      *Task              `json:"task"`
	Running   bool               `json:"running"`
	Cancel    context.CancelFunc `json:"-"`
	StartTime time.Time          `json:"startTime"`
	EndTime   time.Time          `json:"endTime"`
}

type TaskVo struct {
	Task
	ProcessName string    `gorm:"column:process_name;" json:"processName"`
	TargetName  string    `gorm:"column:target_name;" json:"targetName"`
	TriggerName string    `gorm:"column:trigger_name;" json:"triggerName"`
	StartTime   time.Time `json:"startTime"`
	Running     bool      `json:"running"`
}
