package main

import (
	"os"

	"github.com/suhlig/concourse-resource-go"
	euroexchangerates "github.com/suhlig/euro-exchange-rates-resource/euro-exchange-rates"
)

func main() {
	resource := euroexchangerates.ConcourseResource[
		euroexchangerates.Source,
		euroexchangerates.Version,
		euroexchangerates.Params,
	]{}

	if err := concourse.NewRootCommand(&resource).Execute(); err != nil {
		os.Exit(1)
	}
}
