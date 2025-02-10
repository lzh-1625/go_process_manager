package logic

import (
	"time"

	"github.com/lzh-1625/go_process_manager/config"
	"github.com/lzh-1625/go_process_manager/internal/app/model"
	"github.com/lzh-1625/go_process_manager/log"

	"github.com/panjf2000/ants/v2"
)

type loghandler struct {
	antsPool *ants.Pool
}

var (
	Loghandler = new(loghandler)
)

func InitLogHandle() {
	Loghandler.antsPool, _ = ants.NewPool(config.CF.LogHandlerPoolSize, ants.WithNonblocking(true), ants.WithExpiryDuration(3*time.Second), ants.WithPanicHandler(func(i interface{}) {
		log.Logger.Error("es消息储存失败")
	}))
}

func (l *loghandler) AddLog(data model.ProcessLog) {
	if err := l.antsPool.Submit(func() {
		LogLogicImpl.Insert(data.Log, data.Name, data.Using, data.Time)
	}); err != nil {
		log.Logger.Warnw("协程池添加任务失败", "err", err, "当前运行数量", l.antsPool.Running())
	}
}

func (l *loghandler) GetRunning() int {
	return l.antsPool.Running()
}
