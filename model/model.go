package model

import (
	"gorm.io/gorm"
)

type UserBasic struct {
	gorm.Model
	UUID          int64 `gorm:"column:uuid" json:"uuid"`
	Name          string
	Password      string
	Phone         string `valid:"matches(^1[3-9]{1}\\d{9}$)"`
	Email         string `valid:"email"`
	Identity      string
	ClientIp      string
	ClientPort    string
	LoginTime     uint64
	HeartBeatTime uint64
	LogOutTime    uint64 `gorm:"column:login_out_time" json:"login_out_time"`
	DeviceInfo    string
	AccessToken   string
	RefreshToken  string
}

func (table UserBasic) TableName() string {
	return "user_basic"
}
