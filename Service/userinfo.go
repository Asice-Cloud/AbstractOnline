package Service

import (
	"Chat/Model"
	"Chat/utils"
)

func GetUserList() []Model.UserBasic {
	var data []Model.UserBasic
	utils.DB.Find(&data)
	return data
}
