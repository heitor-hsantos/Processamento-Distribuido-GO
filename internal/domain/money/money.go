package money

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	domainerrors "JungleGaming-test/internal/domain/errors"
)

// Currency represents an ISO 4217 code in the major currency registry used by the domain.
type Currency string

const (
	BRL Currency = "BRL"
	USD Currency = "USD"
)

// Money is an immutable value object that represents an amount in minor units.
type Money struct {
	amount   int64
	currency Currency
}

func New(amount int64, currency Currency) (Money, error) {
	if err := validateCurrency(currency); err != nil {
		return Money{}, err
	}
	if amount < 0 {
		return Money{}, fmt.Errorf("money: %w", domainerrors.ErrInvalidAmount)
	}
	return Money{amount: amount, currency: currency}, nil
}

func NewFromString(value string, currency Currency) (Money, error) {
	if err := validateCurrency(currency); err != nil {
		return Money{}, err
	}
	trimmed := strings.TrimSpace(value)
	if trimmed == "" || trimmed == "." || strings.ContainsAny(trimmed, "nNiaA") {
		return Money{}, fmt.Errorf("money: %w", domainerrors.ErrInvalidAmount)
	}
	if strings.Contains(trimmed, "e") || strings.Contains(trimmed, "E") {
		return Money{}, fmt.Errorf("money: %w", domainerrors.ErrInvalidAmount)
	}
	if strings.Count(trimmed, ".") > 1 || strings.Count(trimmed, ",") > 1 {
		return Money{}, fmt.Errorf("money: %w", domainerrors.ErrInvalidAmount)
	}
	if strings.Contains(trimmed, ",") {
		trimmed = strings.ReplaceAll(trimmed, ",", ".")
	}
	parts := strings.Split(trimmed, ".")
	whole := parts[0]
	fraction := ""
	if len(parts) == 2 {
		fraction = parts[1]
	}
	if whole == "" {
		whole = "0"
	}
	if whole[0] == '+' {
		whole = whole[1:]
	}
	if strings.HasPrefix(whole, "-") {
		return Money{}, fmt.Errorf("money: %w", domainerrors.ErrInvalidAmount)
	}
	if len(fraction) > 2 {
		return Money{}, fmt.Errorf("money: %w", domainerrors.ErrInvalidAmount)
	}
	if len(fraction) == 1 {
		fraction += "0"
	}
	if len(fraction) == 0 {
		fraction = "00"
	}
	base, err := strconv.ParseInt(whole, 10, 64)
	if err != nil {
		return Money{}, fmt.Errorf("money: %w", domainerrors.ErrInvalidAmount)
	}
	frac, err := strconv.ParseInt(fraction, 10, 64)
	if err != nil {
		return Money{}, fmt.Errorf("money: %w", domainerrors.ErrInvalidAmount)
	}
	minor := base*100 + frac
	if minor < 0 {
		return Money{}, fmt.Errorf("money: %w", domainerrors.ErrInvalidAmount)
	}
	return Money{amount: minor, currency: currency}, nil
}

func Must(amount int64, currency Currency) Money {
	m, err := New(amount, currency)
	if err != nil {
		panic(err)
	}
	return m
}

func (m Money) Amount() int64      { return m.amount }
func (m Money) Currency() Currency { return m.currency }
func (m Money) IsZero() bool       { return m.amount == 0 }
func (m Money) Equals(other Money) bool {
	return m.amount == other.amount && m.currency == other.currency
}

func (m Money) String() string {
	if m.amount == 0 {
		return "0.00"
	}
	negative := m.amount < 0
	abs := m.amount
	if negative {
		abs = -m.amount
	}
	whole := abs / 100
	frac := abs % 100
	if negative {
		return fmt.Sprintf("-%d.%02d", whole, frac)
	}
	return fmt.Sprintf("%d.%02d", whole, frac)
}

func (m Money) Add(other Money) (Money, error) {
	if err := validateCurrencyPair(m.currency, other.currency); err != nil {
		return Money{}, err
	}
	if m.amount > math.MaxInt64-other.amount {
		return Money{}, fmt.Errorf("money: %w", domainerrors.ErrInvalidAmount)
	}
	return Money{amount: m.amount + other.amount, currency: m.currency}, nil
}

func (m Money) Subtract(other Money) (Money, error) {
	if err := validateCurrencyPair(m.currency, other.currency); err != nil {
		return Money{}, err
	}
	if other.amount > m.amount {
		return Money{}, fmt.Errorf("money: %w", domainerrors.ErrInsufficientFunds)
	}
	return Money{amount: m.amount - other.amount, currency: m.currency}, nil
}

func (m Money) Negate() Money {
	return Money{amount: -m.amount, currency: m.currency}
}

func validateCurrency(currency Currency) error {
	trimmed := strings.TrimSpace(string(currency))
	if trimmed == "" || len(trimmed) != 3 {
		return fmt.Errorf("money: %w", domainerrors.ErrInvalidCurrency)
	}
	for _, r := range trimmed {
		if r < 'A' || r > 'Z' {
			return fmt.Errorf("money: %w", domainerrors.ErrInvalidCurrency)
		}
	}
	return nil
}

func validateCurrencyPair(a, b Currency) error {
	if err := validateCurrency(a); err != nil {
		return err
	}
	if err := validateCurrency(b); err != nil {
		return err
	}
	if a != b {
		return fmt.Errorf("money: %w", domainerrors.ErrInvalidCurrency)
	}
	return nil
}
