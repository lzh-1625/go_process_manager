package service

import (
	"bytes"
	"os"
	"os/exec"
	"strings"

	"github.com/lzh-1625/go_process_manager/config"
	"github.com/lzh-1625/go_process_manager/internal/app/constants"
	"github.com/lzh-1625/go_process_manager/internal/app/model"
	"github.com/lzh-1625/go_process_manager/log"
	"github.com/lzh-1625/go_process_manager/utils"

	"github.com/creack/pty"
)

type ProcessPty struct {
	*ProcessBase
	cacheBytesBuf *bytes.Buffer
	pty           *os.File
}

func (p *ProcessPty) doOnKilled() {
	p.pty.Close()
}

func (p *ProcessPty) Type() constants.TerminalType {
	return constants.TERMINAL_PTY
}

func (p *ProcessPty) Start() (err error) {
	defer func() {
		log.Logger.DeleteAdditionalInfo(1)
		if err != nil {
			p.Config.AutoRestart = false
			p.SetState(constants.PROCESS_WARNNING)
			p.State.Info = "进程启动失败:" + err.Error()
		}
	}()
	log.Logger.AddAdditionalInfo("进程名称", p.Name)
	if ok := p.SetState(constants.PROCESS_START, func() bool {
		return p.State.State != 1
	}); !ok {
		log.Logger.Warnw("进程已在运行，跳过启动")
		return nil
	}
	cmd := exec.Command(p.StartCommand[0], p.StartCommand[1:]...)
	cmd.Dir = p.WorkDir
	pf, err := pty.Start(cmd)
	if err != nil || cmd.Process == nil {
		log.Logger.Errorw("进程启动失败", "err", err)
		return err
	}
	pty.Setsize(pf, &pty.Winsize{
		Rows: 100,
		Cols: 100,
	})
	p.pty = pf
	log.Logger.Infow("进程启动成功", "进程名称", p.Name, "重启次数", p.State.restartTimes)
	p.cmd = cmd
	p.pInit()
	p.push("进程启动成功")
	return nil
}

func (p *ProcessPty) SetTerminalSize(cols, rows int) {
	if cols == 0 || rows == 0 || len(p.ws) != 0 {
		return
	}
	if err := pty.Setsize(p.pty, &pty.Winsize{
		Cols: uint16(cols),
		Rows: uint16(rows),
	}); err != nil {
		log.Logger.Error("设置终端尺寸失败", "err", err)
	}

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
	log.Logger.Debugw("stdout读取线程已启动", "进程名", p.Name, "使用者", p.GetUserString())
	buf := make([]byte, 1024)
	for {
		select {
		case <-p.StopChan:
			{
				log.Logger.Debugw("stdout读取线程已退出", "进程名", p.Name, "使用者", p.GetUserString())
				return
			}
		default:
			{
				n, _ := p.pty.Read(buf)
				p.bufHanle(buf[:n])
				for _, v := range p.ws {
					v.Write(buf[:n])
				}
			}
		}
	}
}

func (p *ProcessPty) ReadCache(ws ConnectInstance) {
	ws.Write(p.cacheBytesBuf.Bytes())
}

func (p *ProcessPty) bufHanle(b []byte) {
	log := strings.TrimSpace(string(b))
	if utils.RemoveANSI(log) != "" {
		p.logReportHandler(log)
	}
	p.cacheBytesBuf.Write(b)
	p.cacheBytesBuf.Next(len(b))
}

func (p *ProcessPty) doOnInit() {
	p.cacheBytesBuf = bytes.NewBuffer(make([]byte, config.CF.ProcessMsgCacheBufLimit))
}

func NewProcessPty(pconfig model.Process) *ProcessBase {
	p := ProcessBase{
		Name:         pconfig.Name,
		StartCommand: strings.Split(pconfig.Cmd, " "),
		WorkDir:      pconfig.Cwd,
	}
	processPty := ProcessPty{
		ProcessBase: &p,
	}
	p.Process = &processPty
	processPty.setProcessConfig(pconfig)
	return &p
}

func RunNewProcessPty(pconfig model.Process) (*ProcessBase, error) {
	processPty := NewProcessPty(pconfig)
	if err := processPty.Start(); err != nil {
		return nil, err
	}
	return processPty, nil
}
