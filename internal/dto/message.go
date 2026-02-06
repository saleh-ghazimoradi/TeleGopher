package dto

import (
	"github.com/saleh-ghazimoradi/TeleGopher/internal/domain"
	"github.com/saleh-ghazimoradi/TeleGopher/internal/helper"
	"time"
)

type MessageResponse struct {
	Id          int64     `json:"id"`
	FromId      int64     `json:"from_id"`
	PrivateId   int64     `json:"private_id"`
	MessageType string    `json:"message_type"`
	Content     string    `json:"content"`
	Delivered   bool      `json:"delivered"`
	Read        bool      `json:"read"`
	CreatedAt   time.Time `json:"created_at"`
}

type CreateMessageRequest struct {
	PrivateId   int64  `json:"private_id"`
	MessageType string `json:"message_type"`
	Content     string `json:"content"`
}

type MessagesListResponse struct {
	Messages    []MessageResponse `json:"messages"`
	Page        int               `json:"page"`
	Limit       int               `json:"limit"`
	HasNextPage bool              `json:"has_next_page"`
}

func (r *CreateMessageRequest) Validate(v *helper.Validator) {
	v.Check(r.PrivateId > 0, "private_id", "must be a valid private chat ID")
	v.Check(r.MessageType != "", "message_type", "must be provided")
	v.Check(isValidMessageType(r.MessageType), "message_type", "must be a valid message type (text, image, file)")
	v.Check(r.Content != "", "content", "must be provided")
	v.Check(len(r.Content) <= 5000, "content", "must not exceed 5000 characters")
}

func isValidMessageType(messageType string) bool {
	switch domain.MessageType(messageType) {
	case domain.MessageTypeText, domain.MessageTypeImage, domain.MessageTypeFile:
		return true
	default:
		return false
	}
}
