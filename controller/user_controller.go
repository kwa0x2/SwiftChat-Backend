package controller

import (
	"fmt"
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
	userService   *service.UserService
	friendService *service.FriendService
	s3Service     *service.S3Service
	socketAdapter *adapter.SocketAdapter
}

func NewUserController(userService *service.UserService, friendService *service.FriendService, s3Service *service.S3Service, socketAdapter *adapter.SocketAdapter) IUserController {
	return &userController{
		userService:   userService,
		friendService: friendService,
		s3Service:     s3Service,
		socketAdapter: socketAdapter,
	}
}

type UsernameUpdateBody struct {
	UserName string `json:"user_name"`
}

func (ctrl *userController) UpdateUsername(ctx *gin.Context) {
	var requestBody UsernameUpdateBody
	session := sessions.Default(ctx)

	if err := ctx.BindJSON(&requestBody); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.NewErrorResponse("JSON Bind Error", err.Error()))
		return
	}

	userSessionInfo, sessionErr := utils.GetUserSessionInfo(ctx)
	if sessionErr != nil {
		ctx.JSON(http.StatusBadRequest, utils.NewErrorResponse("Session Error", sessionErr.Error()))
		return
	}

	if err := ctrl.userService.UpdateUsernameByMail(requestBody.UserName, userSessionInfo.Email); err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.NewErrorResponse("Internal Server Error", "Error updating username by email"))
		return
	}

	session.Set("name", requestBody.UserName)
	session.Save()

	emitData := map[string]interface{}{
		"updated_username": requestBody.UserName,
		"user_email":       userSessionInfo.Email,
	}

	fmt.Println(emitData)

	if err := ctrl.socketAdapter.EmitToFriends("update_username", userSessionInfo.Email, emitData); err == nil {
		ctx.JSON(http.StatusInternalServerError, utils.NewErrorResponse("Internal Server Error", "Failed to emit update username notification to friends"))
		return
	}

	ctx.JSON(http.StatusOK, utils.NewSuccessResponse("OK", "Username successfully updated"))
}

func (ctrl *userController) UploadProfilePhoto(ctx *gin.Context) {
	file, header, err := ctx.Request.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.NewErrorResponse("Form File Error", err.Error()))
		return
	}
	defer file.Close()

	session := sessions.Default(ctx)

	userSessionInfo, sessionErr := utils.GetUserSessionInfo(ctx)
	if sessionErr != nil {
		ctx.JSON(http.StatusBadRequest, utils.NewErrorResponse("Session Error", sessionErr.Error()))
		return
	}

	fileURL, err := ctrl.s3Service.UploadFile(file, header)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.NewErrorResponse("Internal Server Error", "Error uploading file to S3 bucket"))
		return
	}

	if err := ctrl.userService.UpdateUserPhotoByMail(fileURL, userSessionInfo.Email); err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.NewErrorResponse("Internal Server Error", "Error updating user photo"))
		return
	}

	if session.Get("email") != nil {
		session.Set("photo", fileURL)
		session.Save()
	}

	emitData := map[string]interface{}{
		"updated_user_photo": fileURL,
		"user_email":         userSessionInfo.Email,
	}

	if err := ctrl.socketAdapter.EmitToFriends("update_user_photo", userSessionInfo.Email, emitData); err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.NewErrorResponse("Internal Server Error", "Failed to emit update username notification to friends"))
		return
	}

	ctx.JSON(http.StatusOK, fileURL)
}
