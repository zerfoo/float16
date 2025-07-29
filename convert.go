package float16

import (
	"math"

	"github.com/x448/float16"
)

// Global conversion settings
var (
	DefaultConversionMode = ModeIEEE
	DefaultRoundingMode   = RoundNearestEven
)

// ToFloat16 converts a float32 value to Float16 format using default settings
func ToFloat16(f32 float32) Float16 {
	return Float16(float16.Fromfloat32(f32).Bits())
}

// ToFloat16WithMode converts a float32 to Float16 with specified conversion and rounding modes
func ToFloat16WithMode(f32 float32, convMode ConversionMode, roundMode RoundingMode) (Float16, error) {
	if convMode == ModeStrict {
		if math.IsInf(float64(f32), 0) {
			return 0, &Float16Error{Code: ErrInfinity}
		}
		if math.IsNaN(float64(f32)) {
			return 0, &Float16Error{Code: ErrNaN}
		}
		if f32 > 65504.0 || f32 < -65504.0 {
			return 0, &Float16Error{Code: ErrOverflow}
		}
		if f32 != 0 && math.Abs(float64(f32)) < 6.103515625e-05 {
			return 0, &Float16Error{Code: ErrUnderflow}
		}
	}

	f16 := float16.Fromfloat32(f32)
	return Float16(f16.Bits()), nil
}

// ToFloat32 converts a Float16 value to float32 with full precision
func (f Float16) ToFloat32() float32 {
	return float16.Frombits(uint16(f)).Float32()
}

// ToFloat64 converts a Float16 value to float64 with full precision
func (f Float16) ToFloat64() float64 {
	return float64(f.ToFloat32())
}

// FromFloat32 converts a float32 to Float16 (with potential precision loss)
func FromFloat32(f32 float32) Float16 {
	return Float16(float16.Fromfloat32(f32).Bits())
}

// FromFloat64 converts a float64 to Float16 (with potential precision loss)
func FromFloat64(f64 float64) Float16 {
	return Float16(float16.Fromfloat32(float32(f64)).Bits())
}

// FromFloat64WithMode converts a float64 to Float16 with specified modes
func FromFloat64WithMode(f64 float64, convMode ConversionMode, roundMode RoundingMode) (Float16, error) {
	if convMode == ModeStrict {
		if math.IsInf(f64, 0) {
			return 0, &Float16Error{Code: ErrInfinity}
		}
		if math.IsNaN(f64) {
			return 0, &Float16Error{Code: ErrNaN}
		}
		if f64 > 65504.0 || f64 < -65504.0 {
			return 0, &Float16Error{Code: ErrOverflow}
		}
		if f64 != 0 && math.Abs(f64) < 6.103515625e-05 {
			return 0, &Float16Error{Code: ErrUnderflow}
		}
	}

	f16 := FromFloat64(f64)
	if roundMode == RoundTowardZero {
		if f64 > 0 {
			f16 = FromFloat64(math.Floor(f64))
		} else {
			f16 = FromFloat64(math.Ceil(f64))
		}
	} else if roundMode == RoundTowardPositive {
		f16 = FromFloat64(math.Ceil(f64))
	} else if roundMode == RoundTowardNegative {
		f16 = FromFloat64(math.Floor(f64))
	}
	return f16, nil
}

// ToSlice16 converts a slice of float32 to Float16 with optimized performance
func ToSlice16(f32s []float32) []Float16 {
	if len(f32s) == 0 {
		return nil
	}
	res := make([]Float16, len(f32s))
	for i, f := range f32s {
		res[i] = ToFloat16(f)
	}
	return res
}

// ToSlice32 converts a slice of Float16 to float32 with optimized performance
func ToSlice32(f16s []Float16) []float32 {
	if len(f16s) == 0 {
		return nil
	}
	res := make([]float32, len(f16s))
	for i, f := range f16s {
		res[i] = f.ToFloat32()
	}
	return res
}

// ToSlice64 converts a slice of Float16 to float64 with optimized performance
func ToSlice64(f16s []Float16) []float64 {
	if len(f16s) == 0 {
		return nil
	}
	res := make([]float64, len(f16s))
	for i, f := range f16s {
		res[i] = f.ToFloat64()
	}
	return res
}

// FromSlice64 converts a slice of float64 to Float16 with optimized performance
func FromSlice64(f64s []float64) []Float16 {
	if len(f64s) == 0 {
		return nil
	}
	res := make([]Float16, len(f64s))
	for i, f := range f64s {
		res[i] = FromFloat64(f)
	}
	return res
}

// ToSlice16WithMode converts a slice with specified conversion mode
func ToSlice16WithMode(f32s []float32, convMode ConversionMode, roundMode RoundingMode) ([]Float16, []error) {
	if len(f32s) == 0 {
		return nil, nil
	}
	res := make([]Float16, len(f32s))
	errs := []error{}
	for i, f := range f32s {
		r, err := ToFloat16WithMode(f, convMode, roundMode)
		if err != nil {
			errs = append(errs, err)
		}
		res[i] = r
	}
	return res, errs
}

// Integer conversion functions

// FromInt converts an integer to Float16
func FromInt(i int) Float16 {
	return ToFloat16(float32(i))
}

// FromInt32 converts an int32 to Float16
func FromInt32(i int32) Float16 {
	return ToFloat16(float32(i))
}

// FromInt64 converts an int64 to Float16 (with potential precision loss)
func FromInt64(i int64) Float16 {
	return ToFloat16(float32(i))
}

// ToInt converts a Float16 to int (truncated toward zero)
func (f Float16) ToInt() int {
	return int(f.ToFloat32())
}

// ToInt32 converts a Float16 to int32 (truncated toward zero)
func (f Float16) ToInt32() int32 {
	return int32(f.ToFloat32())
}

// ToInt64 converts a Float16 to int64 (truncated toward zero)
func (f Float16) ToInt64() int64 {
	return int64(f.ToFloat32())
}

// Parse converts a string to Float16 (placeholder for future implementation)
func Parse(s string) (Float16, error) {
	// This would implement string parsing - simplified for now
	// In a full implementation, this would parse various float formats
	return PositiveZero, &Float16Error{
		Op:   "parse",
		Msg:  "string parsing not implemented",
		Code: ErrInvalidOperation,
	}
}
func shouldRound(mantissa uint32, shift int, mode RoundingMode) bool {
	if shift <= 0 {
		return false
	}

	// Get the bits that will be discarded
	discardedBits := mantissa & ((1 << shift) - 1)
	guardBit := (mantissa >> (shift - 1)) & 1

	switch mode {
	case RoundNearestEven:
		if guardBit == 0 {
			return false
		}
		// If guard bit is 1, check for exact halfway case
		remainingBits := discardedBits & ((1 << (shift - 1)) - 1)
		if remainingBits != 0 {
			return true // Not exact halfway, round up
		}
		// Exact halfway: round to even (check LSB of result)
		resultLSB := (mantissa >> shift) & 1
		return resultLSB == 1

	case RoundNearestAway:
		return guardBit == 1

	case RoundTowardZero:
		return false

	case RoundTowardPositive:
		return discardedBits != 0

	case RoundTowardNegative:
		return false // This function doesn't know sign, caller must handle

	default:
		return guardBit == 1 // Default to nearest even guard bit behavior
	}
}