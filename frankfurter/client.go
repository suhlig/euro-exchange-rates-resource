package frankfurter

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
)

// Latest fetches the latest rates
//
// TODO Reduce network traffic by passing currencies to retrieve.
//
// [API Documentation]: https://www.frankfurter.app/docs/#latest
func (s ExchangeRatesService) Latest(ctx context.Context) (*ExchangeRates, error) {
	return s.At(ctx, YMD{})
}

// At fetches the rates at the given date
//
// If one or more currencies are passed, only those will be fetched.
// If none are passed, _all_ currencies will be fetched.
//
// [API Documentation]: https://www.frankfurter.app/docs/#historical
func (s ExchangeRatesService) At(ctx context.Context, date YMD) (*ExchangeRates, error) {
	var (
		urlWithPath string
		err         error
	)

	if date.IsZero() {
		urlWithPath, err = url.JoinPath(s.URL, "latest")
	} else {
		urlWithPath, err = url.JoinPath(s.URL, date.String())
	}

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

// Since fetches the rates between the given date and now
//
// TODO Reduce network traffic by passing currencies to retrieve.
//
// [API Documentation]: https://www.frankfurter.app/docs/#timeseries
func (s ExchangeRatesService) Since(ctx context.Context, date YMD) (*History, error) {
	urlWithPath, err := url.JoinPath(s.URL, date.String()+"..")

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
