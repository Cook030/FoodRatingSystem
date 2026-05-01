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

	api := r.Group("/api")
	{
		api.GET("/restaurants/nearby", handler.GetNearbyRestaurants)
		api.GET("/restaurants", handler.GetRestaurants)
		api.GET("/restaurants/:id", handler.GetRestaurantDetail)
		api.GET("/restaurants/:id/ratings", handler.GetRestaurantRatings)

		api.POST("/user/register", handler.Register)
		api.POST("/user/login", handler.Login)

		//读公开，写保护
		auth := api.Group("", middleware.JWTAuth())
		{
			auth.POST("/rating", handler.SubmitRating)
			auth.POST("/restaurants", handler.CreateRestaurant)
		}
	}

	return r
}
