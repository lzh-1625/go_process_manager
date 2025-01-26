package repository

import (
	"github.com/lzh-1625/go_process_manager/internal/app/model"
)

type pushRepository struct{}

var PushRepository = new(pushRepository)

func (p *pushRepository) GetPushList() (result []model.Push) {
	db.Find(&result)
	return
}

func (p *pushRepository) GetPushConfigById(id int) (result model.Push) {
	db.Model(&model.Push{}).Where(&model.Push{Id: int64(id)}).First(&result)
	return
}

func (p *pushRepository) UpdatePushConfig(data model.Push) error {
	return db.Save(&data).Error
}

func (p *pushRepository) AddPushConfig(data model.Push) error {
	return db.Create(&data).Error
}

func (p *pushRepository) DeletePushConfig(id int) error {
	return db.Delete(&model.Push{
		Id: int64(id),
	}).Error
}
