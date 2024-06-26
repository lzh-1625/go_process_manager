package dao

import (
	"msm/log"
	"msm/model"
)

type processDao struct{}

var ProcessDao = new(processDao)

func (p *processDao) GetAllProcessConfig() []model.Process {
	result := []model.Process{}

	tx := db.Find(&result)
	if tx.Error != nil {
		log.Logger.Error(tx.Error)
		return []model.Process{}
	}
	return result
}

func (p *processDao) GetProcessConfigByUser(username string) []model.Process {
	result := []model.Process{}
	tx := db.Debug().Raw(`SELECT p.uuid, p.name FROM permission left join process p where pid =p.uuid and owned  = 1 and account = ?`, username).Scan(&result)
	if tx.Error != nil {
		log.Logger.Error(tx.Error)
		return []model.Process{}
	}
	return result
}

func (p *processDao) UpdateProcessConfig(process model.Process) error {
	tx := db.Save(&process)
	return tx.Error
}

func (p *processDao) AddProcessConfig(process model.Process) (int, error) {
	tx := db.Create(&process)
	return process.Uuid, tx.Error
}

func (p *processDao) DeleteProcessConfig(uuid int) error {
	tx := db.Delete(&model.Process{
		Uuid: uuid,
	})
	return tx.Error
}

func (p *processDao) GetProcessConfigById(uuid int) model.Process {
	result := model.Process{}
	tx := db.Where(&model.Process{Uuid: uuid}).First(&result)
	if tx.Error != nil {
		log.Logger.Error(tx.Error)
		return model.Process{}
	}
	return result
}
