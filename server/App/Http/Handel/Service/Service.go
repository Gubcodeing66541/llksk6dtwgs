package Service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"server/App/Common"
	Service3 "server/App/Http/Request/Service"
	Common3 "server/App/Logic/Common"
	Service2 "server/App/Model/Service"
	"server/Base"
	"strconv"
)

type Service struct{}

// @summary 获取客服基本信息
// @tags 客服信息
// @Param token header string true "认证token"
// @Router /api/service/info [post]
func (Service) Info(c *gin.Context) {
	RoleId := Common.Tools{}.GetRoleId(c)
	var service Service2.Service
	Base.MysqlConn.Find(&service, "service_id = ?", RoleId)

	domain := Common3.Domain{}.GetServiceBind(RoleId)
	Base.MysqlConn.Find(&domain, "bind_service_id = ?", RoleId)

	share := fmt.Sprintf("%s/api/user/ua/join/%s", domain.Domain, service.Code)

	diversion := fmt.Sprintf("%s/api/user/ua/join/%s", domain.Domain, service.Code)
	Common.ApiResponse{}.Success(c, "OK", gin.H{"service": service, "user_share": share, "diversion": diversion})
}

// @summary 修改客服头像和昵称
// @tags 客服信息
// @Param token header string true "认证token"
// @Param name query string true "客服名称"
// @Param head query string false "客服头像"
// @Router /api/service/update [post]
func (Service) Update(c *gin.Context) {
	var req Service3.UpdateServiceDetail
	err := c.ShouldBind(&req)
	if err != nil {
		Common.ApiResponse{}.Error(c, "参数有误", gin.H{})
		return
	}

	roleId := Common.Tools{}.GetRoleId(c)
	Base.MysqlConn.Model(&Service2.Service{}).Where("service_id = ?", roleId).Updates(req.ToUpdate())

	Common.ApiResponse{}.Success(c, "ok", gin.H{"REQ": req})
}

// @summary 创建视频通话
// @tags 客服信息
// @Param token header string true "认证token"
// @Param name query string true "客服名称"
// @Param head query string false "客服头像"
// @Router /api/service/create_live_token/2 [post]
func (Service) CreateLive(c *gin.Context) {
	uid := c.Param("uid")
	userId, _ := strconv.Atoi(uid)
	serviceId := Common.Tools{}.GetRoleId(c)

	// 密钥
	ApiKey := Base.AppConfig.LiveAppKey
	ApiSecret := Base.AppConfig.LiveAppSecret

	// service
	key := Common.Tools{}.CreateActiveCode(serviceId)
	token := Common.Tools{}.CreateToken(
		ApiKey, ApiSecret, Common.Tools{}.GetLiveRoomName(int64(serviceId), int64(userId)), fmt.Sprintf("S:%d", serviceId),
	)
	serviceLink := fmt.Sprintf("%s/api/live/%s", Common3.Domain{}.GetLive(), key)
	Common.RedisTools{}.SetStringByTime(key, token, 60*60)

	// user
	key = Common.Tools{}.CreateActiveCode(serviceId)
	token = Common.Tools{}.CreateToken(
		ApiKey, ApiSecret, Common.Tools{}.GetLiveRoomName(int64(serviceId), int64(userId)), fmt.Sprintf("U:%d", userId),
	)
	userLink := fmt.Sprintf("%s/api/live/%s", Common3.Domain{}.GetLive(), key)
	Common.RedisTools{}.SetStringByTime(key, token, 60*60)

	// response
	Common.ApiResponse{}.Success(c, "ok", gin.H{"serviceLink": serviceLink, "user_link": userLink})
}
