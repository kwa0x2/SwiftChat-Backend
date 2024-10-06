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
	userService *service.UserService
}

func NewAuthController(userService *service.UserService) IAuthController {
	return &authController{
		userService: userService,
	}
}

var (
	stateStore = sync.Map{}
)

func (ctrl *authController) GoogleLogin(ctx *gin.Context) {
	googleConfig := config.GoogleConfig()
	state := uuid.New().String()
	stateStore.Store(state, state)
	url := googleConfig.AuthCodeURL(state)
	ctx.Redirect(http.StatusTemporaryRedirect, url)
}

func (ctrl *authController) GoogleCallback(ctx *gin.Context) {
	code := ctx.Query("code")
	state := ctx.Query("state")

	if _, exists := stateStore.Load(state); !exists {
		ctx.JSON(http.StatusBadRequest, utils.NewErrorResponse("Bad Request", "Invalid state parameter. Please try again"))
		return
	}

	googleConfig := config.GoogleConfig()

	token, err := googleConfig.Exchange(context.Background(), code)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.NewErrorResponse("Internal Server Error", "Code-Token Exchange Failed"))
		return
	}

	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
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

	// id unique degilse
	if !ctrl.userService.IsIdUnique(userData["id"].(string)) {
		user, err := ctrl.userService.GetUserById(userData["id"].(string))
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, utils.NewErrorResponse("Internal Server Error", "Failed to retrieve user by ID"))
			return
		}
		session := sessions.Default(ctx)
		session.Set("id", userData["id"].(string))
		session.Set("name", user.UserName)
		session.Set("email", userData["email"].(string))
		session.Set("photo", user.UserPhoto)
		session.Set("role", user.UserRole)
		session.Save()

		ctx.Redirect(http.StatusTemporaryRedirect, "http://localhost:3000/login")
		return
	}

	jwtClaims := jwt.MapClaims{
		"id":         userData["id"].(string),
		"user_email": userData["email"].(string),
		"user_photo": userData["picture"].(string),
		"user_name":  userData["name"].(string),
		"exp":        time.Now().Add(time.Hour * 2).Unix(),
	}

	tokenString, err := utils.GenerateToken(jwtClaims)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.NewErrorResponse("Internal Server Error", "JWT Token Failed"))
		return
	}

	ctx.Redirect(http.StatusTemporaryRedirect, "http://localhost:3000/createname?token="+tokenString)
}

func (ctrl *authController) CheckAuth(ctx *gin.Context) {
	session := sessions.Default(ctx)

	ctx.JSON(http.StatusOK, gin.H{
		"id":    session.Get("id"),
		"name":  session.Get("name"),
		"email": session.Get("email"),
		"photo": session.Get("photo"),
		"role":  session.Get("role"),
	})
}

func (ctrl *authController) Logout(ctx *gin.Context) {
	session := sessions.Default(ctx)

	session.Clear()
	session.Options(sessions.Options{MaxAge: -1})
	session.Save()

	ctx.SetCookie("connect.sid", "", -1, "/", "localhost", true, true)
}

type SignUpBody struct {
	UserName  string `json:"user_name"`
	UserPhoto string `json:"user_photo"`
}

func (ctrl *authController) SignUp(ctx *gin.Context) {
	var signUpBody SignUpBody
	if err := ctx.BindJSON(&signUpBody); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.NewErrorResponse("JSON Bind Error", err.Error()))
		return
	}

	claims, err := utils.GetClaims(ctx.GetHeader("Authorization"))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.NewErrorResponse("Internal Server Error", "Failed to retrieve JWT claims"))
		return
	}

	userInsertObj := models.User{
		UserID:    claims["id"].(string),
		UserEmail: claims["user_email"].(string),
		UserName:  signUpBody.UserName,
		UserPhoto: signUpBody.UserPhoto,
	}

	session := sessions.Default(ctx)
	session.Set("id", claims["id"].(string))
	session.Set("name", signUpBody.UserName)
	session.Set("email", claims["user_email"].(string))
	session.Set("photo", signUpBody.UserPhoto)
	session.Set("role", "user")
	session.Save()

	if !ctrl.userService.IsUsernameUnique(signUpBody.UserName) {
		ctx.JSON(http.StatusConflict, utils.NewErrorResponse("Conflict", "Username is already taken, must be unique"))
		return
	}

	userData, err := ctrl.userService.Insert(&userInsertObj)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.NewErrorResponse("Internal Server Error", "Failed to insert new user into database"))
		return
	}

	ctx.JSON(http.StatusCreated, utils.NewGetResponse(1, userData))
}
