package ledger

import (
	"time"

	"github.com/google/uuid"

	"JungleGaming-test/internal/domain/money"
)

type EntryType string

const (
	EntryTypeDebit  EntryType = "DEBIT"
	EntryTypeCredit EntryType = "CREDIT"
)

type LedgerEntry struct {
	ID            string
	WalletID      string
	TransactionID string
	EntryType     EntryType
	Amount        money.Money
	CreatedAt     time.Time
}

func NewLedgerEntry(walletID, transactionID string, entryType EntryType, amount money.Money) *LedgerEntry {
	return &LedgerEntry{
		ID:            uuid.New().String(),
		WalletID:      walletID,
		TransactionID: transactionID,
		EntryType:     entryType,
		Amount:        amount,
		CreatedAt:     time.Now().UTC(),
	}
}
