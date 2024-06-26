package loghandler

import (
	"msm/log"
	"msm/model"
	"msm/service/es"
	"time"

	"github.com/panjf2000/ants"
)

type loghandler struct{}

var (
	antsPool   *ants.PoolWithFunc
	Loghandler = new(loghandler)

	logHanleFunc = func(i interface{}) {
		esLog, ok := i.(model.Eslog)
		if !ok {
			log.Logger.Panicw("传入错误参数", "data", esLog)
			return
		}
		es.EsService.Insert(esLog.Log, esLog.Name, esLog.Using, esLog.Time)
	}

	panicHanlderFunc = func(i interface{}) {
		log.Logger.Error("es消息储存失败")
	}
)

func init() {
	antsPool, _ = ants.NewPoolWithFunc(1000, logHanleFunc, ants.WithPanicHandler(panicHanlderFunc), ants.WithExpiryDuration(time.Second*10))
}

func (l *loghandler) AddLog(data model.Eslog) {
	if err := antsPool.Invoke(data); err != nil {
		log.Logger.Errorw("协程池添加任务失败", "err", err, "当前运行数量", antsPool.Running())
	}
}
