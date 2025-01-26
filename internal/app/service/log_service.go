package service

import (
	"github.com/lzh-1625/go_process_manager/config"
	"github.com/lzh-1625/go_process_manager/internal/app/model"
	"github.com/lzh-1625/go_process_manager/internal/app/repository"
)

type LogService interface {
	Search(req model.GetLogReq, filterProcessName ...string) model.LogResp
	Insert(log string, processName string, using string, ts int64)
}

var LogServiceImpl LogService

func InitLog() {
	if config.CF.EsEnable {
		LogServiceImpl = LogEs
	} else {
		LogServiceImpl = LogSqlite
	}
}

type logSqlite struct{}

var LogSqlite = new(logSqlite)

func (l *logSqlite) Search(req model.GetLogReq, filterProcessName ...string) model.LogResp {
	req.FilterName = filterProcessName
	data, total := repository.LogRepository.SearchLog(req)
	return model.LogResp{
		Data:  data,
		Total: total,
	}
}

func (l *logSqlite) Insert(log string, processName string, using string, ts int64) {
	repository.LogRepository.InsertLog(model.ProcessLog{
		Log:   log,
		Name:  processName,
		Using: using,
		Time:  ts,
	})
}

type logEs struct{}

var LogEs = new(logEs)

func (l *logEs) Search(req model.GetLogReq, filterProcessName ...string) model.LogResp {
	return EsService.Search(req, filterProcessName...)
}

func (l *logEs) Insert(log string, processName string, using string, ts int64) {
	EsService.Insert(log, processName, using, ts)
}
