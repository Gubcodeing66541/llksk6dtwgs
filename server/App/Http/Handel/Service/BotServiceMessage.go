package Service

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"math"
	Common2 "server/App/Common"
	Service3 "server/App/Http/Request/Service"
	Service4 "server/App/Logic/Service"
	Service2 "server/App/Model/Service"
	"server/Base"
)

type BotServiceMessage struct{}

func (BotServiceMessage) Create(c *gin.Context) {
	var req Service3.CreateBotServiceMessage
	err := c.ShouldBind(&req)
	if err != nil {
		Common2.ApiResponse{}.Error(c, "参数错误", gin.H{})
		return
	}

	serviceId := Common2.Tools{}.GetServiceId(c)
	req.ServiceId = Common2.Tools{}.GetServiceId(c)
	res := Service4.BotServiceMessage{}.Create(serviceId, req)
	Common2.ApiResponse{}.Success(c, "创建成功", gin.H{"data": res})
}

func (BotServiceMessage) Delete(c *gin.Context) {
	var req Service3.DeleteBotServiceMessage
	err := c.ShouldBind(&req)
	if err != nil {
		Common2.ApiResponse{}.Error(c, "参数错误", gin.H{})
		return
	}

	req.ServiceId = Common2.Tools{}.GetServiceId(c)
	Service4.BotServiceMessage{}.Delete(req.Id, req.ServiceId)
	Common2.ApiResponse{}.Success(c, "删除成功", gin.H{})
}

func (BotServiceMessage) Update(c *gin.Context) {
	var req Service3.UpdateBotServiceMessage
	err := c.ShouldBind(&req)
	if err != nil {
		Common2.ApiResponse{}.Error(c, "参数错误", gin.H{})
		return
	}
	req.ServiceId = Common2.Tools{}.GetServiceId(c)
	Service4.BotServiceMessage{}.Update(req.Id, req.Name, req.MsgType, req.MsgInfo, req.ServiceId, req.Status, req.Title, req.Question)
	Common2.ApiResponse{}.Success(c, "操作成功", gin.H{"id": req.Id})
}

func (BotServiceMessage) List(c *gin.Context) {
	var pageReq Service3.ListBotServiceMessage
	err := c.ShouldBind(&pageReq)
	if err != nil {
		Common2.ApiResponse{}.Error(c, "请提交完整的分页参数", gin.H{"err": err.Error(), "req": pageReq})
		return
	}

	ServiceId := Common2.Tools{}.GetServiceId(c)

	tel := Base.MysqlConn.Model(&Service2.BotServiceMessage{}).Where("service_id = ? and type = ?", ServiceId, pageReq.Type)

	// 计算分页和总数
	var allCount int
	tel.Count(&allCount)
	allPage := math.Ceil(float64(allCount) / float64(pageReq.Limit))

	// 获取分页数据
	var list []Service2.BotServiceMessage
	tel.Offset(pageReq.GetOffset()).Limit(pageReq.GetLimit()).Scan(&list)
	res := gin.H{"count": allCount, "page": allPage, "current_page": pageReq.Page, "list": list}
	Common2.ApiResponse{}.Success(c, "获取成功", res)
}

func (BotServiceMessage) GetById(c *gin.Context) {
	var req Service3.GetByIdBotServiceMessage
	err := c.ShouldBind(&req)
	if err != nil {
		Common2.ApiResponse{}.Error(c, "参数错误", gin.H{})
		return
	}

	req.ServiceId = Common2.Tools{}.GetServiceId(c)
	serviceMessage := Service4.BotServiceMessage{}.GetById(req.ServiceId, req.Id)
	Common2.ApiResponse{}.Success(c, "获取成功", gin.H{"serviceMessage": serviceMessage})
}

func (BotServiceMessage) Swap(c *gin.Context) {
	var req Service3.SwapBotServiceMessage
	err := c.ShouldBind(&req)
	if err != nil {
		Common2.ApiResponse{}.Error(c, "参数错误", gin.H{})
		return
	}

	ServiceId := Common2.Tools{}.GetServiceId(c)
	from := Service4.BotServiceMessage{}.GetById(ServiceId, req.From)
	to := Service4.BotServiceMessage{}.GetById(ServiceId, req.To)
	from.Id, to.Id = to.Id, from.Id
	Base.MysqlConn.Save(&from)
	Base.MysqlConn.Save(&to)
	Common2.ApiResponse{}.Success(c, "修改成功", gin.H{"from": from, "to": to})
}

func (m BotServiceMessage) SetEnable(c *gin.Context) {
	var req Service3.BotServiceMessageId
	err := c.ShouldBind(&req)
	if err != nil {
		Common2.ApiResponse{}.Error(c, "参数错误", gin.H{})
		return
	}
	// SetEnable 设置启用

	serviceID := Common2.Tools{}.GetServiceId(c)
	// 查询对应的 ServiceMessage
	var serviceMessage Service2.BotServiceMessage
	if err := Base.MysqlConn.Where("service_id = ? and id = ?", serviceID, req.Id).First(&serviceMessage).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			Common2.ApiResponse{}.Error(c, "未找到对应的信息", nil)
		} else {
			Common2.ApiResponse{}.Error(c, "数据库查询错误", gin.H{})
		}
		return
	}

	// 切换 Status 值
	newStatus := "enable"
	if serviceMessage.Status == "enable" {
		newStatus = "un_enable"
	}

	// 更新 Status 值到数据库
	if err := Base.MysqlConn.Model(&serviceMessage).Update("status", newStatus).Error; err != nil {
		Common2.ApiResponse{}.Error(c, "状态更新失败", gin.H{})
		return
	}

	// 返回成功响应
	Common2.ApiResponse{}.Success(c, "状态切换成功", gin.H{})

}
