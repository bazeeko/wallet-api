package mysql

import (
	"database/sql"
	"fmt"
	"homework-9/domain"
)

type mysqlWalletRepository struct {
	*sql.DB
}

func NewMysqlWalletRepository(db *sql.DB) domain.MysqlWalletRepository {
	return &mysqlWalletRepository{db}
}

func (m *mysqlWalletRepository) Mine(userID int, name string) error {
	// err := m.QueryRow(`SELECT name, amount FROM wallets WHERE user_id=? AND name=?`, userID, name).
	// 	Scan(&w.Name, &w.Amount)

	// if err != nil {
	// 	return fmt.Errorf("Mine: %w", err)
	// }

	_, err := m.Exec(`UPDATE wallets SET amount=amount+1 WHERE user_id=? AND name=?`,
		userID, name)

	if err != nil {
		return fmt.Errorf("Mine: %w", err)
	}

	return nil
}

func (m *mysqlWalletRepository) Add(userID int, name string) error {
	_, err := m.Exec(`INSERT wallets (user_id, name, amount) VALUES (?, ?, ?)`, userID, name, 0)
	if err != nil {
		return fmt.Errorf("Add: %w", err)
	}

	return nil
}

func (m *mysqlWalletRepository) GetByName(userID int, name string) (domain.Wallet, error) {
	w := domain.Wallet{}

	err := m.QueryRow(`SELECT name, amount FROM wallets WHERE user_id=? AND name=?`, userID, name).
		Scan(&w.Name, &w.Amount)
	if err != nil {
		return domain.Wallet{}, fmt.Errorf("DeleteByName: %w", err)
	}

	return w, nil
}

func (m *mysqlWalletRepository) DeleteByName(userID int, name string) error {
	_, err := m.Exec(`DELETE FROM wallets WHERE user_id=? AND name=?`, userID, name)
	if err != nil {
		return fmt.Errorf("DeleteByName: %w", err)
	}

	return nil
}

func (m *mysqlWalletRepository) GetByUsername(username string) (domain.User, error) {
	u := domain.User{}

	err := m.QueryRow(`SELECT id, username, password FROM users WHERE username=?`, username).
		Scan(&u.ID, &u.Username, &u.Password)
	if err != nil {
		return u, fmt.Errorf("Add: %w", err)
	}

	return u, nil
}
