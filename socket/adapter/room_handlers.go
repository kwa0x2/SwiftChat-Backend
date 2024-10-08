package adapter

import (
	"github.com/zishang520/engine.io/utils"
	"github.com/zishang520/socket.io/socket"
)

// region "handleJoinRoom" handles the event when a socket joins a specific room.
func (adapter *socketAdapter) handleJoinRoom(socketio *socket.Socket, roomData ...any) {
	// Attempt to retrieve the room ID from the provided roomData.
	roomId, ok := roomData[0].(string)
	if !ok {
		utils.Log().Error(`socket message type error socketid: %s `, socketio.Id())
		return
	}
	adapter.Gateway.JoinRoom(socketio, roomId)
}

// endregion
