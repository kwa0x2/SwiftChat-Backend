package types

type MessageType string

const (
	Text        MessageType = "text"
	StarredText MessageType = "starred_text"
	Photo       MessageType = "photo"
	File        MessageType = "file"
)
