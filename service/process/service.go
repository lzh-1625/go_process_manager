package process

import (
	"errors"
	"msm/dao"
	"msm/log"
	"msm/model"
	"strings"
	"sync"
)

type processCtlService struct{}

var processMap sync.Map = sync.Map{}
var ProcessCtlService = new(processCtlService)

func (p *processCtlService) AddProcess(uuid int, prcess Process) {
	processMap.Store(uuid, prcess)
	// processMap.Store("111", prcess)
	// return "111"
}

func (p *processCtlService) KillProcess(uuid int) error {
	value, ok := processMap.Load(uuid)
	if !ok {
		return errors.New("进程不存在")
	}
	result, ok := value.(Process)
	if !ok {
		return errors.New("进程类型错误")
	}
	result.SetAutoRestart(false)
	return result.Kill()
}

func (p *processCtlService) GetProcess(uuid int) (Process, error) {
	process, ok := processMap.Load(uuid)
	if !ok {
		return nil, errors.New("进程获取失败")

	}
	result, ok := process.(Process)
	if !ok {
		return nil, errors.New("进程类型错误")

	}
	return result, nil
}

func (p *processCtlService) KillAllProcess() {
	processMap.Range(func(key, value any) bool {
		value.(Process).Kill()
		return true
	})
}

func (p *processCtlService) DeleteProcess(uuid int) {
	processMap.Delete(uuid)
}

func (p *processCtlService) GetProcessList() []model.ProcessInfo {
	processConfiglist := dao.ProcessDao.GetAllProcessConfig()
	return p.getProcessInfoList(processConfiglist)
}

func (p *processCtlService) GetProcessListByUser(username string) []model.ProcessInfo {
	processConfiglist := dao.ProcessDao.GetProcessConfigByUser(username)
	return p.getProcessInfoList(processConfiglist)
}

func (p *processCtlService) getProcessInfoList(processConfiglist []model.Process) []model.ProcessInfo {
	processInfoList := []model.ProcessInfo{}
	for _, v := range processConfiglist {
		pi := model.ProcessInfo{
			Name: v.Name,
			Uuid: v.Uuid,
		}
		if value, ok := processMap.Load(v.Uuid); ok {
			process := value.(Process)
			pi.State.Info = process.GetStateInfo()
			pi.State.State = process.GetStateState()
			pi.StartTime = process.GetStartTimeFormat()
			pi.User = process.GetWhoUsing()
			pi.Usage.Cpu = process.GetCpuUsage()
			pi.Usage.Mem = process.GetMemUsage()
			pi.Usage.Time = process.GetTimeRecord()
			pi.TermType = process.GetTermType()
		}
		processInfoList = append(processInfoList, pi)
	}
	return processInfoList
}

func (p *processCtlService) ProcessInit() {
	config := dao.ProcessDao.GetAllProcessConfig()
	for _, v := range config {
		if !v.AutoRestart {
			continue
		}
		proc, err := RunNewProcess(v)
		if err != nil {
			log.Logger.Warnw("初始化启动进程失败", v.Name, "name", "err", err)
			continue
		}
		p.AddProcess(v.Uuid, proc)
	}
}

func (p *processCtlService) UpdateProcessConfig(config model.Process) error {
	process, ok := processMap.Load(config.Uuid)
	if !ok {
		return errors.New("进程获取失败")
	}
	result, ok := process.(Process)
	if !ok {
		return errors.New("进程类型错误")
	}
	result.SetConfigLogReport(config.LogReport)
	result.SetConfigStatuPush(config.Push)
	result.SetConfigAutoRestart(config.AutoRestart)
	result.SetStartCommand(strings.Split(config.Cmd, " "))
	result.SetName(config.Name)
	return nil
}
