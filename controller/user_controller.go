package controller

import (
	"github.com/kwa0x2/realtime-chat-backend/service"
)

type UserController struct {
	UserService *service.UserService
}
