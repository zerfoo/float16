package float16

import (
	"fmt"
	"math/bits"
)

// Float16 represents a 16-bit IEEE 754 half-precision floating-point value
type Float16 uint16

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

	// Float64 constants for conversion
	Float64ExponentBias = 1023 // IEEE 754 double precision bias
	Float64MantissaLen  = 52   // Float64 mantissa bits

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

// ConversionMode defines how conversions handle edge cases
type ConversionMode int

const (
	// ModeIEEE uses standard IEEE 754 rounding and special value behavior
	ModeIEEE ConversionMode = iota
	// ModeStrict returns errors for overflow, underflow, and NaN
	ModeStrict
	// ModeFast optimizes for performance, may sacrifice some precision
	ModeFast
	// ModeExact preserves exact values when possible, errors on precision loss
	ModeExact
)

// RoundingMode defines IEEE 754 rounding behavior
type RoundingMode int

const (
	// DefaultRoundingMode rounds to nearest, ties to even (IEEE default)
	DefaultRoundingMode RoundingMode = iota
	// RoundNearestAway rounds to nearest, ties away from zero
	RoundNearestAway
	// RoundTowardZero truncates toward zero
	RoundTowardZero
	// RoundTowardPositive rounds toward +∞
	RoundTowardPositive
	// RoundTowardNegative rounds toward -∞
	RoundTowardNegative
)

// Float16Error represents errors that can occur during Float16 operations
type Float16Error struct {
	Op    string      // Operation that caused the error
	Value interface{} // Input value that caused the error
	Msg   string      // Error message
	Code  ErrorCode   // Specific error code
}

// ErrorCode represents specific error types
type ErrorCode int

const (
	ErrOverflow ErrorCode = iota
	ErrUnderflow
	ErrInvalidOperation
	ErrDivisionByZero
	ErrInexact
	ErrNaN
	ErrInfinity
)

func (e *Float16Error) Error() string {
	if e.Value != nil {
		return fmt.Sprintf("float16.%s: %s (value: %v)", e.Op, e.Msg, e.Value)
	}
	return fmt.Sprintf("float16.%s: %s", e.Op, e.Msg)
}

// Predefined error instances
var (
	ErrOverflowError  = &Float16Error{Code: ErrOverflow, Msg: "value too large for float16"}
	ErrUnderflowError = &Float16Error{Code: ErrUnderflow, Msg: "value too small for float16"}
	ErrNaNError       = &Float16Error{Code: ErrNaN, Msg: "NaN in strict mode"}
	ErrInfinityError  = &Float16Error{Code: ErrInfinity, Msg: "infinity in strict mode"}
	ErrDivByZeroError = &Float16Error{Code: ErrDivisionByZero, Msg: "division by zero"}
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

// CopySign returns a Float16 with the magnitude of f and the sign of sign
func (f Float16) CopySign(sign Float16) Float16 {
	return (f & 0x7FFF) | (sign & SignMask)
}

// Bits returns the underlying uint16 representation
func (f Float16) Bits() uint16 {
	return uint16(f)
}

// FromBits creates a Float16 from its bit representation
func FromBits(bits uint16) Float16 {
	return Float16(bits)
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

// Class returns the IEEE 754 class of the floating-point value
type FloatClass int

const (
	ClassSignalingNaN FloatClass = iota
	ClassQuietNaN
	ClassNegativeInfinity
	ClassNegativeNormal
	ClassNegativeSubnormal
	ClassNegativeZero
	ClassPositiveZero
	ClassPositiveSubnormal
	ClassPositiveNormal
	ClassPositiveInfinity
)

// Class returns the IEEE 754 classification of the Float16 value
func (f Float16) Class() FloatClass {
	if f.IsNaN() {
		// Check if it's a signaling NaN (MSB of mantissa is 0)
		if (f & 0x0200) == 0 {
			return ClassSignalingNaN
		}
		return ClassQuietNaN
	}

	sign := f.Signbit()

	if f.IsInf(0) {
		if sign {
			return ClassNegativeInfinity
		}
		return ClassPositiveInfinity
	}

	if f.IsZero() {
		if sign {
			return ClassNegativeZero
		}
		return ClassPositiveZero
	}

	if f.IsSubnormal() {
		if sign {
			return ClassNegativeSubnormal
		}
		return ClassPositiveSubnormal
	}

	// Normal number
	if sign {
		return ClassNegativeNormal
	}
	return ClassPositiveNormal
}

// Utility functions for bit manipulation

// extractComponents extracts sign, exponent, and mantissa from Float16
func (f Float16) extractComponents() (sign uint16, exp uint16, mant uint16) {
	bits := uint16(f)
	sign = (bits & SignMask) >> 15
	exp = (bits & ExponentMask) >> MantissaLen
	mant = bits & MantissaMask
	return
}

// packComponents packs sign, exponent, and mantissa into Float16
func packComponents(sign, exp, mant uint16) Float16 {
	return Float16((sign << 15) | (exp << MantissaLen) | (mant & MantissaMask))
}

// leadingZeros counts leading zeros in a 10-bit mantissa
func leadingZeros10(x uint16) int {
	if x == 0 {
		return 10
	}
	return bits.LeadingZeros16(x<<6) - 6 // Shift to align with 16-bit and adjust
}
