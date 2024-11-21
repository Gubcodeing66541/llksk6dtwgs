package Service

import (
	"github.com/gin-gonic/gin"
	"server/App/Common"
	"server/App/Http/Request"
	Common2 "server/App/Http/Request/Common"
	Service3 "server/App/Http/Request/Service"
	Service2 "server/App/Model/Service"
	"server/Base"
	"time"
)

type ServiceMessageReply struct{}

func (ServiceMessageReply) List(c *gin.Context) {
	var req Request.PageLimit
	if err := c.ShouldBind(&req); err != nil {
		Common.ApiResponse{}.Error(c, "参数错误", gin.H{})
		return
	}

	var cnt int64
	var list []Service2.ServiceMessageReply
	tel := Base.MysqlConn.Model(&list).Where("service_id = ?", Common.Tools{}.GetRoleId(c))
	tel.Find(&list).Count(&cnt)

	Common.ApiResponse{}.Success(c, "ok", gin.H{"list": list, "total": cnt})
}

func (ServiceMessageReply) Add(c *gin.Context) {
	var req Service3.ServiceMessageReplyReq
	if err := c.ShouldBind(&req); err != nil {
		Common.ApiResponse{}.Error(c, "参数错误", gin.H{})
		return
	}

	Base.MysqlConn.Create(&Service2.ServiceMessageReply{
		ServiceId: Common.Tools{}.GetRoleId(c),
		Type:      req.Type,
		MegType:   req.MegType,
		Title:     req.Title,
		Content:   req.Content,
		CreateAt:  time.Now().Unix(),
		UpdateAt:  time.Now().Unix(),
	})

	Common.ApiResponse{}.Success(c, "ok", gin.H{})
}

func (ServiceMessageReply) Edit(c *gin.Context) {
	var req Service3.ServiceMessageReplyReq
	if err := c.ShouldBind(&req); err != nil {
		Common.ApiResponse{}.Error(c, "参数错误", gin.H{})
		return
	}

	var info Service2.ServiceMessageReply
	Base.MysqlConn.Model(&info).Where("id = ? and service_id = ?", req.Id, Common.Tools{}.GetRoleId(c)).First(&info)
	if info.Id == 0 {
		Common.ApiResponse{}.Error(c, "参数错误", gin.H{})
		return
	}

	Base.MysqlConn.Model(&info).
		Where("id = ? and service_id = ?", req.Id, Common.Tools{}.GetRoleId(c)).Updates(gin.H{
		"type":      req.Type,
		"meg_type":  req.MegType,
		"title":     req.Title,
		"content":   req.Content,
		"update_at": time.Now().Unix(),
	})

	Common.ApiResponse{}.Success(c, "ok", gin.H{})
}

func (ServiceMessageReply) Del(c *gin.Context) {
	var req Common2.IdReq
	if err := c.ShouldBind(&req); err != nil {
		Common.ApiResponse{}.Error(c, "参数错误", gin.H{})
		return
	}

	var info Service2.ServiceMessageReply
	Base.MysqlConn.Model(&info).Where("id = ? and service_id = ?", req.Id, Common.Tools{}.GetRoleId(c)).First(&info)
	if info.Id == 0 {
		Common.ApiResponse{}.Error(c, "参数错误", gin.H{})
		return
	}

	Base.MysqlConn.Delete(&info)
	Common.ApiResponse{}.Success(c, "ok", gin.H{})
}
