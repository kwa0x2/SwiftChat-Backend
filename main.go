package main

import (	"log"

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

var userSockets = make(map[string]string)

func main() {
	config.LoadEnv()
	router := gin.New()
	config.PostgreConnection()
	io := socket.NewServer(nil, nil)
	store := config.RedisSession()
	router.Use(sessions.Sessions("connect.sid", store))

	// io.Of("/chat", nil).On("connection", func(clients ...any) {
	// 	socketio := clients[0].(*socket.Socket)
	// 	ctx := socketio.Request().Context()

	// 	utils.Log().Info(`socket connected %s user id %s`, socketio.Id(), ctx.Value("id").(string))

	// 	userSockets[ctx.Value("id").(string)] = string(socketio.Id())

	// 	socketio.On("disconnect", func(reason ...any) {
	// 		utils.Log().Info(`socket %s disconnected due to %s`, socketio.Id(), reason[0])
	// 	})

	// 	socketio.On("sendMessage", func(args ...interface{}) {

	// 		msg, ok := args[0].(map[string]interface{})
	// 		if !ok {
	// 			utils.Log().Error("Failed to parse sendMessage event data")
	// 			return
	// 		}

	// 		destionation_id := msg["DestionationUserId"].(string)
	// 		message := msg["Message"].(string)
	// 		fmt.Println(userSockets)

	// 		// database insert



	// 		utils.Log().Info("Received message from %s : %s", ctx.Value("id").(string), message)

	// 		io.Of("/chat", nil).To(socket.Room(userSockets[destionation_id])).Emit("chat", map[string]interface{}{
	// 			"sender_id": ctx.Value("id").(string),
	// 			"message":   message,
	// 		})
	// 	})
	// })

	// router.GET("socket.io/*any", middlewares.SessionMiddleware(), gin.WrapH(io.ServeHandler(nil)))
	// router.POST("socket.io/*any", middlewares.SessionMiddleware(), gin.WrapH(io.ServeHandler(nil)))
	// router.StaticFS("/public", http.Dir("./asset"))

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
	routes.SetupSocketIO(router, io)

	if err := router.Run(":9000"); err != nil {
		log.Fatal("failed run app: ", err)
	}
}
