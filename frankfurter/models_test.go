package frankfurter_test

import (
	"encoding/json"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/suhlig/euro-exchange-rates-resource/frankfurter"
)

func TestFrankfurter(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Frankfurter Suite")
}

var _ = Describe("Micros", func() {
	Describe("unmarshaling", func() {
		DescribeTable("converts decimals to micros",
			func(input string, expected frankfurter.Micros) {
				var m frankfurter.Micros

				Expect(json.Unmarshal([]byte(input), &m)).To(Succeed())
				Expect(m).To(Equal(expected))
			},
			Entry("integer", `1`, frankfurter.Micros(1_000_000)),
			Entry("four decimals", `11.3215`, frankfurter.Micros(11_321_500)),
			Entry("trailing zero", `38.5220`, frankfurter.Micros(38_522_000)),
			Entry("less than one", `0.5`, frankfurter.Micros(500_000)),
			Entry("small fraction", `0.000001`, frankfurter.Micros(1)),
			Entry("rounds up", `0.0000005`, frankfurter.Micros(1)),
			Entry("rounds down", `0.0000004`, frankfurter.Micros(0)),
			Entry("negative", `-1.5`, frankfurter.Micros(-1_500_000)),
		)

		DescribeTable("fails on invalid rates",
			func(input string) {
				var m frankfurter.Micros

				Expect(json.Unmarshal([]byte(input), &m)).ToNot(Succeed())
			},
			Entry("scientific notation", `1e6`),
			Entry("no digits", `.`),
			Entry("letters", `1.2.3`),
			Entry("not a number", `"abc"`),
			Entry("overflow", `99999999999999999999`),
		)
	})

	Describe("String", func() {
		DescribeTable("formats micros as trimmed decimals",
			func(micros frankfurter.Micros, expected string) {
				Expect(micros.String()).To(Equal(expected))
			},
			Entry("integer", frankfurter.Micros(1_000_000), "1"),
			Entry("four decimals", frankfurter.Micros(11_321_500), "11.3215"),
			Entry("trims trailing zeros", frankfurter.Micros(38_522_000), "38.522"),
			Entry("zero", frankfurter.Micros(0), "0"),
			Entry("less than one", frankfurter.Micros(500_000), "0.5"),
			Entry("negative", frankfurter.Micros(-1_500_000), "-1.5"),
		)
	})

	It("unmarshals a Rates map", func() {
		var rates frankfurter.Rates

		Expect(json.Unmarshal([]byte(`{"SEK": 11.3215, "USD": 1.0882}`), &rates)).To(Succeed())

		Expect(rates).To(HaveLen(2))
		Expect(rates[frankfurter.Currency("SEK")]).To(Equal(frankfurter.Micros(11_321_500)))
		Expect(rates[frankfurter.Currency("USD")]).To(Equal(frankfurter.Micros(1_088_200)))
	})
})
