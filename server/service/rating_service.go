package service

import (
	"errors"
	"fmt"
	"foodRatingSystem/database"
	"foodRatingSystem/model"
	"foodRatingSystem/repository"
	"time"
)

func SubmitReview(targetrest interface{}, userid uint, stars float64, comment string) error {
	rating := model.Rating{
		UserID:    userid,
		Stars:     stars,
		Comment:   comment,
		CreatedAt: time.Now(),
	}
	var resID uint
	if v, ok := targetrest.(int); ok {
		resID = uint(v)
	} else if v, ok := targetrest.(string); ok {
		var r model.Restaurant
		err := database.DB.Where("name = ?", v).First(&r).Error
		if err != nil {
			return errors.New("找不到餐厅[" + v + "]")
		}
		resID = r.ID
	} else {
		return errors.New("第一个参数必须是餐厅ID(int)或者餐厅名(string)")
	}

	//最终都换成id判断
	rating.RestaurantID = resID

	err := repository.AddRatingAndUpdateScore(rating)
	if err != nil {
		return err
	}

	ClearRatingCache(resID)
	return nil
}

func ClearRatingCache(restaurantID uint) {
	if database.RedisClient == nil {
		return
	}
	//这俩是精确的直接删
	database.RedisClient.Del(database.Ctx, fmt.Sprintf("restaurant:%d", restaurantID))
	database.RedisClient.Del(database.Ctx, fmt.Sprintf("ratings:%d", restaurantID))

	//这三个是处理后的间接的，需要根据餐厅id来匹配，所以一起放在切片里遍历查了再删
	patterns := []string{"recommend:*", "nearby:*", "search:*"}
	for _, pattern := range patterns {
		keys, _ := database.RedisClient.Keys(database.Ctx, pattern).Result()
		if len(keys) > 0 {
			database.RedisClient.Del(database.Ctx, keys...)
		}
	}
}

//当提交评分时，需要清除缓存
// - restaurant:{餐厅ID}
// - ratings:{餐厅ID}
// - recommend:*
// - nearby:*
// - search:*
