package Common

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"server/App/Common"
	Common2 "server/App/Http/Request/Common"
	Common4 "server/App/Logic/Common"
	Service4 "server/App/Logic/Service"
	Common3 "server/App/Model/Common"
	Service2 "server/App/Model/Service"
	"server/Base"
	"strings"
	"time"
)

type Payment struct{}

func (p Payment) sign(Account int64, Time int64) string {
	// sign=md5(UserId+Account+Time+password)
	uid := Base.AppConfig.PayConfig.UserId
	password := Base.AppConfig.PayConfig.Password
	return Common.Tools{}.Md5(fmt.Sprintf("%d%d%d%s", uid, Account, Time, password))
}

func (Payment) createCallbackUrl(Type string, Member string, Day int64, Time int64, userOrderId string) string {
	sign := Common.Tools{}.Md5(fmt.Sprintf("type=%s&member=%s&day=%d&key=%s&time=%d&user_order_id=%s",
		Type, Member, Day, Base.AppConfig.AesKey, Time, userOrderId))
	param := fmt.Sprintf("type=%s&member=%s&day=%d&sign=%s&time=%d&user_order_id=%s", Type, Member, Day, sign, Time, userOrderId)
	return fmt.Sprintf("%s/api/payment/callback?%s", Base.AppConfig.HttpHost, param)
}

func (Payment) Callback(c *gin.Context) {
	var req Common2.CallbackReq
	if err := c.ShouldBind(&req); err != nil {
		Common.ApiResponse{}.Error(c, "请求错误", gin.H{})
		return
	}

	if !req.AuthSign() {
		Common.ApiResponse{}.Error(c, "签名错误", gin.H{})
		return
	}

	var ticket Common3.PaymentTicket
	if err := Base.MysqlConn.Where("id = ? and status = 'wait'", req.Id).First(&ticket).Error; err != nil {
		Common.ApiResponse{}.Error(c, "订单不存在", gin.H{})
		return
	}

	var serverId int
	if ticket.ServiceMember == "" {
		member := Common.Tools{}.CreateActiveMember()
		service := Service4.Service{}.Create(0, member, "user", ticket.Day)
		ticket.ServiceMember = service.Member
		Base.MysqlConn.Model(&ticket).Where("id = ?", ticket.Id).Updates(map[string]interface{}{
			"service_member": service.Member, "status": "success",
		})
		serverId = service.ServiceId
	} else {
		_, _ = Service4.Service{}.Renew(ticket.ServiceMember, ticket.Day, 70)
	}

	_ = Common4.Domain{}.Bind(serverId)
	c.String(200, "success")
}

func (p Payment) Post(c *gin.Context) {
	var req Common2.PostMemberReq
	if err := c.ShouldBind(&req); err != nil {
		Common.ApiResponse{}.Error(c, "请求错误", gin.H{})
		return
	}

	if req.Member != "" {
		var member Service2.Service
		if err := Base.MysqlConn.Where("member = ?", req.Member).First(&member).Error; err != nil {
			Common.ApiResponse{}.Error(c, "用户不存在", gin.H{})
			return
		}
	}

	// 创建订单
	Type := Common.Tools{}.If(req.Member == "", "create_member", "renew_member").(string)
	order, err := Common4.PaymentTicketLogic{}.CreateOrder(req.Member, int(req.Day), Type, 10, int(10*req.Day))
	if err != nil {
		Common.ApiResponse{}.Error(c, err.Error(), gin.H{})
		return
	}

	// 自定义回调地址
	now := time.Now().Unix()
	callback := p.createCallbackUrl(Type, req.Member, int64(req.Day), now, fmt.Sprintf("%d", order.Id))
	returnUrl := fmt.Sprintf("%s/api/payment/show?id=%d&time=%d&sign=%s", Base.AppConfig.HttpHost, order.Id, now, p.createSign(order.Id, now))
	res := Common.Tools{}.HttpPost(Base.AppConfig.PayUrl, gin.H{
		"user_id":    Base.AppConfig.PayConfig.UserId,
		"callback":   callback,
		"account":    req.Day * 10,
		"time":       now,
		"sign":       p.sign(int64(req.Day*10), now),
		"ext":        fmt.Sprintf("自助操作卡密:%s 共%d天", strings.TrimSpace(req.Member), req.Day),
		"return_url": returnUrl,
	})

	var resMap map[string]interface{}
	_ = json.Unmarshal(res, &resMap)

	_, ok := resMap["code"]

	fmt.Println("RESPONSE", resMap)
	if !ok {
		Common.ApiResponse{}.Error(c, "请求失败", gin.H{})
		return
	}

	if resMap["code"].(float64) != 200 {
		Common.ApiResponse{}.Error(c, resMap["msg"].(string), gin.H{})
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, resMap["data"].(map[string]interface{})["pay_url"].(string))
	//Common.ApiResponse{}.Success(c, "ok", resMap)
}

func (p Payment) Return(c *gin.Context) {
	var req Common2.ShowReq
	if err := c.BindQuery(&req); err != nil {
		Common.ApiResponse{}.Error(c, "请求错误", gin.H{"err": err.Error()})
		return
	}

	if !req.CheckAuthSign() {
		Common.ApiResponse{}.Error(c, "签名错误", gin.H{})
		return
	}

	var order Common3.PaymentTicket
	if err := Base.MysqlConn.Where("id = ?", req.Id).First(&order).Error; err != nil {
		Common.ApiResponse{}.Error(c, "订单不存在", gin.H{})
		return
	}

	c.HTML(http.StatusOK, "callback.html", gin.H{"member": order.ServiceMember})
}

func (p Payment) createSign(id int, times int64) string {
	return Common.Tools{}.Md5(fmt.Sprintf("%d%d%s", id, times, Base.AppConfig.AesKey))
}
