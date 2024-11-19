package v1

import (
	"tools-admin/backend/model"
	"tools-admin/backend/pkg/db"
	"tools-admin/backend/pkg/log"
	"tools-admin/backend/service"

	"github.com/gin-gonic/gin"
	"net/http"
)

var dbService = service.NewDatabaseService(db.Db)

// GetDatabases 获取数据库列表
func GetDatabases(c *gin.Context) {
	var req model.DatabaseListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "请求参数错误",
			"data":    err.Error(),
		})
		return
	}

	resp, err := dbService.List(&req)
	if err != nil {
		log.Error("获取数据库列表失败: %v", err)
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "获取数据库列表失败",
			"data":    err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "获取数据库列表成功",
		"data":    resp,
	})
}

// CreateDatabase 创建数据库连接
func CreateDatabase(c *gin.Context) {
	var req model.DatabaseCreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "请求参数错误",
			"data":    err.Error(),
		})
		return
	}

	if err := dbService.Create(&req); err != nil {
		log.Error("创建数据库连接失败: %v", err)
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "创建数据库连接失败",
			"data":    err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "创建数据库连接成功",
	})
}

// TestConnection 测试数据库连接
func TestConnection(c *gin.Context) {
	var req model.DatabaseTestReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "请求参数错误",
			"data":    err.Error(),
		})
		return
	}

	if err := dbService.TestConnection(&req); err != nil {
		log.Error("测试数据库连接失败: %v", err)
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "数据库连接测试失败",
			"data":    err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "数据库连接测试成功",
	})
}

// ExecuteQuery 执行SQL查询
func ExecuteQuery(c *gin.Context) {
	var req model.QueryExecuteReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "请求参数错误",
			"data":    err.Error(),
		})
		return
	}

	resp, err := dbService.ExecuteQuery(&req, c)
	if err != nil {
		log.Error("执行SQL查询失败: %v", err)
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "执行SQL查询失败",
			"data":    err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "执行SQL查询成功",
		"data":    resp,
	})
}

// GetTables 获取数据库表列表
func GetTables(c *gin.Context) {
	var req model.TableListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "请求参数错误",
			"data":    err.Error(),
		})
		return
	}

	resp, err := dbService.GetTables(&req)
	if err != nil {
		log.Error("获取数据库表列表失败: %v", err)
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "获取数据库表列表失败",
			"data":    err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "获取数据库表列表成功",
		"data":    resp,
	})
}

// GetTableSchema 获取表结构
func GetTableSchema(c *gin.Context) {
	var req model.TableSchemaReq
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "请求参数错误",
			"data":    err.Error(),
		})
		return
	}

	resp, err := dbService.GetTableSchema(&req)
	if err != nil {
		log.Error("获取表结构失败: %v", err)
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "获取表结构失败",
			"data":    err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "获取表结构成功",
		"data":    resp,
	})
}

// UpdateDatabase 更新数据库连接
func UpdateDatabase(c *gin.Context) {
	id := c.Param("id")
	var req model.DatabaseUpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "请求参数错误",
			"data":    err.Error(),
		})
		return
	}

	req.ID = id
	if err := dbService.Update(&req); err != nil {
		log.Error("更新数据库连接失败: %v", err)
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "更新数据库连接失败",
			"data":    err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "更新数据库连接成功",
	})
}

// DeleteDatabase 删除数据库连接
func DeleteDatabase(c *gin.Context) {
	id := c.Param("id")
	if err := dbService.Delete(id); err != nil {
		log.Error("删除数据库连接失败: %v", err)
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "删除数据库连接失败",
			"data":    err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "删除数据库连接成功",
	})
}

// TestDatabaseConnection 测试数据库连接
func TestDatabaseConnection(c *gin.Context) {
	id := c.Param("id")
	if err := dbService.TestConnectionByID(id); err != nil {
		log.Error("测试数据库连接失败: %v", err)
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "数据库连接测试失败",
			"data":    err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "数据库连接测试成功",
	})
}
