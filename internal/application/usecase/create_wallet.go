package usecase

import (
	"context"
	"fmt"

	"JungleGaming-test/internal/application/port"
	domainerrors "JungleGaming-test/internal/domain/errors"
	"JungleGaming-test/internal/domain/money"
	"JungleGaming-test/internal/domain/wallet"
	"github.com/google/uuid"
)

type CreateWalletInput struct {
	WalletID       string
	PlayerID       string
	Currency       money.Currency
	InitialBalance money.Money
}

type CreateWalletOutput struct {
	WalletID string
	Balance  int64
	Currency money.Currency
	Version  int64
}

type CreateWalletUseCase struct {
	walletRepo port.WalletRepository
}

func NewCreateWalletUseCase(repo port.WalletRepository) *CreateWalletUseCase {
	return &CreateWalletUseCase{walletRepo: repo}
}

func (uc *CreateWalletUseCase) Execute(ctx context.Context, input CreateWalletInput) (*CreateWalletOutput, error) {
	if input.WalletID == "" {
		input.WalletID = uuid.NewString()
	}
	if input.PlayerID == "" {
		return nil, fmt.Errorf("create-wallet: %w", domainerrors.ErrInvalidOperation)
	}
	if input.Currency == "" {
		input.Currency = money.BRL
	}

	initial := input.InitialBalance
	if initial.Currency() == "" {
		m, err := money.New(0, input.Currency)
		if err != nil {
			return nil, fmt.Errorf("create-wallet: %w", err)
		}
		initial = m
	}
	if initial.Currency() != input.Currency {
		return nil, fmt.Errorf("create-wallet: %w", domainerrors.ErrInvalidCurrency)
	}

	w, err := wallet.NewWallet(input.WalletID, input.PlayerID, initial)
	if err != nil {
		return nil, fmt.Errorf("create-wallet: %w", err)
	}

	if err := uc.walletRepo.Save(ctx, w); err != nil {
		return nil, fmt.Errorf("create-wallet: %w", err)
	}

	return &CreateWalletOutput{
		WalletID: w.ID,
		Balance:  w.Balance.Amount(),
		Currency: w.Currency,
		Version:  w.Version,
	}, nil
}
