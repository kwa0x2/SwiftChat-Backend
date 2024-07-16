package models

type FriendshipStatus string

const (
	Pending   FriendshipStatus = "pending"
	Accepted  FriendshipStatus = "accepted"
	Rejected  FriendshipStatus = "rejected"
	EmailSent FriendshipStatus = "email_sent"
)

type Friendship struct {
	SenderId         string           `json:"sender_id" gorm:"primaryKey;not null"`
	ReceiverId       string           `json:"receiver_id" gorm:"primaryKey;not null"`
	FriendshipStatus FriendshipStatus `json:"friendship_status" gorm:"type:friendship_status;not null;default:pending"`
	User             User             `json:"user" gorm:"foreignkey:SenderId;association_foreignkey:UserID"`
}
