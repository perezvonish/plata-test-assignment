package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"perezvonish/plata-test-assignment/internal/shared/config"
	"time"

	"perezvonish/plata-test-assignment/internal/domain/currency"
)

type GetRateParams struct {
	FromCurrency currency.Currency
	ToCurrency   currency.Currency
}

type ExchangePrice interface {
	GetRate(ctx context.Context, params GetRateParams) (float64, error)
}

type ExchangePriceImpl struct {
	baseUrl     string
	accessToken string
	httpClient  *http.Client
}

type exchangeResponse struct {
	Success bool               `json:"success"`
	Rates   map[string]float64 `json:"rates"`
	Error   *struct {
		Code string `json:"code"`
		Info string `json:"info"`
	} `json:"error"`
}

type ExchangePriceInitParams struct {
	ExchangeApiConfig *config.ExchangeApiConfig
}

func NewExchangePrice(params ExchangePriceInitParams) ExchangePrice {
	return &ExchangePriceImpl{
		baseUrl:     params.ExchangeApiConfig.URL,
		accessToken: params.ExchangeApiConfig.Key,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (e *ExchangePriceImpl) GetRate(ctx context.Context, params GetRateParams) (float64, error) {
	apiUrl, err := url.Parse(fmt.Sprintf("%s/latest", e.baseUrl))
	if err != nil {
		return 0, fmt.Errorf("invalid base url: %w", err)
	}

	query := apiUrl.Query()
	query.Set("access_key", e.accessToken)
	query.Set("base", string(params.FromCurrency))
	query.Set("symbols", string(params.ToCurrency))
	apiUrl.RawQuery = query.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiUrl.String(), nil)

	if err != nil {
		return 0, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := e.httpClient.Do(req)
	if err != nil {
		return 0, fmt.Errorf("http request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("external api returned status: %d", resp.StatusCode)
	}

	var result exchangeResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, fmt.Errorf("failed to decode response: %w", err)
	}

	if !result.Success {
		return 0, fmt.Errorf("api error: %s - %s", result.Error.Code, result.Error.Info)
	}

	rate, ok := result.Rates[string(params.ToCurrency)]
	if !ok {
		return 0, fmt.Errorf("rate for %s not found in response", params.ToCurrency)
	}

	return rate, nil
}
