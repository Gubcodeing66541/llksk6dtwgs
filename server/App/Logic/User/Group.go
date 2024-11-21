package Logic

import (
	Group2 "server/App/Model/Group"
	User2 "server/App/Model/User"
	"server/Base"
	"time"
)

type Group struct{}

func (Group) Join(groupId int, user User2.User) Group2.GroupUser {
	var groupUser Group2.GroupUser
	Base.MysqlConn.Where("group_id = ? and user_id = ?", groupId, user.UserId).Find(&groupUser)
	if groupUser.Id == 0 {
		groupUser = Group2.GroupUser{
			GroupId: groupId, UserId: user.UserId, Rename: "",
			Name: "", Role: "user", Status: "none", CreateTime: time.Now(), UpdateTime: time.Now(),
		}
		Base.MysqlConn.Create(&groupUser)
	}
	return groupUser
}
