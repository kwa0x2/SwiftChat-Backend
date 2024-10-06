package adapter

//func (adapter *SocketAdapter) handleUpdateUsername(connectedUserMail string, args ...any) {
//	data, callback := utils.ExtractArgs(args)
//	if data == nil || callback == nil {
//		utils.SendResponse(callback, "error", "Invalid arguments")
//		return
//	}
//
//	err := adapter.UpdateUsername(data["user_name"].(string), connectedUserMail)
//	if err != nil {
//		utils.LogError(callback, err.Error())
//		return
//	}
//	utils.LogSuccess(callback, "Username updated successfully")
//}
//
//func (adapter *SocketAdapter) UpdateUsername(userName, connectedUserMail string) error {
//	if err := adapter.userService.UpdateUsernameByMail(userName, connectedUserMail); err != nil {
//		return err
//	}
//
//	notifyData := map[string]interface{}{
//		"updated_username": userName,
//		"user_email":       connectedUserMail,
//	}
//
//	friends, err := adapter.friendService.GetFriends(connectedUserMail, true)
//	if err != nil {
//		return err
//	}
//
//	var wg sync.WaitGroup
//	for _, friend := range friends {
//		wg.Add(1)
//		go func(friendEmail string) {
//			defer wg.Done()
//			fmt.Println("friendmail", friendEmail)
//			adapter.EmitToNotificationRoom("update_username", friendEmail, notifyData)
//		}(friend.UserMail)
//	}
//	wg.Wait()
//
//	return nil
//}

//func (adapter *SocketAdapter) handleUpdateUserPhoto(connectedUserMail string, args ...any) {
//	data, callback := utils.ExtractArgs(args)
//	if data == nil || callback == nil {
//		utils.SendResponse(callback, "error", "Invalid arguments")
//		return
//	}
//
//	err := adapter.UpdateUsername(data["user_name"].(string), connectedUserMail)
//	if err != nil {
//		utils.LogError(callback, err.Error())
//		return
//	}
//	utils.LogSuccess(callback, "Username updated successfully")
//}
//
//func (adapter *SocketAdapter) UpdateUserPhoto(userPhoto, connectedUserMail string) error {
//
//	notifyData := map[string]interface{}{
//		"updated_user_photo": userName,
//		"user_email":         connectedUserMail,
//	}
//
//	friends, err := adapter.friendService.GetFriends(connectedUserMail, true)
//	if err != nil {
//		return err
//	}
//
//	var wg sync.WaitGroup
//	for _, friend := range friends {
//		wg.Add(1)
//		go func(friendEmail string) {
//			defer wg.Done()
//			fmt.Println("friendmail", friendEmail)
//			adapter.EmitToNotificationRoom("update_user_photo", friendEmail, notifyData)
//		}(friend.UserMail)
//	}
//	wg.Wait()
//
//	return nil
//}
