package services

import (
	"errors"
	"fmt"
	"github.com/fgiannotti/hubuc_coding_task/core/domain"
)

var UserNotFoundError = func(username string) error { return errors.New(fmt.Sprintf("User %s not found", username)) }

type UsersRepo interface {
	Save(user domain.User) error
	Get(username string) (domain.User, error)
}

type LocalUsersRepo struct {
	DB map[string]domain.User
}

func NewLocalUsersRepo() UsersRepo {
	return &LocalUsersRepo{DB: map[string]domain.User{}}
}

func (l *LocalUsersRepo) Save(user domain.User) error {
	panic("implement me")
}

func (l *LocalUsersRepo) Get(username string) (domain.User, error) {
	user, ok := l.DB[username]
	if !ok {
		return domain.User{}, UserNotFoundError(username)
	}

	return user, nil
}
