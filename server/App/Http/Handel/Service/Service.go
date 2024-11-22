package Service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"server/App/Common"
	Service3 "server/App/Http/Request/Service"
	Common3 "server/App/Logic/Common"
	Common2 "server/App/Model/Common"
	"server/App/Model/Log"
	Service2 "server/App/Model/Service"
	"server/Base"
	"strconv"
	"time"
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

// CheckDomain 切换域名
func (s Service) CheckDomain(c *gin.Context) {
	// 找到客服
	RoleId := Common.Tools{}.GetRoleId(c)
	var service Service2.Service
	Base.MysqlConn.Find(&service, "service_id = ?", RoleId)

	var domain1 Common2.Domain
	Base.MysqlConn.Find(&domain1, "status = ? and  type = ?  and bind_service_id = 0", "enable", "public")
	if domain1.Id == 0 {
		Common.ApiResponse{}.Success(c, "无可分配域名", gin.H{})
		return
	}

	// 找到对应的域名
	domain := Common3.Domain{}.GetServiceBind(RoleId)
	Base.MysqlConn.Find(&domain, "bind_service_id = ?", RoleId)

	// 5 分钟内是否切换过新域名
	var checkDomainLog Log.CheckDomainLog
	Base.MysqlConn.Find(&checkDomainLog, "service_id = ? and  create_time >= ?", RoleId, time.Now().Add(-5*time.Minute).Unix())

	if checkDomainLog.Id != 0 {
		Common.ApiResponse{}.Success(c, "失败", gin.H{})
		return
	}

	// 修改域名下掉
	if domain.BindServiceId > 0 {
		Base.MysqlConn.Model(&domain).
			Where("id = ?", domain.Id).
			Updates(map[string]interface{}{"status": "un_enable", "bind_service_id": 0})

		//增加记录
		Base.MysqlConn.Create(&Log.CheckDomainLog{
			ServiceId:  RoleId,
			CreateTime: time.Now().Unix(),
		})
	}

	domainStruct := &Common3.Domain{}
	err1 := domainStruct.Bind(RoleId)
	if err1 == nil {
		Common.ApiResponse{}.Success(c, "ok", gin.H{})
	}

}
