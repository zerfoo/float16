package float16

import "math"

// BFloat16AddWithMode performs addition with specified arithmetic and rounding modes.
func BFloat16AddWithMode(a, b BFloat16, mode ArithmeticMode, rounding RoundingMode) (BFloat16, error) {
	// Handle NaN propagation: if either operand is NaN, propagate it
	if a.IsNaN() || b.IsNaN() {
		if mode == ModeExactArithmetic {
			return 0, &Float16Error{Op: "bfloat16_add", Msg: "NaN operand in exact mode", Code: ErrNaN}
		}
		return BFloat16QuietNaN, nil
	}

	// Handle zeros
	if a.IsZero() {
		return b, nil
	}
	if b.IsZero() {
		return a, nil
	}

	// Handle infinity cases
	if a.IsInf(0) || b.IsInf(0) {
		if a.IsInf(1) && b.IsInf(-1) || a.IsInf(-1) && b.IsInf(1) {
			if mode == ModeExactArithmetic {
				return 0, &Float16Error{Op: "bfloat16_add", Msg: "infinity - infinity is undefined", Code: ErrInvalidOperation}
			}
			return BFloat16QuietNaN, nil
		}
		if a.IsInf(0) {
			return a, nil
		}
		return b, nil
	}

	if mode == ModeFastArithmetic {
		return BFloat16FromFloat32(a.ToFloat32() + b.ToFloat32()), nil
	}

	// IEEE mode: compute in float32 with specified rounding, handle gradual underflow
	result := a.ToFloat32() + b.ToFloat32()
	bf := BFloat16FromFloat32WithRounding(result, rounding)

	// Gradual underflow: if the float32 result is non-zero but rounds to BFloat16 zero,
	// return the smallest subnormal with the correct sign instead.
	if result != 0 && bf.IsZero() {
		if result > 0 {
			return BFloat16SmallestPosSubnormal, nil
		}
		return BFloat16SmallestNegSubnormal, nil
	}

	return bf, nil
}

// BFloat16SubWithMode performs subtraction with specified arithmetic and rounding modes.
func BFloat16SubWithMode(a, b BFloat16, mode ArithmeticMode, rounding RoundingMode) (BFloat16, error) {
	return BFloat16AddWithMode(a, BFloat16Neg(b), mode, rounding)
}

// BFloat16MulWithMode performs multiplication with specified arithmetic and rounding modes.
func BFloat16MulWithMode(a, b BFloat16, mode ArithmeticMode, rounding RoundingMode) (BFloat16, error) {
	// NaN propagation
	if a.IsNaN() || b.IsNaN() {
		if mode == ModeExactArithmetic {
			return 0, &Float16Error{Op: "bfloat16_mul", Msg: "NaN operand in exact mode", Code: ErrNaN}
		}
		return BFloat16QuietNaN, nil
	}

	aZero := a.IsZero()
	bZero := b.IsZero()

	// 0 * Inf = NaN
	if (aZero && b.IsInf(0)) || (a.IsInf(0) && bZero) {
		if mode == ModeExactArithmetic {
			return 0, &Float16Error{Op: "bfloat16_mul", Msg: "zero times infinity is undefined", Code: ErrInvalidOperation}
		}
		return BFloat16QuietNaN, nil
	}

	// Handle zeros
	if aZero || bZero {
		if a.Signbit() != b.Signbit() {
			return BFloat16NegativeZero, nil
		}
		return BFloat16PositiveZero, nil
	}

	// Handle infinities
	if a.IsInf(0) || b.IsInf(0) {
		if a.Signbit() != b.Signbit() {
			return BFloat16NegativeInfinity, nil
		}
		return BFloat16PositiveInfinity, nil
	}

	if mode == ModeFastArithmetic {
		return BFloat16FromFloat32(a.ToFloat32() * b.ToFloat32()), nil
	}

	// IEEE mode with gradual underflow
	result := a.ToFloat32() * b.ToFloat32()
	bf := BFloat16FromFloat32WithRounding(result, rounding)

	if result != 0 && bf.IsZero() {
		if result > 0 {
			return BFloat16SmallestPosSubnormal, nil
		}
		return BFloat16SmallestNegSubnormal, nil
	}

	return bf, nil
}

// BFloat16DivWithMode performs division with specified arithmetic and rounding modes.
func BFloat16DivWithMode(a, b BFloat16, mode ArithmeticMode, rounding RoundingMode) (BFloat16, error) {
	// NaN propagation
	if a.IsNaN() || b.IsNaN() {
		if mode == ModeExactArithmetic {
			return 0, &Float16Error{Op: "bfloat16_div", Msg: "NaN operand in exact mode", Code: ErrNaN}
		}
		return BFloat16QuietNaN, nil
	}

	// 0 / 0 = NaN
	if a.IsZero() && b.IsZero() {
		if mode == ModeExactArithmetic {
			return 0, &Float16Error{Op: "bfloat16_div", Msg: "zero divided by zero is undefined", Code: ErrInvalidOperation}
		}
		return BFloat16QuietNaN, nil
	}

	// finite / 0 = +/-Inf
	if b.IsZero() {
		if mode == ModeExactArithmetic {
			return 0, &Float16Error{Op: "bfloat16_div", Msg: "division by zero", Code: ErrDivisionByZero}
		}
		if a.Signbit() != b.Signbit() {
			return BFloat16NegativeInfinity, nil
		}
		return BFloat16PositiveInfinity, nil
	}

	// 0 / finite = +/-0
	if a.IsZero() {
		if a.Signbit() != b.Signbit() {
			return BFloat16NegativeZero, nil
		}
		return BFloat16PositiveZero, nil
	}

	// Inf / Inf = NaN
	if a.IsInf(0) && b.IsInf(0) {
		if mode == ModeExactArithmetic {
			return 0, &Float16Error{Op: "bfloat16_div", Msg: "infinity divided by infinity is undefined", Code: ErrInvalidOperation}
		}
		return BFloat16QuietNaN, nil
	}

	// Inf / finite = +/-Inf
	if a.IsInf(0) {
		if a.Signbit() != b.Signbit() {
			return BFloat16NegativeInfinity, nil
		}
		return BFloat16PositiveInfinity, nil
	}

	// finite / Inf = +/-0
	if b.IsInf(0) {
		if a.Signbit() != b.Signbit() {
			return BFloat16NegativeZero, nil
		}
		return BFloat16PositiveZero, nil
	}

	if mode == ModeFastArithmetic {
		return BFloat16FromFloat32(a.ToFloat32() / b.ToFloat32()), nil
	}

	// IEEE mode with gradual underflow
	result := a.ToFloat32() / b.ToFloat32()
	bf := BFloat16FromFloat32WithRounding(result, rounding)

	if result != 0 && bf.IsZero() {
		if result > 0 {
			return BFloat16SmallestPosSubnormal, nil
		}
		return BFloat16SmallestNegSubnormal, nil
	}

	return bf, nil
}

// BFloat16FMA computes a fused multiply-add (a*b + c) for BFloat16 values.
// This is a stub that returns an error; a full implementation is planned for a future phase.
func BFloat16FMA(a, b, c BFloat16) (BFloat16, error) {
	// NaN propagation
	if a.IsNaN() || b.IsNaN() || c.IsNaN() {
		return BFloat16QuietNaN, nil
	}

	// Use float64 FMA for intermediate precision, then round back to BFloat16
	fa := float64(a.ToFloat32())
	fb := float64(b.ToFloat32())
	fc := float64(c.ToFloat32())
	result := math.FMA(fa, fb, fc)

	return BFloat16FromFloat32(float32(result)), nil
}
