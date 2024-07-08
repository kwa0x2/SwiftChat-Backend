package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/kwa0x2/realtime-chat-backend/adapter"
	"github.com/kwa0x2/realtime-chat-backend/controller"
	"github.com/kwa0x2/realtime-chat-backend/gateway"
	"github.com/kwa0x2/realtime-chat-backend/middlewares"
	"github.com/zishang520/socket.io/socket"
)

func AuthRoute(router *gin.Engine, authController *controller.AuthController) {
	authRoutes := router.Group("/api/v1/auth")
	{
		authRoutes.GET("login",authController.GoogleLogin)
		authRoutes.GET("callback", authController.GoogleCallback)
		authRoutes.GET("",middlewares.SessionMiddleware(), authController.CheckAuth)
		authRoutes.POST("logout", authController.Logout)
	}
	
}

func UserRoute(router *gin.Engine, userController *controller.UserController) {
	userRoutes := router.Group("/api/v1/auth")
	userRoutes.Use(middlewares.JwtMiddleware())
	{
		userRoutes.POST("signup", userController.InsertUser)
	}
}

func SetupSocketIO(router *gin.Engine, io *socket.Server) {
	socketGateway := gateway.NewSocketGateway(io)
	socketAdapter := adapter.NewSocketAdapter(socketGateway)

	socketAdapter.HandleConnection()

	router.GET("socket.io/*any", middlewares.SessionMiddleware(), gin.WrapH(io.ServeHandler(nil)))
	router.POST("socket.io/*any", middlewares.SessionMiddleware(), gin.WrapH(io.ServeHandler(nil)))
}


