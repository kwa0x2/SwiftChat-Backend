package controller

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/kwa0x2/realtime-chat-backend/service"
	"github.com/kwa0x2/realtime-chat-backend/socket/adapter"
	"github.com/kwa0x2/realtime-chat-backend/utils"
	"net/http"
	"sync"
)

type UserController struct {
	UserService   *service.UserService
	FriendService *service.FriendService
	S3Service     *service.S3Service
	SocketAdapter *adapter.SocketAdapter
}

type UsernameUpdateBody struct {
	UserName string `json:"user_name"`
}

func (ctrl *UserController) UpdateUsername(ctx *gin.Context) {
	var requestBody UsernameUpdateBody
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

	if err := ctrl.UserService.UpdateUsernameByMail(requestBody.UserName, userMail.(string)); err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.NewErrorResponse("Internal Server Error", "update username by mail error"))
		return
	}

	session.Set("name", requestBody.UserName)
	session.Save()

	emitData := map[string]interface{}{
		"updated_username": requestBody.UserName,
		"user_email":       userMail.(string),
	}

	fmt.Println(emitData)

	friends, err := ctrl.FriendService.GetFriends(userMail.(string), true)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.NewErrorResponse("Internal Server Error", "get friends error"))
		return
	}

	var wg sync.WaitGroup
	for _, friend := range friends {
		wg.Add(1)
		go func(friendEmail string) {
			defer wg.Done()
			fmt.Println("updateusername")
			ctrl.SocketAdapter.EmitToNotificationRoom("update_username", friendEmail, emitData)
		}(friend.UserMail)
	}
	wg.Wait()

	ctx.JSON(http.StatusOK, utils.NewSuccessResponse("OK", "successfully changed"))
}

func (ctrl *UserController) UploadProfilePicture(ctx *gin.Context) {
	file, header, err := ctx.Request.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.NewErrorResponse("Form File Error", err.Error()))
		return
	}
	defer file.Close()

	userMail, exists := ctx.Get("mail")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, utils.NewErrorResponse("Unauthorized", "Authorization required"))
		return
	}

	fileURL, err := ctrl.S3Service.UploadFile(file, header)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.NewErrorResponse("Internal Server Error", "file upload to s3 bucket error"))
		return
	}

	if err := ctrl.UserService.UpdateUserPhotoByMail(fileURL, userMail.(string)); err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.NewErrorResponse("Internal Server Error", "update user photo by mail error"))
		return
	}

	session := sessions.Default(ctx)
	if session.Get("mail") != nil {
		session.Set("photo", fileURL)
		session.Save()
	}

	emitData := map[string]interface{}{
		"updated_user_photo": fileURL,
		"user_email":         userMail.(string),
	}

	friends, err := ctrl.FriendService.GetFriends(userMail.(string), true)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.NewErrorResponse("Internal Server Error", "get friends error"))
		return
	}

	var wg sync.WaitGroup
	for _, friend := range friends {
		wg.Add(1)
		go func(friendEmail string) {
			defer wg.Done()
			ctrl.SocketAdapter.EmitToNotificationRoom("update_user_photo", friendEmail, emitData)
		}(friend.UserMail)
	}
	wg.Wait()

	ctx.JSON(http.StatusOK, fileURL)
}
