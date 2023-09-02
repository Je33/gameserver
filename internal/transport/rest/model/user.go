package model

import "time"

type User struct {
	ID        string    `json:"id"`
	Nickname  string    `json:"nickname"`
	Wallet    string    `json:"wallet"`
	CreatedAt time.Time `json:"createdAt"`
}

type UserAuthReq struct {
	Wallet  string `json:"wallet"`
	Sign    string `json:"sign"`
	Message string `json:"message"`
}

type UserRefreshReq struct {
	RefreshToken string `json:"refresh_token"`
}

type UserAuthRes struct {
	AuthToken    string `json:"auth_token"`
	RefreshToken string `json:"refresh_token"`
}
