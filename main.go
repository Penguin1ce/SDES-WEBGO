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
