package Service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"math"
	"server/App/Common"
	Service3 "server/App/Http/Request/Service"
	"server/App/Http/Response"
	Common2 "server/App/Model/Common"
	Service2 "server/App/Model/Service"
	User2 "server/App/Model/User"
	"server/Base"
	"time"
)

type ServiceRoom struct {
}

func (ServiceRoom) Get(user User2.User, service Service2.Service, Ip string, Drive string, DriveInfo string, maps string) Service2.ServiceRoom {
	var serverRoom Service2.ServiceRoom
	Base.MysqlConn.Find(&serverRoom, "service_id = ? and user_id = ?", service.ServiceId, user.UserId)

	if serverRoom.Id == 0 {
		now := time.Now()
		var hello []Service2.ServiceMessage
		Base.MysqlConn.Find(&hello, "service_id = ? and type = ? and status = 'enable'", service.ServiceId, "hello")
		LateMsg, LateType, lateId := "", "", 0
		for _, item := range hello {
			model := &Common2.Message{
				From: service.ServiceId, To: user.UserId, Type: item.MsgType, Content: item.MsgInfo, Time: time.Now().Unix(),
				SendRole: "hello", CreateTime: now, IsRead: 1, UserId: user.UserId, ServiceId: service.ServiceId}
			Base.MysqlConn.Create(&model)
			LateType, LateMsg = item.MsgType, item.MsgInfo
			lateId = model.Id
		}
		serverRoom = Service2.ServiceRoom{
			ServiceId: service.ServiceId, LateUserReadId: lateId, UserId: user.UserId, LateType: LateType,
			LateMsg: LateMsg, CreateTime: now, UpdateTime: now, Times: time.Now().Unix(), LateId: lateId, Diversion: service.Diversion}
		Base.MysqlConn.Create(&serverRoom)
		Base.MysqlConn.Create(&Service2.ServiceRoomDetail{
			ServiceId: service.ServiceId, UserId: user.UserId, IP: Ip, CreateTime: now, Drive: Drive, DriveInfo: DriveInfo})
	} else {
		fmt.Println("update service room drive_info", DriveInfo, "drive", Drive)
		Base.MysqlConn.Model(&Service2.ServiceRoomDetail{}).Where("user_id = ? and service_id = ?", user.UserId, service.ServiceId).Updates(gin.H{
			"drive": Drive, "drive_info": DriveInfo, "ip": Ip})
	}

	go func(Ip string) {
		userAddr, _ := Common.Tools{}.IPInfo(Ip)
		if userAddr.City != "" {
			Base.MysqlConn.Model(&Service2.ServiceRoomDetail{}).
				Where("service_id = ? and user_id = ?", service.ServiceId, user.UserId).Updates(gin.H{"map": userAddr.Addr})
		}
	}(Ip)
	return serverRoom
}

func (ServiceRoom) GetByUserId(serviceId, userId int) Service2.ServiceRoom {
	var serviceRoom Service2.ServiceRoom
	Base.MysqlConn.Find(&serviceRoom, "service_id = ? and user_id = ?", serviceId, userId)
	return serviceRoom
}

func (ServiceRoom) List(serviceId int, req Service3.ServiceRoomList, isRoomsSearch bool) gin.H {
	var model []Response.ServiceRoom

	where := fmt.Sprintf("service_rooms.service_id = %d ", serviceId)

	if isRoomsSearch {
		where += fmt.Sprintf(" and is_delete = 0 ")
	}

	if req.UserName != "" {
		where += fmt.Sprintf(" and (users.user_name like '%%%s%%' or service_rooms.rename like '%%%s%%')", req.UserName, req.UserName)
	}

	if req.Type == "server_read" {
		where += fmt.Sprintf(" and service_rooms.late_role = 'service' ")
	}

	if req.Type == "server_no_read" {
		where += fmt.Sprintf(" and service_rooms.late_role = 'user' ")
	}

	if req.Type == "server_no_read_count" {
		where += fmt.Sprintf(" and service_rooms.service_no_read > 0 ")
	}

	if req.Type == "top" {
		where += fmt.Sprintf(" and service_rooms.is_top = 1 ")
	}

	if req.Type == "group" {
		where += fmt.Sprintf(" and service_rooms.type = 'group' ")
	}

	if req.Type == "black" {
		where += fmt.Sprintf(" and service_rooms.is_black = %d ", 1)
	} else {
		where += fmt.Sprintf(" and service_rooms.is_black = %d ", 0)
	}

	show := "service_rooms.id,service_rooms.late_type,service_rooms.is_black,service_rooms.is_top,service_rooms.user_no_read,service_rooms.late_msg,service_rooms.service_no_read,service_rooms.update_time,service_rooms.rename,"
	show += "users.user_id,users.user_name,users.user_head"
	join := "join users on service_rooms.user_id = users.user_id"

	tel := Base.MysqlConn.Table("service_rooms").Select(show).Where(where).Joins(join)

	// 计算分页和总数
	var allCount int
	tel.Count(&allCount)
	allPage := int(math.Ceil(float64(allCount) / float64(req.Offset)))

	// 获取分页数据
	tel.Offset((req.Page - 1) * req.Offset).Order("service_rooms.update_time desc").Limit(req.Offset).Scan(&model)

	for key, item := range model {
		userName := fmt.Sprintf("%s:%d", "user", item.UserId)
		model[key].IsOnline = Base.WebsocketHub.UserIdIsOnline(userName)
	}

	return gin.H{"count": allCount, "page": allPage, "current_page": req.Page, "list": model}
}

func (ServiceRoom) ListByServiceManager(ServiceManagerId int, req Service3.ServiceRoomList) gin.H {
	var model []Response.ServiceRoomList
	timeWhere := " where 1= 1"
	if req.StartTime != "" {
		timeWhere = fmt.Sprintf(" and users.create_time >= '%s'", req.StartTime)
	}

	if req.EndTime != "" {
		timeWhere += fmt.Sprintf(" and users.create_time <= '%s'", req.EndTime)
	}

	if req.UserName != "" {
		timeWhere += fmt.Sprintf(" and users.user_name = '%s'", req.UserName)
	}

	if req.ServiceName != "" {
		timeWhere += fmt.Sprintf(" and services.name = '%s'", req.ServiceName)
	}

	if req.ServiceMember != "" {
		timeWhere += fmt.Sprintf(" and services.username = '%s'", req.ServiceMember)
	}

	sql := `
			SELECT  user_name,user_head,
                         service_room_details.*,services.name,services.head as service_head
			FROM  users
				  left join service_room_details on users.user_id =  service_room_details.user_id
				  left join services on service_room_details.service_id = services.service_id
 order by create_time desc
		`
	if ServiceManagerId != 0 {
		sql = fmt.Sprintf(`
			SELECT  user_name,user_head,
                         service_room_details.*,services.name,services.head as service_head
			FROM  (
					  select * from users where user_id in (
						  select * from (
							   select user_id from service_rooms where service_id in (
								   select service_id from services where service_manager_id = %d
							   )
						  ) a
					  )
				  ) as users
				  left join service_room_details on users.user_id =  service_room_details.user_id
				  left join services on service_room_details.service_id = services.service_id 
		`, ServiceManagerId)
	}

	sql = sql + timeWhere
	tel := Base.MysqlConn.Raw(sql)

	// 计算分页和总数
	var allCount int
	telCount := Base.MysqlConn.Table("users")
	if ServiceManagerId != 0 {
		telCount = telCount.Where(`
		user_id in (
			select user_id from service_rooms where service_id in (
				select service_id from services where service_manager_id = ?
			)
		)
	`, ServiceManagerId)
	}
	type CountPage struct {
		Cnt int
	}
	var CountPageCnt CountPage
	sql = fmt.Sprintf("select count(*) as cnt from (%s) tt", sql)
	Base.MysqlConn.Raw(sql).Scan(&CountPageCnt)
	allCount = CountPageCnt.Cnt
	allPage := math.Ceil(float64(allCount) / float64(req.Offset))

	// 获取分页数据
	tel.Offset((req.Page - 1) * req.Offset).Limit(req.Offset).Scan(&model)

	for key, item := range model {
		model[key].IsOnline = Base.WebsocketHub.UserIdIsOnline(Common.Tools{}.GetUserWebSocketId(item.UserId))
	}

	return gin.H{"count": allCount, "page": allPage, "current_page": req.Page, "list": model}
}

func (ServiceRoom) Rename(ServiceId int, userId int, rename string) {
	var serviceRoom Service2.ServiceRoom
	Base.MysqlConn.Model(&serviceRoom).Where("user_id = ? and service_id =?", userId, ServiceId).Update("rename", rename)
}
