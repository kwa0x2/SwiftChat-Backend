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
	userRoutes.Use(middlewares.JwtMiddleware())
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

func FriendRoute(router *gin.Engine, friendController *controller.FriendController) {
	friendRoutes := router.Group("/api/v1/friend")
	{
		friendRoutes.GET("", friendController.GetFriends) // get all
		friendRoutes.GET("blockeds", friendController.GetBlockeds)
		friendRoutes.PUT("block", friendController.Block)
		friendRoutes.DELETE("", friendController.Delete)
	}
}

func RequestRoute(router *gin.Engine, requestController *controller.RequestController){
	requestRoutes := router.Group("/api/v1/request")
	{
		requestRoutes.POST("", requestController.Insert) // send friend req
		requestRoutes.GET("", requestController.GetComingRequests) // get coming req
		requestRoutes.PATCH("accept", requestController.Accept)
		requestRoutes.PATCH("reject", requestController.Reject)
	}
}

func SetupSocketIO(router *gin.Engine, io *socket.Server, messageService *service.MessageService, userService *service.UserService, friendService *service.FriendService) {
	socketGateway := gateway.NewSocketGateway(io)
	socketAdapter := adapter.NewSocketAdapter(socketGateway, messageService, userService, friendService)

	socketAdapter.HandleConnection()

	router.GET("socket.io/*any", middlewares.SessionMiddleware(), gin.WrapH(io.ServeHandler(nil)))
	router.POST("socket.io/*any", middlewares.SessionMiddleware(), gin.WrapH(io.ServeHandler(nil)))
}