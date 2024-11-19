package auth

import (
	"net/http"
	"strings"
	"tools-admin/backend/common/config"
	"tools-admin/backend/pkg/auth"
	"tools-admin/backend/pkg/log"

	"github.com/gin-gonic/gin"
)

func Auth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var token string

		// 先尝试从Token头获取
		token = ctx.Request.Header.Get("Token")
		
		// 如果Token头没有，再尝试从Authorization头获取
		if token == "" {
			authHeader := ctx.Request.Header.Get("Authorization")
			if authHeader != "" {
				parts := strings.SplitN(authHeader, " ", 2)
				if len(parts) == 2 && parts[0] == "Bearer" {
					token = parts[1]
				}
			}
		}

		// 如果都没有token，返回未认证错误
		if token == "" {
			ctx.JSON(http.StatusOK, gin.H{
				"code":    401,
				"message": "请先登录",
			})
			ctx.Abort()
			return
		}

		// 解析token
		claims, err := auth.ParseToken(token, config.Config.Server.JWTSecret)
		if err != nil {
			log.Error("Token解析失败: " + err.Error())
			ctx.JSON(http.StatusOK, gin.H{
				"code":    401,
				"message": "无效的token",
			})
			ctx.Abort()
			return
		}

		// 设置用户信息到上下文
		ctx.Set("userID", claims.UserID)
		ctx.Set("username", claims.Username)
		ctx.Set("roleID", claims.RoleID)
		log.Info("用户认证成功: %s (ID: %d)", claims.Username, claims.UserID)
		ctx.Next()
	}
}
