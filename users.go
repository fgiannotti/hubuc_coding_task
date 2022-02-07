package main

type User struct {
	Name         string `json:"username"`
	Email        string `json:"email"`
	EncryptedPwd string `json:"password"`
}

type UsersRepo interface {
	Save(user User) error
	Get(username string) (User, error)
}

type LocalUsersRepo struct {
	DB map[string]User
}

func NewLocalUsersRepo() UsersRepo {
	return &LocalUsersRepo{DB: map[string]User{}}
}

func (l *LocalUsersRepo) Save(user User) error {
	panic("implement me")
}

func (l *LocalUsersRepo) Get(username string) (User, error) {
	user, ok := l.DB[username]
	if !ok {
		return User{}, UserNotFoundError(username)
	}

	return user, nil
}
