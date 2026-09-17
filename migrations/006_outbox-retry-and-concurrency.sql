-- Feature: outbox com retries, disputa segura entre publishers e recuperação.

ALTER TABLE outbox
    ADD COLUMN event_id UUID,
    ADD COLUMN correlation_id TEXT,
    ADD COLUMN causation_id TEXT,
    ADD COLUMN occurred_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    ADD COLUMN retry_count INT NOT NULL DEFAULT 0,
    ADD COLUMN max_retries INT NOT NULL DEFAULT 10,
    ADD COLUMN next_attempt_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    ADD COLUMN claimed_by TEXT,
    ADD COLUMN claimed_at TIMESTAMPTZ,
    ADD COLUMN last_error TEXT;

ALTER TABLE outbox
    ADD CONSTRAINT ck_outbox_status
        CHECK (status IN ('PENDING', 'PUBLISHED', 'FAILED'));

ALTER TABLE outbox
    ADD CONSTRAINT ck_outbox_retry_count_non_negative
        CHECK (retry_count >= 0);

ALTER TABLE outbox
    ADD CONSTRAINT ck_outbox_max_retries_positive
        CHECK (max_retries > 0);

CREATE UNIQUE INDEX ux_outbox_event_id
    ON outbox(event_id)
    WHERE event_id IS NOT NULL;

CREATE INDEX ix_outbox_dispatch
    ON outbox(status, next_attempt_at, created_at);

CREATE INDEX ix_outbox_claimed
    ON outbox(claimed_by, claimed_at)
    WHERE claimed_by IS NOT NULL;
