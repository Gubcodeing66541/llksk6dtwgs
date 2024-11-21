package Common

import (
	"bytes"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/axgle/mahonia"
	"github.com/gin-gonic/gin"
	"github.com/livekit/protocol/auth"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"server/App/Model/Common"
	"server/Base"
	"strings"
	"time"
)

type Tools struct{}

var CHARS = []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z",
	"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z",
	"1", "2", "3", "4", "5", "6", "7", "8", "9", "0"}

/*
RandAllString  生成随机字符串([a~zA~Z0~9])

	lenNum 长度
*/
func (Tools) RandAllString(lenNum int) string {
	str := strings.Builder{}
	length := len(CHARS)
	for i := 0; i < lenNum; i++ {
		l := CHARS[rand.Intn(length)]
		str.WriteString(l)
	}
	return str.String()
}

func (Tools) Md5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

// 返回一个16位md5加密后的字符串
func (t Tools) Md516(data string) string {
	return t.Md5(data)[8:24]
}

// 返回一个16位md5加密后的字符串
func (t Tools) Md510(data string) string {
	return t.Md5(data)[14:24]
}

func (Tools) GetRoleId(c *gin.Context) int {
	userId, ok := c.Get("role_id")
	if !ok {
		fmt.Print("get_user_id_err:", userId)
		return 0
	}
	return userId.(int)
}

func (Tools) GetRoleType(c *gin.Context) string {
	userId, ok := c.Get("role_type")
	if !ok {
		fmt.Print("get_user_id_err:", userId)
		return ""
	}
	return userId.(string)
}

func (Tools) GetRoleGroupId(c *gin.Context) string {
	groupId, ok := c.Get("group_id")
	if !ok {
		fmt.Print("get_user_groupId_err:", groupId)
		return ""
	}

	return fmt.Sprintf("%s:%d", "group", groupId)
}

func (Tools) GetServiceId(c *gin.Context) int {
	serviceId, ok := c.Get("service_id")
	if !ok {
		fmt.Print("get_user_id_err:", serviceId)
		return 0
	}
	return serviceId.(int)
}

func (Tools) GetCookieToken(c *gin.Context) string {
	token, ok := c.Get("token")
	if !ok {
		fmt.Print("get_cookie_token_err:", token)
		return ""
	}
	return token.(string)
}

func (Tools) GetWebSocketGroupId(c *gin.Context) string {
	userId, idOk := c.Get("role_id")
	roleType, typeOk := c.Get("role_type")
	if !idOk || !typeOk {
		fmt.Print("get_user_id_err:", userId)
		return ""
	}
	return fmt.Sprintf("%s:%d", roleType, userId)
}

func (Tools) GetWebSocketId(c *gin.Context) string {
	userId, idOk := c.Get("role_id")
	roleType, typeOk := c.Get("role_type")
	if !idOk || !typeOk {
		fmt.Print("get_user_id_err:", userId)
		return ""
	}
	return fmt.Sprintf("%s:%d", roleType, userId)
}

func (Tools) GetServiceWebSocketId(ServiceId int) string {
	return fmt.Sprintf("service:%d", ServiceId)
}

func (Tools) GetUserWebSocketId(UserId int) string {
	return fmt.Sprintf("user:%d", UserId)
}

// 创建激活码
func (t Tools) CreateActiveCode(activateId int) string {
	str := fmt.Sprintf("activate_id:%d-time:%s-rand-%d", activateId, time.Now(), rand.Intn(99999))
	str = t.Md516(str)
	return fmt.Sprintf("%s%d", str, activateId)
}

// 创建用户名
func (t Tools) CreateUserName(activateId int) string {
	str := fmt.Sprintf("user_name:%d-time:%s-rand-%d", activateId, time.Now(), rand.Intn(99999))
	str = t.Md510(str)
	return fmt.Sprintf("%s%d", str, activateId)
}

// 创建用户名UUID
func (t Tools) CreateUUID(code string) string {
	str := fmt.Sprintf("user_name:%s-time:%s-rand-%d", code, time.Now(), rand.Intn(99999))
	str = t.Md510(str)
	return fmt.Sprintf("%s", str)
}

func (t Tools) CreateActiveMember() string {
	times := time.Now()
	str := fmt.Sprintf("MEMBER-%s-time:-rand-%d", time.Now(), rand.Intn(99999))
	str = fmt.Sprintf("LM-%02d%02d%02d%s", times.Year(), times.Month(), times.Day(), t.Md516(str))
	return str
}

func (t Tools) CreateOrderId(left string) string {
	times := time.Now()
	str := fmt.Sprintf("MEMBER-%s-time:-rand-%d", time.Now(), rand.Intn(99999))
	str = fmt.Sprintf("%s%02d%02d%02d%s", left, times.Year(), times.Month(), times.Day(), t.Md516(str))
	return str
}

func (t Tools) CreateServiceManagerMember() string {
	times := time.Now()
	str := fmt.Sprintf("service-manager-%s-time:-rand-%d", time.Now(), rand.Intn(99999))
	str = fmt.Sprintf("%02d%02d%02d%s", times.Year(), times.Month(), times.Day(), t.Md516(str))
	return str
}

func (t Tools) CreateUserMember() string {
	times := time.Now()
	str := fmt.Sprintf("USER-%s-time:-rand-%d", time.Now(), rand.Intn(99999))
	str = fmt.Sprintf("U%02d%02d%02d%s", times.Year(), times.Month(), times.Day(), t.Md516(str))
	return str
}

func (Tools) HttpGet(url string) []byte {
	res, err := http.Get(url)
	if err != nil {
		return []byte("")
	}

	defer res.Body.Close()
	str, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return []byte("")
	}
	return str
}

// 解析token
func (Tools) DecodeToken(token string, userInfo interface{}) error {
	if token == "" || token == "null" {
		return errors.New("token is null")
	}

	res, err := Encryption{}.AesDecryptCBC(token)
	if err != nil {
		return errors.New("token decode err")
	}

	err = json.Unmarshal(res, &userInfo)
	if err != nil {
		return errors.New("json decode err")
	}
	return nil
}

// 创建token
func (Tools) EncodeToken(roleId int, roleType string, serviceId int, groupId int, diversion string) string {
	randStr := fmt.Sprintf("%s_ID-%d-time%s,rand-%d", roleType, roleId, time.Now(), rand.Intn(10000))
	userToken := UserAuthToken{
		RoleId: roleId, GroupId: groupId, RoleType: roleType, Diversion: diversion,
		RandStr: randStr, Time: time.Now(), ServiceId: serviceId, Key: Base.AppConfig.Manager.AuthKey,
	}
	token, err := json.Marshal(userToken)
	if err != nil {
		return ""
	}
	return Encryption{}.AesEncryptCBC(token)
}

func (Tools) If(condition bool, ok interface{}, no interface{}) interface{} {
	if condition {
		return ok
	}
	return no
}

func (Tools) SendMsToMq(topic string, message []byte) bool {
	err := Base.Producer.Publish(topic, message)
	if err != nil {
		return false
	}
	return true
}

// 返回：IP地址的信息（格式：字符串的json）
func (Tools) IPInfo(ip string) (Ip, error) {
	client := http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(`https://whois.pconline.com.cn/ipJson.jsp?json=true&ip=` + ip)
	var ipModel Ip

	if err != nil {
		return ipModel, err
	}
	defer resp.Body.Close()
	var buffer [512]byte
	result := bytes.NewBuffer(nil)
	for {
		n, err := resp.Body.Read(buffer[0:])
		result.Write(buffer[0:n])
		if err != nil && err == io.EOF {
			break
		} else if err != nil {
			return ipModel, err
		}
	}

	resStr := mahonia.NewDecoder("gbk").ConvertString(result.String())
	err = json.Unmarshal([]byte(resStr), &ipModel)
	if err != nil {
		return ipModel, errors.New("解析结果错误")
	}
	return ipModel, nil
}

// 根据用户标识返回一个随机的文件名
func (Tools) RandFileName(c *gin.Context) string {
	return fmt.Sprintf(
		"uplod%s%d%s%d",
		Tools{}.GetRoleType(c),
		Tools{}.GetRoleId(c),
		time.Now().Format("060102030405"),
		rand.Intn(999999),
	)
}

func (t Tools) GetDefaultHead() string {
	return Base.AppConfig.HttpHost + "/static/static/service_head.png"
}

/*
//  keyStr 密钥
//  value  消息内容
*/
func (t Tools) HMACSHA1(keyStr, value string) string {

	key := []byte(keyStr)
	mac := hmac.New(sha1.New, key)
	mac.Write([]byte(value))
	//进行base64编码
	res := base64.StdEncoding.EncodeToString(mac.Sum(nil))

	return res
}

func (t Tools) GetHead() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return fmt.Sprintf("%su%d.jpg", Base.AppConfig.HeadImgUrl, r.Intn(3710))
}

// 获取随机昵称
func (t Tools) GetRename() string {
	var renameDb Common.Rename
	Base.MysqlConn.Raw("select * from renames order by rand() limit 1").Scan(&renameDb)
	return renameDb.Rename
}

func (Tools) CreateToken(apiKey string, apiSecret string, roomName string, RoleId string) string {
	//用于连接livekit服务器的认证密钥，livekit.yaml中获取
	canPublish := true
	canSubscribe := true

	//生成认证实体
	grant := auth.NewAccessToken(apiKey, apiSecret).AddGrant(&auth.VideoGrant{
		RoomJoin:     true,
		Room:         roomName,
		CanPublish:   &canPublish,
		CanSubscribe: &canSubscribe,
	})

	//设置实体对象
	jwt, err := grant.SetIdentity(RoleId).SetValidFor(time.Hour * 24).ToJWT()
	if err != nil {
		fmt.Println(err)
	}
	return jwt
}

func (Tools) GetLiveRoomName(serviceId int64, userId int64) string {
	return fmt.Sprintf("S:%d_U:%d", serviceId, userId)
}

// 创建服务组code
func (Tools) CreateDiversionCode(serviceId int64) string {
	return fmt.Sprintf("%d%d", serviceId, time.Now().Unix())
}

func (Tools) HttpPost(url string, data interface{}) []byte {
	jsonStr, _ := json.Marshal(data)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	if err != nil {
		return []byte("")
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return []byte("")
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	return body
}
