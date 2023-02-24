package entity

type Session struct {
	UserID     string `json:"uid"`
	SessionId  string `json:"sid"`
	LastUpdate string `json:"last_update"`
	ExpireAt   string `json:"expire_at"`
}
