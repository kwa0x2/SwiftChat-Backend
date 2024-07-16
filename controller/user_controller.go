package controller

import (
	"fmt"
	"net/http"

	"github.com/gin-contrib/sessions"
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
		UserEmail: claims["user_email"].(string),
		UserName:  insertUser.Username,
		UserPhoto: claims["user_photo"].(string),
	}

	session := sessions.Default(ctx)
	session.Set("id", claims["id"].(string))
	session.Set("name", claims["user_name"].(string))
	session.Set("mail", claims["user_email"].(string))
	session.Set("photo", claims["user_photo"].(string))
	session.Set("role", "user")
	session.Save()

	fmt.Println(session.ID())

	if !ctrl.UserService.IsUsernameUnique(insertUser.Username) {
		ctx.JSON(http.StatusConflict, helpers.NewErrorResponse(http.StatusConflict, "Conflict", "Username must be unique"))
		return
	}

	userdata,err := ctrl.UserService.Insert(&user)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helpers.NewErrorResponse(http.StatusInternalServerError, "Internal Server Error", err.Error()))
		return
	}

	ctx.JSON(http.StatusCreated, userdata)
}

func (ctrl *UserController) GetAll(ctx *gin.Context) {
	usersData, err := ctrl.UserService.GetAll()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helpers.NewErrorResponse(http.StatusInternalServerError, "Internal Server Error", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, helpers.NewGetResponse(http.StatusOK, "OK", len(usersData), usersData))
}

func (ctrl *UserController) GetByEmail(ctx *gin.Context) {
	usersData, err := ctrl.UserService.GetByEmail("asdasdsa@gmail.com")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helpers.NewErrorResponse(http.StatusInternalServerError, "Internal Server Error", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, helpers.NewGetResponse(http.StatusOK, "OK", 1, usersData))
}
