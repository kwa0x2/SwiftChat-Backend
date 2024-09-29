package controller

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
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
		}
		responseData = append(responseData, responseItem)

	}
	ctx.JSON(http.StatusOK, responseData)
}
