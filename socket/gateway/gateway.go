package gateway

import (
	"github.com/zishang520/socket.io/socket"
)

type ISocketGateway interface {
	OnConnection(callback func(socketio *socket.Socket))
	EmitRoom(room, event string, data interface{})
	JoinRoom(socketio *socket.Socket, room string)
	Emit(event string, data interface{})
	EmitToNotificationRoom(notifyAction, receiverMail string, notifyObj any)
	EmitToRoomId(notifyAction, roomId string, notifyObj any)
}

type socketGateway struct {
	Server    *socket.Server
	namespace string
}

func NewSocketGateway(server *socket.Server, namespace string) ISocketGateway {
	return &socketGateway{
		Server:    server,
		namespace: namespace,
	}
}

// region "OnConnection" sets up a listener for new socket connections.
func (g *socketGateway) OnConnection(callback func(socketio *socket.Socket)) {
	// When a new connection is established, the provided callback is invoked.
	g.Server.Of(g.namespace, nil).On("connection", func(clients ...any) {
		socketio := clients[0].(*socket.Socket) // Extract the socket from the clients.
		callback(socketio)                      // Invoke the callback with the connected socket.
	})
}

// endregion

// region "EmitRoom" sends an event with data to all sockets in a specific room.
func (g *socketGateway) EmitRoom(room, event string, data interface{}) {
	g.Server.Of(g.namespace, nil).To(socket.Room(room)).Emit(event, data)
}

// endregion

// region "Emit" sends an event with data to all connected sockets.
func (g *socketGateway) Emit(event string, data interface{}) {
	g.Server.Of(g.namespace, nil).Emit(event, data)
}

// endregion

// region "JoinRoom" adds a socket to a specified room.
func (g *socketGateway) JoinRoom(socketio *socket.Socket, room string) {
	socketio.Join(socket.Room(room)) // Add the socket to the specified room.
}

// endregion

// region "EmitToNotificationRoom" sends a notification action with data to a specific user's notification room.
func (g *socketGateway) EmitToNotificationRoom(notifyAction, receiverMail string, notifyObj any) {
	data := map[string]interface{}{
		"action": notifyAction,
		"data":   notifyObj,
	}

	g.EmitRoom("notification", receiverMail, data)
}

// endregion

// region "EmitToRoomId" sends a notification action with data to a specific room by room ID.
func (g *socketGateway) EmitToRoomId(notifyAction, roomId string, notifyObj any) {
	data := map[string]interface{}{
		"action": notifyAction,
		"data":   notifyObj,
	}

	g.Emit(roomId, data)
}

// endregion
