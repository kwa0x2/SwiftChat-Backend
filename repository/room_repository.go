package repository

import (
	"github.com/google/uuid"
	"github.com/kwa0x2/realtime-chat-backend/models"
	"github.com/kwa0x2/realtime-chat-backend/types"
	"gorm.io/gorm"
	"time"
)

type RoomRepository struct {
	DB *gorm.DB
}

// region CREATE ROOM REPOSITORY
func (r *RoomRepository) CreateRoom(tx *gorm.DB, room *models.Room) (*models.Room, error) {
	db := r.DB
	if tx != nil {
		db = tx
	}
	if err := db.Create(&room).Error; err != nil {
		return nil, err
	}
	return room, nil
}

//endregion

func (r *RoomRepository) Update(tx *gorm.DB, room *models.Room) error {
	db := r.DB
	if tx != nil {
		db = tx
	}
	if err := db.Model(&models.Room{}).Where("room_id = ?", room.RoomID).Updates(&models.Room{LastMessage: room.LastMessage, LastMessageID: room.LastMessageID}).Error; err != nil {
		return err
	}
	return nil
}

type ChatList struct {
	RoomID        uuid.UUID          `json:"room_id"`
	LastMessage   string             `json:"last_message"`
	UpdatedAt     time.Time          `json:"updatedAt" gorm:"column:updatedAt"`
	UserName      string             `json:"user_name"`
	UserPhoto     string             `json:"user_photo"`
	UserEmail     string             `json:"user_email"`
	FriendStatus  types.FriendStatus `json:"friend_status"`
	CreatedAt     time.Time          `json:"createdAt" gorm:"column:createdAt"`
	LastMessageID uuid.UUID          `json:"last_message_id" gorm:"type:uuid"`
	DeletedAt     gorm.DeletedAt     `json:"message_deleted_at"`
}

func (r *RoomRepository) GetChatList(userId, userMail string) ([]*ChatList, error) {
	var chatLists []*ChatList

	if err := r.DB.Model(&models.Room{}).Debug().
		Select(`DISTINCT ON ("ROOM".room_id) "ROOM".room_id, "ROOM".last_message_id,"ROOM".last_message_id, "ROOM"."updatedAt", "USER".user_name, "USER".user_photo,"USER"."createdAt", "USER".user_email, "FRIEND".friend_status, "MESSAGE".message AS last_message, "MESSAGE"."deletedAt" AS message_deleted_at`).
		Joins(`INNER JOIN "USER_ROOM" ON "ROOM".room_id = "USER_ROOM".room_id`).
		Joins(`LEFT JOIN "USER_ROOM" ur2 ON "ROOM".room_id = ur2.room_id AND ur2.user_id != ?`, userId).
		Joins(`LEFT JOIN "USER" ON ur2.user_id = "USER".user_id`).
		Joins(`LEFT JOIN "FRIEND" ON (("USER".user_email = "FRIEND".user_mail AND ? = "FRIEND".user_mail2) OR ("USER".user_email = "FRIEND".user_mail2 AND ? = "FRIEND".user_mail))`, userMail, userMail).
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
