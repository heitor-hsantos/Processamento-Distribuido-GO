package usecase

import (
	"context"
	"fmt"

	"JungleGaming-test/internal/application/port"
)

type ReconcileInput struct{ WalletID string }

type ReconcileOutput struct {
	WalletID string
	Status   string
}

type ReconcileUseCase struct {
	walletRepo port.WalletRepository
}

func NewReconcileUseCase(repo port.WalletRepository) *ReconcileUseCase {
	return &ReconcileUseCase{walletRepo: repo}
}

func (uc *ReconcileUseCase) Execute(ctx context.Context, input ReconcileInput) (*ReconcileOutput, error) {
	w, err := uc.walletRepo.FindByID(ctx, input.WalletID)
	if err != nil {
		return nil, fmt.Errorf("reconcile: %w", err)
	}
	return &ReconcileOutput{WalletID: w.ID, Status: "ok"}, nil
}
