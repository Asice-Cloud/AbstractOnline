package service

import (
	"Chat/config"
	"Chat/model"
	"gorm.io/gorm/clause"
)

// Search all user
func GetUserList() []model.UserBasic {
	var data []model.UserBasic
	config.DB.Model(&model.UserBasic{}).Find(&data)
	return data
}

// Create new user
func CreatUser(user model.UserBasic) (rep interface{}, err error) {
	tx := config.DB.Begin()
	var exist_user model.UserBasic
	result := tx.Where("name=?", user.Name).First(&exist_user)
	if result.Error == nil {
		tx.Rollback()
		return -1, nil
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
	tx := config.DB.Begin()
	result := tx.Where("id = ?", user.ID).First(&existingUser)
	if result.Error != nil {
		return result.Error
	}
	result = tx.Delete(&existingUser) // Change this line
	if result.Error != nil {
		return result.Error
	}
	tx.Commit()
	return nil
}

// Update user
func UpdateUser(user model.UserBasic) (rep interface{}, err error) {
	tx := config.DB.Begin()
	var exist model.UserBasic
	result := tx.Where("id=?", user.ID).First(&exist)
	if result.Error != nil {
		tx.Rollback()
		return -1, result.Error
	}
	result = tx.Model(&exist).Updates(map[string]interface{}{
		"Name":     user.Name,
		"Password": user.Password,
	})
	if result.Error != nil {
		tx.Rollback()
		return nil, result.Error
	}
	tx.Commit()
	return exist.ID, nil
}
