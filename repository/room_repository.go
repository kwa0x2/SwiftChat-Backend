package repository

import (
	"github.com/google/uuid"
	"github.com/kwa0x2/swiftchat-backend/models"
	"github.com/kwa0x2/swiftchat-backend/types"
	"gorm.io/gorm"
	"time"
)

type IRoomRepository interface {
	Create(tx *gorm.DB, room *models.Room) (*models.Room, error)
	Update(tx *gorm.DB, whereRoom *models.Room, updateRoom *models.Room) error
	GetChatList(userId, userEmail string) ([]*ChatList, error)
	GetDB() *gorm.DB
}
type roomRepository struct {
	DB *gorm.DB
}

func NewRoomRepository(db *gorm.DB) IRoomRepository {
	return &roomRepository{
		DB: db,
	}
}

// region "Create" adds a new room to the database
func (r *roomRepository) Create(tx *gorm.DB, room *models.Room) (*models.Room, error) {
	db := r.DB
	if tx != nil {
		db = tx // Use the provided transaction if available
	}
	if err := db.Create(&room).Error; err != nil {
		return nil, err
	}
	return room, nil
}

//endregion

// region "Update" modifies the fields of a room in the database based on specified conditions
func (r *roomRepository) Update(tx *gorm.DB, whereRoom *models.Room, updateRoom *models.Room) error {
	db := r.DB
	if tx != nil {
		db = tx // Use the provided transaction if available
	}
	if err := db.Model(&models.Room{}).Where(whereRoom).Updates(updateRoom).Error; err != nil {
		return err
	}
	return nil
}

// endregion

// region "GetChatList" DTO
type ChatList struct {
	RoomID           uuid.UUID          `json:"room_id"`                           // Unique identifier for the room
	LastMessage      string             `json:"last_message"`                      // Last message in the chat
	UpdatedAt        time.Time          `json:"updatedAt" gorm:"column:updatedAt"` // Last update timestamp
	UserName         string             `json:"user_name"`                         // Name of the user associated with the room
	UserPhoto        string             `json:"user_photo"`                        // Photo of the user associated with the room
	UserEmail        string             `json:"user_email"`                        // Email of the user associated with the room
	FriendStatus     types.FriendStatus `json:"friend_status"`                     // Status of the friendship with the user
	CreatedAt        time.Time          `json:"createdAt" gorm:"column:createdAt"` // Room creation timestamp
	LastMessageID    uuid.UUID          `json:"last_message_id" gorm:"type:uuid"`  // Identifier for the last message
	MessageDeletedAt gorm.DeletedAt     `json:"message_deleted_at"`                // Timestamp when the last message was deleted
	MessageType      types.MessageType  `json:"message_type" gorm:"type:message_type;not null"`
}

// endregion

// region "GetChatList" retrieves the list of chat rooms for a user, including last message details
func (r *roomRepository) GetChatList(userId, userEmail string) ([]*ChatList, error) {
	// GetChatList fetches all chat rooms associated with a user, including the last message details.
	// It returns a slice of ChatList and an error if the retrieval fails.

	var chatLists []*ChatList

	if err := r.DB.Model(&models.Room{}).Debug().
		Select(`DISTINCT ON ("ROOM".room_id) "ROOM".room_id, "ROOM".last_message_id, "ROOM"."updatedAt", "USER".user_name, "USER".user_photo,"USER"."createdAt", "USER".user_email, "FRIEND".friend_status, "MESSAGE".message AS last_message,"MESSAGE".message_type,  "MESSAGE"."deletedAt" AS message_deleted_at`).
		Joins(`INNER JOIN "USER_ROOM" ON "ROOM".room_id = "USER_ROOM".room_id`).
		Joins(`LEFT JOIN "USER_ROOM" ur2 ON "ROOM".room_id = ur2.room_id AND ur2.user_id != ?`, userId).
		Joins(`LEFT JOIN "USER" ON ur2.user_id = "USER".user_id`).
		Joins(`LEFT JOIN "FRIEND" ON (("USER".user_email = "FRIEND".user_mail AND ? = "FRIEND".user_mail2) OR ("USER".user_email = "FRIEND".user_mail2 AND ? = "FRIEND".user_mail))`, userEmail, userEmail).
		Joins(`LEFT JOIN "MESSAGE" ON "ROOM".last_message_id = "MESSAGE".message_id`).
		Where(`"USER_ROOM".user_id = ?`, userId).
		Where(`"MESSAGE".room_id IS NOT NULL`).
		Where(`"ROOM"."deletedAt" IS NULL`).
		Order(`"ROOM".room_id, "ROOM"."updatedAt" DESC`).
		Scan(&chatLists).Error; err != nil {
		return nil, err
	}

	return chatLists, nil
}

// endregion

// region "GetDB" returns the underlying gorm.DB instance
func (r *roomRepository) GetDB() *gorm.DB {
	return r.DB // Return the database instance
}

// endregion
