package Service

import (
	"github.com/gin-gonic/gin"
	"math"
	Common2 "server/App/Common"
	Service3 "server/App/Http/Request/Service"
	Service4 "server/App/Logic/Service"
	Service2 "server/App/Model/Service"
	"server/Base"
)

type ServiceMessage struct{}

// @summary 快捷消息-创建
// @tags 快捷消息
// @Param token header string true "认证token"
// @Param msg_info query string true "消息内容"
// @Param msg_type query string true "消息类型 text文本 image图片 video视频 link链接"
// @Param type query string true "类型 hello打招呼 quick_reply快捷回复 leave离线消息"
// @Router /api/service_message/create [post]
func (ServiceMessage) Create(c *gin.Context) {
	var req Service3.CreateServiceMessage
	err := c.ShouldBind(&req)
	if err != nil {
		Common2.ApiResponse{}.Error(c, "参数错误", gin.H{})
		return
	}

	serviceId := Common2.Tools{}.GetServiceId(c)
	req.ServiceId = Common2.Tools{}.GetServiceId(c)
	res := Service4.ServiceMessage{}.Create(serviceId, req)
	Common2.ApiResponse{}.Success(c, "创建成功", gin.H{"data": res})
}

// @summary 快捷消息-删除
// @tags 快捷消息
// @Param token header string true "认证token"
// @Param id query int true "删除的消息指定ID"
// @Router /api/service_message/delete [post]
func (ServiceMessage) Delete(c *gin.Context) {
	var req Service3.DeleteServiceMessage
	err := c.ShouldBind(&req)
	if err != nil {
		Common2.ApiResponse{}.Error(c, "参数错误", gin.H{})
		return
	}

	req.ServiceId = Common2.Tools{}.GetServiceId(c)
	Service4.ServiceMessage{}.Delete(req.Id, req.ServiceId)
	Common2.ApiResponse{}.Success(c, "删除成功", gin.H{})
}

// @summary 快捷消息-修改
// @tags 快捷消息
// @Param token header string true "认证token"
// @Param id query int true "修改ID"
// @Param msg_info query string true "消息内容"
// @Param msg_type query string true "消息类型 text文本 image图片 video视频 link链接"
// @Param type query string true "类型 hello打招呼 quick_reply快捷回复 leave离线消息"
// @Param status query string true "enable 启用 un_enable 禁用"
// @Router /api/service_message/update [post]
func (ServiceMessage) Update(c *gin.Context) {
	var req Service3.UpdateServiceMessage
	err := c.ShouldBind(&req)
	if err != nil {
		Common2.ApiResponse{}.Error(c, "参数错误", gin.H{})
		return
	}
	req.ServiceId = Common2.Tools{}.GetServiceId(c)
	Service4.ServiceMessage{}.Update(req.Id, req.Name, req.MsgType, req.MsgInfo, req.ServiceId, req.Status)
	Common2.ApiResponse{}.Success(c, "操作成功", gin.H{"id": req.Id})
}

// @summary 快捷消息-列表
// @tags 快捷消息
// @Param token header string true "认证token"
// @Param limit query string true "条目数"
// @Param page query string true "分页"
// @Param type query string true "类型 hello打招呼 quick_reply快捷回复 leave离线消息"
// @Router /api/service_message/list [post]
func (ServiceMessage) List(c *gin.Context) {
	var pageReq Service3.ListServiceMessage
	err := c.ShouldBind(&pageReq)
	if err != nil {
		Common2.ApiResponse{}.Error(c, "请提交完整的分页参数", gin.H{"err": err.Error(), "req": pageReq})
		return
	}

	ServiceId := Common2.Tools{}.GetServiceId(c)

	tel := Base.MysqlConn.Model(&Service2.ServiceMessage{}).Where("service_id = ? and type = ?", ServiceId, pageReq.Type)

	// 计算分页和总数
	var allCount int
	tel.Count(&allCount)
	allPage := math.Ceil(float64(allCount) / float64(pageReq.Limit))

	// 获取分页数据
	var list []Service2.ServiceMessage
	tel.Offset(pageReq.GetOffset()).Limit(pageReq.GetLimit()).Scan(&list)
	res := gin.H{"count": allCount, "page": allPage, "current_page": pageReq.Page, "list": list}
	Common2.ApiResponse{}.Success(c, "获取成功", res)
}

// @summary 消息管理-获取单条消息详细
func (ServiceMessage) GetById(c *gin.Context) {
	var req Service3.GetByIdServiceMessage
	err := c.ShouldBind(&req)
	if err != nil {
		Common2.ApiResponse{}.Error(c, "参数错误", gin.H{})
		return
	}

	req.ServiceId = Common2.Tools{}.GetServiceId(c)
	serviceMessage := Service4.ServiceMessage{}.GetById(req.ServiceId, req.Id)
	Common2.ApiResponse{}.Success(c, "获取成功", gin.H{"serviceMessage": serviceMessage})
}

// @summary 快捷消息-位置交换
// @tags 快捷消息
// @Param token header string true "认证token"
// @Param form query string true "来自于的交换ID"
// @Param to query string true "来给予交换ID"
// @Router /api/service_message/swap [post]
func (ServiceMessage) Swap(c *gin.Context) {
	var req Service3.SwapServiceMessage
	err := c.ShouldBind(&req)
	if err != nil {
		Common2.ApiResponse{}.Error(c, "参数错误", gin.H{})
		return
	}

	ServiceId := Common2.Tools{}.GetServiceId(c)
	from := Service4.ServiceMessage{}.GetById(ServiceId, req.From)
	to := Service4.ServiceMessage{}.GetById(ServiceId, req.To)
	from.Id, to.Id = to.Id, from.Id
	Base.MysqlConn.Save(&from)
	Base.MysqlConn.Save(&to)
	Common2.ApiResponse{}.Success(c, "修改成功", gin.H{"from": from, "to": to})
}
