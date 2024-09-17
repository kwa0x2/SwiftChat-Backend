package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kwa0x2/realtime-chat-backend/service"
	"github.com/kwa0x2/realtime-chat-backend/utils"
)

type MessageController struct {
	MessageService *service.MessageService
}

type PrivateConversationBody struct {
	MessageSenderID   string `json:"message_sender_id"`
	MessageReceiverID string `json:"message_receiver_id"`
}

func (ctrl *MessageController) GetPrivateConversation(ctx *gin.Context) {
	var privateConversationBody PrivateConversationBody

	if err := ctx.BindJSON(&privateConversationBody); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.NewErrorResponse(http.StatusBadRequest, "Bad Request", err.Error()))
		return
	}

	data, err := ctrl.MessageService.GetPrivateConversation(privateConversationBody.MessageSenderID, privateConversationBody.MessageReceiverID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.NewErrorResponse(http.StatusInternalServerError, "Internal Server Error", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, utils.NewGetResponse(http.StatusOK, "OK", len(data), data))
}

type MessageHistoryBody struct {
	RoomID string `json:"room_id"`
}

func (ctrl *MessageController) GetMessageHistory(ctx *gin.Context) {
	var messageHistoryBody MessageHistoryBody

	if err := ctx.BindJSON(&messageHistoryBody); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.NewErrorResponse(http.StatusBadRequest, "Bad Request", err.Error()))
		return
	}

	data, err := ctrl.MessageService.GetMessageHistoryByRoomID(messageHistoryBody.RoomID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.NewErrorResponse(http.StatusInternalServerError, "Internal Server Error", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, data)
}
