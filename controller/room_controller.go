package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/kwa0x2/swiftchat-backend/service"
	"github.com/kwa0x2/swiftchat-backend/types"
	"github.com/kwa0x2/swiftchat-backend/utils"
	"net/http"
)

type IRoomController interface {
	GetOrCreateRoom(ctx *gin.Context)
	GetChatList(ctx *gin.Context)
}

type roomController struct {
	RoomService     service.IRoomService
	UserRoomService service.IUserRoomService
	UserService     service.IUserService
	FriendService   service.IFriendService
}

func NewRoomController(roomService service.IRoomService, userRoomService service.IUserRoomService, userService service.IUserService, friendService service.IFriendService) IRoomController {
	return &roomController{
		RoomService:     roomService,
		UserRoomService: userRoomService,
		UserService:     userService,
		FriendService:   friendService,
	}
}

// region ActionBody represents the structure of the request body for certain actions.
type ActionBody struct {
	Email  string              `json:"email"`
	Status types.RequestStatus `json:"status"`
}

// endregion

// region "GetOrCreateRoom" handles the request to retrieve or create a chat room.
func (ctrl *roomController) GetOrCreateRoom(ctx *gin.Context) {
	var actionBody ActionBody

	// Bind JSON request body to ActionBody struct.
	if err := ctx.BindJSON(&actionBody); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.NewErrorResponse("JSON Bind Error", err.Error()))
		return
	}

	// Get user session information.
	userSessionInfo, userSessionErr := utils.GetUserSessionInfo(ctx)
	if userSessionErr != nil {
		ctx.JSON(http.StatusBadRequest, utils.NewErrorResponse("Session Error", userSessionErr.Error()))
		return
	}

	// Fetch the user by email to find the other participant in the chat.
	user, userErr := ctrl.UserService.GetByEmail(actionBody.Email)
	if userErr != nil {
		ctx.JSON(http.StatusInternalServerError, utils.NewErrorResponse("Internal Server Errors", "Error retrieving user by email"))
		return
	}

	// Check if a room already exists between the current user and the fetched user.
	room, roomErr := ctrl.UserRoomService.GetRoom(userSessionInfo.ID, user.UserID)
	if roomErr != nil {
		ctx.JSON(http.StatusInternalServerError, utils.NewErrorResponse("Internal Server Error", "Error retrieving private room"))
		return
	}

	var roomId string

	// If no room exists, create a new room and add users to it.
	if room == "" {
		newRoomId, newRoomErr := ctrl.RoomService.CreateAndAddUsers(userSessionInfo.ID, user.UserID)
		if newRoomErr != nil {
			// If there's an error creating the room, return an internal server error response.
			ctx.JSON(http.StatusInternalServerError, utils.NewErrorResponse("Internal Server Error", "Error creating and adding users to room"))
			return
		}
		roomId = newRoomId // Set the new room ID.
	} else {
		roomId = room // Use the existing room ID.
	}

	ctx.JSON(http.StatusOK, gin.H{
		"room_id": roomId,
	})
}

// endregion

// region "GetChatList" handles the request to retrieve the user's chat list.
func (ctrl *roomController) GetChatList(ctx *gin.Context) {
	// Get user session information.
	userSessionInfo, err := utils.GetUserSessionInfo(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.NewErrorResponse("Session Error", err.Error()))
		return
	}

	// Fetch the user's chat list using their session information.
	chatListData, chatListErr := ctrl.RoomService.GetChatList(userSessionInfo.ID, userSessionInfo.Email)
	if chatListErr != nil {
		ctx.JSON(http.StatusInternalServerError, utils.NewErrorResponse("Internal Server Errors", "Error retrieving chat list"))
		return
	}

	ctx.JSON(http.StatusOK, utils.NewGetResponse(len(chatListData), chatListData))
}

// endregion
