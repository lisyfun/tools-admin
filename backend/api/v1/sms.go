package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetSMSTemplates 获取短信模板列表
func GetSMSTemplates(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "获取短信模板列表成功",
		"data": []interface{}{}, // TODO: 实现短信模板列表查询
	})
}

// CreateSMSTemplate 创建短信模板
func CreateSMSTemplate(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "创建短信模板成功",
		"data": nil,
	})
}

// UpdateSMSTemplate 更新短信模板
func UpdateSMSTemplate(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "更新短信模板成功",
		"data": nil,
	})
}

// DeleteSMSTemplate 删除短信模板
func DeleteSMSTemplate(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "删除短信模板成功",
		"data": nil,
	})
}

// GetSMSRecipients 获取短信接收人列表
func GetSMSRecipients(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "获取短信接收人列表成功",
		"data": []interface{}{}, // TODO: 实现短信接收人列表查询
	})
}

// CreateSMSRecipient 创建短信接收人
func CreateSMSRecipient(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "创建短信接收人成功",
		"data": nil,
	})
}

// UpdateSMSRecipient 更新短信接收人
func UpdateSMSRecipient(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "更新短信接收人成功",
		"data": nil,
	})
}

// DeleteSMSRecipient 删除短信接收人
func DeleteSMSRecipient(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "删除短信接收人成功",
		"data": nil,
	})
}
