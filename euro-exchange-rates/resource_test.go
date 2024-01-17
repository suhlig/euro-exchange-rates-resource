package euroexchangerates_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/suhlig/concourse-resource-go"
	xr "github.com/suhlig/euro-exchange-rates-resource/euro-exchange-rates"
)

var _ = Describe("Check", func() {
	var (
		err      error
		server   *httptest.Server
		resource concourse.Resource[xr.Source, xr.Version, xr.Params]
		request  concourse.CheckRequest[xr.Source, xr.Version]
		response concourse.CheckResponse[xr.Version]
	)

	BeforeEach(func() {
		server = httptest.NewServer(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintln(w, `
					{
						"amount": 1.0,
						"base": "EUR",
						"date": "2024-01-16",
						"rates": { "SEK": 11.3215, "USD": 1.0882 }
					}
				`)
			}))

		resource = xr.ConcourseResource[xr.Source, xr.Version, xr.Params]{
			HttpClient: server.Client(),
		}

		request = concourse.CheckRequest[xr.Source, xr.Version]{}
		response = concourse.CheckResponse[xr.Version]{}
	})

	AfterEach(func() {
		server.Close()
	})

	JustBeforeEach(func(ctx SpecContext) {
		err = resource.Check(ctx, request, &response, GinkgoWriter)
	})

	Context("happy day", func() {
		BeforeEach(func() {
			request.Source.URL = server.URL
		})

		It("works", func() {
			Expect(err).ToNot(HaveOccurred())
		})

		It("has exactly one version", func() {
			Expect(response).To(HaveLen(1))
		})

		It("has the expected version", func() {
			Expect(response[0].Date).To(Equal("2024-01-16"))
		})
	})

	Context("empty source config", func() {
		It("fails", func() {
			Expect(err).To(HaveOccurred())
		})
	})
})
