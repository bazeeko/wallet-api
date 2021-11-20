package main

import (
	"database/sql"
	_userHttpDelivery "homework-9/user/delivery/http"

	_userHttpDeliveryMiddleware "homework-9/user/delivery/http/middleware"
	_userMysqlRepo "homework-9/user/repository/mysql"
	_userUsecase "homework-9/user/usecase"
	_walletHttpDelivery "homework-9/wallet/delivery/http"
	_walletLocalRepo "homework-9/wallet/repository/local"
	_walletMysqlRepo "homework-9/wallet/repository/mysql"
	_walletUsecase "homework-9/wallet/usecase"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
)

func connectDB(config string) (*sql.DB, error) {
	conn, err := sql.Open("mysql", config)
	if err != nil {
		return nil, err
	}

	err = conn.Ping()
	if err != nil {
		return nil, err
	}

	// _, err = conn.Exec(`CREATE DATABASE IF NOT EXISTS h9database`)
	// if err != nil {
	// 	return nil, err
	// }

	// _, err = conn.Exec(`USE h9database`)
	// if err != nil {
	// 	return nil, err
	// }

	_, err = conn.Exec(`CREATE TABLE IF NOT EXISTS users (
		id BIGINT NOT NULL UNIQUE,
		username VARCHAR(40) NOT NULL UNIQUE,
		password VARCHAR(40) NOT NULL,
		PRIMARY KEY (id)
	);`)

	if err != nil {
		return nil, err
	}

	_, err = conn.Exec(`CREATE TABLE IF NOT EXISTS wallets (
		id BIGINT NOT NULL AUTO_INCREMENT UNIQUE,
		user_id BIGINT NOT NULL,
		name VARCHAR(40) NOT NULL,
		amount BIGINT,
		PRIMARY KEY (id),
		FOREIGN KEY (user_id) REFERENCES users(id)
	);`)

	if err != nil {
		return nil, err
	}

	// создаём первого пользователя
	conn.Exec(`INSERT users (id, username, password) VALUES (?, ?, ?)`,
		0, "admin", "pass")

	return conn, nil
}

func main() {
	// config := "root:password@tcp(127.0.0.1:3306)/"
	config := "root:password@(mysqldb)/h9homework"
	conn, err := connectDB(config)
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	e := echo.New()
	middl := _userHttpDeliveryMiddleware.InitMiddleware()
	e.Use(middl.ExecTime)

	userRepo := _userMysqlRepo.NewMysqlUserRepository(conn)

	walletMysqlRepo := _walletMysqlRepo.NewMysqlWalletRepository(conn)
	walletLocalRepo := _walletLocalRepo.NewLocalWalletRepository()

	userUsecase := _userUsecase.NewUserUsecase(userRepo, walletMysqlRepo)
	walletUsecase := _walletUsecase.NewWalletUsecase(walletMysqlRepo, walletLocalRepo)

	_userHttpDelivery.NewUserHandler(e, userRepo)
	_walletHttpDelivery.NewWalletHandler(e, walletUsecase, userUsecase)

	log.Fatalln(e.Start(":8080"))
}

// CREATE USER
// curl -u admin:pass -X POST -v --raw 'localhost:8080/app/user/2?username=superuser&password=superpassword'

// GET USER
// curl -u admin:pass -v --raw 'localhost:8080/app/user/0'
// curl -u admin:pass -v --raw 'localhost:8080/app/user/2'

// CREATE WALLET
// curl -u superuser:superpassword -X POST -v --raw 'localhost:8080/app/wallet/superwallet'
// curl -u admin:pass -X POST -v --raw 'localhost:8080/app/wallet/megawallet'

// START MINING
// curl -u admin:pass -X OPTIONS -v --raw 'localhost:8080/app/wallet/megawallet/start'
// curl -u superuser:superpassword -X OPTIONS -v --raw 'localhost:8080/app/wallet/superwallet/start'

// STOP MINING
// curl -u admin:pass -X OPTIONS -v --raw 'localhost:8080/app/wallet/megawallet/stop'
// curl -u superuser:superpassword -X OPTIONS -v --raw 'localhost:8080/app/wallet/superwallet/stop'
