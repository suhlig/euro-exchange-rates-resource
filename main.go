package main

import (
	"net/http"
	"os"

	"github.com/suhlig/concourse-resource-go"
	xr "github.com/suhlig/euro-exchange-rates-resource/euro-exchange-rates"
)

func main() {
	resource := xr.ConcourseResource[
		xr.Source,
		xr.Version,
		xr.Params,
	]{
		HttpClient: http.DefaultClient,
	}

	if err := concourse.NewRootCommand(&resource).Execute(); err != nil {
		os.Exit(1)
	}
}
