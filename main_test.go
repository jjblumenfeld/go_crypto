package main

import (
	"testing"

	"github.com/shopspring/decimal"
)

func Test_getAmountToBuy(t *testing.T) {
	USD := decimal.NewFromInt(0)
	allocation := decimal.NewFromFloat(0.5)
	ratePerDollar := decimal.NewFromFloat(0.8)

	expectedAmountToBuy := decimal.NewFromFloat(0.0)

	amountToBuy, err := getAmountToBuy(USD, allocation, ratePerDollar)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if !amountToBuy.Equal(expectedAmountToBuy) {
		t.Errorf("Expected amountToBuy to be %s, but got %s", expectedAmountToBuy.String(), amountToBuy.String())
	}
}
