package euroexchangerates_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/suhlig/concourse-resource-go"
	xr "github.com/suhlig/euro-exchange-rates-resource/euro-exchange-rates"
)

var _ = Describe("Check", Label("integration"), func() {
	var (
		err      error
		resource concourse.Resource[xr.Source, xr.Version, xr.Params]
		request  concourse.CheckRequest[xr.Source, xr.Version]
		response concourse.CheckResponse[xr.Version]
	)

	BeforeEach(func() {
		resource = xr.ConcourseResource[
			xr.Source,
			xr.Version,
			xr.Params,
		]{}
		request = concourse.CheckRequest[xr.Source, xr.Version]{}
		response = concourse.CheckResponse[xr.Version]{}
	})

	JustBeforeEach(func(ctx SpecContext) {
		err = resource.Check(ctx, request, &response, GinkgoWriter)
	})

	Context("happy day", func() {
		BeforeEach(func() {
			request.Source.URL = "https://api.frankfurter.app"
		})

		It("works", func() {
			Expect(err).ToNot(HaveOccurred())
		})

		It("has a version", func() {
			Expect(response).ToNot(BeEmpty())
		})
	})

	Context("empty source config", func() {
		It("fails", func() {
			Expect(err).To(HaveOccurred())
		})
	})
})
