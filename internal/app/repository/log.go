package repository

import (
	"github.com/lzh-1625/go_process_manager/internal/app/model"
	"github.com/lzh-1625/go_process_manager/log"
)

type logRepository struct{}

var LogRepository = new(logRepository)

func (l *logRepository) InsertLog(data model.ProcessLog) {
	if err := db.Create(&data).Error; err != nil {
		log.Logger.Errorw("日志插入失败", "err", err)
	}
}

func (l *logRepository) SearchLog(query model.GetLogReq) (result []model.ProcessLog, total int64) {
	tx := db.Model(&model.ProcessLog{}).Where(&model.ProcessLog{
		Name:  query.Match.Name,
		Using: query.Match.Using,
	})
	if query.Match.Log != "" {
		tx.Where("log like ?", "%"+query.Match.Log+"%")
	}
	if query.Sort == "desc" {
		tx.Order("time desc")
	}
	if query.TimeRange.StartTime != 0 {
		tx.Where("time > ?", query.TimeRange.StartTime)
	}
	if query.TimeRange.EndTime != 0 {
		tx.Where("time < ?", query.TimeRange.EndTime)
	}
	if len(query.FilterName) != 0 {
		tx.Where("name in ?", query.FilterName)
	}
	tx.Count(&total)
	tx.Limit(query.Page.Size)
	tx.Offset(query.Page.From)
	tx.Find(&result)
	return
}
