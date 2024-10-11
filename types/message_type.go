package types

type MessageType string

const (
	Text  MessageType = "text"
	Photo MessageType = "photo"
	File  MessageType = "file"
)
