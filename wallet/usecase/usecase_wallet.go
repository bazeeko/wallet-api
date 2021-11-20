package usecase

import (
	"context"
	"fmt"
	"homework-9/domain"
	"log"
	"time"
)

type walletUsecase struct {
	mysqlWalletRepo domain.MysqlWalletRepository
	localWalletRepo domain.LocalWalletRepository
}

func NewWalletUsecase(mwr domain.MysqlWalletRepository, lwr domain.LocalWalletRepository) domain.WalletUsecase {
	return &walletUsecase{
		mysqlWalletRepo: mwr,
		localWalletRepo: lwr,
	}
}

func (uc *walletUsecase) Add(userID int, name string) error {
	w, _ := uc.mysqlWalletRepo.GetByName(userID, name)
	if w != (domain.Wallet{}) {
		return fmt.Errorf("wallet '%s' already exists", name)
	}

	err := uc.mysqlWalletRepo.Add(userID, name)
	if err != nil {
		return err
	}

	return nil
}

func (uc *walletUsecase) GetByName(userID int, name string) (domain.Wallet, error) {
	w, _ := uc.mysqlWalletRepo.GetByName(userID, name)
	if w == (domain.Wallet{}) {
		return domain.Wallet{}, fmt.Errorf("user with id %d doesn't have wallet '%s'", userID, name)
	}

	return w, nil
}

func (uc *walletUsecase) DeleteByName(userID int, name string) error {
	w, _ := uc.mysqlWalletRepo.GetByName(userID, name)
	if w == (domain.Wallet{}) {
		return fmt.Errorf("the wallet '%s' does not exist", name)
	}

	if uc.localWalletRepo.IsMining(userID, name) {
		return fmt.Errorf("the wallet '%s' is currently mining", name)
	}

	err := uc.mysqlWalletRepo.DeleteByName(userID, name)
	if err != nil {
		return err
	}

	return nil
}

func (uc *walletUsecase) Mine(userID int, name string) error {
	w, _ := uc.mysqlWalletRepo.GetByName(userID, name)
	if w == (domain.Wallet{}) {
		return fmt.Errorf("the wallet '%s' doesn't exist", name)
	}

	if err := uc.localWalletRepo.AddMiningWallet(userID, w.Name); err != nil {
		fmt.Println(err)
		return err
	}

	go func() {

		quit := make(chan struct{})

	LOOP:
		for {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()

			select {
			case <-quit:
				cancel()
				break LOOP
			case <-ctx.Done():
				if err := uc.mysqlWalletRepo.Mine(userID, w.Name); err != nil {
					break LOOP
				}
				uc.localWalletRepo.CheckMiningDone(userID, w.Name, quit)
				cancel()
			}
		}
		func() {
			err := uc.localWalletRepo.DeleteMiningWallet(userID, name)
			if err != nil {
				log.Printf("DeleteMiningWallet: %s", err)
			}
		}()
	}()

	return nil
}

func (uc *walletUsecase) StopMining(userID int, name string) error {
	w, _ := uc.mysqlWalletRepo.GetByName(userID, name)
	if w == (domain.Wallet{}) {
		return fmt.Errorf("the wallet '%s' doesn't exist", name)
	}

	if !uc.localWalletRepo.IsMining(userID, w.Name) {
		return fmt.Errorf("the wallet '%s' is not mining", w.Name)
	}

	err := uc.localWalletRepo.StopMining(userID, w.Name)
	if err != nil {
		return err
	}

	return nil
}
