package controller

import (
	"errors"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/kwa0x2/realtime-chat-backend/helpers"
	"github.com/kwa0x2/realtime-chat-backend/models"
	"github.com/kwa0x2/realtime-chat-backend/service"
	"gorm.io/gorm"
)

type FriendshipController struct {
	FriendshipService *service.FriendshipService
	UserService       *service.UserService
}

type SendFriendRequestBody struct {
	Email string `json:"email"`
}

func (ctrl *FriendshipController) SendFriendRequest(ctx *gin.Context) {
	var requestBody SendFriendRequestBody
	session := sessions.Default(ctx)

	if err := ctx.BindJSON(&requestBody); err != nil {
		ctx.JSON(http.StatusBadRequest, helpers.NewErrorResponse(http.StatusBadRequest, "Bad Request", err.Error()))
		return
	}

	userId := session.Get("id")
	if userId == nil {
		ctx.JSON(http.StatusBadRequest, helpers.NewErrorResponse(http.StatusBadRequest, "Bad Request", "UserId not found"))
		return
	}

	user, err := ctrl.UserService.GetByEmail(requestBody.Email)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		ctx.JSON(http.StatusOK, helpers.NewErrorResponse(http.StatusOK, "OK", "sent mail"))
		return
	}

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helpers.NewErrorResponse(http.StatusInternalServerError, "Internal Server Error", err.Error()))
		return
	}

	var friendshipObj models.Friendship

	friendshipObj.SenderId = userId.(string)
	friendshipObj.ReceiverId = user.UserID

	data, err := ctrl.FriendshipService.SendFriendRequest(&friendshipObj)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helpers.NewErrorResponse(http.StatusInternalServerError, "Internal Server Error", err.Error()))
		return
	}

	ctx.JSON(http.StatusCreated, helpers.NewSuccessResponse(http.StatusCreated, "CREATED", data))
}

func (ctrl *FriendshipController) GetComingRequests(ctx *gin.Context) {
	session := sessions.Default(ctx)

	userId := session.Get("id")
	if userId == nil {
		ctx.JSON(http.StatusBadRequest, helpers.NewErrorResponse(http.StatusBadRequest, "Bad Request", "UserId not found"))
		return
	}

	data, err := ctrl.FriendshipService.GetComingRequests(userId.(string))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helpers.NewErrorResponse(http.StatusInternalServerError, "Internal Server Error", err.Error()))
	}

	var responseData []map[string]interface{}
	for _, item := range data {
		responseItem := map[string]interface{}{
			"sender_id": item.SenderId,
			"user_name":   item.User.UserName,
			"user_photo":  item.User.UserPhoto,
		}
		responseData = append(responseData, responseItem)
	}

	ctx.JSON(http.StatusOK, responseData)
}

func (ctrl *FriendshipController) GetFriends(ctx *gin.Context) {
	session := sessions.Default(ctx)

	userId := session.Get("id")
	if userId == nil {
		ctx.JSON(http.StatusBadRequest, helpers.NewErrorResponse(http.StatusBadRequest, "Bad Request", "UserId not found"))
		return
	}

	data, err := ctrl.FriendshipService.GetFriends(userId.(string))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helpers.NewErrorResponse(http.StatusInternalServerError, "Internal Server Error", err.Error()))
	}

	var responseData []map[string]interface{}
	for _, item := range data {
		responseItem := map[string]interface{}{
			"sender_id": item.SenderId,
			"user_name":   item.User.UserName,
			"user_photo":  item.User.UserPhoto,
		}
		responseData = append(responseData, responseItem)

	}
	ctx.JSON(http.StatusOK, responseData)
}

func (ctrl *FriendshipController) GetBlockeds(ctx *gin.Context) {
	session := sessions.Default(ctx)

	userId := session.Get("id")
	if userId == nil {
		ctx.JSON(http.StatusBadRequest, helpers.NewErrorResponse(http.StatusBadRequest, "Bad Request", "UserId not found"))
		return
	}

	data, err := ctrl.FriendshipService.GetBlockeds(userId.(string))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helpers.NewErrorResponse(http.StatusInternalServerError, "Internal Server Error", err.Error()))
	}

	var responseData []map[string]interface{}
	for _, item := range data {
		responseItem := map[string]interface{}{
			"receiver_id": item.ReceiverId,
			"user_name":   item.User.UserName,
			"user_photo":  item.User.UserPhoto,
		}
		responseData = append(responseData, responseItem)

	}
	ctx.JSON(http.StatusOK, responseData)
}

type ActionBody struct {
	FriendId string `json:"friend_id"`
}

func (ctrl *FriendshipController) Block(ctx *gin.Context) {
	var requestBody ActionBody
	session := sessions.Default(ctx)

	if err := ctx.BindJSON(&requestBody); err != nil {
		ctx.JSON(http.StatusBadRequest, helpers.NewErrorResponse(http.StatusBadRequest, "Bad Request", err.Error()))
		return
	}
	userId := session.Get("id")
	if userId == nil {
		ctx.JSON(http.StatusBadRequest, helpers.NewErrorResponse(http.StatusBadRequest, "Bad Request", "UserId not found"))
		return
	}

	var friendshipObj models.Friendship

	friendshipObj.SenderId = requestBody.FriendId
	friendshipObj.ReceiverId = userId.(string)

	if err := ctrl.FriendshipService.Block(&friendshipObj); err != nil {
		ctx.JSON(http.StatusInternalServerError, helpers.NewErrorResponse(http.StatusInternalServerError, "Internal Server Error", err.Error()))
	}

	ctx.JSON(http.StatusOK, helpers.NewSuccessResponse(http.StatusOK, "OK", "success"))
}

func (ctrl *FriendshipController) Delete(ctx *gin.Context) {
	var requestBody ActionBody
	session := sessions.Default(ctx)

	if err := ctx.BindJSON(&requestBody); err != nil {
		ctx.JSON(http.StatusBadRequest, helpers.NewErrorResponse(http.StatusBadRequest, "Bad Request", err.Error()))
		return
	}
	userId := session.Get("id")
	if userId == nil {
		ctx.JSON(http.StatusBadRequest, helpers.NewErrorResponse(http.StatusBadRequest, "Bad Request", "UserId not found"))
		return
	}

	var friendshipObj models.Friendship

	friendshipObj.SenderId = requestBody.FriendId
	friendshipObj.ReceiverId = userId.(string)

	if err := ctrl.FriendshipService.Delete(&friendshipObj); err != nil {
		ctx.JSON(http.StatusInternalServerError, helpers.NewErrorResponse(http.StatusInternalServerError, "Internal Server Error", err.Error()))
	}

	ctx.JSON(http.StatusNoContent, helpers.NewSuccessResponse(http.StatusNoContent, "No Content", "success"))
}

func (ctrl *FriendshipController) Accept(ctx *gin.Context) {
	var requestBody ActionBody
	session := sessions.Default(ctx)

	if err := ctx.BindJSON(&requestBody); err != nil {
		ctx.JSON(http.StatusBadRequest, helpers.NewErrorResponse(http.StatusBadRequest, "Bad Request", err.Error()))
		return
	}
	userId := session.Get("id")
	if userId == nil {
		ctx.JSON(http.StatusBadRequest, helpers.NewErrorResponse(http.StatusBadRequest, "Bad Request", "UserId not found"))
		return
	}

	var friendshipObj models.Friendship

	friendshipObj.SenderId = requestBody.FriendId
	friendshipObj.ReceiverId = userId.(string)

	if err := ctrl.FriendshipService.Accept(&friendshipObj); err != nil {
		ctx.JSON(http.StatusInternalServerError, helpers.NewErrorResponse(http.StatusInternalServerError, "Internal Server Error", err.Error()))
	}

	ctx.JSON(http.StatusOK, helpers.NewSuccessResponse(http.StatusOK, "OK", "success"))
}

func (ctrl *FriendshipController) Reject(ctx *gin.Context) {
	var requestBody ActionBody
	session := sessions.Default(ctx)

	if err := ctx.BindJSON(&requestBody); err != nil {
		ctx.JSON(http.StatusBadRequest, helpers.NewErrorResponse(http.StatusBadRequest, "Bad Request", err.Error()))
		return
	}
	userId := session.Get("id")
	if userId == nil {
		ctx.JSON(http.StatusBadRequest, helpers.NewErrorResponse(http.StatusBadRequest, "Bad Request", "UserId not found"))
		return
	}

	var friendshipObj models.Friendship

	friendshipObj.SenderId = requestBody.FriendId
	friendshipObj.ReceiverId = userId.(string)

	if err := ctrl.FriendshipService.Reject(&friendshipObj); err != nil {
		ctx.JSON(http.StatusInternalServerError, helpers.NewErrorResponse(http.StatusInternalServerError, "Internal Server Error", err.Error()))
	}

	ctx.JSON(http.StatusOK, helpers.NewSuccessResponse(http.StatusOK, "OK", "success"))
}
