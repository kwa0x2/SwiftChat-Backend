package main

import (
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/kwa0x2/realtime-chat-backend/config"
	"github.com/kwa0x2/realtime-chat-backend/controller"
	"github.com/kwa0x2/realtime-chat-backend/middlewares"
	"github.com/kwa0x2/realtime-chat-backend/repository"
	"github.com/kwa0x2/realtime-chat-backend/routes"
	"github.com/kwa0x2/realtime-chat-backend/service"
	"github.com/zishang520/engine.io/utils"
	"github.com/zishang520/socket.io/socket"
)

func main() {
	config.LoadEnv()
	router := gin.New()
	config.PostgreConnection()
	io := socket.NewServer(nil, nil)
	store := config.RedisSession()
	router.Use(sessions.Sessions("connect.sid", store))

	io.Of("/chat",nil).On("connection", func(clients ...any) {
		socket := clients[0].(*socket.Socket)

		utils.Log().Info(`socket %s connected`, socket.Id())

		socket.Emit("foo", "bar")

		socket.On("disconnect", func(reason ...any) {
			utils.Log().Info(`socket %s disconnected due to %s`, socket.Id(), reason[0])
		})
	})

	router.GET("socket.io/*any", middlewares.SessionMiddleware(), gin.WrapH(io.ServeHandler(nil)))
	router.POST("socket.io/*any", middlewares.SessionMiddleware(), gin.WrapH(io.ServeHandler(nil)))
	router.StaticFS("/public", http.Dir("./asset"))

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	userRepository := &repository.UserRepository{DB: config.DB}
	userService := &service.UserService{UserRepository: userRepository}
	userController := &controller.UserController{UserService: userService}

	authRepository := &repository.AuthRepository{DB: config.DB}
	authService := &service.AuthService{AuthRepository: authRepository}
	authController := &controller.AuthController{AuthService: authService}

	routes.UserRoute(router, userController)
	routes.AuthRoute(router, authController)

	if err := router.Run(":9000"); err != nil {
		log.Fatal("failed run app: ", err)
	}
}
