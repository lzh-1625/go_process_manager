package process

import (
	"bufio"
	"fmt"
	"io"
	"msm/config"
	"msm/log"
	"msm/model"
	"os/exec"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

type ProcessStd struct {
	ProcessBase
	cacheLine []string
	stdin     io.WriteCloser
	stdout    *bufio.Scanner
}

func (p *ProcessStd) Kill() error {
	return p.cmd.Process.Kill()
}

func (p *ProcessStd) watchDog() {
	state, _ := p.cmd.Process.Wait()
	close(p.stopChan)
	p.State.State = 0
	if state.ExitCode() != 0 {
		log.Logger.Infow("进程停止", "进程名称", p.Name, "exitCode", state.ExitCode(), "进程类型", "std")
		p.push(fmt.Sprintf("进程停止,退出码 %d", state.ExitCode()))
		if p.Config.AutoRestart {
			p.ReStart()
		}
	} else {
		log.Logger.Infow("进程正常退出", "进程名称", p.Name)
		p.push("进程正常退出")
	}
}

func (p *ProcessStd) WriteBytes(input []byte) (err error) {
	p.logReportHandler(config.CF.ProcessInputPrefix + string(input))
	_, err = p.stdin.Write(append(input, '\n'))
	return
}

func (p *ProcessStd) Write(input string) (err error) {
	p.logReportHandler(config.CF.ProcessInputPrefix + input)
	_, err = p.stdin.Write([]byte(input + "\n"))
	return
}

func (p *ProcessStd) ReStart() {
	if p.State.restartTimes > config.CF.ProcessRestartsLimit {
		log.Logger.Warnw("重启次数达到上限", "name", p.Name, "limit", config.CF.ProcessRestartsLimit)
		p.State.State = 2
		p.State.Info = "重启次数异常"
		p.push("进程重启次数达到上限")
		return
	}
	cmd := exec.Command(p.StartCommand[0], p.StartCommand[1:]...) // 替换为你要执行的命令及参数
	cmd.Dir = p.cmd.Dir
	out, err := cmd.StdoutPipe()
	if err != nil {
		log.Logger.Errorw("重启失败，输出管道获取失败", "err", err)
		p.Config.AutoRestart = false
		return
	}
	p.stdout = bufio.NewScanner(out)
	p.stdin, err = cmd.StdinPipe()
	if err != nil {
		log.Logger.Errorw("重启失败，输入管道获取失败", "err", err)
		p.Config.AutoRestart = false
		return
	}
	err = cmd.Start()
	if err != nil {
		log.Logger.Errorw("重启失败，进程启动出错:", "err", err)
		p.Config.AutoRestart = false
		return
	}
	p.State.restartTimes++
	log.Logger.Infow("进程启动成功", "进程名称", p.Name, "重启次数", p.State.restartTimes)
	p.cmd = cmd
	p.pInit()
	p.push("进程启动成功")

}

func (p *ProcessStd) pInit() {
	log.Logger.Infow("创建进程成功")
	p.Control.changControlChan = make(chan int)
	p.stopChan = make(chan struct{})
	p.State.State = 1
	p.Pid = p.cmd.Process.Pid
	p.State.startTime = time.Now()
	p.cacheLine = make([]string, config.CF.ProcessMsgCacheLinesLimit)
	p.InitPerformanceStatus()
	p.initPsutil()
	go p.watchDog()
	go p.readInit()
	go p.monitorHanler()
}

func (p *ProcessStd) ReadCache(ws *websocket.Conn) {
	for _, line := range p.cacheLine {
		ws.WriteMessage(websocket.BinaryMessage, []byte(line))
	}
}

func (p *ProcessStd) readInit() {
	var output string
	log.Logger.Debugw("stdout读取线程已启动", "进程名", p.Name, "使用者", p.WhoUsing)
	for {
		select {
		case <-p.stopChan:
			{
				p.IsUsing.Store(false)
				p.WhoUsing = ""
				log.Logger.Debugw("stdout读取线程已退出", "进程名", p.Name, "使用者", p.WhoUsing)
				return
			}
		default:
			{
				output = p.Read()
				if p.IsUsing.Load() && output != "" {
					p.ws.wsMux.Lock()
					p.ws.wsConnect.WriteMessage(websocket.BinaryMessage, []byte(output))
					p.ws.wsMux.Unlock()
				}
			}
		}
	}
}
func (p *ProcessStd) Read() string {
	if p.stdout.Scan() {
		output := p.stdout.Text()
		p.logReportHandler(output)
		p.cacheLine = p.cacheLine[1:]
		p.cacheLine = append(p.cacheLine, output)
		return output
	}
	return ""
}

func RunNewProcessStd(pconfig model.Process) (*ProcessStd, error) {
	args := strings.Split(pconfig.Cmd, " ")
	cmd := exec.Command(args[0], args[1:]...) // 替换为你要执行的命令及参数

	processStd := ProcessStd{
		ProcessBase: ProcessBase{
			Name:         pconfig.Name,
			StartCommand: args,
		},
	}
	cmd.Dir = pconfig.Cwd
	out, err := cmd.StdoutPipe()
	if err != nil {
		log.Logger.Errorw("输出管道获取失败", "err", err)
		return nil, err
	}
	processStd.stdout = bufio.NewScanner(out)
	processStd.stdin, err = cmd.StdinPipe()
	if err != nil {
		log.Logger.Errorw("输入管道获取失败", "err", err)
		return nil, err
	}
	err = cmd.Start()
	if err != nil || cmd.Process == nil {
		log.Logger.Error("进程启动出错:", err)
		return nil, err
	}
	log.Logger.Infow("创建进程成功", "config", pconfig)
	processStd.cmd = cmd
	processStd.SetTermType("std")
	processStd.pInit()
	processStd.setProcessConfig(pconfig)
	return &processStd, nil
}
