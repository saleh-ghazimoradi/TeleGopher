package dto

import "time"

type PrivateRequest struct {
	ReceiverId uint `json:"receiver_id"`
}

type PrivateResponse struct {
	Id        uint      `json:"id"`
	User1Id   uint      `json:"user1_id"`
	User2Id   uint      `json:"user2_id"`
	CreatedAt time.Time `json:"created_at"`
}
