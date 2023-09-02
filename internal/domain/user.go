package domain

import (
	"time"
)

type User struct {
	ID        string
	Nickname  string
	Wallet    string
	CreatedAt time.Time
}

type UserAuthReq struct {
	Wallet  string
	Sign    string
	Message string
}
