package repository

import (
	"github.com/lzh-1625/go_process_manager/internal/app/model"
)

type configRepository struct{}

var ConfigRepository = new(configRepository)

func (c *configRepository) GetConfigValue(key string) (string, error) {
	var result string
	if err := db.Model(&model.Config{}).Select("value").Where(&model.Config{Key: key}).First(&result).Error; err != nil {
		return "", err
	}
	return result, nil
}

func (c *configRepository) SetConfigValue(key, value string) error {
	config := model.Config{Key: key}
	updateData := model.Config{Value: &value}
	err := db.Model(&config).Where(&config).Assign(updateData).FirstOrCreate(&config).Error
	if err != nil {
		return err
	}
	return nil
}
