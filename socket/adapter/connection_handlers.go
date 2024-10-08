package adapter

import (
	"fmt"
)

// region "handleDisconnect" removes a user's email from the online users list when they disconnect.
func (adapter *socketAdapter) handleDisconnect(email string) {
	// Iterate through the list of online user emails.
	for i, existingEmail := range adapter.onlineUserEmails {
		// Check if the current email matches the disconnected user's email.
		if existingEmail == email {
			// Remove the email from the online users list.
			adapter.onlineUserEmails = append(adapter.onlineUserEmails[:i], adapter.onlineUserEmails[i+1:]...)
			fmt.Println(email, " is offline") // Log the user's disconnection.

			// Emit the updated online users list to notify all clients.
			adapter.Gateway.Emit("onlineUsers", adapter.onlineUserEmails)

			break // Exit the loop after removing the email.
		}
	}

	// Emit the updated online users list again (redundant here, can be removed).
	adapter.Gateway.Emit("onlineUsers", adapter.onlineUserEmails)
}

// endregion
