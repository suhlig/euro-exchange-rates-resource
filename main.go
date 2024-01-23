package main

import (
	"net/http"
	"os"

	_ "time/tzdata" // https://pkg.go.dev/time/tzdata#pkg-overview

	"github.com/suhlig/concourse-resource-go"
	xr "github.com/suhlig/euro-exchange-rates-resource/euro-exchange-rates"
)

func main() {
	resource := xr.ConcourseResource[xr.Source, xr.Version, xr.Params]{
		HttpClient: http.DefaultClient,
	}

	if err := concourse.NewRootCommand(&resource, "Euro Exchange Rates resource").Execute(); err != nil {
		os.Exit(1)
	}
}
