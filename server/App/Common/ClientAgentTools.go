package Common

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"strings"
)

type ClientAgentTools struct{}

// 是否是微信检测
func (ClientAgentTools) IsWechat(c *gin.Context) bool {
	userAgent := c.GetHeader("User-Agent")
	if strings.Index(userAgent, "MicroMessenger") != -1 {
		return true
	}
	return false
}

// 是否是微信检测
func (ClientAgentTools) IsDouYin(c *gin.Context) bool {
	userAgent := c.GetHeader("User-Agent")
	if strings.Index(userAgent, "aweme") != -1 {
		return true
	}
	return false
}

// 是否支付宝检测
func (ClientAgentTools) IsAlipay(c *gin.Context) bool {
	userAgent := c.GetHeader("User-Agent")
	if strings.Index(userAgent, "Alipay") != -1 {
		return true
	}
	return false
}

func (ClientAgentTools) IsMobile(c *gin.Context) bool {
	userAgent := c.GetHeader("User-Agent")
	if len(userAgent) == 0 {
		return false
	}

	isMobile := false
	mobileKeywords := []string{"Mobile", "Android", "Silk/", "Kindle",
		"BlackBerry", "Opera Mini", "Opera Mobi"}

	for i := 0; i < len(mobileKeywords); i++ {
		if strings.Contains(userAgent, mobileKeywords[i]) {
			isMobile = true
			break
		}
	}

	return isMobile
}

func (ClientAgentTools) GetDrive(c *gin.Context) string {
	userAgent := c.GetHeader("User-Agent")
	fmt.Println("userAgent:", userAgent)
	if strings.Index(userAgent, "Windows") != -1 {
		return "Windows"
	}
	if strings.Index(userAgent, "Android") != -1 {
		return "Android"
	}
	if strings.Index(userAgent, "iPhone") != -1 {
		return "iPhone"
	}
	if strings.Index(userAgent, "iPod") != -1 {
		return "iPod"
	}
	if strings.Index(userAgent, "iPad") != -1 {
		return "iPad"
	}
	if strings.Index(userAgent, "Windows Phone") != -1 {
		return "Windows Phone"
	}
	if strings.Index(userAgent, "MQQBrowser") != -1 {
		return "QQ浏览器"
	}
	if strings.Index(userAgent, "iPhone") != -1 {
		return "iPhone"
	}
	return "未知"
}

func (ClientAgentTools) GetDriveInfo(c *gin.Context) string {
	return c.GetHeader("User-Agent")
}
