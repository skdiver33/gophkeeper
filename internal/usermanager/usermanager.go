package usermanager

import (
	"context"
	"errors"
	"fmt"

	"github.com/skdiver33/gophkeeper/model"
)

type UserStorageInterface interface {
	AddUser(ctx context.Context, user *model.User) (int, error)
	GetUser(ctx context.Context, login string, password string) (*model.User, error)
}

type UManager struct {
	Storage UserStorageInterface
}

var (
	ErrInternal             = errors.New("internal error")
	ErrUserAlreadyExist     = errors.New("error. User already exist")
	ErrUserWithCredNotFound = errors.New("error. User with login and password not found")
)

func (um *UManager) UserRegister(ctx context.Context, data model.User) (*model.User, error) {

	//find user in storage if exist return error if not found add user to storage
	user, err := um.Storage.GetUser(ctx, data.Login, data.Password)
	if err != nil {
		return nil, fmt.Errorf("error get user from storage %w,%w", ErrInternal, err)
	}
	if user != nil {
		return nil, ErrUserAlreadyExist
	}
	newUser := &model.User{Login: data.Login, Password: data.Password}

	id, err := um.Storage.AddUser(ctx, newUser)
	if err != nil {
		return nil, fmt.Errorf("error add user to storage %w,%w", ErrInternal, err)
	}
	newUser.ID = id
	return newUser, nil
}

func (um *UManager) UserAuth(context.Context, *model.User) (string, error) {

	return "", nil
}
