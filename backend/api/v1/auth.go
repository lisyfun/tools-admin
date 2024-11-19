package v1

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"tools-admin/backend/common/config"
	"tools-admin/backend/internal/model"
	"tools-admin/backend/internal/service"
	"tools-admin/backend/pkg/db"
	"tools-admin/backend/pkg/log"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{
		authService: service.NewAuthService(
			config.Config.Server.JWTSecret,
			time.Duration(config.Config.Server.JWTExpire)*time.Second,
		),
	}
}

type ResetPasswordRequest struct {
	Username    string `json:"username" binding:"required"`
	NewPassword string `json:"newPassword" binding:"required,min=6,max=20"`
}

// Login 用户登录
func (h *AuthHandler) Login(c *gin.Context) {
	var req service.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "请求参数错误",
		})
		return
	}

	resp, err := h.authService.Login(&req)
	if err != nil {
		switch err {
		case service.ErrInvalidCredentials:
			c.JSON(http.StatusOK, gin.H{
				"code":    401,
				"message": "用户名或密码错误",
			})
		case service.ErrUserDisabled:
			c.JSON(http.StatusOK, gin.H{
				"code":    403,
				"message": "用户已被禁用",
			})
		default:
			c.JSON(http.StatusOK, gin.H{
				"code":    500,
				"message": "系统错误",
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "登录成功",
		"data":    resp,
	})
}

// Register 用户注册
func (h *AuthHandler) Register(c *gin.Context) {
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "请求参数错误",
		})
		return
	}

	if err := h.authService.Register(&user); err != nil {
		if err == service.ErrUserAlreadyExists {
			c.JSON(http.StatusOK, gin.H{
				"code":    400,
				"message": "用户名已存在",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "系统错误",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "注册成功",
	})
}

// GetUserInfo 获取用户信息
func (h *AuthHandler) GetUserInfo(c *gin.Context) {
	userID, _ := c.Get("userID")
	var user model.User
	if err := db.Db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "获取用户信息失败",
		})
		return
	}

	// 清除敏感信息
	user.Password = ""

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "获取成功",
		"data":    user,
	})
}

// Logout 用户登出
func (h *AuthHandler) Logout(c *gin.Context) {
	// 由于使用的是 JWT，服务端不需要做特殊处理
	// 客户端只需要删除本地存储的 token 即可
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "登出成功",
	})
}

// ResetPassword 重置密码
func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 1, "message": "无效的请求参数"})
		return
	}

	err := h.authService.ResetPassword(req.Username, req.NewPassword)
	if err != nil {
		log.Error("重置密码失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "message": "重置密码失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "密码重置成功"})
}
