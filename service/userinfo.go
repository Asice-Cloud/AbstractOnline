package service

import (
	"Chat/config"
	"Chat/model"
	"Chat/pkg"
	"gorm.io/gorm"

	"gorm.io/gorm/clause"
)

// direct a user:
func FindUserByName(name string) *gorm.DB {
	var exist_data model.UserBasic
	return config.DB.Model(&model.UserBasic{}).Where("name = ?", name).First(&exist_data)
}

func FinduUserByPhone(phone string) *gorm.DB {
	var exist_data model.UserBasic
	return config.DB.Model(&model.UserBasic{}).Where("name = ?", phone).First(&exist_data)
}

func FinduUserByEmail(email string) *gorm.DB {
	var exist_data model.UserBasic
	return config.DB.Model(&model.UserBasic{}).Where("name = ?", email).First(&exist_data)
}

// login
func Login(name string, password string) (rep interface{}, err error) {
	var exist_data model.UserBasic
	result := config.DB.Model(&model.UserBasic{}).Where("name = ? AND password = ?", name, password).First(&exist_data)
	if result.Error != nil {
		return nil, result.Error
	}
	// Generate JWT token
	atoken, rtoken, err := pkg.GenToken(uint64(exist_data.ID), name)
	exist_data.AccessToken = atoken
	exist_data.RefreshToken = rtoken
	return exist_data, nil
}

// Create new user
func CreatUser(user model.UserBasic) (rep interface{}, err error) {
	tx := config.DB.Begin()
	result := FindUserByName(user.Name)
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
	if err := user.OptimisticLock(tx); err != nil {
		tx.Rollback()
		return nil, err
	}
	var exist model.UserBasic
	result := tx.Where("id=?", user.ID).First(&exist)
	if result.Error != nil {
		tx.Rollback()
		return -1, result.Error
	}
	result = tx.Model(&exist).Updates(&model.UserBasic{
		Name:     user.Name,
		Password: user.Password,
		Phone:    user.Phone,
		Email:    user.Email,
	})
	if result.Error != nil {
		tx.Rollback()
		return nil, result.Error
	}
	aToken, rToken, err := pkg.GenToken(uint64(exist.ID), exist.Name)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()
	exist.AccessToken = aToken
	exist.RefreshToken = rToken
	return exist, nil
}
