package frankfurter

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
)

type ExchangeRatesService struct {
	HttpClient *http.Client
	URL        string
}

type ExchangeRates struct {
	Date   string
	Amount float32
	Base   string
	Rates  map[string]float32
}

func (s ExchangeRatesService) Latest(ctx context.Context) (*ExchangeRates, error) {
	return s.At(ctx, "latest")
}

func (s ExchangeRatesService) At(ctx context.Context, date string) (*ExchangeRates, error) {
	path, err := url.JoinPath(s.URL, date)

	if err != nil {
		return nil, err
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, path, nil)

	if err != nil {
		return nil, err
	}

	httpResponse, err := s.HttpClient.Do(request)

	if err != nil {
		return nil, err
	}

	var rates ExchangeRates

	err = json.NewDecoder(httpResponse.Body).Decode(&rates)

	if err != nil {
		return nil, err
	}

	return &rates, nil
}
