package money

import (
	"math"
	"testing"
)

func TestMoneyAddAndSubtract(t *testing.T) {
	base, err := New(100, BRL)
	if err != nil {
		t.Fatalf("unexpected error creating base money: %v", err)
	}

	bonus, err := New(25, BRL)
	if err != nil {
		t.Fatalf("unexpected error creating bonus money: %v", err)
	}

	sum, err := base.Add(bonus)
	if err != nil {
		t.Fatalf("unexpected error adding money: %v", err)
	}
	if sum.Amount() != 125 {
		t.Fatalf("expected 125, got %d", sum.Amount())
	}

	diff, err := base.Subtract(bonus)
	if err != nil {
		t.Fatalf("unexpected error subtracting money: %v", err)
	}
	if diff.Amount() != 75 {
		t.Fatalf("expected 75, got %d", diff.Amount())
	}
}

func TestMoneyRejectsNegativeAmount(t *testing.T) {
	if _, err := New(-1, BRL); err == nil {
		t.Fatal("expected negative amount validation to fail")
	}
}

func TestMoneyRejectsCurrencyMismatch(t *testing.T) {
	left, _ := New(10, BRL)
	right, _ := New(3, USD)
	if _, err := left.Add(right); err == nil {
		t.Fatal("expected currency mismatch to fail")
	}
}

func TestMoneyOverflowOnAddition(t *testing.T) {
	left, _ := New(math.MaxInt64, BRL)
	right, _ := New(1, BRL)
	if _, err := left.Add(right); err == nil {
		t.Fatal("expected addition overflow to fail")
	}
}

func TestMoneyOverflowOnSubtraction(t *testing.T) {
	left, _ := New(10, BRL)
	right, _ := New(11, BRL)
	if _, err := left.Subtract(right); err == nil {
		t.Fatal("expected subtraction underflow to fail")
	}
}

func TestMoneyParsesDecimalStrings(t *testing.T) {
	value, err := NewFromString("25.00", BRL)
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}
	if value.Amount() != 2500 {
		t.Fatalf("expected 2500 cents, got %d", value.Amount())
	}
	if value.String() != "25.00" {
		t.Fatalf("expected 25.00 string, got %s", value.String())
	}
}

func TestMoneyRejectsScientificNotation(t *testing.T) {
	if _, err := NewFromString("1e2", BRL); err == nil {
		t.Fatal("expected scientific notation to fail")
	}
}
