package service

import (
	"errors"
	"slices"
	"strings"
	"sync"

	"github.com/lzh-1625/go_process_manager/internal/app/constants"
	"github.com/lzh-1625/go_process_manager/internal/app/model"
	"github.com/lzh-1625/go_process_manager/internal/app/repository"
	"github.com/lzh-1625/go_process_manager/log"
	"github.com/lzh-1625/go_process_manager/utils"
)

type processCtlService struct {
	processMap sync.Map
}

var (
	ProcessCtlService = new(processCtlService)
)

func (p *processCtlService) AddProcess(uuid int, process *ProcessBase) {
	p.processMap.Store(uuid, process)
}

func (p *processCtlService) KillProcess(uuid int) error {
	value, ok := p.processMap.Load(uuid)
	if !ok {
		return errors.New("进程不存在")
	}
	result, ok := value.(*ProcessBase)
	if !ok {
		return errors.New("进程类型错误")
	}
	if result.State.State != 1 {
		return nil
	}
	result.State.manualStopFlag = true
	return result.Kill()
}

func (p *processCtlService) GetProcess(uuid int) (*ProcessBase, error) {
	process, ok := p.processMap.Load(uuid)
	if !ok {
		return nil, errors.New("进程获取失败")
	}
	result, ok := process.(*ProcessBase)
	if !ok {
		return nil, errors.New("进程类型错误")

	}
	return result, nil
}

func (p *processCtlService) KillAllProcess() {
	wg := sync.WaitGroup{}
	p.processMap.Range(func(key, value any) bool {
		process := value.(*ProcessBase)
		if process.State.State != 1 {
			return true
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			process.State.manualStopFlag = true
			process.Kill()
		}()
		return true
	})
	wg.Wait()
}

func (p *processCtlService) KillAllProcessByUserName(userName string) {
	stopPermissionProcess := repository.PermissionRepository.GetProcessNameByPermission(userName, constants.OPERATION_STOP)
	wg := sync.WaitGroup{}
	p.processMap.Range(func(key, value any) bool {
		process := value.(*ProcessBase)
		if !slices.Contains(stopPermissionProcess, process.Name) {
			return true
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			process.Kill()
		}()
		return true
	})
	wg.Wait()
}

func (p *processCtlService) DeleteProcess(uuid int) {
	p.processMap.Delete(uuid)
}

func (p *processCtlService) GetProcessList() []model.ProcessInfo {
	processConfiglist := repository.ProcessRepository.GetAllProcessConfig()
	return p.getProcessInfoList(processConfiglist)
}

func (p *processCtlService) GetProcessListByUser(username string) []model.ProcessInfo {
	processConfiglist := repository.ProcessRepository.GetProcessConfigByUser(username)
	return p.getProcessInfoList(processConfiglist)
}

func (p *processCtlService) getProcessInfoList(processConfiglist []model.Process) []model.ProcessInfo {
	processInfoList := []model.ProcessInfo{}
	for _, v := range processConfiglist {
		pi := model.ProcessInfo{
			Name: v.Name,
			Uuid: v.Uuid,
		}
		if value, ok := p.processMap.Load(v.Uuid); ok {
			process := value.(*ProcessBase)
			pi.State.Info = process.State.Info
			pi.State.State = process.State.State
			pi.StartTime = process.GetStartTimeFormat()
			pi.User = process.GetUserString()
			pi.Usage.Cpu = process.performanceStatus.cpu
			pi.Usage.Mem = process.performanceStatus.mem
			pi.Usage.Time = process.performanceStatus.time
			pi.TermType = process.Type()
			pi.CgroupEnable = process.cgroup.enable
			pi.CpuLimit = process.Config.cpuLimit
			pi.MemoryLimit = process.Config.memoryLimit
		}
		processInfoList = append(processInfoList, pi)
	}
	return processInfoList
}

func (p *processCtlService) ProcessStartAll() {
	p.processMap.Range(func(key, value any) bool {
		process := value.(*ProcessBase)
		err := process.Start()
		if err != nil {
			log.Logger.Errorw("进程启动失败", "name", process.Name)
		}
		return true
	})
}

func (p *processCtlService) RunPrcessById(id int) (*ProcessBase, error) {
	config := repository.ProcessRepository.GetProcessConfigById(id)
	proc, err := p.RunNewProcess(config)
	if err != nil {
		log.Logger.Warnw("初始化启动进程失败", config.Name, "name", "err", err)
		return nil, err
	}
	p.AddProcess(id, proc)
	return proc, nil
}

func (p *processCtlService) ProcessInit() {
	config := repository.ProcessRepository.GetAllProcessConfig()
	for _, v := range config {

		proc, err := p.NewProcess(v)
		if err != nil {
			log.Logger.Warnw("初始化启动进程失败", v.Name, "name", "err", err)
			continue
		}
		if v.AutoRestart {
			err := proc.Start()
			if err != nil {
				log.Logger.Warnw("初始化启动进程失败", v.Name, "name", "err", err)
				continue
			}
		}
		p.AddProcess(v.Uuid, proc)
	}
}

func (p *processCtlService) ProcesStartAllByUsername(userName string) {
	startPermissionProcess := repository.PermissionRepository.GetProcessNameByPermission(userName, constants.OPERATION_START)
	p.processMap.Range(func(key, value any) bool {
		process := value.(*ProcessBase)
		if !slices.Contains(startPermissionProcess, process.Name) {
			return true
		}
		err := process.Start()
		if err != nil {
			log.Logger.Errorw("进程启动失败", "name", process.Name)
		}
		return true
	})
}

func (p *processCtlService) UpdateProcessConfig(config model.Process) error {
	process, ok := p.processMap.Load(config.Uuid)
	if !ok {
		return errors.New("进程获取失败")
	}
	result, ok := process.(*ProcessBase)
	if !ok {
		return errors.New("进程类型错误")
	}
	if !result.Lock.TryLock() {
		return errors.New("进程当前正在被使用")
	}
	defer result.Lock.Unlock()
	result.Config.logReport = config.LogReport
	result.Config.PushIds = utils.JsonStrToStruct[[]int](config.PushIds)
	result.Config.cgroupEnable = config.CgroupEnable
	result.Config.memoryLimit = config.MemoryLimit
	result.Config.cpuLimit = config.CpuLimit
	result.Config.AutoRestart = config.AutoRestart
	result.Config.compulsoryRestart = config.CompulsoryRestart
	result.StartCommand = strings.Split(config.Cmd, " ")
	result.WorkDir = config.Cwd
	result.Name = config.Name
	return nil
}

func (p *processCtlService) RunNewProcess(config model.Process) (proc *ProcessBase, err error) {
	switch config.TermType {
	case constants.TERMINAL_STD:
		proc, err = RunNewProcessStd(config)
	case constants.TERMINAL_PTY:
		proc, err = RunNewProcessPty(config)
	default:
		err = errors.New("终端类型错误")
	}
	return
}

func (p *processCtlService) NewProcess(config model.Process) (proc *ProcessBase, err error) {
	switch config.TermType {
	case constants.TERMINAL_STD:
		proc = NewProcessStd(config)
	case constants.TERMINAL_PTY:
		proc = NewProcessPty(config)
	default:
		err = errors.New("终端类型错误")
	}
	return
}
