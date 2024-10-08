package controller

import (
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/kwa0x2/swiftchat-backend/models"
	"github.com/kwa0x2/swiftchat-backend/socket/gateway"
	"github.com/kwa0x2/swiftchat-backend/types"
	"gorm.io/gorm"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kwa0x2/swiftchat-backend/service"
	"github.com/kwa0x2/swiftchat-backend/utils"
)

type IRequestController interface {
	GetRequests(ctx *gin.Context)
	Patch(ctx *gin.Context)
	SendFriend(ctx *gin.Context)
}

type requestController struct {
	RequestService service.IRequestService
	FriendService  service.IFriendService
	UserService    service.IUserService
	SocketGateway  gateway.ISocketGateway
	ResendService  service.IResendService
}

func NewRequestController(RequestService service.IRequestService, friendService service.IFriendService,
	userService service.IUserService, socketGateway gateway.ISocketGateway, resendService service.IResendService) IRequestController {
	return &requestController{
		RequestService: RequestService,
		FriendService:  friendService,
		UserService:    userService,
		SocketGateway:  socketGateway,
		ResendService:  resendService,
	}
}

// region "GetRequests" handles the request to retrieve incoming friend requests.
func (ctrl *requestController) GetRequests(ctx *gin.Context) {
	// Get user session information.
	userSessionInfo, err := utils.GetUserSessionInfo(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.NewErrorResponse("Session Error", err.Error()))
		return
	}

	// Retrieve requests for the user.
	data, GetReqErr := ctrl.RequestService.GetRequests(userSessionInfo.Email)
	if GetReqErr != nil {
		ctx.JSON(http.StatusInternalServerError, utils.NewErrorResponse("Internal Server Error", "Error retrieving incoming requests"))
		return
	}

	// Prepare response data from the retrieved incoming requests.
	var responseData []map[string]interface{}
	for _, item := range data {
		responseItem := map[string]interface{}{
			"sender_mail": item.SenderMail,     // The email of the sender.
			"user_name":   item.User.UserName,  // The name of the user who sent the request.
			"user_photo":  item.User.UserPhoto, // The photo of the user who sent the request.
		}
		// Add each response item to the response data array.
		responseData = append(responseData, responseItem)
	}

	ctx.JSON(http.StatusOK, utils.NewGetResponse(len(responseData), responseData))
}

// endregion

// region "Patch" handles the request to update a friend request's status.
func (ctrl *requestController) Patch(ctx *gin.Context) {
	var requestBody ActionBody

	// Bind JSON request body to ActionBody struct.
	if err := ctx.BindJSON(&requestBody); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.NewErrorResponse("JSON Bind Error", err.Error()))
		return
	}

	// Get user session information.
	userSessionInfo, err := utils.GetUserSessionInfo(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.NewErrorResponse("Session Error", err.Error()))
		return
	}

	// Update the friendship request status.
	data, UpdateErr := ctrl.RequestService.UpdateFriendshipRequest(userSessionInfo.Email, requestBody.Email, requestBody.Status)
	if UpdateErr != nil {
		ctx.JSON(http.StatusInternalServerError, utils.NewErrorResponse("Internal Server Error", "Error updating friendship request"))
		return
	}

	// Emit notification about the updated friendship request to the socket gateway.
	ctrl.SocketGateway.EmitToNotificationRoom("update_friendship_request", requestBody.Email, data)
	ctx.JSON(http.StatusOK, utils.NewSuccessResponse("Success", "Friendship request updated successfully"))
}

// endregion

// region "SendFriend" handles the request to send a friend request.
func (ctrl *requestController) SendFriend(ctx *gin.Context) {
	var actionBody ActionBody

	// Bind JSON request body to ActionBody struct.
	if err := ctx.BindJSON(&actionBody); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.NewErrorResponse("JSON Bind Error", err.Error()))
		return
	}

	// Get user session information.
	userSessionInfo, err := utils.GetUserSessionInfo(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.NewErrorResponse("Session Error", err.Error()))
		return
	}

	requestObj := models.Request{
		SenderMail:   userSessionInfo.Email, // Sender's email from session.
		ReceiverMail: actionBody.Email,      // Receiver's email from request body.
	}

	var pgErr *pgconn.PgError
	existingFriend, GetSpecificFriendErr := ctrl.FriendService.GetSpecificFriend(requestObj.SenderMail, requestObj.ReceiverMail)
	if GetSpecificFriendErr != nil && !errors.Is(GetSpecificFriendErr, gorm.ErrRecordNotFound) {
		ctx.JSON(http.StatusInternalServerError, utils.NewErrorResponse("Internal Server Error", "Unable to retrieve friend status"))
		return
	}

	// Check for existing friendship status.
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

	// Check if the receiver's email exists.
	if isEmailExists := ctrl.UserService.IsEmailExists(requestObj.ReceiverMail); !isEmailExists {
		// If not, create a new friend request.
		if createErr := ctrl.RequestService.Create(nil, &requestObj); err != nil {
			if errors.As(createErr, &pgErr) && pgErr.Code == "23505" {
				ctx.JSON(http.StatusConflict, utils.NewErrorResponse("Already Sent", "Duplicate friend request"))
				return
			}
			ctx.JSON(http.StatusInternalServerError, utils.NewErrorResponse("Internal Server Error", "Failed to insert friend request"))
			return
		}

		// Send email notification about the friend request.
		_, SendEmailErr := ctrl.ResendService.SendEmail(requestObj.ReceiverMail, "You have received a new friend request from the SwiftChat app!", "friend_request")
		if SendEmailErr != nil {
			ctx.JSON(http.StatusInternalServerError, utils.NewErrorResponse("Internal Server Error", "Failed to send email"))
			return
		}

		ctx.JSON(http.StatusOK, utils.NewSuccessResponse("Email Sent", "Friend request email sent"))
		return
	}

	// If the email exists, insert and return user information.
	data, dataErr := ctrl.RequestService.InsertAndReturnUser(&requestObj)
	if dataErr != nil {
		if errors.As(dataErr, &pgErr) && pgErr.Code == "23505" {
			ctx.JSON(http.StatusConflict, utils.NewErrorResponse("Already Sent", "Duplicate friend request"))
			return
		}
		ctx.JSON(http.StatusInternalServerError, utils.NewErrorResponse("Internal Server Error", "Failed to insert friend request"))
		return
	}

	// Emit a notification about the friend request to the socket gateway.
	ctrl.SocketGateway.EmitToNotificationRoom("friend_request", requestObj.ReceiverMail, data)
	ctx.JSON(http.StatusOK, utils.NewSuccessResponse("Friend Sent", "Friend request successfully sent"))
}

// endregion
