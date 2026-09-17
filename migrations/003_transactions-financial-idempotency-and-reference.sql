-- Feature: reforçar idempotência financeira e referências entre transações externas.

ALTER TABLE wager_transactions
    ADD COLUMN external_transaction_id TEXT,
    ADD COLUMN payload_hash TEXT,
    ADD COLUMN player_id UUID,
    ADD COLUMN round_id TEXT,
    ADD COLUMN game_id TEXT,
    ADD COLUMN failure_code TEXT,
    ADD COLUMN processed_balance_cents BIGINT;

ALTER TABLE wager_transactions
    ADD CONSTRAINT ck_wager_transactions_currency_iso
        CHECK (currency ~ '^[A-Z]{3}$');

ALTER TABLE wager_transactions
    ADD CONSTRAINT ck_wager_transactions_type
        CHECK (transaction_type IN ('OPENING', 'BET', 'WIN', 'LOSS', 'REFUND', 'ROLLBACK'));

ALTER TABLE wager_transactions
    ADD CONSTRAINT ck_wager_transactions_state
        CHECK (transaction_state IN ('PENDING', 'PENDING_REFERENCE', 'PROCESSED', 'REJECTED', 'FAILED'));

ALTER TABLE wager_transactions
    ADD CONSTRAINT ck_wager_transactions_positive_amount
        CHECK (amount_cents >= 0);

ALTER TABLE wager_transactions
    ADD CONSTRAINT ck_wager_transactions_payload_hash_presence
        CHECK ((idempotency_key IS NULL AND payload_hash IS NULL) OR (idempotency_key IS NOT NULL AND payload_hash IS NOT NULL));

CREATE UNIQUE INDEX ux_wager_transactions_provider_external
    ON wager_transactions(provider_id, external_transaction_id)
    WHERE external_transaction_id IS NOT NULL;

CREATE INDEX ix_wager_transactions_reference_lookup
    ON wager_transactions(provider_id, reference_id)
    WHERE reference_id IS NOT NULL;
