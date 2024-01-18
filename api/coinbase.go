package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/shopspring/decimal"
)

const CoinBaseURL = "https://api.coinbase.com/v2/exchange-rates?currency=USD"

type CoinBaseResponse struct {
	Data struct {
		Currency string            `json:"currency"`
		Rates    map[string]string `json:"rates"`
	} `json:"data"`
}

func FetchUSDCryptoRates(URL string) (CoinBaseResponse, error) {
	// Send GET request to Coinbase API
	resp, err := http.Get(URL)
	if err != nil {
		return CoinBaseResponse{}, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return CoinBaseResponse{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var exchangeRates CoinBaseResponse
	err = json.NewDecoder(resp.Body).Decode(&exchangeRates)
	if err != nil {
		return CoinBaseResponse{}, err
	}

	return exchangeRates, nil
}

func ParseRateFromCoinBaseResponse(exchangeRates CoinBaseResponse, rateKey string) (decimal.Decimal, error) {
	emptyRate, _ := decimal.NewFromString("")
	rateStr, ok := exchangeRates.Data.Rates[rateKey]
	if !ok {
		return emptyRate, fmt.Errorf("%s rate not found", rateKey) // Return 0 instead of nil
	}

	rate, err := decimal.NewFromString(rateStr)
	if err != nil {
		return emptyRate, fmt.Errorf("failed to parse %s rate: %v", rateStr, err) // Return 0 instead of nil
	}
	return rate, nil
}
