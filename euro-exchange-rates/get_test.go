package euroexchangerates_test

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/suhlig/concourse-resource-go"
	xr "github.com/suhlig/euro-exchange-rates-resource/euro-exchange-rates"
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

	Context("requested version exists", func() {
		BeforeEach(func() {
			request.Source.URL = server.URL
			responseBody = `
				{
					"amount": 1.0,
					"base": "EUR",
					"date": "2024-01-16",
					"rates": { "SEK": 11.3215, "USD": 1.0882, "THB": 38.522 }
				}
			`
		})

		It("works", func() {
			Expect(err).ToNot(HaveOccurred())
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
