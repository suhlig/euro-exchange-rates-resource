package frankfurter

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type ExchangeRatesService struct {
	HttpClient *http.Client
	URL        string
}

type ExchangeRates struct {
	Date   YMD
	Amount float32
	Base   Currency
	Rates  Rates
}

type History struct {
	Amount float32
	Base   Currency
	Start  YMD     `json:"start_date"`
	End    YMD     `json:"end_date"`
	Rates  RatesAt `json:"rates"`
}

type Rates map[Currency]float32
type RatesAt map[YMD]Rates
type Currency string

// UnmarshalJSON provides custom unmarshaling as we cannot naiively unmarshal a map with time.Time keys.
func (ra *RatesAt) UnmarshalJSON(raw []byte) error {
	var rates map[string]Rates

	err := json.Unmarshal(raw, &rates)

	if err != nil {
		return fmt.Errorf("could not unmarshal rates: %w", err)
	}

	if *ra == nil {
		*ra = make(RatesAt)
	}

	for k, v := range rates {
		ymd, err := NewYMD(k)

		if err != nil {
			return err
		}

		(*ra)[ymd] = v
	}

	return nil
}

// YMD specializes time.Time encoded as YYYY-MM-DD
// The time zone is hard-coded to Europe/Berlin, which is the same as Frankfurt.
// The time is hard-coded to 16:00 because this is what the ECB specifies:
//
// "The reference rates are usually updated at around 16:00 CET every working day, except on TARGET closing days."
//
// from https://www.ecb.europa.eu/stats/policy_and_exchange_rates/euro_reference_exchange_rates/html/index.en.html
type YMD time.Time

func NewYMD(s string) (YMD, error) {
	frankfurt, err := time.LoadLocation("Europe/Berlin")

	if err != nil {
		return YMD{}, err
	}

	result, err := time.ParseInLocation(time.DateOnly, s, frankfurt)

	if err != nil {
		return YMD{}, fmt.Errorf("unable to interpret '%s' as YYYY-MM-DD format: %w", s, err)
	}

	result = result.Add(16 * time.Hour)

	return YMD(result), nil
}

func (d YMD) String() string {
	return time.Time(d).Format(time.DateOnly)
}

func (d YMD) Before(u YMD) bool {
	return time.Time(d).Before(time.Time(u))
}

func (d YMD) IsZero() bool {
	return time.Time(d).IsZero()
}

func (d YMD) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Quote(d.String())), nil
}

func (d *YMD) UnmarshalJSON(data []byte) error {
	unquoted, err := strconv.Unquote(string(data))

	if err != nil {
		return fmt.Errorf("unable to unquote '%s': %w", data, err)
	}

	s, err := NewYMD(unquoted)

	if err != nil {
		return err
	}

	*d = s

	return nil
}
