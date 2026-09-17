package usecase_test

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"

	"JungleGaming-test/internal/application/usecase"
	"JungleGaming-test/internal/domain/ledger"
	"JungleGaming-test/internal/domain/money"
	"JungleGaming-test/internal/domain/outbox"
	"JungleGaming-test/internal/domain/transaction"
	"JungleGaming-test/internal/domain/wallet"
)

type mockWalletRepo struct {
	wallet *wallet.Wallet
}

func (m *mockWalletRepo) Save(ctx context.Context, w *wallet.Wallet) error {
	bal, _ := money.New(w.Balance.Amount(), w.Currency)
	m.wallet.Balance = bal
	return nil
}

func (m *mockWalletRepo) FindByID(ctx context.Context, id string) (*wallet.Wallet, error) {
	bal, _ := money.New(m.wallet.Balance.Amount(), m.wallet.Currency)
	return &wallet.Wallet{
		ID:       m.wallet.ID,
		PlayerID: m.wallet.PlayerID,
		Balance:  bal,
		Currency: m.wallet.Currency,
		Version:  m.wallet.Version,
	}, nil
}
func (m *mockWalletRepo) FindByPlayerID(ctx context.Context, id string) (*wallet.Wallet, error) {
	return nil, nil
}

type mockTrxRepo struct {
	tx map[string]*transaction.WagerTransaction
}

func (m *mockTrxRepo) Save(ctx context.Context, t *transaction.WagerTransaction) error {
	m.tx[t.ID] = t
	return nil
}
func (m *mockTrxRepo) FindByID(ctx context.Context, id string) (*transaction.WagerTransaction, error) {
	return nil, nil
}
func (m *mockTrxRepo) FindByIdempotencyKey(ctx context.Context, key string) (*transaction.WagerTransaction, error) {
	return nil, nil
}

type mockLedgerRepo struct{}

func (m *mockLedgerRepo) Save(ctx context.Context, entry *ledger.LedgerEntry) error { return nil }

type mockOutboxRepo struct{}

func (m *mockOutboxRepo) Save(ctx context.Context, record *outbox.OutboxRecord) error { return nil }
func (m *mockOutboxRepo) FindPending(ctx context.Context, limit int) ([]*outbox.OutboxRecord, error) {
	return nil, nil
}
func (m *mockOutboxRepo) MarkAsPublished(ctx context.Context, id string) error { return nil }

type mockUoW struct {
	mu sync.Mutex
}

func (m *mockUoW) Do(ctx context.Context, fn func(ctx context.Context) error) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	return fn(ctx)
}

func TestProcessBet_Concurrency(t *testing.T) {
	initialAmount, _ := money.New(1000, "USD") // 10.00 USD
	w, _ := wallet.NewWallet("wallet-1", "player-1", initialAmount)

	walletRepo := &mockWalletRepo{wallet: w}
	trxRepo := &mockTrxRepo{tx: make(map[string]*transaction.WagerTransaction)}
	uow := &mockUoW{}

	uc := usecase.NewProcessBetUseCase(walletRepo, trxRepo, &mockLedgerRepo{}, &mockOutboxRepo{}, uow)

	const numWorkers = 50
	const betAmount int64 = 50 // 0.50 USD
	var wg sync.WaitGroup

	var processed int32
	var rejected int32

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			input := usecase.ProcessBetInput{
				WalletID:       "wallet-1",
				ProviderID:     "provider-1",
				TransactionID:  fmt.Sprintf("trx-%d", id),
				Amount:         betAmount,
				Currency:       "USD",
				IdempotencyKey: fmt.Sprintf("idem-key-%d", id),
			}

			out, err := uc.Execute(context.Background(), input)
			if err != nil {
				t.Errorf("worker %d failed: %v", id, err)
				return
			}

			if out.State == transaction.StateProcessed {
				atomic.AddInt32(&processed, 1)
			} else if out.State == transaction.StateRejected {
				atomic.AddInt32(&rejected, 1)
			}
		}(i)
	}

	wg.Wait()

	// Initial = 1000, 50 workers * 50 = 2500 total demand
	// Since we only have 1000, we expect exactly 20 successful withdrawals.
	// Balance should be exactly 0, not negative.
	if walletRepo.wallet.Balance.Amount() != 0 {
		t.Errorf("expected balance to be 0, got %d", walletRepo.wallet.Balance.Amount())
	}

	if processed != 20 {
		t.Errorf("expected exactly 20 processed bets, got %d", processed)
	}

	if rejected != 30 {
		t.Errorf("expected exactly 30 rejected bets, got %d", rejected)
	}
}
