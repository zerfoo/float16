package float16

import (
	"fmt"
	"math"
	"unsafe"
)

// Global conversion settings
var (
	DefaultConversionMode = ModeIEEE
	DefaultRoundingMode   = RoundNearestEven
)

// ToFloat16 converts a float32 value to Float16 format using default settings
func ToFloat16(f32 float32) Float16 {
	// Handle special cases first
	if f32 == 0 {
		if math.Signbit(float64(f32)) {
			return NegativeZero
		}
		return PositiveZero
	}

	bits := math.Float32bits(f32)
	sign := (bits >> 31) & 0x1
	exp32 := (bits >> 23) & 0xFF
	mant32 := bits & 0x007FFFFF

	// Special cases: NaN and Inf
	if exp32 == 0xFF {
		if mant32 == 0 {
			// Infinity
			if sign != 0 {
				return NegativeInfinity
			}
			return PositiveInfinity
		}
		// NaN - preserve sign and some payload
		nanMant := uint16((mant32 >> 13) & 0x03FF)
		if nanMant == 0 {
			nanMant = 1 // Ensure it's a NaN, not infinity
		}
		return Float16((uint16(sign) << 15) | 0x7C00 | nanMant)
	}

	// For subnormal float32, we need to convert to float16 subnormal
	if exp32 == 0 {
		// The value is subnormal, so we need to convert it to a float16 subnormal
		// or flush to zero if it's too small.
		// A float32 subnormal is value * 2^-126.
		// We need to convert it to a float16 subnormal, which is value * 2^-14.
		// So we need to shift the mantissa by 126 - 14 = 112 bits.
		// Since the float32 mantissa has 23 bits, we will lose a lot of precision.
		// We can approximate this by converting the float32 to float64 and then to float16.
		return FromFloat64(float64(f32))
	}
	// Normal number in float16
	exp16 := exp32 - 127 + 15
	if exp16 >= 31 {
		// Overflow
		if sign != 0 {
			return NegativeInfinity
		}
		return PositiveInfinity
	}

	if exp16 <= 0 {
		// Underflow to subnormal
		shift := uint(1 - exp16)
		mant32 |= 0x800000 // Add implicit leading 1
		mant16 := uint16(mant32 >> (shift + 13))
		// Rounding
		roundBit := (mant32 >> (shift + 12)) & 1
		if roundBit != 0 {
			mant16++
		}
		return Float16((uint16(sign) << 15) | mant16)
	}

	// Extract mantissa bits (10 bits) with rounding
	mant16 := uint16((mant32 + 0x1000) >> 13)

	// Check for overflow in mantissa (due to rounding)
	if (mant16 & 0x0400) != 0 {
		mant16 = 0
		exp16++
	}

	// Check for overflow after rounding
	if exp16 >= 31 {
		if sign != 0 {
			return NegativeInfinity
		}
		return PositiveInfinity
	}

	// Combine sign, exponent, and mantissa
	return Float16((uint16(sign) << 15) | (uint16(exp16) << 10) | (mant16 & 0x03FF))
}

// ToFloat16WithMode converts a float32 to Float16 with specified conversion and rounding modes
func ToFloat16WithMode(f32 float32, convMode ConversionMode, roundMode RoundingMode) (Float16, error) {
	// Handle special cases first for performance
	if f32 == 0.0 {
		// Use math.Float32bits to check the sign bit directly
		bits := math.Float32bits(f32)
		if (bits & (1 << 31)) != 0 { // Check sign bit
			return NegativeZero, nil
		}
		return PositiveZero, nil
	}

	if math.IsInf(float64(f32), 0) {
		if convMode == ModeStrict {
			return 0, &Float16Error{
				Op:    "convert",
				Value: f32,
				Msg:   "infinity not allowed in strict mode",
				Code:  ErrInfinity,
			}
		}
		if f32 > 0 {
			return PositiveInfinity, nil
		}
		return NegativeInfinity, nil
	}

	if math.IsNaN(float64(f32)) {
		if convMode == ModeStrict {
			return 0, &Float16Error{
				Op:    "convert",
				Value: f32,
				Msg:   "NaN not allowed in strict mode",
				Code:  ErrNaN,
			}
		}
		// Return quiet NaN, preserving sign
		if math.Signbit(float64(f32)) {
			return NegativeQNaN, nil
		}
		return QuietNaN, nil
	}

	// Extract IEEE 754 float32 components using bit manipulation for performance
	bits := math.Float32bits(f32)
	sign32 := bits >> 31
	exp32 := int32((bits >> 23) & 0xFF)
	mant32 := bits & 0x7FFFFF

	// Handle subnormal float32 inputs (exp32 == 0) specially
	isSubnormal := exp32 == 0
	if isSubnormal && mant32 == 0 {
		// Zero
		if sign32 != 0 {
			return NegativeZero, nil
		}
		return PositiveZero, nil
	}

	// Convert exponent from float32 bias to float16 bias
	exp16 := exp32 - Float32ExponentBias + ExponentBias

	// Handle overflow to infinity
	if exp16 >= ExponentInfinity {
		if convMode == ModeStrict || convMode == ModeExact {
			return 0, &Float16Error{
				Op:    "convert",
				Value: f32,
				Msg:   "overflow: value too large for float16",
				Code:  ErrOverflow,
			}
		}
		// Clamp to infinity
		if sign32 != 0 {
			return NegativeInfinity, nil
		}
		return PositiveInfinity, nil
	}

	// Handle underflow cases
	if exp16 <= 0 {
		// Check if we can represent as subnormal
		shift := 1 - exp16 // Number of bits to shift right to make it subnormal

		// For subnormal float32 inputs (exp32 == 0), we don't add the implicit leading 1
		// For normal float32 inputs, we add the implicit leading 1
		if exp32 != 0 {
			mant32 |= 0x800000 // Add implicit leading 1
		}

		// Calculate total shift needed for denormalization
		// For subnormal float16, the exponent is 0, so we need to shift right by:
		// 1. The difference in exponent biases (127 - 15 = 112)
		// 2. Plus 1 to account for the implicit leading 1 in float32
		// 3. Minus 1 because we're already accounting for the subnormal shift in the exponent
		// This simplifies to: (127 - 15) + 1 - 1 = 112
		totalShift := Float32ExponentBias - ExponentBias + 1 - 1

		// For subnormal float16, we need to shift right by an additional (1 - exp16)
		// But since exp16 is 0 for subnormals, this becomes (1 - 0) = 1
		totalShift += int(shift)

		// Check for complete underflow (beyond what we can represent even with subnormals)
		// The smallest positive subnormal float16 is 2^-24, which requires 24 bits of precision
		// If we need to shift more than 23 bits, we'll lose all precision
		if totalShift >= 24 { // 23 (mantissa) + 1 (round bit)
			if convMode == ModeStrict || convMode == ModeExact {
				return 0, &Float16Error{
					Op:    "convert",
					Value: f32,
					Msg:   "underflow: value too small for float16",
					Code:  ErrUnderflow,
				}
			}
			// Flush to zero
			if sign32 != 0 {
				return NegativeZero, nil
			}
			return PositiveZero, nil
		}

		// For subnormals, we need to shift right by totalShift, keeping extra bits for rounding
		// We'll keep one extra bit for the round bit and one for the sticky bit
		extraBits := 2
		if totalShift > 22 { // If we're shifting more than 22 bits, we won't have enough bits left
			extraBits = 0
		} else if totalShift > 21 { // Only room for round bit
			extraBits = 1
		}

		// Extract the bits we'll keep, plus extra bits for rounding
		var mant16 uint16
		if mant32 != 0 {
			// For subnormal float32 inputs, we don't add the implicit leading 1
			// For normal float32 inputs, we add the implicit leading 1
			if exp32 != 0 {
				mant32 |= 0x800000 // Add implicit leading 1 for normal numbers
			}
			mant16 = uint16((mant32 >> (totalShift - extraBits)))
		}

		// Check if we need to round
		roundBit := uint32(0)
		stickyBit := uint32(0)

		if totalShift > extraBits {
			roundBit = (mant32 >> (totalShift - extraBits - 1)) & 0x1
		}

		if totalShift > extraBits+1 {
			stickyMask := (uint32(1) << (totalShift - extraBits - 1)) - 1
			stickyBit = mant32 & stickyMask
			if stickyBit != 0 {
				stickyBit = 1
			}
		}

		// Apply rounding
		if shouldRound(uint32(mant16), int(roundBit|stickyBit), roundMode) {
			mant16++
			// Check for carry that would require renormalization
			if mant16 > 0x3FF {
				mant16 = 0x200 // 1.0 * 2^-10 (smallest normal)
				exp16 = 1      // Exponent for 2^-14
				// No need to check for overflow here since we're dealing with subnormals
			}
		}

		// For subnormals, the exponent is 0
		exp16 = 0

		// If we have a normal result after rounding, adjust exponent and mantissa
		if mant16 >= 0x400 {
			// This can happen due to rounding up from a value just below the normal range
			exp16 = 1
			mant16 >>= 1

			// If we're still in the normal range, we're done
			if mant16 < 0x400 {
				return packComponents(uint16(sign32), uint16(exp16), mant16), nil
			}

			// If we still have a value >= 0x400, it means we rounded up to the next power of two
			// This should only happen if we had a value very close to the next power of two
			// and we rounded up due to the rounding mode
			mant16 >>= 1
			exp16++
		}

		return packComponents(uint16(sign32), uint16(exp16), mant16), nil
	}

	// Normal number conversion
	// Extract top 10 bits of mantissa for float16
	mant16 := mant32 >> (Float32MantissaLen - MantissaLen)

	// Apply rounding
	if shouldRound(mant32, Float32MantissaLen-MantissaLen, roundMode) {
		mant16++
		// Handle mantissa overflow
		if mant16 >= (1 << MantissaLen) {
			mant16 = 0
			exp16++
			// Check for exponent overflow after rounding
			if exp16 >= ExponentInfinity {
				if convMode == ModeStrict || convMode == ModeExact {
					return 0, &Float16Error{
						Op:    "convert",
						Value: f32,
						Msg:   "overflow after rounding",
						Code:  ErrOverflow,
					}
				}
				if sign32 != 0 {
					return NegativeInfinity, nil
				}
				return PositiveInfinity, nil
			}
		}
	}

	return packComponents(uint16(sign32), uint16(exp16), uint16(mant16)), nil
}

// shouldRound determines if rounding should occur based on the rounding mode
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

// ToFloat32 converts a Float16 value to float32 with full precision
func (f Float16) ToFloat32() float32 {
	// Handle special cases
	if f.IsZero() {
		if f.Signbit() {
			return math.Float32frombits(0x80000000) // -0.0
		}
		return 0.0
	}

	if f.IsNaN() {
		// Preserve NaN payload and sign
		sign := uint32(0)
		if f.Signbit() {
			sign = 0x80000000
		}
		payload := uint32(f & MantissaMask)
		// Ensure it's a quiet NaN and preserve payload
		return math.Float32frombits(sign | 0x7FC00000 | (payload << (Float32MantissaLen - MantissaLen)))
	}

	if f.IsInf(0) {
		if f.Signbit() {
			return float32(math.Inf(-1))
		}
		return float32(math.Inf(1))
	}

	// Extract components
	sign, exp16, mant16 := f.extractComponents()

	if exp16 == 0 {
		// Subnormal number
		if mant16 == 0 {
			// Zero (already handled above, but for completeness)
			if sign != 0 {
				return math.Float32frombits(0x80000000)
			}
			return 0.0
		}

		// For subnormal numbers, the value is: sign * 0.mantissa * 2^-14
		// We need to convert this to a normalized float32: sign * 1.mantissa * 2^e
		// The smallest positive subnormal is 2^-24 (0x0001 = 2^-14 * 2^-10)
		// The largest subnormal is just under 2^-14 (0x03FF = (1-2^-10) * 2^-14)

		// Handle the case where mantissa is zero (0.0 or -0.0)
		if mant16 == 0 {
			if sign != 0 {
				return math.Float32frombits(0x80000000) // -0.0
			}
			return 0.0 // +0.0
		}

		// For subnormal numbers, we need to normalize the mantissa
		// The mantissa is in the range [0x001, 0x3FF] for subnormals
		// We need to find the position of the leading 1 bit

		// Count leading zeros in the 10-bit mantissa
		leadingZeros := leadingZeros10(mant16)
		if leadingZeros < 0 || leadingZeros > 9 {
			// Should never happen due to leadingZeros10 implementation, but be defensive
			if sign != 0 {
				return math.Float32frombits(0x80000000) // -0.0
			}
			return 0.0 // +0.0
		}

		// The number of positions to shift left to normalize (1 to 10)
		shift := leadingZeros + 1

		// Shift the mantissa left to normalize it (make the leading 1 explicit)
		// For example, for 0x0001 (2^-24):
		//   mant16 = 0x0001 = 0b0000000001
		//   leadingZeros = 9, shift = 10
		//   mant16 <<= 10 = 0x0400 = 0b10000000000 (11 bits, but we keep only 10)
		mant16 <<= shift
		// Keep only the 10 LSBs (mantissa part)
		mant16 &= 0x3FF

		// For subnormal numbers, the exponent is -14 (1 - ExponentBias)
		// After normalization, we need to adjust the exponent by (shift - 1)
		// So the final exponent is: -14 - (shift - 1) = -15 + shift
		// Then we add the float32 bias (127) to get the biased exponent
		// For 0x0001: exp32 = 127 - 15 + 10 = 122 (which is correct for 2^-24)
		exp32 := int32(Float32ExponentBias - 15 + shift)
		if exp32 <= 0 || exp32 >= 255 {
			// Underflow to zero or overflow to infinity
			if exp32 >= 255 {
				// Infinity
				if sign != 0 {
					return float32(math.Inf(-1))
				}
				return float32(math.Inf(1))
			}
			// Zero
			if sign != 0 {
				return math.Float32frombits(0x80000000)
			}
			return 0.0
		}

		// Shift mantissa to float32 position
		mant32 := uint32(mant16) << (Float32MantissaLen - MantissaLen)

		// Combine into IEEE 754 float32
		bits := (uint32(sign) << 31) | (uint32(exp32) << 23) | mant32
		return math.Float32frombits(bits)
	}

	// Normal number
	// Convert exponent from float16 bias to float32 bias
	exp32 := int32(exp16) - ExponentBias + Float32ExponentBias

	// Shift mantissa to float32 position
	mant32 := uint32(mant16) << (Float32MantissaLen - MantissaLen)

	// Combine into IEEE 754 float32
	bits := (uint32(sign) << 31) | (uint32(exp32) << 23) | mant32
	return math.Float32frombits(bits)
}

// ToFloat64 converts a Float16 value to float64 with full precision
func (f Float16) ToFloat64() float64 {
	// Handle special cases
	if f.IsZero() {
		if f.Signbit() {
			return math.Copysign(0.0, -1.0)
		}
		return 0.0
	}

	if f.IsNaN() {
		sign := uint64(0)
		if f.Signbit() {
			sign = 0x8000000000000000
		}
		payload := uint64(f & MantissaMask)
		return math.Float64frombits(sign | 0x7FF8000000000000 | (payload << (Float64MantissaLen - MantissaLen)))
	}

	if f.IsInf(0) {
		if f.Signbit() {
			return math.Inf(-1)
		}
		return math.Inf(1)
	}

	// Extract components
	sign, exp16, mant16 := f.extractComponents()

	if exp16 == 0 { // Subnormal
		// val = sign * 0.mantissa * 2^-14
		// smallest subnormal: 1 * 2^-10 * 2^-14 = 2^-24
		// largest subnormal: (1023/1024) * 2^-14
		val := float64(mant16) * math.Pow(2, -24)
		if sign != 0 {
			return -val
		}
		return val
	}

	// Normal number
	exp64 := int64(exp16) - ExponentBias + Float64ExponentBias
	mant64 := uint64(mant16) << (Float64MantissaLen - MantissaLen)

	bits := (uint64(sign) << 63) | (uint64(exp64) << 52) | mant64
	return math.Float64frombits(bits)
}

// FromFloat32 converts a float32 to Float16 (with potential precision loss)
func FromFloat32(f32 float32) Float16 {
	// Handle special cases first
	if f32 == 0.0 {
		if math.Signbit(float64(f32)) {
			return NegativeZero
		}
		return PositiveZero
	}

	if math.IsInf(float64(f32), 0) {
		if f32 > 0 {
			return PositiveInfinity
		}
		return NegativeInfinity
	}

	if math.IsNaN(float64(f32)) {
		if math.Signbit(float64(f32)) {
			return NegativeQNaN
		}
		return QuietNaN
	}

	// Extract bits from float32
	bits := math.Float32bits(f32)
	sign := uint16((bits >> 31) & 0x1)
	exp32 := int16((bits >> 23) & 0xFF)
	mant32 := bits & 0x007FFFFF // 23-bit mantissa

	// Handle subnormal float32 values
	if exp32 == 0 && mant32 != 0 {
		// For very small subnormal numbers, we might need to return the smallest subnormal or zero
		if f32 < math.SmallestNonzeroFloat32/2 {
			if sign != 0 {
				return NegativeZero
			}
			return PositiveZero
		}

		// For subnormals, we'll convert them to the smallest subnormal in float16
		// or zero, depending on their magnitude
		if f32 < math.SmallestNonzeroFloat32/1024 {
			if sign != 0 {
				return 0x8001 // Negative smallest subnormal
			}
			return 0x0001 // Smallest positive subnormal
		}

		// For other subnormals, we'll scale them to the float16 subnormal range
		scaled := f32 * float32(1<<10) // Scale up to get more precision
		return FromFloat32(scaled)
	}

	// Convert exponent from float32 bias (127) to float16 bias (15)
	exp16 := int32(exp32) - 127 + 15

	// Check for overflow/underflow
	if exp16 >= 0x1F {
		// Overflow - return infinity with correct sign
		return Float16((sign << 15) | 0x7C00)
	}

	// Handle normal numbers
	if exp16 > 0 {
		// Normal number - extract top 10 bits of mantissa with rounding
		mant16 := uint16((mant32 >> (23 - 10)) & 0x3FF)
		roundBit := (mant32 >> (23 - 10 - 1)) & 0x1
		mant16 += uint16(roundBit)

		// Check for mantissa overflow (due to rounding)
		if (mant16 & 0x400) != 0 {
			mant16 >>= 1
			exp16++
			// Check for exponent overflow after rounding
			if exp16 >= 0x1F {
				return Float16((sign << 15) | 0x7C00)
			}
		}

		// Combine sign, exponent, and mantissa
		return Float16((sign << 15) | (uint16(exp16) << 10) | (mant16 & 0x3FF))
	}

	// Handle underflow - convert to subnormal or flush to zero
	shift := 1 - exp16    // Number of bits to shift right
	if shift > 10+1+127 { // 10 mantissa bits + 1 for the implicit leading 1 + 127 for float32 exponent range
		// Too small to represent, flush to zero
		if sign != 0 {
			return NegativeZero
		}
		return PositiveZero
	}

	// Convert to subnormal
	// Add the implicit leading 1 and shift to get the mantissa
	mant16 := uint16((0x800000|mant32)>>(shift+23-10-1)) >> 1

	// For very small numbers, ensure we don't lose all precision
	if mant16 == 0 {
		mant16 = 1 // Smallest subnormal
	}

	return Float16((sign << 15) | (mant16 & 0x3FF))
}

// FromFloat64 converts a float64 to Float16 (with potential precision loss)
func FromFloat64(f64 float64) Float16 {
	// Handle special cases first
	if f64 == 0.0 {
		if math.Signbit(f64) {
			return NegativeZero
		}
		return PositiveZero
	}

	if math.IsInf(f64, 0) {
		if f64 > 0 {
			return PositiveInfinity
		}
		return NegativeInfinity
	}

	if math.IsNaN(f64) {
		if math.Signbit(f64) {
			return NegativeQNaN
		}
		return QuietNaN
	}

	// Extract bits from float64
	bits := math.Float64bits(f64)
	sign := uint16((bits >> 63) & 0x1)
	exp64 := int16((bits >> 52) & 0x7FF)
	mant64 := bits & 0x000F_FFFF_FFFF_FFFF // 52-bit mantissa

	// Handle subnormal float64 values
	if exp64 == 0 && mant64 != 0 {
		// For very small subnormal numbers, we might need to return the smallest subnormal
		// or zero, depending on the value
		if f64 < math.SmallestNonzeroFloat32 {
			// This is smaller than the smallest float32 subnormal
			if sign != 0 {
				return NegativeZero
			}
			return PositiveZero
		}

		// Convert through float32 for better handling of subnormals
		return FromFloat32(float32(f64))
	}

	// Convert exponent from float64 bias (1023) to float16 bias (15)
	exp16 := int32(exp64) - 1023 + 15

	// Check for overflow/underflow
	if exp16 >= 0x1F {
		// Overflow - return infinity with correct sign
		return Float16((sign << 15) | 0x7C00)
	}

	// Handle normal numbers
	if exp16 > 0 {
		// Normal number - extract top 10 bits of mantissa with rounding
		// Add 1 to the 11th bit for rounding
		mant16 := uint16((mant64 >> (52 - 10)) & 0x3FF)
		roundBit := (mant64 >> (52 - 10 - 1)) & 0x1
		mant16 += uint16(roundBit)

		// Check for mantissa overflow (due to rounding)
		if (mant16 & 0x400) != 0 {
			mant16 >>= 1
			exp16++
			// Check for exponent overflow after rounding
			if exp16 >= 0x1F {
				return Float16((sign << 15) | 0x7C00)
			}
		}

		// Combine sign, exponent, and mantissa
		return Float16((sign << 15) | (uint16(exp16) << 10) | (mant16 & 0x3FF))
	}

	// Handle underflow - convert to subnormal or flush to zero
	shift := 1 - exp16     // Number of bits to shift right
	if shift > 10+1+1022 { // 10 mantissa bits + 1 for the implicit leading 1 + 1022 for float64 exponent range
		// Too small to represent, flush to zero
		if sign != 0 {
			return NegativeZero
		}
		return PositiveZero
	}

	// Convert to subnormal
	// Add the implicit leading 1 and shift to get the mantissa
	mant16 := uint16((0x8000000000000|mant64)>>(shift+52-10-1)) >> 1

	// For very small numbers, ensure we don't lose all precision
	if mant16 == 0 {
		mant16 = 1 // Smallest subnormal
	}

	return Float16((sign << 15) | (mant16 & 0x3FF))
}

// FromFloat64WithMode converts a float64 to Float16 with specified modes
func FromFloat64WithMode(f64 float64, convMode ConversionMode, roundMode RoundingMode) (Float16, error) {
	// Handle special cases first
	if f64 == 0.0 {
		if math.Signbit(f64) {
			return NegativeZero, nil
		}
		return PositiveZero, nil
	}

	// In strict mode, check for special values that might be disallowed
	if convMode == ModeStrict {
		if math.IsInf(f64, 0) {
			return 0, &Float16Error{
				Op:    "FromFloat64WithMode",
				Value: float32(f64),
				Msg:   "infinity not allowed in strict mode",
				Code:  ErrInfinity,
			}
		}
		if math.IsNaN(f64) {
			return 0, &Float16Error{
				Op:    "FromFloat64WithMode",
				Value: float32(f64),
				Msg:   "NaN not allowed in strict mode",
				Code:  ErrNaN,
			}
		}
		// Check for overflow/underflow in strict mode
		if f64 > 65504.0 || f64 < -65504.0 {
			return 0, &Float16Error{
				Op:    "FromFloat64WithMode",
				Value: float32(f64),
				Msg:   "value out of range for float16",
				Code:  ErrOverflow,
			}
		}
		// Check for underflow in strict mode (smaller than smallest normal)
		if f64 != 0 && math.Abs(f64) < 6.103515625e-05 { // 2^-14
			return 0, &Float16Error{
				Op:    "FromFloat64WithMode",
				Value: float32(f64),
				Msg:   "value underflows float16 in strict mode",
				Code:  ErrUnderflow,
			}
		}
	}

	// Extract bits from float64
	bits := math.Float64bits(f64)
	sign := uint16((bits >> 63) & 0x1)
	exp64 := int16((bits >> 52) & 0x7FF)
	mant64 := bits & 0x000F_FFFF_FFFF_FFFF // 52-bit mantissa

	// Handle subnormal float64 values
	if exp64 == 0 && mant64 != 0 {
		exp64 = 1 - 1023 // Subnormal exponent for float64
		for (mant64 & (1 << 52)) == 0 {
			mant64 <<= 1
			exp64--
		}
	}

	// Convert exponent from float64 bias (1023) to float16 bias (15)
	exp16 := int32(exp64) - 1023 + 15

	// Check for overflow/underflow (for non-strict mode)
	if exp16 >= 0x1F {
		// Overflow - return infinity with correct sign
		if convMode == ModeStrict {
			return 0, &Float16Error{
				Op:    "FromFloat64WithMode",
				Value: float32(f64),
				Msg:   "value overflows float16 in strict mode",
				Code:  ErrOverflow,
			}
		}
		return Float16((sign << 15) | 0x7C00), nil
	}

	// Handle normal numbers
	if exp16 > 0 {
		// Extract mantissa bits with guard and round bits
		mant16 := uint16((mant64 >> (52 - 10)) & 0x3FF)
		guardBit := (mant64 >> (52 - 10 - 1)) & 0x1
		roundBit := (mant64 >> (52 - 10 - 2)) & 0x1
		stickyBit := uint16(0)
		if (mant64 & ((1 << (52 - 10 - 2)) - 1)) != 0 {
			stickyBit = 1
		}

		// Apply rounding based on mode
		roundUp := false
		switch roundMode {
		case RoundNearestEven: // Round to nearest, ties to even (IEEE default)
			if guardBit == 1 && (roundBit == 1 || stickyBit == 1 || (mant16&0x1) == 1) {
				roundUp = true
			}
		case RoundTowardZero: // Truncate toward zero
			// No rounding needed
		case RoundTowardPositive: // Round toward +∞
			if sign == 0 && (guardBit == 1 || roundBit == 1 || stickyBit == 1) {
				roundUp = true
			}
		case RoundTowardNegative: // Round toward -∞
			if sign == 1 && (guardBit == 1 || roundBit == 1 || stickyBit == 1) {
				roundUp = true
			}
		case RoundNearestAway: // Round to nearest, ties away from zero
			if guardBit == 1 && (roundBit == 1 || stickyBit == 1) {
				roundUp = true
			}
		}

		if roundUp {
			mant16++
			// Check for mantissa overflow (due to rounding)
			if (mant16 & 0x400) != 0 {
				mant16 >>= 1
				exp16++
				// Check for exponent overflow after rounding
				if exp16 >= 0x1F {
					if convMode == ModeStrict {
						return 0, &Float16Error{
							Op:    "FromFloat64WithMode",
							Value: float32(f64),
							Msg:   "value overflows float16 in strict mode after rounding",
							Code:  ErrOverflow,
						}
					}
					return Float16((sign << 15) | 0x7C00), nil
				}
			}
		}

		// Combine sign, exponent, and mantissa
		return Float16((sign << 15) | (uint16(exp16) << 10) | (mant16 & 0x3FF)), nil
	}

	// Handle underflow - convert to subnormal or flush to zero
	shift := 1 - exp16 // Number of bits to shift right
	if shift > 10+1 {  // 10 mantissa bits + 1 for the implicit leading 1
		// Too small to represent, flush to zero
		if sign != 0 {
			return NegativeZero, nil
		}
		return PositiveZero, nil
	}

	// Convert to subnormal
	// Add the implicit leading 1 and shift to get the mantissa
	mant16 := uint16((0x8000000000000 | mant64) >> (shift + 52 - 10))
	// Get guard and round bits
	guardBit := (mant64 >> (shift + 52 - 10 - 1)) & 0x1
	roundBit := (mant64 >> (shift + 52 - 10 - 2)) & 0x1
	stickyBit := uint16(0)
	if (mant64 & ((1 << (shift + 52 - 10 - 2)) - 1)) != 0 {
		stickyBit = 1
	}

	// Apply rounding based on mode
	roundUp := false
	switch roundMode {
	case RoundNearestEven: // Round to nearest, ties to even (IEEE default)
		if guardBit == 1 && (roundBit == 1 || stickyBit == 1 || (mant16&0x1) == 1) {
			roundUp = true
		}
	case RoundTowardZero: // Truncate toward zero
		// No rounding needed
	case RoundTowardPositive: // Round toward +∞
		if sign == 0 && (guardBit == 1 || roundBit == 1 || stickyBit == 1) {
			roundUp = true
		}
	case RoundTowardNegative: // Round toward -∞
		if sign == 1 && (guardBit == 1 || roundBit == 1 || stickyBit == 1) {
			roundUp = true
		}
	case RoundNearestAway: // Round to nearest, ties away from zero
		if guardBit == 1 && (roundBit == 1 || stickyBit == 1) {
			roundUp = true
		}
	}

	if roundUp {
		mant16++
	}

	// Check if rounding caused overflow to normal number
	if (mant16 & 0x400) != 0 {
		return Float16((sign << 15) | 0x0400), nil // Smallest normal number
	}

	return Float16((sign << 15) | (mant16 & 0x3FF)), nil
}

// Batch conversion functions optimized for performance

// ToSlice16 converts a slice of float32 to Float16 with optimized performance
func ToSlice16(f32s []float32) []Float16 {
	if len(f32s) == 0 {
		return nil
	}

	result := make([]Float16, len(f32s))

	// Use unsafe pointer arithmetic for better performance
	// This avoids bounds checking in the inner loop
	src := (*float32)(unsafe.Pointer(&f32s[0]))
	dst := (*Float16)(unsafe.Pointer(&result[0]))

	for i := 0; i < len(f32s); i++ {
		srcPtr := (*float32)(unsafe.Pointer(uintptr(unsafe.Pointer(src)) + uintptr(i)*unsafe.Sizeof(float32(0))))
		dstPtr := (*Float16)(unsafe.Pointer(uintptr(unsafe.Pointer(dst)) + uintptr(i)*unsafe.Sizeof(Float16(0))))
		*dstPtr = ToFloat16(*srcPtr)
	}

	return result
}

// ToSlice32 converts a slice of Float16 to float32 with optimized performance
func ToSlice32(f16s []Float16) []float32 {
	if len(f16s) == 0 {
		return nil
	}

	result := make([]float32, len(f16s))

	// Use unsafe pointer arithmetic for better performance
	src := (*Float16)(unsafe.Pointer(&f16s[0]))
	dst := (*float32)(unsafe.Pointer(&result[0]))

	for i := 0; i < len(f16s); i++ {
		srcPtr := (*Float16)(unsafe.Pointer(uintptr(unsafe.Pointer(src)) + uintptr(i)*unsafe.Sizeof(Float16(0))))
		dstPtr := (*float32)(unsafe.Pointer(uintptr(unsafe.Pointer(dst)) + uintptr(i)*unsafe.Sizeof(float32(0))))
		*dstPtr = (*srcPtr).ToFloat32()
	}

	return result
}

// ToSlice64 converts a slice of Float16 to float64 with optimized performance
func ToSlice64(f16s []Float16) []float64 {
	if len(f16s) == 0 {
		return nil
	}

	result := make([]float64, len(f16s))

	for i, f16 := range f16s {
		result[i] = f16.ToFloat64()
	}

	return result
}

// FromSlice64 converts a slice of float64 to Float16 with optimized performance
func FromSlice64(f64s []float64) []Float16 {
	if len(f64s) == 0 {
		return nil
	}

	result := make([]Float16, len(f64s))

	for i, f64 := range f64s {
		result[i] = FromFloat64(f64)
	}

	return result
}

// SIMD-friendly batch conversion with error handling
// ToSlice16WithMode converts a slice with specified conversion mode
func ToSlice16WithMode(f32s []float32, convMode ConversionMode, roundMode RoundingMode) ([]Float16, []error) {
	if len(f32s) == 0 {
		return nil, nil
	}

	result := make([]Float16, len(f32s))
	var errors []error

	for i, f32 := range f32s {
		f16, err := ToFloat16WithMode(f32, convMode, roundMode)
		result[i] = f16
		if err != nil {
			if errors == nil {
				errors = make([]error, 0, len(f32s))
			}
			// Store error with index information
			indexedErr := &Float16Error{
				Op:    fmt.Sprintf("convert[%d]", i),
				Value: f32,
				Msg:   err.Error(),
				Code:  err.(*Float16Error).Code,
			}
			errors = append(errors, indexedErr)
		}
	}

	return result, errors
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
