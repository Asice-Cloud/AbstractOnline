package service

import (
	"Chat/config"
	"Chat/model"
)

// Search all user
func GetUserList() []model.UserBasic {
	var data []model.UserBasic
	config.DB.Model(&model.UserBasic{}).Find(&data)
	return data
}
