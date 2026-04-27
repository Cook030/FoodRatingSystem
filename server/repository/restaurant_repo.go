package repository

import (
	"foodRatingSystem/database"
	"foodRatingSystem/model"
)

func GetAllRestaurants() ([]model.Restaurant, error) {
	var restaurants []model.Restaurant
	err := database.DB.Find(&restaurants).Error
	if err != nil {
		return nil, err
	}
	for i := range restaurants {
		var count int64
		database.DB.Model(&model.Rating{}). // 去查 rating 评分表
							Where("restaurant_id = ?", restaurants[i].ID). // 条件：只查属于这家餐厅的评分
							Count(&count)                                  // 统计一共有多少条，结果放进 count 变量里
		//SELECT COUNT(*) FROM ratings WHERE restaurant_id = $
		restaurants[i].ReviewCount = int(count)
	}
	return restaurants, nil
}

func GetRestaurantByID(id int) (*model.Restaurant, error) {
	var r model.Restaurant
	//把查找到的第一个符合要求的数据存到r里
	err := database.DB.First(&r, id).Error
	if err != nil {
		return nil, err
	}
	var count int64
	database.DB.Model(&model.Rating{}).Where("restaurant_id = ?", r.ID).Count(&count)
	r.ReviewCount = int(count)
	return &r, nil
}

func CreateRestaurant(rest model.Restaurant) (*model.Restaurant, error) {
	err := database.DB.Create(&rest).Error //自动在restaurants表里新建一行
	if err != nil {
		return nil, err
	}
	return &rest, nil
}
