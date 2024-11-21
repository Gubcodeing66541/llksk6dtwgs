package Common

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"math/rand"
	"net/http"
	"server/App/Common"
	"server/Base"
	"server/Base/WebSocket"
	"time"
)

type WebSocketConnect struct{}

// 升级websocket
var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (WebSocketConnect) Conn(c *gin.Context) {
	//升级get请求为webSocket协议
	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Print("websocket-connect:err", err.Error())
		return
	}

	roleId := Common.Tools{}.GetRoleId(c)
	roleType := Common.Tools{}.GetRoleType(c)
	GroupId := Common.Tools{}.GetRoleGroupId(c)

	serviceId := Common.Tools{}.GetServiceId(c)

	Event{}.OnConnect(roleType, roleId, GroupId, serviceId, ws)
}

type Event struct{}

func (Event) OnConnect(roleType string, roleId int, groupId string, serviceId int, ws *websocket.Conn) {

	username := fmt.Sprintf("%s:%d", roleType, roleId)

	connId := fmt.Sprintf("user_id_%s_time_%s_rand_number_%d", username, time.Now(), rand.Intn(999999))
	connId = Common.Tools{}.Md5(connId)

	conn := WebSocket.Connect{UserId: username, Conn: ws, ConnId: connId}
	Base.WebsocketHub.AddUser(conn)

	//if groupId != "" {
	//	Base.WebsocketHub.JoinGroup(conn, groupId, true)
	//}

	if roleType == "user" {
		defer Common.ApiResponse{}.SendMsgToService(serviceId, "leave", gin.H{"user_id": roleId})
		Common.ApiResponse{}.SendMsgToService(serviceId, "online", gin.H{"user_id": roleId})
	} else {
		//Base.WebsocketHub.BindGroup(connId, 0)
	}

	defer ws.Close()
	defer Base.WebsocketHub.DelUser(conn)
	defer Event{}.OnClose(conn)

	for {
		err := ws.SetReadDeadline(time.Now().Add(20 * time.Second))
		if err != nil {
			return
		}
		_, message, err := ws.ReadMessage()
		if err != nil {
			break
		}

		Event{}.OnMessage(conn, message)
	}
}

func (Event) OnMessage(conn WebSocket.Connect, message []byte) {
}

func (Event) OnClose(conn WebSocket.Connect) {

}
