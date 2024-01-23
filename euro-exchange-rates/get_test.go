package euroexchangerates_test

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/suhlig/concourse-resource-go"
	xr "github.com/suhlig/euro-exchange-rates-resource/euro-exchange-rates"
	"github.com/suhlig/euro-exchange-rates-resource/frankfurter"
)

var _ = Describe("Get", func() {
	var (
		err      error
		request  concourse.GetRequest[xr.Source, xr.Version, xr.Params]
		response *concourse.Response[xr.Version]
		inputDir string
	)

	BeforeEach(func() {
		inputDir = GinkgoT().TempDir()
	})

	JustBeforeEach(func(ctx SpecContext) {
		response, err = resource.Get(ctx, request, GinkgoWriter, inputDir)
	})

	Context("requested version does not exist", func() {
		BeforeEach(func() {
			request.Source.URL = server.URL

			beforeEcbEvenExisted, e := frankfurter.NewYMD("1998-05-30")
			Expect(e).ToNot(HaveOccurred())
			request.Version = xr.Version{Date: beforeEcbEvenExisted}

			responseBody = `
				{
					"amount": 1.0,
					"base": "EUR",
					"date": "2024-01-15",
					"rates": { "SEK": 11.3215, "USD": 1.0882, "THB": 38.522 }
				}
			`
		})

		It("fails", func() {
			Expect(err).To(HaveOccurred())
		})

		It("has a useful error message", func() {
			Expect(err).To(MatchError(ContainSubstring("not available")))
		})
	})

	Context("requested version exists", func() {
		BeforeEach(func() {
			request.Source.URL = server.URL

			midJanuary, e := frankfurter.NewYMD("2024-01-15")
			Expect(e).ToNot(HaveOccurred())
			request.Version = xr.Version{Date: midJanuary}

			responseBody = `
				{
					"amount": 1.0,
					"base": "EUR",
					"date": "2024-01-15",
					"rates": { "SEK": 11.3215, "USD": 1.0882, "THB": 38.522 }
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
			Expect(requestURL.Path).To(Equal("/2024-01-15"))
		})

		Context("SEK currency requested", func() {
			var currencyFile string

			BeforeEach(func() {
				currencyFile = filepath.Join(inputDir, "SEK")
			})

			It("creates the file", func() {
				Expect(currencyFile).To(BeAnExistingFile())
			})

			It("writes the expected content", func() {
				content, err := os.ReadFile(currencyFile)
				Expect(err).ToNot(HaveOccurred())

				Expect(string(content)).To(Equal("11.3215"))
			})
		})

		Context("USD currency requested", func() {
			var currencyFile string

			BeforeEach(func() {
				currencyFile = filepath.Join(inputDir, "USD")
			})

			It("creates the file", func() {
				Expect(currencyFile).To(BeAnExistingFile())
			})

			It("writes the expected content", func() {
				content, err := os.ReadFile(currencyFile)
				Expect(err).ToNot(HaveOccurred())

				Expect(string(content)).To(Equal("1.0882"))
			})
		})

		Context("THB currency requested", func() {
			var currencyFile string

			BeforeEach(func() {
				currencyFile = filepath.Join(inputDir, "THB")
			})

			It("creates the file", func() {
				Expect(currencyFile).To(BeAnExistingFile())
			})

			It("writes the expected content", func() {
				content, err := os.ReadFile(currencyFile)
				Expect(err).ToNot(HaveOccurred())

				Expect(string(content)).To(Equal("38.522"))
			})
		})

		Context("Metadata", func() {
			It("is not empty", func() {
				Expect(response.Metadata).ToNot(BeEmpty())
			})

			It("has the expected length", func() {
				Expect(response.Metadata).To(HaveLen(3))
			})
		})
	})
})
