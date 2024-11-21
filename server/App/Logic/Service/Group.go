package Service

import (
	"errors"
	Service3 "server/App/Http/Request/Service"
	Group2 "server/App/Model/Group"
	"server/Base"
)

type Group struct{}

func (Group) UpdateUser(serviceId int, req Service3.UpdateGroupUserReq) error {
	// 检测自己是否有权限
	var group Group2.Group
	Base.MysqlConn.Find(&group, "service_id = ?", serviceId)
	if group.GroupId == 0 {
		return errors.New("无操作权限")
	}
	Base.MysqlConn.Model(&Group2.GroupUser{}).Where("service_id = ? and user_id = ?", serviceId, req.UserId).Updates(req.ToMap())
	return nil
}
