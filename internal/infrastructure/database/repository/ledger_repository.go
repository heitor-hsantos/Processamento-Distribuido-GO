package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"JungleGaming-test/internal/domain/ledger"
	"JungleGaming-test/internal/infrastructure/database/transaction"
)

type LedgerRepositoryImpl struct {
	db *pgxpool.Pool
}

func NewLedgerRepository(db *pgxpool.Pool) *LedgerRepositoryImpl {
	return &LedgerRepositoryImpl{db: db}
}

func (r *LedgerRepositoryImpl) Save(ctx context.Context, entry *ledger.LedgerEntry) error {
	q := transaction.GetQuery(ctx, r.db)
	query := `
		INSERT INTO wallet_ledger_entries (id, wallet_id, transaction_id, entry_type, amount_cents, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := q.Exec(ctx, query, entry.ID, entry.WalletID, entry.TransactionID, entry.EntryType, entry.Amount.Amount(), entry.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to save ledger entry: %w", err)
	}
	return nil
}
