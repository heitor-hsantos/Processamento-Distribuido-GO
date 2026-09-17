package wallet

import (
	"fmt"

	domainerrors "JungleGaming-test/internal/domain/errors"
	"JungleGaming-test/internal/domain/money"
)

type Wallet struct {
	ID       string
	PlayerID string
	Balance  money.Money
	Currency money.Currency
	Version  int64
}

func NewWallet(id, playerID string, initial money.Money) (*Wallet, error) {
	if id == "" {
		return nil, fmt.Errorf("wallet: %w", domainerrors.ErrInvalidOperation)
	}
	if playerID == "" {
		return nil, fmt.Errorf("wallet: %w", domainerrors.ErrInvalidOperation)
	}
	if err := validateWalletCurrency(initial); err != nil {
		return nil, err
	}
	return &Wallet{
		ID:       id,
		PlayerID: playerID,
		Balance:  initial,
		Currency: initial.Currency(),
		Version:  1,
	}, nil
}

func (w *Wallet) CanWithdraw(amount money.Money) bool {
	if err := validateWalletCurrency(amount); err != nil {
		return false
	}
	if w.Currency != amount.Currency() {
		return false
	}
	return w.Balance.Amount() >= amount.Amount()
}

func (w *Wallet) Withdraw(amount money.Money) error {
	if err := validateWalletCurrency(amount); err != nil {
		return fmt.Errorf("wallet: %w", err)
	}
	if !w.CanWithdraw(amount) {
		return fmt.Errorf("wallet: %w", domainerrors.ErrInsufficientFunds)
	}
	updated, err := w.Balance.Subtract(amount)
	if err != nil {
		return fmt.Errorf("wallet: %w", err)
	}
	w.Balance = updated
	w.Version++
	return nil
}

func (w *Wallet) Deposit(amount money.Money) error {
	if err := validateWalletCurrency(amount); err != nil {
		return fmt.Errorf("wallet: %w", err)
	}
	if w.Currency != amount.Currency() {
		return fmt.Errorf("wallet: %w", domainerrors.ErrInvalidCurrency)
	}
	updated, err := w.Balance.Add(amount)
	if err != nil {
		return fmt.Errorf("wallet: %w", err)
	}
	w.Balance = updated
	w.Version++
	return nil
}

func validateWalletCurrency(amount money.Money) error {
	if err := validateCurrency(amount.Currency()); err != nil {
		return fmt.Errorf("wallet: %w", err)
	}
	return nil
}

func validateCurrency(currency money.Currency) error {
	if currency == "" {
		return fmt.Errorf("wallet: %w", domainerrors.ErrInvalidCurrency)
	}
	if len(currency) != 3 {
		return fmt.Errorf("wallet: %w", domainerrors.ErrInvalidCurrency)
	}
	return nil
}
