package adapter

import (
	"github.com/kwa0x2/realtime-chat-backend/gateway"
	"github.com/kwa0x2/realtime-chat-backend/models"
	"github.com/kwa0x2/realtime-chat-backend/service"
	"github.com/zishang520/engine.io/utils"
	"github.com/zishang520/socket.io/socket"
)

type SocketAdapter struct {
	gateway        gateway.SocketGateway
	userSockets    map[string]string
	messageService *service.MessageService
}

func NewSocketAdapter(gateway gateway.SocketGateway, messageService *service.MessageService) *SocketAdapter {
	return &SocketAdapter{gateway: gateway, userSockets: make(map[string]string), messageService: messageService}
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

			var message models.Message

			message.MessageContent = data["message"].(string)
			message.MessageSenderID = connectedUserID
			message.MessageReceiverID = data["receiver_id"].(string)

			addedMessageData, err := adapter.messageService.InsertMessage(&message)
			if err != nil {
				utils.Log().Error(`while addding message error`)
				return
			}

			utils.Log().Info("Added and sended message %+v\n", addedMessageData)

			adapter.gateway.Emit("chat", adapter.userSockets[message.MessageReceiverID], map[string]interface{}{
				"sender_id": message.MessageSenderID,
				"message":   message.MessageContent,
			})

		})
	})
}
