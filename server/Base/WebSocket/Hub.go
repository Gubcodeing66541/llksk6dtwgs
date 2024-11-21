package WebSocket

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"sync"
)

var Lock sync.RWMutex

type Hub struct {
	UserListMap     map[string]map[string]Connect
	ServiceBindUser map[string]int
}

func (h *Hub) AddUser(c Connect) {
	Lock.Lock()
	defer Lock.Unlock()
	if _, ok := h.UserListMap[c.UserId]; !ok {
		h.UserListMap[c.UserId] = map[string]Connect{}
	}

	h.UserListMap[c.UserId][c.ConnId] = c
}

func (h *Hub) DelUser(c Connect) {

	Lock.Lock()
	defer Lock.Unlock()

	if _, ok := h.UserListMap[c.UserId][c.ConnId]; !ok {
		return
	}

	delete(h.UserListMap[c.UserId], c.ConnId)

	if len(h.UserListMap[c.UserId]) == 0 {
		delete(h.UserListMap, c.UserId)
		delete(h.ServiceBindUser, c.UserId)
	}
}

func (h *Hub) SendToUserId(userId string, message []byte) {
	Lock.RLock()
	defer Lock.RUnlock()
	if _, ok := h.UserListMap[userId]; !ok {
		return
	}

	var userConn Connect
	for _, userConn = range h.UserListMap[userId] {
		err := userConn.Conn.WriteMessage(1, message)
		if err != nil {
			fmt.Print("websocket-err: SendToUserId", err.Error())
		}
	}
}

func (h *Hub) SendToConnId(userId string, connId string, message []byte) {
	Lock.RLock()
	defer Lock.RUnlock()
	// 不在线则跳过发送
	if _, ok := h.UserListMap[userId][connId]; !ok {
		return
	}

	err := h.UserListMap[userId][connId].Conn.WriteMessage(1, message)
	if err != nil {
		fmt.Print("websocket-err: SendToConnId", err.Error())
	}
	return
}

func (h *Hub) UserIdIsOnline(userId string) int {
	Lock.RLock()
	defer Lock.RUnlock()
	_, ok := h.UserListMap[userId]
	if ok {
		return 1
	}
	return 0
}

func (h *Hub) GetAllStatus() map[string]interface{} {
	Lock.RLock()
	defer Lock.RUnlock()

	return map[string]interface{}{"h": gin.H{
		"ServiceBindUser": h.ServiceBindUser,
	}}
}

func (h *Hub) BindUser(ServiceUserId string, UserId int) {
	Lock.Lock()
	defer Lock.Unlock()
	h.ServiceBindUser[ServiceUserId] = UserId
}

func (h *Hub) GetBindUser(ServiceId string) int {
	Lock.RLock()
	defer Lock.RUnlock()
	userId, ok := h.ServiceBindUser[ServiceId]
	if ok {
		return userId
	}
	return 0
}

// Run 启动服务
func (h *Hub) Run() {

}
