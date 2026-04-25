package domain

import (
	"time"
)

type User struct {
	Id                   uint   `gorm:"primaryKey"`
	Name                 string `gorm:"not null"`
	Email                string `gorm:"uniqueIndex;not null"`
	Password             string `gorm:"not null"`
	RefreshTokenWeb      *string
	RefreshTokenWebAt    *time.Time
	RefreshTokenMobile   *string
	RefreshTokenMobileAt *time.Time
	Version              uint `gorm:"default:1;not null"`
	CreatedAt            time.Time
	UpdatedAt            time.Time
}

func (u *User) ToMap() map[string]any {
	return map[string]any{
		"id":    u.Id,
		"name":  u.Name,
		"email": u.Email,
	}
}
