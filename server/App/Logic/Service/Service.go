package Service

import (
	"errors"
	"math/rand"
	"server/App/Common"
	Common2 "server/App/Model/Common"
	Group2 "server/App/Model/Group"
	Service2 "server/App/Model/Service"
	"server/Base"
	"time"
)

type Service struct{}

func (Service) Create(agentId int64, Member string, Type string, Day int) Service2.Service {
	addTime := time.Now()
	if Day == 0 {
		addTime = addTime.Add(60 * 30 * time.Second)
	} else {
		addTime = addTime.AddDate(0, 0, Day)
	}
	service := Service2.Service{
		Member:           Member,
		AgentId:          agentId,
		Name:             "lAIMI",
		Type:             Type,
		Head:             Common.Tools{}.GetDefaultHead(),
		TimeOut:          addTime,
		CreateTime:       time.Now(),
		FirstLoginStatus: 0,
		UpdateTime:       time.Now(),
		UserDefault:      "default",
	}
	Base.MysqlConn.Create(&service)
	service.Code = Common.Tools{}.CreateActiveCode(service.ServiceId)
	service.Diversion = Common.Tools{}.CreateDiversionCode(int64(service.ServiceId))
	Base.MysqlConn.Save(&service)

	// 如过是群聊则创建群组
	if Type == "group" {
		Base.MysqlConn.Create(&Group2.Group{
			ServiceId:  service.ServiceId,
			GroupName:  "群组",
			GroupHead:  Common.Tools{}.GetDefaultHead(),
			Status:     "none",
			Code:       "G" + Common.Tools{}.CreateActiveCode(service.ServiceId),
			Notice:     "",
			Hello:      "",
			NoReadCnt:  0,
			LateTime:   time.Now(),
			CreateTime: time.Now(),
			UpdateTime: time.Now(),
		})
	}

	return service
}

func (Service) Renew(Member string, Day int, Price int) (Service2.Service, error) {
	var service Service2.Service
	Base.MysqlConn.Find(&service, "member = ?", Member)
	if service.ServiceId == 0 {
		return service, errors.New("账号不存在")
	}

	if service.TimeOut.After(time.Now()) {
		service.TimeOut = service.TimeOut.AddDate(0, 0, Day)
	} else {
		service.TimeOut = time.Now().AddDate(0, 0, Day)
	}
	service.UpdateTime = time.Now()
	Base.MysqlConn.Save(&service)
	Base.MysqlConn.Create(&Common2.Order{ServiceId: service.ServiceId, Day: Day, Price: Price, Money: Price * Day, CreateTime: time.Now(), Type: "renew"})
	return service, nil
}

func (Service) Get(code string) (Service2.Service, error) {
	var service Service2.Service
	Base.MysqlConn.Find(&service, "code = ?", code)
	if service.ServiceId == 0 {
		return service, errors.New("客服不存在")
	}
	return service, nil
}

func (s Service) IdGet(serviceId int) (Service2.Service, error) {
	var model Service2.Service
	Base.MysqlConn.Find(&model, "service_id = ?", serviceId)
	return model, nil
}

func (s Service) MemberGet(serviceId string) (Service2.Service, error) {
	var model Service2.Service
	Base.MysqlConn.Find(&model, "member = ?", serviceId)
	return model, nil
}

// GetByDiversion 通过分流码获取客服
func (Service) GetByDiversion(code string, userId int) (Service2.Service, error) {

	// 检测是否已经分配过客服
	var serviceRoom Service2.ServiceRoom
	Base.MysqlConn.Find(&serviceRoom, "user_id = ? and diversion = ?", userId, code)
	if serviceRoom.ServiceId != 0 {
		var service Service2.Service
		Base.MysqlConn.Find(&service, "service_id = ?", serviceRoom.ServiceId)
		return service, nil
	}

	// 获取分流码对应的客服
	var service []Service2.Service
	Base.MysqlConn.Find(&service, "diversion = ?", code)
	var action Service2.Service

	if len(service) == 0 {
		return action, errors.New("客服不存在")
	}

	// 寻找在线并且未过期的客服
	var onlineService []Service2.Service
	for _, v := range service {
		isOnline := Base.WebsocketHub.UserIdIsOnline(Common.Tools{}.GetServiceWebSocketId(v.ServiceId))
		if v.TimeOut.After(time.Now()) && isOnline == 1 {
			onlineService = append(onlineService, v)
		}
	}

	// 如果有在线的客服，就随机分配一个
	if len(onlineService) >= 1 {
		if len(onlineService) == 1 {
			action = onlineService[0]
			return action, nil
		}
		action = onlineService[rand.Intn(len(onlineService))]
	} else {
		// 如果没有在线的客服，就随机分配一个
		if len(service) == 1 {
			action = service[0]
			return action, nil
		}
		action = service[rand.Intn(len(service))]
	}

	return action, nil
}
