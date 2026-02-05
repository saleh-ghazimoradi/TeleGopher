package dto

import (
	"time"
)

type PrivateResponse struct {
	Id        int64     `json:"id"`
	User1Id   int64     `json:"user1_id"`
	User2Id   int64     `json:"user2_id"`
	CreatedAt time.Time `json:"created_at"`
}
