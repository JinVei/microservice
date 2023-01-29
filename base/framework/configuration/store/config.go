package store

type config struct {
	Addr     []string `json:"addrs,omitempty"`
	Username string   `json:"username,omitempty"`
	Password string   `json:"password,omitempty"`
}
