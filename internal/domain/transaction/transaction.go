package transaction

import (
	"fmt"
	"time"

	domainerrors "JungleGaming-test/internal/domain/errors"
	"JungleGaming-test/internal/domain/money"
)

type Kind string

type State string

const (
	KindOpening  Kind = "OPENING"
	KindBet      Kind = "BET"
	KindWin      Kind = "WIN"
	KindLoss     Kind = "LOSS"
	KindRefund   Kind = "REFUND"
	KindRollback Kind = "ROLLBACK"
)

const (
	StatePending          State = "PENDING"
	StatePendingReference State = "PENDING_REFERENCE"
	StateProcessed        State = "PROCESSED"
	StateRejected         State = "REJECTED"
	StateFailed           State = "FAILED"
)

type WagerTransaction struct {
	ID             string
	WalletID       string
	ProviderID     string
	Type           Kind
	State          State
	Amount         money.Money
	Currency       money.Currency
	IdempotencyKey string
	ReferenceID    string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func NewWagerTransaction(id, walletID, providerID string, kind Kind, amount money.Money, idempotencyKey string) (*WagerTransaction, error) {
	if id == "" || walletID == "" || providerID == "" {
		return nil, fmt.Errorf("transaction: %w", domainerrors.ErrInvalidOperation)
	}
	if kind == "" {
		return nil, fmt.Errorf("transaction: %w", domainerrors.ErrInvalidOperation)
	}
	if idempotencyKey == "" {
		return nil, fmt.Errorf("transaction: %w", domainerrors.ErrInvalidOperation)
	}
	if amount.Currency() == "" {
		return nil, fmt.Errorf("transaction: %w", domainerrors.ErrInvalidCurrency)
	}
	if amount.Amount() < 0 {
		return nil, fmt.Errorf("transaction: %w", domainerrors.ErrInvalidAmount)
	}

	now := time.Now().UTC()
	return &WagerTransaction{
		ID:             id,
		WalletID:       walletID,
		ProviderID:     providerID,
		Type:           kind,
		State:          StatePending,
		Amount:         amount,
		Currency:       amount.Currency(),
		IdempotencyKey: idempotencyKey,
		CreatedAt:      now,
		UpdatedAt:      now,
	}, nil
}

func (t *WagerTransaction) TransitionTo(next State) error {
	if !isValidTransition(t.State, next) {
		return fmt.Errorf("transaction: %w", domainerrors.ErrInvalidOperation)
	}
	t.State = next
	t.UpdatedAt = time.Now().UTC()
	return nil
}

func (t *WagerTransaction) IsDebit() bool {
	return t.Type == KindBet || t.Type == KindRollback
}

func (t *WagerTransaction) IsCredit() bool {
	return t.Type == KindWin || t.Type == KindRefund || t.Type == KindOpening
}

func (t *WagerTransaction) IsFinal() bool {
	return t.State == StateProcessed || t.State == StateRejected || t.State == StateFailed
}

func isValidTransition(current, next State) bool {
	if current == "" {
		current = StatePending
	}

	switch current {
	case StatePending:
		return next == StatePending || next == StatePendingReference || next == StateProcessed || next == StateRejected || next == StateFailed
	case StatePendingReference:
		return next == StatePendingReference || next == StateProcessed || next == StateRejected || next == StateFailed
	case StateProcessed, StateRejected, StateFailed:
		return false
	default:
		return false
	}
}
