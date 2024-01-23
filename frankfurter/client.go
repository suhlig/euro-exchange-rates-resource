package frankfurter

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
)

// Latest fetches the latest rates
//
// [API Documentation]: https://www.frankfurter.app/docs/#latest
func (s ExchangeRatesService) Latest(ctx context.Context, currencies ...Currency) (*ExchangeRates, error) {
	return s.At(ctx, YMD{}, currencies...)
}

// At fetches the rates at the given date
//
// If one or more currencies are passed, only those will be fetched.
// If none are passed, _all_ currencies will be fetched.
//
// If the returned version is not the one requested, it fails. While Frankfurter returns the closest rate,
// Concourse explicitly states that the resource must fail if the requested version is not available.
//
// [API Documentation]: https://www.frankfurter.app/docs/#historical
func (s ExchangeRatesService) At(ctx context.Context, date YMD, currencies ...Currency) (*ExchangeRates, error) {
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

	if len(currencies) > 0 {
		query := url.Values{}
		query.Add("to", strings.Join(mapFunc(currencies, func(c Currency) string { return string(c) }), ","))
		urlWithPath = urlWithPath + "?" + query.Encode()
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, urlWithPath, nil)

	if err != nil {
		return nil, err
	}

	request.Header.Set("User-Agent", "Concourse Euro Exchange Rates Resource; https://github.com/suhlig/euro-exchange-rates-resource")

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
// [API Documentation]: https://www.frankfurter.app/docs/#timeseries
func (s ExchangeRatesService) Since(ctx context.Context, date YMD, currencies ...Currency) (*History, error) {
	urlWithPath, err := url.JoinPath(s.URL, date.String()+"..")

	if err != nil {
		return nil, err
	}

	if len(currencies) > 0 {
		query := url.Values{}
		query.Add("to", strings.Join(mapFunc(currencies, func(c Currency) string { return string(c) }), ","))
		urlWithPath = urlWithPath + "?" + query.Encode()
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, urlWithPath, nil)

	if err != nil {
		return nil, err
	}

	request.Header.Set("User-Agent", "Concourse Euro Exchange Rates Resource; https://github.com/suhlig/euro-exchange-rates-resource")

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

// https://stackoverflow.com/a/71624929
func mapFunc[T, U any](ts []T, f func(T) U) []U {
	us := make([]U, len(ts))

	for i := range ts {
		us[i] = f(ts[i])
	}

	return us
}
