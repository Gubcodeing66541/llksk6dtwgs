package Admin

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"math"
	Common2 "server/App/Common"
	"server/App/Http/Request/Admin"
	"server/App/Http/Response"
	Common3 "server/App/Logic/Common"
	"server/App/Model/Common"
	"server/Base"
	"strings"
)

type Domain struct{}

func (Domain) List(c *gin.Context) {
	var pageReq Admin.DomainListLimit
	err := c.ShouldBind(&pageReq)
	if err != nil {
		Common2.ApiResponse{}.Error(c, "请提交完整的分页参数", gin.H{"err": err.Error()})
		return
	}

	tel := Base.MysqlConn.Table("domains").
		Select("domains.*,service_id,member,name,head,time_out").
		Joins("left JOIN services ON services.service_id = domains.bind_service_id")

	tel = tel.Where("domains.type = ? ", pageReq.Type)

	if pageReq.Domain != "" {
		tel = tel.Where("domain like ?", "%"+pageReq.Domain+"%")
	}

	if pageReq.Username != "" {
		tel = tel.Where("services.username like ? ", "%"+pageReq.Username+"%")
	}

	if pageReq.IsBindService == 1 {
		tel = tel.Where("bind_service_id != 0")
	}

	if pageReq.IsBindService == 2 {
		tel = tel.Where("bind_service_id = 0")
	}

	// 计算分页和总数
	var allCount int
	tel.Count(&allCount)
	allPage := math.Ceil(float64(allCount) / float64(pageReq.Offset))

	// 获取分页数据
	var list []Response.DomainRes
	tel.Offset((pageReq.Page - 1) * pageReq.Offset).Limit(pageReq.Offset).Scan(&list)
	res := gin.H{"count": allCount, "page": allPage, "current_page": pageReq.Page, "list": list}
	Common2.ApiResponse{}.Success(c, "获取成功", res)

}

func (Domain) QueryById(c *gin.Context) {
	var req Admin.QueryById
	err := c.ShouldBind(&req)
	if err != nil {
		Common2.ApiResponse{}.Error(c, "参数错误", gin.H{})
		return
	}
	domain := Common3.Domain{}.QueryById(req.DomainId)
	if domain.Id == 0 {
		Common2.ApiResponse{}.Error(c, "未查询到该域名", gin.H{})
		return
	}
	Common2.ApiResponse{}.Success(c, "id查询域名", gin.H{"domains": domain})
}

func (Domain) Delete(c *gin.Context) {
	var req Admin.DomainDelete
	err := c.ShouldBind(&req)
	if err != nil {
		Common2.ApiResponse{}.Error(c, "参数错误", gin.H{})
		return
	}

	err2 := Common3.Domain{}.Delete(req.DomainId)
	if err2 != nil {
		Common2.ApiResponse{}.Error(c, err2.Error(), gin.H{})
		return
	}
	Common2.ApiResponse{}.Success(c, "删除域名", gin.H{})
}

func (Domain) Update(c *gin.Context) {
	var req Admin.DomainUpdate
	err := c.ShouldBind(&req)
	if err != nil {
		Common2.ApiResponse{}.Error(c, "参数错误", gin.H{})
		return
	}

	err2 := Common3.Domain{}.Update(req.Id, strings.TrimSpace(req.Domain), req.Type, req.Status)
	if err2 != nil {
		Common2.ApiResponse{}.Error(c, err2.Error(), gin.H{})
		return
	}
	Common2.ApiResponse{}.Success(c, "修改域名", gin.H{})
}

func (Domain) Create(c *gin.Context) {
	var req Admin.DomainSave
	err := c.ShouldBind(&req)
	if err != nil {
		Common2.ApiResponse{}.Error(c, "参数错误", gin.H{})
		return
	}

	Common3.Domain{}.Create(strings.TrimSpace(req.Domain), req.TypeEd, "enable")
	Common2.ApiResponse{}.Success(c, "保存域名", gin.H{})
}

func (Domain) EnableDisable(c *gin.Context) {
	var req Admin.DomainEnableDisable
	err := c.ShouldBind(&req)
	if err != nil {
		Common2.ApiResponse{}.Error(c, "参数错误", gin.H{})
		return
	}

	Common3.Domain{}.EnableDisable(req.Id, req.Status)
	Common2.ApiResponse{}.Success(c, "状态已跟新", gin.H{})
}

func (Domain) UnBind(c *gin.Context) {
	var req Admin.DomainDelete
	err := c.ShouldBind(&req)
	if err != nil {
		Common2.ApiResponse{}.Error(c, "参数错误", gin.H{})
		return
	}

	DomainCount := Common3.Domain{}.GetNoUsePrivateNum()
	if DomainCount == 0 {
		Common2.ApiResponse{}.Error(c, "可用域名不足，无法解绑", gin.H{})
		return
	}

	domain := Common3.Domain{}.Get(req.DomainId)
	update := map[string]interface{}{"bind_service_id": 0, "status": "down"}
	Base.MysqlConn.Model(&Common.Domain{}).Where("id = ?", req.DomainId).Updates(update)
	_ = Common3.Domain{}.Bind(domain.BindServiceId)

	Common2.ApiResponse{}.Success(c, "域名已解绑并下架", gin.H{})

	param := fmt.Sprintf("?service_id=%d&type=%s&content=%s", domain.BindServiceId, "ban", "域名拦截提醒")
	Common2.Tools{}.HttpGet("http://127.0.0.1/api/socket/send_to_service_socket" + param)
}
