package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"JungleGaming-test/internal/domain/outbox"
	"JungleGaming-test/internal/infrastructure/database/transaction"
)

type OutboxRepositoryImpl struct {
	db *pgxpool.Pool
}

func NewOutboxRepository(db *pgxpool.Pool) *OutboxRepositoryImpl {
	return &OutboxRepositoryImpl{db: db}
}

func (r *OutboxRepositoryImpl) Save(ctx context.Context, record *outbox.OutboxRecord) error {
	q := transaction.GetQuery(ctx, r.db)
	query := `
		INSERT INTO outbox (id, event_type, aggregate_id, payload, created_at, status)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := q.Exec(ctx, query, record.ID, record.EventType, record.AggregateID, record.Payload, record.CreatedAt, record.Status)
	if err != nil {
		return fmt.Errorf("failed to save outbox record: %w", err)
	}
	return nil
}

func (r *OutboxRepositoryImpl) FindPending(ctx context.Context, limit int) ([]*outbox.OutboxRecord, error) {
	q := transaction.GetQuery(ctx, r.db)
	query := `
		SELECT id, event_type, aggregate_id, payload, created_at, published_at, status
		FROM outbox
		WHERE status = 'PENDING'
		ORDER BY created_at ASC
		LIMIT $1
	`
	rows, err := q.Query(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to find pending outbox records: %w", err)
	}
	defer rows.Close()

	var records []*outbox.OutboxRecord
	for rows.Next() {
		var rec outbox.OutboxRecord
		if err := rows.Scan(&rec.ID, &rec.EventType, &rec.AggregateID, &rec.Payload, &rec.CreatedAt, &rec.PublishedAt, &rec.Status); err != nil {
			return nil, fmt.Errorf("failed to scan outbox record: %w", err)
		}
		records = append(records, &rec)
	}
	return records, nil
}

func (r *OutboxRepositoryImpl) MarkAsPublished(ctx context.Context, id string) error {
	q := transaction.GetQuery(ctx, r.db)
	query := `
		UPDATE outbox
		SET status = 'PUBLISHED', published_at = $1
		WHERE id = $2
	`
	_, err := q.Exec(ctx, query, time.Now().UTC(), id)
	if err != nil {
		return fmt.Errorf("failed to mark outbox record as published: %w", err)
	}
	return nil
}
