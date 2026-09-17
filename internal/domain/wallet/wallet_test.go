package wallet

import (
	"testing"

	"JungleGaming-test/internal/domain/money"
)

func TestWalletDepositAndWithdraw(t *testing.T) {
	w, err := NewWallet("wallet-1", "player-1", money.Must(100, money.BRL))
	if err != nil {
		t.Fatalf("unexpected wallet creation error: %v", err)
	}

	if err := w.Deposit(money.Must(25, money.BRL)); err != nil {
		t.Fatalf("unexpected deposit error: %v", err)
	}
	if w.Balance.Amount() != 125 {
		t.Fatalf("expected balance 125, got %d", w.Balance.Amount())
	}

	if err := w.Withdraw(money.Must(30, money.BRL)); err != nil {
		t.Fatalf("unexpected withdraw error: %v", err)
	}
	if w.Balance.Amount() != 95 {
		t.Fatalf("expected balance 95, got %d", w.Balance.Amount())
	}
}

func TestWalletRejectsInsufficientFunds(t *testing.T) {
	w, err := NewWallet("wallet-2", "player-2", money.Must(50, money.BRL))
	if err != nil {
		t.Fatalf("unexpected wallet creation error: %v", err)
	}

	if err := w.Withdraw(money.Must(51, money.BRL)); err == nil {
		t.Fatal("expected insufficient funds error")
	}
}

func TestWalletRejectsCurrencyMismatch(t *testing.T) {
	w, err := NewWallet("wallet-3", "player-3", money.Must(50, money.BRL))
	if err != nil {
		t.Fatalf("unexpected wallet creation error: %v", err)
	}

	if err := w.Deposit(money.Must(10, money.USD)); err == nil {
		t.Fatal("expected currency mismatch to fail")
	}
}
