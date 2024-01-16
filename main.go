package main

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/homeport/concourse-resource-go"
)

func main() {
	r := Resource[Source, Version, Params]{}

	if err := concourse.NewRootCommand(&r).Execute(); err != nil {
		os.Exit(1)
	}
}

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
	fmt.Fprintf(log, "Fetching recent exchange rates from %s\n", request.Source.URL)

	*response = append(*response, Version{Date: "2024-01-15"})

	return nil
}

func (r Resource[S, V, P]) Get(ctx context.Context, request concourse.GetRequest[Source, Version, Params], response *concourse.Response[Version], log io.Writer, destination string) error {
	fmt.Fprintf(log, "Fetching exchange rates for %s from %s and placing them in %s\n", request.Params.Currencies, request.Source.URL, destination)

	response.Version = Version{Date: "2024-01-15"}
	response.Metadata = append(response.Metadata, concourse.NameValuePair{Name: "currencies available", Value: "42"})

	return nil
}

func (r Resource[S, V, P]) Put(ctx context.Context, request concourse.PutRequest[Source, Params], response *concourse.Response[Version], log io.Writer, source string) error {
	fmt.Fprintf(log, "This resource does nothing on put\n")
	return nil
}
