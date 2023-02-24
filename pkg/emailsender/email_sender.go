package emailsender

import "context"

type Config struct {
	SID        string `json:"secretID"`
	SK         string `json:"secretKey"`
	FromEmail  string `json:"fromEmail"`
	ReplyTo    string `json:"replyTo"`
	TemplateID uint64 `json:"templateID"`
	Subject    string `json:"subject"`
}

type EmailSender interface {
	Send(ctx context.Context, email, parameter string) error
}

func New(conf Config) (EmailSender, error) {
	s, err := NewTEmailSender(conf)
	return s, err
}
