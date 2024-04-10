package service

import (
	"Chat/model"
	"Chat/utils"
	"gorm.io/gorm/clause"
)

// Search all user
func GetUserList() []model.UserBasic {
	var data []model.UserBasic
	utils.DB.Model(&model.UserBasic{}).Find(&data)
	return data
}

// Add new user
func CreatUser(user model.UserBasic) (rep interface{}, err error) {
	tx := utils.DB.Begin()
	var exist_user model.UserBasic
	result := tx.Model(&model.UserBasic{}).Where("name=?", user.Name).First(&exist_user)
	if result.Error == nil {
		tx.Rollback()
		return "Star already exist", nil
	}
	data := &model.UserBasic{
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
func DeleteUser(user model.UserBasic) error {
	var existingUser model.UserBasic
	tx := utils.DB.Begin()
	result := tx.Model(&model.UserBasic{}).Where("id = ?", user.ID).First(&existingUser)
	if result.Error != nil {
		return result.Error
	}
	result = tx.Model(&model.UserBasic{}).Where("id=?", existingUser.ID).Unscoped().Delete(&existingUser)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// Update user
func UpdateUser(user model.UserBasic) (rep interface{}, err error) {
	tx := utils.DB.Begin()
	var exist model.UserBasic
	result := tx.Model(&model.UserBasic{}).Where("id=?", user.ID).First(&exist)
	if result.Error != nil {
		tx.Rollback()
		return
	}
	tx.Model(&exist).Updates(model.UserBasic{Name: user.Name, Password: user.Password})
	tx.Commit()
	return exist.ID, nil
}
