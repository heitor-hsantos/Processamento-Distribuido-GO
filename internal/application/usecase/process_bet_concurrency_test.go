package usecase

import (
	"sync"
	"testing"

	"JungleGaming-test/internal/domain/money"
	"JungleGaming-test/internal/domain/wallet"
)

func TestWalletOperationsRemainConsistentUnderSerializedAccess(t *testing.T) {
	w, err := wallet.NewWallet("wallet-concurrency-1", "player-concurrency-1", money.Must(10000, money.BRL))
	if err != nil {
		t.Fatalf("wallet creation failed: %v", err)
	}

	var mu sync.Mutex
	var successful int
	for i := 0; i < 5; i++ {
		mu.Lock()
		err = w.Withdraw(money.Must(2000, money.BRL))
		mu.Unlock()
		if err == nil {
			successful++
		}
	}
	if successful != 5 {
		t.Fatalf("expected 5 successful withdrawals, got %d", successful)
	}
	if w.Balance.Amount() != 0 {
		t.Fatalf("expected zero balance after 5 withdrawals, got %d", w.Balance.Amount())
	}
}
