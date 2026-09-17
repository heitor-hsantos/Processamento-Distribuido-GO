package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"JungleGaming-test/internal/application/port"
	"JungleGaming-test/internal/application/usecase"
	domainerrors "JungleGaming-test/internal/domain/errors"
	"JungleGaming-test/internal/domain/money"
)

type WalletHandler struct {
	CreateWallet *usecase.CreateWalletUseCase
	WalletRepo   port.WalletRepository
}

type createWalletRequest struct {
	PlayerID       string            `json:"playerId"`
	InitialBalance *moneyInputAmount `json:"initialBalance"`
}

type moneyInputAmount struct {
	Amount   string `json:"amount"`
	Currency string `json:"currency"`
}

func (h *WalletHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createWalletRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid wallet payload", http.StatusBadRequest)
		return
	}
	if req.PlayerID == "" {
		http.Error(w, "playerId is required", http.StatusBadRequest)
		return
	}
	initial := money.Money{}
	if req.InitialBalance != nil {
		parsedCur := money.Currency(req.InitialBalance.Currency)
		m, err := money.NewFromString(req.InitialBalance.Amount, parsedCur)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		initial = m
	}
	if initial.Currency() == "" {
		initial, _ = money.New(0, money.BRL)
	}
	resp, err := h.CreateWallet.Execute(r.Context(), usecase.CreateWalletInput{
		PlayerID:       req.PlayerID,
		Currency:       initial.Currency(),
		InitialBalance: initial,
	})
	if err != nil {
		if errors.Is(err, domainerrors.ErrInvalidOperation) || errors.Is(err, domainerrors.ErrInvalidCurrency) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	balanceMoney, err := money.New(resp.Balance, resp.Currency)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"id":       resp.WalletID,
		"playerId": req.PlayerID,
		"balance":  map[string]string{"amount": balanceMoney.String(), "currency": string(resp.Currency)},
		"version":  resp.Version,
	})
}

func (h *WalletHandler) Get(w http.ResponseWriter, r *http.Request) {
	walletID := strings.TrimPrefix(r.URL.Path, "/wallets/")
	if walletID == "" || walletID == "/wallets" {
		http.Error(w, "wallet id is required", http.StatusBadRequest)
		return
	}
	if h.WalletRepo == nil {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{"id": walletID, "status": "ok"})
		return
	}
	wallet, err := h.WalletRepo.FindByID(r.Context(), walletID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"id":       wallet.ID,
		"playerId": wallet.PlayerID,
		"balance":  map[string]string{"amount": wallet.Balance.String(), "currency": string(wallet.Currency)},
		"version":  wallet.Version,
	})
}
