package api

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func Test_FetchUSDCryptoRates(t *testing.T) {

	fetch_tests := map[string]struct {
		name                  string
		path                  string
		server_response       string
		expected_return_value CoinBaseResponse
		expected_error        error
	}{
		"/valid_response": {
			path: "/valid_response",
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
			expected_error: nil,
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(fetch_tests[r.URL.Path].server_response))
	}))
	defer server.Close()

	for _, test := range fetch_tests {
		rv, error := FetchUSDCryptoRates(server.URL + test.path)
		if error != nil {
			t.Errorf("Expected no error, got %v", error)
		}
		if !reflect.DeepEqual(rv, test.expected_return_value) {
			t.Errorf("Expected %v, got %v", test.expected_return_value, rv)
		}
	}
}
