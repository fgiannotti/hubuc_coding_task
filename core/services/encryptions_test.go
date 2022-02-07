package services

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBcryptGeneratesHashedPwd(t *testing.T) {
	pwd := "pwd"
	encrypted, err := NewBcryptEncryptionsService().Encrypt(pwd)

	assert.Nil(t, err)
	assert.NotEmpty(t, encrypted)
}


func TestBcryptCompareMatchesPwd(t *testing.T) {
	pwd := "random"
	match, err := NewBcryptEncryptionsService().Compare("$2a$04$VJiidviKaccvQlLaJWdcveISlwP3Ze3D6btRBF02KrhPSAAAA", pwd)

	assert.NotNil(t, err)
	assert.False(t, match)
}


func TestBcryptCompareNotMatchesPwd(t *testing.T) {

	match, err := NewBcryptEncryptionsService().Compare("$2a$04$VJiidviKaccvQlLaJWdcveISlwP3Ze3D6btRBF02KrhPSJ8kJE5Gm", "pwd")

	assert.Equal(t, err, nil)
	assert.True(t, match)
}
