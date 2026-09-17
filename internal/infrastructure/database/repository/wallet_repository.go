package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"JungleGaming-test/internal/domain/money"
	"JungleGaming-test/internal/domain/wallet"
	"JungleGaming-test/internal/infrastructure/database/transaction"
)

type WalletRepository struct {
	db *pgxpool.Pool
}

func NewWalletRepository(db *pgxpool.Pool) *WalletRepository {
	return &WalletRepository{db: db}
}

func (r *WalletRepository) Save(ctx context.Context, w *wallet.Wallet) error {
	if w == nil {
		return fmt.Errorf("wallet repository: invalid wallet")
	}

	q := transaction.GetQuery(ctx, r.db)
	_, err := q.Exec(ctx, `
		INSERT INTO wallets (id, player_id, currency, balance_cents, version, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
		ON CONFLICT (id) DO UPDATE
		SET player_id = EXCLUDED.player_id,
		    currency = EXCLUDED.currency,
		    balance_cents = EXCLUDED.balance_cents,
		    version = EXCLUDED.version,
		    updated_at = NOW()
	`, w.ID, w.PlayerID, w.Currency, w.Balance.Amount(), w.Version)
	if err != nil {
		return fmt.Errorf("wallet repository: %w", err)
	}
	return nil
}

func (r *WalletRepository) FindByID(ctx context.Context, id string) (*wallet.Wallet, error) {
	if id == "" {
		return nil, fmt.Errorf("wallet repository: invalid id")
	}

	q := transaction.GetQuery(ctx, r.db)
	var (
		walletID, playerID, currency string
		balanceCents                 int64
		version                      int64
	)
	// We use FOR UPDATE to acquire a row-level lock during the transaction
	row := q.QueryRow(ctx, `SELECT id, player_id, currency, balance_cents, version FROM wallets WHERE id = $1 FOR UPDATE`, id)
	if err := row.Scan(&walletID, &playerID, &currency, &balanceCents, &version); err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("wallet repository: wallet not found")
		}
		return nil, fmt.Errorf("wallet repository: %w", err)
	}

	amount, err := money.New(balanceCents, money.Currency(currency))
	if err != nil {
		return nil, fmt.Errorf("wallet repository: %w", err)
	}
	return &wallet.Wallet{ID: walletID, PlayerID: playerID, Balance: amount, Currency: money.Currency(currency), Version: version}, nil
}

func (r *WalletRepository) FindByPlayerID(ctx context.Context, playerID string) (*wallet.Wallet, error) {
	if playerID == "" {
		return nil, fmt.Errorf("wallet repository: invalid player id")
	}

	q := transaction.GetQuery(ctx, r.db)
	var (
		walletID, storedPlayerID, currency string
		balanceCents                       int64
		version                            int64
	)
	row := q.QueryRow(ctx, `SELECT id, player_id, currency, balance_cents, version FROM wallets WHERE player_id = $1 LIMIT 1`, playerID)
	if err := row.Scan(&walletID, &storedPlayerID, &currency, &balanceCents, &version); err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("wallet repository: wallet not found")
		}
		return nil, fmt.Errorf("wallet repository: %w", err)
	}
	amount, err := money.New(balanceCents, money.Currency(currency))
	if err != nil {
		return nil, fmt.Errorf("wallet repository: %w", err)
	}
	return &wallet.Wallet{ID: walletID, PlayerID: storedPlayerID, Balance: amount, Currency: money.Currency(currency), Version: version}, nil
}
