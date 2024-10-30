package model

import (
	"gorm.io/gorm"
)

type UserBasic struct {
	gorm.Model
	UUID        string `gorm:"column:uuid" json:"uuid"`
	Name        string
	Password    string
	Phone       string `valid:"matches(^1[3-9]{1}\\d{9}$)"`
	Email       string `valid:"email"`
	Identity    string
	ClientIp    string
	ClientPort  string
	Salt        string
	DeviceInfo  string
	AccessToken string
}

func (table UserBasic) TableName() string {
	return "user_basic"
}
