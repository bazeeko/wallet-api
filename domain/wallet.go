package domain

type Wallet struct {
	Name   string `json:"name"`
	Amount int64  `json:"amount"`
	// MiningDone chan struct{} `json:"-"`
}

type MysqlWalletRepository interface {
	Add(userdID int, name string) error
	GetByName(userID int, name string) (Wallet, error)
	DeleteByName(userID int, name string) error
	Mine(userID int, name string) error
}

type LocalWalletRepository interface {
	IsMining(userID int, name string) bool
	AddMiningWallet(userID int, name string) error
	DeleteMiningWallet(userID int, name string) error
	StopMining(userID int, name string) error
	CheckMiningDone(userID int, name string, quit chan struct{})
}

type WalletUsecase interface {
	Add(userID int, name string) error
	GetByName(userID int, name string) (Wallet, error)
	DeleteByName(userID int, name string) error
	Mine(userID int, name string) error
	StopMining(userID int, name string) error
}
