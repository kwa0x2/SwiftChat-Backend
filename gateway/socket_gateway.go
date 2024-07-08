package gateway

import (
	"github.com/zishang520/socket.io/socket"
)


type SocketGateway interface {
	OnConnection(callback func(socketio *socket.Socket))
	Emit(event, room string, data interface{})
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

func (g *socketGateway) Emit(event, room string, data interface{}){
	g.server.Of("/chat", nil).To(socket.Room(room)).Emit(event,data)
}
