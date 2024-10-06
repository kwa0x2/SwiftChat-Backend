package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/kwa0x2/realtime-chat-backend/service"
	"github.com/kwa0x2/realtime-chat-backend/types"
	"github.com/kwa0x2/realtime-chat-backend/utils"
	"net/http"
)

type IRoomController interface {
	GetOrCreatePrivateRoom(ctx *gin.Context)
	GetChatList(ctx *gin.Context)
}

type roomController struct {
	roomService     *service.RoomService
	userRoomService *service.UserRoomService
	userService     *service.UserService
	friendService   *service.FriendService
}

func NewRoomController(roomService *service.RoomService, userRoomService *service.UserRoomService, userService *service.UserService, friendService *service.FriendService) IRoomController {
	return &roomController{
		roomService:     roomService,
		userRoomService: userRoomService,
		userService:     userService,
		friendService:   friendService,
	}
}

type ActionBody struct {
	Email  string              `json:"email"`
	Status types.RequestStatus `json:"status"`
}

func (ctrl *roomController) GetOrCreatePrivateRoom(ctx *gin.Context) {
	var actionBody ActionBody

	if err := ctx.BindJSON(&actionBody); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.NewErrorResponse("JSON Bind Error", err.Error()))
		return
	}

	userSessionInfo, sessionErr := utils.GetUserSessionInfo(ctx)
	if sessionErr != nil {
		ctx.JSON(http.StatusBadRequest, utils.NewErrorResponse("Session Error", sessionErr.Error()))
		return
	}

	user, userErr := ctrl.userService.GetByEmail(actionBody.Email)
	if userErr != nil {
		ctx.JSON(http.StatusInternalServerError, utils.NewErrorResponse("Internal Server Errors", "Error retrieving user by email"))
		return
	}

	room, roomErr := ctrl.userRoomService.GetPrivateRoom(userSessionInfo.ID, user.UserID)
	if roomErr != nil {
		ctx.JSON(http.StatusInternalServerError, utils.NewErrorResponse("Internal Server Error", "Error retrieving private room"))
		return
	}

	var roomId string

	if room == "" {
		newRoomId, newRoomErr := ctrl.roomService.CreateAndAddUsers(userSessionInfo.ID, user.UserID, "private")
		if newRoomErr != nil {
			ctx.JSON(http.StatusInternalServerError, utils.NewErrorResponse("Internal Server Error", "Error creating and adding users to room"))
			return
		}
		roomId = newRoomId
	} else {
		roomId = room
	}

	ctx.JSON(http.StatusOK, gin.H{
		"room_id": roomId,
	})
}

func (ctrl *roomController) GetChatList(ctx *gin.Context) {
	userSessionInfo, err := utils.GetUserSessionInfo(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.NewErrorResponse("Session Error", err.Error()))
		return
	}

	chatListData, chatListErr := ctrl.roomService.GetChatList(userSessionInfo.ID, userSessionInfo.Email)
	if chatListErr != nil {
		ctx.JSON(http.StatusInternalServerError, utils.NewErrorResponse("Internal Server Errors", "Error retrieving chat list"))
		return
	}

	ctx.JSON(http.StatusOK, utils.NewGetResponse(len(chatListData), chatListData))
}
