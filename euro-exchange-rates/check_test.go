package euroexchangerates_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/suhlig/concourse-resource-go"
	xr "github.com/suhlig/euro-exchange-rates-resource/euro-exchange-rates"
)

var _ = Describe("Check", func() {
	var (
		err      error
		request  concourse.CheckRequest[xr.Source, xr.Version]
		response concourse.CheckResponse[xr.Version]
	)

	JustBeforeEach(func(ctx SpecContext) {
		response, err = resource.Check(ctx, request, GinkgoWriter)
	})

	Context("no version given", func() {
		BeforeEach(func() {
			request.Source.URL = server.URL
			responseBody = `
				{
					"amount": 1.0,
					"base": "EUR",
					"date": "2024-01-16",
					"rates": { "SEK": 11.3215, "USD": 1.0882 }
				}
			`
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
