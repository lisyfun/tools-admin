package router

import (
	v1 "tools-admin/backend/api/v1"
	"tools-admin/backend/common/config"
	"tools-admin/backend/middleware/auth"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func InitRouter(r *gin.Engine) {
	// 配置CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * 60 * 60, // 12小时
	}))

	// API v1
	apiV1 := r.Group("/api/v1")
	{
		// 登录相关路由
		authHandler := v1.NewAuthHandler()
		apiV1.POST("/login", authHandler.Login)
		apiV1.POST("/logout", authHandler.Logout)
		apiV1.GET("/user/info", auth.JWTAuthMiddleware(config.Config.Server.JWTSecret), authHandler.GetUserInfo)
		apiV1.POST("/user/reset-password", authHandler.ResetPassword)

		// 需要认证的路由
		v1Group := apiV1.Group("")
		v1Group.Use(auth.JWTAuthMiddleware(config.Config.Server.JWTSecret))
		{
			// 菜单相关路由
			v1Group.GET("/menus", v1.GetMenus)
			v1Group.POST("/menu", v1.CreateMenu)
			v1Group.PUT("/menu/:id", v1.UpdateMenu)
			v1Group.DELETE("/menu/:id", v1.DeleteMenu)

			// Dashboard相关路由
			dashboardApi := &v1.DashboardApi{}
			dashboard := v1Group.Group("/dashboard")
			{
				dashboard.GET("/overview", dashboardApi.GetOverview)
				dashboard.GET("/task-chart", dashboardApi.GetTaskChart)
				dashboard.GET("/sms-chart", dashboardApi.GetSmsChart)
			}

			// 任务管理
			taskAPI := v1Group.Group("/task")
			{
				taskAPI.GET("", v1.GetTasks)
				taskAPI.GET("/:id", v1.GetTaskById)
				taskAPI.POST("", v1.CreateTask)
				taskAPI.PUT("/:id", v1.UpdateTask)
				taskAPI.DELETE("/:id", v1.DeleteTask)
				taskAPI.DELETE("/batch", v1.BatchDeleteTasks)
				taskAPI.GET("/:id/logs", v1.GetTaskLogs)
				taskAPI.POST("/:id/run", v1.RunTask)
				taskAPI.GET("/cron-patterns", v1.GetCommonCronPatterns)
				taskAPI.PATCH("/:id/status", v1.UpdateTaskStatus)
				taskAPI.GET("/next-run-times", v1.GetNextRunTimes) // 新增：获取下次执行时间
			}

			// 用户相关路由
			v1Group.GET("/users", v1.GetUsers)
			v1Group.POST("/user", v1.CreateUser)
			v1Group.PUT("/user/:id", v1.UpdateUser)
			v1Group.DELETE("/user/:id", v1.DeleteUser)
		}
	}
}
