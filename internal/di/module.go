package di

import (
	"JungleGaming-test/internal/infrastructure/messaging/consumer"
	"JungleGaming-test/internal/infrastructure/messaging/publisher"
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Module(
		"jungle-gaming",
		fx.Provide(NewConfig),
		fx.Provide(NewDBPool),
		fx.Provide(NewWalletRepository),
		fx.Provide(NewTransactionRepository),
		fx.Provide(NewLedgerRepository),
		fx.Provide(NewOutboxRepository),
		fx.Provide(NewInboxRepository),
		fx.Provide(NewUnitOfWork),
		fx.Provide(NewCreateWalletUseCase),
		fx.Provide(NewProcessBetUseCase),
		fx.Provide(NewReconcileUseCase),
		fx.Provide(NewOIDCValidator),
		fx.Provide(NewBetConsumer),
		fx.Provide(NewOutboxPublisher),
		fx.Provide(NewWalletHandler),
		fx.Provide(NewTransactionHandler),
		fx.Provide(NewHealthHandler),
		fx.Provide(NewHTTPHandler),
		fx.Provide(NewHTTPServer),
		fx.Invoke(RegisterDBLifecycle),
		fx.Invoke(func(lc fx.Lifecycle, bet *consumer.BetConsumer, outbox *publisher.OutboxPublisher) {
			RegisterLifecycle(lc, bet, outbox)
		}),
		fx.Invoke(RegisterServerLifecycle),
	)
}
