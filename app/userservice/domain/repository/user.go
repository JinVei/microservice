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

func (u *UserRepository) GetUserByEmail(ctx context.Context, email string) (entity.Users, error) {
	var user entity.Users
	_, err := u.xorm.Context(ctx).Where("email = ?", email).Get(&user)
	return user, err
}

func (u *UserRepository) CreateUser(ctx context.Context, user *entity.Users) error {
	if _, err := u.xorm.InsertOne(user); err != nil {
		return err
	}
	return nil
}
