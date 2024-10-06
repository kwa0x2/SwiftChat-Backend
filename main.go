package main

import (
	"encoding/gob"
	"github.com/kwa0x2/realtime-chat-backend/socket/adapter"
	"github.com/kwa0x2/realtime-chat-backend/socket/gateway"
	"github.com/resend/resend-go/v2"
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/kwa0x2/realtime-chat-backend/config"
	"github.com/kwa0x2/realtime-chat-backend/controller"
	"github.com/kwa0x2/realtime-chat-backend/repository"
	"github.com/kwa0x2/realtime-chat-backend/routes"
	"github.com/kwa0x2/realtime-chat-backend/service"
	"github.com/zishang520/socket.io/socket"
)

func init() {
	gob.Register(time.Time{})
}
func main() {
	config.LoadEnv()
	router := gin.New()
	config.PostgreConnection()
	config.InitS3()
	socketServer := socket.NewServer(nil, nil)
	store := config.RedisSession()
	router.Use(sessions.Sessions("connect.sid", store))
	resendClient := resend.NewClient(os.Getenv("RESEND_API_KEY"))

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	s3Service := &service.S3Service{}
	resendService := &service.ResendService{ResendClient: resendClient}

	userRepository := &repository.UserRepository{DB: config.DB}
	userService := &service.UserService{UserRepository: userRepository}

	userRoomRepository := &repository.UserRoomRepository{DB: config.DB}
	userRoomService := &service.UserRoomService{UserRoomRepository: userRoomRepository}

	roomRepository := &repository.RoomRepository{DB: config.DB}
	roomService := &service.RoomService{RoomRepository: roomRepository, UserRoomService: userRoomService}

	messageRepository := &repository.MessageRepository{DB: config.DB}
	messageService := &service.MessageService{MessageRepository: messageRepository, RoomService: roomService}

	friendRepository := &repository.FriendRepository{DB: config.DB}
	friendService := &service.FriendService{FriendRepository: friendRepository}

	requestRepository := &repository.RequestRepository{DB: config.DB}
	requestService := &service.RequestService{RequestRepository: requestRepository, FriendService: friendService, UserService: userService}

	socketGateway := &gateway.SocketGateway{Server: socketServer}
	socketAdapter := &adapter.SocketAdapter{Gateway: socketGateway, MessageService: messageService, FriendService: friendService}

	userController := controller.NewUserController(userService, friendService, s3Service, socketAdapter)
	authController := controller.NewAuthController(userService)
	roomController := controller.NewRoomController(roomService, userRoomService, userService, friendService)
	messageController := controller.NewMessageController(messageService)
	friendController := controller.NewFriendController(friendService, userService, requestService, socketAdapter)
	requestController := controller.NewRequestController(requestService, friendService, userService, socketAdapter, resendService)

	routes.UserRoute(router, userController)
	routes.AuthRoute(router, authController)
	routes.MessageRoute(router, messageController)
	routes.FriendRoute(router, friendController)
	routes.RequestRoute(router, requestController)
	routes.RoomRoute(router, roomController)
	routes.SetupSocketIO(router, socketServer, socketAdapter)

	if err := router.Run(":9000"); err != nil {
		log.Fatal("failed run app: ", err)
	}
}
