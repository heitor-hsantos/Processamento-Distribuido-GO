-- Feature: proteger evolução de versão da carteira e formato de moeda.

ALTER TABLE wallets
    ALTER COLUMN version SET DEFAULT 1;

ALTER TABLE wallets
    ADD CONSTRAINT ck_wallets_currency_iso
        CHECK (currency ~ '^[A-Z]{3}$');

ALTER TABLE wallets
    ADD CONSTRAINT ck_wallets_version_positive
        CHECK (version >= 1);
