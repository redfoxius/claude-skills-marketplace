package gateway

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type OpenExchangeRatesClient struct {
	apiKey string
	http   *http.Client
}

func NewOpenExchangeRatesClient(apiKey string) *OpenExchangeRatesClient {
	return &OpenExchangeRatesClient{apiKey: apiKey, http: &http.Client{}}
}

func (c *OpenExchangeRatesClient) FetchRate(base, quote string) (float64, error) {
	url := fmt.Sprintf("https://openexchangerates.org/api/latest.json?app_id=%s&base=%s", c.apiKey, base)

	resp, err := c.http.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var payload struct {
		Rates map[string]float64 `json:"rates"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return 0, err
	}

	rate, ok := payload.Rates[quote]
	if !ok {
		return 0, fmt.Errorf("no rate for %s", quote)
	}
	return rate, nil
}
