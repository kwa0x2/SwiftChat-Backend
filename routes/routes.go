package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/kwa0x2/swiftchat-backend/controller"
	"github.com/kwa0x2/swiftchat-backend/middlewares"
	"github.com/kwa0x2/swiftchat-backend/socket/adapter"
	"github.com/zishang520/socket.io/socket"
)

// region Auth Routes
func AuthRoute(router *gin.Engine, authController controller.IAuthController) {
	authRoutes := router.Group("/api/v1/auth")
	{

		authRoutes.GET("login/google", authController.GoogleLogin)                          // Google login
		authRoutes.GET("login/google/callback", authController.GoogleCallback)              // Google OAuth callback
		authRoutes.POST("logout", authController.Logout)                                    // Logout user
		authRoutes.POST("signup", authController.SignUp)                                    // Sign up user
		authRoutes.GET("status", middlewares.SessionMiddleware(), authController.CheckAuth) // Check authentication status
	}
}

// endregion

// region User Routes
func UserRoute(router *gin.Engine, userController controller.IUserController) {
	userRoutes := router.Group("/api/v1/user")
	{
		//userRoutes.PATCH("username", middlewares.SessionMiddleware(), userController.UpdateUsername)
		//userRoutes.PATCH("update-profile-photo", middlewares.CombinedAuthMiddleware(), userController.UploadProfilePhoto)

		userRoutes.PATCH("username", middlewares.SessionMiddleware(), userController.UpdateUsername)               // Update username
		userRoutes.PATCH("profile-photo", middlewares.CombinedAuthMiddleware(), userController.UploadProfilePhoto) // Update profile photo

	}
}

// endregion

// region Message Routes
func MessageRoute(router *gin.Engine, messageController controller.IMessageController) {
	messageRoutes := router.Group("/api/v1/messages")
	messageRoutes.Use(middlewares.SessionMiddleware())
	{
		messageRoutes.POST("history", messageController.GetMessageHistory) // Get chat history
	}
}

// endregion

// region Friend Routes
func FriendRoute(router *gin.Engine, friendController controller.IFriendController) {
	friendRoutes := router.Group("/api/v1/friends")
	friendRoutes.Use(middlewares.SessionMiddleware())
	{
		friendRoutes.GET("", friendController.GetFriends)             // Get all friends
		friendRoutes.GET("blocked", friendController.GetBlockedUsers) // Get blocked friends
		friendRoutes.PATCH("block", friendController.Block)           // Block a friend
		friendRoutes.DELETE("", friendController.Delete)              // Delete a friend
	}
}

// endregion

func RequestRoute(router *gin.Engine, requestController controller.IRequestController) {
	requestRoutes := router.Group("/api/v1/friend-requests")
	requestRoutes.Use(middlewares.SessionMiddleware())
	{
		requestRoutes.POST("", requestController.SendFriend) // Send a friend request
		requestRoutes.GET("", requestController.GetRequests) // Get incoming friend requests
		requestRoutes.PATCH("", requestController.Patch)     // Update friend request status

	}
}

// region Room Routes
func RoomRoute(router *gin.Engine, roomController controller.IRoomController) {
	roomRoutes := router.Group("/api/v1/rooms")
	roomRoutes.Use(middlewares.SessionMiddleware())
	{
		roomRoutes.POST("check", roomController.GetOrCreateRoom) // Check or create a room
		roomRoutes.GET("chat-list", roomController.GetChatList)  // Get chat list
	}
}

// endregion

// region File Routes
func FileRoute(router *gin.Engine, fileController controller.IFileController) {
	fileRoutes := router.Group("/api/v1/files")
	fileRoutes.Use(middlewares.SessionMiddleware())
	{
		fileRoutes.POST("upload", fileController.UploadFile) // Upload a file
	}
}

// endregion

// region Socket.IO Setup
func SetupSocketIO(router *gin.Engine, server *socket.Server, socketAdapter adapter.ISocketAdapter) {
	socketAdapter.HandleConnection()

	// Handle Socket.IO routes with session middleware
	router.GET("socket.io/*any", middlewares.SessionMiddleware(), gin.WrapH(server.ServeHandler(nil)))
	router.POST("socket.io/*any", middlewares.SessionMiddleware(), gin.WrapH(server.ServeHandler(nil)))
}

// endregion
