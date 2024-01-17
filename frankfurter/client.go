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

type Rates map[Currency]float32
type Currency string
type YMD string

type ExchangeRates struct {
	Date   YMD
	Amount float32
	Base   Currency
	Rates  Rates
}

type History struct {
	Amount float32
	Base   Currency
	Start  YMD `json:"start_date"`
	End    YMD `json:"end_date"`
	Rates  map[YMD]Rates
}

func (s ExchangeRatesService) Latest(ctx context.Context) (*ExchangeRates, error) {
	return s.At(ctx, "latest")
}

func (s ExchangeRatesService) At(ctx context.Context, date string) (*ExchangeRates, error) {
	urlWithPath, err := url.JoinPath(s.URL, date)

	if err != nil {
		return nil, err
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, urlWithPath, nil)

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

func (s ExchangeRatesService) Since(ctx context.Context, date string) (*History, error) {
	urlWithPath, err := url.JoinPath(s.URL, date+"..")

	if err != nil {
		return nil, err
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, urlWithPath, nil)

	if err != nil {
		return nil, err
	}

	httpResponse, err := s.HttpClient.Do(request)

	if err != nil {
		return nil, err
	}

	var history History

	err = json.NewDecoder(httpResponse.Body).Decode(&history)

	if err != nil {
		return nil, err
	}

	return &history, nil
}
