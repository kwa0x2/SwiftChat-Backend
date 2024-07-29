package controller

import (
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

	existRoomId, err := ctrl.UserRoomService.GetPrivateRoom(userId.(string), user.UserID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.NewErrorResponse(http.StatusInternalServerError, "Internal Server Error", err.Error()))
		return
	}

	var roomId string

	if existRoomId == "" {
		newRoomId, err := ctrl.RoomService.CreateAndAddUsers(userId.(string), user.UserID, "private")
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, utils.NewErrorResponse(http.StatusInternalServerError, "Internal Server Error", err.Error()))
			return
		}
		roomId = newRoomId
	} else {
		roomId = existRoomId
	}

	ctx.JSON(http.StatusOK, gin.H{
		"room_id": roomId,
	})
}
