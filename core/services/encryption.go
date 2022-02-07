package services

import "golang.org/x/crypto/bcrypt"

type Encryptions interface {
	Encrypt(pwd string) (string, error)
	Compare(encryptedPwd string, pwd string) (bool,error)
}
type BcryptEncryptionsService struct {
	cost int
}

func (es *BcryptEncryptionsService) Compare(encryptedPwd string, pwd string) (bool, error) {
	panic("implement me")
}

func NewBcryptEncryptionsService() Encryptions {
	return &BcryptEncryptionsService{cost: bcrypt.DefaultCost}
}

func (es *BcryptEncryptionsService) Encrypt(pwd string) (string, error) {
	encrypted, err := bcrypt.GenerateFromPassword([]byte(pwd), es.cost)
	if err != nil {
		return "", err
	}
	return string(encrypted), nil
}
