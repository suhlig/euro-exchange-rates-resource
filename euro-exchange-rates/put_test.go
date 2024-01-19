package euroexchangerates_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/suhlig/concourse-resource-go"
	xr "github.com/suhlig/euro-exchange-rates-resource/euro-exchange-rates"
)

var _ = Describe("Put", func() {
	var (
		err      error
		request  concourse.PutRequest[xr.Source, xr.Params]
		response concourse.Response[xr.Version]
		inputDir string
	)

	JustBeforeEach(func(ctx SpecContext) {
		err = resource.Put(ctx, request, &response, GinkgoWriter, inputDir)
	})

	It("works", func() {
		Expect(err).ToNot(HaveOccurred())
	})
})
