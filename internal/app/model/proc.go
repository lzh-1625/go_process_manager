package model

import "github.com/lzh-1625/go_process_manager/internal/app/constants"

type ProcessInfo struct {
	Name         string                 `json:"name"`
	Uuid         int                    `json:"uuid"`
	StartTime    string                 `json:"startTime"`
	User         string                 `json:"user"`
	Usage        Usage                  `json:"usage"`
	State        State                  `json:"state"`
	TermType     constants.TerminalType `json:"termType"`
	CgroupEnable bool                   `json:"cgroupEnable"`
	MemoryLimit  *float32               `json:"memoryLimit"`
	CpuLimit     *float32               `json:"cpuLimit"`
}

type Usage struct {
	Cpu  []float64 `json:"cpu"`
	Mem  []float64 `json:"mem"`
	Time []string  `json:"time"`
}

type State struct {
	State constants.ProcessState `json:"state"`
	Info  string                 `json:"info"`
}
