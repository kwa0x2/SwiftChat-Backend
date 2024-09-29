package adapter

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/kwa0x2/realtime-chat-backend/models"
	"github.com/kwa0x2/realtime-chat-backend/service"
	"github.com/kwa0x2/realtime-chat-backend/socket/gateway"
	"github.com/kwa0x2/realtime-chat-backend/types"
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
	resendService  *service.ResendService
}

func NewSocketAdapter(gateway gateway.SocketGateway, messageService *service.MessageService, userService *service.UserService, friendService *service.FriendService, requestService *service.RequestService, resendService *service.ResendService) *SocketAdapter {
	return &SocketAdapter{gateway: gateway, userSockets: make(map[string]string), messageService: messageService, userService: userService, friendService: friendService, requestService: requestService, resendService: resendService}
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

			callback, ok := args[1].(func([]interface{}, error))
			if !ok {
				utils.Log().Error(`callback function type error socketid: %s`, socketio.Id())
				return
			}

			status := adapter.SendMessage(&messageObj, connectedUserMail, data["other_user_email"].(string))

			response := []interface{}{map[string]interface{}{"status": status}}
			callback(response, nil)
		})

		socketio.On("sendFriend", func(args ...any) {
			fmt.Println("status2")

			receiverMail, ok := args[0].(string)
			if !ok {
				utils.Log().Error(`socket message type error socketid: %s`, socketio.Id())
				return
			}

			if receiverMail != connectedUserMail {
				requestObj := models.Request{
					SenderMail:   connectedUserMail,
					ReceiverMail: receiverMail,
				}

				callback, ok := args[1].(func([]interface{}, error))
				if !ok {
					utils.Log().Error(`callback function type error socketid: %s`, socketio.Id())
					return
				}

				status := adapter.SendFriend(&requestObj, receiverMail)
				fmt.Println("status")
				response := []interface{}{map[string]interface{}{"status": status}}
				callback(response, nil)
			}
		})

		socketio.On("deleteMessage", func(args ...any) {
			data, ok := args[0].(map[string]interface{}) // message id, other_user_email ve room id
			if !ok {
				utils.Log().Error(`socket message type error socketid: %s`, socketio.Id())
				return
			}

			adapter.DeleteMessage(data["other_user_email"].(string), data["room_id"].(string), data["message_id"].(string))
		})

		socketio.On("editMessage", func(args ...any) {
			data, ok := args[0].(map[string]interface{}) // message id,new message,  other_user_email ve room id
			if !ok {
				utils.Log().Error(`socket message type error socketid: %s`, socketio.Id())
				return
			}

			adapter.EditMessage(data["other_user_email"].(string), data["room_id"].(string), data["message_id"].(string), data["edited_message"].(string))
		})

		socketio.On("blockFriend", func(args ...any) {
			data, ok := args[0].(map[string]interface{}) // friend_mail, friend_name
			if !ok {
				utils.Log().Error(`socket message type error socketid: %s`, socketio.Id())
				return
			}

			friendObj := models.Friend{
				UserMail:  data["friend_mail"].(string),
				UserMail2: connectedUserMail,
			}

			callback, ok := args[1].(func([]interface{}, error))
			if !ok {
				utils.Log().Error(`callback function type error socketid: %s`, socketio.Id())
				return
			}

			responseData := adapter.BlockFriend(&friendObj, data["friend_name"].(string))

			response := []interface{}{responseData}
			callback(response, nil)
		})

		socketio.On("deleteFriend", func(args ...any) {
			data, ok := args[0].(map[string]interface{}) // user_mail, user_name
			if !ok {
				utils.Log().Error(`socket message type error socketid: %s`, socketio.Id())
				return
			}

			friendObj := models.Friend{
				UserMail:  data["user_mail"].(string),
				UserMail2: connectedUserMail,
			}

			callback, ok := args[1].(func([]interface{}, error))
			if !ok {
				utils.Log().Error(`callback function type error socketid: %s`, socketio.Id())
				return
			}

			status := adapter.DeleteFriend(&friendObj, data["user_name"].(string))

			response := []interface{}{map[string]interface{}{"status": status}}
			callback(response, nil)
		})

		socketio.On("updateFriendshipRequest", func(args ...any) {
			data, ok := args[0].(map[string]interface{}) // sender mail ve status
			if !ok {
				utils.Log().Error(`socket message type error socketid: %s`, socketio.Id())
				return
			}

			requestObj := models.Request{
				SenderMail:    data["sender_mail"].(string),
				ReceiverMail:  connectedUserMail,
				RequestStatus: types.RequestStatus(data["status"].(string)),
			}

			callback, ok := args[1].(func([]interface{}, error))
			if !ok {
				utils.Log().Error(`callback function type error socketid: %s`, socketio.Id())
				return
			}

			status := adapter.UpdateFriendshipRequest(&requestObj)

			response := []interface{}{map[string]interface{}{"status": status}}
			callback(response, nil)
		})
	})
}

func (adapter *SocketAdapter) JoinRoom(socketio *socket.Socket, room string) {
	adapter.gateway.JoinRoom(socketio, room)
	utils.Log().Info("User %s joined room %s", socketio.Id(), room)
}

func (adapter *SocketAdapter) SendMessage(messageObj *models.Message, senderMail, receiverMail string) string {
	isBlocked, err := adapter.friendService.IsBlocked(senderMail, receiverMail)
	if err != nil {
		utils.Log().Error(`error while get blocked status `)
		return "error"
	}
	if isBlocked != false {
		utils.Log().Error(`friend is blocked `)
		return "error"
	}

	addedMessageData, messageErr := adapter.messageService.InsertAndUpdateRoom(messageObj)
	if messageErr != nil {
		utils.Log().Error(`error while adding message `)
		return "error"
	}

	utils.Log().Info("Added and sended message %+v\n", addedMessageData)
	adapter.EmitToRoomId("new_message", messageObj.RoomID.String(), addedMessageData)

	notifyData := map[string]interface{}{
		"room_id":   addedMessageData.RoomID,
		"message":   addedMessageData.Message,
		"sender_id": addedMessageData.SenderID,
		"updatedAt": addedMessageData.UpdatedAt,
	}
	fmt.Println("receiver", receiverMail)
	adapter.EmitToNotificationRoom("new_message", receiverMail, notifyData)
	return "success"
}

func (adapter *SocketAdapter) SendFriend(request *models.Request, receiverMail string) string {
	var pgErr *pgconn.PgError

	if isEmailExists := adapter.userService.IsEmailExists(receiverMail); !isEmailExists {
		if err := adapter.requestService.Insert(nil, request); err != nil {
			if errors.As(err, &pgErr) && pgErr.Code == "23505" {
				return "duplicate"
			}
			return "error"
		}

		_, err := adapter.resendService.SendMail(receiverMail, "You have received a new friend request from the SwiftChat app!", "friend_request")
		if err != nil {
			return "error"
		}

		return "email_sent"
	}

	data, err := adapter.requestService.InsertAndReturnUser(request)
	if err != nil {
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return "duplicate"
		}

		return "error"
	}

	adapter.EmitToNotificationRoom("friend_request", receiverMail, data)
	return "friend_sent"
}

func (adapter *SocketAdapter) EmitToNotificationRoom(notifyAction, receiverMail string, notifyObj any) {
	data := map[string]interface{}{
		"action": notifyAction,
		"data":   notifyObj,
	}

	utils.Log().Info("notify %+v\n mail:%s", data, receiverMail)

	adapter.gateway.EmitRoom("notification", receiverMail, data)
}

func (adapter *SocketAdapter) EmitToRoomId(notifyAction, roomId string, notifyObj any) {
	data := map[string]interface{}{
		"action": notifyAction,
		"data":   notifyObj,
	}

	utils.Log().Info("notify %+v\n roomId:%s", data, roomId)

	adapter.gateway.Emit(roomId, data)
}

func (adapter *SocketAdapter) DeleteMessage(connectedUserMail, roomId, messageId string) {
	if err := adapter.messageService.DeleteById(messageId); err != nil {
		utils.Log().Error(`error while deleting message `)
		return
	}

	utils.Log().Info("deleted message %+v\n", messageId)

	adapter.EmitToRoomId("delete_message", roomId, messageId)

	notifyData := map[string]interface{}{
		"room_id":    roomId,
		"message_id": messageId,
	}

	adapter.EmitToNotificationRoom("delete_message", connectedUserMail, notifyData)
}

func (adapter *SocketAdapter) EditMessage(connectedUserMail, roomId, messageId, editedMessage string) {
	if err := adapter.messageService.UpdateMessageByIdBody(messageId, editedMessage); err != nil {
		utils.Log().Error(`error while editing message `)
		return
	}

	utils.Log().Info("edited message %+v\n", messageId)

	notifyData := map[string]interface{}{
		"message_id":     messageId,
		"edited_message": editedMessage,
	}

	adapter.EmitToRoomId("edit_message", roomId, notifyData)

	adapter.EmitToNotificationRoom("edit_message", connectedUserMail, notifyData)
}

func (adapter *SocketAdapter) UpdateFriendshipRequest(request *models.Request) string {
	data, err := adapter.requestService.UpdateFriendshipRequest(request)
	if err != nil {
		utils.Log().Error(`error while editing message `)
		return "error"
	}

	utils.Log().Info("update friendship request %s", request)

	adapter.EmitToNotificationRoom("update_friendship_request", request.SenderMail, data)
	return "success"
}

func (adapter *SocketAdapter) BlockFriend(friend *models.Friend, friendName string) map[string]interface{} {
	friendObj := map[string]interface{}{
		"status":        "",
		"friend_status": "",
	}

	friendStatus, err := adapter.friendService.Block(friend)
	if err != nil {
		utils.Log().Error(`error while editing message `)
		friendObj["status"] = "error"
		return friendObj
	}

	data := map[string]interface{}{
		"friend_name":   friendName,
		"friend_mail":   friend.UserMail2,
		"friend_status": friendStatus,
	}

	utils.Log().Info("blocked user %s", friend)

	adapter.EmitToNotificationRoom("blocked_friend", friend.UserMail, data)
	friendObj["status"] = "success"
	friendObj["friend_status"] = friendStatus

	return friendObj
}

func (adapter *SocketAdapter) DeleteFriend(friend *models.Friend, userName string) string {
	if err := adapter.friendService.Delete(friend); err != nil {
		utils.Log().Error(`error while editing message `)
		return "error"
	}

	data := map[string]interface{}{
		"user_email": friend.UserMail2,
	}

	utils.Log().Info("deleted friend %s", friend)

	adapter.EmitToNotificationRoom("deleted_friend", friend.UserMail, data)
	return "success"
}
