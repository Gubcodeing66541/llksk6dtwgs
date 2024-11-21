package Common

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"server/Base"
)

type RedisTools struct{}

func (RedisTools) GetString(key string) string {
	conn := Base.RedisPool.Get()
	defer conn.Close()

	res, err := redis.String(conn.Do("Get", key))
	if err != nil {
		fmt.Println("Get", res, err.Error())
		return ""
	}
	return res
}

func (RedisTools) SetString(key string, val string) {
	fmt.Println("set redis key ", key, val)
	conn := Base.RedisPool.Get()
	defer conn.Close()
	rsp, err := conn.Do("Set", key, val)
	if err != nil {
		fmt.Println("Set", rsp, err.Error())
	}
	rsp, err = conn.Do("expire", key, 60)
	if err != nil {
		fmt.Println("expire", rsp, err.Error())
	}
}

func (RedisTools) Det(key string) {
	conn := Base.RedisPool.Get()
	defer conn.Close()
	rsp, err := conn.Do("del", key)
	if err != nil {
		fmt.Println("Set", rsp, err.Error())
	}
}

func (RedisTools) SetInt(key string, val int) int {
	conn := Base.RedisPool.Get()
	defer conn.Close()

	r, err := redis.Int(conn.Do("Set", key, val))
	if err != nil {
		fmt.Println("set  failed,", err)
		return -1
	}
	return r
}

func (RedisTools) GetInt(key string) int {
	conn := Base.RedisPool.Get()
	defer conn.Close()

	r, err := redis.Int(conn.Do("Get", key))
	if err != nil {
		fmt.Println("get  failed,", err)
		return -1
	}
	return r
}

func (RedisTools) SetStringByTime(key string, val string, time int) {
	fmt.Println("set redis key ", key, val)
	conn := Base.RedisPool.Get()
	defer conn.Close()
	rsp, err := conn.Do("Set", key, val)
	if err != nil {
		fmt.Println("Set", rsp, err.Error())
	}
	rsp, err = conn.Do("expire", key, time)
	if err != nil {
		fmt.Println("expire", rsp, err.Error())
	}
}
