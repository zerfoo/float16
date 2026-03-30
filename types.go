package float16

import (
	"fmt"
)

// ErrorCode represents specific error categories for float16 operations
type ErrorCode int

const (
	ErrInvalidOperation ErrorCode = iota
	ErrNaN
	ErrInfinity
	ErrOverflow
	ErrUnderflow
	ErrDivisionByZero
	ErrNotImplemented
)

// Float16Error provides detailed error information for float16 operations
type Float16Error struct {
	Op   string
	Msg  string
	Code ErrorCode
}

func (e *Float16Error) Error() string {
	if e == nil {
		return "<nil>"
	}
	if e.Op != "" {
		return fmt.Sprintf("float16 %s: %s", e.Op, e.Msg)
	}
	return "float16: " + e.Msg
}

// RoundingMode controls how results are rounded during conversion/arithmetic
type RoundingMode int

const (
	// Round to nearest, ties to even
	RoundNearestEven RoundingMode = iota
	// Round toward zero (truncate)
	RoundTowardZero
	// Round toward +Inf
	RoundTowardPositive
	// Round toward -Inf
	RoundTowardNegative
	// Round to nearest, ties away from zero
	RoundNearestAway
)

// ConversionMode controls error reporting behavior for conversions
type ConversionMode int

const (
	// ModeIEEE performs IEEE-style conversion, saturating to Inf/0 with no errors
	ModeIEEE ConversionMode = iota
	// ModeStrict reports errors for NaN, Inf, overflow, and underflow
	ModeStrict
)

// Float16 represents a 16-bit IEEE 754 half-precision floating-point value
type Float16 uint16

// Bits returns the IEEE 754 half-precision bit pattern of f
func (f Float16) Bits() uint16 { return uint16(f) }

// FromBits constructs a Float16 from its IEEE 754 half-precision bit pattern
func FromBits(b uint16) Float16 { return Float16(b) }

// IEEE 754 half-precision format constants
const (
	SignMask     = 0x8000 // 0b1000000000000000 - Sign bit mask
	ExponentMask = 0x7C00 // 0b0111110000000000 - Exponent bits mask
	MantissaMask = 0x03FF // 0b0000001111111111 - Mantissa bits mask
	MantissaLen  = 10     // Number of mantissa bits
	ExponentLen  = 5      // Number of exponent bits

	// Exponent bias and limits for IEEE 754 half-precision
	// bias = 2^(exponent_bits-1) - 1 = 2^4 - 1 = 15
	ExponentBias = 15 // Bias for 5-bit exponent
	ExponentMax  = 31 // Maximum exponent value (11111 binary)
	ExponentMin  = 0  // Minimum exponent value

	// Normalized exponent range
	ExponentNormalMin = 1  // Minimum normalized exponent
	ExponentNormalMax = 30 // Maximum normalized exponent (infinity at 31)

	// Float32 constants for conversion
	Float32ExponentBias = 127 // IEEE 754 single precision bias
	Float32ExponentLen  = 8   // Float32 exponent bits
	Float32MantissaLen  = 23  // Float32 mantissa bits

	// Special exponent values
	ExponentZero     = 0  // Zero and subnormal numbers
	ExponentInfinity = 31 // Infinity and NaN
)

// Special values following IEEE 754 half-precision standard
const (
	PositiveZero     Float16 = 0x0000 // +0.0
	NegativeZero     Float16 = 0x8000 // -0.0
	PositiveInfinity Float16 = 0x7C00 // +∞
	NegativeInfinity Float16 = 0xFC00 // -∞

	// Largest finite values
	MaxValue Float16 = 0x7BFF // Largest positive finite value (~65504)
	MinValue Float16 = 0xFBFF // Largest negative finite value (~-65504)

	// Smallest normalized positive value
	SmallestNormal Float16 = 0x0400 // 2^-14 ≈ 6.103515625e-05

	// Largest subnormal value
	LargestSubnormal Float16 = 0x03FF // (1023/1024) * 2^-14 ≈ 6.097555161e-05

	// Smallest positive subnormal value
	SmallestSubnormal Float16 = 0x0001 // 2^-24 ≈ 5.960464478e-08

	// Common NaN representations
	QuietNaN     Float16 = 0x7E00 // Quiet NaN (most significant mantissa bit set)
	SignalingNaN Float16 = 0x7D00 // Signaling NaN
	NegativeQNaN Float16 = 0xFE00 // Negative quiet NaN
)

// IsZero returns true if the Float16 value represents zero (positive or negative)
func (f Float16) IsZero() bool {
	return (f & 0x7FFF) == 0
}

// IsInf returns true if the Float16 value represents infinity
// If sign > 0, returns true only for positive infinity
// If sign < 0, returns true only for negative infinity
// If sign == 0, returns true for either infinity
func (f Float16) IsInf(sign int) bool {
	if (f & 0x7FFF) != PositiveInfinity {
		return false
	}
	if sign == 0 {
		return true
	}
	return (sign > 0) == ((f & SignMask) == 0)
}

// IsNaN returns true if the Float16 value represents NaN (Not a Number)
func (f Float16) IsNaN() bool {
	exp := (f & ExponentMask) >> MantissaLen
	mant := f & MantissaMask
	return exp == ExponentInfinity && mant != 0
}

// IsFinite returns true if the Float16 value is finite (not infinity or NaN)
func (f Float16) IsFinite() bool {
	exp := (f & ExponentMask) >> MantissaLen
	return exp != ExponentInfinity
}

// IsNormal returns true if the Float16 value is normalized (not zero, subnormal, infinite, or NaN)
func (f Float16) IsNormal() bool {
	exp := (f & ExponentMask) >> MantissaLen
	return exp != ExponentZero && exp != ExponentInfinity
}

// IsSubnormal returns true if the Float16 value is subnormal (denormalized)
func (f Float16) IsSubnormal() bool {
	exp := (f & ExponentMask) >> MantissaLen
	mant := f & MantissaMask
	return exp == ExponentZero && mant != 0
}

// FloatClass enumerates the IEEE 754 classification of a Float16 value
type FloatClass int

const (
	ClassPositiveZero FloatClass = iota
	ClassNegativeZero
	ClassPositiveSubnormal
	ClassNegativeSubnormal
	ClassPositiveNormal
	ClassNegativeNormal
	ClassPositiveInfinity
	ClassNegativeInfinity
	ClassQuietNaN
	ClassSignalingNaN
)

// Class returns the IEEE 754 classification of the value
func (f Float16) Class() FloatClass {
	bits := uint16(f)
	sign := (bits & SignMask) != 0
	exp := (bits & ExponentMask) >> MantissaLen
	mant := bits & MantissaMask

	switch exp {
	case ExponentZero:
		if mant == 0 {
			if sign {
				return ClassNegativeZero
			}
			return ClassPositiveZero
		}
		if sign {
			return ClassNegativeSubnormal
		}
		return ClassPositiveSubnormal
	case ExponentInfinity:
		if mant == 0 {
			if sign {
				return ClassNegativeInfinity
			}
			return ClassPositiveInfinity
		}
		// NaN: distinguish quiet vs signaling by top mantissa bit (bit 9)
		if (mant & (1 << (MantissaLen - 1))) != 0 {
			return ClassQuietNaN
		}
		return ClassSignalingNaN
	default:
		if sign {
			return ClassNegativeNormal
		}
		return ClassPositiveNormal
	}
}

// Sign returns the sign of the Float16 value: 1 for positive, -1 for negative, 0 for zero
func (f Float16) Sign() int {
	if f.IsZero() {
		return 0
	}
	if (f & SignMask) != 0 {
		return -1
	}
	return 1
}

// Signbit returns true if the Float16 value has a negative sign bit
func (f Float16) Signbit() bool {
	return (f & SignMask) != 0
}

// Abs returns the absolute value of the Float16
func (f Float16) Abs() Float16 {
	return f & 0x7FFF // Clear sign bit
}

// Neg returns the negation of the Float16
func (f Float16) Neg() Float16 {
	return f ^ SignMask // Flip sign bit
}

// CopySign returns a value with the magnitude of f and the sign of s
func (f Float16) CopySign(s Float16) Float16 {
	// Clear sign bit of f, then OR with sign bit of s
	return (f & ^Float16(SignMask)) | (s & Float16(SignMask))
}

// ToInt converts Float16 to int (truncates toward zero)
func (f Float16) ToInt() int {
	return int(f.ToFloat32())
}

// String returns a string representation of the Float16 value
func (f Float16) String() string {
	if f.IsNaN() {
		if f.Signbit() {
			return "-NaN"
		}
		return "NaN"
	}
	if f.IsInf(0) {
		if f.Signbit() {
			return "-Inf"
		}
		return "+Inf"
	}
	return fmt.Sprintf("%.6g", f.ToFloat32())
}

// GoString returns a Go syntax representation of the Float16 value
func (f Float16) GoString() string {
	return fmt.Sprintf("float16.FromBits(0x%04x)", uint16(f))
}

func (f Float16) ToInt32() int32 {
	return int32(f.ToFloat32())
}

// ToInt64 converts Float16 to int64 (truncates toward zero)
func (f Float16) ToInt64() int64 {
	return int64(f.ToFloat32())
}
