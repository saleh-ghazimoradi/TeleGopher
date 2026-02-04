package domain

import "time"

type MessageType string

const (
	MessageTypeText  MessageType = "text"
	MessageTypeImage MessageType = "image"
	MessageTypeFile  MessageType = "file"
)

type Message struct {
	Id          int64
	FromId      int64
	PrivateId   int64
	MessageType MessageType
	Content     string
	Delivered   bool
	Read        bool
	CreatedAt   time.Time
	Version     int
}
