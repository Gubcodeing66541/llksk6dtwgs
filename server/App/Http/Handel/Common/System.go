package Common

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"path"
	"server/App/Common"
	Common2 "server/App/Logic/Common"
	"server/Base"
	"time"
)

type System struct{}

// @summary 系统默认文件上传
// @tags 公共接口
// @Param token header string true "认证token"
// @Param image  formData file true "文件参数"
// @Router /api/system/upload [post]
func (s System) Upload(c *gin.Context) {
	f, err := c.FormFile("image")
	if err != nil {
		Common.ApiResponse{}.Error(c, "请选择需要上传的文件", gin.H{"err": err.Error()})
		return
	}

	if f == nil {
		Common.ApiResponse{}.Error(c, "请选择需要上传的文件.", gin.H{})
		return
	}

	FileMap := map[string]string{"image/png": "png", "image/gif": "gif", "image/jpeg": "jpg", "video/mp4": "mp4", "video/ogg": "ogg", "video/webm": "webm"}
	if FileMap[f.Header.Get("Content-Type")] == "" {
		Common.ApiResponse{}.Error(c, "只允许上传png gif jpg格式图片和mp4,ogg,webm格式的视频", gin.H{})
		return
	}

	if f.Size > 1024*1024*60 {
		Common.ApiResponse{}.Error(c, "文件大小不能超出60mb", gin.H{"size": f.Size})
		return
	}

	// 拼接保存路径 将读取的文件保存在服务端
	rootPath := fmt.Sprintf("/static/upload/%s", time.Now().Format("20060102"))
	fileName := fmt.Sprintf("%s.%s", Common.Tools{}.RandFileName(c), FileMap[f.Header.Get("Content-Type")])
	dst := path.Join("."+rootPath, fileName)
	_ = os.MkdirAll("."+rootPath, os.ModePerm)

	// 保存文件
	err = c.SaveUploadedFile(f, dst)
	if err != nil {
		Common.ApiResponse{}.Error(c, "error", gin.H{"err": err.Error(), "file": dst})
		return
	}

	filePath := rootPath + "/" + fileName
	Common.ApiResponse{}.Success(c, "OK", gin.H{"file_name": filePath, "file_type": FileMap[f.Header.Get("Content-Type")]})
}

func (System) Live(c *gin.Context) {
	key := c.Param("key")
	token := Common.RedisTools{}.GetString(key)

	if token == "" {
		c.String(http.StatusOK, "视频通话已过期，请重新发起新的连线通话~")
		return
	}
	link := fmt.Sprintf("%s/live/index.html", Common2.Domain{}.GetLive())
	c.HTML(http.StatusOK, "localstorage.html", gin.H{
		"key":  "live_token",
		"val":  token,
		"url":  link,
		"key2": "live_host",
		"val2": Base.AppConfig.LiveAppHost,
	})
}
func (System) AliUpload(c *gin.Context) {
	oss := Base.AppConfig.Oss.Ali
	url, _ := Common.Oss{}.GetAliToken(oss.Region, oss.AccessKeyId, oss.AccessKeySecret)
	Common.ApiResponse{}.Success(c, "OK", gin.H{"url": url})
}
