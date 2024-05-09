package controller

import (

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kwa0x2/realtime-chat-backend/helpers"
	"github.com/kwa0x2/realtime-chat-backend/models"
	"github.com/kwa0x2/realtime-chat-backend/service"
)


type UserController struct {
	UserService *service.UserService
}

type InsertUser struct {
	Username string `json:"username"`
}

func (ctrl *UserController) InsertUser(ctx *gin.Context) {
	var insertUser InsertUser
	if err := ctx.BindJSON(&insertUser); err != nil {
		ctx.JSON(http.StatusBadRequest, helpers.NewErrorResponse(http.StatusBadRequest, "Bad Request", err.Error()))
		return
	}

	claims, err := helpers.GetClaims(ctx.GetHeader("Authorization"))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helpers.NewErrorResponse(http.StatusInternalServerError, "Internal Server Error", err.Error()))
		return
	}

	user := models.User{
		UserID:    claims["id"].(string),
		UserEmail: claims["email"].(string),
		UserName:  insertUser.Username,
		UserPhoto: claims["photo"].(string),
	}

	if !ctrl.UserService.IsUsernameUnique(insertUser.Username) {
		ctx.JSON(http.StatusConflict, helpers.NewErrorResponse(http.StatusConflict, "Conflict", "Username must be unique"))
		return
	}

	if err := ctrl.UserService.InsertUser(&user); err != nil {
		ctx.JSON(http.StatusInternalServerError, helpers.NewErrorResponse(http.StatusInternalServerError, "Internal Server Error", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, helpers.NewSuccessResponse(http.StatusOK, "OK", "Succesfully inserted"))
}
