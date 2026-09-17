package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"JungleGaming-test/internal/application/usecase"
	domainerrors "JungleGaming-test/internal/domain/errors"
	"JungleGaming-test/internal/domain/money"
)

type TransactionHandler struct {
	ProcessBet *usecase.ProcessBetUseCase
}

type betRequest struct {
	TransactionID  string `json:"transactionId"`
	WalletID       string `json:"walletId"`
	ProviderID     string `json:"providerId"`
	Amount         int64  `json:"amount"`
	Currency       string `json:"currency"`
	IdempotencyKey string `json:"idempotencyKey"`
}

func (h *TransactionHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req betRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	if req.TransactionID == "" || req.WalletID == "" || req.ProviderID == "" {
		http.Error(w, "transactionId, walletId and providerId are required", http.StatusBadRequest)
		return
	}
	if req.IdempotencyKey == "" {
		req.IdempotencyKey = req.ProviderID + ":" + req.TransactionID
	}
	if req.Currency == "" {
		req.Currency = string(money.BRL)
	}

	input := usecase.ProcessBetInput{
		TransactionID:  req.TransactionID,
		WalletID:       req.WalletID,
		ProviderID:     req.ProviderID,
		Amount:         req.Amount,
		Currency:       money.Currency(req.Currency),
		IdempotencyKey: req.IdempotencyKey,
	}

	out, err := h.ProcessBet.Execute(r.Context(), input)
	if err != nil {
		if errors.Is(err, domainerrors.ErrInsufficientFunds) || errors.Is(err, domainerrors.ErrInvalidAmount) || errors.Is(err, domainerrors.ErrInvalidCurrency) || errors.Is(err, domainerrors.ErrInvalidOperation) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"transactionId": out.TransactionID,
		"state":         out.State,
	})
}

func (h *TransactionHandler) List(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/wagering/transactions/")
	if path == "" {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{"transactionId": path, "status": "ok"})
}
