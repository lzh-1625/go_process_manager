package model

import "time"

type User struct {
	Account    string    `json:"account" gorm:"primaryKey;column:account" `
	Password   string    `json:"password" gorm:"column:password" `
	Role       int       `json:"role" gorm:"column:role" `
	CreateTime time.Time `json:"createTime" gorm:"column:create_time" `
	Remark     string    `json:"remark" gorm:"column:remark" `
}

func (*User) TableName() string {
	return "users"
}
