package main

import (
	"fmt"
	"tools-admin/backend/common/config"
	"tools-admin/backend/middleware/cors"
	"tools-admin/backend/router"

	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化配置
	server := config.Config.Server
	r := gin.Default()

	// 跨域
	r.Use(cors.Cors())

	// 初始化路由
	router.InitRouter(r)

	// 启动应用
	if err := r.Run(":" + server.Port); err != nil {
		fmt.Println("Failed to run server on port ", server.Port, ":", err)
	}
}
