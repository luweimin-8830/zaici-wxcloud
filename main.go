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

	// 初始化配置服务和处理器
	configService := service.NewConfigService()
	configHandler := handler.NewConfigHandler(configService)

	// 初始化聊天服务和处理器
	chatService := service.NewChatService()
	chatHandler := handler.NewChatHandler(chatService)

	// 初始化详情记录服务和处理器
	detailService := service.NewDetailService()
	detailHandler := handler.NewDetailHandler(detailService, userService)

	// 初始化匹配服务和处理器
	matchService := service.NewMatchService()
	matchHandler := handler.NewMatchHandler(matchService)

	// 初始化在线服务和处理器
	onlineService := service.NewOnlineService()
	onlineHandler := handler.NewOnlineHandler(onlineService)

	// 初始化门店服务和处理器
	shopHandler := handler.NewShopHandler()

	// 初始化媒体服务和处理器
	mediaHandler := handler.NewMediaHandler()

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

	// 配置相关接口 /api/config/*
	configGroup := r.Group("/api/config")
	{
		configGroup.GET("/getDistance", configHandler.GetDistance)
		configGroup.POST("/saveDistance", configHandler.SaveDistance)
	}

	// 聊天相关接口 /api/chat/*
	chatGroup := r.Group("/api/chat")
	{
		chatGroup.POST("/get", chatHandler.GetMessages)
		chatGroup.POST("/send", chatHandler.SendMessage)
		chatGroup.POST("/update", chatHandler.UpdateState)
		chatGroup.POST("/updateState", chatHandler.UpdateStateByID)
		chatGroup.POST("/saveBlock", chatHandler.SaveBlock)
		chatGroup.POST("/delInfo", chatHandler.DeleteInfo)
		chatGroup.POST("/check", chatHandler.CheckContent)
	}

	// 微信内容安全回调
	r.POST("/censor", chatHandler.Censor)

	// GoEasy Webhook 回调
	r.POST("/webhook", chatHandler.Webhook)

	// 详情记录相关接口 /api/detailRecord/*
	detailGroup := r.Group("/api/detailRecord")
	{
		detailGroup.POST("/get", detailHandler.GetDetail)
		detailGroup.POST("/save", detailHandler.SaveDetail)
		detailGroup.POST("/saveSeat", detailHandler.SaveSeat)
	}

	// 匹配相关接口 /api/match/*
	matchGroup := r.Group("/api/match")
	{
		matchGroup.POST("/get", matchHandler.GetMatches)
		matchGroup.POST("/add", matchHandler.AddMatch)
		matchGroup.POST("/del", matchHandler.DeleteMatch)
		matchGroup.POST("/getLikeMatch", matchHandler.GetLikeMatch)
		matchGroup.POST("/sendMessage", matchHandler.SendMessage)
	}

// ... existing code ...
	// 在线相关接口 /api/online/*
	onlineGroup := r.Group("/api/online")
	{
		onlineGroup.GET("/status", onlineHandler.GetStatus)
		onlineGroup.POST("/save", onlineHandler.SaveOnline)
		onlineGroup.POST("/update", onlineHandler.UpdateOnline)
		onlineGroup.GET("/history", onlineHandler.GetHistory)
		onlineGroup.POST("/shop", onlineHandler.GetShop)
	}
// ... existing code ...

	// 门店相关接口 /api/shop/*
	shopGroup := r.Group("/api/shop")
	{
		shopGroup.GET("/detail", shopHandler.GetDetail)
		shopGroup.POST("/", shopHandler.GetNearList)
		shopGroup.POST("/save", shopHandler.Save)
		shopGroup.POST("/update", shopHandler.Update)
		shopGroup.POST("/del", shopHandler.Delete)
		shopGroup.GET("/admin", shopHandler.Admin)
		shopGroup.POST("/admin", shopHandler.Admin)
	}

	// 媒体相关接口 /api/media/*
	mediaGroup := r.Group("/api/media")
	{
		mediaGroup.GET("/list", mediaHandler.ListPictures)
		mediaGroup.POST("/create", mediaHandler.CreatePicture)
		mediaGroup.GET("/get/:id", mediaHandler.GetPictureByID)
		mediaGroup.POST("/update/:id", mediaHandler.UpdatePicture)
		mediaGroup.POST("/delete/:id", mediaHandler.DeletePicture)
		mediaGroup.POST("/checkStatus/:id", mediaHandler.UpdatePictureSecCheckStatus)
		mediaGroup.GET("/getByHash", mediaHandler.GetPictureByHash)
		mediaGroup.POST("/getHash", mediaHandler.GetHash)       // 兼容旧接口
		mediaGroup.POST("/startCensor", mediaHandler.StartCensor) // 从 index.js 迁移
	}

	log.Fatal(r.Run(":80"))
}
