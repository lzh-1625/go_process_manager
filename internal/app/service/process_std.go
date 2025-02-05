package service

import (
	"bufio"
	"io"
	"os/exec"
	"strings"

	"github.com/lzh-1625/go_process_manager/config"
	"github.com/lzh-1625/go_process_manager/internal/app/constants"
	"github.com/lzh-1625/go_process_manager/internal/app/model"
	"github.com/lzh-1625/go_process_manager/log"
	"github.com/lzh-1625/go_process_manager/utils"
)

type ProcessStd struct {
	*ProcessBase
	cacheLine []string
	stdin     io.WriteCloser
	stdout    *bufio.Scanner
}

func (p *ProcessStd) Type() constants.TerminalType {
	return constants.TERMINAL_STD
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

func (p *ProcessStd) Start() (err error) {
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

	out, err := cmd.StdoutPipe()
	if err != nil {
		log.Logger.Errorw("启动失败，输出管道获取失败", "err", err)
		return err
	}
	p.stdout = bufio.NewScanner(out)
	p.stdin, err = cmd.StdinPipe()
	if err != nil {
		log.Logger.Errorw("启动失败，输入管道获取失败", "err", err)
		return err
	}
	err = cmd.Start()
	if err != nil {
		log.Logger.Errorw("启动失败，进程启动出错:", "err", err)
		return err
	}
	log.Logger.Infow("进程启动成功", "重启次数", p.State.restartTimes)
	p.cmd = cmd
	p.pInit()
	p.push("进程启动成功")
	return nil
}

func (p *ProcessStd) doOnInit() {
	p.cacheLine = make([]string, config.CF.ProcessMsgCacheLinesLimit)
}

func (p *ProcessStd) ReadCache(ws ConnectInstance) {
	for _, line := range p.cacheLine {
		ws.WriteString(line)
	}
}

func (p *ProcessStd) doOnKilled() {
	// 不执行如何操作
}

func (p *ProcessStd) SetTerminalSize(cols, rows int) {
	log.Logger.Debug("当前终端不支持修改尺寸")
}

func (p *ProcessStd) readInit() {
	var output string
	log.Logger.Debugw("stdout读取线程已启动", "进程名", p.Name, "使用者", p.GetUserString())
	for {
		select {
		case <-p.StopChan:
			{
				log.Logger.Debugw("stdout读取线程已退出", "进程名", p.Name, "使用者", p.GetUserString())
				return
			}
		default:
			{
				output = p.Read()
				if len(p.ws) == 0 {
					continue
				}
				p.wsLock.Lock()
				for _, v := range p.ws {
					v.WriteString(output)
				}
				p.wsLock.Unlock()
			}
		}
	}
}
func (p *ProcessStd) Read() string {
	if p.stdout.Scan() {
		output := utils.RemoveNotValidUtf8InString(p.stdout.Text())
		p.logReportHandler(output)
		p.cacheLine = p.cacheLine[1:]
		p.cacheLine = append(p.cacheLine, output)
		return output
	}
	return ""
}

func NewProcessStd(pconfig model.Process) *ProcessBase {
	p := ProcessBase{
		Name:         pconfig.Name,
		StartCommand: strings.Split(pconfig.Cmd, " "),
		WorkDir:      pconfig.Cwd,
	}
	processStd := ProcessStd{
		ProcessBase: &p,
	}
	p.Process = &processStd
	processStd.setProcessConfig(pconfig)
	return &p
}

func RunNewProcessStd(pconfig model.Process) (*ProcessBase, error) {
	processStd := NewProcessStd(pconfig)
	if err := processStd.Start(); err != nil {
		return nil, err
	}
	return processStd, nil
}
