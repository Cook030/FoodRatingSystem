package repository

import (
	"foodRatingSystem/database"
	"foodRatingSystem/model"

	"gorm.io/gorm"
)

func AddRatingAndUpdateScore(r model.Rating) error {
	//开启数据库事务（Transaction）
	return database.DB.Transaction(func(tx *gorm.DB) error {
		err := tx.Create(&r).Error //新增一条评价
		if err != nil {
			return err
		}
		var avg float64
		var count int64
		tx.Model(&model.Rating{}).Where("restaurant_id = ?", r.RestaurantID).Count(&count)                                //统计评价总数
		tx.Model(&model.Rating{}).Where("restaurant_id = ?", r.RestaurantID).Select("COALESCE(AVG(stars), 0)").Scan(&avg) //计算平均分
		result := tx.Model(&model.Restaurant{}).Where("id = ?", r.RestaurantID).Updates(map[string]interface{}{
			"avg_score":    avg,
			"review_count": count,
		}) //更新平均分和评价数
		return result.Error
	})
}

func GetRatingsByRestaurantID(restaurantID int) ([]model.Rating, error) {
	var ratings []model.Rating
	err := database.DB.Where("restaurant_id = ?", restaurantID).Order("created_at DESC").Find(&ratings).Error
	return ratings, err
}
