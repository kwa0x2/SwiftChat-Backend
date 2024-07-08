package adapter

import (
	"github.com/kwa0x2/realtime-chat-backend/gateway"
	"github.com/zishang520/engine.io/utils"
	"github.com/zishang520/socket.io/socket"
)

type SocketAdapter struct {
	gateway gateway.SocketGateway
	userSockets map[string]string
}

func NewSocketAdapter(gateway gateway.SocketGateway) *SocketAdapter {
	return &SocketAdapter{gateway: gateway, userSockets: make(map[string]string)}
}

func (adapter *SocketAdapter) HandleConnection() {
	adapter.gateway.OnConnection(func(socketio *socket.Socket) {
		ctx := socketio.Request().Context()
		senderID := ctx.Value("id").(string)

		utils.Log().Info(`socket connected %s user id %s`, socketio.Id(), senderID)

		adapter.userSockets[senderID] = string(socketio.Id())


		socketio.On("disconnect", func(reason ...any) {
			utils.Log().Info(`socket disconnected %s user id %s`, socketio.Id(), senderID)
		})

		socketio.On("sendMessage", func(args ...any) {
			data, ok := args[0].(map[string]interface{})
			if !ok {
				utils.Log().Error(`socket message type error %s user id %s`, socketio.Id(), senderID)
				return
			}

			receiverID := data["receiver_id"].(string)
			message := data["message"].(string)

			utils.Log().Info("Received message from %s : %s", senderID, message)

			adapter.gateway.Emit("chat", adapter.userSockets[receiverID], map[string]interface{}{
				"sender_id": senderID,
				"message":   message,
			})

			

		})
	})
}
