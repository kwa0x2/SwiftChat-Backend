package adapter

import (
	"fmt"
	"github.com/kwa0x2/swiftchat-backend/service"
	"github.com/kwa0x2/swiftchat-backend/socket/gateway"
	"github.com/zishang520/socket.io/socket"
	"sync"
)

type ISocketAdapter interface {
	HandleConnection()
	EmitToFriendsAndSentRequests(event, userEmail string, emitData interface{}) error
}

type socketAdapter struct {
	Gateway          gateway.ISocketGateway
	onlineUserEmails []string
	MessageService   service.IMessageService
	FriendService    service.IFriendService
	RequestService   service.IRequestService
	mux              sync.RWMutex
}

func NewSocketAdapter(gateway gateway.ISocketGateway, messageService service.IMessageService, friendService service.IFriendService, requestService service.IRequestService) ISocketAdapter {
	return &socketAdapter{
		Gateway:        gateway,
		MessageService: messageService,
		FriendService:  friendService,
		RequestService: requestService,
	}
}

// region "HandleConnection" manages user connections
func (adapter *socketAdapter) HandleConnection() {
	adapter.Gateway.OnConnection(func(socketio *socket.Socket) {
		ctx := socketio.Request().Context()
		connectedUserID := ctx.Value("id").(string)
		connectedUserMail := ctx.Value("email").(string)
		fmt.Println(connectedUserMail, " is  online")

		if !adapter.emailExists(connectedUserMail) {
			adapter.onlineUserEmails = append(adapter.onlineUserEmails, connectedUserMail)
		}

		adapter.Gateway.Emit("onlineUsers", adapter.onlineUserEmails) // Broadcast online users

		socketio.On("disconnect", func(...any) {
			adapter.handleDisconnect(connectedUserMail)
		})

		socketio.On("joinRoom", func(roomData ...any) {
			adapter.handleJoinRoom(socketio, roomData...)
		})

		socketio.On("sendMessage", func(args ...any) {
			adapter.handleSendMessage(connectedUserID, connectedUserMail, args...)
		})

		socketio.On("deleteMessage", func(args ...any) {
			adapter.handleDeleteMessage(args...)
		})

		socketio.On("editMessage", func(args ...any) {
			adapter.handleEditMessage(args...)
		})

		socketio.On("updateMessageType", func(args ...any) {
			adapter.handleUpdateMessageType(args...)
		})
		socketio.On("readMessage", func(args ...any) {
			adapter.handleReadMessage(connectedUserID, args...)
		})
	})
}

// endregion

// region "EmitToFriendsAndSentRequests" sends an event to all friends and sent requests of the specified user with the provided data.
func (adapter *socketAdapter) EmitToFriendsAndSentRequests(event, userEmail string, emitData interface{}) error {
	// Retrieve the list of friends for the given userEmail.
	friends, err := adapter.FriendService.GetFriends(userEmail, true)
	if err != nil {
		return err
	}

	requests, ReqErr := adapter.RequestService.GetSentRequests(userEmail)
	if ReqErr != nil {
		return ReqErr
	}

	emailSet := make(map[string]struct{})

	// Add friends' emails to the map.
	for _, friend := range friends {
		emailSet[friend.UserMail] = struct{}{}
	}

	// Add requests' sender emails to the map.
	for _, request := range requests {
		emailSet[request.ReceiverMail] = struct{}{}
	}

	// Prepare a WaitGroup to synchronize goroutines.
	var wg sync.WaitGroup
	for email := range emailSet {
		wg.Add(1) // Increment the WaitGroup counter.
		go func(email string) {
			defer wg.Done() // Decrement the counter when the goroutine completes.
			// Emit the event to the friend's notification room with the provided data.
			adapter.Gateway.EmitToNotificationRoom(event, email, emitData)
		}(email) // Pass the unique email to the goroutine.
	}
	wg.Wait() // Wait for all goroutines to finish.

	return nil // Return nil indicating success.
}

// endregion

// region "emailExists" checks if an email is already in the online user list
func (adapter *socketAdapter) emailExists(email string) bool {
	for _, existingEmail := range adapter.onlineUserEmails {
		if existingEmail == email {
			return true
		}
	}
	return false
}

// endregion
