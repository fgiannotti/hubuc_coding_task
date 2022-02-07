package domain

type User struct {
	Name         string `json:"username"`
	Email        string `json:"email"`
	EncryptedPwd string `json:"password"`
}