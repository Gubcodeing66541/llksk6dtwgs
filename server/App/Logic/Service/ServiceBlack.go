package Service

import (
	"github.com/gin-gonic/gin"
	Service2 "server/App/Model/Service"
	"server/Base"
	"time"
)

type ServiceBlack struct{}

// 如果是ip拉黑吧所有用户加入到黑名单中
type users struct {
	UserId int `json:"user_id"`
}

// Ip拉黑
func (ServiceBlack) IpBlack(IsBlack int, serviceId int, Ip string) {
	// 取消拉黑操作
	if IsBlack == 0 {
		var blackList []Service2.ServiceBlack
		Base.MysqlConn.Find(&blackList, "service_id = ? and type = 'ip' and ip = ?", serviceId, Ip)
		Base.MysqlConn.Delete(&Service2.ServiceBlack{}, "service_id = ? and type = 'ip' and ip = ?", serviceId, Ip)
		for key, _ := range blackList {
			Base.MysqlConn.Model(&Service2.ServiceRoom{}).Where("service_id = ? and user_id = ?", serviceId, blackList[key].UserId).Update(gin.H{"is_black": 0})
		}
		return
	}

	//拉黑操作
	var IpToUsers []users
	Base.MysqlConn.Raw("select user_id from service_rooms where service_id = ? and late_ip = ? and user_id not in ( select user_id from service_blacks where service_id = ?)", serviceId, Ip, serviceId).Scan(&IpToUsers)
	for key, _ := range IpToUsers {
		Base.MysqlConn.Create(&Service2.ServiceBlack{
			ServiceId: serviceId, Type: "ip",
			UserId:  IpToUsers[key].UserId,
			TimeOut: time.Now(),
			Ip:      Ip, CreateTime: time.Now(),
		})
		Base.MysqlConn.Model(&Service2.ServiceRoom{}).
			Where("service_id = ? and user_id = ?", serviceId, IpToUsers[key].UserId).
			Update("is_black", 1)
	}
}

func (ServiceBlack) UserBlack(IsBlack int, serviceId int, Ip string, UserId int) {
	if IsBlack == 1 {
		Base.MysqlConn.Create(&Service2.ServiceBlack{
			TimeOut:   time.Now(),
			ServiceId: serviceId, Type: "user", UserId: UserId, Ip: Ip, CreateTime: time.Now(),
		})
		Base.MysqlConn.Model(&Service2.ServiceRoom{}).
			Where("service_id = ? and user_id = ?", serviceId, UserId).
			Update("is_black", 1)
		return
	}
	Base.MysqlConn.Delete(&Service2.ServiceBlack{}, "service_id = ? and type = 'user' and user_id = ?", serviceId, UserId)
	Base.MysqlConn.Model(&Service2.ServiceRoom{}).
		Where("service_id = ? and user_id = ?", serviceId, UserId).
		Update("is_black", 0)
}

// Ip拉黑
func (ServiceBlack) ServiceManagerIpBlack(roleId int, Ip string, Day int, serviceManagerId int) {
	//拉黑操作
	Base.MysqlConn.Create(&Service2.ServiceBlack{
		ServiceId: 0, Type: "ip", UserId: 0, Day: Day, TimeOut: time.Now().AddDate(0, 0, Day),
		Ip: Ip, CreateTime: time.Now(),
	})

	var services []Service2.Service
	Base.MysqlConn.Find(&services, "service_manager_id=?", roleId)

	//循环新增所有拉黑的用户
	for _, v := range services {
		serviceId := v.ServiceId

		//查询所有的serviceRoom
		var serviceRooms []Service2.ServiceRoom
		Base.MysqlConn.Find(&serviceRooms, "service_id=? and late_ip=?", serviceId, Ip)

		for _, serviceRoom := range serviceRooms {
			Base.MysqlConn.Create(&Service2.ServiceBlack{
				ServiceId: serviceId, Type: "ip", UserId: serviceRoom.UserId, Day: Day, TimeOut: time.Now().AddDate(0, 0, Day),
				Ip: Ip, CreateTime: time.Now(),
			})
		}

	}

	//循环拉黑所有该ip的用户
	for _, v := range services {
		serviceId := v.ServiceId
		Base.MysqlConn.Model(&Service2.ServiceRoom{}).Where("service_id=? and late_ip=?", serviceId, Ip).Updates(&Service2.ServiceRoom{IsBlack: 1})

	}

}

func (ServiceBlack) ServiceManagerUserBlack(UserId int, serviceManagerId int, Day int) {
	Base.MysqlConn.Create(&Service2.ServiceBlack{
		ServiceId: 0, Type: "user", UserId: UserId, Ip: "Ip", CreateTime: time.Now(),
		Day: Day, TimeOut: time.Now().AddDate(0, 0, Day),
	})
}
