package euroexchangerates

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"sort"

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
	Date string `json:"date" validate:"required,datetime=2006-01-02"`
}

type Params struct {
	Currencies []string `json:"currencies"`
}

func (r ConcourseResource[S, V, P]) Name() string {
	return "Euro Exchange Rates"
}

func (r ConcourseResource[S, V, P]) Check(ctx context.Context, request concourse.CheckRequest[Source, Version], response *concourse.CheckResponse[Version], log io.Writer) error {
	service := frankfurter.ExchangeRatesService{URL: request.Source.URL, HttpClient: r.HttpClient}

	if request.Version.Date == "" {
		fmt.Fprintf(log, "Fetching latest exchange rates\n")
		rates, err := service.Latest(ctx)

		if err != nil {
			return fmt.Errorf("unable to fetch latest rate from %s: %w", request.Source.URL, err)
		}

		*response = concourse.CheckResponse[Version]{Version{Date: string(rates.Date)}}
	} else {
		fmt.Fprintf(log, "Fetching exchange rates since %s\n", request.Version)
		history, err := service.Since(ctx, request.Version.Date)

		if err != nil {
			return fmt.Errorf("unable to fetch rates since %s from %s: %w", request.Version, request.Source.URL, err)
		}

		for date, _ := range history.Rates {
			*response = append(*response, Version{Date: string(date)})
		}

		// Cannot use sort.Sort(response) here because
		// a) T comparable is not enough for sorting (needs to be cmp.Ordered), and
		// b) implementing Less is not possible in a generic way because only certain built-in types satisfy https://pkg.go.dev/cmp@master#Ordered
		sort.Slice(*response, func(i, j int) bool {
			return (*response)[i].Date > (*response)[j].Date
		})
	}

	return nil
}

func (r ConcourseResource[S, V, P]) Get(ctx context.Context, request concourse.GetRequest[Source, Version, Params], response *concourse.Response[Version], log io.Writer, destination string) error {
	fmt.Fprintf(log, "Fetching exchange rates for %s as of %s and placing them in %s\n", request.Params.Currencies, request.Version, destination)

	rates, err := frankfurter.ExchangeRatesService{
		URL:        request.Source.URL,
		HttpClient: r.HttpClient,
	}.At(ctx, request.Version.Date)

	if err != nil {
		return fmt.Errorf("unable to fetch rates as of %s from %s: %w", request.Version.Date, request.Source.URL, err)
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

		os.WriteFile(path.Join(destination, string(c)), []byte(fmt.Sprintf("%.3f", rate)), 0755)
	}

	response.Version = Version{Date: request.Version.Date}

	for _, c := range currencies {
		response.Metadata = append(response.Metadata, concourse.NameValuePair{Name: string(c), Value: fmt.Sprintf("%.3f", rates.Rates[c])})
	}

	return nil
}

func (r ConcourseResource[S, V, P]) Put(ctx context.Context, request concourse.PutRequest[Source, Params], response *concourse.Response[Version], log io.Writer, source string) error {
	fmt.Fprintf(log, "This resource does nothing on put\n")
	return nil
}

// https://stackoverflow.com/a/71624929
func mapFunc[T, U any](ts []T, f func(T) U) []U {
	us := make([]U, len(ts))

	for i := range ts {
		us[i] = f(ts[i])
	}

	return us
}
