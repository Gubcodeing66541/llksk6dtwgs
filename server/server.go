package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

//
//import (
//	"fmt"
//	"github.com/gin-gonic/gin"
//	swaggerFiles "github.com/swaggo/files"
//	ginSwagger "github.com/swaggo/gin-swagger"
//	"server/App"
//	"server/Base"
//	_ "server/docs" //依赖项必须导入
//)
//
//// @title 系统API文档`
//// @version 1.0`
//// @description 系统api `
//// @description 后端：`
//func main() {
//
//	//启动初始化
//	Base.Base{}.Init()
//
//	// 启动web服务
//	HttpServer := gin.Default()
//
//	App.HttpRoute{}.BindRoute(HttpServer)
//
//	if Base.AppConfig.Debug {
//		HttpServer.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
//	}
//
//	////监听消息队列启动
//	//if Base.AppConfig.Model != "dev" {
//	//	Server.NsqServer{}.Run("group")
//	//}
//
//	//启动服务
//	_ = HttpServer.Run(fmt.Sprintf(":%d", Base.AppConfig.Port))
//}

// GeoLocationResponse 用于映射 API 返回的 JSON 数据
type GeoLocationResponse struct {
	Status  int `json:"status"`
	Content struct {
		AddressDetail struct {
			Province string `json:"province"`
			City     string `json:"city"`
		} `json:"address_detail"`
	} `json:"content"`
}

func getGeoLocation(ip string) string {
	apiKey := "WrjgBGDGi2gS3a7a9B2P2f5U5WaE4FBh" // 替换为你申请的 API Key
	url := fmt.Sprintf("http://api.map.baidu.com/location/ip?ak=%s&ip=%s", apiKey, ip)

	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	// 打印原始的响应体
	fmt.Println("Raw Response Body:", string(body))

	// 将返回的 JSON 数据解析到结构体中
	var geoResp GeoLocationResponse
	if err := json.Unmarshal(body, &geoResp); err != nil {
		log.Fatal(err)
	}

	// 打印详细的地址信息
	return fmt.Sprintf("Province: %s, City: %s", geoResp.Content.AddressDetail.Province, geoResp.Content.AddressDetail.City)
}

func main() {
	ip := "182.143.151.43" // 示例 IP 地址
	data := getGeoLocation(ip)
	fmt.Println(data)
}
