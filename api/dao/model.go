package dao

import "gorm.io/gorm"

type UserInfo struct {
	gorm.Model
	UserName string
	PassWord string
}

func (UserInfo) TableName() string {
	return "userinfo"
}
