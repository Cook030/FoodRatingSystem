package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// 先经过中间件 → 再进入路由 → 最后执行 Handler
func Cors() gin.HandlerFunc {
	//Cors() 是创建函数， Cors() 返回的匿名函数才是处理每次请求的处理函数
	return func(c *gin.Context) {
		method := c.Request.Method // 获取当前HTTP 请求的方法，比如 GET、POST、OPTIONS、PUT 等
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")

		// 放行所有 OPTIONS 方法
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		c.Next()
	}
}

//浏览器请求
//↓
//【1. CORS 中间件】（处理跨域）
//↓
//c.Next() 继续往下
//↓
//【2. Gin 路由】
//匹配 /api/restaurants/nearby
//↓
//【3.  Handler】
//GetNearbyRestaurants
//↓
//返回 JSON 给前端
