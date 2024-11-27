package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"server/App"
	"server/App/Server"
	"server/Base"
	_ "server/docs" //依赖项必须导入
)

// @title 系统API文档`
// @version 1.0`
// @description 系统api `
// @description 后端：`
func main() {

	//启动初始化
	Base.Base{}.Init()

	// 启动web服务
	HttpServer := gin.Default()

	App.HttpRoute{}.BindRoute(HttpServer)

	if Base.AppConfig.Debug {
		HttpServer.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	//监听消息队列启动
	if Base.AppConfig.Model != "dev" {
		Server.NsqServer{}.Run("group")
	}

	//启动服务
	_ = HttpServer.Run(fmt.Sprintf(":%d", Base.AppConfig.Port))
}
