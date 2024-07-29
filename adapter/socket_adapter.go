package adapter

import (
	"github.com/google/uuid"
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
	userService    *service.UserService
	friendService  *service.FriendService
}

func NewSocketAdapter(gateway gateway.SocketGateway, messageService *service.MessageService, userService *service.UserService, friendService *service.FriendService) *SocketAdapter {
	return &SocketAdapter{gateway: gateway, userSockets: make(map[string]string), messageService: messageService, userService: userService, friendService: friendService}
}

func (adapter *SocketAdapter) HandleConnection() {
	adapter.gateway.OnConnection(func(socketio *socket.Socket) {
		ctx := socketio.Request().Context()
		connectedUserID := ctx.Value("id").(string)

		utils.Log().Info("new connection established socketid: %s userid: %s", socketio.Id(), connectedUserID)

		adapter.userSockets[connectedUserID] = string(socketio.Id())

		socketio.On("joinRoom", func(roomData ...any) {
			roomId, ok := roomData[0].(string)
			if !ok {
				utils.Log().Error(`socket message type error socketid: %s `, socketio.Id())
				return
			}
			adapter.JoinRoom(socketio, roomId)
		})

		socketio.On("sendMessage", func(args ...any) {
			data, ok := args[0].(map[string]interface{})
			if !ok {
				utils.Log().Error(`socket message type error socketid: %s`, socketio.Id())
				return
			}

			roomID, err := uuid.Parse(data["room_id"].(string))
			if err != nil {
				utils.Log().Error("invalid room_id format")
				return
			}

			var messageObj models.Message

			messageObj.SenderID = connectedUserID
			messageObj.Message = data["message"].(string)
			messageObj.RoomID = roomID

			adapter.SendMessage(&messageObj)
		})

	})
}

func (adapter *SocketAdapter) JoinRoom(socketio *socket.Socket, room string) {
	adapter.gateway.JoinRoom(socketio, room)
	utils.Log().Info("User %s joined room %s", socketio.Id(), room)
}

func (adapter *SocketAdapter) SendMessage(messageObj *models.Message) {
	addedMessageData, err := adapter.messageService.InsertAndUpdateRoom(messageObj)
	if err != nil {
		utils.Log().Error(`error while adding message `)
		return
	}

	utils.Log().Info("Added and sended message %+v\n", addedMessageData)
	adapter.gateway.Emit(messageObj.RoomID.String(), addedMessageData)
}
