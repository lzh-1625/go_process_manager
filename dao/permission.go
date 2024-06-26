package dao

import (
	"msm/log"
	"msm/model"

	"gorm.io/gorm"
)

type permissionDao struct{}

var PermissionDao = new(permissionDao)

func (p *permissionDao) GetPermssionList(account string) []model.PermissionPo {
	result := []model.PermissionPo{}
	if err := db.Raw(`SELECT p.name ,p.uuid as pid,p2.owned ,p2."start" ,p2.stop ,p2.terminal  
	FROM users u full join process p left join permission p2 on p2.account == u.account and p2.pid =p.uuid WHERE u.account = ? or u.account ISNULL`, account).Find(&result); err.Error != nil {
		log.Logger.Warnw("权限查询失败", "err", err)
	}

	return result
}

func (p *permissionDao) EditPermssion(data model.Permission) error {
	if db.Model(&model.Permission{}).Where("account = ? and pid = ?", data.Account, data.Pid).First(nil).Error == gorm.ErrRecordNotFound {
		db.Omit("name").Create(&model.Permission{
			Account: data.Account,
			Pid:     data.Pid,
		})
	}
	return db.Debug().Model(&model.Permission{}).Where("account = ? and pid = ?", data.Account, data.Pid).Updates(map[string]interface{}{
		"owned":    data.Owned,
		"start":    data.Start,
		"stop":     data.Stop,
		"terminal": data.Terminal,
	}).Error
}

func (p *permissionDao) GetPermission(user string, pid int) (result model.Permission) {
	db.Debug().Model(&model.Permission{}).Where("account = ? and pid = ?", user, pid).First(&result)
	return
}
