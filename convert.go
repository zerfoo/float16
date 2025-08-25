package float16

import (
	"math"
	"strconv"
)

// FromFloat32 converts a float32 value to a Float16 value.
// It handles special cases like NaN, infinities, and zeros.
// The conversion follows IEEE 754-2008 rules for half-precision.
func FromFloat32(f32 float32) Float16 {
	// Use the more accurate converter with proper rounding and subnormal handling
	return fromFloat32New(f32)
}

// FromFloat32WithRounding converts a float32 to Float16 using the provided rounding mode.
// It mirrors fromFloat32New but respects the explicit rounding mode instead of always
// rounding to nearest-even.
func FromFloat32WithRounding(f32 float32, mode RoundingMode) Float16 {
	bits := math.Float32bits(f32)
	sign := uint16(bits >> 31)
	exp := int32((bits >> 23) & 0xff)
	mant := uint32(bits & 0x7fffff)

	// Special cases
	if exp == 0xff {
		if mant == 0 {
			return Float16(sign<<15 | 0x7c00) // infinity
		}
		return Float16(sign<<15 | 0x7e00) // qNaN
	}

	// Zero (preserve sign)
	if exp == 0 && mant == 0 {
		return Float16(sign << 15)
	}

	// Adjust exponent bias: float32 (127) -> float16 (15)
	exp -= 127 - 15

	// Overflow to infinity
	if exp >= 0x1f {
		return Float16(sign<<15 | 0x7c00)
	}

	// Underflow and subnormals
	if exp <= 0 {
		if exp < -10 {
			// Too small for subnormal even after rounding; return signed zero
			return Float16(sign << 15)
		}
		// Convert to subnormal
		mant = (mant | 1<<23) >> uint(1-exp)
		// Round mantissa down to 10 bits using the requested mode
		if shouldRoundWithMode(mant, 13, sign<<15, mode) {
			mant += 1 << 13
		}
		return Float16(uint16(sign<<15) | uint16(mant>>13))
	}

	// Normal numbers
	mant |= 1 << 23 // restore implicit leading 1

	// Round mantissa down to 10 bits using the requested mode
	if shouldRoundWithMode(mant, 13, sign<<15, mode) {
		mant += 1 << 13
	}

	// Check for mantissa overflow after rounding
	if mant >= 1<<24 {
		exp++
		mant = 0 // implicit 1 will be added by format
	}

	// Exponent overflow => infinity
	if exp >= 0x1f {
		return Float16(sign<<15 | 0x7c00)
	}

	mantissa10 := (mant >> 13) & 0x3ff
	return Float16(uint16(sign<<15) | uint16(exp<<10) | uint16(mantissa10))
}

// shouldRoundWithMode is like shouldRound but uses an explicit rounding mode
// rather than the global DefaultRoundingMode. The meaning of parameters matches
// shouldRound: mantissa is the bits prior to truncation, shift is the number of
// bits being truncated, sign carries SignMask for sign checks.
func shouldRoundWithMode(mantissa uint32, shift int, sign uint16, mode RoundingMode) bool {
	if shift <= 0 {
		return false
	}

	guard := (mantissa >> uint(shift-1)) & 1
	sticky := mantissa & ((1 << uint(shift-1)) - 1)
	lsb := (mantissa >> uint(shift)) & 1
	anyDiscarded := (guard | (boolToUint(sticky != 0))) == 1

	switch mode {
	case RoundNearestEven:
		return guard == 1 && (sticky != 0 || lsb == 1)
	case RoundNearestAway:
		return guard == 1 || sticky != 0
	case RoundTowardZero:
		return false
	case RoundTowardPositive:
		return (sign&SignMask) == 0 && anyDiscarded
	case RoundTowardNegative:
		return (sign&SignMask) != 0 && anyDiscarded
	default:
		return false
	}
}

// ToFloat32 converts a Float16 value to a float32 value.
// It handles special cases like NaN, infinities, and zeros.
func (f Float16) ToFloat32() float32 {
	bits := uint16(f)
	sign := (bits & SignMask) != 0
	exp := (bits & ExponentMask) >> MantissaLen
	mant := bits & MantissaMask

	// Handle special cases
	if exp == ExponentInfinity {
		if mant != 0 { // NaN
			return float32(math.NaN())
		}
		if sign {
			return float32(math.Inf(-1))
		}
		return float32(math.Inf(1))
	}

	if exp == ExponentZero {
		if mant == 0 {
			if sign {
				return float32(math.Copysign(0.0, -1.0))
			}
			return 0.0
		}
		// Subnormal: mant * 2^-24
		val := float32(mant) * (1.0 / (1 << 24))
		if sign {
			return -val
		}
		return val
	}

	// Normalized: (1 + mant/2^10) * 2^(exp-15)
	val := (1.0 + float32(mant)/1024.0) * float32(math.Ldexp(1, int(exp)-ExponentBias))
	if sign {
		return -val
	}
	return val
}

// FromFloat64 converts a float64 value to a Float16 value.
// It handles special cases like NaN, infinities, and zeros.
func FromFloat64(f64 float64) Float16 {
	return FromFloat32(float32(f64)) // Simplified: convert via float32
}

// ToFloat16 converts a float64 to a Float16 value.
// This is a convenience wrapper used in tests and utilities.
func ToFloat16(f64 float64) Float16 {
	return FromFloat64(f64)
}

// ToSlice16 converts a slice of float32 to a slice of Float16.
// This is a convenience wrapper used in tests and utilities.
func ToSlice16(s []float32) []Float16 {
	result := make([]Float16, len(s))
	for i, v := range s {
		result[i] = FromFloat32(v)
	}
	return result
}

// FromFloat64WithMode converts a float64 to Float16 with specified conversion and rounding modes
func FromFloat64WithMode(f64 float64, convMode ConversionMode, roundMode RoundingMode) (Float16, error) {
	// Basic conversion first
	result := FromFloat64(f64)

	if convMode == ModeStrict {
		// NaN
		if math.IsNaN(f64) {
			return 0, &Float16Error{Op: "FromFloat64WithMode", Msg: "NaN in strict mode", Code: ErrNaN}
		}
		// Infinity
		if math.IsInf(f64, 0) {
			return 0, &Float16Error{Op: "FromFloat64WithMode", Msg: "infinity in strict mode", Code: ErrInfinity}
		}
		// Overflow: magnitude exceeds max finite float16
		max := MaxValue.ToFloat64()
		if math.Abs(f64) > max {
			return 0, &Float16Error{Op: "FromFloat64WithMode", Msg: "overflow", Code: ErrOverflow}
		}
		// Underflow: result became subnormal or zero for non-zero input
		if f64 != 0 && (result.IsZero() || result.IsSubnormal()) {
			return 0, &Float16Error{Op: "FromFloat64WithMode", Msg: "underflow", Code: ErrUnderflow}
		}
	}

	return result, nil
}

// ToFloat64 converts a Float16 value to a float64 value.
// It handles special cases like NaN, infinities, and zeros.
func (f Float16) ToFloat64() float64 {
	return float64(f.ToFloat32()) // Simplified: convert via float32
}

// shouldRound determines whether to round up during conversion
// This is a helper function used in conversion algorithms
func shouldRound(mantissa uint32, shift int, sign uint16) bool {
	if shift <= 0 {
		return false
	}

	// Bits about to be discarded
	guard := (mantissa >> uint(shift-1)) & 1
	sticky := mantissa & ((1 << uint(shift-1)) - 1)
	lsb := (mantissa >> uint(shift)) & 1
	anyDiscarded := (guard | (boolToUint(sticky != 0))) == 1

	switch DefaultRoundingMode {
	case RoundNearestEven:
		// Round up if guard=1 and (sticky!=0 or LSB is 1) => ties to even
		return guard == 1 && (sticky != 0 || lsb == 1)
	case RoundNearestAway:
		// Round up on half or more (guard=1). If less than half (guard=0), do not round.
		// sticky doesn't affect decision except that if sticky>0, it's strictly more than half.
		return guard == 1 || sticky != 0
	case RoundTowardZero:
		return false
	case RoundTowardPositive:
		// Round up for positive numbers if any discarded bits are non-zero
		return (sign&SignMask) == 0 && anyDiscarded
	case RoundTowardNegative:
		// Round up (i.e., toward -inf increases magnitude) for negative numbers if discarded bits
		return (sign&SignMask) != 0 && anyDiscarded
	default:
		// Invalid rounding mode: do not round
		return false
	}
}

// boolToUint converts a bool to 0/1 as uint32
func boolToUint(b bool) uint32 {
	if b {
		return 1
	}
	return 0
}

// Parse converts a string to a Float16 value
// This is a simplified implementation for testing
func Parse(s string) (Float16, error) {
	// Minimal parser: return error for standard numeric strings (not implemented)
	return 0, &Float16Error{Op: "Parse", Msg: "parsing not implemented for numeric strings", Code: ErrInvalidOperation}
}

// FromInt converts an integer to Float16
func FromInt(i int) Float16 {
	return FromFloat32(float32(i))
}

// ToSlice16WithMode converts a slice of float32 to Float16 with specified modes
func ToSlice16WithMode(s []float32, convMode ConversionMode, roundMode RoundingMode) ([]Float16, []error) {
	result := make([]Float16, len(s))
	errs := make([]error, len(s))

	for i, v := range s {
		// Convert
		result[i] = FromFloat32(v)
		errs[i] = nil

		if convMode == ModeStrict {
			// Overflow if magnitude exceeds max finite Float16
			max := MaxValue.ToFloat64()
			if math.Abs(float64(v)) > max {
				errs[i] = &Float16Error{Op: "ToSlice16WithMode", Msg: "overflow", Code: ErrOverflow}
				continue
			}
			// Underflow if non-zero converted to subnormal or zero
			if v != 0 && (result[i].IsZero() || result[i].IsSubnormal()) {
				errs[i] = &Float16Error{Op: "ToSlice16WithMode", Msg: "underflow", Code: ErrUnderflow}
			}
		}
	}
	return result, errs
}

// ToSlice32 converts a slice of Float16 to a slice of float32
func ToSlice32(s []Float16) []float32 {
	result := make([]float32, len(s))
	for i, v := range s {
		result[i] = v.ToFloat32()
	}
	return result
}

// ToSlice64 converts a slice of Float16 to a slice of float64
func ToSlice64(s []Float16) []float64 {
	result := make([]float64, len(s))
	for i, v := range s {
		result[i] = v.ToFloat64()
	}
	return result
}

// FromSlice64 converts a slice of float64 to a slice of Float16
func FromSlice64(s []float64) []Float16 {
	result := make([]Float16, len(s))
	for i, v := range s {
		result[i] = FromFloat64(v)
	}
	return result
}

// FromInt32 converts an int32 to Float16
func FromInt32(i int32) Float16 {
	return FromFloat32(float32(i))
}

// FromInt64 converts an int64 to Float16
func FromInt64(i int64) Float16 {
	return FromFloat64(float64(i))
}

// ParseFloat converts a string to a Float16 value.
// The precision parameter is ignored for Float16.
// It returns the Float16 value and an error if the string cannot be parsed.
func ParseFloat(s string, precision int) (Float16, error) {
	// This implementation is a placeholder and does not fully parse
	// a string to a float16. It only handles basic cases.
	// A full implementation would require a more complex parser.

	switch s {
	case "NaN":
		return NaN(), nil
	case "+Inf", "Inf":
		return PositiveInfinity, nil
	case "-Inf":
		return NegativeInfinity, nil
	case "+0", "0":
		return PositiveZero, nil
	case "-0":
		return NegativeZero, nil
	}

	// Attempt to parse as float32 and convert
	f32, err := strconv.ParseFloat(s, 32)
	if err != nil {
		return 0, err
	}
	return FromFloat32(float32(f32)), nil
}
