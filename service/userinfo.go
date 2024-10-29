package service

import (
	"Abstract/config"
	"Abstract/model"
	"Abstract/utils"
	"fmt"
	"github.com/jaevor/go-nanoid"
	"gorm.io/gorm/clause"
	"log"
	"math/rand"
)

// search user
func FindByName(name string) model.UserBasic {
	var exist_data model.UserBasic
	config.DB.Model(&model.UserBasic{}).Where("name = ?", name).First(&exist_data)
	return exist_data
}
func FindUserByNameAndPwd(name, password string) model.UserBasic {
	var exist_data model.UserBasic
	config.DB.Model(&model.UserBasic{}).Where("name = ? and password = ?", name, password).First(&exist_data)
	return exist_data
}

// Create new user
func CreatUser(user model.UserBasic) (rep interface{}, err error) {
	tx := config.DB.Begin()
	var exist_data model.UserBasic
	result := config.DB.Model(&model.UserBasic{}).Where("name = ?", user.Name).First(&exist_data)
	if result.Error == nil {
		tx.Rollback()
		return -1, nil
	}
	salt := fmt.Sprintf("%06d", rand.Int31())
	user.Password = utils.MakePassword(user.Password, salt)

	gen, err := nanoid.Canonic()
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}
	data := &model.UserBasic{
		UUID:     gen(),
		Name:     user.Name,
		Password: user.Password,
		Salt:     salt,
	}
	result = tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&data)
	if result.Error != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()
	return data.ID, nil
}

// login
func Login(name string, password string) (rep interface{}, err error) {
	var exist_data model.UserBasic
	result := config.DB.Model(&model.UserBasic{}).Where("name = ? AND pass_word = ?", name, password).First(&exist_data)
	if result.Error != nil {
		return nil, result.Error
	}
	// Generate JWT token
	//atoken, rtoken, err := utils.GenToken(uint64(exist_data.ID), name)
	//exist_data.AccessToken = atoken
	//exist_data.RefreshToken = rtoken
	return exist_data, nil
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

	tx.Commit()
	//exist.AccessToken = aToken
	//exist.RefreshToken = rToken
	return exist, nil
}
