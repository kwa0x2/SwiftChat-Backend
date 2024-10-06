package controller

import (
	"github.com/kwa0x2/realtime-chat-backend/models"
	"github.com/kwa0x2/realtime-chat-backend/socket/adapter"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/kwa0x2/realtime-chat-backend/service"
	"github.com/kwa0x2/realtime-chat-backend/utils"
)

type FriendController struct {
	FriendService  *service.FriendService
	UserService    *service.UserService
	RequestService *service.RequestService
	SocketAdapter  *adapter.SocketAdapter
}

func (ctrl *FriendController) GetFriends(ctx *gin.Context) {
	session := sessions.Default(ctx)

	userMail := session.Get("mail")
	if userMail == nil {
		ctx.JSON(http.StatusBadRequest, utils.NewErrorResponse("Session Error", "UserMail not found"))
		return
	}

	data, err := ctrl.FriendService.GetFriends(userMail.(string), false)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.NewErrorResponse("Internal Server Error", "arkadaslar cekilirken sorun olustu"))
		return
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
		ctx.JSON(http.StatusBadRequest, utils.NewErrorResponse("Session Error", "UserMail not found"))
		return
	}

	data, err := ctrl.FriendService.GetBlocked(userMail.(string))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.NewErrorResponse("Internal Server Error", "engellenek kullanicilar cekilirken sorun olustu"))
		return
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

func (ctrl *FriendController) Block(ctx *gin.Context) {
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

	friend := models.Friend{
		UserMail:  actionBody.Mail,   // blocklanan
		UserMail2: userMail.(string), // blocklayan
	}

	friendStatus, err := ctrl.FriendService.Block(&friend)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.NewErrorResponse("Internal Server Error", "kullanici blocklanirken bir sorun olustu"))
		return
	}

	data := map[string]interface{}{
		"friend_mail":   friend.UserMail2,
		"friend_status": friendStatus,
	}

	ctrl.SocketAdapter.EmitToNotificationRoom("blocked_friend", friend.UserMail, data)
	ctx.JSON(http.StatusOK, utils.NewSuccessResponse("Success", "basariyla engellendi"))
}

func (ctrl *FriendController) Delete(ctx *gin.Context) {
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

	friendObj := models.Friend{
		UserMail:  actionBody.Mail,
		UserMail2: userMail.(string),
	}

	if err := ctrl.FriendService.Delete(&friendObj); err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.NewErrorResponse("Delete Friend Error", "kullanici silinirken hata olustu"))
		return
	}

	data := map[string]interface{}{
		"user_email": userMail.(string),
	}

	ctrl.SocketAdapter.EmitToNotificationRoom("deleted_friend", actionBody.Mail, data)
	ctx.JSON(http.StatusOK, utils.NewSuccessResponse("Success", "basariyla silindi"))
}
