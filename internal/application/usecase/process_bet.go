package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"

	"JungleGaming-test/internal/application/port"
	"JungleGaming-test/internal/domain/ledger"
	"JungleGaming-test/internal/domain/money"
	"JungleGaming-test/internal/domain/outbox"
	"JungleGaming-test/internal/domain/transaction"
)

type ProcessBetInput struct {
	WalletID       string
	ProviderID     string
	TransactionID  string
	Amount         int64
	Currency       money.Currency
	IdempotencyKey string
}

type ProcessBetOutput struct {
	TransactionID string
	State         transaction.State
}

type ProcessBetUseCase struct {
	walletRepo port.WalletRepository
	trxRepo    port.TransactionRepository
	ledgerRepo port.LedgerRepository
	outboxRepo port.OutboxRepository
	uow        port.UnitOfWork
}

func NewProcessBetUseCase(
	walletRepo port.WalletRepository,
	trxRepo port.TransactionRepository,
	ledgerRepo port.LedgerRepository,
	outboxRepo port.OutboxRepository,
	uow port.UnitOfWork,
) *ProcessBetUseCase {
	return &ProcessBetUseCase{
		walletRepo: walletRepo,
		trxRepo:    trxRepo,
		ledgerRepo: ledgerRepo,
		outboxRepo: outboxRepo,
		uow:        uow,
	}
}

func (uc *ProcessBetUseCase) Execute(ctx context.Context, input ProcessBetInput) (*ProcessBetOutput, error) {
	amount, err := money.New(input.Amount, input.Currency)
	if err != nil {
		return nil, fmt.Errorf("process-bet: %w", err)
	}

	var finalTrx *transaction.WagerTransaction

	err = uc.uow.Do(ctx, func(ctxWithTx context.Context) error {
		w, err := uc.walletRepo.FindByID(ctxWithTx, input.WalletID)
		if err != nil {
			return fmt.Errorf("find wallet: %w", err)
		}

		trx, err := transaction.NewWagerTransaction(input.TransactionID, input.WalletID, input.ProviderID, transaction.KindBet, amount, input.IdempotencyKey)
		if err != nil {
			return fmt.Errorf("new transaction: %w", err)
		}

		if !w.CanWithdraw(amount) {
			if err := trx.TransitionTo(transaction.StateRejected); err != nil {
				return err
			}
			if err := uc.trxRepo.Save(ctxWithTx, trx); err != nil {
				return err
			}
			finalTrx = trx
			return nil
		}

		if err := w.Withdraw(amount); err != nil {
			return fmt.Errorf("withdraw: %w", err)
		}

		if err := trx.TransitionTo(transaction.StateProcessed); err != nil {
			return err
		}

		if err := uc.walletRepo.Save(ctxWithTx, w); err != nil {
			return fmt.Errorf("save wallet: %w", err)
		}
		if err := uc.trxRepo.Save(ctxWithTx, trx); err != nil {
			return fmt.Errorf("save transaction: %w", err)
		}

		// Save Ledger Entry
		lEntry := ledger.NewLedgerEntry(w.ID, trx.ID, ledger.EntryTypeDebit, amount)
		if err := uc.ledgerRepo.Save(ctxWithTx, lEntry); err != nil {
			return fmt.Errorf("save ledger: %w", err)
		}

		// Save Outbox Event
		payload, _ := json.Marshal(map[string]interface{}{
			"transaction_id": trx.ID,
			"wallet_id":      w.ID,
			"amount":         amount.Amount(),
			"currency":       amount.Currency(),
			"state":          trx.State,
		})
		outboxRec := &outbox.OutboxRecord{
			ID:          uuid.New().String(),
			EventType:   "BET_PROCESSED",
			AggregateID: trx.ID,
			Payload:     payload,
			CreatedAt:   time.Now().UTC(),
			Status:      "PENDING",
		}
		if err := uc.outboxRepo.Save(ctxWithTx, outboxRec); err != nil {
			return fmt.Errorf("save outbox: %w", err)
		}

		finalTrx = trx
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("process-bet: %w", err)
	}

	return &ProcessBetOutput{TransactionID: finalTrx.ID, State: finalTrx.State}, nil
}
