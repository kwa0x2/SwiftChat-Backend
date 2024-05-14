package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kwa0x2/realtime-chat-backend/service"
)

type ChatController struct {
	ChatService *service.ChatService
}

func (ctrl *ChatController) SessionTest(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}