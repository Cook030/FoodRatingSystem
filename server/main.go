package main

import (
	"fmt"
	"foodRatingSystem/config"
	"foodRatingSystem/database"
	"foodRatingSystem/router"
)

func main() {
	config.LoadConfig()
	database.Connectdb()
	database.ConnectRedis()

	r := router.SetupRouter()
	fmt.Println("服务已启动，监听端口 :8080")
	r.Run(":8080")
}
