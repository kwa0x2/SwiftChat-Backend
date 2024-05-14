package main

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/kwa0x2/realtime-chat-backend/config"
	"github.com/kwa0x2/realtime-chat-backend/controller"
	"github.com/kwa0x2/realtime-chat-backend/repository"
	"github.com/kwa0x2/realtime-chat-backend/routes"
	"github.com/kwa0x2/realtime-chat-backend/service"
)

func main() {
	config.LoadEnv()
	router := gin.New()
	config.PostgreConnection()
	store := config.RedisConnection()
	if store == nil {
		panic("Redis bağlantısı başarısız")
	}
	
	router.Use(sessions.Sessions("mysession", store))
	
	userRepository := &repository.UserRepository{
		DB: config.DB,
	}

	userService := &service.UserService{
		UserRepository: userRepository,
	}

	userController := &controller.UserController{
		UserService: userService,
	}

	authRepository := &repository.AuthRepository{
		DB: config.DB,
	}

	authService := &service.AuthService{
		AuthRepository: authRepository,
	}

	authController := &controller.AuthController{
		AuthService: authService,
	}


	chatRepository := &repository.ChatRepository{
		DB: config.DB,
	}

	chatService := &service.ChatService{
		ChatRepository: chatRepository,
	}

	chatController := &controller.ChatController{
		ChatService: chatService,
	}

	routes.UserRoute(router, userController)
	routes.AuthRoute(router, authController)
	routes.ChatRoute(router, chatController)

	router.Run(":9000")
}
