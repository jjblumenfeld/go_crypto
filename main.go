package main

import (
	"fmt"
	"os"

	"github.com/jjblumenfeld/go/crypto/api"
	"github.com/shopspring/decimal"
)

func getAmountToBuy(USD decimal.Decimal, allocation decimal.Decimal, ratePerDollar decimal.Decimal) (amountToBuy decimal.Decimal, err error) {
	amountToBuy = USD.Mul(allocation).Mul(ratePerDollar)
	return amountToBuy, nil
}

func main() {
	// Get the totalUSD from command line arguments
	if len(os.Args) < 2 {
		fmt.Println("Please provide the value of total USD as a command line argument")
		return
	}

	totalUSDStr := os.Args[1]
	totalUSD, err := decimal.NewFromString(totalUSDStr)
	if err != nil {
		fmt.Println("Invalid value for totalUSD:", totalUSDStr)
		return
	}

	type cryptoCurrency struct {
		name       string
		allocation decimal.Decimal
	}

	// Define the crypto currencies we will buy
	cryptoCurrenciesToBuy := map[string]cryptoCurrency{
		"ETH": {
			name:       "Ethereum",
			allocation: decimal.NewFromFloat(0.3),
		},
		"BTC": {
			name:       "Bitcoin",
			allocation: decimal.NewFromFloat(0.7),
		},
	}

	coinBaseResponse, err := api.FetchUSDCryptoRates(api.CoinBaseURL)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Create a slice to store the name of the currency and how much to buy.
	type currencyToBuy struct {
		name   string
		rate   decimal.Decimal
		amount decimal.Decimal
	}

	var results []currencyToBuy
	for symbol, currency := range cryptoCurrenciesToBuy {
		rate, err := api.ParseRateFromCoinBaseResponse(coinBaseResponse, symbol)
		if err != nil {
			fmt.Println(err)
			return
		}

		amountToBuy, err := getAmountToBuy(totalUSD, currency.allocation, rate)
		if err != nil {
			fmt.Println(err)
			return
		}

		results = append(results, currencyToBuy{
			name:   currency.name,
			rate:   rate,
			amount: amountToBuy,
		})

	}

	fmt.Printf("With $%s, you can buy:\n", totalUSD)
	for _, result := range results {
		fmt.Printf("%s, of %s at the current rate of $1USD = %s %s\n", result.amount, result.name, result.rate, result.name)
	}
}
