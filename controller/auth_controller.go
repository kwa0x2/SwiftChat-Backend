package controller

import (
	"context"
	"encoding/json"

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

	if !ctrl.AuthService.IsIdUnique(userData["id"].(string)) {
		ctx.JSON(http.StatusAccepted, helpers.NewSuccessResponse(http.StatusAccepted, "Accepted", userData["id"].(string)))
		return
	}

	tokenString, err := helpers.GenerateToken(userData["id"].(string), userData["email"].(string), userData["picture"].(string))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helpers.NewErrorResponse(http.StatusInternalServerError, "Internal Server Error", "JWT Token Failed"))
		return
	}

	ctx.JSON(http.StatusOK, helpers.NewLoginResponse(http.StatusOK, "Continue", "Login successfully", tokenString))
}
