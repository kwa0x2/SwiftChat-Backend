package adapter

import "github.com/zishang520/engine.io/utils"

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
