package repository

import (
	"errors"

	"github.com/lzh-1625/go_process_manager/internal/app/model"
	"github.com/lzh-1625/go_process_manager/log"

	"gorm.io/gorm"
)

type processRepository struct{}

var ProcessRepository = new(processRepository)

func (p *processRepository) GetAllProcessConfig() []model.Process {
	result := []model.Process{}

	tx := db.Find(&result)
	if tx.Error != nil {
		log.Logger.Error(tx.Error)
		return []model.Process{}
	}
	return result
}

func (p *processRepository) GetProcessConfigByUser(username string) []model.Process {
	result := []model.Process{}
	tx := db.Raw(`SELECT p.* FROM permission left join process p where pid =p.uuid and owned  = 1 and account = ?`, username).Scan(&result)
	if tx.Error != nil {
		log.Logger.Error(tx.Error)
		return []model.Process{}
	}
	return result
}

func (p *processRepository) UpdateProcessConfig(process model.Process) error {
	tx := db.Save(&process)
	return tx.Error
}

func (p *processRepository) AddProcessConfig(process model.Process) (int, error) {
	var existingProcess model.Process
	err := db.Model(&model.Process{}).Where("name = ?", process.Name).First(&existingProcess).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		log.Logger.Error(err)
		return 0, err
	}

	if err == nil {

		return 0, errors.New("process name already exists")
	}

	tx := db.Create(&process)
	if tx.Error != nil {
		log.Logger.Error(tx.Error)
		return 0, tx.Error
	}

	return process.Uuid, nil
}

func (p *processRepository) DeleteProcessConfig(uuid int) error {
	return db.Delete(&model.Process{
		Uuid: uuid,
	}).Error
}

func (p *processRepository) GetProcessConfigById(uuid int) model.Process {
	result := model.Process{}
	tx := db.Model(&model.Process{}).Where(&model.Process{Uuid: uuid}).First(&result)
	if tx.Error != nil {
		log.Logger.Error(tx.Error)
		return model.Process{}
	}
	return result
}
