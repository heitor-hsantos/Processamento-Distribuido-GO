package transaction

import (
	"testing"

	"JungleGaming-test/internal/domain/money"
)

func TestWagerTransactionTransitions(t *testing.T) {
	tx, err := NewWagerTransaction("txn-1", "wallet-1", "provider-1", KindBet, money.Must(50, money.BRL), "idem-1")
	if err != nil {
		t.Fatalf("unexpected create error: %v", err)
	}

	if err := tx.TransitionTo(StateProcessed); err != nil {
		t.Fatalf("unexpected processed transition error: %v", err)
	}
	if tx.State != StateProcessed {
		t.Fatalf("expected processed state, got %s", tx.State)
	}
}

func TestWagerTransactionRejectsInvalidStateTransition(t *testing.T) {
	tx, err := NewWagerTransaction("txn-2", "wallet-2", "provider-2", KindBet, money.Must(20, money.BRL), "idem-2")
	if err != nil {
		t.Fatalf("unexpected create error: %v", err)
	}

	if err := tx.TransitionTo(StateProcessed); err != nil {
		t.Fatalf("unexpected processed transition error: %v", err)
	}
	if err := tx.TransitionTo(StatePending); err == nil {
		t.Fatal("expected invalid transition back to pending to fail")
	}
}

func TestWagerTransactionAllowsReferenceWait(t *testing.T) {
	tx, err := NewWagerTransaction("txn-3", "wallet-3", "provider-3", KindRollback, money.Must(30, money.BRL), "idem-3")
	if err != nil {
		t.Fatalf("unexpected create error: %v", err)
	}
	if err := tx.TransitionTo(StatePendingReference); err != nil {
		t.Fatalf("unexpected pending reference transition error: %v", err)
	}
	if err := tx.TransitionTo(StateProcessed); err != nil {
		t.Fatalf("unexpected processed transition after reference wait: %v", err)
	}
}
