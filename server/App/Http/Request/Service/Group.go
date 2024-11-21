package Service

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"server/App/Http/Request"
)

type GroupIdReq struct {
	GroupId int `json:"group_id" binding:"required,number"`
}

func (GroupIdReq) GetErr(err validator.ValidationErrors) string {
	for _, e := range err {
		if e.Field() == "GroupId" {
			return "请确认选择的组"
		}
	}
	return "参数有误"
}

type CreateGroupReq struct {
	GroupName string `json:"group_name" binding:"required"`
	GroupHead string `json:"group_head" binding:"required"`
}

func (CreateGroupReq) GetErr(err validator.ValidationErrors) string {
	for _, e := range err {
		if e.Field() == "GroupName" {
			return "请输入名称"
		}
		if e.Field() == "group_head" {
			return "请选择头像"
		}
	}
	return "参数有误"
}

type UpdateGroupUserReq struct {
	UserId int    `json:"user_id" binding:"required"`
	Rename string `json:"rename"`
	Role   string `json:"role" binding:"required"`   // user用户 manager管理员
	Status string `json:"status" binding:"required"` // none无状态 stop禁言 ban禁止登入
}

func (req *UpdateGroupUserReq) ToMap() map[string]interface{} {
	return map[string]interface{}{"rename": req.Rename, "role": req.Role, "status": req.Status}
}

func (UpdateGroupUserReq) GetErr(err validator.ValidationErrors) string {
	for _, e := range err {
		if e.Field() == "UserId" {
			return "请选择用户"
		}
		if e.Field() == "Role" {
			if e.Value().(string) != "user" && e.Value().(string) != "manager" {
				return "只允许user和manager两种角色"
			}
			return "请选择角色"
		}
		if e.Field() == "Status" {
			if e.Value().(string) != "none" && e.Value().(string) != "stop" {
				return "只允许none和stop两种状态"
			}
			return "请确定发言状态"
		}
	}
	return "参数有误"
}

type GroupUpdateReq struct {
	GroupId   int    `json:"group_id" binding:"required,number"`
	GroupName string `json:"group_name" binding:"required"`
	GroupHead string `json:"group_head" binding:"required"`
	Status    string `json:"status" binding:"required"`
	Notice    string `json:"notice"`
	Hello     string `json:"hello"`
}

func (GroupUpdateReq) GetErr(err validator.ValidationErrors) string {
	for _, e := range err {
		if e.Field() == "GroupId" {
			return "请确认选择的组"
		}
		if e.Field() == "GroupName" {
			return "请输入名称"
		}
		if e.Field() == "GroupHead" {
			return "请填写头像"
		}
		if e.Field() == "Status" {
			return "请填写状态"
		}
		return fmt.Sprintf(e.Error())
	}
	return "参数有误"
}

type GroupMessageReq struct {
	Request.PageLimit
}

func (GroupMessageReq) GetErr(err validator.ValidationErrors) string {
	for _, e := range err {
		return fmt.Sprintf(e.Error())
	}
	return "参数有误"
}

type GroupUserIdReq struct {
	UserId int `json:"user_id" binding:"required,number"`
}

func (GroupUserIdReq) GetErr(err validator.ValidationErrors) string {
	for _, e := range err {
		if e.Field() == "UserId" {
			return "请确认选择的组"
		}
	}
	return "参数有误"
}

type UpdateGroupReq struct {
	GroupName string `json:"group_name"  binding:"required"`
	GroupHead string `json:"group_head" binding:"required"`
	Status    string `json:"status" binding:"required"` // none无状态 stop禁言
	Notice    string `json:"notice"`                    // 公告
	Hello     string `json:"hello"`                     // 打招呼
}

func (UpdateGroupReq) GetErr(err validator.ValidationErrors) string {
	for _, e := range err {
		if e.Field() == "GroupName" {
			return "GroupName"
		}
		if e.Field() == "GroupHead" {
			return "GroupHead"
		}
		if e.Field() == "Status" {
			if e.Value().(string) != "none" && e.Value().(string) != "stop" {
				return "只允许none和stop两种状态"
			}
			return "请确定发言状态"
		}
		if e.Field() == "Notice" {
			return "Notice"
		}
		if e.Field() == "Hello" {
			return "Hello"
		}
	}
	return "参数有误"
}

type GroupUserSearchReq struct {
	UserName string `json:"user_name"`
}

func (GroupUserSearchReq) GetErr(err validator.ValidationErrors) string {
	for _, e := range err {
		if e.Field() == "UserName" {
			return "请确认选择UserName"
		}
		return fmt.Sprintf(e.Error())
	}
	return "参数有误"
}

type GroupDeleteMsgReq struct {
	Id int `json:"id" binding:"required,number"`
}

func (GroupDeleteMsgReq) GetErr(err validator.ValidationErrors) string {
	for _, e := range err {
		if e.Field() == "Id" {
			return "请选择消息"
		}
	}
	return "参数有误"
}
