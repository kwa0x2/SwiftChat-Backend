package controller

import (
	"github.com/google/uuid"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kwa0x2/realtime-chat-backend/service"
	"github.com/kwa0x2/realtime-chat-backend/utils"
)

type IMessageController interface {
	GetMessageHistory(ctx *gin.Context)
}

type messageController struct {
	MessageService service.IMessageService
}

func NewMessageController(messageService service.IMessageService) IMessageController {
	return &messageController{
		MessageService: messageService,
	}
}

// region MessageHistoryBody defines the structure for the request body to get message history.
type MessageHistoryBody struct {
	RoomID uuid.UUID `json:"room_id"` // Unique identifier for the chat room.
}

// endregion

// region "GetMessageHistory" handles the request to retrieve message history for a specific room.
func (ctrl *messageController) GetMessageHistory(ctx *gin.Context) {
	var messageHistoryBody MessageHistoryBody

	// Bind JSON request body to the MessageHistoryBody struct.
	if err := ctx.BindJSON(&messageHistoryBody); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.NewErrorResponse("JSON Bind Error", err.Error()))
		return
	}

	// Retrieve message history data using the provided room ID.
	messageHistoryData, err := ctrl.MessageService.GetMessageHistoryByRoomID(messageHistoryBody.RoomID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.NewErrorResponse("Internal Server Error", "Error retrieving message history by room id."))
		return
	}

	ctx.JSON(http.StatusOK, utils.NewGetResponse(len(messageHistoryData), messageHistoryData))
}

// endregion
