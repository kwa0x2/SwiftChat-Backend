package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kwa0x2/realtime-chat-backend/service"
	"github.com/kwa0x2/realtime-chat-backend/utils"
)

type IMessageController interface {
	GetMessageHistory(ctx *gin.Context)
}

type messageController struct {
	messageService *service.MessageService
}

func NewMessageController(messageService *service.MessageService) IMessageController {
	return &messageController{
		messageService: messageService,
	}
}

type MessageHistoryBody struct {
	RoomID string `json:"room_id"`
}

func (ctrl *messageController) GetMessageHistory(ctx *gin.Context) {
	var messageHistoryBody MessageHistoryBody

	if err := ctx.BindJSON(&messageHistoryBody); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.NewErrorResponse("JSON Bind Error", err.Error()))
		return
	}

	messageHistoryData, err := ctrl.messageService.GetMessageHistoryByRoomID(messageHistoryBody.RoomID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.NewErrorResponse("Internal Server Error", "Error retrieving message history by room id."))
		return
	}

	ctx.JSON(http.StatusOK, utils.NewGetResponse(len(messageHistoryData), messageHistoryData))
}
