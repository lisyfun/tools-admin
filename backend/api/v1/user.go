package v1

import (
	"strconv"
	"tools-admin/backend/model"
	"tools-admin/backend/pkg/db"
	"tools-admin/backend/pkg/log"
	"tools-admin/backend/service"

	"github.com/gin-gonic/gin"
	"net/http"
)

var userService = service.NewUserService(db.Db)

// GetUsers 获取用户列表
func GetUsers(c *gin.Context) {
	var req model.UserListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "请求参数错误",
			"data":    err.Error(),
		})
		return
	}

	resp, err := userService.List(&req)
	if err != nil {
		log.Error("获取用户列表失败: %v", err)
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "获取用户列表失败",
			"data":    err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "获取用户列表成功",
		"data":    resp,
	})
}

// CreateUser 创建用户
func CreateUser(c *gin.Context) {
	var req model.UserCreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "请求参数错误",
			"data":    err.Error(),
		})
		return
	}

	if err := userService.Create(&req); err != nil {
		log.Error("创建用户失败: %v", err)
		code := 500
		msg := "创建用户失败"
		if err == service.ErrUsernameExists {
			code = 400
			msg = err.Error()
		}
		c.JSON(http.StatusOK, gin.H{
			"code":    code,
			"message": msg,
			"data":    err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "创建用户成功",
	})
}

// UpdateUser 更新用户
func UpdateUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "无效的用户ID",
		})
		return
	}

	var req model.UserUpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "请求参数错误",
			"data":    err.Error(),
		})
		return
	}

	if err := userService.Update(uint(id), &req); err != nil {
		log.Error("更新用户失败: %v", err)
		code := 500
		msg := "更新用户失败"
		if err == service.ErrUserNotFound {
			code = 400
			msg = err.Error()
		}
		c.JSON(http.StatusOK, gin.H{
			"code":    code,
			"message": msg,
			"data":    err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "更新用户成功",
	})
}

// DeleteUser 删除用户
func DeleteUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "无效的用户ID",
		})
		return
	}

	if err := userService.Delete(uint(id)); err != nil {
		log.Error("删除用户失败: %v", err)
		code := 500
		msg := "删除用户失败"
		if err == service.ErrUserNotFound {
			code = 400
			msg = err.Error()
		}
		c.JSON(http.StatusOK, gin.H{
			"code":    code,
			"message": msg,
			"data":    err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "删除用户成功",
	})
}

// ResetUser 重置密码
func ResetUser(c *gin.Context) {
	var req model.ResetPasswordReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "请求参数错误",
			"data":    err.Error(),
		})
		return
	}

	if err := userService.ResetPassword(&req); err != nil {
		log.Error("重置密码失败: %v", err)
		code := 500
		msg := "重置密码失败"
		if err == service.ErrUserNotFound {
			code = 400
			msg = err.Error()
		}
		c.JSON(http.StatusOK, gin.H{
			"code":    code,
			"message": msg,
			"data":    err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "重置密码成功",
	})
}

func Register(v1 *gin.RouterGroup) {
	userGroup := v1.Group("/users")
	{
		userGroup.POST("", CreateUser)           // 创建用户
		userGroup.PUT("/:id", UpdateUser)        // 更新用户
		userGroup.DELETE("/:id", DeleteUser)     // 删除用户
		userGroup.GET("/:id", GetUsers)          // 获取用户
		userGroup.GET("", GetUsers)              // 获取用户列表
		userGroup.POST("/reset", ResetUser)      // 重置密码
	}
}
