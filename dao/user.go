package dao

import (
	"errors"
	"msm/config"
	"msm/model"
	"msm/utils"
	"time"
)

type userDao struct{}

var UserDao = new(userDao)

func (u *userDao) GetUserByName(name string) model.User {
	var result model.User
	db.Where("account = ?", name).First(&result)
	return result
}

func (u *userDao) CreateUser(user model.User) error {
	if len(user.Password) < config.CF.UserPassWordMinLength {
		return errors.New("密码小于最小长度")
	}
	user.Password = utils.Md5(user.Password)
	user.CreateTime = time.Now()
	tx := db.Create(&user)
	return tx.Error
}

func (u *userDao) UpdatePassword(name string, password string) error {
	if len(password) < config.CF.UserPassWordMinLength {
		return errors.New("新密码太短")
	}
	tx := db.Model(&model.User{}).Where("account = ?", name).Update("password", utils.Md5(password))
	return tx.Error
}

func (u *userDao) DeleteUser(name string) error {
	if err := db.Where("account = ?", name).First(&model.User{}).Error; err != nil {
		return err
	}
	tx := db.Delete(&model.User{Account: name})
	return tx.Error
}

func (u *userDao) GetUserList() []model.User {
	result := []model.User{}
	db.Find(&result)
	return result
}
