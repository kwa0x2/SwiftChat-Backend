package adapter

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/kwa0x2/realtime-chat-backend/models"
	"github.com/kwa0x2/realtime-chat-backend/utils"
)

func (adapter *SocketAdapter) handleSendMessage(connectedUserID, connectedUserMail string, args ...any) {
	data, callback := utils.ExtractArgs(args)
	if data == nil || callback == nil {
		utils.SendResponse(callback, "error", "Invalid arguments")
		return
	}

	roomID, err := uuid.Parse(data["room_id"].(string))
	if err != nil {
		utils.LogError(callback, "invalid room_id format")
		return
	}

	messageObj := models.Message{
		SenderID: connectedUserID,
		Message:  data["message"].(string),
		RoomID:   roomID,
	}

	addedMessageId, err := adapter.SendMessage(&messageObj, connectedUserMail, data["user_email"].(string))
	if err != nil {
		utils.LogError(callback, err.Error())
		return
	}

	utils.LogSuccess(callback, addedMessageId)
}

func (adapter *SocketAdapter) SendMessage(messageObj *models.Message, senderMail, receiverMail string) (string, error) {
	isBlocked, err := adapter.FriendService.IsBlocked(senderMail, receiverMail)
	if err != nil {
		return "", err
	}
	if isBlocked {
		return "", errors.New("friend is blocked")
	}

	addedMessageData, messageErr := adapter.MessageService.InsertAndUpdateRoom(messageObj)
	if messageErr != nil {
		return "", messageErr
	}

	adapter.EmitToRoomId("new_message", messageObj.RoomID.String(), addedMessageData)

	notifyData := map[string]interface{}{
		"room_id":    addedMessageData.RoomID,
		"message":    addedMessageData.Message,
		"message_id": addedMessageData.MessageID,
		"updatedAt":  addedMessageData.UpdatedAt,
	}

	adapter.EmitToNotificationRoom("new_message", receiverMail, notifyData)
	return addedMessageData.MessageID.String(), nil
}

func (adapter *SocketAdapter) handleDeleteMessage(args ...any) {
	data, callback := utils.ExtractArgs(args)
	if data == nil || callback == nil {
		utils.SendResponse(callback, "error", "Invalid arguments")
		return
	}

	err := adapter.DeleteMessage(data["user_email"].(string), data["room_id"].(string), data["message_id"].(string))
	if err != nil {
		utils.LogError(callback, err.Error())
		return
	}

	utils.LogSuccess(callback, "Message deleted successfully")
}

func (adapter *SocketAdapter) DeleteMessage(connectedUserMail, roomId, messageId string) error {
	if err := adapter.MessageService.DeleteById(messageId); err != nil {
		return err
	}

	adapter.EmitToRoomId("delete_message", roomId, messageId)

	notifyData := map[string]interface{}{
		"room_id":    roomId,
		"message_id": messageId,
	}
	adapter.EmitToNotificationRoom("delete_message", connectedUserMail, notifyData)
	return nil
}

func (adapter *SocketAdapter) handleEditMessage(args ...any) {
	data, callback := utils.ExtractArgs(args)
	if data == nil || callback == nil {
		utils.SendResponse(callback, "error", "Invalid arguments")
		return
	}

	err := adapter.EditMessage(data["user_email"].(string), data["room_id"].(string), data["message_id"].(string), data["edited_message"].(string))
	if err != nil {
		utils.LogError(callback, err.Error())
		return
	}
	utils.LogSuccess(callback, "Message edited successfully")
}

func (adapter *SocketAdapter) EditMessage(connectedUserMail, roomId, messageId, editedMessage string) error {
	if err := adapter.MessageService.UpdateMessageByIdBody(messageId, editedMessage); err != nil {
		return err
	}

	notifyData := map[string]interface{}{
		"message_id":     messageId,
		"edited_message": editedMessage,
	}

	adapter.EmitToRoomId("edit_message", roomId, notifyData)
	adapter.EmitToNotificationRoom("edit_message", connectedUserMail, notifyData)
	return nil
}

func (adapter *SocketAdapter) handleStarMessage(args ...any) {
	data, callback := utils.ExtractArgs(args)
	if data == nil || callback == nil {
		utils.SendResponse(callback, "error", "Invalid arguments")
		return
	}

	err := adapter.StarMessage(data["room_id"].(string), data["message_id"].(string))
	if err != nil {
		utils.LogError(callback, err.Error())
		return
	}
	utils.LogSuccess(callback, "Message starred successfully")
}

func (adapter *SocketAdapter) StarMessage(roomId, messageId string) error {
	if err := adapter.MessageService.StarMessageById(messageId); err != nil {
		return err
	}

	notifyData := map[string]interface{}{
		"message_id": messageId,
	}

	adapter.EmitToRoomId("star_message", roomId, notifyData)
	return nil
}

func (adapter *SocketAdapter) handleReadMessage(connectedUserID string, args ...any) {
	fmt.Println("read message")

	data, callback := utils.ExtractArgs(args)
	if data == nil || callback == nil {
		utils.SendResponse(callback, "error", "Invalid arguments")
		return
	}
	messageId, okMessage := data["message_id"].(string)

	if !okMessage || messageId == "" {
		err := adapter.ReadMessage(connectedUserID, data["room_id"].(string), nil)
		if err != nil {
			utils.LogError(callback, err.Error())
			return
		}
		utils.LogSuccess(callback, "Room read successfully without message ID")
		return
	}

	err := adapter.ReadMessage(connectedUserID, data["room_id"].(string), &messageId)
	if err != nil {
		utils.LogError(callback, err.Error())
		return
	}
	utils.LogSuccess(callback, "Message read successfully with message ID")
}

func (adapter *SocketAdapter) ReadMessage(connectedUserID, roomId string, messageId *string) error {
	if err := adapter.MessageService.ReadMessageByRoomId(connectedUserID, roomId, messageId); err != nil {
		return err
	}

	notifyData := map[string]interface{}{
		"room_id": roomId,
	}

	adapter.EmitToRoomId("read_message", roomId, notifyData)
	return nil
}
