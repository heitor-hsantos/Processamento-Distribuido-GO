-- Feature: deduplicação durável da inbox por consumidor + hash estável.

ALTER TABLE inbox
    ADD COLUMN consumer_name TEXT NOT NULL DEFAULT 'bet-consumer';

ALTER TABLE inbox
    ADD CONSTRAINT ck_inbox_payload_hash_not_empty
        CHECK (length(payload_hash) > 0);

ALTER TABLE inbox
    ADD CONSTRAINT ck_inbox_consumer_name_not_empty
        CHECK (length(consumer_name) > 0);

DROP INDEX IF EXISTS inbox_message_id_key;

CREATE UNIQUE INDEX ux_inbox_consumer_message
    ON inbox(consumer_name, message_id);

CREATE INDEX ix_inbox_handled
    ON inbox(handled, consumed_at);
