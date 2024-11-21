package Service

import "github.com/go-playground/validator/v10"

type UserId struct {
	UserId int `json:"user_id" uri:"user_id" form:"user_id" `
}

type RoomTop struct {
	UserId int `json:"user_id" uri:"user_id" form:"user_id" `
	Top    int `json:"top" uri:"top" form:"top" `
}

type RoomBlack struct {
	Id      int    `json:"id"`
	UserId  int    `json:"user_id" uri:"user_id" form:"user_id" `
	IsBlack int    `json:"is_black" uri:"is_black" form:"is_black" `
	Day     int    `json:"day" uri:"day" form:"day" `
	Ip      string `json:"ip" uri:"ip" form:"ip" `
	Type    string `json:"type" uri:"type" form:"type" `
}

// ServiceRoomList all  所有
// server_read  已回复
// server_no_read 未回复
// server_no_read_count 未读
// top 置顶
// black 拉黑
type ServiceRoomList struct {
	IsClearStatus int    `json:"is_clear_status"  uri:"is_clear_status" form:"is_clear_status"`
	UserName      string `json:"user_name" uri:"user_name" form:"user_name"`
	UserNoRead    int    `json:"user_no_read" uri:"user_no_read" form:"user_no_read"`
	ServiceNoRead int    `json:"service_no_read" uri:"service_no_read" form:"service_no_read" `
	StartTime     string `json:"start_time" uri:"start_time" form:"start_time" `
	EndTime       string `json:"end_time" uri:"end_time" form:"end_time" `
	ServiceName   string `json:"service_name" uri:"service_name" form:"service_name" `
	ServiceMember string `json:"service_member" uri:"service_member" form:"service_member" `
	Type          string `json:"type" uri:"type" form:"type"`
	Top           int    `json:"is_top" uri:"is_top" form:"is_top" `
	Page          int    `json:"page" uri:"page" xml:"page" form:"page"`
	Offset        int    `json:"offset" uri:"offset" xml:"offset" form:"offset"`
}

type ServiceRoomDetail struct {
	Rename string `json:"rename" uri:"rename" form:"rename"`
	UserId int    `json:"user_id" uri:"user_id" form:"user_id"`
	Drive  string `json:"drive"   uri:"drive" form:"drive"`
	IP     string `json:"ip"   uri:"ip" form:"ip"`
	Map    string `json:"map"   uri:"map" form:"map"`
	Wechat string `json:"wechat"   uri:"wechat" form:"wechat"`
	Mobile string `json:"mobile"  uri:"mobile" form:"mobile"`
	Tag    string `json:"tag"  uri:"tag" form:"tag"`
}

func (ServiceRoomDetail) GetErr(err validator.ValidationErrors) string {
	for _, e := range err {
		if e.Field() == "Rename" {
			return "Rename"
		}
		if e.Field() == "UserId" {
			return "UserId"
		}
		if e.Field() == "Drive" {
			return "Drive"
		}
		if e.Field() == "IP" {
			return "IP"
		}
		if e.Field() == "Map" {
			return "Map"
		}
		if e.Field() == "Wechat" {
			return "Wechat"
		}
		if e.Field() == "Mobile" {
			return "Mobile"
		}
		if e.Field() == "Tag" {
			return "Tag"
		}
	}
	return "参数有误"
}

type UpdateServiceRoomRename struct {
	UserId int    `json:"user_id" uri:"user_id" form:"user_id" binding:"required" `
	Rename string `json:"rename" uri:"rename" form:"rename" binding:"required" `
}

type DelService struct {
	Id        int `json:"id" uri:"id" form:"id" binding:"required" `
	ServiceId int `json:"service_id" uri:"service_id" form:"service_id"  `
}
