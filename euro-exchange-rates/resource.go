package euroexchangerates

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"sort"
	"strconv"
	"strings"

	"github.com/suhlig/concourse-resource-go"
	"github.com/suhlig/euro-exchange-rates-resource/frankfurter"
	"golang.org/x/exp/maps"
)

type ConcourseResource[S Source, V Version, P Params] struct {
	HttpClient *http.Client
}

type Source struct {
	URL        string                 `json:"url" validate:"required,http_url"`
	Currencies []frankfurter.Currency `json:"currencies"`
	Verbose    bool
}

type Version struct {
	Date frankfurter.YMD `json:"date" validate:"required"`
}

func (v Version) String() string {
	return v.Date.String()
}

type Params struct{}

func (r ConcourseResource[S, V, P]) Check(ctx context.Context, request concourse.CheckRequest[Source, Version], log io.Writer) (concourse.CheckResponse[Version], error) {
	if request.Source.Verbose {
		r.HttpClient.Transport = RequestResponseLogger{Writer: log}
	}

	service := frankfurter.ExchangeRatesService{URL: request.Source.URL, HttpClient: r.HttpClient}

	var response concourse.CheckResponse[Version]

	if request.Version.Date.IsZero() {
		fmt.Fprintf(log, "Fetching latest exchange rates\n")
		rates, err := service.Latest(ctx, request.Source.Currencies...)

		if err != nil {
			return nil, fmt.Errorf("unable to fetch latest rate from %s: %w", request.Source.URL, err)
		}

		response = concourse.CheckResponse[Version]{Version{Date: rates.Date}}
	} else {
		fmt.Fprintf(log, "Fetching exchange rates since %s\n", request.Version)
		history, err := service.Since(ctx, request.Version.Date, request.Source.Currencies...)

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
			return response[i].Date.Before(response[j].Date)
		})
	}

	return response, nil
}

func (r ConcourseResource[S, V, P]) Get(ctx context.Context, request concourse.GetRequest[Source, Version, Params], log io.Writer, destination string) (*concourse.Response[Version], error) {
	if request.Source.Verbose {
		r.HttpClient.Transport = RequestResponseLogger{Writer: log}
	}

	if len(request.Source.Currencies) == 0 {
		fmt.Fprintf(log, "Fetching all exchange rates as of %s and placing them in %s\n", request.Version, destination)
	} else {
		fmt.Fprintf(log, "Fetching exchange rates for %s as of %s and placing them in %s\n", request.Source.Currencies, request.Version, destination)
	}

	rates, err := frankfurter.ExchangeRatesService{
		URL:        request.Source.URL,
		HttpClient: r.HttpClient,
	}.At(ctx, request.Version.Date, request.Source.Currencies...)

	if err != nil {
		return nil, fmt.Errorf("unable to fetch rates as of %s from %s: %w", request.Version.Date, request.Source.URL, err)
	}

	if !request.Version.Date.Equal(rates.Date) {
		return nil, fmt.Errorf("requested version %s is not available; closest is %s", request.Version.Date, rates.Date)
	}

	for c := range rates.Rates {
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

	for c := range rates.Rates {
		response.Metadata = append(response.Metadata, concourse.NameValuePair{Name: string(c), Value: rateString(rates.Rates[c])})
	}

	return &response, nil
}

func (r ConcourseResource[S, V, P]) Put(ctx context.Context, request concourse.PutRequest[Source, Params], log io.Writer, source string) (*concourse.Response[Version], error) {
	fmt.Fprintf(log, "This resource does nothing on put\n")
	return &concourse.Response[Version]{}, nil
}

// https://stackoverflow.com/a/40555281
func rateString(rate float32) string {
	return strconv.FormatFloat(float64(rate), 'f', -1, 32)
}

type RequestResponseLogger struct {
	Writer io.Writer
}

func (t RequestResponseLogger) RoundTrip(req *http.Request) (*http.Response, error) {
	dumpRequest(t.Writer, req)

	resp, err := http.DefaultTransport.RoundTrip(req)

	if err != nil {
		return nil, err
	}

	err = dumpResponse(t.Writer, resp)

	if err != nil {
		return nil, err
	}

	return resp, err
}

func dumpRequest(w io.Writer, req *http.Request) {
	fmt.Fprintf(w, "> %s %s\n", req.Method, req.URL)

	for k, v := range req.Header {
		fmt.Fprintf(w, "> %s: %v\n", k, strings.Join(v, ", "))
	}

	// we don't send a body; no need to log it
}

func dumpResponse(w io.Writer, resp *http.Response) error {
	var responseBody bytes.Buffer
	_, err := responseBody.ReadFrom(resp.Body)

	if err != nil {
		return err
	}

	// preserve the body for downstream reading
	resp.Body = io.NopCloser(bytes.NewReader(responseBody.Bytes()))

	fmt.Fprintf(w, "< %d\n", resp.StatusCode)

	for k, v := range resp.Header {
		fmt.Fprintf(w, "< %s: %v\n", k, strings.Join(v, ", "))
	}

	_, err = io.Copy(w, &responseBody)

	if err != nil {
		return err
	}

	fmt.Fprintln(w)

	return nil
}
