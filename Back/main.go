package main

import (
	"GoWorkingProject/Back/model"
	"GoWorkingProject/Back/router"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.DebugMode)

	if err := model.InitMySQL(); err != nil {
		log.Fatalf("❌ 初始化 MySQL 失败: %v", err)
	}

	if err := model.InitRedis(); err != nil {
		log.Fatalf("❌ 初始化 Redis 失败: %v", err)
	}

	r := gin.Default()
	
	// 静态文件服务，用于访问上传的头像
	r.Static("/uploads", "./uploads")
	
	router.RegisterRoutes(r)

	log.Println("🚀 GWP 启动成功，监听 :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("❌ 启动服务失败: %v", err)
	}
}
