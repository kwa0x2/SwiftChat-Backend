package controller

import (
	"github.com/kwa0x2/realtime-chat-backend/models"
	"github.com/kwa0x2/realtime-chat-backend/socket/adapter"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kwa0x2/realtime-chat-backend/service"
	"github.com/kwa0x2/realtime-chat-backend/utils"
)

type IFriendController interface {
	GetFriends(ctx *gin.Context)
	GetBlockedUsers(ctx *gin.Context)
	Block(ctx *gin.Context)
	Delete(ctx *gin.Context)
}

type friendController struct {
	friendService  *service.FriendService
	userService    *service.UserService
	requestService *service.RequestService
	socketAdapter  *adapter.SocketAdapter
}

func NewFriendController(friendService *service.FriendService, userService *service.UserService, requestService *service.RequestService, socketAdapter *adapter.SocketAdapter) *friendController {
	return &friendController{
		friendService:  friendService,
		userService:    userService,
		requestService: requestService,
		socketAdapter:  socketAdapter,
	}
}

func (ctrl *friendController) GetFriends(ctx *gin.Context) {
	userSessionInfo, err := utils.GetUserSessionInfo(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.NewErrorResponse("Session Error", err.Error()))
		return
	}

	data, err := ctrl.friendService.GetFriends(userSessionInfo.Email, false)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.NewErrorResponse("Internal Server Error", "Error retrieving friends"))
		return
	}

	var responseData []map[string]interface{}
	for _, item := range data {

		responseItem := map[string]interface{}{
			"friend_mail": item.UserMail,
			"user_name":   item.User.UserName,
			"user_photo":  item.User.UserPhoto,
		}
		responseData = append(responseData, responseItem)

	}
	ctx.JSON(http.StatusOK, utils.NewGetResponse(len(responseData), responseData))
}

func (ctrl *friendController) GetBlockedUsers(ctx *gin.Context) {
	userSessionInfo, err := utils.GetUserSessionInfo(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.NewErrorResponse("Session Error", err.Error()))
		return
	}

	data, err := ctrl.friendService.GetBlocked(userSessionInfo.Email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.NewErrorResponse("Internal Server Error", "Error retrieving blocked users"))
		return
	}

	var responseData []map[string]interface{}
	for _, item := range data {
		responseItem := map[string]interface{}{
			"blocked_mail": item.UserMail,
			"user_name":    item.User.UserName,
		}
		responseData = append(responseData, responseItem)

	}
	ctx.JSON(http.StatusOK, utils.NewGetResponse(len(responseData), responseData))
}

func (ctrl *friendController) Block(ctx *gin.Context) {
	var actionBody ActionBody

	if err := ctx.BindJSON(&actionBody); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.NewErrorResponse("JSON Bind Error", err.Error()))
		return
	}

	userSessionInfo, err := utils.GetUserSessionInfo(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.NewErrorResponse("Session Error", err.Error()))
		return
	}

	friend := models.Friend{
		UserMail:  actionBody.Email,      // User being blocked
		UserMail2: userSessionInfo.Email, // User blocking
	}

	friendStatus, err := ctrl.friendService.Block(&friend)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.NewErrorResponse("Internal Server Error", "Error blocking the user"))
		return
	}

	data := map[string]interface{}{
		"friend_mail":   friend.UserMail2,
		"friend_status": friendStatus,
	}

	ctrl.socketAdapter.EmitToNotificationRoom("blocked_friend", friend.UserMail, data)
	ctx.JSON(http.StatusOK, utils.NewSuccessResponse("Success", "User has been successfully blocked"))
}

func (ctrl *friendController) Delete(ctx *gin.Context) {
	var actionBody ActionBody

	if err := ctx.BindJSON(&actionBody); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.NewErrorResponse("JSON Bind Error", err.Error()))
		return
	}

	userSessionInfo, err := utils.GetUserSessionInfo(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.NewErrorResponse("Session Error", err.Error()))
		return
	}

	friendDeleteObj := models.Friend{
		UserMail:  actionBody.Email,
		UserMail2: userSessionInfo.Email,
	}

	if err := ctrl.friendService.Delete(&friendDeleteObj); err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.NewErrorResponse("Delete Friend Error", "Error delete the user"))
		return
	}

	notifyData := map[string]interface{}{
		"user_email": userSessionInfo.Email,
	}

	ctrl.socketAdapter.EmitToNotificationRoom("deleted_friend", actionBody.Email, notifyData)
	ctx.JSON(http.StatusOK, utils.NewSuccessResponse("Success", "User has been successfully deleted"))
}
