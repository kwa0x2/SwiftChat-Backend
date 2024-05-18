package controller

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/kwa0x2/realtime-chat-backend/service"
)

type ChatController struct {
	ChatService *service.ChatService
}

func (ctrl *ChatController) Login(ctx *gin.Context) {
	session := sessions.Default(ctx)

	session.Set("user_id",11111)
	session.Set("user_email","asdadsa@gmail.com")
	session.Set("user_auth",1)
	session.Save()


	ctx.JSON(http.StatusOK, gin.H{
		"status": "true",
	})
}

func (ctrl *ChatController) Auth(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"message": "pong",
	})
}