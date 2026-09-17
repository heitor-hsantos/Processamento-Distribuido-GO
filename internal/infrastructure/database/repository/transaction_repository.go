package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"JungleGaming-test/internal/domain/money"
	"JungleGaming-test/internal/domain/transaction"
	dbtx "JungleGaming-test/internal/infrastructure/database/transaction"
)

type TransactionRepository struct {
	db *pgxpool.Pool
}

func NewTransactionRepository(db *pgxpool.Pool) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) Save(ctx context.Context, t *transaction.WagerTransaction) error {
	q := dbtx.GetQuery(ctx, r.db)
	_, err := q.Exec(ctx, `
		INSERT INTO wager_transactions (id, wallet_id, provider_id, transaction_type, transaction_state, amount_cents, currency, idempotency_key, reference_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW(), NOW())
		ON CONFLICT (id) DO UPDATE
		SET transaction_state = EXCLUDED.transaction_state,
		    reference_id = EXCLUDED.reference_id,
		    updated_at = NOW()
	`, t.ID, t.WalletID, t.ProviderID, t.Type, t.State, t.Amount.Amount(), t.Amount.Currency(), t.IdempotencyKey, t.ReferenceID)
	if err != nil {
		return fmt.Errorf("transaction repository save: %w", err)
	}
	return nil
}

func (r *TransactionRepository) FindByID(ctx context.Context, id string) (*transaction.WagerTransaction, error) {
	q := dbtx.GetQuery(ctx, r.db)
	return r.scanRow(q.QueryRow(ctx, `
		SELECT id, wallet_id, provider_id, transaction_type, transaction_state, amount_cents, currency, idempotency_key, reference_id 
		FROM wager_transactions WHERE id = $1
	`, id))
}

func (r *TransactionRepository) FindByIdempotencyKey(ctx context.Context, key string) (*transaction.WagerTransaction, error) {
	q := dbtx.GetQuery(ctx, r.db)
	return r.scanRow(q.QueryRow(ctx, `
		SELECT id, wallet_id, provider_id, transaction_type, transaction_state, amount_cents, currency, idempotency_key, reference_id 
		FROM wager_transactions WHERE idempotency_key = $1
	`, key))
}

func (r *TransactionRepository) scanRow(row pgx.Row) (*transaction.WagerTransaction, error) {
	var t transaction.WagerTransaction
	var amountCents int64
	var currency string
	var refID *string
	err := row.Scan(&t.ID, &t.WalletID, &t.ProviderID, &t.Type, &t.State, &amountCents, &currency, &t.IdempotencyKey, &refID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("transaction not found")
		}
		return nil, err
	}
	amt, err := money.New(amountCents, money.Currency(currency))
	if err != nil {
		return nil, err
	}
	t.Amount = amt
	if refID != nil {
		t.ReferenceID = *refID
	}
	return &t, nil
}
