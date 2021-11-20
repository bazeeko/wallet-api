package usercase

import (
	"fmt"
	"homework-9/domain"
)

type userUsecase struct {
	userRepo        domain.MysqlUserRepository
	mysqlWalletRepo domain.MysqlWalletRepository
}

func NewUserUsecase(ur domain.MysqlUserRepository, mwr domain.MysqlWalletRepository) domain.UserUsecase {
	return &userUsecase{ur, mwr}
}

func (uc *userUsecase) Add(u domain.User) error {
	user, _ := uc.userRepo.GetById(u.ID)

	fmt.Println("AAAAAAA")

	if user.ID == u.ID {
		return fmt.Errorf("user id already exists")
	}

	err := uc.userRepo.Add(u)
	if err != nil {
		return err
	}

	return nil
}

func (uc *userUsecase) GetById(id int) (domain.User, error) {
	fmt.Println("in getbyid")
	user, _ := uc.userRepo.GetById(id)
	fmt.Println("aftergetbyid repo")
	if user.ID != id {
		return domain.User{}, fmt.Errorf("user doesn't exists")
	}

	fmt.Println("user in gbi", user)

	// wallets, err := uc.mysqlWalletRepo.GetWalletsByUserID(id)
	// if err != nil {
	// 	fmt.Println("wallet error", err)
	// 	return domain.User{}, err
	// }

	// fmt.Println("wallets", wallets)
	// user.Wallets = wallets

	fmt.Println("user wallets", user.Wallets)
	return user, nil
}

func (uc *userUsecase) GetByUsername(username string) (domain.User, error) {
	user, _ := uc.userRepo.GetByUsername(username)
	if user.Username != username {
		return domain.User{}, fmt.Errorf("user doesn't exists")
	}

	return user, nil
}
