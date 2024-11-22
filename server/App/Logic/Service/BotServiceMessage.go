package Service

import (
	Service3 "server/App/Http/Request/Service"
	Service2 "server/App/Model/Service"
	"server/Base"
	"time"
)

// BotServiceMessage 快捷消息
type BotServiceMessage struct {
}

// Create  创建消息
func (BotServiceMessage) Create(serviceId int, req Service3.CreateBotServiceMessage) interface{} {
	return Base.MysqlConn.Create(&Service2.BotServiceMessage{
		ServiceId: serviceId, Status: "enable",
		MsgInfo: req.MsgInfo, MsgType: req.MsgType,
		Title: req.Title, Question: req.Question,

		Type: req.Type, CreateTime: time.Now(), Name: req.Name})
}

// Delete 删除招呼
func (BotServiceMessage) Delete(id int, serviceId int) {
	Base.MysqlConn.Delete(&Service2.BotServiceMessage{}, "id = ? and service_id = ? ", id, serviceId)
}

// Update 修改招呼
func (BotServiceMessage) Update(id int, name, msgType string, msgInfo string, serviceId int, status string, title string, question string) interface{} {
	var serviceMessage Service2.BotServiceMessage
	return Base.MysqlConn.Model(&serviceMessage).Where("id = ? and service_id = ?", id, serviceId).Updates(Service2.BotServiceMessage{
		MsgType: msgType, Name: name, MsgInfo: msgInfo, Status: status, Title: title, Question: question,
	})
}

// List 招呼列表
func (BotServiceMessage) List(serviceId int, typeStr string) []Service2.BotServiceMessage {
	var serviceMessage []Service2.BotServiceMessage
	Base.MysqlConn.Find(&serviceMessage, "service_id = ? and type = ?", serviceId, typeStr)
	return serviceMessage
}

func (BotServiceMessage) GetById(serviceId int, id int) Service2.BotServiceMessage {
	var serviceMessage Service2.BotServiceMessage
	Base.MysqlConn.Find(&serviceMessage, "service_id = ? and id = ?", serviceId, id)
	return serviceMessage
}
