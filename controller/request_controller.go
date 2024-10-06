package controller

import (
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/kwa0x2/realtime-chat-backend/models"
	"github.com/kwa0x2/realtime-chat-backend/socket/adapter"
	"github.com/kwa0x2/realtime-chat-backend/types"
	"gorm.io/gorm"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kwa0x2/realtime-chat-backend/service"
	"github.com/kwa0x2/realtime-chat-backend/utils"
)

type IRequestController interface {
	GetComingRequests(ctx *gin.Context)
	PatchUpdateRequest(ctx *gin.Context)
	SendFriend(ctx *gin.Context)
}

type requestController struct {
	requestService *service.RequestService
	friendService  *service.FriendService
	userService    *service.UserService
	socketAdapter  *adapter.SocketAdapter
	resendService  *service.ResendService
}

func NewRequestController(requestService *service.RequestService, friendService *service.FriendService,
	userService *service.UserService, socketAdapter *adapter.SocketAdapter, resendService *service.ResendService) IRequestController {
	return &requestController{
		requestService: requestService,
		friendService:  friendService,
		userService:    userService,
		socketAdapter:  socketAdapter,
		resendService:  resendService,
	}
}

func (ctrl *requestController) GetComingRequests(ctx *gin.Context) {
	userSessionInfo, err := utils.GetUserSessionInfo(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.NewErrorResponse("Session Error", err.Error()))
		return
	}

	data, err := ctrl.requestService.GetComingRequests(userSessionInfo.Email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.NewErrorResponse("Internal Server Error", "Error retrieving incoming requests"))
		return
	}

	var responseData []map[string]interface{}
	for _, item := range data {
		responseItem := map[string]interface{}{
			"sender_mail": item.SenderMail,
			"user_name":   item.User.UserName,
			"user_photo":  item.User.UserPhoto,
		}
		responseData = append(responseData, responseItem)
	}

	ctx.JSON(http.StatusOK, utils.NewGetResponse(len(responseData), responseData))
}

func (ctrl *requestController) PatchUpdateRequest(ctx *gin.Context) {
	var requestBody ActionBody

	if err := ctx.BindJSON(&requestBody); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.NewErrorResponse("JSON Bind Error", err.Error()))
		return
	}

	userSessionInfo, err := utils.GetUserSessionInfo(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.NewErrorResponse("Session Error", err.Error()))
		return
	}

	requestObj := models.Request{
		SenderMail:    requestBody.Email,
		ReceiverMail:  userSessionInfo.Email,
		RequestStatus: requestBody.Status,
	}

	data, err := ctrl.requestService.UpdateFriendshipRequest(&requestObj)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.NewErrorResponse("Internal Server Error", "Error updating friendship request"))
		return
	}

	ctrl.socketAdapter.EmitToNotificationRoom("update_friendship_request", requestBody.Email, data)
	ctx.JSON(http.StatusOK, utils.NewSuccessResponse("Success", "Friendship request updated successfully"))
}

func (ctrl *requestController) SendFriend(ctx *gin.Context) {
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

	requestObj := models.Request{
		SenderMail:   userSessionInfo.Email,
		ReceiverMail: actionBody.Email,
	}

	var pgErr *pgconn.PgError
	existingFriend, err := ctrl.friendService.GetFriend(requestObj.SenderMail, requestObj.ReceiverMail)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		ctx.JSON(http.StatusInternalServerError, utils.NewErrorResponse("Internal Server Error", "Unable to retrieve friend status"))
		return
	}

	if existingFriend != nil {
		fmt.Println("status", existingFriend.FriendStatus)
		if existingFriend.FriendStatus == types.Friend {
			ctx.JSON(http.StatusConflict, utils.NewErrorResponse("Already Friend", "Users are already friends"))
			return
		} else if existingFriend.FriendStatus == types.BlockFirstSecond || existingFriend.FriendStatus == types.BlockSecondFirst {
			ctx.JSON(http.StatusConflict, utils.NewErrorResponse("Blocked User", "Users are blocked"))
			return
		}
	}

	if isEmailExists := ctrl.userService.IsEmailExists(requestObj.ReceiverMail); !isEmailExists {
		if err := ctrl.requestService.Insert(nil, &requestObj); err != nil {
			if errors.As(err, &pgErr) && pgErr.Code == "23505" {
				ctx.JSON(http.StatusConflict, utils.NewErrorResponse("Already Sent", "Duplicate friend request"))
				return
			}
			ctx.JSON(http.StatusInternalServerError, utils.NewErrorResponse("Internal Server Error", "Failed to insert friend request"))
			return
		}

		_, err := ctrl.resendService.SendMail(requestObj.ReceiverMail, "You have received a new friend request from the SwiftChat app!", "friend_request")
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, utils.NewErrorResponse("Internal Server Error", "Failed to send email"))
			return
		}

		ctx.JSON(http.StatusOK, utils.NewSuccessResponse("Email Sent", "Friend request email sent"))
		return
	}

	data, err := ctrl.requestService.InsertAndReturnUser(&requestObj)
	if err != nil {
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			ctx.JSON(http.StatusConflict, utils.NewErrorResponse("Already Sent", "Duplicate friend request"))
			return
		}
		ctx.JSON(http.StatusInternalServerError, utils.NewErrorResponse("Internal Server Error", "Failed to insert friend request"))
		return
	}

	ctrl.socketAdapter.EmitToNotificationRoom("friend_request", requestObj.ReceiverMail, data)
	ctx.JSON(http.StatusOK, utils.NewSuccessResponse("Friend Sent", "Friend request successfully sent"))
}
