package frankfurter

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type ExchangeRatesService struct {
	HttpClient *http.Client
	URL        string
}

type ExchangeRates struct {
	Date   YMD
	Amount Micros
	Base   Currency
	Rates  Rates
}

type History struct {
	Amount Micros
	Base   Currency
	Start  YMD     `json:"start_date"`
	End    YMD     `json:"end_date"`
	Rates  RatesAt `json:"rates"`
}

type Rates map[Currency]Micros
type RatesAt map[YMD]Rates
type Currency string

// MicrosPerUnit is the number of micros per currency unit.
const MicrosPerUnit = 1_000_000

// Micros stores an exchange rate as an integer number of micros
// (1/1,000,000 of a currency unit) instead of a floating point number.
type Micros int64

// String formats the rate as a decimal with trailing zeros removed,
// e.g. 11_321_500 micros becomes "11.3215".
func (m Micros) String() string {
	negative := m < 0

	if negative {
		m = -m
	}

	whole := m / MicrosPerUnit
	fraction := m % MicrosPerUnit

	if fraction == 0 {
		return strconv.FormatInt(int64(whole), 10)
	}

	fractionStr := strings.TrimRight(fmt.Sprintf("%06d", int64(fraction)), "0")
	result := strconv.FormatInt(int64(whole), 10) + "." + fractionStr

	if negative {
		result = "-" + result
	}

	return result
}

// UnmarshalJSON converts a decimal rate, e.g. "11.3215", to micros by
// rounding to the nearest integer micro. It does so without intermediate
// floating point arithmetic so that the conversion is exact.
func (m *Micros) UnmarshalJSON(data []byte) error {
	raw := string(data)

	negative := false

	if strings.HasPrefix(raw, "-") {
		negative = true
		raw = raw[1:]
	}

	if strings.ContainsAny(raw, "eE") {
		return fmt.Errorf("unable to parse %q as micros: scientific notation is not supported", data)
	}

	whole, fraction, _ := strings.Cut(raw, ".")

	digits := whole + fraction

	if len(digits) == 0 {
		return fmt.Errorf("unable to parse %q as micros", data)
	}

	for _, digit := range digits {
		if digit < '0' || digit > '9' {
			return fmt.Errorf("unable to parse %q as micros", data)
		}
	}

	// Round to the nearest micro, i.e. keep up to 6 fractional digits.
	roundUp := len(fraction) > 6 && fraction[6] >= '5'

	if len(fraction) > 6 {
		fraction = fraction[:6]
	}

	// Pad the fractional part to exactly 6 digits.
	fraction += strings.Repeat("0", max(0, 6-len(fraction)))

	value, err := strconv.ParseInt(whole+fraction, 10, 64)

	if err != nil {
		return fmt.Errorf("unable to parse %q as micros: %w", data, err)
	}

	if roundUp {
		value++
	}

	if negative {
		value = -value
	}

	*m = Micros(value)

	return nil
}

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

func (d YMD) Equal(other YMD) bool {
	return time.Time(d).Equal(time.Time(other))
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
