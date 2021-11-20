package mysql

import (
	"database/sql"
	"fmt"
	"homework-9/domain"
)

type mysqlUserRepository struct {
	*sql.DB
}

func NewMysqlUserRepository(db *sql.DB) domain.MysqlUserRepository {
	return &mysqlUserRepository{db}
}

func (m *mysqlUserRepository) Add(u domain.User) error {
	_, err := m.Exec(`INSERT users (id, username, password) VALUES (?, ?, ?)`,
		u.ID, u.Username, u.Password)
	if err != nil {
		return fmt.Errorf("Add: %w", err)
	}

	return nil
}

func (m *mysqlUserRepository) GetById(id int) (domain.User, error) {
	u := domain.User{}

	err := m.QueryRow(`SELECT id, username, password FROM users WHERE id=?`, id).
		Scan(&u.ID, &u.Username, &u.Password)
	if err != nil {
		return u, fmt.Errorf("Add: %w", err)
	}

	rows, err := m.Query(`SELECT name FROM wallets WHERE user_id=?`, id)
	if err != nil {
		return domain.User{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var walletName string

		err := rows.Scan(&walletName)
		if err != nil {
			return domain.User{}, err
		}

		u.Wallets = append(u.Wallets, walletName)
	}

	return u, nil
}

func (m *mysqlUserRepository) GetByUsername(username string) (domain.User, error) {
	u := domain.User{}

	err := m.QueryRow(`SELECT id, username, password FROM users WHERE username=?`, username).
		Scan(&u.ID, &u.Username, &u.Password)
	if err != nil {
		return u, fmt.Errorf("GetByUsername: %w", err)
	}

	return u, nil
}
