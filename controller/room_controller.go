package controller

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/kwa0x2/realtime-chat-backend/service"
	"github.com/kwa0x2/realtime-chat-backend/utils"
	"net/http"
)

type RoomController struct {
	RoomService     *service.RoomService
	UserRoomService *service.UserRoomService
	UserService     *service.UserService
	FriendService   *service.FriendService
}

func (ctrl *RoomController) GetOrCreatePrivateRoom(ctx *gin.Context) {
	var actionBody ActionBody
	session := sessions.Default(ctx)

	if err := ctx.BindJSON(&actionBody); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.NewErrorResponse(http.StatusBadRequest, "Bad Request", err.Error()))
		return
	}

	userId := session.Get("id")
	if userId == nil {
		ctx.JSON(http.StatusBadRequest, utils.NewErrorResponse(http.StatusBadRequest, "Bad Request", "UserId not found"))
		return
	}

	user, err := ctrl.UserService.GetByEmail(actionBody.Mail)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.NewErrorResponse(http.StatusInternalServerError, "Internal Server Errors", err.Error()))
		return
	}

	room, roomErr := ctrl.UserRoomService.GetPrivateRoom(userId.(string), user.UserID)
	if roomErr != nil {
		ctx.JSON(http.StatusInternalServerError, utils.NewErrorResponse(http.StatusInternalServerError, "Internal Server Error", err.Error()))
		return
	}

	//data, err := ctrl.FriendService.GetBlocked(userMail.(string))
	//if err != nil {
	//	ctx.JSON(http.StatusInternalServerError, utils.NewErrorResponse(http.StatusInternalServerError, "Internal Server Error", err.Error()))
	//}

	var roomId string

	if room == "" {
		newRoomId, newRoomErr := ctrl.RoomService.CreateAndAddUsers(userId.(string), user.UserID, "private")
		if newRoomErr != nil {
			ctx.JSON(http.StatusInternalServerError, utils.NewErrorResponse(http.StatusInternalServerError, "Internal Server Error", err.Error()))
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

func (ctrl *RoomController) GetChatList(ctx *gin.Context) {
	session := sessions.Default(ctx)

	userId := session.Get("id")
	if userId == nil {
		ctx.JSON(http.StatusBadRequest, utils.NewErrorResponse(http.StatusBadRequest, "Bad Request", "UserId not found"))
		return
	}
	fmt.Printf("userid %s", userId.(string))
	data, err := ctrl.RoomService.GetChatList(userId.(string))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.NewErrorResponse(http.StatusInternalServerError, "Internal Server Errors", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, data)
}
