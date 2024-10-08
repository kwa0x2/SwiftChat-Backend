package controller

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/kwa0x2/realtime-chat-backend/models"
	"sync"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-contrib/sessions"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kwa0x2/realtime-chat-backend/config"
	"github.com/kwa0x2/realtime-chat-backend/service"
	"github.com/kwa0x2/realtime-chat-backend/utils"
)

type IAuthController interface {
	GoogleLogin(ctx *gin.Context)
	GoogleCallback(ctx *gin.Context)
	CheckAuth(ctx *gin.Context)
	Logout(ctx *gin.Context)
	SignUp(ctx *gin.Context)
}

type authController struct {
	UserService service.IUserService
}

func NewAuthController(UserService service.IUserService) IAuthController {
	return &authController{
		UserService: UserService,
	}
}

// In-memory storage for state during Google login
var (
	stateStore = sync.Map{}
)

// region "GoogleLogin" starts the Google authentication flow
func (ctrl *authController) GoogleLogin(ctx *gin.Context) {
	googleConfig := config.GoogleConfig()           // Get Google OAuth configuration
	state := uuid.New().String()                    // Generate a unique state
	stateStore.Store(state, state)                  // Store state in memory
	url := googleConfig.AuthCodeURL(state)          // Generate the authorization URL
	ctx.Redirect(http.StatusTemporaryRedirect, url) // Redirect user to Google for login
}

// endregion

// region "GoogleCallback" handles the redirect from Google after login
func (ctrl *authController) GoogleCallback(ctx *gin.Context) {
	code := ctx.Query("code")   // Get the authorization code from the query parameters
	state := ctx.Query("state") // Get the state parameter from the query parameters

	// Validate the state parameter to prevent CSRF attacks
	if _, exists := stateStore.Load(state); !exists {
		ctx.JSON(http.StatusBadRequest, utils.NewErrorResponse("Bad Request", "Invalid state parameter. Please try again"))
		return
	}

	googleConfig := config.GoogleConfig() // Get Google OAuth configuration

	// Exchange the authorization code for a token
	token, err := googleConfig.Exchange(context.Background(), code)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.NewErrorResponse("Internal Server Error", "Code-Token Exchange Failed"))
		return
	}

	// Fetch user information from Google
	resp, respErr := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if respErr != nil {
		ctx.JSON(http.StatusInternalServerError, utils.NewErrorResponse("Internal Server Error", "User data fetch failed"))
		return
	}
	defer resp.Body.Close()

	var userData map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&userData)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.NewErrorResponse("Internal Server Error", "JSON Parsing Failed"))
		return
	}

	// Check if the user ID from Google is unique in our system
	if !ctrl.UserService.IsIdUnique(userData["id"].(string)) {
		// If the user ID is not unique, it means the user already exists in our database
		user, getUserErr := ctrl.UserService.GetUserById(userData["id"].(string))
		if getUserErr != nil {
			ctx.JSON(http.StatusInternalServerError, utils.NewErrorResponse("Internal Server Error", "Failed to retrieve user by ID"))
			return
		}

		// If the user exists, save their data in the session
		session := sessions.Default(ctx)
		session.Set("id", userData["id"].(string))
		session.Set("name", user.UserName)
		session.Set("email", userData["email"].(string))
		session.Set("photo", user.UserPhoto)
		session.Set("role", user.UserRole)
		session.Save()

		// Redirect the user to the login page since they already exist in the system
		ctx.Redirect(http.StatusTemporaryRedirect, "http://localhost:3000/login")
		return
	}

	// If the user ID is unique, we proceed to create JWT claims for the new user
	jwtClaims := jwt.MapClaims{
		"id":         userData["id"].(string),
		"user_email": userData["email"].(string),
		"user_photo": userData["picture"].(string),
		"user_name":  userData["name"].(string),
		"exp":        time.Now().Add(time.Hour * 2).Unix(),
	}

	// Generate JWT token
	tokenString, tokenErr := utils.GenerateToken(jwtClaims)
	if tokenErr != nil {
		ctx.JSON(http.StatusInternalServerError, utils.NewErrorResponse("Internal Server Error", "JWT Token Failed"))
		return
	}

	// Redirect to the create name page and include the JWT token as a query parameter
	ctx.Redirect(http.StatusTemporaryRedirect, "http://localhost:3000/createname?token="+tokenString)
}

// endregion

// region "CheckAuth" checks the authentication status and returns user session data
func (ctrl *authController) CheckAuth(ctx *gin.Context) {
	session := sessions.Default(ctx)

	// Return user data in JSON format
	ctx.JSON(http.StatusOK, gin.H{
		"id":    session.Get("id"),
		"name":  session.Get("name"),
		"email": session.Get("email"),
		"photo": session.Get("photo"),
		"role":  session.Get("role"),
	})
}

// endregion

// region "Logout" logs out the user and clears the session
func (ctrl *authController) Logout(ctx *gin.Context) {
	session := sessions.Default(ctx)

	// Clear the session
	session.Clear()
	session.Options(sessions.Options{MaxAge: -1}) // Set session age to -1 to delete it
	session.Save()                                // Save session changes

	// Clear the cookie for session identification
	ctx.SetCookie("connect.sid", "", -1, "/", "localhost", true, true)
}

// endregion

// region "SignUpBody" struct for parsing sign-up requests
type SignUpBody struct {
	UserName  string `json:"user_name"`
	UserPhoto string `json:"user_photo"`
}

// endregion

// region "SignUp" registers a new user
func (ctrl *authController) SignUp(ctx *gin.Context) {
	var signUpBody SignUpBody
	if err := ctx.BindJSON(&signUpBody); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.NewErrorResponse("JSON Bind Error", err.Error()))
		return
	}

	// Retrieve JWT claims from the authorization header
	claims, err := utils.GetClaims(ctx.GetHeader("Authorization"))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.NewErrorResponse("Internal Server Error", "Failed to retrieve JWT claims"))
		return
	}

	// Create a user object for insertion
	userInsertObj := models.User{
		UserID:    claims["id"].(string),
		UserEmail: claims["user_email"].(string),
		UserName:  signUpBody.UserName,
		UserPhoto: signUpBody.UserPhoto,
	}

	// Save user data in session
	session := sessions.Default(ctx)
	session.Set("id", claims["id"].(string))
	session.Set("name", signUpBody.UserName)
	session.Set("email", claims["user_email"].(string))
	session.Set("photo", signUpBody.UserPhoto)
	session.Set("role", "user")
	session.Save()

	// Check if the username is unique
	if !ctrl.UserService.IsUsernameUnique(signUpBody.UserName) {
		ctx.JSON(http.StatusConflict, utils.NewErrorResponse("Conflict", "Username is already taken, must be unique"))
		return
	}

	// Insert the new user into the database
	userData, createErr := ctrl.UserService.Create(&userInsertObj)
	if createErr != nil {
		ctx.JSON(http.StatusInternalServerError, utils.NewErrorResponse("Internal Server Error", "Failed to insert new user into database"))
		return
	}

	ctx.JSON(http.StatusCreated, userData)
}

// endregion
