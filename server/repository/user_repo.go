package repository

import (
	"errors"
	"foodRatingSystem/database"
	"foodRatingSystem/model"

	"gorm.io/gorm"
)

func CreateUser(user *model.User) (*model.User, error) {
	err := database.DB.Create(&user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func GetUserByUsername(user_name string) (*model.User, error) {
	var u model.User
	err := database.DB.Where("user_name = ?", user_name).First(&u).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &u, nil
}

func GetUserByID(user_id uint) (*model.User, error) {
	var u model.User
	err := database.DB.First(&u, user_id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &u, nil
}
