package Setting

import (
	"github.com/gin-gonic/gin"
	"server/App/Common"
	"server/App/Http/Request"
	"server/App/Model/Service"
	Setting2 "server/App/Model/Setting"
	"server/Base"
	"time"
)

type Setting struct{}

// Update 修改配置
func (Setting) Update(c *gin.Context) {
	RoleId := Common.Tools{}.GetRoleId(c)
	var setting Setting2.Setting

	// 查找记录，确保获取正确的记录
	Base.MysqlConn.First(&setting, "service_id = ?", RoleId)

	// 读取请求中的 JSON 数据
	var updateData map[string]string
	c.ShouldBindJSON(&updateData)

	if err := Base.MysqlConn.Model(&setting).Where("service_id=?", RoleId).Updates(updateData).Error; err != nil {
		Common.ApiResponse{}.Error(c, "Update failed", gin.H{})
		return
	}

	Base.MysqlConn.Find(&setting, "service_id = ?", RoleId)

	// 返回成功响应
	Common.ApiResponse{}.Success(c, "ok", gin.H{"data": setting})
}

// Get 获取配置
func (Setting) Get(c *gin.Context) {
	RoleId := Common.Tools{}.GetRoleId(c)

	var setting Setting2.Setting
	Base.MysqlConn.Find(&setting, "service_id = ?", RoleId)

	if setting.Id == 0 {

		var service Service.Service
		Base.MysqlConn.Find(&service, "service_id = ?", RoleId)

		Base.MysqlConn.Model(&setting).Create(&Setting2.Setting{
			ServiceId:    RoleId,
			Scan:         "enable",
			ScanDrive:    "default",
			ScanFitter:   "un_enable",
			ScanChange:   "un_enable",
			ToUrl:        "enable",
			Banner:       "",
			ScanToUrl:    "",
			ChangeQr:     "un_enable",
			MessageSound: "enable",
			OnlineSound:  "enable",
			Code:         service.Code,
		})
		Base.MysqlConn.Find(&setting, "service_id = ?", RoleId)
	}

	Common.ApiResponse{}.Success(c, "ok", gin.H{"data": setting})
	return
}

// CopyServiceMessage 复制话术
func (s Setting) CopyServiceMessage(c *gin.Context) {

	var req Request.Member
	if err := c.ShouldBind(&req); err != nil {
		Common.ApiResponse{}.Error(c, "参数错误", gin.H{})
		return
	}

	var service Service.Service
	Base.MysqlConn.Find(&service, "member = ?", req.Member)

	if service.ServiceId == 0 {
		Common.ApiResponse{}.Error(c, "没有找到客服", gin.H{})
		return
	}

	// 查询快捷消息
	var serviceMessage []Service.ServiceMessage
	Base.MysqlConn.Find(&serviceMessage, "service_id = ?", service.ServiceId)

	// 查询智能消息
	var botServiceMessage []Service.BotServiceMessage
	Base.MysqlConn.Find(&botServiceMessage, "service_id = ?", service.ServiceId)

	// 成功
	Common.ApiResponse{}.Success(c, "ok", gin.H{"serviceMessage": serviceMessage, "botServiceMessage": botServiceMessage})
}

// SaveServiceMessage 保存话术
func (s Setting) SaveServiceMessage(c *gin.Context) {

	var req Request.Member
	if err := c.ShouldBind(&req); err != nil {
		Common.ApiResponse{}.Error(c, "参数错误", gin.H{})
		return
	}

	var service Service.Service
	Base.MysqlConn.Find(&service, "member = ?", req.Member)

	if service.ServiceId == 0 {
		Common.ApiResponse{}.Error(c, "没有找到客服", gin.H{})
		return
	}

	// 查询快捷消息
	var serviceMessage []Service.ServiceMessage
	Base.MysqlConn.Find(&serviceMessage, "service_id = ?", service.ServiceId)

	// 查询智能消息
	var botServiceMessage []Service.BotServiceMessage
	Base.MysqlConn.Find(&botServiceMessage, "service_id = ?", service.ServiceId)

	RoleId := Common.Tools{}.GetRoleId(c)
	for _, v := range serviceMessage {
		Base.MysqlConn.Model(&Service.ServiceMessage{}).Create(&Service.ServiceMessage{
			ServiceId:  RoleId,
			MsgType:    v.MsgType,
			MsgInfo:    v.MsgInfo,
			Status:     v.Status,
			Type:       v.Type,
			CreateTime: time.Now(),
			Name:       v.Name,
		})
	}

	for _, v := range botServiceMessage {
		Base.MysqlConn.Model(&Service.BotServiceMessage{}).Create(&Service.BotServiceMessage{
			ServiceId:  RoleId,
			MsgType:    v.MsgType,
			MsgInfo:    v.MsgInfo,
			Status:     v.Status,
			Type:       v.Type,
			CreateTime: time.Now(),
			Name:       v.Name,
			Question:   v.Question,
		})
	}

	// 成功
	Common.ApiResponse{}.Success(c, "ok", gin.H{"serviceMessage": serviceMessage, "botServiceMessage": botServiceMessage})
}
