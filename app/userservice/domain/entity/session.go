package entity

type Session struct {
	UserID    string `json:"uid"`
	SessionId string `json:"sid"`
	// TODO: customer session data
	LastUpdate string `json:"last_update"`
	ExpireAt   string `json:"expire_at"`
}
