package router

import (
	"foodRatingSystem/handler"
	"foodRatingSystem/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default() //创建总路由 r 返回一个engine指针

	// 开启跨域
	r.Use(middleware.Cors()) //获得一个engine实例

	api := r.Group("/api") //用 r 创建分组 /api
	{
		api.GET("/restaurants/nearby", handler.GetNearbyRestaurants)
		api.GET("/restaurants", handler.GetRestaurants)
		api.GET("/restaurants/:id", handler.GetRestaurantDetail)
		api.GET("/restaurants/:id/ratings", handler.GetRestaurantRatings)

		api.POST("/rating", handler.SubmitRating)
		api.POST("/restaurants", handler.CreateRestaurant)

		api.POST("/user/register", handler.Register)
		api.POST("/user/login", handler.Login)
	}

	return r
}
