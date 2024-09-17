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
	requestService *service.RequestService
}

func NewSocketAdapter(gateway gateway.SocketGateway, messageService *service.MessageService, userService *service.UserService, friendService *service.FriendService, requestService *service.RequestService) *SocketAdapter {
	return &SocketAdapter{gateway: gateway, userSockets: make(map[string]string), messageService: messageService, userService: userService, friendService: friendService, requestService: requestService}
}

func (adapter *SocketAdapter) HandleConnection() {
	adapter.gateway.OnConnection(func(socketio *socket.Socket) {
		ctx := socketio.Request().Context()
		connectedUserID := ctx.Value("id").(string)
		connectedUserMail := ctx.Value("mail").(string)

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

			messageObj := models.Message{
				SenderID: connectedUserID,
				Message:  data["message"].(string),
				RoomID:   roomID,
			}

			adapter.SendMessage(&messageObj, connectedUserMail, data["other_user_email"].(string))

		})

		socketio.On("sendFriend", func(emailData ...any) {
			receiverMail, ok := emailData[0].(string)
			if !ok {
				utils.Log().Error(`socket message type error socketid: %s`, socketio.Id())
				return
			}

			requestObj := models.Request{
				SenderMail:   connectedUserMail,
				ReceiverMail: receiverMail,
			}

			adapter.SendFriend(&requestObj, receiverMail)
		})

	})
}

func (adapter *SocketAdapter) JoinRoom(socketio *socket.Socket, room string) {
	adapter.gateway.JoinRoom(socketio, room)
	utils.Log().Info("User %s joined room %s", socketio.Id(), room)
}

func (adapter *SocketAdapter) SendMessage(messageObj *models.Message, senderMail, receiverMail string) {
	isBlocked, err := adapter.friendService.IsBlocked(senderMail, receiverMail)
	if err != nil {
		utils.Log().Error(`error while get blocked status `)
		return
	}
	if isBlocked != false {
		utils.Log().Error(`friend is blocked `)
		return
	}

	addedMessageData, messageErr := adapter.messageService.InsertAndUpdateRoom(messageObj)
	if messageErr != nil {
		utils.Log().Error(`error while adding message `)
		return
	}

	utils.Log().Info("Added and sended message %+v\n", addedMessageData)
	adapter.gateway.Emit(messageObj.RoomID.String(), addedMessageData)

	notifyData := map[string]interface{}{
		"room_id":   addedMessageData.RoomID,
		"message":   addedMessageData.Message,
		"sender_id": addedMessageData.SenderID,
		"updatedAt": addedMessageData.UpdatedAt,
	}

	adapter.SendNotification("message", receiverMail, notifyData)
}

func (adapter *SocketAdapter) SendFriend(request *models.Request, receiverMail string) {
	data, err := adapter.requestService.InsertAndReturnUser(request)
	if err != nil {
		utils.Log().Error(`error while sending friend request `)
		return
	}

	utils.Log().Info("successfully send friend request %s %+v\n", receiverMail, data)
	adapter.SendNotification("friend_request", receiverMail, data)

}

func (adapter *SocketAdapter) SendNotification(notifyType, receiverMail string, notifyObj any) {
	data := map[string]interface{}{
		"notification_type": notifyType,
		"data":              notifyObj,
	}

	utils.Log().Info("notify %+v\n mail:%s", data, receiverMail)

	adapter.gateway.EmitRoom("notification", receiverMail, data)
}
