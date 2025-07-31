package float16

import (
	"fmt"
	"math"
)

// Global arithmetic settings
var (
	DefaultArithmeticMode = ModeIEEEArithmetic
	DefaultRounding       = DefaultRoundingMode
)

// ArithmeticMode defines the precision/performance trade-off for arithmetic operations
type ArithmeticMode int

const (
	// ModeIEEE provides full IEEE 754 compliance with proper rounding
	ModeIEEEArithmetic ArithmeticMode = iota
	// ModeFastArithmetic optimizes for speed, may sacrifice some precision
	ModeFastArithmetic
	// ModeExactArithmetic provides exact results when possible, errors on precision loss
	ModeExactArithmetic
)

// Add performs addition of two Float16 values
func Add(a, b Float16) Float16 {
	result, _ := AddWithMode(a, b, DefaultArithmeticMode, DefaultRounding)
	return result
}

// AddWithMode performs addition with specified arithmetic and rounding modes
func AddWithMode(a, b Float16, mode ArithmeticMode, rounding RoundingMode) (Float16, error) {
	// Handle special cases first for performance
	if a.IsZero() {
		return b, nil
	}
	if b.IsZero() {
		return a, nil
	}

	// Handle NaN cases
	if a.IsNaN() || b.IsNaN() {
		if mode == ModeExactArithmetic {
			return 0, &Float16Error{
				Op:   "add",
				Msg:  "NaN operand in exact mode",
				Code: ErrNaN,
			}
		}
		// Return a quiet NaN
		return QuietNaN, nil
	}

	// Handle infinity cases
	if a.IsInf(0) || b.IsInf(0) {
		if a.IsInf(1) && b.IsInf(-1) {
			// +∞ + (-∞) = NaN
			if mode == ModeExactArithmetic {
				return 0, &Float16Error{
					Op:   "add",
					Msg:  "infinity - infinity is undefined",
					Code: ErrInvalidOperation,
				}
			}
			return QuietNaN, nil
		}
		if a.IsInf(-1) && b.IsInf(1) {
			// (-∞) + (+∞) = NaN
			if mode == ModeExactArithmetic {
				return 0, &Float16Error{
					Op:   "add",
					Msg:  "infinity - infinity is undefined",
					Code: ErrInvalidOperation,
				}
			}
			return QuietNaN, nil
		}
		// Return the infinity
		if a.IsInf(0) {
			return a, nil
		}
		return b, nil
	}

	// For high performance, convert to float32, compute, and convert back
	// This approach is faster than implementing full IEEE 754 arithmetic in float16
	if mode == ModeFastArithmetic {
		f32a := a.ToFloat32()
		f32b := b.ToFloat32()
		result := f32a + f32b
		return NewConverter(ModeIEEE, rounding).ToFloat16(result), nil
	}

	// Full IEEE 754 implementation for exact mode
	return addIEEE754(a, b, rounding)
}

// Sub performs subtraction of two Float16 values
func Sub(a, b Float16) Float16 {
	result, _ := SubWithMode(a, b, DefaultArithmeticMode, DefaultRounding)
	return result
}

// SubWithMode performs subtraction with specified arithmetic and rounding modes
func SubWithMode(a, b Float16, mode ArithmeticMode, rounding RoundingMode) (Float16, error) {
	// Subtraction is addition with negated second operand
	return AddWithMode(a, b.Neg(), mode, rounding)
}

// Mul performs multiplication of two Float16 values
func Mul(a, b Float16) Float16 {
	result, _ := MulWithMode(a, b, DefaultArithmeticMode, DefaultRounding)
	return result
}

// MulWithMode performs multiplication with specified arithmetic and rounding modes
func MulWithMode(a, b Float16, mode ArithmeticMode, rounding RoundingMode) (Float16, error) {
	// Handle special cases
	// Check for zero times infinity cases first
	aIsZero := a.IsZero()
	bIsInf := b.IsInf(0)
	if (aIsZero && bIsInf) || (a.IsInf(0) && b.IsZero()) {
		// 0 * ∞ = NaN
		if mode == ModeExactArithmetic {
			return 0, &Float16Error{
				Op:   "mul",
				Msg:  "zero times infinity is undefined",
				Code: ErrInvalidOperation,
			}
		}
		return QuietNaN, nil
	}

	// Handle zero cases
	if aIsZero || b.IsZero() {
		// Handle sign of zero result: 0 * anything = ±0
		signA := a.Signbit()
		signB := b.Signbit()
		if signA != signB {
			return NegativeZero, nil
		}
		return PositiveZero, nil
	}

	// Handle NaN cases
	if a.IsNaN() || b.IsNaN() {
		if mode == ModeExactArithmetic {
			return 0, &Float16Error{
				Op:   "mul",
				Msg:  "NaN operand in exact mode",
				Code: ErrNaN,
			}
		}
		return QuietNaN, nil
	}

	// Handle infinity cases
	if a.IsInf(0) || b.IsInf(0) {
		// Check for 0 * ∞ which is NaN
		if (a.IsInf(0) && b.IsZero()) || (a.IsZero() && b.IsInf(0)) {
			if mode == ModeExactArithmetic {
				return 0, &Float16Error{
					Op:   "mul",
					Msg:  "zero times infinity is undefined",
					Code: ErrInvalidOperation,
				}
			}
			return QuietNaN, nil
		}

		// ∞ * finite = ±∞ (sign depends on operand signs)
		signA := a.Signbit()
		signB := b.Signbit()
		if signA != signB {
			return NegativeInfinity, nil
		}
		return PositiveInfinity, nil
	}

	// For high performance, use float32 arithmetic
	if mode == ModeFastArithmetic {
		f32a := a.ToFloat32()
		f32b := b.ToFloat32()
		result := f32a * f32b
		return NewConverter(ModeIEEE, rounding).ToFloat16(result), nil
	}

	// Full IEEE 754 implementation
	return addIEEE754(a, b, rounding)
}

// Div performs division of two Float16 values
func Div(a, b Float16) Float16 {
	result, _ := DivWithMode(a, b, DefaultArithmeticMode, DefaultRounding)
	return result
}

// DivWithMode performs division with specified arithmetic and rounding modes
func DivWithMode(a, b Float16, mode ArithmeticMode, rounding RoundingMode) (Float16, error) {
	// Handle division by zero
	if b.IsZero() {
		if a.IsZero() {
			// 0/0 = NaN
			if mode == ModeExactArithmetic {
				return 0, &Float16Error{
					Op:   "div",
					Msg:  "zero divided by zero is undefined",
					Code: ErrInvalidOperation,
				}
			}
			return QuietNaN, nil
		}
		// finite/0 = ±∞
		if mode == ModeExactArithmetic {
			return 0, &Float16Error{
				Op:   "div",
				Msg:  "division by zero",
				Code: ErrDivisionByZero,
			}
		}
		signA := a.Signbit()
		signB := b.Signbit()
		if signA != signB {
			return NegativeInfinity, nil
		}
		return PositiveInfinity, nil
	}

	// Handle zero dividend
	if a.IsZero() {
		// 0/finite = ±0
		signA := a.Signbit()
		signB := b.Signbit()
		if signA != signB {
			return NegativeZero, nil
		}
		return PositiveZero, nil
	}

	// Handle infinity cases
	if a.IsInf(0) || b.IsInf(0) {
		if a.IsInf(0) && b.IsInf(0) {
			// ∞/∞ = NaN
			if mode == ModeExactArithmetic {
				return 0, &Float16Error{
					Op:   "div",
					Msg:  "infinity divided by infinity is undefined",
					Code: ErrInvalidOperation,
				}
			}
			return QuietNaN, nil
		}

		if a.IsInf(0) {
			// ∞/finite = ±∞
			signA := a.Signbit()
			signB := b.Signbit()
			if signA != signB {
				return NegativeInfinity, nil
			}
			return PositiveInfinity, nil
		}

		// finite/∞ = ±0
		signA := a.Signbit()
		signB := b.Signbit()
		if signA != signB {
			return NegativeZero, nil
		}
		return PositiveZero, nil
	}

	// Handle NaN cases
	if a.IsNaN() || b.IsNaN() {
		if mode == ModeExactArithmetic {
			return 0, &Float16Error{
				Op:   "div",
				Msg:  "NaN operand in exact mode",
				Code: ErrNaN,
			}
		}
		return QuietNaN, nil
	}

	// Handle infinity cases
	if a.IsInf(0) && b.IsInf(0) {
		// ∞/∞ = NaN
		if mode == ModeExactArithmetic {
			return 0, &Float16Error{
				Op:   "div",
				Msg:  "infinity divided by infinity is undefined",
				Code: ErrInvalidOperation,
			}
		}
		return QuietNaN, nil
	}

	if a.IsInf(0) {
		// ∞/finite = ±∞
		signA := a.Signbit()
		signB := b.Signbit()
		if signA != signB {
			return NegativeInfinity, nil
		}
		return PositiveInfinity, nil
	}

	if b.IsInf(0) {
		// finite/∞ = ±0
		signA := a.Signbit()
		signB := b.Signbit()
		if signA != signB {
			return NegativeZero, nil
		}
		return PositiveZero, nil
	}

	// For high performance, use float32 arithmetic
	if mode == ModeFastArithmetic {
		f32a := a.ToFloat32()
		f32b := b.ToFloat32()
		result := f32a / f32b
		return NewConverter(ModeIEEE, rounding).ToFloat16(result), nil
	}

	// Full IEEE 754 implementation
	return addIEEE754(a, b, rounding)
}

// IEEE 754 compliant arithmetic implementations

// addIEEE754 implements full IEEE 754 addition
func addIEEE754(a, b Float16, rounding RoundingMode) (Float16, error) {
	// For addition, we can use the simpler approach of converting to float32
	// since the intermediate precision is sufficient for exact float16 results
	f32a := a.ToFloat32()
	f32b := b.ToFloat32()
	result := f32a + f32b
	return NewConverter(ModeIEEE, rounding).ToFloat16WithMode(result)
}

// mulIEEE754 implements full IEEE 754 multiplication
func mulIEEE754(a, b Float16, rounding RoundingMode) (Float16, error) {
	// For multiplication, we can use the simpler approach of converting to float32
	// since the intermediate precision is sufficient for exact float16 results
	f32a := a.ToFloat32()
	f32b := b.ToFloat32()
	result := f32a * f32b
	return NewConverter(ModeIEEE, rounding).ToFloat16WithMode(result)
}

// divIEEE754 implements full IEEE 754 division
func divIEEE754(a, b Float16, rounding RoundingMode) (Float16, error) {
	// For division, we can use the simpler approach of converting to float32
	// since the intermediate precision is sufficient for exact float16 results
	f32a := a.ToFloat32()
	f32b := b.ToFloat32()
	result := f32a / f32b

	// Use the provided rounding mode for the conversion back to Float16
	return NewConverter(ModeExact, rounding).ToFloat16WithMode(result)
}

// Comparison operations

// Equal returns true if two Float16 values are equal
func Equal(a, b Float16) bool {
	// Handle NaN: NaN != NaN
	if a.IsNaN() || b.IsNaN() {
		return false
	}
	// Handle zero: +0 == -0
	if a.IsZero() && b.IsZero() {
		return true
	}
	return a == b
}

// Less returns true if a < b
func Less(a, b Float16) bool {
	// Handle NaN: any comparison with NaN is false
	if a.IsNaN() || b.IsNaN() {
		return false
	}

	// Handle zero: -0 == +0 for comparison
	if a.IsZero() && b.IsZero() {
		return false
	}

	// Handle signs
	signA := a.Signbit()
	signB := b.Signbit()

	if signA && !signB {
		return true // negative < positive
	}
	if !signA && signB {
		return false // positive > negative
	}

	// Same sign: compare magnitudes
	if signA {
		// Both negative: larger magnitude is smaller value
		return a > b
	} else {
		// Both positive: smaller magnitude is smaller value
		return a < b
	}
}

// Greater returns true if a > b
func Greater(a, b Float16) bool {
	return Less(b, a)
}

// LessEqual returns true if a <= b
func LessEqual(a, b Float16) bool {
	return Less(a, b) || Equal(a, b)
}

// GreaterEqual returns true if a >= b
func GreaterEqual(a, b Float16) bool {
	return Greater(a, b) || Equal(a, b)
}

// Min returns the smaller of two Float16 values
func Min(a, b Float16) Float16 {
	// Handle NaN: return the non-NaN value, or NaN if both are NaN
	if a.IsNaN() {
		return b
	}
	if b.IsNaN() {
		return a
	}
	// Handle -0 and +0
	if a.IsZero() && b.IsZero() {
		if a.Signbit() {
			return a // a is -0
		}
		return b // b is -0, or both are +0
	}
	if Less(a, b) {
		return a
	}
	return b
}

// Max returns the larger of two Float16 values
func Max(a, b Float16) Float16 {
	// Handle NaN: return the non-NaN value, or NaN if both are NaN
	if a.IsNaN() {
		return b
	}
	if b.IsNaN() {
		return a
	}

	if Greater(a, b) {
		return a
	}
	return b
}

// Batch operations for high-performance computing

// AddSlice performs element-wise addition of two Float16 slices
func AddSlice(a, b []Float16) []Float16 {
	if len(a) != len(b) {
		panic("float16: slice length mismatch")
	}

	result := make([]Float16, len(a))
	for i := range a {
		result[i] = Add(a[i], b[i])
	}
	return result
}

// SubSlice performs element-wise subtraction of two Float16 slices
func SubSlice(a, b []Float16) []Float16 {
	if len(a) != len(b) {
		panic("float16: slice length mismatch")
	}

	result := make([]Float16, len(a))
	for i := range a {
		result[i] = Sub(a[i], b[i])
	}
	return result
}

// MulSlice performs element-wise multiplication of two Float16 slices
func MulSlice(a, b []Float16) []Float16 {
	if len(a) != len(b) {
		panic("float16: slice length mismatch")
	}

	result := make([]Float16, len(a))
	for i := range a {
		product := Mul(a[i], b[i])
		result[i] = product
		// Debug print
		fmt.Printf("MulSlice: a[%d]=%v (0x%04X), b[%d]=%v (0x%04X), product=%v (0x%04X)\n", i, a[i], uint16(a[i]), i, b[i], uint16(b[i]), product, uint16(product))
	}
	fmt.Printf("MulSlice: result=%v\n", result)
	return result
}

// DivSlice performs element-wise division of two Float16 slices
func DivSlice(a, b []Float16) []Float16 {
	if len(a) != len(b) {
		panic("float16: slice length mismatch")
	}

	result := make([]Float16, len(a))
	for i := range a {
		result[i] = Div(a[i], b[i])
	}
	return result
}

// ScaleSlice multiplies each element in the slice by a scalar
func ScaleSlice(s []Float16, scalar Float16) []Float16 {
	result := make([]Float16, len(s))
	for i := range s {
		result[i] = Mul(s[i], scalar)
	}
	return result
}

// SumSlice returns the sum of all elements in the slice
func SumSlice(s []Float16) Float16 {
	var sum Float16 = PositiveZero
	for _, v := range s {
		sum = Add(sum, v)
	}
	return sum
}

// DotProduct computes the dot product of two Float16 slices
func DotProduct(a, b []Float16) Float16 {
	if len(a) != len(b) {
		panic("float16: slice length mismatch")
	}

	var sum Float16 = PositiveZero
	for i := range a {
		product := Mul(a[i], b[i])
		sum = Add(sum, product)
	}
	return sum
}

// Norm2 computes the L2 norm (Euclidean norm) of a Float16 slice
func Norm2(s []Float16) Float16 {
	var sumSquares Float16 = PositiveZero
	for _, v := range s {
		square := Mul(v, v)
		sumSquares = Add(sumSquares, square)
	}
	return NewConverter(DefaultConversionMode, DefaultRoundingMode).FromFloat64(math.Sqrt(sumSquares.ToFloat64()))
}
