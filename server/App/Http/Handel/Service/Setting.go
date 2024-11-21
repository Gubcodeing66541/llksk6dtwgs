package Service

import (
	"github.com/gin-gonic/gin"
	"server/App/Common"
	Service2 "server/App/Http/Request/Service"
	Service4 "server/App/Model/Service"
	"server/Base"
)

type Setting struct{}

// BindDiversion 绑定
func (Setting) BindDiversion(c *gin.Context) {
	var req Service2.BindDiversionReq
	if err := c.ShouldBind(&req); err != nil {
		Common.ApiResponse{}.Error(c, req.GetError(err), gin.H{})
		return
	}

	// 获取服务信息
	var service Service4.Service
	Base.MysqlConn.Find(&service, "member = ?", req.Member)
	if req.Member != "" && service.ServiceId == 0 {
		Common.ApiResponse{}.Error(c, "搜索的客服激活码不存在，请输入正确的激活码", gin.H{})
		return
	}

	RoleId := Common.Tools{}.GetRoleId(c)
	if RoleId == service.ServiceId {
		Common.ApiResponse{}.Error(c, "不能绑定自己", gin.H{})
		return
	}

	Diversion := Common.Tools{}.If(req.Member == "", Common.Tools{}.CreateDiversionCode(int64(RoleId)), service.Diversion).(string)

	Base.MysqlConn.Model(&Service4.Service{}).Where("service_id = ?", RoleId).Update("diversion", Diversion)
	Common.ApiResponse{}.Success(c, "ok", gin.H{})
}
