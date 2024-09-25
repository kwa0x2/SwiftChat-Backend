package gateway

import (
	"github.com/zishang520/socket.io/socket"
)

type SocketGateway interface {
	OnConnection(callback func(socketio *socket.Socket))
	EmitRoom(room, event string, data interface{})
	JoinRoom(socketio *socket.Socket, room string)
	Emit(event string, data interface{})
}

type socketGateway struct {
	server *socket.Server
}

func NewSocketGateway(server *socket.Server) SocketGateway {
	return &socketGateway{server: server}
}

func (g *socketGateway) OnConnection(callback func(socketio *socket.Socket)) {
	g.server.Of("/chat", nil).On("connection", func(clients ...any) {
		socketio := clients[0].(*socket.Socket)
		callback(socketio)
	})
}

func (g *socketGateway) EmitRoom(room, event string, data interface{}) {
	g.server.Of("/chat", nil).To(socket.Room(room)).Emit(event, data)
}

func (g *socketGateway) Emit(event string, data interface{}) {
	g.server.Of("/chat", nil).Emit(event, data)
}

func (g *socketGateway) JoinRoom(socketio *socket.Socket, room string) {
	socketio.Join(socket.Room(room))
}
