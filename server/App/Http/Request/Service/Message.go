package Service

type SendToUser struct {
	To   int    `form:"to" json:"to" uri:"to" xml:"to" binding:"required"`
	Type string `form:"type" json:"type" uri:"type" xml:"type" binding:"required"`
	Info string `form:"info" json:"info" uri:"info" xml:"info" binding:"required"`
}

type MsgPage struct {
	Id int `form:"id" json:"id" uri:"id" xml:"id"`
}

type MsgUserId struct {
	UserId int `json:"user_id"`
}

type MsgList struct {
	ServiceId int `json:"service_id" uri:"service_id" xml:"service_id" form:"service_id"`
	UserId    int `json:"user_id" uri:"user_id" xml:"user_id" form:"user_id"`
	Page      int `json:"page" uri:"page" xml:"page" form:"page"`
	Offset    int `json:"offset" uri:"offset" xml:"offset" form:"offset"`
}

type ServiceSendMessage struct {
	UserId  int    `json:"user_id" form:"user_id" json:"user_id" uri:"user_id" xml:"user_id"`
	Type    string `json:"type" form:"type" json:"type" uri:"type" xml:"type"`
	Content string `json:"content" form:"content" json:"content" uri:"content" xml:"content"`
}

type ServiceSendMessageGroup struct {
	UserId  []int  `json:"user_id" form:"user_id" json:"user_id" uri:"user_id" xml:"user_id"`
	Type    string `json:"type" form:"type" json:"type" uri:"type" xml:"type"`
	Content string `json:"content" form:"content" json:"content" uri:"content" xml:"content"`
}

type UserSendMessage struct {
	Type           string `json:"type" form:"type" json:"type" uri:"type" xml:"type"`
	Content        string `json:"content" form:"content" json:"content" uri:"content" xml:"content"`
	ServiceContent string `json:"service_content" form:"service_content" uri:"service_content" xml:"service_content"`
}

type UpdateServiceDetail struct {
	Head        string `json:"head" form:"head" json:"head" uri:"head" xml:"head"`
	Name        string `json:"name" form:"name" json:"name" uri:"name" xml:"name"`
	UserDefault string `json:"user_default" form:"user_default" json:"user_default" uri:"user_default" xml:"user_default"`
}

func (u *UpdateServiceDetail) ToUpdate() map[string]string {
	return map[string]string{
		"head":         u.Head,
		"name":         u.Name,
		"user_default": u.UserDefault,
	}
}

type RemoveMsg struct {
	UserId int `json:"user_id" form:"user_id" json:"user_id" uri:"user_id" xml:"user_id"`
	Id     int `json:"id" form:"id" json:"id" uri:"id" xml:"id"`
}

type RemoveLateMsg struct {
	Type string `json:"type" form:"type" json:"type" uri:"type" xml:"type" binding:"required"` // week month all
}

func (*RemoveLateMsg) GerErr() string {
	return "type 参数错误"
}
