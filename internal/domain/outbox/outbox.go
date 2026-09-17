package outbox

import "time"

type InboxRecord struct {
	ID          string
	MessageID   string
	PayloadHash string
	Payload     []byte
	ConsumedAt  time.Time
	Handled     bool
}

type OutboxRecord struct {
	ID          string
	EventType   string
	AggregateID string
	Payload     []byte
	CreatedAt   time.Time
	PublishedAt *time.Time
	Status      string
}
