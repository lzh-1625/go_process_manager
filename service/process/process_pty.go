package process

import (
	"bytes"
	"fmt"
	"msm/config"
	"msm/log"
	"msm/model"
	"msm/utils"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/creack/pty"
	"github.com/gorilla/websocket"
)

type ProcessPty struct {
	ProcessBase
	cacheBytesBuf *bytes.Buffer
	pty           *os.File
}

func (p *ProcessPty) Kill() error {
	if err := p.cmd.Process.Kill(); err != nil {
		log.Logger.Errorw("进程杀死失败", "err", err, "state", p.State.State)
		return err
	}
	return p.pty.Close()
}

func (p *ProcessPty) watchDog() {
	state, _ := p.cmd.Process.Wait()
	close(p.stopChan)
	p.State.State = 0
	p.pty.Close()
	if state.ExitCode() != 0 {
		log.Logger.Infow("进程停止", "进程名称", p.Name, "exitCode", state.ExitCode(), "进程类型", "pty")
		p.push(fmt.Sprintf("进程停止,退出码 %d", state.ExitCode()))
		if p.Config.AutoRestart {
			p.ReStart()
		}
	} else {
		log.Logger.Infow("进程正常退出", "进程名称", p.Name)
		p.push("进程正常退出")
	}
}

func (p *ProcessPty) ReStart() {
	if p.State.restartTimes > config.CF.ProcessRestartsLimit {
		log.Logger.Warnw("重启次数达到上限", "name", p.Name, "limit", config.CF.ProcessRestartsLimit)
		p.State.State = 2
		p.State.Info = "重启次数异常"
		p.push("进程重启次数达到上限")
		return
	}
	cmd := exec.Command(p.StartCommand[0], p.StartCommand[1:]...)
	cmd.Dir = p.cmd.Dir
	pf, err := pty.Start(cmd)
	if err != nil || p.cmd.Process == nil {
		log.Logger.Error("进程启动出错:", err)
		return
	}
	pty.Setsize(pf, &pty.Winsize{
		Rows: 100,
		Cols: 100,
	})
	p.pty = pf
	p.State.restartTimes++
	log.Logger.Infow("进程启动成功", "进程名称", p.Name, "重启次数", p.State.restartTimes)
	p.cmd = cmd
	p.pInit()
	p.push("进程启动成功")
}

func (p *ProcessPty) WriteBytes(input []byte) (err error) {
	p.logReportHandler(config.CF.ProcessInputPrefix + string(input))
	_, err = p.pty.Write(input)
	return
}

func (p *ProcessPty) Write(input string) (err error) {
	p.logReportHandler(config.CF.ProcessInputPrefix + input)
	_, err = p.pty.Write([]byte(input))
	return
}

func (p *ProcessPty) readInit() {
	log.Logger.Debugw("stdout读取线程已启动", "进程名", p.Name, "使用者", p.WhoUsing)
	buf := make([]byte, 1024)
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
				n, _ := p.pty.Read(buf)
				p.bufHanle(buf[:n])
				if p.IsUsing.Load() {
					p.ws.wsMux.Lock()
					p.ws.wsConnect.WriteMessage(websocket.TextMessage, buf[:n])
					p.ws.wsMux.Unlock()
				}
			}
		}
	}
}

func (p *ProcessPty) ReadCache(ws *websocket.Conn) {
	ws.WriteMessage(websocket.TextMessage, p.cacheBytesBuf.Bytes())
}

func (p *ProcessPty) bufHanle(b []byte) {
	log := strings.TrimSpace(string(b))
	if utils.RemoveANSI(log) != "" {
		p.logReportHandler(log)
	}
	p.cacheBytesBuf.Write(b)
	p.cacheBytesBuf.Next(len(b))
}

func (p *ProcessPty) pInit() {
	p.SetTermType("pty")
	p.Control.changControlChan = make(chan int)
	p.stopChan = make(chan struct{})
	p.State.State = 1
	p.Pid = p.cmd.Process.Pid
	p.State.startTime = time.Now()
	p.cacheBytesBuf = bytes.NewBuffer(make([]byte, config.CF.ProcessMsgCacheBufLimit))
	p.InitPerformanceStatus()
	p.initPsutil()
	go p.readInit()
	go p.monitorHanler()
	go p.watchDog()
}

func RunNewProcessPty(pconfig model.Process) (*ProcessPty, error) {
	args := strings.Split(pconfig.Cmd, " ")
	cmd := exec.Command(args[0], args[1:]...) // 替换为你要执行的命令及参数

	processPty := ProcessPty{
		ProcessBase: ProcessBase{
			Name:         pconfig.Name,
			StartCommand: args,
		},
	}
	cmd.Dir = pconfig.Cwd
	pf, err := pty.Start(cmd)
	if err != nil || cmd.Process == nil {
		log.Logger.Error("进程启动出错:", err)
		return nil, err
	}
	pty.Setsize(pf, &pty.Winsize{
		Rows: 100,
		Cols: 100,
	})
	processPty.pty = pf
	processPty.cmd = cmd
	log.Logger.Infow("创建进程成功")
	processPty.setProcessConfig(pconfig)
	processPty.pInit()
	return &processPty, nil
}
