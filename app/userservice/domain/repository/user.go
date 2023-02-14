package repository

import (
	"context"

	"github.com/jinvei/microservice/app/userservice/domain"
	"github.com/jinvei/microservice/app/userservice/domain/entity"
	"xorm.io/xorm"
)

type UserRepository struct {
	xorm *xorm.Engine
}

func NewUserRepository(xorm *xorm.Engine) domain.IUserRepository {
	return &UserRepository{
		xorm: xorm,
	}
}

func (u *UserRepository) GetUserByEmail(ctx context.Context, email string) (entity.User, error) {
	var user entity.User
	_, err := u.xorm.Context(ctx).Where("email = ?", email).Get(&user)
	return user, err
}
