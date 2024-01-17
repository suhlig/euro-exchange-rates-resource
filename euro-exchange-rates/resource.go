package euroexchangerates

import (
	"context"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/suhlig/concourse-resource-go"
	"github.com/suhlig/euro-exchange-rates-resource/frankfurter"
	"golang.org/x/exp/maps"
)

type ConcourseResource[S Source, V Version, P Params] struct{}

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
	fmt.Fprintf(log, "Fetching most recent exchange rates since %s\n", request.Version)
	rates, err := frankfurter.ExchangeRatesService{URL: request.Source.URL}.Latest()

	if err != nil {
		return fmt.Errorf("unable to fetch rates from %s: %w", request.Source.URL, err)
	}

	*response = append(*response, Version{Date: rates.Date})

	return nil
}

func (r ConcourseResource[S, V, P]) Get(ctx context.Context, request concourse.GetRequest[Source, Version, Params], response *concourse.Response[Version], log io.Writer, destination string) error {
	fmt.Fprintf(log, "Fetching exchange rates for %s as of %s and placing them in %s\n", request.Params.Currencies, request.Version, destination)

	rates, err := frankfurter.ExchangeRatesService{URL: request.Source.URL}.At(request.Version.Date)

	if err != nil {
		return fmt.Errorf("unable to fetch rates as of %s from %s: %w", request.Version.Date, request.Source.URL, err)
	}

	var currencies []string

	if len(request.Params.Currencies) == 0 {
		currencies = maps.Keys(rates.Rates)
	} else {
		currencies = request.Params.Currencies
	}

	for _, c := range currencies {
		rate, found := rates.Rates[c]

		if !found {
			fmt.Fprintf(log, "Warning: requested currency %s is not part of the response. Available curencies are: %v\n", c, maps.Keys(rates.Rates))
			continue
		}

		os.WriteFile(path.Join(destination, c), []byte(fmt.Sprintf("%.3f", rate)), 0755)
	}

	response.Version = Version{Date: request.Version.Date}

	for _, c := range currencies {
		response.Metadata = append(response.Metadata, concourse.NameValuePair{Name: c, Value: fmt.Sprintf("%.3f", rates.Rates[c])})
	}

	return nil
}

func (r ConcourseResource[S, V, P]) Put(ctx context.Context, request concourse.PutRequest[Source, Params], response *concourse.Response[Version], log io.Writer, source string) error {
	fmt.Fprintf(log, "This resource does nothing on put\n")
	return nil
}
