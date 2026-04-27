package handler

import (
	"foodRatingSystem/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 前端请求 → Handler接收 → 解析参数 → 调用Service → 拿到结果 → 返回JSON
// handler只做收参数、调service、返回json
func GetNearbyRestaurants(c *gin.Context) {
	latStr := c.Query("lat")
	lonStr := c.Query("lon")

	lat, _ := strconv.ParseFloat(latStr, 64) //将字符串解析为浮点数
	lon, _ := strconv.ParseFloat(lonStr, 64) //64 (对应 float64)

	data, err := service.GetNearbyRestaurants(lat, lon)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		//gin.H{}等价于
		//c.JSON(http.StatusInternalServerError,map[string]any{}{"error":err.Error()})
		return
	}

	c.JSON(http.StatusOK, data)
}

// 搜索框
func GetRestaurants(c *gin.Context) {
	latStr := c.Query("lat")
	lonStr := c.Query("lon")
	search := c.Query("search")
	sortBy := c.DefaultQuery("sort", "distance")

	lat, _ := strconv.ParseFloat(latStr, 64)
	lon, _ := strconv.ParseFloat(lonStr, 64)

	data, err := service.GetRestaurants(lat, lon, search, sortBy)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, data)
}

func GetRestaurantDetail(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的餐厅ID"})
		return
	}

	restaurant, err := service.GetRestaurantByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "餐厅不存在"})
		return
	}

	c.JSON(http.StatusOK, restaurant)
}

func GetRestaurantRatings(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的餐厅ID"})
		return
	}

	ratings, err := service.GetRatingsByRestaurantID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, ratings)
}

func CreateRestaurant(c *gin.Context) {
	var req struct {
		Name      string  `json:"name"`
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
		Category  string  `json:"category"`
	}

	err := c.ShouldBindJSON(&req) //接收前端 JSON，通过结构体后面 反引号 `` 里面的 json:"xxx"对应上，自动把值塞进 req 结构体，格式不对就返回错误
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	restaurant, err := service.CreateRestaurant(req.Name, req.Latitude, req.Longitude, req.Category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, restaurant)
}
