package Base

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/nsqio/go-nsq"
	"io/ioutil"
	"math/rand"
	"os"
	"server/App/Model/Agent"
	"server/App/Model/Common"
	GroupModel "server/App/Model/Group"
	"server/App/Model/Log"
	Service2 "server/App/Model/Service"
	"server/App/Model/Setting"
	"server/App/Model/User"
	"server/Base/Config"
	"server/Base/Nsq"
	"server/Base/WebSocket"
	"time"
)

type Base struct{}

var AppConfig Config.App

var MysqlConn *gorm.DB

var WebsocketHub WebSocket.Hub

var RedisPool *redis.Pool //创建redis连接池

var Producer *nsq.Producer

func (b Base) Init() {
	rand.Seed(time.Now().Unix())

	b.initConfig()
	if AppConfig.Model != "dev" {
		b.InitConsumer()
	}

	b.initMysql()
	b.initWebSocketHub()
	b.initRedis()
	b.initSqlDate()
}

func (b Base) InitConsumer() {
	Producer = Nsq.NsqConsumer{}.CreateProducer(AppConfig.Mq.Nsq.Host)
	//Nsq.NsqConsumer{}.InitConsumer(Constant.Topic, Constant.Channel, AppConfig.Mq.Nsq.Host)
}

func (b Base) initWebSocketHub() {
	WebsocketHub = WebSocket.Hub{
		UserListMap:     map[string]map[string]WebSocket.Connect{},
		ServiceBindUser: map[string]int{},
	}
	WebsocketHub.Run()
}

// 配置初始化
func (b Base) initConfig() {
	file, err := os.Open("./config.json")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	res, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Print(err.Error())
		panic("json配置文件打开失败")
	}

	err = json.Unmarshal(res, &AppConfig)
	if err != nil {
		fmt.Print(err.Error())
		panic("json 配置解析异常")
	}
}

// mysql 初始化
func (b Base) initMysql() {

	var err error
	c := AppConfig.Database.Mysql
	connStr := fmt.Sprintf("%s:%s@(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", c.Username, c.Password, c.Host, c.Port, c.Database)
	MysqlConn, err = gorm.Open("mysql", connStr)
	if err != nil {
		fmt.Println(connStr)
		fmt.Print(err.Error())
		panic("mysql 初始异常")
	}

	MysqlConn.DB().SetMaxIdleConns(10)
	MysqlConn.DB().SetMaxOpenConns(100)

	MysqlConn.DB().SetConnMaxLifetime(20 * time.Second)
	MysqlConn.DB().SetMaxOpenConns(100)
	//MysqlConn.LogMode(AppConfig.Debug)

	auto := MysqlConn.Set("gorm:table_options", "ENGINE=InnoDB CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci")
	auto.AutoMigrate(&User.User{}, &User.UserLoginLog{}, &User.UserAuthMap{})
	auto.AutoMigrate(&Service2.Service{}, &Service2.ServiceBlack{}, &Service2.ServiceRoom{}, &Service2.ServiceMessage{}, &Service2.ServiceRoomDetail{}, &Service2.BotServiceMessage{})
	auto.AutoMigrate(&Common.Domain{}, &Common.Message{}, &Common.Order{}, &Common.Rename{})
	auto.AutoMigrate(&GroupModel.Group{}, &GroupModel.GroupUser{}, &GroupModel.GroupMessage{})
	auto.AutoMigrate(&Common.PaymentTicket{}, &Common.Ip{})
	auto.AutoMigrate(&Agent.Agent{}, &Agent.AgentAccountLog{}, &Log.CheckDomainLog{})

	auto.AutoMigrate(&Setting.Setting{}, &User.OnlyId{})

}

func (b Base) initRedis() {
	redisConfig := AppConfig.Database.Redis
	config := fmt.Sprintf("%s:%d", redisConfig.Host, redisConfig.Port)
	RedisPool = &redis.Pool{ //实例化一个连接池
		MaxIdle:     16,  //最初的连接数量
		MaxActive:   0,   //连接池最大连接数量,不确定可以用0（0表示自动定义），按需分配
		IdleTimeout: 300, //连接关闭时间 300秒 （300秒不使用自动关闭）
		Dial: func() (redis.Conn, error) { //要连接的redis数据库
			option := redis.DialPassword(redisConfig.Password)
			return redis.Dial("tcp", config, option)
		},
	}
}

func (b Base) initSqlDate() {
	var rename Common.Rename
	err := MysqlConn.First(&rename).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		fmt.Println("initSqlDate", err.Error())
		return
	}
	if rename.Id != 0 {
		fmt.Println("RENAME", rename)
		return
	}

	filename := "./rename.md"
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("open file err", err.Error())
		return
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	count := 0
	for scanner.Scan() {
		lineText := scanner.Text() // 获取当前行的文本内容
		fmt.Println(lineText)      // 输出每一行的内容
		renameData := Common.Rename{
			Rename: lineText,
		}
		MysqlConn.Create(&renameData)
		if count >= 500 {
			count = 0
			time.Sleep(time.Second)
		}
		count++
		fmt.Println(lineText)
	}

	if err := scanner.Err(); err != nil {
		panic("读取命名文件时发生错误")
	}
}
