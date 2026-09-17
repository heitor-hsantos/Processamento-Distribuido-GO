package di

import (
	"context"

	"JungleGaming-test/config"
	"JungleGaming-test/internal/application/port"
	"JungleGaming-test/internal/application/usecase"
	"JungleGaming-test/internal/infrastructure/auth"
	"JungleGaming-test/internal/infrastructure/messaging/consumer"
	"JungleGaming-test/internal/infrastructure/messaging/publisher"
	"go.uber.org/fx"
)

func NewOIDCValidator(cfg config.Config) *auth.OIDCValidator {
	return auth.NewOIDCValidator(auth.OIDCConfig{IssuerURL: cfg.IssuerURL, Audience: cfg.Audience})
}

func NewBetConsumer(cfg config.Config, process *usecase.ProcessBetUseCase, inbox port.InboxRepository, uow port.UnitOfWork) (*consumer.BetConsumer, error) {
	return consumer.NewBetConsumer(cfg, process, inbox, uow)
}

func NewOutboxPublisher(cfg config.Config, repo port.OutboxRepository) (*publisher.OutboxPublisher, error) {
	return publisher.NewOutboxPublisher(cfg, repo)
}

func RegisterLifecycle(lc fx.Lifecycle, workers ...interface{}) {
	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			for _, worker := range workers {
				switch w := worker.(type) {
				case *consumer.BetConsumer:
					_ = w.Start(context.Background())
				case *publisher.OutboxPublisher:
					_ = w.Start(context.Background())
				}
			}
			return nil
		},
		OnStop: func(context.Context) error {
			for _, worker := range workers {
				switch w := worker.(type) {
				case *consumer.BetConsumer:
					_ = w.Stop(context.Background())
				case *publisher.OutboxPublisher:
					_ = w.Stop(context.Background())
				}
			}
			return nil
		},
	})
}
