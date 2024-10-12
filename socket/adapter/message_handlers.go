package adapter

import (
	"errors"
	"github.com/getsentry/sentry-go"
	"github.com/google/uuid"
	"github.com/kwa0x2/swiftchat-backend/models"
	"github.com/kwa0x2/swiftchat-backend/types"
	"github.com/kwa0x2/swiftchat-backend/utils"
)

// region "handleSendMessage" processes sending a message in a chat room.
func (adapter *socketAdapter) handleSendMessage(connectedUserID, connectedUserMail string, args ...any) {
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

	// Create a message object with sender ID, message content, and room ID.
	messageObj := models.Message{
		SenderID:       connectedUserID,
		MessageContent: data["message_content"].(string),
		RoomID:         roomID,
		MessageType:    types.MessageType(data["message_type"].(string)),
	}

	addedMessageId, sendErr := adapter.SendMessage(&messageObj, connectedUserMail, data["user_email"].(string))
	if sendErr != nil {
		utils.LogError(callback, sendErr.Error())
		return
	}

	utils.LogSuccess(callback, addedMessageId)
}

// endregion

// region "SendMessage" handles the actual sending of a message and notification to the recipient.
func (adapter *socketAdapter) SendMessage(messageObj *models.Message, senderMail, receiverMail string) (string, error) {
	// Check if the sender has blocked the receiver.
	isBlocked, err := adapter.FriendService.IsBlocked(senderMail, receiverMail)
	if err != nil {
		return "", err
	}
	if isBlocked {
		return "", errors.New("friend is blocked") // Return error if blocked.
	}

	// Insert the message and update the room.
	addedMessageData, messageErr := adapter.MessageService.InsertAndUpdateRoom(messageObj)
	if messageErr != nil {
		return "", messageErr
	}

	// Prepare notification data to send to the recipient.
	notifyData := map[string]interface{}{
		"room_id":         addedMessageData.RoomID,
		"message_content": addedMessageData.MessageContent,
		"message_id":      addedMessageData.MessageID,
		"updatedAt":       addedMessageData.UpdatedAt,
		"message_type":    addedMessageData.MessageType,
	}

	// Emit new message event to the chat room.
	adapter.Gateway.EmitToRoomId("new_message", messageObj.RoomID.String(), addedMessageData)
	// Emit notification of the new message to the recipient.
	adapter.Gateway.EmitToNotificationRoom("new_message", receiverMail, notifyData)
	return addedMessageData.MessageID.String(), nil
}

// endregion

// region "handleDeleteMessage" processes message deletion requests.
func (adapter *socketAdapter) handleDeleteMessage(args ...any) {
	data, callback := utils.ExtractArgs(args)
	if data == nil || callback == nil {
		utils.SendResponse(callback, "error", "Invalid arguments")
		return
	}

	messageID, err := uuid.Parse(data["message_id"].(string))
	if err != nil {
		utils.LogError(callback, "invalid message_id format")
		return
	}

	deleteErr := adapter.DeleteMessage(data["user_email"].(string), data["room_id"].(string), messageID)
	if deleteErr != nil {
		utils.LogError(callback, deleteErr.Error())
		return
	}

	utils.LogSuccess(callback, "Message deleted successfully")
}

// endregion

// region "DeleteMessage" handles the deletion of a message from the database and notifies clients.
func (adapter *socketAdapter) DeleteMessage(connectedUserMail, roomId string, messageId uuid.UUID) error {
	// Delete the message by its ID.
	if err := adapter.MessageService.DeleteById(messageId); err != nil {
		return err
	}

	notifyData := map[string]interface{}{
		"room_id":    roomId,
		"message_id": messageId,
	}

	// Emit message deletion event to the chat room.
	adapter.Gateway.EmitToRoomId("delete_message", roomId, messageId)
	// Emit notification of the deleted message to the user.
	adapter.Gateway.EmitToNotificationRoom("delete_message", connectedUserMail, notifyData)
	return nil
}

// endregion

// region "handleEditMessage" processes message edit requests.
func (adapter *socketAdapter) handleEditMessage(args ...any) {
	data, callback := utils.ExtractArgs(args)
	if data == nil || callback == nil {
		utils.SendResponse(callback, "error", "Invalid arguments")
		return
	}

	messageID, err := uuid.Parse(data["message_id"].(string))
	if err != nil {
		sentry.CaptureException(err)
		utils.LogError(callback, "invalid message_id format")
		return
	}

	editErr := adapter.EditMessage(data["user_email"].(string), data["room_id"].(string), data["edited_message"].(string), messageID)
	if editErr != nil {
		sentry.CaptureException(editErr)
		utils.LogError(callback, editErr.Error())
		return
	}
	utils.LogSuccess(callback, "Message edited successfully")
}

// endregion

// region "EditMessage" updates a message in the database and notifies clients of the change.
func (adapter *socketAdapter) EditMessage(connectedUserMail, roomId, editedMessage string, messageId uuid.UUID) error {
	// Update the message by its ID.
	if err := adapter.MessageService.UpdateMessageById(messageId, editedMessage); err != nil {
		sentry.CaptureException(err)
		return err
	}

	// Prepare notification data for the edited message.
	notifyData := map[string]interface{}{
		"message_id":     messageId,
		"edited_message": editedMessage,
	}

	// Emit message edit event to the chat room.
	adapter.Gateway.EmitToRoomId("edit_message", roomId, notifyData)
	// Emit notification of the edited message to the user.
	adapter.Gateway.EmitToNotificationRoom("edit_message", connectedUserMail, notifyData)
	return nil
}

// endregion

// region "handleUpdateMessageStarred" processes requests to update message type.
func (adapter *socketAdapter) handleUpdateMessageStarred(args ...any) {
	data, callback := utils.ExtractArgs(args)
	if data == nil || callback == nil {
		utils.SendResponse(callback, "error", "Invalid arguments")
		return
	}

	messageID, err := uuid.Parse(data["message_id"].(string))
	if err != nil {
		sentry.CaptureException(err)
		utils.LogError(callback, "invalid message_id format")
		return
	}

	starErr := adapter.UpdateMessageStarred(data["room_id"].(string), messageID, data["message_starred"].(bool))
	if starErr != nil {
		sentry.CaptureException(starErr)
		utils.LogError(callback, starErr.Error())
		return
	}
	utils.LogSuccess(callback, "Message starred boolean updated successfully")
}

// endregion

// region "updateMessageStarred" updates the message type and notifies clients.
func (adapter *socketAdapter) UpdateMessageStarred(roomId string, messageId uuid.UUID, messageStarred bool) error {
	// Update the message type in the database.
	if err := adapter.MessageService.UpdateMessageStarredById(messageId, messageStarred); err != nil {
		sentry.CaptureException(err)
		return err
	}

	notifyData := map[string]interface{}{
		"message_id":      messageId,
		"message_starred": messageStarred,
	}

	adapter.Gateway.EmitToRoomId("updated_message_starred", roomId, notifyData)
	return nil
}

// endregion

// region "handleReadMessage" processes marking messages as read.
func (adapter *socketAdapter) handleReadMessage(connectedUserID string, args ...any) {
	data, callback := utils.ExtractArgs(args)
	if data == nil || callback == nil {
		utils.SendResponse(callback, "error", "Invalid arguments")
		return
	}
	messageId, okMessage := data["message_id"].(string)

	if !okMessage || messageId == "" { // If no message ID is provided.
		err := adapter.ReadMessage(connectedUserID, data["room_id"].(string), nil) // Read the room without specific message.
		if err != nil {
			sentry.CaptureException(err)
			utils.LogError(callback, err.Error())
			return
		}
		utils.LogSuccess(callback, "Room read successfully without message ID")
		return
	}

	// Attempt to read the specific message.
	err := adapter.ReadMessage(connectedUserID, data["room_id"].(string), &messageId)
	if err != nil {
		sentry.CaptureException(err)
		utils.LogError(callback, err.Error())
		return
	}
	utils.LogSuccess(callback, "Message read successfully with message ID")
}

// endregion

// region "ReadMessage" marks message(s) as read in the database and notifies clients.
func (adapter *socketAdapter) ReadMessage(connectedUserID, roomId string, messageId *string) error {
	if err := adapter.MessageService.ReadMessageByRoomId(connectedUserID, roomId, messageId); err != nil {
		return err
	}

	notifyData := map[string]interface{}{
		"room_id": roomId,
	}

	adapter.Gateway.EmitToRoomId("read_message", roomId, notifyData)
	return nil
}

// endregion
