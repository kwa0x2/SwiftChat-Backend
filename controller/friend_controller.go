package controller

import (
	"fmt"
	"github.com/kwa0x2/realtime-chat-backend/types"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/kwa0x2/realtime-chat-backend/models"
	"github.com/kwa0x2/realtime-chat-backend/service"
	"github.com/kwa0x2/realtime-chat-backend/utils"
)

type FriendController struct {
	FriendService *service.FriendService
	UserService   *service.UserService
}

func (ctrl *FriendController) GetFriends(ctx *gin.Context) {
	session := sessions.Default(ctx)

	userMail := session.Get("mail")
	if userMail == nil {
		ctx.JSON(http.StatusBadRequest, utils.NewErrorResponse(http.StatusBadRequest, "Bad Request", "UserId not found"))
		return
	}

	data, err := ctrl.FriendService.GetFriends(userMail.(string))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.NewErrorResponse(http.StatusInternalServerError, "Internal Server Error", err.Error()))
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
	ctx.JSON(http.StatusOK, responseData)
}

func (ctrl *FriendController) GetBlocked(ctx *gin.Context) {
	session := sessions.Default(ctx)

	userMail := session.Get("mail")
	if userMail == nil {
		ctx.JSON(http.StatusBadRequest, utils.NewErrorResponse(http.StatusBadRequest, "Bad Request", "UserId not found"))
		return
	}

	data, err := ctrl.FriendService.GetBlocked(userMail.(string))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.NewErrorResponse(http.StatusInternalServerError, "Internal Server Error", err.Error()))
	}

	var responseData []map[string]interface{}
	for _, item := range data {
		responseItem := map[string]interface{}{
			"blocked_mail": item.UserMail,
			"user_name":    item.User.UserName,
			"user_photo":   item.User.UserPhoto,
		}
		responseData = append(responseData, responseItem)

	}
	ctx.JSON(http.StatusOK, responseData)
}

type ActionBody struct {
	Mail   string              `json:"mail"`
	Status types.RequestStatus `json:"status"`
}

func (ctrl *FriendController) Block(ctx *gin.Context) {
	var actionBody ActionBody
	session := sessions.Default(ctx)

	if err := ctx.BindJSON(&actionBody); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.NewErrorResponse(http.StatusBadRequest, "Bad Request", err.Error()))
		return
	}
	userMail := session.Get("mail")
	if userMail == nil {
		ctx.JSON(http.StatusBadRequest, utils.NewErrorResponse(http.StatusBadRequest, "Bad Request", "UserId not found"))
		return
	}

	var friendObj models.Friend

	friendObj.UserMail = actionBody.Mail    // blocklanan
	friendObj.UserMail2 = userMail.(string) // blocklayan

	fmt.Print(actionBody.Mail)

	if err := ctrl.FriendService.Block(&friendObj); err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.NewErrorResponse(http.StatusInternalServerError, "Internal Server Error", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, utils.NewSuccessResponse(http.StatusOK, "OK", "success"))
}

func (ctrl *FriendController) Delete(ctx *gin.Context) {
	var actionBody ActionBody
	session := sessions.Default(ctx)

	if err := ctx.BindJSON(&actionBody); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.NewErrorResponse(http.StatusBadRequest, "Bad Request", err.Error()))
		return
	}

	userMail := session.Get("mail")
	if userMail == nil {
		ctx.JSON(http.StatusBadRequest, utils.NewErrorResponse(http.StatusBadRequest, "Bad Request", "UserId not found"))
		return
	}

	var friendObj models.Friend
	friendObj.UserMail = actionBody.Mail
	friendObj.UserMail2 = userMail.(string)

	if err := ctrl.FriendService.Delete(&friendObj); err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.NewErrorResponse(http.StatusInternalServerError, "Internal Server Error", err.Error()))
		return
	}

	ctx.JSON(http.StatusNoContent, utils.NewSuccessResponse(http.StatusNoContent, "No Content", "success"))
}
