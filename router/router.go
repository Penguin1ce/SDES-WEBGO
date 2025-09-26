package router

import (
	"SDES/controller"

	"github.com/gin-gonic/gin"
)

func InitRouter(r *gin.Engine) {
	// 静态文件服务
	r.Static("/static", "./static")

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
