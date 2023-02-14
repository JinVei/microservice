package domain

import (
	"context"

	"github.com/jinvei/microservice/app/userservice/domain/entity"
)

type IUserRepository interface {
	GetUserByEmail(ctx context.Context, email string) (entity.User, error)
	// SetUserSession()
	// GetUserSessions()
}
