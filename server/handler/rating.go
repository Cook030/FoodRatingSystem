package handler

import (
	"foodRatingSystem/database"
	"foodRatingSystem/model"
	"foodRatingSystem/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RatingRequest struct {
	RestaurantID   int     `json:"restaurant_id"`
	RestaurantName string  `json:"restaurant_name"`
	Username       string  `json:"username"`
	Stars          float64 `json:"stars"`
	Comment        string  `json:"comment"`
}

func SubmitRating(c *gin.Context) {
	var req RatingRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数格式错误"})
		return
	}

	userID := c.GetString("user_id")
	userName := c.GetString("user_name")

	var user model.User
	if err := database.DB.Where("id = ? AND user_name = ?", userID, userName).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户不存在"})
		return
	}

	var targetrest interface{}
	if req.RestaurantID > 0 {
		targetrest = req.RestaurantID
	} else {
		targetrest = req.RestaurantName
	}

	err = service.SubmitReview(targetrest, user.ID, req.Stars, req.Comment)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "评价成功！"})
}
