package main

import (
	"fmt"
	"log"
	"wxcloudrun-golang/db"
	"wxcloudrun-golang/handler"
	"wxcloudrun-golang/service"

	"github.com/gin-gonic/gin"
)

func main() {
	if err := db.Init(); err != nil {
		panic(fmt.Sprintf("mysql init failed with %+v", err))
	}

	r := gin.Default()

	// 初始化用户服务和处理器
	userService := service.NewUserService()
	userHandler := handler.NewUserHandler(userService)

	// 初始化 Banner 服务和处理器
	bannerService := service.NewBannerService()
	bannerHandler := handler.NewBannerHandler(bannerService)

	// 用户相关接口 /api/user/*
	userGroup := r.Group("/api/user")
	{
		userGroup.POST("/", userHandler.GetOrCreateUser)
		userGroup.POST("/save", userHandler.SaveUser)
		userGroup.POST("/superLike", userHandler.UseSuperLike)
		userGroup.POST("/addSuperLike", userHandler.AddSuperLike)
		userGroup.POST("/addAdmin", userHandler.AddAdmin)
	}

	// Banner 相关接口 /api/banner/*
	bannerGroup := r.Group("/api/banner")
	{
		bannerGroup.GET("/", bannerHandler.ListBanners)
		bannerGroup.GET("/detail", bannerHandler.GetBannerDetail)
		bannerGroup.POST("/save", bannerHandler.SaveBanner)
		bannerGroup.POST("/del", bannerHandler.DeleteBanner)
	}

	log.Fatal(r.Run(":80"))
}
