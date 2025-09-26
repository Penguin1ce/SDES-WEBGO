package router

import (
	"SDES/controller"

	"github.com/gin-gonic/gin"
)

func InitRouter(r *gin.Engine) {
	// 静态文件服务
	r.Static("/static", "./static")

	// 启用CORS中间件
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})
	// 主页路由
	r.GET("/", func(c *gin.Context) {
		c.File("./static/index.html")
	})

	// API路由
	baseApi := r.Group("/api")
	{
		baseApi.POST("/encrypt", controller.EncryptHandler)
		baseApi.POST("/decrypt", controller.DecryptHandler)
	}
}
