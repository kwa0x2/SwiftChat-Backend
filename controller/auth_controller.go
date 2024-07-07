package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-contrib/sessions"
	"github.com/google/uuid"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kwa0x2/realtime-chat-backend/config"
	"github.com/kwa0x2/realtime-chat-backend/helpers"
	"github.com/kwa0x2/realtime-chat-backend/service"
)

type AuthController struct {
	AuthService *service.AuthService
	State       string
}

func (ctrl *AuthController) GoogleLogin(ctx *gin.Context) {
	googleConfig := config.GoogleConfig()
	ctrl.State = uuid.New().String()
	url := googleConfig.AuthCodeURL(ctrl.State)
	ctx.Redirect(http.StatusTemporaryRedirect, url)
}

func (ctrl *AuthController) GoogleCallback(ctx *gin.Context) {
	expectedState := ctx.Query("state")
	if expectedState != ctrl.State {
		ctx.String(http.StatusBadRequest, "States don't Match!!!")
		return
	}

	code := ctx.Query("code")

	googleConfig := config.GoogleConfig()

	token, err := googleConfig.Exchange(context.Background(), code)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helpers.NewErrorResponse(http.StatusInternalServerError, "Internal Server Error", "Code-Token Exchange Failed"))
		return
	}

	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helpers.NewErrorResponse(http.StatusInternalServerError, "Internal Server Error", "User data fetch failed"))
		return
	}
	defer resp.Body.Close()

	var userData map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&userData)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helpers.NewErrorResponse(http.StatusInternalServerError, "Internal Server Error", "JSON Parsing Failed"))
		return
	}

	// id unique degilse
	if !ctrl.AuthService.IsIdUnique(userData["id"].(string)) {

		username := ctrl.AuthService.GetUserName(userData["id"].(string))

		session := sessions.Default(ctx)
		session.Set("id", userData["id"].(string))
		session.Set("name", username)
		session.Set("mail", userData["email"].(string))
		session.Set("photo", userData["picture"].(string))
		session.Set("role", "user")
		session.Save()
		fmt.Println(session.ID())

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

	tokenString, err := helpers.GenerateToken(jwtClaims)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helpers.NewErrorResponse(http.StatusInternalServerError, "Internal Server Error", "JWT Token Failed"))
		return
	}

	ctx.Redirect(http.StatusTemporaryRedirect, "http://localhost:3000/createname?token="+tokenString)
}

func (ctrl *AuthController) CheckAuth(ctx *gin.Context) {
	session := sessions.Default(ctx)

	ctx.JSON(http.StatusOK, gin.H{
		"id":    session.Get("id"),
		"name":  session.Get("name"),
		"mail":  session.Get("mail"),
		"photo": session.Get("photo"),
		"role":  session.Get("role"),
	})
}

func (ctrl *AuthController) Logout(ctx *gin.Context) {
	session := sessions.Default(ctx)

	session.Clear()
	session.Options(sessions.Options{MaxAge: -1})
	session.Save()
	
	ctx.SetCookie("connect.sid","",-1,"/","localhost",true,true)
}