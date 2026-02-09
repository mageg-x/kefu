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
	appController := &controllers.AppController{}
	visitorController := &controllers.VisitorController{}
	agentController := &controllers.AgentController{}
	// API 路由组
	api := r.Group("/api/v1")
	{
		// 不需要认证的路由
		api.POST("/login", userController.Login)
		api.GET("/config", appController.GetConfig)

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

			// App 管理路由
			app := auth.Group("/apps")
			{
				app.GET("/list", appController.GetApps)
				app.POST("/create", appController.CreateApp)
				app.PUT("/update", appController.UpdateApp)
				app.DELETE("/delete", appController.DeleteApp)
			}
		}
	}

	r.GET("/ws/chat", visitorController.WSHandler)
	r.GET("/ws/agent", middleware.AuthMiddleware(), agentController.WSHandler)

	return r
}
