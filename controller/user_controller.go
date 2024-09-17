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
}

type UsernameUpdateBody struct {
	UserName string `json:"user_name"`
}

func (ctrl *UserController) Update(ctx *gin.Context) {
	var requestBody UsernameUpdateBody
	session := sessions.Default(ctx)

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
	session.Clear()
	session.Set("name", "asdasdas")
	session.Save()
	//if err := session.Save(); err != nil {
	//	ctx.JSON(http.StatusInternalServerError, utils.NewErrorResponse(http.StatusInternalServerError, "Internal Server Error", err.Error()))
	//	return
	//}
	fmt.Println(session.Get("name"))
	ctx.JSON(http.StatusOK, utils.NewSuccessResponse(http.StatusOK, "OK", "success"))
}
