package euroexchangerates_test

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/suhlig/concourse-resource-go"
	xr "github.com/suhlig/euro-exchange-rates-resource/euro-exchange-rates"
	"github.com/suhlig/euro-exchange-rates-resource/frankfurter"
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
			request.Source.Currencies = []frankfurter.Currency{
				frankfurter.Currency("SEK"),
				frankfurter.Currency("USD"),
			}
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

		It("has the expected request path", func() {
			Expect(requestURL.Path).To(Equal("/latest"))
		})

		It("has the expected request query", func() {
			Expect(requestURL.Query().Get("to")).To(Equal("SEK,USD"))
		})

		It("has exactly one version", func() {
			Expect(response).To(HaveLen(1))
		})

		It("has the expected version", func() {
			Expect(
				time.Time(response[0].Date),
			).To(BeTemporally("==",
				time.Date(2024, 1, 16, 16, 0, 0, 0, time.FixedZone("Europe/Frankfurt", 60*60)),
				time.Minute))
		})
	})

	Context("version given", func() {
		BeforeEach(func() {
			request.Source.URL = server.URL

			midJanuary, e := frankfurter.NewYMD("2024-01-15")
			Expect(e).ToNot(HaveOccurred())
			request.Version = xr.Version{Date: midJanuary}

			responseBody = `
				{
					"amount": 1.0,
					"base": "EUR",
					"start_date": "2024-01-15",
					"end_date": "2024-01-17",
					"rates": {
						"2024-01-15": {
							"FOO": 7.4575,
							"BAR": 1.0887
						},
						"2024-01-16": {
							"FOO": 7.4585,
							"BAR": 1.089
						},
						"2024-01-17": {
							"FOO": 7.4574,
							"BAR": 1.0872
						}
					}
				}
			`
		})

		Context("currencies configured", func() {
			BeforeEach(func() {
				request.Source.Currencies = []frankfurter.Currency{
					frankfurter.Currency("FOO"),
					frankfurter.Currency("BAR"),
				}
			})

			It("has the expected request query", func() {
				Expect(requestURL.Query().Get("to")).To(Equal("FOO,BAR"))
			})
		})

		It("works", func() {
			Expect(err).ToNot(HaveOccurred())
		})

		It("has the expected request path", func() {
			Expect(requestURL.Path).To(Equal("/2024-01-15.."))
		})

		It("has three versions", func() {
			Expect(response).To(HaveLen(3))
		})

		Context("latest version", func() {
			var latest xr.Version

			BeforeEach(func() {
				latest = response[0]
			})

			It("has the expected version", func() {
				Expect(
					time.Time(latest.Date),
				).To(BeTemporally("==",
					time.Date(2024, 1, 15, 16, 0, 0, 0, time.FixedZone("Europe/Frankfurt", 60*60)),
					time.Minute))
			})
		})

		Context("oldest version", func() {
			var oldest xr.Version

			BeforeEach(func() {
				oldest = response[len(response)-1]
			})

			It("has the expected version", func() {
				Expect(
					time.Time(oldest.Date),
				).To(BeTemporally("==",
					time.Date(2024, 1, 17, 16, 0, 0, 0, time.FixedZone("Europe/Frankfurt", 60*60)),
					time.Minute))
			})
		})
	})

	Context("empty source config", func() {
		It("fails", func() {
			Expect(err).To(HaveOccurred())
		})
	})
})
