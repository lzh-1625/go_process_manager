package model

import "github.com/lzh-1625/go_process_manager/internal/app/constants"

type Process struct {
	Uuid              int                    `gorm:"primaryKey;autoIncrement;column:uuid" json:"uuid"`
	Name              string                 `gorm:"column:name" json:"name"`
	Cmd               string                 `gorm:"column:args" json:"cmd"`
	Cwd               string                 `gorm:"column:cwd" json:"cwd"`
	AutoRestart       bool                   `gorm:"column:auto_restart" json:"autoRestart"`
	CompulsoryRestart bool                   `gorm:"column:compulsory_restart" json:"compulsoryRestart"`
	PushIds           string                 `gorm:"column:push_ids" json:"pushIds"`
	LogReport         bool                   `gorm:"column:log_report" json:"logReport"`
	TermType          constants.TerminalType `gorm:"column:term_type" json:"termType"`
	CgroupEnable      bool                   `gorm:"column:cgroup_enable" json:"cgroupEnable"`
	MemoryLimit       *float32               `gorm:"column:memory_limit" json:"memoryLimit"`
	CpuLimit          *float32               `gorm:"column:cpu_limit" json:"cpuLimit"`
}

func (*Process) TableName() string {
	return "process"
}
