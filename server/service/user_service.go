package service

import (
	"errors"
	"fmt"
	"foodRatingSystem/database"
	"foodRatingSystem/model"
	"foodRatingSystem/repository"

	"golang.org/x/crypto/bcrypt"
)

func Register(user *model.User) (*model.User, error) {
	username := user.UserName
	password := user.PasswordHash
	if username == "" || password == "" {
		return nil, errors.New("用户名/密码为空！")
	}

	var count int64
	database.DB.Model(&model.User{}).Count(&count)
	user.UserID = fmt.Sprintf("user_%d", count+1)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("密码加密失败")
	}
	user.PasswordHash = string(hashedPassword)

	u, err := repository.CreateUser(user)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func Login(username, password string) (*model.User, error) {
	user, err := repository.GetUserByUsername(username)
	if err != nil {
		return nil, errors.New("用户不存在")
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, errors.New("密码错误")
	}

	return user, nil
}
