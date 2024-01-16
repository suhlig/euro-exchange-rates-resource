package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/homeport/concourse-resource-go"
	"github.com/homeport/euro-exchange-rates-resource/frankfurter"
	"golang.org/x/exp/maps"
)

func main() {
	r := Resource[Source, Version, Params]{}

	if err := concourse.NewRootCommand(&r).Execute(); err != nil {
		os.Exit(1)
	}
}

// TODO The Resource type does not have any members. Do we really need it?
type Resource[S Source, V Version, P Params] struct{}

type Source struct {
	URL string `json:"url"`
}

type Version struct {
	Date string `json:"date"`
}

type Params struct {
	Currencies []string `json:"currencies"`
}

func (r Resource[S, V, P]) Check(ctx context.Context, request concourse.Request[Source], response *[]Version, log io.Writer) error {
	fmt.Fprintf(log, "Fetching most recent exchange rates\n")
	rates, err := frankfurter.ExchangeRatesService{URL: request.Source.URL}.Latest()

	if err != nil {
		return fmt.Errorf("unable to fetch rates from %s: %w", request.Source.URL, err)
	}

	*response = append(*response, Version{Date: rates.Date})

	return nil
}

func (r Resource[S, V, P]) Get(ctx context.Context, request concourse.GetRequest[Source, Version, Params], response *concourse.Response[Version], log io.Writer, destination string) error {
	fmt.Fprintf(log, "Fetching exchange rates for %s as of %s and placing them in %s\n", request.Params.Currencies, request.Version.Date, destination)

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
	response.Metadata = append(response.Metadata, concourse.NameValuePair{Name: "currencies available", Value: fmt.Sprintf("%d", len(rates.Rates))})

	return nil
}

func (r Resource[S, V, P]) Put(ctx context.Context, request concourse.PutRequest[Source, Params], response *concourse.Response[Version], log io.Writer, source string) error {
	fmt.Fprintf(log, "This resource does nothing on put\n")
	return nil
}
