package services

import (
	"errors"
	"fmt"
	"github.com/fgiannotti/hubuc_coding_task/core/domain"
	"go.uber.org/zap"
)

var UserNotFoundError = errors.New("User not found")

type UsersRepo interface {
	Save(user domain.User) error
	Get(username string) (domain.User, error)
}

type LocalUsersRepo struct {
	DB     map[string]domain.User
	logger *zap.SugaredLogger
}

func NewLocalUsersRepo(logger *zap.SugaredLogger) UsersRepo {
	return &LocalUsersRepo{DB: map[string]domain.User{}, logger: logger}
}

// Save - Always overwrites entry
func (l *LocalUsersRepo) Save(user domain.User) error {
	l.DB[user.Name] = user
	return nil
}

func (l *LocalUsersRepo) Get(username string) (domain.User, error) {
	user, ok := l.DB[username]
	if !ok {
		l.logger.Info("User not found in map")
		return domain.User{}, fmt.Errorf("%w: username %s", UserNotFoundError, username)
	}

	return user, nil
}
