package process

import (
	"errors"
	"msm/config"
	"msm/log"
	"msm/model"
	loghandler "msm/service/log"
	"msm/service/push"
	"os/exec"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
	pu "github.com/shirou/gopsutil/process"
)

type Process interface {
	ReadCache(*websocket.Conn)
	GetName() string
	SetName(string)
	GetTermType() string
	SetTermType(string)
	SetIsUsing(bool)
	GetWhoUsing() string
	SetWhoUsing(string)
	SetStartCommand([]string)
	GetControlController() string
	SetControlController(string)
	ChangControlChan() chan int
	StopChan() chan struct{}
	SetConfigLogReport(bool)
	SetConfigStatuPush(bool)
	SetConfigAutoRestart(bool)
	GetStateInfo() string
	GetStateState() uint8
	Kill() error
	SetWsConn(*websocket.Conn)
	Write(string) error
	WriteBytes([]byte) error
	GetStartTimeFormat() string
	VerifyControl() bool
	ResetRestartTimes()
	InitPerformanceStatus()
	ProcessControl(string)
	AddCpuUsage(float64)
	AddMemUsage(float64)
	AddRecordTime()
	GetTimeRecord() []string
	GetMemUsage() []float64
	GetCpuUsage() []float64
	monitorHanler()
	initPsutil()
	SetAutoRestart(bool)
	TryLock() bool
	Unlock()
	ReStart()
}

type ProcessBase struct {
	Name         string
	termType     string
	Pid          int
	cmd          *exec.Cmd
	IsUsing      atomic.Bool
	StartCommand []string
	Lock         sync.Mutex
	WhoUsing     string
	stopChan     chan struct{}
	Control      struct {
		Controller       string
		changControlChan chan int
		changControlTime time.Time
	}
	ws struct {
		wsConnect *websocket.Conn
		wsMux     sync.RWMutex
	}
	Config struct {
		AutoRestart bool
		statuPush   bool
		logReport   bool
	}
	State struct {
		startTime    time.Time
		Info         string
		State        uint8 //0 为未运行，1为运作中，2为异常状态
		restartTimes int
	}
	performanceStatus struct {
		cpu  []float64
		mem  []float64
		time []string
	}
	monitor struct {
		enable      bool
		ProcessBase *pu.Process
	}
}

func (p *ProcessBase) GetTermType() string {
	return p.termType
}

func (p *ProcessBase) SetTermType(s string) {
	p.termType = s
}

func (p *ProcessBase) GetStateInfo() string {
	return p.State.Info
}

func (p *ProcessBase) GetStateState() uint8 {
	return p.State.State
}

func (p *ProcessBase) SetAutoRestart(data bool) {
	p.Config.AutoRestart = data
}

func (p *ProcessBase) GetWhoUsing() string {
	return p.WhoUsing
}

func (p *ProcessBase) GetControlController() string {
	return p.Control.Controller
}

func (p *ProcessBase) SetControlController(c string) {
	p.Control.Controller = c
}

func (p *ProcessBase) SetWsConn(ws *websocket.Conn) {
	p.ws.wsConnect = ws
}

func (p *ProcessBase) logReportHandler(log string) {
	if config.CF.EsEnable && p.Config.logReport && len([]rune(log)) > config.CF.LogMinLenth {
		loghandler.Loghandler.AddLog(model.Eslog{
			Log:   log,
			Using: p.WhoUsing,
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
	if p.State.State == 1 && p.IsUsing.Load() {
		p.Control.changControlChan <- 0
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
}

func (p *ProcessBase) ResetRestartTimes() {
	p.State.restartTimes = 0
}

func (p *ProcessBase) push(message string) {
	if p.Config.statuPush {
		messagePlaceholders := map[string]string{
			"{$name}":    p.Name,
			"{$user}":    p.WhoUsing,
			"{$message}": message,
			"{$status}":  strconv.Itoa(int(p.State.State)),
		}
		push.PushService.Push(messagePlaceholders)
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
func (p *ProcessBase) GetCpuUsage() []float64 {
	return p.performanceStatus.cpu
}

func (p *ProcessBase) GetMemUsage() []float64 {
	return p.performanceStatus.mem
}

func (p *ProcessBase) GetTimeRecord() []string {
	return p.performanceStatus.time
}

func (p *ProcessBase) monitorHanler() {
	defer log.Logger.Infow("性能监控结束", "name", p.Name, "pid", p.Pid)
	for {
		if !p.monitor.enable {
			return
		}
		select {
		case <-time.After(time.Minute * time.Duration(config.CF.PerformanceInfoInterval)):
			if p.State.State != 1 {
				log.Logger.Debugw("进程状态异常，跳过监控数据获取", "name", p.Name)
				p.AddCpuUsage(0)
				p.AddMemUsage(0)
				p.AddRecordTime()
				continue
			}
			ProcessBase := p.monitor.ProcessBase
			cpuPercent, err := ProcessBase.CPUPercent()
			if err != nil {
				log.Logger.Errorw("CPU使用率获取失败", "err", err)
				return
			}
			memInfo, err := ProcessBase.MemoryInfo()
			if err != nil {
				log.Logger.Errorw("内存使用率获取失败", "err", err)
				return
			}
			p.AddRecordTime()
			p.AddCpuUsage(cpuPercent)
			p.AddMemUsage(float64(memInfo.RSS / 1000))
			log.Logger.Debugw("进程资源使用率获取成功", "pid", p.Pid, "name", p.Name, "cpu", cpuPercent, "mem", memInfo.RSS)
		case <-p.stopChan:
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
		p.monitor.ProcessBase = pup
	}
}

func (p *ProcessBase) SetConfigLogReport(b bool) {
	p.Config.logReport = b
}

func (p *ProcessBase) SetConfigAutoRestart(b bool) {
	p.Config.AutoRestart = b
}

func (p *ProcessBase) SetConfigStatuPush(b bool) {
	p.Config.statuPush = b
}

func (p *ProcessBase) SetName(s string) {
	p.Name = s
}

func (p *ProcessBase) SetStartCommand(cmd []string) {
	p.StartCommand = cmd
}

func (p *ProcessBase) ChangControlChan() chan int {
	return p.Control.changControlChan
}

func (p *ProcessBase) SetIsUsing(b bool) {
	p.IsUsing.Store(b)
}

func (p *ProcessBase) GetName() string {
	return p.Name
}

func (p *ProcessBase) SetWhoUsing(s string) {
	p.WhoUsing = s
}

func (p *ProcessBase) StopChan() chan struct{} {
	return p.stopChan
}

func (p *ProcessBase) TryLock() bool {
	return p.Lock.TryLock()
}

func (p *ProcessBase) Unlock() {
	p.Lock.Unlock()
}

func RunNewProcess(config model.Process) (proc Process, err error) {
	switch config.TermType {
	case "std":
		proc, err = RunNewProcessStd(config)
	case "pty":
		proc, err = RunNewProcessPty(config)
	default:
		err = errors.New("终端类型错误")
	}
	return
}
