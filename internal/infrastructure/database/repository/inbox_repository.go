package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"JungleGaming-test/internal/domain/outbox"
	"JungleGaming-test/internal/infrastructure/database/transaction"
)

type InboxRepositoryImpl struct {
	db *pgxpool.Pool
}

func NewInboxRepository(db *pgxpool.Pool) *InboxRepositoryImpl {
	return &InboxRepositoryImpl{db: db}
}

func (r *InboxRepositoryImpl) Save(ctx context.Context, record *outbox.InboxRecord) error {
	q := transaction.GetQuery(ctx, r.db)
	query := `
		INSERT INTO inbox (id, message_id, payload_hash, payload, consumed_at, handled)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := q.Exec(ctx, query, record.ID, record.MessageID, record.PayloadHash, record.Payload, record.ConsumedAt, record.Handled)
	if err != nil {
		return fmt.Errorf("failed to save inbox record: %w", err)
	}
	return nil
}

func (r *InboxRepositoryImpl) FindByMessageID(ctx context.Context, messageID string) (*outbox.InboxRecord, error) {
	q := transaction.GetQuery(ctx, r.db)
	query := `
		SELECT id, message_id, payload_hash, payload, consumed_at, handled
		FROM inbox
		WHERE message_id = $1
	`
	var rec outbox.InboxRecord
	err := q.QueryRow(ctx, query, messageID).Scan(&rec.ID, &rec.MessageID, &rec.PayloadHash, &rec.Payload, &rec.ConsumedAt, &rec.Handled)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find inbox record: %w", err)
	}
	return &rec, nil
}
