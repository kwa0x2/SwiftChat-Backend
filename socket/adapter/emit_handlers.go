package adapter

import (
	"github.com/zishang520/engine.io/utils"
	"sync"
)

func (adapter *SocketAdapter) EmitToNotificationRoom(notifyAction, receiverMail string, notifyObj any) {
	data := map[string]interface{}{
		"action": notifyAction,
		"data":   notifyObj,
	}

	utils.Log().Info("notify %+v\n mail:%s", data, receiverMail)

	adapter.Gateway.EmitRoom("notification", receiverMail, data)
}

func (adapter *SocketAdapter) EmitToRoomId(notifyAction, roomId string, notifyObj any) {
	data := map[string]interface{}{
		"action": notifyAction,
		"data":   notifyObj,
	}

	utils.Log().Info("notify %+v\n roomId:%s", data, roomId)

	adapter.Gateway.Emit(roomId, data)
}

func (adapter *SocketAdapter) EmitToFriends(event, userEmail string, emitData interface{}) error {
	friends, err := adapter.FriendService.GetFriends(userEmail, true)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	for _, friend := range friends {
		wg.Add(1)
		go func(friendEmail string) {
			defer wg.Done()
			adapter.EmitToNotificationRoom(event, friendEmail, emitData)
		}(friend.UserMail)
	}
	wg.Wait()

	return nil
}
