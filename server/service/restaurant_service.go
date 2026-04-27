package service

import (
	"encoding/json"
	"fmt"
	"foodRatingSystem/database"
	"foodRatingSystem/model"
	"foodRatingSystem/repository"
	"foodRatingSystem/utils"
	"math"
	"sort"
	"strings"
	"time"
)

type RestaurantWithDistance struct {
	model.Restaurant
	Distance float64 `json:"distance"`
}

func GetNearbyRestaurants(userLat, userLon float64) ([]RestaurantWithDistance, error) {
	//1.设置一个缓存键 nearby:用户纬度:用户经度
	cacheKey := fmt.Sprintf("nearby:%.2f:%.2f", userLat, userLon)
	//2.从缓存查
	cacheData, err := database.RedisClient.Get(database.Ctx, cacheKey).Result()
	//3.如果查到了返回
	if err == nil {
		var results []RestaurantWithDistance
		json.Unmarshal([]byte(cacheData), &results)
		return results, nil
	}

	//4.如果没查到，开始用数据库查，然后存入缓存
	rests, err := repository.GetAllRestaurants() //获得所有餐厅，包括评价数量
	if err != nil {
		return nil, err
	}
	var rwd []RestaurantWithDistance
	for _, rest := range rests {
		resLat := rest.Latitude
		resLon := rest.Longitude
		dist := utils.Distance(userLat, userLon, resLat, resLon)
		rwd = append(rwd, RestaurantWithDistance{
			Restaurant: rest,
			Distance:   dist,
		})
	}
	sort.Slice(rwd, func(i, j int) bool {
		return rwd[i].Distance < rwd[j].Distance
	})

	//5.缓存没查到，那就把上面查到的数据存到缓存
	data, _ := json.Marshal(rwd)
	database.RedisClient.Set(database.Ctx, cacheKey, data, 2*time.Hour)

	return rwd, nil
}

type RestaurantWithScore struct {
	model.Restaurant
	Distance   float64 `json:"distance"`
	FinalScore float64 `json:"final_score"`
}

func GetRecommendedRestaurants(userLat, userLon float64) ([]RestaurantWithScore, error) {
	//因为有两个及以上的参数所以用fmt.Sprintf格式化
	//据指定的格式和参数，生成并返回一个格式化后的字符串（不会直接打印到控制台）
	cacheKey := fmt.Sprintf("recommend:%.2f:%.2f", userLat, userLon)
	cacheData, err := database.RedisClient.Get(database.Ctx, cacheKey).Result()
	if err == nil {
		var results []RestaurantWithScore
		json.Unmarshal([]byte(cacheData), &results)
		return results, nil
	}

	rests, err := repository.GetAllRestaurants()
	if err != nil {
		return nil, err
	}

	var results []RestaurantWithScore
	for _, rest := range rests {
		dist := utils.Distance(userLat, userLon, rest.Latitude, rest.Longitude)

		// --- 核心算法实现 ---
		// 1. 评分权重 (0.6)
		scorePart := rest.AverageScore * 0.6

		// 2. 距离权重 (0.3): 距离越小，(1/(dist+1)) 越大
		distPart := (1.0 / (dist + 1.0)) * 0.3

		// 3. 人气权重 (0.1): 使用 log10 平滑处理，防止大店评价数霸榜
		reviewPart := math.Log10(float64(rest.ReviewCount)+1.0) * 0.1

		finalScore := scorePart + distPart + reviewPart

		results = append(results, RestaurantWithScore{
			Restaurant: rest,
			Distance:   dist,
			FinalScore: finalScore,
		})
	}

	// 按照综合得分从高到低排序 (降序)
	sort.Slice(results, func(i, j int) bool {
		return results[i].FinalScore > results[j].FinalScore
	})

	data, _ := json.Marshal(results)
	database.RedisClient.Set(database.Ctx, cacheKey, data, 2*time.Hour)

	return results, nil
}

// 搜索框相关逻辑
func GetRestaurants(userLat, userLon float64, search, sortBy string) ([]RestaurantWithScore, error) {
	cacheKey := fmt.Sprintf("search:%.2f:%.2f:%s:%s", userLat, userLon, search, sortBy)
	cacheData, err := database.RedisClient.Get(database.Ctx, cacheKey).Result()
	if err == nil {
		var results []RestaurantWithScore
		json.Unmarshal([]byte(cacheData), &results)
		return results, nil
	}

	rests, err := repository.GetAllRestaurants()
	if err != nil {
		return nil, err
	}

	var results []RestaurantWithScore
	for _, rest := range rests {
		if search != "" && !isContainsKeywords(rest.Name, search) {
			continue
		}

		dist := utils.Distance(userLat, userLon, rest.Latitude, rest.Longitude)

		scorePart := rest.AverageScore * 0.6
		distPart := (1.0 / (dist + 1.0)) * 0.3
		reviewPart := math.Log10(float64(rest.ReviewCount)+1.0) * 0.1

		finalScore := scorePart + distPart + reviewPart

		results = append(results, RestaurantWithScore{
			Restaurant: rest,
			Distance:   dist,
			FinalScore: finalScore,
		})
	}

	switch sortBy {
	case "score":
		sort.Slice(results, func(i, j int) bool {
			return results[i].AverageScore > results[j].AverageScore
		})
	case "reviews":
		sort.Slice(results, func(i, j int) bool {
			return results[i].ReviewCount > results[j].ReviewCount
		})
	case "recommended":
		sort.Slice(results, func(i, j int) bool {
			return results[i].FinalScore > results[j].FinalScore
		})
	default:
		sort.Slice(results, func(i, j int) bool {
			return results[i].Distance < results[j].Distance
		})
	}

	data, _ := json.Marshal(results)
	database.RedisClient.Set(database.Ctx, cacheKey, data, 2*time.Hour)

	return results, nil
}

func isContainsKeywords(s, substr string) bool {
	s = strings.ToLower(s)
	substr = strings.ToLower(substr)
	return strings.Contains(s, substr)
}

func GetRestaurantByID(id int) (*model.Restaurant, error) {
	cacheKey := fmt.Sprintf("restaurant:%d", id)
	cacheData, err := database.RedisClient.Get(database.Ctx, cacheKey).Result()
	if err == nil {
		var rest model.Restaurant
		json.Unmarshal([]byte(cacheData), &rest)
		return &rest, nil
	}

	rest, err := repository.GetRestaurantByID(id)
	if err != nil {
		return nil, err
	}
	data, _ := json.Marshal(rest)
	database.RedisClient.Set(database.Ctx, cacheKey, data, 10*time.Minute)
	return rest, nil
}

func GetRatingsByRestaurantID(restaurantID int) ([]model.Rating, error) {
	cacheKey := fmt.Sprintf("ratings:%d", restaurantID)
	cacheData, err := database.RedisClient.Get(database.Ctx, cacheKey).Result()
	if err == nil {
		var ratings []model.Rating
		json.Unmarshal([]byte(cacheData), &ratings)
		return ratings, nil
	}

	ratings, err := repository.GetRatingsByRestaurantID(restaurantID)
	if err != nil {
		return nil, err
	}

	data, _ := json.Marshal(ratings)
	database.RedisClient.Set(database.Ctx, cacheKey, data, 10*time.Minute)

	return ratings, nil
}

// 创建餐厅会改变 PostgreSQL 数据库里的数据，对应的 Redis 缓存就要删掉
func CreateRestaurant(name string, lat, lon float64, category string) (*model.Restaurant, error) {
	rest := model.Restaurant{
		Name:      name,
		Latitude:  lat,
		Longitude: lon,
		Category:  category,
		// 没写平均分 结构体默认赋0
	}

	result, err := repository.CreateRestaurant(rest)
	if err == nil {
		ClearListCache()
		restaurantCacheKey := fmt.Sprintf("restaurant:%d", result.ID)
		database.RedisClient.Del(database.Ctx, restaurantCacheKey)
	}
	return result, err
}

func ClearListCache() {
	if database.RedisClient == nil {
		return
	}
	patterns := []string{"recommend:*", "nearby:*", "search:*"}
	for _, pattern := range patterns {
		keys, _ := database.RedisClient.Keys(database.Ctx, pattern).Result()
		if len(keys) > 0 {
			database.RedisClient.Del(database.Ctx, keys...)
		}
	}
}

// 当创建新餐厅时
// ✅ 需要清除：因为列表里会多一个餐厅
// - recommend:*
// - nearby:*
// - search:*

// ❌ 不需要清除：因为
// - restaurant:{新ID}  ← 根本还没被缓存过（刚创建）
// - ratings:{新ID}     ← 根本还没被缓存过（还没评分）
// - restaurant:{旧ID}  ← 旧餐厅数据没变，为什么要删？
