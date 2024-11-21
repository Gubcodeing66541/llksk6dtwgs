package Logic

import (
	"fmt"
	"server/App/Common"
	Common2 "server/App/Logic/Common"
	Service2 "server/App/Model/Service"
	User2 "server/App/Model/User"
	"server/Base"
	"time"
)

type User struct{}

func (User) OpenIdToUser(openId string) User2.User {
	var user User2.User
	Base.MysqlConn.Find(&user, "open_id = ?", openId)
	return user
}

func (User) UnionIdToUser(unionId string) User2.User {
	var user User2.User
	Base.MysqlConn.Find(&user, "union_id = ?", unionId)
	return user
}

func (User) UserIdToUser(userId int) User2.User {
	var user User2.User
	Base.MysqlConn.Find(&user, "user_id = ?", userId)
	return user
}

// CookieUUIDToUser 通过系统生成的uuid获取用户信息
func (User) CookieUUIDToUser(cookieUuid string) User2.User {
	var user User2.User
	var userMap User2.UserAuthMap
	Base.MysqlConn.Find(&userMap, "cookie_uid = ?", cookieUuid)
	Base.MysqlConn.Find(&user, "user_Id = ?", userMap.UserId)
	return user
}

// CheckCookieUUIDToUser 通过系统生成的uuid获取用户信息
func (User) CheckCookieUUIDToUser(userMap User2.User, JoinUuid string) User2.User {
	var userModel User2.User
	if userMap.UserId != 0 {
		userModel = User{}.UserIdToUser(userMap.UserId)
	} else {
		username := fmt.Sprintf("%s", Common.Tools{}.GetRename())
		userModel = User{}.CreateUser("", username, Common.Tools{}.GetHead(), 0, "")
		Base.MysqlConn.Create(User2.UserAuthMap{CookieUid: JoinUuid, UserId: userModel.UserId, Action: "show"})
	}
	return userModel
}

func (User) CreateUser(openId string, Name string, Header string, Sex int, unionId string) User2.User {
	now := time.Now()
	user := User2.User{OpenId: openId, UserName: Name, UserHead: Header, CreateTime: now, UpdateTime: now, UnionId: unionId}
	Base.MysqlConn.Create(&user)
	return user
}

func (User) CreateWebUser(Name string, Header string, Token string) User2.User {
	now := time.Now()
	user := User2.User{UserName: Name, UserHead: Header, CreateTime: now, UpdateTime: now, Token: Token}
	Base.MysqlConn.Create(&user)
	return user
}

func (User) HandelLeaveMessage(serviceId int, userId int) {
	// 如果客服在线则不用管了
	ServiceIsOnline := Base.WebsocketHub.UserIdIsOnline(Common.Tools{}.GetServiceWebSocketId(serviceId))
	if ServiceIsOnline == 1 {
		return
	}

	var leaveMsg []Service2.ServiceMessage
	Base.MysqlConn.Find(&leaveMsg, "service_id = ? and type = 'leave'", serviceId)

	for key, v := range leaveMsg {
		if v.Status == "enable" {
			_ = Common2.Message{}.SendToUser(serviceId, userId, leaveMsg[key].MsgType, leaveMsg[key].MsgInfo, true)
		}
	}
}
