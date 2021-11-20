package local

import (
	"fmt"
	"homework-9/domain"
	"sync"
)

type localWalletRepository struct {
	activeWallets map[string](chan struct{})
	sync.Mutex
}

func NewLocalWalletRepository() domain.LocalWalletRepository {
	return &localWalletRepository{activeWallets: make(map[string](chan struct{}))}
}

func (wc *localWalletRepository) IsMining(userID int, name string) bool {
	wc.Lock()
	defer wc.Unlock()

	name = fmt.Sprintf("%s_%d", name, userID)

	_, exists := wc.activeWallets[name]

	return exists
}

func (wc *localWalletRepository) AddMiningWallet(userID int, name string) error {
	wc.Lock()
	defer wc.Unlock()

	name = fmt.Sprintf("%s_%d", name, userID)

	_, exists := wc.activeWallets[name]
	if exists {
		return fmt.Errorf("the wallet '%s' is already mining", name)
	}

	// wallet.MiningDone = make(chan struct{})
	wc.activeWallets[name] = make(chan struct{})

	return nil
}

func (wc *localWalletRepository) DeleteMiningWallet(userID int, name string) error {
	wc.Lock()
	defer wc.Unlock()

	name = fmt.Sprintf("%s_%d", name, userID)

	_, exists := wc.activeWallets[name]
	if !exists {
		return fmt.Errorf("the wallet '%s' is not mining", name)
	}

	delete(wc.activeWallets, name)

	return nil
}

func (wc *localWalletRepository) StopMining(userID int, name string) error {
	wc.Lock()

	name = fmt.Sprintf("%s_%d", name, userID)

	_, exists := wc.activeWallets[name]
	if !exists {
		return fmt.Errorf("the wallet '%s' is not mining", name)
	}

	defer wc.Unlock()

	// wc.activeWallets[name].MiningDone <- struct{}{}

	wc.activeWallets[name] <- struct{}{}

	return nil
}

func (wc *localWalletRepository) CheckMiningDone(userID int, name string, quit chan struct{}) {
	// wc.Lock()
	// defer wc.Unlock()

	name = fmt.Sprintf("%s_%d", name, userID)

	select {
	case <-wc.activeWallets[name]:
		go func() {
			quit <- struct{}{}
		}()
	default:
		break
	}
}
