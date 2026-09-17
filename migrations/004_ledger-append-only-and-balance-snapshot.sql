-- Feature: transformar ledger em apêndice auditável com snapshot de saldo.

ALTER TABLE wallet_ledger_entries
    ADD COLUMN balance_before_cents BIGINT,
    ADD COLUMN balance_after_cents BIGINT;

ALTER TABLE wallet_ledger_entries
    ADD CONSTRAINT ck_wallet_ledger_entries_positive_amount
        CHECK (amount_cents >= 0);

ALTER TABLE wallet_ledger_entries
    ADD CONSTRAINT ux_wallet_ledger_entries_wallet_transaction
        UNIQUE (wallet_id, transaction_id);

CREATE OR REPLACE FUNCTION prevent_wallet_ledger_entries_mutation()
RETURNS trigger
LANGUAGE plpgsql
AS $$
BEGIN
    RAISE EXCEPTION 'wallet_ledger_entries is append-only';
END;
$$;

DROP TRIGGER IF EXISTS trg_wallet_ledger_entries_block_update ON wallet_ledger_entries;
CREATE TRIGGER trg_wallet_ledger_entries_block_update
    BEFORE UPDATE ON wallet_ledger_entries
    FOR EACH ROW
    EXECUTE FUNCTION prevent_wallet_ledger_entries_mutation();

DROP TRIGGER IF EXISTS trg_wallet_ledger_entries_block_delete ON wallet_ledger_entries;
CREATE TRIGGER trg_wallet_ledger_entries_block_delete
    BEFORE DELETE ON wallet_ledger_entries
    FOR EACH ROW
    EXECUTE FUNCTION prevent_wallet_ledger_entries_mutation();
