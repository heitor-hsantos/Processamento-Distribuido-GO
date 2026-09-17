package router

import (
	"net/http"

	"JungleGaming-test/internal/infrastructure/auth"
	"JungleGaming-test/internal/infrastructure/http/handler"
	"JungleGaming-test/internal/infrastructure/http/middleware"
)

func NewRouter(walletHandler *handler.WalletHandler, trxHandler *handler.TransactionHandler, healthHandler *handler.HealthHandler, validator *auth.OIDCValidator) http.Handler {
	mux := http.NewServeMux()
	authMid := middleware.Auth(validator)

	if healthHandler == nil {
		healthHandler = &handler.HealthHandler{}
	}

	mux.HandleFunc("/health/live", healthHandler.Live)
	mux.HandleFunc("/health/ready", healthHandler.Ready)

	mux.Handle("/wallets", authMid(middleware.Idempotency(http.HandlerFunc(walletHandler.Create))))
	mux.Handle("/wallets/", authMid(http.HandlerFunc(walletHandler.Get)))
	mux.Handle("/wagering/transactions", authMid(middleware.Idempotency(http.HandlerFunc(trxHandler.Create))))
	mux.Handle("/wagering/transactions/", authMid(http.HandlerFunc(trxHandler.List)))

	return mux
}
