package Common

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

type Cookie struct {
}

func (Cookie) GetCookie(c *gin.Context) string {
	cookie, err := c.Cookie("key_cookie")
	if err != nil {
		fmt.Printf("get cookie err:%s", err.Error())
		return ""
	}
	return cookie
}

func (Cookie) SetCookie(c *gin.Context, key string, val string) {
	c.SetCookie(key, val, 60*24*24*1024, "/", "/", false, true)
}
