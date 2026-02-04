package domain

import "time"

type Private struct {
	Id        int64
	User1Id   int64
	User2Id   int64
	CreatedAt time.Time
	Version   int
}
