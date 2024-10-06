package adapter

import (
	"fmt"
	"github.com/kwa0x2/realtime-chat-backend/service"
	"github.com/kwa0x2/realtime-chat-backend/socket/gateway"
	"github.com/zishang520/engine.io/utils"
	"github.com/zishang520/socket.io/socket"
	"sync"
)

type SocketAdapter struct {
	Gateway          *gateway.SocketGateway
	onlineUserEmails []string
	MessageService   *service.MessageService
	FriendService    *service.FriendService
	mux              sync.RWMutex
}

func (adapter *SocketAdapter) HandleConnection() {
	adapter.Gateway.OnConnection(func(socketio *socket.Socket) {
		ctx := socketio.Request().Context()
		connectedUserID := ctx.Value("id").(string)
		connectedUserMail := ctx.Value("mail").(string)
		fmt.Println(connectedUserMail, " is  online")

		if !adapter.emailExists(connectedUserMail) {
			adapter.onlineUserEmails = append(adapter.onlineUserEmails, connectedUserMail)
		}

		adapter.broadcastOnlineUsers()

		socketio.On("disconnect", func(...any) {
			adapter.handleDisconnect(connectedUserMail)
		})

		socketio.On("joinRoom", func(roomData ...any) {
			adapter.handleJoinRoom(socketio, roomData...)
		})

		socketio.On("sendMessage", func(args ...any) {
			adapter.handleSendMessage(connectedUserID, connectedUserMail, args...)
		})

		socketio.On("deleteMessage", func(args ...any) {
			adapter.handleDeleteMessage(args...)
		})

		socketio.On("editMessage", func(args ...any) {
			adapter.handleEditMessage(args...)
		})

		socketio.On("starMessage", func(args ...any) {
			adapter.handleStarMessage(args...)
		})
		socketio.On("readMessage", func(args ...any) {
			adapter.handleReadMessage(connectedUserID, args...)
		})
	})
}

func (adapter *SocketAdapter) handleJoinRoom(socketio *socket.Socket, roomData ...any) {
	roomId, ok := roomData[0].(string)
	if !ok {
		utils.Log().Error(`socket message type error socketid: %s `, socketio.Id())
		return
	}
	adapter.JoinRoom(socketio, roomId)
}

func (adapter *SocketAdapter) JoinRoom(socketio *socket.Socket, room string) {
	adapter.Gateway.JoinRoom(socketio, room)
	utils.Log().Info("User %s joined room %s", socketio.Id(), room)
}

func (adapter *SocketAdapter) broadcastOnlineUsers() {

	adapter.Gateway.Emit("onlineUsers", adapter.onlineUserEmails)

}

func (adapter *SocketAdapter) emailExists(email string) bool {
	for _, existingEmail := range adapter.onlineUserEmails {
		if existingEmail == email {
			return true
		}
	}
	return false
}

func (adapter *SocketAdapter) handleDisconnect(email string) {

	for i, existingEmail := range adapter.onlineUserEmails {
		if existingEmail == email {
			adapter.onlineUserEmails = append(adapter.onlineUserEmails[:i], adapter.onlineUserEmails[i+1:]...)
			fmt.Println(email, " is offline")
			adapter.Gateway.Emit("onlineUsers", adapter.onlineUserEmails)

			break
		}
	}

	adapter.broadcastOnlineUsers()
}
