package euroexchangerates

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"sort"
	"strconv"

	"github.com/suhlig/concourse-resource-go"
	"github.com/suhlig/euro-exchange-rates-resource/frankfurter"
	"golang.org/x/exp/maps"
)

type ConcourseResource[S Source, V Version, P Params] struct {
	HttpClient *http.Client
}

type Source struct {
	URL string `json:"url" validate:"required,http_url"`
}

type Version struct {
	Date frankfurter.YMD `json:"date" validate:"required"`
}

type Params struct {
	Currencies []string `json:"currencies"`
}

func (r ConcourseResource[S, V, P]) Check(ctx context.Context, request concourse.CheckRequest[Source, Version], log io.Writer) (concourse.CheckResponse[Version], error) {
	service := frankfurter.ExchangeRatesService{URL: request.Source.URL, HttpClient: r.HttpClient}

	var response concourse.CheckResponse[Version]

	if request.Version.Date.IsZero() {
		fmt.Fprintf(log, "Fetching latest exchange rates\n")
		rates, err := service.Latest(ctx)

		if err != nil {
			return nil, fmt.Errorf("unable to fetch latest rate from %s: %w", request.Source.URL, err)
		}

		response = concourse.CheckResponse[Version]{Version{Date: rates.Date}}
	} else {
		fmt.Fprintf(log, "Fetching exchange rates since %s\n", request.Version)
		history, err := service.Since(ctx, request.Version.Date)

		if err != nil {
			return nil, fmt.Errorf("unable to fetch rates since %s from %s: %w", request.Version, request.Source.URL, err)
		}

		for date := range history.Rates {
			response = append(response, Version{Date: date})
		}

		// Cannot use sort.Sort(response) here because
		// a) [T comparable] is not enough for sorting (needs to be [T cmp.Ordered]), and
		// b) implementing Less is not possible in a generic way because only certain built-in types satisfy https://pkg.go.dev/cmp@master#Ordered
		sort.Slice(response, func(i, j int) bool {
			return response[i].Date.After(response[j].Date)
		})
	}

	return response, nil
}

func (r ConcourseResource[S, V, P]) Get(ctx context.Context, request concourse.GetRequest[Source, Version, Params], log io.Writer, destination string) (*concourse.Response[Version], error) {
	fmt.Fprintf(log, "Fetching exchange rates for %s as of %s and placing them in %s\n", request.Params.Currencies, request.Version, destination)

	rates, err := frankfurter.ExchangeRatesService{
		URL:        request.Source.URL,
		HttpClient: r.HttpClient,
	}.At(ctx, request.Version.Date)

	if err != nil {
		return nil, fmt.Errorf("unable to fetch rates as of %s from %s: %w", request.Version.Date, request.Source.URL, err)
	}

	var currencies []frankfurter.Currency

	if len(request.Params.Currencies) == 0 {
		currencies = maps.Keys(rates.Rates)
	} else {
		currencies = mapFunc(request.Params.Currencies, func(c string) frankfurter.Currency {
			return frankfurter.Currency(c)
		})
	}

	for _, c := range currencies {
		rate, found := rates.Rates[c]

		if !found {
			fmt.Fprintf(log, "Warning: requested currency %s is not part of the response. Available curencies are: %v\n", c, maps.Keys(rates.Rates))
			continue
		}

		os.WriteFile(path.Join(destination, string(c)), []byte(rateString(rate)), 0755)
	}

	response := concourse.Response[Version]{
		Version: Version{Date: request.Version.Date},
	}

	for _, c := range currencies {
		response.Metadata = append(response.Metadata, concourse.NameValuePair{Name: string(c), Value: rateString(rates.Rates[c])})
	}

	return &response, nil
}

func (r ConcourseResource[S, V, P]) Put(ctx context.Context, request concourse.PutRequest[Source, Params], log io.Writer, source string) (*concourse.Response[Version], error) {
	fmt.Fprintf(log, "This resource does nothing on put\n")
	return &concourse.Response[Version]{}, nil
}

// https://stackoverflow.com/a/71624929
func mapFunc[T, U any](ts []T, f func(T) U) []U {
	us := make([]U, len(ts))

	for i := range ts {
		us[i] = f(ts[i])
	}

	return us
}

// https://stackoverflow.com/a/40555281
func rateString(rate float32) string {
	return strconv.FormatFloat(float64(rate), 'f', -1, 32)
}
