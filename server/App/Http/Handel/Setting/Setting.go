package Setting

import (
	"github.com/gin-gonic/gin"
	"server/App/Common"
	Setting2 "server/App/Model/Setting"
	"server/Base"
)

type Setting struct{}

// Update 修改配置
func (Setting) Update(c *gin.Context) {
	RoleId := Common.Tools{}.GetRoleId(c)
	var setting Setting2.Setting

	// 查找记录，确保获取正确的记录
	if err := Base.MysqlConn.First(&setting, "service_id = ?", RoleId).Error; err != nil {
		Common.ApiResponse{}.Error(c, "Record not found", gin.H{})
		return
	}

	if setting.Id == 0 {

	}

	// 读取请求中的 JSON 数据
	var updateData map[string]interface{}
	if err := c.ShouldBindJSON(&updateData); err != nil {
		Common.ApiResponse{}.Error(c, "Invalid data", gin.H{})
		return
	}

	if err := Base.MysqlConn.Model(&setting).Where("service_id=?", RoleId).Updates(updateData).Error; err != nil {
		Common.ApiResponse{}.Error(c, "Update failed", gin.H{})
		return
	}

	// 返回成功响应
	Common.ApiResponse{}.Success(c, "ok", gin.H{})
}

// Get 获取配置
func (Setting) Get(c *gin.Context) {
	RoleId := Common.Tools{}.GetRoleId(c)
	var setting Setting2.Setting
	Base.MysqlConn.Find(&setting, "service_id = ?", RoleId)

	Common.ApiResponse{}.Success(c, "ok", gin.H{"data": setting})
	return
}
