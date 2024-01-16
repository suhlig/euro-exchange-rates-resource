package frankfurter

import (
	"encoding/json"
	"net/http"
	"net/url"
)

type ExchangeRatesService struct {
	URL string
}

type ExchangeRates struct {
	Date   string
	Amount float32
	Base   string
	Rates  map[string]float32
}

func (s ExchangeRatesService) Latest() (*ExchangeRates, error) {
	return s.At("latest")
}

func (s ExchangeRatesService) At(date string) (*ExchangeRates, error) {
	path, err := url.JoinPath(s.URL, date)

	if err != nil {
		return nil, err
	}

	httpResponse, err := http.Get(path)

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
