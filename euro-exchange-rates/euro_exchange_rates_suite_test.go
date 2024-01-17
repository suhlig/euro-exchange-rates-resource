package euroexchangerates_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestEuroExchangeRates(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "EuroExchangeRates Suite")
}
