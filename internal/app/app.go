package app

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/kwa0x2/realtime-chat-backend/config"
	"github.com/kwa0x2/realtime-chat-backend/internal/di"
	"github.com/kwa0x2/realtime-chat-backend/routes"
	"github.com/resend/resend-go/v2"
	"github.com/zishang520/socket.io/socket"
	"os"
)

type App struct {
	Router       *gin.Engine    // Gin router for handling HTTP requests
	Socket       *socket.Server // Socket.IO server for real-time communication
	ResendClient *resend.Client // Resend client for sending emails
}

// region "NewApp" initializes a new App instance and configures the necessary components.
func NewApp() *App {
	config.LoadEnv()                                              // Load environment variables and configure services
	config.PostgreConnection()                                    // Initialize PostgreSQL connection
	config.InitS3()                                               // Initialize S3 storage
	router := gin.New()                                           // Create a new Gin engine
	socketServer := socket.NewServer(nil, nil)                    // Create a new Socket.IO server
	resendClient := resend.NewClient(os.Getenv("RESEND_API_KEY")) // Initialize the Resend client with the API key from environment variables
	store := config.RedisSession()                                // Initialize Redis session store

	// Middleware for sessions and CORS
	router.Use(sessions.Sessions("connect.sid", store)) // Use session management middleware
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},                            // Allow requests from this origin
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"}, // Allowed HTTP methods
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},          // Allowed headers
		ExposeHeaders:    []string{"Content-Length"},                                   // Exposed headers to the client
		AllowCredentials: true,                                                         // Allow credentials in requests
	}))

	return &App{ // Return a new App instance with the configured components
		Router:       router,
		Socket:       socketServer,
		ResendClient: resendClient,
	}
}

// endregion

// region "SetupRoutes" initializes the application's routes and associates them with the controllers.
func (a *App) SetupRoutes() {
	container := di.NewContainer(a.Socket, a.ResendClient) // Create a new dependency injection container

	// Setup routes for various controllers
	routes.UserRoute(a.Router, container.UserController)
	routes.AuthRoute(a.Router, container.AuthController)
	routes.MessageRoute(a.Router, container.MessageController)
	routes.FriendRoute(a.Router, container.FriendController)
	routes.RequestRoute(a.Router, container.RequestController)
	routes.RoomRoute(a.Router, container.RoomController)
	routes.SetupSocketIO(a.Router, a.Socket, container.SocketAdapter) // Setup Socket.IO routes
}

// endregion

// region "Run" starts the HTTP server on the specified port.
func (a *App) Run() error {
	return a.Router.Run(":9000") // Start the server on port 9000
}

// endregion
