package api

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/shopspring/decimal"
)

func Test_FetchUSDCryptoRates(t *testing.T) {

	fetch_tests := map[string]struct {
		server_response       string
		expected_return_value CoinBaseResponse
		expected_error        bool
	}{
		"/valid_response": {
			server_response: `{
				"data": {
					"currency": "USD",
					"rates": {
						"ETH": "13.5043889264010804",
						"BTC": "0.0010000000000000"
					}
				}
			}`,
			expected_return_value: CoinBaseResponse{
				Data: struct {
					Currency string            `json:"currency"`
					Rates    map[string]string `json:"rates"`
				}{
					Currency: "USD",
					Rates: map[string]string{
						"ETH": "13.5043889264010804",
						"BTC": "0.0010000000000000",
					},
				},
			},
			expected_error: false,
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(fetch_tests[r.URL.Path].server_response))
	}))
	defer server.Close()

	for key, test := range fetch_tests {
		rv, err := FetchUSDCryptoRates(server.URL + key)
		has_err := err != nil
		if has_err != test.expected_error {
			t.Errorf("%s expected error to be %v, got %v", key, test.expected_error, err)
		}
		if !reflect.DeepEqual(rv, test.expected_return_value) {
			t.Errorf("%s, expected %v, got %v", key, test.expected_return_value, rv)
		}
	}
}

func Test_ParseRateFromCoinBaseResponse(t *testing.T) {
	empty_return_value, _ := decimal.NewFromString("")

	parse_tests := map[string]struct {
		exchange_rates               CoinBaseResponse
		rate_key                     string
		expected_return_value_string string
		expected_error               bool
	}{
		"valid_response": {
			exchange_rates: CoinBaseResponse{
				Data: struct {
					Currency string            `json:"currency"`
					Rates    map[string]string `json:"rates"`
				}{
					Currency: "USD",
					Rates: map[string]string{
						"ETH": "13.5043889264010804",
						"BTC": "0.0010000000000000",
					},
				},
			},
			rate_key:                     "ETH",
			expected_return_value_string: "13.5043889264010804",
			expected_error:               false,
		},
		"invalid_rate_key": {
			exchange_rates: CoinBaseResponse{
				Data: struct {
					Currency string            `json:"currency"`
					Rates    map[string]string `json:"rates"`
				}{
					Currency: "USD",
					Rates: map[string]string{
						"ETH": "13.5043889264010804",
						"BTC": "0.0010000000000000",
					},
				},
			},
			rate_key:                     "DOGE",
			expected_return_value_string: empty_return_value.String(),
			expected_error:               true,
		},
		"padded_rate_value": {
			exchange_rates: CoinBaseResponse{
				Data: struct {
					Currency string            `json:"currency"`
					Rates    map[string]string `json:"rates"`
				}{
					Currency: "USD",
					Rates: map[string]string{
						"ETH": "13.5043889264010804",
						"BTC": "000.0010000000000000",
					},
				},
			},
			rate_key:                     "BTC",
			expected_return_value_string: "0.001",
			expected_error:               false,
		},
	}

	for key, test := range parse_tests {
		rv, err := ParseRateFromCoinBaseResponse(test.exchange_rates, test.rate_key)
		has_err := err != nil
		if has_err != test.expected_error {
			t.Errorf("%s expected error %v, got %v", key, test.expected_error, err)
		}
		if rv.String() != test.expected_return_value_string {
			t.Errorf("%s, expected %v, got %v", key, test.expected_return_value_string, rv)
		}
	}
}
