package controller

import (
	"github.com/kwa0x2/swiftchat-backend/socket/gateway"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kwa0x2/swiftchat-backend/service"
	"github.com/kwa0x2/swiftchat-backend/utils"
)

type IFriendController interface {
	GetFriends(ctx *gin.Context)
	GetBlockedUsers(ctx *gin.Context)
	Block(ctx *gin.Context)
	Delete(ctx *gin.Context)
}

type friendController struct {
	FriendService service.IFriendService
	SocketGateway gateway.ISocketGateway
}

func NewFriendController(FriendService service.IFriendService, socketGateway gateway.ISocketGateway) IFriendController {
	return &friendController{
		FriendService: FriendService,
		SocketGateway: socketGateway,
	}
}

// region "GetFriends" retrieves the list of friends for the authenticated user.
func (ctrl *friendController) GetFriends(ctx *gin.Context) {
	// Get user session information
	userSessionInfo, err := utils.GetUserSessionInfo(ctx)
	if err != nil {
		utils.HandleErrorWithSentry(ctx, err, nil)
		ctx.JSON(http.StatusBadRequest, utils.NewErrorResponse("Session Error", err.Error()))
		return
	}

	// Retrieve the user's friends, with the second parameter 'false' indicating that only friends (not unfriended users) should be fetched.
	friends, GetFriendsErr := ctrl.FriendService.GetFriends(userSessionInfo.Email, false)
	if GetFriendsErr != nil {
		utils.HandleErrorWithSentry(ctx, GetFriendsErr, map[string]interface{}{"user_email": userSessionInfo.Email})
		ctx.JSON(http.StatusInternalServerError, utils.NewErrorResponse("Internal Server Error", "Error retrieving friends"))
		return
	}

	// Prepare the response data by mapping the friend information
	var responseData []map[string]interface{}
	for _, item := range friends {
		responseItem := map[string]interface{}{
			"friend_email": item.UserEmail,      // Friend's email
			"user_name":    item.User.UserName,  // Friend's username
			"user_photo":   item.User.UserPhoto, // Friend's profile photo
		}
		responseData = append(responseData, responseItem) // Append the formatted item
	}

	ctx.JSON(http.StatusOK, utils.NewGetResponse(len(responseData), responseData))
}

// endregion

// region "GetBlockedUsers" retrieves the list of users blocked by the authenticated user.
func (ctrl *friendController) GetBlockedUsers(ctx *gin.Context) {
	// Get user session information
	userSessionInfo, err := utils.GetUserSessionInfo(ctx)
	if err != nil {
		utils.HandleErrorWithSentry(ctx, err, nil)
		ctx.JSON(http.StatusBadRequest, utils.NewErrorResponse("Session Error", err.Error()))
		return
	}

	// Fetch blocked users from the FriendService
	blockedUsers, GetBlockedUsersErr := ctrl.FriendService.GetBlockedUsers(userSessionInfo.Email)
	if GetBlockedUsersErr != nil {
		utils.HandleErrorWithSentry(ctx, GetBlockedUsersErr, map[string]interface{}{"user_email": userSessionInfo.Email})
		ctx.JSON(http.StatusInternalServerError, utils.NewErrorResponse("Internal Server Error", "Error retrieving blocked users"))
		return
	}

	// Prepare the response data by mapping the blocked user information
	var responseData []map[string]interface{}
	for _, item := range blockedUsers {
		responseItem := map[string]interface{}{
			"blocked_email": item.UserEmail,     // Blocked user's email
			"user_name":     item.User.UserName, // Blocked user's username
		}
		responseData = append(responseData, responseItem) // Append the formatted item
	}

	ctx.JSON(http.StatusOK, utils.NewGetResponse(len(responseData), responseData))
}

// endregion

// region "Block" blocks a user and sends a notification using Socket.IO.
func (ctrl *friendController) Block(ctx *gin.Context) {
	var actionBody ActionBody

	// Bind the JSON body to the actionBody struct
	if err := ctx.BindJSON(&actionBody); err != nil {
		utils.HandleErrorWithSentry(ctx, err, nil)
		ctx.JSON(http.StatusBadRequest, utils.NewErrorResponse("JSON Bind Error", err.Error()))
		return
	}

	// Get user session information
	userSessionInfo, userSessionErr := utils.GetUserSessionInfo(ctx)
	if userSessionErr != nil {
		utils.HandleErrorWithSentry(ctx, userSessionErr, nil)
		ctx.JSON(http.StatusBadRequest, utils.NewErrorResponse("Session Error", userSessionErr.Error()))
		return
	}

	// Block the user using the FriendService
	blockErr := ctrl.FriendService.Block(actionBody.Email, userSessionInfo.Email)
	if blockErr != nil {
		utils.HandleErrorWithSentry(ctx, blockErr, map[string]interface{}{"user_email": userSessionInfo.Email})
		ctx.JSON(http.StatusInternalServerError, utils.NewErrorResponse("Internal Server Error", "Error blocking the user"))
		return
	}

	// Prepare the notification data
	notifyData := map[string]interface{}{
		"friend_email":  userSessionInfo.Email, // Email of the user who is blocking
		"friend_status": "blocked",             // Status of the blocking action
	}

	// Emit a socket event to notify about the blocked friend
	ctrl.SocketGateway.EmitToNotificationRoom("blocked_friend", actionBody.Email, notifyData)
	ctx.JSON(http.StatusOK, utils.NewSuccessResponse("Success", "User has been successfully blocked"))
}

// endregion

// region "Delete" removes a friend from the user's friend list and sends a notification using Socket.IO.
func (ctrl *friendController) Delete(ctx *gin.Context) {
	var actionBody ActionBody

	// Bind the JSON body to the actionBody struct
	if err := ctx.BindJSON(&actionBody); err != nil {
		utils.HandleErrorWithSentry(ctx, err, nil)
		ctx.JSON(http.StatusBadRequest, utils.NewErrorResponse("JSON Bind Error", err.Error()))
		return
	}

	// Get user session information
	userSessionInfo, userSessionErr := utils.GetUserSessionInfo(ctx)
	if userSessionErr != nil {
		utils.HandleErrorWithSentry(ctx, userSessionErr, nil)
		ctx.JSON(http.StatusBadRequest, utils.NewErrorResponse("Session Error", userSessionErr.Error()))
		return
	}

	// Delete the friend using the FriendService
	if err := ctrl.FriendService.Delete(actionBody.Email, userSessionInfo.Email); err != nil {
		utils.HandleErrorWithSentry(ctx, err, map[string]interface{}{"user_email": userSessionInfo.Email})
		ctx.JSON(http.StatusInternalServerError, utils.NewErrorResponse("Delete Friend Error", "Error delete the user"))
		return
	}

	// Prepare the notification data for the deleted friend
	notifyData := map[string]interface{}{
		"user_email": userSessionInfo.Email, // Email of the user who deleted the friend
	}

	// Emit a socket event to notify about the deleted friend
	ctrl.SocketGateway.EmitToNotificationRoom("deleted_friend", actionBody.Email, notifyData)
	ctx.JSON(http.StatusOK, utils.NewSuccessResponse("Success", "User has been successfully deleted"))
}

// endregion
