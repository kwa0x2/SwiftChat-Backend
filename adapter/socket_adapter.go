package adapter

import (

	"github.com/kwa0x2/realtime-chat-backend/gateway"
	"github.com/kwa0x2/realtime-chat-backend/models"
	"github.com/kwa0x2/realtime-chat-backend/service"
	"github.com/zishang520/engine.io/utils"
	"github.com/zishang520/socket.io/socket"
)

type SocketAdapter struct {
	gateway           gateway.SocketGateway
	userSockets       map[string]string
	messageService    *service.MessageService
	userService       *service.UserService
	friendshipService *service.FriendshipService
}

func NewSocketAdapter(gateway gateway.SocketGateway, messageService *service.MessageService, userService *service.UserService, friendshipService *service.FriendshipService) *SocketAdapter {
	return &SocketAdapter{gateway: gateway, userSockets: make(map[string]string), messageService: messageService, userService: userService, friendshipService: friendshipService}
}

func (adapter *SocketAdapter) HandleConnection() {
	adapter.gateway.OnConnection(func(socketio *socket.Socket) {
		ctx := socketio.Request().Context()
		connectedUserID := ctx.Value("id").(string)

		utils.Log().Info(`socket connected %s user id %s`, socketio.Id(), connectedUserID)

		adapter.userSockets[connectedUserID] = string(socketio.Id())

		socketio.On("disconnect", func(reason ...any) {
			utils.Log().Info(`socket disconnected %s user id %s`, socketio.Id(), connectedUserID)
		})

		socketio.On("sendMessage", func(args ...any) {
			data, ok := args[0].(map[string]interface{})
			if !ok {
				utils.Log().Error(`socket message type error %s user id %s`, socketio.Id(), connectedUserID)
				return
			}

			var messageObj models.Message

			messageObj.MessageContent = data["message"].(string)
			messageObj.MessageSenderID = connectedUserID
			messageObj.MessageReceiverID = data["DestionationUserId"].(string)

			addedMessageData, err := adapter.messageService.Insert(&messageObj)
			if err != nil {
				utils.Log().Error(`while addding message error`)
				return
			}

			utils.Log().Info("Added and sended message %+v\n", addedMessageData)

			// direkt eklenen veri donucek
			adapter.gateway.Emit("chat", adapter.userSockets[messageObj.MessageReceiverID], map[string]interface{}{
				"sender_id": messageObj.MessageSenderID,
				"message":   messageObj.MessageContent,
			})

		})

		// socketio.On("sendFriendship", func(args ...any) {
		// 	data, ok := args[0].(map[string]interface{})
		// 	if !ok {
		// 		utils.Log().Error(`socket message type error %s user id %s`, socketio.Id(), connectedUserID)
		// 		return
		// 	}
		// 	utils.Log().Printf(`socket message type error %s user id %s email %s`, socketio.Id(), connectedUserID, data["email"].(string))

		// 	user, err := adapter.userService.GetByEmail(data["email"].(string))

		// 	if errors.Is(err, gorm.ErrRecordNotFound) {
		// 		utils.Log().Printf("User with email %s not found", data["email"].(string))

		// 		adapter.gateway.Emit("response", adapter.userSockets[connectedUserID], map[string]interface{}{
		// 			"status":  "success",
		// 			"message": "Friendship email sent successfully",
		// 		})
		// 		return
		// 	}
			
		// 	if err != nil {
		// 		utils.Log().Error(`while getting user error`)
		// 		return
		// 	}

		// 	var friendshipObj models.Friendship

		// 	friendshipObj.SenderId = connectedUserID
		// 	friendshipObj.ReceiverId = user.UserID
		// 	friendshipObj.FriendshipStatus = "pending"

		// 	friendshipStatus, err := adapter.friendshipService.SendFriendRequest(&friendshipObj)
		// 	if err != nil {
		// 		utils.Log().Error(`while adding friendship error`)
		// 		return
		// 	}

		// 	adapter.gateway.Emit("friendship", adapter.userSockets[friendshipObj.ReceiverId], map[string]interface{}{
		// 		"income_friendship_sender_id": friendshipObj.SenderId,
		// 		"income_friendship_status":    friendshipStatus,
		// 	})

		// 	adapter.gateway.Emit("response", adapter.userSockets[connectedUserID], map[string]interface{}{
		// 		"status":  "success",
		// 		"message": "Friendship request sent successfully",
		// 	})

		// })
	})
}
