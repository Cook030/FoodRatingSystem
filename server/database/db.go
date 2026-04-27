package database

import (
	"fmt"

	"foodRatingSystem/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connectdb() {
	dsn := "host=localhost user=postgres password=root dbname=postgres port=5432 sslmode=disable"
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("无法连接数据库：", err)
		return
	}
	fmt.Println("成功连接到数据库")

	DB.AutoMigrate(&model.Restaurant{}, &model.Rating{}, &model.User{})
	fmt.Println("数据库迁移完成")
}
