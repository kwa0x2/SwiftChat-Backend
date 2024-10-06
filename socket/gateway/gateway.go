package gateway

import (
	"github.com/zishang520/socket.io/socket"
)

type ISocketGateway interface {
	OnConnection(callback func(socketio *socket.Socket))
	EmitRoom(room, event string, data interface{})
	JoinRoom(socketio *socket.Socket, room string)
	Emit(event string, data interface{})
}

type SocketGateway struct {
	Server *socket.Server
}

func (g *SocketGateway) OnConnection(callback func(socketio *socket.Socket)) {
	g.Server.Of("/chat", nil).On("connection", func(clients ...any) {
		socketio := clients[0].(*socket.Socket)
		callback(socketio)
	})
}

func (g *SocketGateway) OnDisconnect(callback func(socketio *socket.Socket)) {
	g.Server.Of("/chat", nil).On("disconnect", func(clients ...any) {
		socketio := clients[0].(*socket.Socket)
		callback(socketio)
	})
}

func (g *SocketGateway) EmitRoom(room, event string, data interface{}) {
	g.Server.Of("/chat", nil).To(socket.Room(room)).Emit(event, data)
}

func (g *SocketGateway) Emit(event string, data interface{}) {
	g.Server.Of("/chat", nil).Emit(event, data)
}

func (g *SocketGateway) JoinRoom(socketio *socket.Socket, room string) {
	socketio.Join(socket.Room(room))
}
