package di

import (
	"context"
	"net/http"

	"JungleGaming-test/config"
	"JungleGaming-test/internal/application/port"
	"JungleGaming-test/internal/application/usecase"
	"JungleGaming-test/internal/infrastructure/auth"
	"JungleGaming-test/internal/infrastructure/database/repository"
	"JungleGaming-test/internal/infrastructure/database/transaction"
	infrahandler "JungleGaming-test/internal/infrastructure/http/handler"
	infrarouter "JungleGaming-test/internal/infrastructure/http/router"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"
)

func NewConfig() config.Config { return config.Load() }

func NewDBPool(cfg config.Config) (*pgxpool.Pool, error) {
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, cfg.PostgresURL)
	if err != nil {
		return nil, err
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}
	return pool, nil
}

func RegisterDBLifecycle(lc fx.Lifecycle, pool *pgxpool.Pool) {
	lc.Append(fx.Hook{
		OnStop: func(context.Context) error {
			if pool != nil {
				pool.Close()
			}
			return nil
		},
	})
}

func NewWalletRepository(db *pgxpool.Pool) port.WalletRepository {
	return repository.NewWalletRepository(db)
}

func NewTransactionRepository(db *pgxpool.Pool) port.TransactionRepository {
	return repository.NewTransactionRepository(db)
}

func NewLedgerRepository(db *pgxpool.Pool) port.LedgerRepository {
	return repository.NewLedgerRepository(db)
}

func NewOutboxRepository(db *pgxpool.Pool) port.OutboxRepository {
	return repository.NewOutboxRepository(db)
}

func NewInboxRepository(db *pgxpool.Pool) port.InboxRepository {
	return repository.NewInboxRepository(db)
}

func NewUnitOfWork(db *pgxpool.Pool) port.UnitOfWork { return transaction.NewUnitOfWork(db) }

func NewCreateWalletUseCase(repo port.WalletRepository) *usecase.CreateWalletUseCase {
	return usecase.NewCreateWalletUseCase(repo)
}

func NewProcessBetUseCase(
	walletRepo port.WalletRepository,
	trxRepo port.TransactionRepository,
	ledgerRepo port.LedgerRepository,
	outboxRepo port.OutboxRepository,
	uow port.UnitOfWork,
) *usecase.ProcessBetUseCase {
	return usecase.NewProcessBetUseCase(walletRepo, trxRepo, ledgerRepo, outboxRepo, uow)
}

func NewReconcileUseCase(repo port.WalletRepository) *usecase.ReconcileUseCase {
	return usecase.NewReconcileUseCase(repo)
}

func NewWalletHandler(create *usecase.CreateWalletUseCase, repo port.WalletRepository) *infrahandler.WalletHandler {
	return &infrahandler.WalletHandler{CreateWallet: create, WalletRepo: repo}
}

func NewTransactionHandler(process *usecase.ProcessBetUseCase) *infrahandler.TransactionHandler {
	return &infrahandler.TransactionHandler{ProcessBet: process}
}

func NewHealthHandler() *infrahandler.HealthHandler {
	return &infrahandler.HealthHandler{}
}

func NewHTTPHandler(walletHandler *infrahandler.WalletHandler, trxHandler *infrahandler.TransactionHandler, healthHandler *infrahandler.HealthHandler, validator *auth.OIDCValidator) http.Handler {
	return infrarouter.NewRouter(walletHandler, trxHandler, healthHandler, validator)
}

func NewHTTPServer(cfg config.Config, handler http.Handler) *http.Server {
	srv := &http.Server{
		Addr:    ":8080",
		Handler: handler,
	}
	return srv
}

func RegisterServerLifecycle(lc fx.Lifecycle, srv *http.Server) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				_ = srv.ListenAndServe()
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return srv.Shutdown(ctx)
		},
	})
}
