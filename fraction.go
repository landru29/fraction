package fraction

import (
	"fmt"
	"math"
	"strconv"
)

const (
	bitSize = 64

	defaultRenderingDeep = 65536
)

// Number is a fraction number.
type Number struct {
	integer    int64
	floating   float64
	isNegative bool
}

// Raw is the raw value of the number.
func (n Number) Raw() float64 {
	return (float64(n.integer) + n.floating) * map[bool]float64{
		true:  -1.0,
		false: 1.0,
	}[n.isNegative]
}

// String implements the pflag.Value interface.
func (n Number) String() string {
	if n.IsZero() {
		return "0"
	}

	numerator, denominator, exact := n.Render(defaultRenderingDeep)

	return fmt.Sprintf(
		"%s %d / %d %s",
		n.Sign(),
		numerator+n.Numerator(denominator),
		denominator,
		map[bool]string{true: "!", false: ""}[exact],
	)
}

// Numerator is the integer part.
func (n Number) Numerator(denominator int64) int64 {
	return n.integer * denominator
}

// Sign is the negative sign.
func (n Number) Sign() string {
	return map[bool]string{true: "-", false: ""}[n.isNegative]
}

// IsZero checks if the number is zero.
func (n Number) IsZero() bool {
	return n.integer == 0 && n.floating == 0
}

// Set implements the pflag.Value interface.
func (n *Number) Set(val string) error {
	num, err := strconv.ParseFloat(val, bitSize)
	if err != nil {
		return err
	}

	*n = Number{
		isNegative: math.Signbit(num),
		integer:    int64(math.Abs(num)),
		floating:   math.Abs(num) - math.Floor(math.Abs(num)),
	}

	return nil
}

// Type implements the pflag.Value interface.
func (n Number) Type() string {
	return "number"
}

// Render computes the fraction
func (n Number) Render(deep int) (int64, int64, bool) {
	numerator0, denominator0, _, _, exact := approximate(
		n.floating,
		0,
		1,
		1,
		1,
		deep,
	)

	return numerator0, denominator0, exact
}

func approximate(
	value float64,
	numerator0 int64,
	denominator0 int64,
	numerator1 int64,
	denominator1 int64,
	step int,
) (int64, int64, int64, int64, bool) {
	mediantNumerator := numerator0 + numerator1
	mediantDenominator := denominator0 + denominator1

	mediant := float64(mediantNumerator) / float64(mediantDenominator)

	if step < 0 {
		return mediantNumerator, mediantDenominator, mediantNumerator, mediantDenominator, mediant == value
	}

	switch {
	case value < mediant:
		return approximate(value, numerator0, denominator0, mediantNumerator, mediantDenominator, step-1)
	case value > mediant:
		return approximate(value, mediantNumerator, mediantDenominator, numerator1, denominator1, step-1)
	default:
		return mediantNumerator, mediantDenominator, mediantNumerator, mediantDenominator, true
	}
}
