package service

import (
	"github.com/lzh-1625/go_process_manager/log"
)

func (p *ProcessBase) initCgroup() {
	log.Logger.Debugw("不支持cgroup")
}
