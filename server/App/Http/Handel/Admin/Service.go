package Admin

import (
	"github.com/gin-gonic/gin"
	"server/App/Common"
	"server/App/Http/Request"
	"server/App/Http/Request/Admin"
	Common3 "server/App/Logic/Common"
	Service22 "server/App/Logic/Service"
	Common2 "server/App/Model/Common"
	Service2 "server/App/Model/Service"
	"server/Base"
	"strings"
	"time"
)

type Service struct{}

func (Service) Create(c *gin.Context) {
	var req Admin.ServiceCreateReq
	err := c.ShouldBind(&req)
	if err != nil {
		Common.ApiResponse{}.Error(c, "请填写需要创建的账号数量", gin.H{"err": err.Error()})
		return
	}

	var memberList []string
	for i := 0; i < req.Number; i++ {
		member := Service22.Service{}.Create(0, Common.Tools{}.CreateActiveMember(), req.Type, req.Day)
		if member.ServiceId == 0 {
			continue
		}

		// 创建订单
		Base.MysqlConn.Create(&Common2.Order{
			ServiceId:  member.ServiceId,
			Day:        req.Day,
			Price:      req.Price,
			Money:      req.Day * req.Price,
			Type:       "create",
			CreateTime: time.Now(),
		})
		memberList = append(memberList, member.Member)

		// 绑定域名
		_ = Common3.Domain{}.Bind(member.ServiceId)

	}

	Common.ApiResponse{}.Success(c, "ok", gin.H{"members": memberList})
}

func (Service) Renew(c *gin.Context) {
	var req Admin.ServiceRenewReq
	err := c.ShouldBind(&req)
	if err != nil {
		Common.ApiResponse{}.Error(c, "请填写完整", gin.H{})
		return
	}

	var serviceModel Service2.Service
	var ErrMember []string
	for _, member := range req.Member {
		member = strings.TrimSpace(member)
		_, err := Service22.Service{}.Renew(member, req.Day, req.Price)
		if err != nil {
			ErrMember = append(ErrMember, member)
		}

		// 绑定域名
		serviceModel, err = Service22.Service{}.MemberGet(member)
		if err == nil {
			_ = Common3.Domain{}.Bind(serviceModel.ServiceId)
		}

	}

	Common.ApiResponse{}.Success(c, "ok", gin.H{})
}

func (Service) List(c *gin.Context) {
	var req Request.PageLimit
	err := c.ShouldBind(&req)
	if err != nil {
		Common.ApiResponse{}.Error(c, "请填写完整", gin.H{})
		return
	}
	var model []Service2.Service
	tel := Base.MysqlConn.Model(&Service2.Service{})
	if req.Type != "" {
		tel = tel.Where("type = ?", req.Type)
	}
	if req.Member != "" {
		tel = tel.Where("member like ?", "%"+req.Member+"%")
	}
	var Count int
	tel.Count(&Count)
	Common.DbHelp{}.ModelByPage(tel, req.Limit, req.Page).Order("service_id desc").Find(&model)
	Common.ApiResponse{}.Success(c, "ok", gin.H{"service": model, "count": Count, "current_page": req.Page})
}

func (Service) Order(c *gin.Context) {
	var req Admin.OrderReq
	err := c.ShouldBind(&req)
	if err != nil {
		Common.ApiResponse{}.Error(c, "请填写完整", gin.H{})
		return
	}
}
