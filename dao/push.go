package dao

import (
	"msm/model"
)

type pushDao struct{}

var PushDao = new(pushDao)

func (p *pushDao) GetPushList() (result []model.Push) {
	db.Find(&result)
	return
}

func (p *pushDao) GetPushConfigById(id int) (result model.Push) {
	db.Where("id = ?", id).First(&result)
	return
}

func (p *pushDao) UpdatePushConfig(data model.Push) error {
	return db.Save(&data).Error
}

func (p *pushDao) AddPushConfig(data model.Push) error {
	return db.Create(&data).Error
}

func (p *pushDao) DeletePushConfig(id int) error {
	return db.Delete(&model.Push{
		Id: int64(id),
	}).Error
}
