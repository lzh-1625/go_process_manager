package service

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/lzh-1625/go_process_manager/config"
	"github.com/lzh-1625/go_process_manager/internal/app/constants"
	"github.com/lzh-1625/go_process_manager/internal/app/middle"
	"github.com/lzh-1625/go_process_manager/internal/app/model"
	"github.com/lzh-1625/go_process_manager/log"

	"github.com/containerd/cgroups/v3"
	"github.com/containerd/cgroups/v3/cgroup1"
	"github.com/containerd/cgroups/v3/cgroup2"
	"github.com/opencontainers/runtime-spec/specs-go"
	pu "github.com/shirou/gopsutil/process"
)

type Process interface {
	ReadCache(ConnectInstance)
	Kill() error
	Write(string) error
	WriteBytes([]byte) error
	readInit()
	doOnInit()
	doOnKilled()
	Start() error
	Type() constants.TerminalType
	SetTerminalSize(int, int)
}

type ProcessBase struct {
	Process
	Name         string
	Pid          int
	cmd          *exec.Cmd
	StartCommand []string
	WorkDir      string
	Lock         sync.Mutex
	StopChan     chan struct{}
	Control      struct {
		Controller       string
		changControlTime time.Time
	}
	ws     map[string]ConnectInstance
	wsLock sync.Mutex
	Config struct {
		AutoRestart       bool
		compulsoryRestart bool
		statuPush         bool
		logReport         bool
		cgroupEnable      bool
		memoryLimit       *float32
		cpuLimit          *float32
	}
	State struct {
		startTime      time.Time
		Info           string
		State          constants.ProcessState //0 为未运行，1为运作中，2为异常状态
		stateLock      sync.Mutex
		restartTimes   int
		manualStopFlag bool
	}
	performanceStatus struct {
		cpu  []float64
		mem  []float64
		time []string
	}
	monitor struct {
		enable bool
		pu     *pu.Process
	}
	cgroup struct {
		enable bool
		delete func() error
	}
}
type ConnectInstance interface {
	Write([]byte)
	WriteString(string)
	Cancel()
}

func (p *ProcessBase) watchDog() {
	state, _ := p.cmd.Process.Wait()
	if p.cgroup.enable && p.cgroup.delete != nil {
		err := p.cgroup.delete()
		if err != nil {
			log.Logger.Errorw("cgroup删除失败", "err", err, "进程名称", p.Name)
		}
	}
	close(p.StopChan)
	p.doOnKilled()
	p.SetState(constants.PROCESS_STOP)
	if state.ExitCode() != 0 {
		log.Logger.Infow("进程停止", "进程名称", p.Name, "exitCode", state.ExitCode(), "进程类型", p.Type())
		p.push(fmt.Sprintf("进程停止,退出码 %d", state.ExitCode()))
	} else {
		log.Logger.Infow("进程正常退出", "进程名称", p.Name)
		p.push("进程正常退出")
	}
	if !p.Config.AutoRestart || p.State.manualStopFlag { // 不重启或手动关闭
		return
	}
	if p.Config.compulsoryRestart { // 强制重启
		p.Start()
		return
	}
	if state.ExitCode() == 0 { // 正常退出
		return
	}
	if p.State.restartTimes < config.CF.ProcessRestartsLimit { // 重启次数未达限制
		p.Start()
		p.State.restartTimes++
		return
	}
	log.Logger.Warnw("重启次数达到上限", "name", p.Name, "limit", config.CF.ProcessRestartsLimit)
	p.SetState(constants.PROCESS_WARNNING)
	p.State.Info = "重启次数异常"
	p.push("进程重启次数达到上限")
}

func (p *ProcessBase) pInit() {
	log.Logger.Infow("创建进程成功")
	p.StopChan = make(chan struct{})
	p.State.manualStopFlag = false
	p.Pid = p.cmd.Process.Pid
	p.State.startTime = time.Now()
	p.ws = make(map[string]ConnectInstance)
	p.doOnInit()
	p.InitPerformanceStatus()
	p.initPsutil()
	p.initCgroup()
	go p.watchDog()
	go p.readInit()
	go p.monitorHanler()
	middle.ProcessWaitCond.Trigger()
}

// fn 函数执行成功的情况下对state赋值
func (p *ProcessBase) SetState(state constants.ProcessState, fn ...func() bool) bool {
	p.State.stateLock.Lock()
	defer p.State.stateLock.Unlock()
	for _, v := range fn {
		if !v() {
			return false
		}
	}
	p.State.State = state
	middle.ProcessWaitCond.Trigger()
	go TaskService.RunTaskByTriggerEvent(p.Name, state)
	return true
}

func (p *ProcessBase) GetUserString() string {
	return strings.Join(p.GetUserList(), ";")
}

func (p *ProcessBase) GetUserList() []string {
	userList := make([]string, 0, len(p.ws))
	for i := range p.ws {
		userList = append(userList, i)
	}
	return userList
}

func (p *ProcessBase) HasWsConn(userName string) bool {
	return p.ws[userName] != nil
}

func (p *ProcessBase) AddConn(user string, c ConnectInstance) {
	if p.ws[user] != nil {
		log.Logger.Error("已存在连接")
		return
	}
	p.wsLock.Lock()
	defer p.wsLock.Unlock()
	p.ws[user] = c
}

func (p *ProcessBase) DeleteConn(user string) {
	p.wsLock.Lock()
	defer p.wsLock.Unlock()
	delete(p.ws, user)
}

func (p *ProcessBase) logReportHandler(log string) {
	if p.Config.logReport && len([]rune(log)) > config.CF.LogMinLenth {
		Loghandler.AddLog(model.ProcessLog{
			Log:   log,
			Using: p.GetUserString(),
			Name:  p.Name,
			Time:  time.Now().UnixMilli(),
		})
	}
}

func (p *ProcessBase) GetStartTimeFormat() string {
	return p.State.startTime.Format(time.DateTime)
}

func (p *ProcessBase) ProcessControl(name string) {
	p.Control.changControlTime = time.Now()
	p.Control.Controller = name
	for _, ws := range p.ws {
		ws.Cancel()
	}
}

// 没人在使用或控制时间过期
func (p *ProcessBase) VerifyControl() bool {
	return p.Control.Controller == "" || p.Control.changControlTime.Unix() < time.Now().Unix()-config.CF.ProcessExpireTime
}

func (p *ProcessBase) setProcessConfig(pconfig model.Process) {
	p.Config.AutoRestart = pconfig.AutoRestart
	p.Config.logReport = pconfig.LogReport
	p.Config.statuPush = pconfig.Push
	p.Config.compulsoryRestart = pconfig.CompulsoryRestart
	p.Config.cgroupEnable = pconfig.CgroupEnable
	p.Config.memoryLimit = pconfig.MemoryLimit
	p.Config.cpuLimit = pconfig.CpuLimit
}

func (p *ProcessBase) ResetRestartTimes() {
	p.State.restartTimes = 0
}

func (p *ProcessBase) push(message string) {
	if p.Config.statuPush {
		messagePlaceholders := map[string]string{
			"{$name}":    p.Name,
			"{$user}":    p.GetUserString(),
			"{$message}": message,
			"{$status}":  strconv.Itoa(int(p.State.State)),
		}
		PushService.Push(messagePlaceholders)
	}
}

func (p *ProcessBase) InitPerformanceStatus() {
	p.performanceStatus.cpu = make([]float64, config.CF.PerformanceInfoListLength)
	p.performanceStatus.mem = make([]float64, config.CF.PerformanceInfoListLength)
	p.performanceStatus.time = make([]string, config.CF.PerformanceInfoListLength)
}

func (p *ProcessBase) AddCpuUsage(usage float64) {
	p.performanceStatus.cpu = append(p.performanceStatus.cpu[1:], usage)
}

func (p *ProcessBase) AddMemUsage(usage float64) {
	p.performanceStatus.mem = append(p.performanceStatus.mem[1:], usage)
}

func (p *ProcessBase) AddRecordTime() {
	p.performanceStatus.time = append(p.performanceStatus.time[1:], time.Now().Format(time.DateTime))
}

func (p *ProcessBase) initCgroup() {
	if !p.Config.cgroupEnable {
		log.Logger.Debugw("不启用cgroup")
		return
	}
	switch cgroups.Mode() {
	case cgroups.Unavailable:
		log.Logger.Warnw("当前系统不支持cgroup")
	case cgroups.Legacy, cgroups.Hybrid:
		log.Logger.Debugw("启用cgroupv1")
		p.initCgroupV1()
	case cgroups.Unified:
		log.Logger.Debugw("启用cgroupv2")
		p.initCgroupV2()
	}
}

func (p *ProcessBase) initCgroupV1() {
	resources := &specs.LinuxResources{}
	if p.Config.cpuLimit != nil {
		period := uint64(config.CF.CgroupPeriod)
		quota := int64(float32(config.CF.CgroupPeriod) * *p.Config.cpuLimit * 0.01)
		cpuResources := &specs.LinuxCPU{
			Period: &period,
			Quota:  &quota,
		}
		resources.CPU = cpuResources
	}
	if p.Config.memoryLimit != nil {
		limit := int64(*p.Config.memoryLimit * 1024 * 1024)
		memResources := &specs.LinuxMemory{
			Limit: &limit,
		}
		if config.CF.CgroupSwapLimit {
			memResources.Swap = &limit
		}
		resources.Memory = memResources
	}
	control, err := cgroup1.New(cgroup1.StaticPath("/"+p.Name), resources)
	if err != nil {
		log.Logger.Errorw("启用cgroup失败", "err", err, "name", p.Name)
		return
	}
	control.AddProc(uint64(p.cmd.Process.Pid))
	p.cgroup.delete = control.Delete
	p.cgroup.enable = true
}

func (p *ProcessBase) initCgroupV2() {
	resources := &cgroup2.Resources{}
	if p.Config.cpuLimit != nil {
		period := uint64(config.CF.CgroupPeriod)
		quota := int64(float32(config.CF.CgroupPeriod) * *p.Config.cpuLimit * 0.01)
		resources.CPU = &cgroup2.CPU{
			Max: cgroup2.NewCPUMax(&quota, &period),
		}
	}
	if p.Config.memoryLimit != nil {
		limit := int64(*p.Config.memoryLimit * 1024 * 1024)
		memResources := &cgroup2.Memory{
			Max: &limit,
		}
		if config.CF.CgroupSwapLimit {
			memResources.Swap = &limit
		}
		resources.Memory = memResources
	}
	control, err := cgroup2.NewSystemd("/", p.Name+".slice", -1, resources)
	if err != nil {
		log.Logger.Errorw("启用cgroup失败", "err", err, "name", p.Name)
		return
	}
	control.AddProc(uint64(p.cmd.Process.Pid))
	p.cgroup.delete = control.DeleteSystemd
	p.cgroup.enable = true
}

func (p *ProcessBase) monitorHanler() {
	if !p.monitor.enable {
		return
	}
	log.Logger.AddAdditionalInfo("name", p.Name)
	log.Logger.AddAdditionalInfo("pid", p.Pid)
	defer log.Logger.Infow("性能监控结束")
	ticker := time.NewTicker(time.Second * time.Duration(config.CF.PerformanceInfoInterval))
	defer ticker.Stop()
	for {
		if p.State.State != 1 {
			log.Logger.Debugw("进程未在运行", "state", p.State.State)
			return
		}
		cpuPercent, err := p.monitor.pu.CPUPercent()
		if err != nil {
			log.Logger.Errorw("CPU使用率获取失败", "err", err)
			return
		}
		memInfo, err := p.monitor.pu.MemoryInfo()
		if err != nil {
			log.Logger.Errorw("内存使用率获取失败", "err", err)
			return
		}
		p.AddRecordTime()
		p.AddCpuUsage(cpuPercent)
		p.AddMemUsage(float64(memInfo.RSS) / 1024.0)
		// log.Logger.Debugw("进程资源使用率获取成功", "cpu", cpuPercent, "mem", memInfo.RSS)
		select {
		case <-ticker.C:
		case <-p.StopChan:
			return
		}
	}
}

func (p *ProcessBase) initPsutil() {
	pup, err := pu.NewProcess(int32(p.Pid))
	if err != nil {
		p.monitor.enable = false
		log.Logger.Debug("pu进程获取失败")
	} else {
		p.monitor.enable = true
		log.Logger.Debug("pu进程获取成功")
		p.monitor.pu = pup
	}
}

func (p *ProcessBase) Kill() error {
	p.cmd.Process.Signal(syscall.SIGINT)
	select {
	case <-p.StopChan:
		{
			return nil
		}
	case <-time.After(time.Second * time.Duration(config.CF.KillWaitTime)):
		{
			log.Logger.Debugw("进程kill超时,强制停止进程", "name", p.Name)
			return p.cmd.Process.Kill()
		}
	}
}
