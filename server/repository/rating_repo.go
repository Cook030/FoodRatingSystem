package repository

import (
	"foodRatingSystem/database"
	"foodRatingSystem/model"

	"gorm.io/gorm"
)

func AddRatingAndUpdateScore(r model.Rating) error {
	//Transaction 开启数据库事务
	return database.DB.Transaction(func(tx *gorm.DB) error {
		err := tx.Create(&r).Error //新增一条评价
		if err != nil {
			return err
		}
		var avg float64
		var count int64
		tx.Model(&model.Rating{}).Where("restaurant_id = ?", r.RestaurantID).Count(&count)
		tx.Model(&model.Rating{}).Where("restaurant_id = ?", r.RestaurantID).Select("COALESCE(AVG(stars), 0)").Scan(&avg)
		result := tx.Model(&model.Restaurant{}).Where("id = ?", r.RestaurantID).Updates(map[string]interface{}{
			"avg_score":    avg,
			"review_count": count,
		}) //更新平均分和评价数
		return result.Error
	})
}

func GetRatingsByRestaurantID(restaurantID int) ([]model.Rating, error) {
	var ratings []model.Rating
	err := database.DB.Preload("User").Where("restaurant_id = ?", restaurantID).Order("created_at DESC").Find(&ratings).Error
	return ratings, err
}
