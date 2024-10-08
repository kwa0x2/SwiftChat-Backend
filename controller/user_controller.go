package controller

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/kwa0x2/realtime-chat-backend/service"
	"github.com/kwa0x2/realtime-chat-backend/socket/adapter"
	"github.com/kwa0x2/realtime-chat-backend/utils"
	"net/http"
)

type IUserController interface {
	UpdateUsername(ctx *gin.Context)
	UploadProfilePhoto(ctx *gin.Context)
}

type userController struct {
	UserService   service.IUserService
	FriendService service.IFriendService
	S3Service     service.IS3Service
	SocketAdapter adapter.ISocketAdapter
}

func NewUserController(userService service.IUserService, friendService service.IFriendService, s3Service service.IS3Service, socketAdapter adapter.ISocketAdapter) IUserController {
	return &userController{
		UserService:   userService,
		FriendService: friendService,
		S3Service:     s3Service,
		SocketAdapter: socketAdapter,
	}
}

// region UsernameUpdateBody represents the structure of the request body for updating the username.
type UsernameUpdateBody struct {
	UserName string `json:"user_name"` // The new username for the user.
}

// endregion

// region "UpdateUsername" handles the request to update the user's username.
func (ctrl *userController) UpdateUsername(ctx *gin.Context) {
	var requestBody UsernameUpdateBody
	session := sessions.Default(ctx)

	// Bind JSON request body to UsernameUpdateBody struct.
	if err := ctx.BindJSON(&requestBody); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.NewErrorResponse("JSON Bind Error", err.Error()))
		return
	}

	// Get user session information.
	userSessionInfo, sessionErr := utils.GetUserSessionInfo(ctx)
	if sessionErr != nil {
		ctx.JSON(http.StatusBadRequest, utils.NewErrorResponse("Session Error", sessionErr.Error()))
		return
	}

	// Update the user's username in the database using their email.
	if err := ctrl.UserService.UpdateUserNameByMail(requestBody.UserName, userSessionInfo.Email); err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.NewErrorResponse("Internal Server Error", "Error updating username by email"))
		return
	}

	// Update the session with the new username.
	session.Set("name", requestBody.UserName)
	session.Save()

	// Prepare data to emit to friends regarding the username update.
	emitData := map[string]interface{}{
		"updated_username": requestBody.UserName,
		"user_email":       userSessionInfo.Email,
	}

	// Emit the username update notification to friends using the socket adapter.
	if err := ctrl.SocketAdapter.EmitToFriendsAndSentRequests("update_username", userSessionInfo.Email, emitData); err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.NewErrorResponse("Internal Server Error", "Failed to emit update username notification to friends"))
		return
	}

	ctx.JSON(http.StatusOK, utils.NewSuccessResponse("OK", "Username successfully updated"))
}

// endregion

// region "UploadProfilePhoto" handles the request to upload a user's profile photo.
func (ctrl *userController) UploadProfilePhoto(ctx *gin.Context) {
	// Retrieve the file from the form data.
	file, header, err := ctx.Request.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.NewErrorResponse("Form File Error", err.Error()))
		return
	}
	defer file.Close() // Ensure the file is closed after processing.

	session := sessions.Default(ctx)

	// Get user session information.
	userSessionInfo, userSessionErr := utils.GetUserSessionInfo(ctx)
	if userSessionErr != nil {
		ctx.JSON(http.StatusBadRequest, utils.NewErrorResponse("Session Error", userSessionErr.Error()))
		return
	}

	// Upload the file to the S3 bucket and retrieve the file URL.
	fileURL, UploadErr := ctrl.S3Service.UploadFile(file, header)
	if UploadErr != nil {
		ctx.JSON(http.StatusInternalServerError, utils.NewErrorResponse("Internal Server Error", "Error uploading file to S3 bucket"))
		return
	}

	// Update the user's photo in the database using their email.
	if UpdateErr := ctrl.UserService.UpdateUserPhotoByMail(fileURL, userSessionInfo.Email); UpdateErr != nil {
		ctx.JSON(http.StatusInternalServerError, utils.NewErrorResponse("Internal Server Error", "Error updating user photo"))
		return
	}

	// Update the session with the new photo URL.
	if session.Get("email") != nil {
		session.Set("photo", fileURL)
		session.Save()
	}

	// Prepare data to emit to friends regarding the photo update.
	emitData := map[string]interface{}{
		"updated_user_photo": fileURL,
		"user_email":         userSessionInfo.Email,
	}

	// Emit the photo update notification to friends using the socket adapter.
	if EmitErr := ctrl.SocketAdapter.EmitToFriendsAndSentRequests("update_user_photo", userSessionInfo.Email, emitData); EmitErr != nil {
		ctx.JSON(http.StatusInternalServerError, utils.NewErrorResponse("Internal Server Error", "Failed to emit update user photo notification to friends"))
		return
	}

	ctx.JSON(http.StatusOK, fileURL)
}

// endregion
