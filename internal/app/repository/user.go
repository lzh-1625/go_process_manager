package repository

import (
	"time"

	"github.com/lzh-1625/go_process_manager/internal/app/model"
	"github.com/lzh-1625/go_process_manager/utils"
)

type userRepository struct{}

var UserRepository = new(userRepository)

func (u *userRepository) GetUserByName(name string) model.User {
	var result model.User
	db.Model(&model.User{}).Where(&model.User{Account: name}).First(&result)
	return result
}

func (u *userRepository) CreateUser(user model.User) error {
	user.Password = utils.Md5(user.Password)
	user.CreateTime = time.Now()
	tx := db.Create(&user)
	return tx.Error
}

func (u *userRepository) UpdatePassword(name string, password string) error {
	tx := db.Model(&model.User{}).Where(&model.User{Account: name}).Updates(&model.User{Password: utils.Md5(password)})
	return tx.Error
}

func (u *userRepository) DeleteUser(name string) error {
	if err := db.Model(&model.User{}).Where(&model.User{Account: name}).First(&model.User{}).Error; err != nil {
		return err
	}
	tx := db.Delete(&model.User{Account: name})
	return tx.Error
}

func (u *userRepository) GetUserList() []model.User {
	result := []model.User{}
	db.Find(&result)
	return result
}
