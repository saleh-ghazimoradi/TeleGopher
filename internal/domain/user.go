package domain

import "time"

type User struct {
	Id                   int64
	Name                 string
	Email                string
	Password             string
	RefreshTokenWeb      *string
	RefreshTokenWebAt    *time.Time
	RefreshTokenMobile   *string
	RefreshTokenMobileAt *time.Time
	CreatedAt            time.Time
	Version              int
}

func (u *User) ToMap() map[string]any {
	return map[string]any{
		"id":    u.Id,
		"name":  u.Name,
		"email": u.Email,
	}
}
