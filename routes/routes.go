package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/kwa0x2/realtime-chat-backend/controller"
	"github.com/kwa0x2/realtime-chat-backend/middlewares"
	"github.com/kwa0x2/realtime-chat-backend/socket/adapter"
	"github.com/zishang520/socket.io/socket"
)

func AuthRoute(router *gin.Engine, authController controller.IAuthController) {
	authRoutes := router.Group("/api/v1/auth")
	{
		authRoutes.GET("login", authController.GoogleLogin)
		authRoutes.POST("logout", authController.Logout)
		authRoutes.POST("signup", authController.SignUp)
		authRoutes.GET("callback", authController.GoogleCallback)
		authRoutes.GET("", middlewares.SessionMiddleware(), authController.CheckAuth)
	}
}

func UserRoute(router *gin.Engine, userController controller.IUserController) {
	userRoutes := router.Group("/api/v1/user")
	{
		userRoutes.PATCH("username", middlewares.SessionMiddleware(), userController.UpdateUsername)
		userRoutes.POST("upload-profile-photo", middlewares.CombinedAuthMiddleware(), userController.UploadProfilePhoto)

	}
}

func MessageRoute(router *gin.Engine, messageController controller.IMessageController) {
	messageRoutes := router.Group("/api/v1/message")
	messageRoutes.Use(middlewares.SessionMiddleware())
	{
		messageRoutes.POST("history", messageController.GetMessageHistory)
	}
}

func FriendRoute(router *gin.Engine, friendController controller.IFriendController) {
	friendRoutes := router.Group("/api/v1/friend")
	{
		friendRoutes.GET("", friendController.GetFriends) // get all
		friendRoutes.GET("blocked", friendController.GetBlockedUsers)
		friendRoutes.PATCH("block", friendController.Block)
		friendRoutes.DELETE("", friendController.Delete)
	}
}

func RequestRoute(router *gin.Engine, requestController controller.IRequestController) {
	requestRoutes := router.Group("/api/v1/request")
	{
		requestRoutes.POST("", requestController.SendFriend) // send friend req
		requestRoutes.GET("", requestController.GetRequests) // get coming req
		requestRoutes.PATCH("", requestController.Patch)

	}
}

func RoomRoute(router *gin.Engine, roomController controller.IRoomController) {
	roomRoutes := router.Group("/api/v1/room")
	{
		roomRoutes.POST("check", roomController.GetOrCreatePrivateRoom)
		roomRoutes.GET("chatlist", roomController.GetChatList)
	}
}

func SetupSocketIO(router *gin.Engine, server *socket.Server, socketAdapter adapter.ISocketAdapter) {
	socketAdapter.HandleConnection()

	router.GET("socket.io/*any", middlewares.SessionMiddleware(), gin.WrapH(server.ServeHandler(nil)))
	router.POST("socket.io/*any", middlewares.SessionMiddleware(), gin.WrapH(server.ServeHandler(nil)))
}
