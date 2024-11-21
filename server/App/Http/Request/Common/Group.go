package Common

import "github.com/go-playground/validator/v10"

type GroupSocketMessage struct {
	GroupId    int    `json:"group_id"  form:"group_id" uri:"group_id" xml:"group_id"`
	Type       string `json:"type"  form:"type" uri:"type" xml:"type" binding:"required"`
	Content    string `json:"content"  form:"content" uri:"content" xml:"content" binding:"required"`
	IsJustPush int    `json:"is_push"  form:"is_push" uri:"is_push" xml:"is_push"` // 是否只推送
}

func (GroupSocketMessage) GetErr(err validator.ValidationErrors) string {
	for _, e := range err {
		if e.Field() == "GroupId" {
			return "请选择组"
		}
		if e.Field() == "Type" {
			return "请选择类型"
		}
		if e.Field() == "Content" {
			return "请输入发送的消息"
		}
	}
	return "请求错误"
}
