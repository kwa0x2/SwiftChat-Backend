package di

import (
	"github.com/kwa0x2/swiftchat-backend/config"
	"github.com/kwa0x2/swiftchat-backend/controller"
	"github.com/kwa0x2/swiftchat-backend/repository"
	"github.com/kwa0x2/swiftchat-backend/service"
	"github.com/kwa0x2/swiftchat-backend/socket/adapter"
	"github.com/kwa0x2/swiftchat-backend/socket/gateway"
	"github.com/resend/resend-go/v2"
	"github.com/zishang520/socket.io/socket"
)

type Container struct {
	UserController    controller.IUserController
	AuthController    controller.IAuthController
	RoomController    controller.IRoomController
	MessageController controller.IMessageController
	FriendController  controller.IFriendController
	RequestController controller.IRequestController
	FileController    controller.IFileController
	SocketAdapter     adapter.ISocketAdapter
}

// region "NewContainer" initializes a new DI container, wiring up all dependencies.
func NewContainer(socketServer *socket.Server, resendClient *resend.Client) *Container {
	s3Service := service.NewS3Service()                     // S3 service for file storage
	resendService := service.NewResendService(resendClient) // Resend service for email handling

	userRepository := repository.NewUserRepository(config.DB) // User repository for data access
	userService := service.NewUserService(userRepository)     // User service for business logic

	userRoomRepository := repository.NewUserRoomRepository(config.DB) // User-Room repository for data access
	userRoomService := service.NewUserRoomService(userRoomRepository) // User-Room service for business logic

	roomRepository := repository.NewRoomRepository(config.DB)              // Room repository for data access
	roomService := service.NewRoomService(roomRepository, userRoomService) // Room service for business logic

	messageRepository := repository.NewMessageRepository(config.DB)             // Message repository for data access
	messageService := service.NewMessageService(messageRepository, roomService) // Message service for business logic

	friendRepository := repository.NewFriendRepository(config.DB) // Friend repository for data access
	friendService := service.NewFriendService(friendRepository)   // Friend service for business logic

	requestRepository := repository.NewRequestRepository(config.DB)                            // Request repository for data access
	requestService := service.NewRequestService(requestRepository, friendService, userService) // Request service for business logic

	socketGateway := gateway.NewSocketGateway(socketServer, "/chat")                                        // Initialize the socket gateway for handling socket connections
	socketAdapter := adapter.NewSocketAdapter(socketGateway, messageService, friendService, requestService) // Socket adapter for emitting events

	// Return a new Container with all initialized controllers and the socket adapter
	return &Container{
		UserController:    controller.NewUserController(userService, friendService, s3Service, socketAdapter),
		AuthController:    controller.NewAuthController(userService),
		RoomController:    controller.NewRoomController(roomService, userRoomService, userService, friendService),
		MessageController: controller.NewMessageController(messageService),
		FriendController:  controller.NewFriendController(friendService, socketGateway),
		RequestController: controller.NewRequestController(requestService, friendService, userService, socketGateway, resendService),
		FileController:    controller.NewFileController(s3Service),
		SocketAdapter:     socketAdapter,
	}
}

// endregion
