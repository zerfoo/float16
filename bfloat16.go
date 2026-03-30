package float16

import (
	"fmt"
	"math"
)

// BFloat16 represents a 16-bit "Brain Floating Point" format value
// Used by Google Brain, TensorFlow, and various ML frameworks
// Format: 1 sign bit, 8 exponent bits, 7 mantissa bits
type BFloat16 uint16

// BFloat16 format constants
const (
	BFloat16SignMask     = 0x8000 // 0b1000000000000000 - Sign bit mask
	BFloat16ExponentMask = 0x7F80 // 0b0111111110000000 - Exponent bits mask
	BFloat16MantissaMask = 0x007F // 0b0000000001111111 - Mantissa bits mask
	BFloat16MantissaLen  = 7      // Number of mantissa bits
	BFloat16ExponentLen  = 8      // Number of exponent bits

	// Exponent bias and limits for BFloat16
	// bias = 2^(exponent_bits-1) - 1 = 2^7 - 1 = 127 (same as Float32)
	BFloat16ExponentBias = 127 // Bias for 8-bit exponent
	BFloat16ExponentMax  = 255 // Maximum exponent value
	BFloat16ExponentMin  = 0   // Minimum exponent value

	// Normalized exponent range
	BFloat16ExponentNormalMin = 1   // Minimum normalized exponent
	BFloat16ExponentNormalMax = 254 // Maximum normalized exponent (infinity at 255)

	// Special exponent values
	BFloat16ExponentZero     = 0   // Zero and subnormal numbers
	BFloat16ExponentInfinity = 255 // Infinity and NaN
)

// Special BFloat16 values
const (
	BFloat16PositiveZero     BFloat16 = 0x0000 // +0.0
	BFloat16NegativeZero     BFloat16 = 0x8000 // -0.0
	BFloat16PositiveInfinity BFloat16 = 0x7F80 // +∞
	BFloat16NegativeInfinity BFloat16 = 0xFF80 // -∞
	BFloat16QuietNaN         BFloat16 = 0x7FC0 // Quiet NaN
	BFloat16SignalingNaN     BFloat16 = 0x7F81 // Signaling NaN

	// Largest finite values
	BFloat16MaxValue    BFloat16 = 0x7F7F // Largest positive normal
	BFloat16MinValue    BFloat16 = 0xFF7F // Largest negative normal (most negative)
	BFloat16SmallestPos BFloat16 = 0x0080 // Smallest positive normal
	BFloat16SmallestNeg BFloat16 = 0x8080 // Smallest negative normal

	// Smallest subnormal values
	BFloat16SmallestPosSubnormal BFloat16 = 0x0001 // Smallest positive subnormal
	BFloat16SmallestNegSubnormal BFloat16 = 0x8001 // Smallest negative subnormal
)

// FromBits creates a BFloat16 from its bit representation
func BFloat16FromBits(bits uint16) BFloat16 {
	return BFloat16(bits)
}

// Bits returns the bit representation of the BFloat16
func (b BFloat16) Bits() uint16 {
	return uint16(b)
}

// FromFloat32 converts a float32 to BFloat16 using round-to-nearest-even
// BFloat16 is essentially a truncated float32, so conversion is straightforward
func BFloat16FromFloat32(f float32) BFloat16 {
	return BFloat16FromFloat32WithRounding(f, RoundNearestEven)
}

// BFloat16FromFloat32WithRounding converts a float32 to BFloat16 with the specified rounding mode.
func BFloat16FromFloat32WithRounding(f float32, mode RoundingMode) BFloat16 {
	bits := math.Float32bits(f)
	sign := bits & 0x80000000
	exp := bits & 0x7F800000
	mant := bits & 0x007FFFFF

	// Handle special cases: NaN, Inf, Zero
	if exp == 0x7F800000 { // Inf or NaN
		return BFloat16(bits >> 16) // Preserve Inf/NaN
	}
	if exp == 0 && mant == 0 { // Zero
		return BFloat16(sign >> 16) // Preserve signed zero
	}

	// Extract the high 16 bits (sign, exponent, and 7 MSBs of mantissa)
	bfloat16Bits := bits >> 16

	// Check the bit at position 15 of the original float32 bits (the first bit to be truncated)
	roundBit := (bits >> 15) & 0x1

	// Check if any of the bits from position 0 to 14 are non-zero
	stickyBits := bits & 0x7FFF

	switch mode {
	case RoundNearestEven:
		// Round to nearest, ties to even
		if roundBit == 1 { // If the bit to be rounded is 1
			if stickyBits != 0 { // If there are any non-zero bits after the roundBit
				bfloat16Bits++ // Round up
			} else { // It's a tie (roundBit is 1, stickyBits are 0)
				if (bfloat16Bits & 0x1) == 1 { // If the LSB of the truncated BFloat16 is odd
					bfloat16Bits++ // Round up to make it even
				}
			}
		}
	case RoundTowardZero:
		// Truncate, which is effectively rounding toward zero
		// No action needed as the shift already truncates
	case RoundTowardPositive:
		// Round toward +Inf
		if sign == 0 && (roundBit == 1 || stickyBits != 0) { // Positive and needs rounding up
			bfloat16Bits++
		}
	case RoundTowardNegative:
		// Round toward -Inf
		if sign != 0 && (roundBit == 1 || stickyBits != 0) { // Negative and needs rounding down
			bfloat16Bits++
		}
	case RoundNearestAway:
		// Round to nearest, ties away from zero
		if roundBit == 1 {
			bfloat16Bits++
		}
	default:
		// Default to RoundNearestEven if an unknown mode is provided
		if roundBit == 1 && (stickyBits != 0 || (bfloat16Bits&0x1) == 1) {
			bfloat16Bits++
		}
	}

	return BFloat16(bfloat16Bits)
}

// ToFloat32 converts BFloat16 to float32
func (b BFloat16) ToFloat32() float32 {
	// Expand back to 32 bits by shifting left 16 positions
	// The lower 16 bits become zero
	bits := uint32(b) << 16
	return math.Float32frombits(bits)
}

// FromFloat64 converts a float64 to BFloat16
func BFloat16FromFloat64(f float64) BFloat16 {
	return BFloat16FromFloat64WithRounding(f, RoundNearestEven)
}

// BFloat16FromFloat64WithRounding converts a float64 to BFloat16 with the specified rounding mode.
func BFloat16FromFloat64WithRounding(f float64, mode RoundingMode) BFloat16 {
	return BFloat16FromFloat32WithRounding(float32(f), mode)
}

// BFloat16FromFloat32WithMode converts a float32 to BFloat16 with specified conversion and rounding modes.
func BFloat16FromFloat32WithMode(f32 float32, convMode ConversionMode, roundMode RoundingMode) (BFloat16, error) {
	// First, perform the rounding conversion
	b := BFloat16FromFloat32WithRounding(f32, roundMode)

	// Check for special values and ranges
	if math.IsNaN(float64(f32)) {
		if convMode == ModeStrict {
			return 0, &Float16Error{Op: "BFloat16FromFloat32", Msg: "NaN conversion in strict mode", Code: ErrNaN}
		}
		return BFloat16QuietNaN, nil
	}

	if math.IsInf(float64(f32), 0) {
		if convMode == ModeStrict {
			return 0, &Float16Error{Op: "BFloat16FromFloat32", Msg: "Inf conversion in strict mode", Code: ErrInfinity}
		}
		// Already handled by BFloat16FromFloat32WithRounding, which preserves Inf
		return b, nil
	}

	// Check for overflow/underflow against BFloat16's finite range
	// Convert BFloat16 max/min values to float32 for comparison
	bf16Max := BFloat16MaxValue.ToFloat32()
	bf16Min := BFloat16MinValue.ToFloat32()
	bf16SmallestNormalPos := BFloat16SmallestPos.ToFloat32()

	if f32 > bf16Max || f32 < bf16Min {
		if convMode == ModeStrict {
			return 0, &Float16Error{Op: "BFloat16FromFloat32", Msg: "overflow in strict mode", Code: ErrOverflow}
		}
		// ModeIEEE: saturate to infinity
		if f32 > 0 {
			return BFloat16PositiveInfinity, nil
		}
		return BFloat16NegativeInfinity, nil
	}

	// Check for underflow to zero (denormalized numbers or zero)
	// If the original float32 is non-zero but smaller than the smallest normal BFloat16
	// and the result after rounding is zero, it's an underflow.
	if f32 != 0 && math.Abs(float64(f32)) < float64(bf16SmallestNormalPos) && b.IsZero() {
		if convMode == ModeStrict {
			return 0, &Float16Error{Op: "BFloat16FromFloat32", Msg: "underflow in strict mode", Code: ErrUnderflow}
		}
		// ModeIEEE: saturate to zero (already handled by rounding to zero)
		return b, nil
	}

	return b, nil
}

// BFloat16FromFloat64WithMode converts a float64 to BFloat16 with specified conversion and rounding modes.
func BFloat16FromFloat64WithMode(f64 float64, convMode ConversionMode, roundMode RoundingMode) (BFloat16, error) {
	// Convert float64 to float32 first, then use the float32 conversion logic
	// This might lose precision for float64, but BFloat16 is based on float32's exponent range.
	b, err := BFloat16FromFloat32WithMode(float32(f64), convMode, roundMode)
	if err != nil {
		return 0, err
	}
	return b, nil
}

// String returns a string representation of the BFloat16
func (b BFloat16) String() string {
	if b.IsNaN() {
		return "NaN"
	}
	if b.IsInf(1) {
		return "+Inf"
	}
	if b.IsInf(-1) {
		return "-Inf"
	}
	return fmt.Sprintf("%g", b.ToFloat32())
}

// Classification methods

// Class returns the IEEE 754 classification of the BFloat16 value
func (b BFloat16) Class() FloatClass {
	bits := uint16(b)
	sign := (bits & BFloat16SignMask) != 0
	exp := (bits & BFloat16ExponentMask) >> BFloat16MantissaLen
	mant := bits & BFloat16MantissaMask

	switch exp {
	case BFloat16ExponentZero:
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
	case BFloat16ExponentInfinity:
		if mant == 0 {
			if sign {
				return ClassNegativeInfinity
			}
			return ClassPositiveInfinity
		}
		// NaN: distinguish quiet vs signaling by top mantissa bit (bit 6 of BFloat16 mantissa)
		// BFloat16 mantissa is 7 bits, so top bit is bit 6 (0-indexed)
		if (mant & (1 << (BFloat16MantissaLen - 1))) != 0 {
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

// IsZero returns true if the BFloat16 is zero (positive or negative)
func (b BFloat16) IsZero() bool {
	return (uint16(b) &^ BFloat16SignMask) == 0
}

// IsInf reports whether b is an infinity, according to sign
func (b BFloat16) IsInf(sign int) bool {
	bits := uint16(b)
	exp := (bits & BFloat16ExponentMask) >> BFloat16MantissaLen
	mantissa := bits & BFloat16MantissaMask

	if exp != BFloat16ExponentInfinity || mantissa != 0 {
		return false
	}

	if sign == 0 {
		return true
	}

	signBit := (bits & BFloat16SignMask) != 0
	return (sign > 0) == !signBit
}

// IsNaN reports whether b is an IEEE 754 "not-a-number" value
func (b BFloat16) IsNaN() bool {
	bits := uint16(b)
	exp := (bits & BFloat16ExponentMask) >> BFloat16MantissaLen
	mantissa := bits & BFloat16MantissaMask

	return exp == BFloat16ExponentInfinity && mantissa != 0
}

// IsFinite reports whether b is neither infinite nor NaN
func (b BFloat16) IsFinite() bool {
	return !b.IsInf(0) && !b.IsNaN()
}

// IsNormal reports whether b is a normal number
func (b BFloat16) IsNormal() bool {
	if b.IsNaN() || b.IsInf(0) || b.IsZero() {
		return false
	}

	bits := uint16(b)
	exp := (bits & BFloat16ExponentMask) >> BFloat16MantissaLen
	return exp != BFloat16ExponentZero
}

// IsSubnormal reports whether b is a subnormal number
func (b BFloat16) IsSubnormal() bool {
	if b.IsNaN() || b.IsInf(0) || b.IsZero() {
		return false
	}

	bits := uint16(b)
	exp := (bits & BFloat16ExponentMask) >> BFloat16MantissaLen
	mantissa := bits & BFloat16MantissaMask
	return exp == BFloat16ExponentZero && mantissa != 0
}

// Signbit reports whether b is negative or negative zero
func (b BFloat16) Signbit() bool {
	return (uint16(b) & BFloat16SignMask) != 0
}

// CopySign returns a value with the magnitude of f and the sign of s
func (b BFloat16) CopySign(s BFloat16) BFloat16 {
	// Clear sign bit of b, then OR with sign bit of s
	return (b &^ BFloat16SignMask) | (s & BFloat16SignMask)
}

// Arithmetic operations

// BFloat16Add adds two BFloat16 values
func BFloat16Add(a, b BFloat16) BFloat16 {
	return BFloat16FromFloat32(a.ToFloat32() + b.ToFloat32())
}

// BFloat16Sub subtracts two BFloat16 values
func BFloat16Sub(a, b BFloat16) BFloat16 {
	return BFloat16FromFloat32(a.ToFloat32() - b.ToFloat32())
}

// BFloat16Mul multiplies two BFloat16 values
func BFloat16Mul(a, b BFloat16) BFloat16 {
	return BFloat16FromFloat32(a.ToFloat32() * b.ToFloat32())
}

// BFloat16Div divides two BFloat16 values
func BFloat16Div(a, b BFloat16) BFloat16 {
	return BFloat16FromFloat32(a.ToFloat32() / b.ToFloat32())
}

// Comparison operations

// BFloat16Equal returns true if a equals b
func BFloat16Equal(a, b BFloat16) bool {
	// Handle NaN case
	if a.IsNaN() || b.IsNaN() {
		return false
	}

	// Handle zero case (positive and negative zero are equal)
	if a.IsZero() && b.IsZero() {
		return true
	}

	return a == b
}

// BFloat16Less returns true if a < b
func BFloat16Less(a, b BFloat16) bool {
	return a.ToFloat32() < b.ToFloat32()
}

// BFloat16LessEqual returns true if a <= b
func BFloat16LessEqual(a, b BFloat16) bool {
	return a.ToFloat32() <= b.ToFloat32()
}

// BFloat16Greater returns true if a > b
func BFloat16Greater(a, b BFloat16) bool {
	return a.ToFloat32() > b.ToFloat32()
}

// BFloat16GreaterEqual returns true if a >= b
func BFloat16GreaterEqual(a, b BFloat16) bool {
	return a.ToFloat32() >= b.ToFloat32()
}

// Utility functions

// BFloat16Abs returns the absolute value of b
func BFloat16Abs(b BFloat16) BFloat16 {
	return BFloat16(uint16(b) &^ BFloat16SignMask)
}

// BFloat16Neg returns the negation of b
func BFloat16Neg(b BFloat16) BFloat16 {
	return BFloat16(uint16(b) ^ BFloat16SignMask)
}

// BFloat16Min returns the smaller of a or b
func BFloat16Min(a, b BFloat16) BFloat16 {
	if a.IsNaN() || b.IsNaN() {
		return BFloat16QuietNaN
	}
	if BFloat16Less(a, b) {
		return a
	}
	return b
}

// BFloat16Max returns the larger of a or b
func BFloat16Max(a, b BFloat16) BFloat16 {
	if a.IsNaN() || b.IsNaN() {
		return BFloat16QuietNaN
	}
	if BFloat16Greater(a, b) {
		return a
	}
	return b
}

// Cross-conversion between Float16 and BFloat16

// ToBFloat16 converts a Float16 to BFloat16
func (f Float16) ToBFloat16() BFloat16 {
	return BFloat16FromFloat32(f.ToFloat32())
}

// ToFloat16 converts a BFloat16 to Float16
func (b BFloat16) ToFloat16() Float16 {
	return FromFloat32(b.ToFloat32())
}

// BFloat16FromFloat16 converts a Float16 to BFloat16
func BFloat16FromFloat16(f Float16) BFloat16 {
	return f.ToBFloat16()
}

// Float16FromBFloat16 converts a BFloat16 to Float16
func Float16FromBFloat16(b BFloat16) Float16 {
	return b.ToFloat16()
}

// Batch operations for high-performance computing

// BFloat16AddSlice performs element-wise addition of two BFloat16 slices
func BFloat16AddSlice(a, b []BFloat16) []BFloat16 {
	if len(a) != len(b) {
		panic("float16: slice length mismatch")
	}

	result := make([]BFloat16, len(a))
	for i := range a {
		result[i] = BFloat16Add(a[i], b[i])
	}
	return result
}

// BFloat16SubSlice performs element-wise subtraction of two BFloat16 slices
func BFloat16SubSlice(a, b []BFloat16) []BFloat16 {
	if len(a) != len(b) {
		panic("float16: slice length mismatch")
	}

	result := make([]BFloat16, len(a))
	for i := range a {
		result[i] = BFloat16Sub(a[i], b[i])
	}
	return result
}

// BFloat16MulSlice performs element-wise multiplication of two BFloat16 slices
func BFloat16MulSlice(a, b []BFloat16) []BFloat16 {
	if len(a) != len(b) {
		panic("float16: slice length mismatch")
	}

	result := make([]BFloat16, len(a))
	for i := range a {
		result[i] = BFloat16Mul(a[i], b[i])
	}
	return result
}

// BFloat16DivSlice performs element-wise division of two BFloat16 slices
func BFloat16DivSlice(a, b []BFloat16) []BFloat16 {
	if len(a) != len(b) {
		panic("float16: slice length mismatch")
	}

	result := make([]BFloat16, len(a))
	for i := range a {
		result[i] = BFloat16Div(a[i], b[i])
	}
	return result
}

// BFloat16ScaleSlice multiplies each element in the slice by a scalar
func BFloat16ScaleSlice(s []BFloat16, scalar BFloat16) []BFloat16 {
	result := make([]BFloat16, len(s))
	for i := range s {
		result[i] = BFloat16Mul(s[i], scalar)
	}
	return result
}

// BFloat16SumSlice returns the sum of all elements in the slice
func BFloat16SumSlice(s []BFloat16) BFloat16 {
	sum := BFloat16PositiveZero
	for _, v := range s {
		sum = BFloat16Add(sum, v)
	}
	return sum
}

// BFloat16ToSlice32 converts a slice of BFloat16 values to float32
func BFloat16ToSlice32(s []BFloat16) []float32 {
	result := make([]float32, len(s))
	for i, v := range s {
		result[i] = v.ToFloat32()
	}
	return result
}

// BFloat16FromSlice32 converts a slice of float32 values to BFloat16
func BFloat16FromSlice32(s []float32) []BFloat16 {
	result := make([]BFloat16, len(s))
	for i, v := range s {
		result[i] = BFloat16FromFloat32(v)
	}
	return result
}

// BFloat16ToSlice64 converts a slice of BFloat16 values to float64
func BFloat16ToSlice64(s []BFloat16) []float64 {
	result := make([]float64, len(s))
	for i, v := range s {
		result[i] = float64(v.ToFloat32())
	}
	return result
}

// BFloat16FromSlice64 converts a slice of float64 values to BFloat16
func BFloat16FromSlice64(s []float64) []BFloat16 {
	result := make([]BFloat16, len(s))
	for i, v := range s {
		result[i] = BFloat16FromFloat32(float32(v))
	}
	return result
}

// Convenience constants for common BFloat16 values
var (
	BFloat16Zero  = BFloat16PositiveZero
	BFloat16One   = BFloat16FromFloat32(1.0)
	BFloat16Two   = BFloat16FromFloat32(2.0)
	BFloat16Half  = BFloat16FromFloat32(0.5)
	BFloat16E     = BFloat16FromFloat32(float32(math.E))
	BFloat16Pi    = BFloat16FromFloat32(float32(math.Pi))
	BFloat16Sqrt2 = BFloat16FromFloat32(float32(math.Sqrt2))
)
