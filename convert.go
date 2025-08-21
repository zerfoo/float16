package float16

import (
	"fmt"
	"math"
	"strconv"
)

// FromFloat32 converts a float32 value to a Float16 value.
// It handles special cases like NaN, infinities, and zeros.
// The conversion follows IEEE 754-2008 rules for half-precision.
func FromFloat32(f32 float32) Float16 {
	f32Bits := math.Float32bits(f32)
	sign := f32Bits & 0x80000000
	exp := (f32Bits >> 23) & 0xFF
	mant := f32Bits & 0x7FFFFF

	var f16Bits uint16

	if exp == 0xFF { // NaN or Infinity
		if mant != 0 { // NaN
			f16Bits = 0x7E00 // Quiet NaN
		} else { // Infinity
			if sign != 0 {
				f16Bits = 0xFC00 // Negative Infinity
			} else {
				f16Bits = 0x7C00 // Positive Infinity
			}
		}
	} else if exp == 0 { // Zero or Denormalized
		if mant == 0 { // Zero
			if sign != 0 {
				f16Bits = 0x8000 // Negative Zero
			} else {
				f16Bits = 0x0000 // Positive Zero
			}
		} else { // Denormalized float32, convert to float16 denormalized or zero
			// Shift mantissa right to align with float16 denormalized range
			// This is a simplified approach and might not be perfectly accurate for all denormals
			// A more robust implementation would involve proper rounding and handling of underflow
			f16Bits = uint16(mant >> 13) // Shift 23 - 10 = 13 bits
			if sign != 0 {
				f16Bits |= 0x8000
			}
		}
	} else { // Normalized float32
		exp16 := int(exp) - 127 + 15 // Adjust bias
		if exp16 >= 31 { // Overflow, convert to infinity
			if sign != 0 {
				f16Bits = 0xFC00 // Negative Infinity
			} else {
				f16Bits = 0x7C00 // Positive Infinity
			}
		} else if exp16 <= 0 { // Underflow, convert to denormalized or zero
			// This is a simplified approach. Proper denormalization requires
			// shifting the mantissa and potentially losing precision.
			if exp16 == 0 { // Smallest normalized float32 maps to float16 denormal
				f16Bits = uint16(mant >> 13) // Shift 23 - 10 = 13 bits
			} else { // Smaller than smallest denormal, convert to zero
				f16Bits = 0x0000
			}
			if sign != 0 {
				f16Bits |= 0x8000
			}
		} else { // Normalized float16
			f16Bits = uint16(exp16<<10) | uint16(mant>>13) // Shift 23 - 10 = 13 bits
			if sign != 0 {
				f16Bits |= 0x8000
			}
		}
	}
	return Float16(f16Bits)
}

// ToFloat32 converts a Float16 value to a float32 value.
// It handles special cases like NaN, infinities, and zeros.
func (f Float16) ToFloat32() float32 {
	f16Bits := uint16(f)
	sign := uint32(f16Bits & 0x8000) << 16 // Shift to float32 sign position
	exp := (f16Bits >> 10) & 0x1F
	mant := f16Bits & 0x3FF

	var f32Bits uint32

	if exp == 0x1F { // NaN or Infinity
		if mant != 0 { // NaN
			f32Bits = 0x7FC00000 // Quiet NaN
		} else { // Infinity
			if sign != 0 {
				f32Bits = 0xFF800000 // Negative Infinity
			} else {
				f32Bits = 0x7F800000 // Positive Infinity
			}
		}
	} else if exp == 0 { // Zero or Denormalized
		if mant == 0 { // Zero
			if sign != 0 {
				f32Bits = 0x80000000 // Negative Zero
			} else {
				f32Bits = 0x00000000 // Positive Zero
			}
		} else { // Denormalized float16, convert to float32 denormalized
			// Shift mantissa left to align with float32 denormalized range
			// This is a simplified approach and might not be perfectly accurate for all denormals
			// A more robust implementation would involve proper scaling
			f32Bits = uint32(mant) << 13 // Shift 10 + 13 = 23 bits
			if sign != 0 {
				f32Bits |= 0x80000000
			}
		}
	} else { // Normalized float16
		exp32 := uint32(int(exp) - 15 + 127) // Adjust bias
		f32Bits = sign | (exp32 << 23) | (uint32(mant) << 13) // Shift 10 + 13 = 23 bits
	}
	return math.Float32frombits(f32Bits)
}

// FromFloat64 converts a float64 value to a Float16 value.
// It handles special cases like NaN, infinities, and zeros.
func FromFloat64(f64 float64) Float16 {
	return FromFloat32(float32(f64)) // Simplified: convert via float32
}

// ToFloat64 converts a Float16 value to a float64 value.
// It handles special cases like NaN, infinities, and zeros.
func (f Float16) ToFloat64() float64 {
	return float64(f.ToFloat32()) // Simplified: convert via float32
}

// FromBits creates a Float16 from its raw uint16 bit representation.
func FromBits(bits uint16) Float16 {
	return Float16(bits)
}

// Bits returns the raw uint16 bit representation of a Float16.
func (f Float16) Bits() uint16 {
	return uint16(f)
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


