package float16

import (
	"math"

	"github.com/x448/float16"
)

// Converter holds the conversion and rounding modes for float16 operations.
type Converter struct {
	ConversionMode ConversionMode
	RoundingMode   RoundingMode
}

// NewConverter creates a new Converter with the specified modes.
func NewConverter(convMode ConversionMode, roundMode RoundingMode) *Converter {
	return &Converter{
		ConversionMode: convMode,
		RoundingMode:   roundMode,
	}
}

// ToFloat16 converts a float32 value to Float16 format using the Converter's settings.
func (c *Converter) ToFloat16(f32 float32) Float16 {
	return Float16(float16.Fromfloat32(f32).Bits())
}

// ToFloat16WithMode converts a float32 to Float16 with specified conversion and rounding modes
func (c *Converter) ToFloat16WithMode(f32 float32) (Float16, error) {
	convMode := c.ConversionMode
	roundMode := c.RoundingMode
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
func (c *Converter) FromFloat32(f32 float32) Float16 {
	return c.ToFloat16(f32)
}

// FromFloat64 converts a float64 to Float16 (with potential precision loss)
func (c *Converter) FromFloat64(f64 float64) Float16 {
	return c.ToFloat16(float32(f64))
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

	return c.ToFloat16WithMode(float32(f64))
}

// ToSlice16 converts a slice of float32 to Float16 with optimized performance
func (c *Converter) ToSlice16(f32s []float32) []Float16 {
	if len(f32s) == 0 {
		return nil
	}
	res := make([]Float16, len(f32s))
	for i, f := range f32s {
		res[i] = c.ToFloat16(f)
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
func (c *Converter) FromSlice64(f64s []float64) []Float16 {
	if len(f64s) == 0 {
		return nil
	}
	res := make([]Float16, len(f64s))
	for i, f := range f64s {
		res[i] = c.FromFloat64(f)
	}
	return res
}

// ToSlice16WithMode converts a slice with specified conversion mode
func (c *Converter) ToSlice16WithMode(f32s []float32) ([]Float16, []error) {
	if len(f32s) == 0 {
		return nil, nil
	}
	res := make([]Float16, len(f32s))
	errs := []error{}
	for i, f := range f32s {
		r, err := c.ToFloat16WithMode(f)
		if err != nil {
			errs = append(errs, err)
		}
		res[i] = r
	}
	return res, errs
}

// Integer conversion functions

// FromInt converts an integer to Float16
func (c *Converter) FromInt(i int) Float16 {
	return c.ToFloat16(float32(i))
}

// FromInt32 converts an int32 to Float16
func (c *Converter) FromInt32(i int32) Float16 {
	return c.ToFloat16(float32(i))
}

// FromInt64 converts an int64 to Float16 (with potential precision loss)
func (c *Converter) FromInt64(i int64) Float16 {
	return c.ToFloat16(float32(i))
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
func (c *Converter) Parse(s string) (Float16, error) {
	// This would implement string parsing - simplified for now
	// In a full implementation, this would parse various float formats
	return PositiveZero, &Float16Error{
		Op:   "parse",
		Msg:  "string parsing not implemented",
		Code: ErrInvalidOperation,
	}
}
func (c *Converter) shouldRound(mantissa uint32, shift int, sign uint16) bool {
	switch c.RoundingMode {
	case RoundNearestEven:
		// If the value is exactly halfway, round to the nearest even number.
		if mantissa&(1<<uint(shift-1)) != 0 && mantissa&((1<<uint(shift-1))-1) == 0 {
			return (mantissa>>uint(shift))&1 != 0
		}
		// Otherwise, round to the nearest number.
		return mantissa&(1<<uint(shift-1)) != 0
	case RoundNearestAway:
		return mantissa&(1<<uint(shift-1)) != 0
	case RoundTowardZero:
		return false
	case RoundTowardPositive:
		return sign == 0 && mantissa&((1<<uint(shift))-1) != 0
	case RoundTowardNegative:
		return sign != 0 && mantissa&((1<<uint(shift))-1) != 0
	}
	return false
}
