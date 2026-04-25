package dto

import (
	"time"
)

type MessageRequest struct {
	PrivateId   uint   `json:"private_id"`
	MessageType string `json:"message_type"`
	Content     string `json:"content"`
}

type MessageResponse struct {
	Id          uint      `json:"id"`
	FromId      uint      `json:"from_id"`
	PrivateId   uint      `json:"private_id"`
	MessageType string    `json:"message_type"`
	Content     string    `json:"content"`
	Delivered   bool      `json:"delivered"`
	Read        bool      `json:"read"`
	CreatedAt   time.Time `json:"created_at"`
}

type MessageListResponse struct {
	Messages    []MessageResponse `json:"messages"`
	Page        int               `json:"page"`
	Limit       int               `json:"limit"`
	HasNextPage bool              `json:"has_next_page"`
}
