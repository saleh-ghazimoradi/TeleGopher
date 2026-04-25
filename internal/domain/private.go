package domain

import "time"

type Private struct {
	Id        uint `gorm:"primaryKey"`
	User1Id   uint `gorm:"not null;index:idx_privates_user1_id"`
	User2Id   uint `gorm:"not null;index:idx_privates_user2_id"`
	CreatedAt time.Time
	Version   int `gorm:"not null;default:1"`

	User1 User `gorm:"foreignKey:User1Id;references:Id;constraint:OnDelete:CASCADE"`
	User2 User `gorm:"foreignKey:User2Id;references:Id;constraint:OnDelete:CASCADE"`
}
