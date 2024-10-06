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

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/kwa0x2/realtime-chat-backend/service"
	"github.com/kwa0x2/realtime-chat-backend/utils"
)

type RequestController struct {
	RequestService *service.RequestService
	FriendService  *service.FriendService
	UserService    *service.UserService
	SocketAdapter  *adapter.SocketAdapter
	ResendService  *service.ResendService
}

func (ctrl *RequestController) GetComingRequests(ctx *gin.Context) {
	session := sessions.Default(ctx)

	userMail := session.Get("mail")
	if userMail == nil {
		ctx.JSON(http.StatusBadRequest, utils.NewErrorResponse("Session Error", "UserMail not found"))
		return
	}

	data, err := ctrl.RequestService.GetComingRequests(userMail.(string))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.NewErrorResponse("Internal Server Error", "gelen mesajlar alinirken hata"))
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

	ctx.JSON(http.StatusOK, responseData)
}

func (ctrl *RequestController) PatchUpdateRequest(ctx *gin.Context) {
	var requestBody ActionBody
	session := sessions.Default(ctx)

	if err := ctx.BindJSON(&requestBody); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.NewErrorResponse("JSON Bind Error", err.Error()))
		return
	}

	userMail := session.Get("mail")
	if userMail == nil {
		ctx.JSON(http.StatusBadRequest, utils.NewErrorResponse("Session Error", "UserMail not found"))
		return
	}

	requestObj := models.Request{
		SenderMail:    requestBody.Mail,
		ReceiverMail:  userMail.(string),
		RequestStatus: requestBody.Status,
	}

	data, err := ctrl.RequestService.UpdateFriendshipRequest(&requestObj)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.NewErrorResponse("Internal Server Error", "friendship req update edilirken hata"))
		return
	}

	ctrl.SocketAdapter.EmitToNotificationRoom("update_friendship_request", requestBody.Mail, data)
	ctx.JSON(http.StatusOK, utils.NewSuccessResponse("Success", "success"))
}

func (ctrl *RequestController) SendFriend(ctx *gin.Context) {
	var actionBody ActionBody
	session := sessions.Default(ctx)

	if err := ctx.BindJSON(&actionBody); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.NewErrorResponse("JSON Bind Error", err.Error()))
		return
	}

	userMail := session.Get("mail")
	if userMail == nil {
		ctx.JSON(http.StatusBadRequest, utils.NewErrorResponse("Session Error", "UserMail not found"))
		return
	}

	requestObj := models.Request{
		SenderMail:   userMail.(string),
		ReceiverMail: actionBody.Mail,
	}

	var pgErr *pgconn.PgError
	existingFriend, err := ctrl.FriendService.GetFriend(requestObj.SenderMail, requestObj.ReceiverMail)
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

	if isEmailExists := ctrl.UserService.IsEmailExists(requestObj.ReceiverMail); !isEmailExists {
		if err := ctrl.RequestService.Insert(nil, &requestObj); err != nil {
			if errors.As(err, &pgErr) && pgErr.Code == "23505" {
				ctx.JSON(http.StatusConflict, utils.NewErrorResponse("Already Sent", "Duplicate friend request"))
				return
			}
			ctx.JSON(http.StatusInternalServerError, utils.NewErrorResponse("Internal Server Error", "Failed to insert friend request"))
			return
		}

		_, err := ctrl.ResendService.SendMail(requestObj.ReceiverMail, "You have received a new friend request from the SwiftChat app!", "friend_request")
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, utils.NewErrorResponse("Internal Server Error", "Failed to send email"))
			return
		}

		ctx.JSON(http.StatusOK, utils.NewSuccessResponse("Email Sent", "Friend request email sent"))
		return
	}

	data, err := ctrl.RequestService.InsertAndReturnUser(&requestObj)
	if err != nil {
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			ctx.JSON(http.StatusConflict, utils.NewErrorResponse("Already Sent", "Duplicate friend request"))
			return
		}
		ctx.JSON(http.StatusInternalServerError, utils.NewErrorResponse("Internal Server Error", "Failed to insert friend request"))
		return
	}

	ctrl.SocketAdapter.EmitToNotificationRoom("friend_request", requestObj.ReceiverMail, data)
	ctx.JSON(http.StatusOK, utils.NewSuccessResponse("Friend Sent", "Friend request successfully sent"))
}
