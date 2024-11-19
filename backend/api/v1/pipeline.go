package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetPipelines 获取流水线列表
func GetPipelines(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "获取流水线列表成功",
		"data": []interface{}{}, // TODO: 实现流水线列表查询
	})
}

// CreatePipeline 创建流水线
func CreatePipeline(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "创建流水线成功",
		"data": nil,
	})
}

// UpdatePipeline 更新流水线
func UpdatePipeline(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "更新流水线成功",
		"data": nil,
	})
}

// DeletePipeline 删除流水线
func DeletePipeline(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "删除流水线成功",
		"data": nil,
	})
}

// RunPipeline 运行流水线
func RunPipeline(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"message": "流水线运行成功",
		"data": nil,
	})
}
