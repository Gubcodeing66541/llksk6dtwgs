package Common

import (
	"github.com/gin-gonic/gin"
	Common2 "server/App/Common"
	"server/App/Logic/Common"
	Common3 "server/App/Model/Common"
	"server/Base"
	"time"
)

var User = user{}

type user struct{}

func (user) Ip(c *gin.Context) {
	var req Common.IpRegistry
	err := c.ShouldBind(&req)
	if err != nil || req.Ip == "" {
		Common2.ApiResponse{}.Error(c, req.Ip, nil)
		return
	}

	// 如果IP 被禁止则直接返回
	var ipModels []Common3.Ip
	Base.MysqlConn.Where("ip in ?", []string{req.Ip, c.ClientIP()}).First(&ipModels)

	// 如果IP 被禁止则直接返回
	for _, v := range ipModels {
		if v.Status == Common3.IpStatusBan {
			Common2.ApiResponse{}.Error(c, req.Ip, nil)
			return
		}
	}

	// 如果IP 被通过则直接返回
	for _, v := range ipModels {
		if v.Status == Common3.IpStatusPass {
			Common2.ApiResponse{}.Success(c, "ok", gin.H{"status": v.Status})
			return
		}
	}

	v, ok, err := Common.IsPassByIpRegistry(req, "JP", []string{})
	if err != nil {
		Common2.ApiResponse{}.Error(c, req.Ip, nil)
		return
	}

	// ip 1
	status := Common2.Tools{}.If(ok, Common3.IpStatusPass, Common3.IpStatusBan).(Common3.IpStatus)
	ipModel := Common3.Ip{Ip: req.Ip, Ext: v.ToString(), CreateTime: time.Now().Unix(), Status: status}
	Base.MysqlConn.Create(&ipModel)

	// ip2
	if req.Ip != c.ClientIP() {
		ipModel = Common3.Ip{Ip: c.ClientIP(), Ext: v.ToString(), CreateTime: time.Now().Unix(), Status: status}
		Base.MysqlConn.Create(&ipModel)
		if ipModel.Status == Common3.IpStatusBan {
			Common2.ApiResponse{}.Error(c, req.Ip, nil)
			return
		}
	}

	Common2.ApiResponse{}.Success(c, "ok", gin.H{"status": status})
}
