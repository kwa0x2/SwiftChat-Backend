package main

import (
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
	config.Connection()

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


	routes.UserRoute(router, userController)
	routes.AuthRoute(router, authController)

	router.Run(":9000")
}
