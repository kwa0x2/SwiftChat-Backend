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
		authRoutes.GET("login", authController.GoogleLogin)
		authRoutes.POST("logout", authController.Logout)
		authRoutes.POST("signup", authController.SignUp)
		authRoutes.GET("callback", authController.GoogleCallback)
		authRoutes.GET("", middlewares.SessionMiddleware(), authController.CheckAuth)
	}
}

func UserRoute(router *gin.Engine, userController *controller.UserController) {
	userRoutes := router.Group("/api/v1/user")
	{
		userRoutes.PATCH("username", middlewares.SessionMiddleware(), userController.UpdateUsername)
		userRoutes.POST("upload-profile-picture", middlewares.CombinedAuthMiddleware(), userController.UploadProfilePicture)

	}
}

func MessageRoute(router *gin.Engine, messageController *controller.MessageController) {
	messageRoutes := router.Group("/api/v1/message")
	messageRoutes.Use(middlewares.SessionMiddleware())
	{
		messageRoutes.POST("conversation/private", messageController.GetPrivateConversation)
		messageRoutes.POST("history", messageController.GetMessageHistory)
		messageRoutes.DELETE(":messageId", messageController.DeleteById)
		messageRoutes.PATCH("", messageController.UpdateMessageByIdBody)
	}
}

func FriendRoute(router *gin.Engine, friendController *controller.FriendController) {
	friendRoutes := router.Group("/api/v1/friend")
	{
		friendRoutes.GET("", friendController.GetFriends) // get all
		friendRoutes.GET("blocked", friendController.GetBlocked)
		friendRoutes.PATCH("block", friendController.Block)
		friendRoutes.DELETE("", friendController.Delete)
	}
}

func RequestRoute(router *gin.Engine, requestController *controller.RequestController) {
	requestRoutes := router.Group("/api/v1/request")
	{
		requestRoutes.POST("", requestController.Insert)           // send friend req
		requestRoutes.GET("", requestController.GetComingRequests) // get coming req
		requestRoutes.PATCH("", requestController.PatchUpdateRequest)
	}
}

func RoomRoute(router *gin.Engine, roomController *controller.RoomController) {
	roomRoutes := router.Group("/api/v1/room")
	{
		roomRoutes.POST("check", roomController.GetOrCreatePrivateRoom)
		roomRoutes.GET("chatlist", roomController.GetChatList)
	}
}

func SetupSocketIO(router *gin.Engine, io *socket.Server, messageService *service.MessageService, userService *service.UserService, friendService *service.FriendService, requestService *service.RequestService) {
	socketGateway := gateway.NewSocketGateway(io)
	socketAdapter := adapter.NewSocketAdapter(socketGateway, messageService, userService, friendService, requestService)

	socketAdapter.HandleConnection()

	router.GET("socket.io/*any", middlewares.SessionMiddleware(), gin.WrapH(io.ServeHandler(nil)))
	router.POST("socket.io/*any", middlewares.SessionMiddleware(), gin.WrapH(io.ServeHandler(nil)))
}
