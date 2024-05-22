package controller

import (
	"context"
	"encoding/json"
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
		session := sessions.Default(ctx)
		session.Set("user_id", userData["id"].(string))
		session.Set("user_authority_id", 2)
		session.Save()
		ctx.JSON(http.StatusAccepted, helpers.NewLoginResponse(http.StatusAccepted, "Accepted", "Login successfully"))
		return
	}

	jwtClaims := jwt.MapClaims{
		"id":    userData["id"].(string),
		"email": userData["email"].(string),
		"photo": userData["picture"].(string),
		"exp":   time.Now().Add(time.Hour * 2).Unix(),
	}

	tokenString, err := helpers.GenerateToken(jwtClaims)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helpers.NewErrorResponse(http.StatusInternalServerError, "Internal Server Error", "JWT Token Failed"))
		return
	}

	ctx.JSON(http.StatusOK, helpers.NewSignUpResponse(http.StatusOK, "Continue", "Login successfully", tokenString))
}
