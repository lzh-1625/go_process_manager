package dao

import (
	"msm/model"

	"gorm.io/gorm"
)

type configDao struct{}

var ConfigDao = new(configDao)

func (c *configDao) GetConfigValue(key string) (string, error) {
	var result string
	if err := db.Model(&model.Config{}).Select("value").Where("key = ?", key).First(&result).Error; err != nil {
		return "", err
	}
	return result, nil
}

func (c *configDao) SetConfigValue(key, value string) error {
	if db.Model(&model.Config{}).Where("key = ?", key).First(nil).Error == gorm.ErrRecordNotFound {
		return db.Create(&model.Config{
			Key:   key,
			Value: value,
		}).Error
	} else {
		return db.Model(&model.Config{}).Where("key = ?", key).Updates(model.Config{Value: value}).Error
	}
}
