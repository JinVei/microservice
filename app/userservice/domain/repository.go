package domain

import (
	"context"

	"github.com/jinvei/microservice/app/userservice/domain/entity"
)

type IUserRepository interface {
	GetUserByEmail(ctx context.Context, email string) (entity.Users, error)
	CreateUser(ctx context.Context, user *entity.Users) error
}
