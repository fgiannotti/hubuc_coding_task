package services

import "golang.org/x/crypto/bcrypt"

type Encryptions interface {
	Encrypt(pwd string) (string, error)
	Compare(encryptedPwd string, pwd string) (bool,error)
}
type BcryptEncryptionsService struct {
	cost int
}

func NewBcryptEncryptionsService() Encryptions {
	return &BcryptEncryptionsService{cost: bcrypt.MinCost}
}

func (es *BcryptEncryptionsService) Compare(encryptedPwd string, pwd string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(encryptedPwd), []byte(pwd))
	if err != nil {
		return false, err
	}
	return true, nil
}

func (es *BcryptEncryptionsService) Encrypt(pwd string) (string, error) {
	encrypted, err := bcrypt.GenerateFromPassword([]byte(pwd), es.cost)
	if err != nil {
		return "", err
	}
	return string(encrypted), nil
}
