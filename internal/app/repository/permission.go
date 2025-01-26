package repository

import (
	"reflect"

	"github.com/lzh-1625/go_process_manager/internal/app/constants"
	"github.com/lzh-1625/go_process_manager/internal/app/model"
	"github.com/lzh-1625/go_process_manager/log"

	"gorm.io/gorm"
)

type permissionRepository struct{}

var PermissionRepository = new(permissionRepository)

func (p *permissionRepository) GetPermssionList(account string) []model.PermissionPo {
	result := []model.PermissionPo{}
	if err := db.Raw(`SELECT p.name ,p.uuid as pid,p2.owned ,p2."start" ,p2.stop ,p2.terminal,p2.log ,p2.write    
	FROM users u full join process p left join permission p2 on p2.account == u.account and p2.pid =p.uuid 
	WHERE u.account = ? or u.account ISNULL`, account).Find(&result); err.Error != nil {
		log.Logger.Warnw("权限查询失败", "err", err)
	}

	return result
}

func (p *permissionRepository) EditPermssion(data model.Permission) error {
	if db.Model(&model.Permission{}).Where(&model.Permission{Account: data.Account, Pid: data.Pid}).First(nil).Error == gorm.ErrRecordNotFound {
		db.Omit("name").Create(&model.Permission{
			Account: data.Account,
			Pid:     data.Pid,
		})
	}
	return db.Model(&model.Permission{}).Where(&model.Permission{Account: data.Account, Pid: data.Pid}).Updates(map[string]interface{}{
		"owned":    data.Owned,
		"start":    data.Start,
		"stop":     data.Stop,
		"terminal": data.Terminal,
		"log":      data.Log,
		"write":    data.Write,
	}).Error
}

func (p *permissionRepository) GetPermission(user string, pid int) (result model.Permission) {
	db.Model(&model.Permission{}).Where(&model.Permission{Account: user, Pid: int32(pid)}).First(&result)
	return
}

func (p *permissionRepository) GetProcessNameByPermission(user string, op constants.OprPermission) (result []string) {
	query := model.PermissionPo{Account: user, Owned: true}
	reflect.ValueOf(&query).Elem().FieldByName(string(op)).SetBool(true)
	db.Model(&model.Permission{}).Select("name").Joins("right join process p on p.uuid = pid").Where(query).Find(&result)
	return
}
