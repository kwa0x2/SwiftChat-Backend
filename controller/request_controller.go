package controller

import (
	"github.com/kwa0x2/realtime-chat-backend/types"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/kwa0x2/realtime-chat-backend/models"
	"github.com/kwa0x2/realtime-chat-backend/service"
	"github.com/kwa0x2/realtime-chat-backend/utils"
)

type RequestController struct {
	RequestService *service.RequestService
	FriendService  *service.FriendService
	UserService    *service.UserService
}

type ActionBody struct {
	Mail   string              `json:"mail"`
	Status types.RequestStatus `json:"status"`
}

func (ctrl *RequestController) Insert(ctx *gin.Context) {
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

	if actionBody.Mail == userMail.(string) {
		ctx.JSON(http.StatusBadRequest, utils.NewErrorResponse(http.StatusBadRequest, "Bad Request", "You can't send friend request to yourself"))
		return
	}

	var requestObj models.Request

	requestObj.SenderMail = userMail.(string)
	requestObj.ReceiverMail = actionBody.Mail

	if err := ctrl.RequestService.Insert(nil, &requestObj); err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.NewErrorResponse(http.StatusInternalServerError, "Internal Server Error", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, utils.NewSuccessResponse(http.StatusOK, "OK", "success"))
}

func (ctrl *RequestController) GetComingRequests(ctx *gin.Context) {
	session := sessions.Default(ctx)

	userMail := session.Get("mail")
	if userMail == nil {
		ctx.JSON(http.StatusBadRequest, utils.NewErrorResponse(http.StatusBadRequest, "Bad Request", "UserId not found"))
		return
	}

	data, err := ctrl.RequestService.GetComingRequests(userMail.(string))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.NewErrorResponse(http.StatusInternalServerError, "Internal Server Error", err.Error()))
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
