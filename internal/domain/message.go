package domain

import "time"

type MessageType string

const (
	MessageTypeText  MessageType = "text"
	MessageTypeImage MessageType = "image"
	MessageTypeFile  MessageType = "file"
)

type Message struct {
	Id          uint        `gorm:"primaryKey"`
	FromId      uint        `gorm:"not null;index:idx_messages_from_id"`
	PrivateId   *uint       `gorm:"index:idx_messages_private_id"`
	MessageType MessageType `gorm:"not null"`
	Content     string      `gorm:"not null"`
	Delivered   bool        `gorm:"not null;default:false"`
	Read        bool        `gorm:"not null;default:false"`
	CreatedAt   time.Time
	Version     int `gorm:"not null;default:1"`

	From    User     `gorm:"foreignKey:FromId;references:Id;constraint:OnDelete:CASCADE"`
	Private *Private `gorm:"foreignKey:PrivateId;references:Id;constraint:OnDelete:CASCADE"`
}
