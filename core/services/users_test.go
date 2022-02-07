package services

import (
	"github.com/fgiannotti/hubuc_coding_task/core/domain"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetUserSuccess(t *testing.T) {
	expectedUsr := domain.User{Name: "usr1", Email: "test@franco.com", EncryptedPwd: "encrypted"}
	usersRepo := &LocalUsersRepo{
		DB: map[string]domain.User{"usr1": expectedUsr},
	}
	usr, err := usersRepo.Get("usr1")
	assert.Equal(t, usr, expectedUsr)
	assert.Nil(t, err)
}


func TestGetUserNoMatchError(t *testing.T) {
	emptyMap := map[string]domain.User{}
	usersRepo := &LocalUsersRepo{
		DB: emptyMap,
	}
	_, err := usersRepo.Get("randommmm")
	assert.NotNil(t, err)
}
