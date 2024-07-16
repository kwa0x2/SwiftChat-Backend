package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/kwa0x2/realtime-chat-backend/adapter"
	"github.com/kwa0x2/realtime-chat-backend/controller"
	"github.com/kwa0x2/realtime-chat-backend/gateway"
	"github.com/kwa0x2/realtime-chat-backend/middlewares"
	"github.com/kwa0x2/realtime-chat-backend/service"
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
	userRoutes := router.Group("/api/v1/user")
	// userRoutes.Use(middlewares.JwtMiddleware())
	{
		userRoutes.POST("signup", userController.InsertUser)
		userRoutes.GET("", userController.GetAll)
		userRoutes.GET("get", userController.GetByEmail)

	}
}

func MessageRoute(router *gin.Engine, messageController *controller.MessageController) {
	messageRoutes := router.Group("/api/v1/message")
	messageRoutes.Use(middlewares.SessionMiddleware())
	{
		messageRoutes.POST("conversation/private", messageController.GetPrivateConversation)
	}
}

func FriendshipRoute(router *gin.Engine, friendshipController *controller.FriendshipController) {
	friendshipRoutes := router.Group("/api/v1/friendship")
	{
		friendshipRoutes.POST("", friendshipController.SendFriendRequest)
		friendshipRoutes.GET("coming", friendshipController.GetComingRequests)
		friendshipRoutes.GET("friends", friendshipController.GetFriends)
		friendshipRoutes.GET("blockeds", friendshipController.GetBlockeds)
		friendshipRoutes.PUT("block", friendshipController.Block)
		friendshipRoutes.DELETE("", friendshipController.Delete)
		friendshipRoutes.PUT("accept", friendshipController.Accept)
		friendshipRoutes.DELETE("reject", friendshipController.Reject)
	}
}

func SetupSocketIO(router *gin.Engine, io *socket.Server, messageService *service.MessageService, userService *service.UserService, friendshipService *service.FriendshipService) {
	socketGateway := gateway.NewSocketGateway(io)
	socketAdapter := adapter.NewSocketAdapter(socketGateway, messageService, userService, friendshipService)

	socketAdapter.HandleConnection()

	router.GET("socket.io/*any", middlewares.SessionMiddleware(), gin.WrapH(io.ServeHandler(nil)))
	router.POST("socket.io/*any", middlewares.SessionMiddleware(), gin.WrapH(io.ServeHandler(nil)))
}