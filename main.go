package main

import (
	"SDES/router"
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("S-DES Web 服务器启动中...")

	// 设置Gin模式
	gin.SetMode(gin.DebugMode)

	// 创建Gin路由器
	r := gin.Default()

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

	router.InitRouter(r)

	// 启动服务器
	fmt.Println("服务器启动在 http://localhost:8080")
	fmt.Println("API端点:")
	fmt.Println("  POST /api/encrypt - 加密")
	fmt.Println("  POST /api/decrypt - 解密")

	err := r.Run(":8080")
	if err != nil {
		return
	}
}
