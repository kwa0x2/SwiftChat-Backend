package controller

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/kwa0x2/realtime-chat-backend/service"
	"github.com/kwa0x2/realtime-chat-backend/utils"
	"net/http"
)

type UserController struct {
	UserService *service.UserService
	S3Service   *service.S3Service
}

type UsernameUpdateBody struct {
	UserName string `json:"user_name"`
}

func (ctrl *UserController) UpdateUsername(ctx *gin.Context) {
	var requestBody UsernameUpdateBody
	session := sessions.Default(ctx)
	fmt.Println("ilk", session.Get("name"))

	if err := ctx.BindJSON(&requestBody); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.NewErrorResponse(http.StatusBadRequest, "Bad Request", err.Error()))
		return
	}

	userMail := session.Get("mail")
	if userMail == nil {
		ctx.JSON(http.StatusBadRequest, utils.NewErrorResponse(http.StatusBadRequest, "Bad Request", "UserId not found"))
		return
	}

	if err := ctrl.UserService.UpdateUsernameByMail(requestBody.UserName, userMail.(string)); err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.NewErrorResponse(http.StatusInternalServerError, "Bad Request", err.Error()))
		return
	}

	session.Set("name", requestBody.UserName)
	session.Save()

	fmt.Println("son", session.Get("name"))
	ctx.JSON(http.StatusOK, utils.NewSuccessResponse(http.StatusOK, "OK", "success"))
}

func (ctrl *UserController) UploadProfilePicture(ctx *gin.Context) {
	file, header, err := ctx.Request.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.NewErrorResponse(http.StatusBadRequest, "Bad Request", err.Error()))
		return
	}
	defer file.Close()
	session := sessions.Default(ctx)

	userMail := session.Get("mail")
	if userMail == nil {
		ctx.JSON(http.StatusBadRequest, utils.NewErrorResponse(http.StatusBadRequest, "Bad Request", "UserId not found"))
		return
	}

	fileURL, err := ctrl.S3Service.UploadFile(file, header)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.NewErrorResponse(http.StatusInternalServerError, "Internal Server Error", err.Error()))
		return
	}

	if err := ctrl.UserService.UpdateUserPhotoByMail(fileURL, userMail.(string)); err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.NewErrorResponse(http.StatusInternalServerError, "Bad Request", err.Error()))
		return
	}

	session.Set("photo", fileURL)
	session.Save()

	ctx.JSON(http.StatusOK, utils.NewSuccessResponse(http.StatusOK, "OK", "File uploaded successfully"))
}
