package Service

import (
	"Chat/Model"
	"Chat/utils"
	"gorm.io/gorm/clause"
)

// Search all user
func GetUserList() []Model.UserBasic {
	var data []Model.UserBasic
	utils.DB.Model(&Model.UserBasic{}).Find(&data)
	return data
}

// Add new user
func CreatUser(user Model.UserBasic) (rep interface{}, err error) {
	tx := utils.DB.Begin()
	var exist_user Model.UserBasic
	result := tx.Model(&Model.UserBasic{}).Where("name=?", user.Name).First(&exist_user)
	if result.Error == nil {
		tx.Rollback()
		return "Star already exist", nil
	}
	data := &Model.UserBasic{
		Name:     user.Name,
		Password: user.Password,
	}
	result = tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&data)
	if result.Error != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()
	return data.ID, nil
}

// Delete existing user
func DeleteUser(user Model.UserBasic) error {
	var existingUser Model.UserBasic
	tx := utils.DB.Begin()
	result := tx.Model(&Model.UserBasic{}).Where("id = ?", user.ID).First(&existingUser)
	if result.Error != nil {
		return result.Error
	}
	result = tx.Model(&Model.UserBasic{}).Where("id=?", existingUser.ID).Unscoped().Delete(&existingUser)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// Update user
func UpdateUser(user Model.UserBasic) (rep interface{}, err error) {
	tx := utils.DB.Begin()
	var exist Model.UserBasic
	result := tx.Model(&Model.UserBasic{}).Where("id=?", user.ID).First(&exist)
	if result.Error != nil {
		tx.Rollback()
		return
	}
	tx.Model(&exist).Updates(Model.UserBasic{Name: user.Name, Password: user.Password})
	tx.Commit()
	return exist.ID, nil
}
