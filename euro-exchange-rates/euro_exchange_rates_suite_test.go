package euroexchangerates_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/suhlig/concourse-resource-go"
	xr "github.com/suhlig/euro-exchange-rates-resource/euro-exchange-rates"
)

func TestEuroExchangeRates(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "EuroExchangeRates Suite")
}

var (
	server       *httptest.Server
	resource     concourse.Resource[xr.Source, xr.Version, xr.Params]
	responseBody string
)

var _ = BeforeEach(func() {
	server = httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, responseBody)
		}))

	resource = xr.ConcourseResource[xr.Source, xr.Version, xr.Params]{
		HttpClient: server.Client(),
	}
})

var _ = AfterEach(func() {
	responseBody = "" // make sure we are not re-using it
	server.Close()
})
