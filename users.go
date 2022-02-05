package main

type User struct {
	ID           string `json:"id"`
	Name         string `json:"username"`
	Email        string `json:"email"`
	EncryptedPwd string `json:"password"`
}

type UsersRepo interface {
	save(user User) error
	get(userId string) (User, error)
}

type LocalUsersRepo struct {
	DB map[string]User
}

func NewLocalUsersRepo() UsersRepo {
	return &LocalUsersRepo{DB: map[string]User{}}
}

func (l *LocalUsersRepo) save(user User) error {
	panic("implement me")
}

func (l *LocalUsersRepo) get(username string) (User, error) {
	user, ok := l.DB[username]
	if !ok {
		return User{}, UserNotFoundError(username)
	}

	return user, nil
}
