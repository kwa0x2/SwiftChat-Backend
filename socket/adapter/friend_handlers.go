package adapter

//func (adapter *SocketAdapter) handleSendFriend(connectedUserMail string, args ...any) {
//	data, callback := utils.ExtractArgs(args)
//	if data == nil || callback == nil {
//		return
//	}
//
//	receiverEmail, ok := data["receiver_email"].(string)
//	if !ok || receiverEmail == connectedUserMail {
//		utils.LogError(callback, "socket message type error")
//		return
//	}
//
//	requestObj := models.Request{
//		SenderMail:   connectedUserMail,
//		ReceiverMail: receiverEmail,
//	}
//
//	status := adapter.SendFriend(&requestObj)
//	utils.SendResponse(callback, status, "")
//}
//
//func (adapter *SocketAdapter) SendFriend(request *models.Request) string {
//	var pgErr *pgconn.PgError
//	existingFriend, err := adapter.friendService.GetFriend(request.SenderMail, request.ReceiverMail)
//	if err != nil {
//		if !errors.Is(err, gorm.ErrRecordNotFound) {
//			return "error"
//
//		}
//	}
//	if existingFriend != nil && existingFriend.FriendStatus == "friend" {
//		return "already_friend"
//	}
//
//	if isEmailExists := adapter.userService.IsEmailExists(request.ReceiverMail); !isEmailExists {
//		if err := adapter.requestService.Insert(nil, request); err != nil {
//			if errors.As(err, &pgErr) && pgErr.Code == "23505" {
//				return "duplicate"
//			}
//			return "error"
//		}
//
//		_, err := adapter.resendService.SendMail(request.ReceiverMail, "You have received a new friend request from the SwiftChat app!", "friend_request")
//		if err != nil {
//			return "error"
//		}
//
//		return "email_sent"
//	}
//
//	data, err := adapter.requestService.InsertAndReturnUser(request)
//	if err != nil {
//		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
//			return "duplicate"
//		}
//
//		return "error"
//	}
//
//	adapter.EmitToNotificationRoom("friend_request", request.ReceiverMail, data)
//	return "friend_sent"
//}

//func (adapter *SocketAdapter) handleBlockFriend(connectedUserMail string, args ...any) {
//	data, callback := utils.ExtractArgs(args)
//	if data == nil || callback == nil {
//		utils.SendResponse(callback, "error", "Invalid arguments")
//		return
//	}
//
//	friendObj := models.Friend{
//		UserMail:  data["friend_mail"].(string),
//		UserMail2: connectedUserMail,
//	}
//
//	responseData := adapter.BlockFriend(&friendObj, data["friend_name"].(string))
//	utils.SendResponse(callback, responseData["status"].(string), "friend blocked")
//}
//
//func (adapter *SocketAdapter) BlockFriend(friend *models.Friend, friendName string) map[string]interface{} {
//	friendObj := map[string]interface{}{
//		"status":        "",
//		"friend_status": "",
//	}
//
//	friendStatus, err := adapter.friendService.Block(friend)
//	if err != nil {
//		friendObj["status"] = "error"
//		return friendObj
//	}
//
//	data := map[string]interface{}{
//		"friend_name":   friendName,
//		"friend_mail":   friend.UserMail2,
//		"friend_status": friendStatus,
//	}
//
//	adapter.EmitToNotificationRoom("blocked_friend", friend.UserMail, data)
//	friendObj["status"] = "success"
//	friendObj["friend_status"] = friendStatus
//
//	return friendObj
//}

//func (adapter *SocketAdapter) handleDeleteFriend(connectedUserMail string, args ...any) {
//	data, callback := utils.ExtractArgs(args)
//	if data == nil || callback == nil {
//		utils.SendResponse(callback, "error", "Invalid arguments")
//		return
//	}
//
//	friendObj := models.Friend{
//		UserMail:  data["user_mail"].(string),
//		UserMail2: connectedUserMail,
//	}
//
//	err := adapter.DeleteFriend(&friendObj)
//	if err != nil {
//		utils.LogError(callback, err.Error())
//		return
//	}
//	utils.LogSuccess(callback, "Message edited successfully")
//}
//
//func (adapter *SocketAdapter) DeleteFriend(friend *models.Friend) error { // hem unblock hemde arkadaslikta cikarmada kullaniliyor
//	if err := adapter.friendService.Delete(friend); err != nil {
//		return err
//	}
//
//	data := map[string]interface{}{
//		"user_email": friend.UserMail2,
//	}
//
//	adapter.EmitToNotificationRoom("deleted_friend", friend.UserMail, data)
//	return nil
//}

//func (adapter *SocketAdapter) handleUpdateFriendshipRequest(connectedUserMail string, args ...any) {
//	data, callback := utils.ExtractArgs(args)
//	if data == nil || callback == nil {
//		utils.SendResponse(callback, "error", "Invalid arguments")
//		return
//	}
//
//	requestObj := models.Request{
//		SenderMail:    data["sender_mail"].(string),
//		ReceiverMail:  connectedUserMail,
//		RequestStatus: types.RequestStatus(data["status"].(string)),
//	}
//
//	err := adapter.UpdateFriendshipRequest(&requestObj)
//	if err != nil {
//		utils.LogError(callback, err.Error())
//		return
//	}
//	utils.LogSuccess(callback, "Message edited successfully")
//} // gelen arkadaslik istekleri islemlerinde kullaniliyor
//
//func (adapter *SocketAdapter) UpdateFriendshipRequest(request *models.Request) error {
//	data, err := adapter.requestService.UpdateFriendshipRequest(request)
//	if err != nil {
//		return err
//	}
//
//	adapter.EmitToNotificationRoom("update_friendship_request", request.SenderMail, data)
//	return nil
//}
