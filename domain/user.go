package domain

type User struct {
	ID       int      `json:"id"`
	Username string   `json:"username"`
	Password string   `json:"-"`
	Wallets  []string `json:"wallets"`
}

type MysqlUserRepository interface {
	Add(User) error
	GetById(id int) (User, error)
	GetByUsername(username string) (User, error)
}

type UserUsecase interface {
	Add(User) error
	GetById(id int) (User, error)
	GetByUsername(username string) (User, error)
}
