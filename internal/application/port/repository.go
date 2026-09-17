package port

import (
	"context"

	"JungleGaming-test/internal/domain/ledger"
	"JungleGaming-test/internal/domain/outbox"
	"JungleGaming-test/internal/domain/transaction"
	"JungleGaming-test/internal/domain/wallet"
)

type WalletRepository interface {
	Save(context.Context, *wallet.Wallet) error
	FindByID(context.Context, string) (*wallet.Wallet, error)
	FindByPlayerID(context.Context, string) (*wallet.Wallet, error)
}

type TransactionRepository interface {
	Save(context.Context, *transaction.WagerTransaction) error
	FindByID(context.Context, string) (*transaction.WagerTransaction, error)
	FindByIdempotencyKey(context.Context, string) (*transaction.WagerTransaction, error)
}

type LedgerRepository interface {
	Save(context.Context, *ledger.LedgerEntry) error
}

type OutboxRepository interface {
	Save(context.Context, *outbox.OutboxRecord) error
	FindPending(context.Context, int) ([]*outbox.OutboxRecord, error)
	MarkAsPublished(context.Context, string) error
}

type InboxRepository interface {
	Save(context.Context, *outbox.InboxRecord) error
	FindByMessageID(context.Context, string) (*outbox.InboxRecord, error)
}

type UnitOfWork interface {
	Do(ctx context.Context, fn func(ctx context.Context) error) error
}
