package router

import (
	"github.com/gin-gonic/gin"

	"kefu-server/controllers"
	"kefu-server/middleware"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.Use(middleware.CORS())

	// 创建控制器实例
	userController := &controllers.UserController{}
	// API 路由组
	api := r.Group("/api/v1")
	{
		// 不需要认证的路由
		api.POST("/login", userController.Login)

		// 需要认证的路由
		auth := api.Group("/")
		auth.Use(middleware.AuthMiddleware())
		{
			// 用户管理路由
			auth.POST("/logout", userController.Logout)
			user := auth.Group("/user")
			{
				user.GET("/info", userController.GetUserInfo)
			}
		}
	}

	return r
}
