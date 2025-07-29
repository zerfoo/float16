package float16

import (
	"math"
	"testing"
)

func TestMathExtra(t *testing.T) {
	testCases := []struct {
		name     string
		f        Float16
		f2       Float16
		f3       Float16
		expected Float16
		op       interface{}
	}{
		// Asin
		{"Asin(0)", FromFloat32(0), 0, 0, FromFloat32(0), Asin},
		{"Asin(0.5)", FromFloat32(0.5), 0, 0, ToFloat16(float32(math.Asin(0.5))), Asin},
		{"Asin(-0.5)", FromFloat32(-0.5), 0, 0, ToFloat16(float32(math.Asin(-0.5))), Asin},
		{"Asin(1)", FromFloat32(1), 0, 0, ToFloat16(float32(math.Asin(1))), Asin},
		{"Asin(-1)", FromFloat32(-1), 0, 0, ToFloat16(float32(math.Asin(-1))), Asin},
		{"Asin(2)", FromFloat32(2), 0, 0, QuietNaN, Asin},
		{"Asin(-2)", FromFloat32(-2), 0, 0, QuietNaN, Asin},
		{"Asin(NaN)", QuietNaN, 0, 0, QuietNaN, Asin},

		// Acos
		{"Acos(0)", FromFloat32(0), 0, 0, ToFloat16(float32(math.Acos(0))), Acos},
		{"Acos(0.5)", FromFloat32(0.5), 0, 0, ToFloat16(float32(math.Acos(0.5))), Acos},
		{"Acos(-0.5)", FromFloat32(-0.5), 0, 0, ToFloat16(float32(math.Acos(-0.5))), Acos},
		{"Acos(1)", FromFloat32(1), 0, 0, ToFloat16(float32(math.Acos(1))), Acos},
		{"Acos(-1)", FromFloat32(-1), 0, 0, ToFloat16(float32(math.Acos(-1))), Acos},
		{"Acos(2)", FromFloat32(2), 0, 0, QuietNaN, Acos},
		{"Acos(-2)", FromFloat32(-2), 0, 0, QuietNaN, Acos},
		{"Acos(NaN)", QuietNaN, 0, 0, QuietNaN, Acos},

		// Atan
		{"Atan(0)", FromFloat32(0), 0, 0, FromFloat32(0), Atan},
		{"Atan(1)", FromFloat32(1), 0, 0, ToFloat16(float32(math.Atan(1))), Atan},
		{"Atan(-1)", FromFloat32(-1), 0, 0, ToFloat16(float32(math.Atan(-1))), Atan},
		{"Atan(Inf)", PositiveInfinity, 0, 0, Div(Pi, FromInt(2)), Atan},
		{"Atan(-Inf)", NegativeInfinity, 0, 0, Div(Pi, FromInt(2)).Neg(), Atan},
		{"Atan(NaN)", QuietNaN, 0, 0, QuietNaN, Atan},

		// Atan2
		{"Atan2(0, 0)", FromFloat32(0), FromFloat32(0), 0, FromFloat32(0), Atan2},
		{"Atan2(1, 1)", FromFloat32(1), FromFloat32(1), 0, ToFloat16(float32(math.Atan2(1, 1))), Atan2},
		{"Atan2(-1, 1)", FromFloat32(-1), FromFloat32(1), 0, ToFloat16(float32(math.Atan2(-1, 1))), Atan2},
		{"Atan2(NaN, 1)", QuietNaN, FromFloat32(1), 0, QuietNaN, Atan2},
		{"Atan2(1, NaN)", FromFloat32(1), QuietNaN, 0, QuietNaN, Atan2},

		// Sinh
		{"Sinh(0)", FromFloat32(0), 0, 0, FromFloat32(0), Sinh},
		{"Sinh(1)", FromFloat32(1), 0, 0, ToFloat16(float32(math.Sinh(1))), Sinh},
		{"Sinh(-1)", FromFloat32(-1), 0, 0, ToFloat16(float32(math.Sinh(-1))), Sinh},
		{"Sinh(Inf)", PositiveInfinity, 0, 0, PositiveInfinity, Sinh},
		{"Sinh(-Inf)", NegativeInfinity, 0, 0, NegativeInfinity, Sinh},
		{"Sinh(NaN)", QuietNaN, 0, 0, QuietNaN, Sinh},

		// Cosh
		{"Cosh(0)", FromFloat32(0), 0, 0, FromInt(1), Cosh},
		{"Cosh(1)", FromFloat32(1), 0, 0, ToFloat16(float32(math.Cosh(1))), Cosh},
		{"Cosh(-1)", FromFloat32(-1), 0, 0, ToFloat16(float32(math.Cosh(-1))), Cosh},
		{"Cosh(Inf)", PositiveInfinity, 0, 0, PositiveInfinity, Cosh},
		{"Cosh(-Inf)", NegativeInfinity, 0, 0, PositiveInfinity, Cosh},
		{"Cosh(NaN)", QuietNaN, 0, 0, QuietNaN, Cosh},

		// Tanh
		{"Tanh(0)", FromFloat32(0), 0, 0, FromFloat32(0), Tanh},
		{"Tanh(1)", FromFloat32(1), 0, 0, ToFloat16(float32(math.Tanh(1))), Tanh},
		{"Tanh(-1)", FromFloat32(-1), 0, 0, ToFloat16(float32(math.Tanh(-1))), Tanh},
		{"Tanh(Inf)", PositiveInfinity, 0, 0, FromInt(1), Tanh},
		{"Tanh(-Inf)", NegativeInfinity, 0, 0, FromInt(-1), Tanh},
		{"Tanh(NaN)", QuietNaN, 0, 0, QuietNaN, Tanh},

		// RoundToEven
		{"RoundToEven(0)", FromFloat32(0), 0, 0, FromFloat32(0), RoundToEven},
		{"RoundToEven(0.5)", FromFloat32(0.5), 0, 0, FromFloat32(0), RoundToEven},
		{"RoundToEven(1.5)", FromFloat32(1.5), 0, 0, FromFloat32(2), RoundToEven},
		{"RoundToEven(2.5)", FromFloat32(2.5), 0, 0, FromFloat32(2), RoundToEven},
		{"RoundToEven(NaN)", QuietNaN, 0, 0, QuietNaN, RoundToEven},
		{"RoundToEven(Inf)", PositiveInfinity, 0, 0, PositiveInfinity, RoundToEven},

		// Remainder
		{"Remainder(5, 3)", FromFloat32(5), FromFloat32(3), 0, FromFloat32(-1), Remainder},
		{"Remainder(-5, 3)", FromFloat32(-5), FromFloat32(3), 0, FromFloat32(1), Remainder},
		{"Remainder(5, -3)", FromFloat32(5), FromFloat32(-3), 0, FromFloat32(-1), Remainder},
		{"Remainder(0, 1)", FromFloat32(0), FromFloat32(1), 0, FromFloat32(0), Remainder},
		{"Remainder(1, 0)", FromFloat32(1), FromFloat32(0), 0, QuietNaN, Remainder},
		{"Remainder(NaN, 1)", QuietNaN, FromFloat32(1), 0, QuietNaN, Remainder},
		{"Remainder(1, NaN)", FromFloat32(1), QuietNaN, 0, QuietNaN, Remainder},
		{"Remainder(Inf, 1)", PositiveInfinity, FromFloat32(1), 0, QuietNaN, Remainder},
		{"Remainder(1, Inf)", FromFloat32(1), PositiveInfinity, 0, FromFloat32(1), Remainder},

		// Clamp
		{"Clamp(0.5, 0, 1)", FromFloat32(0.5), FromFloat32(0), FromFloat32(1), FromFloat32(0.5), Clamp},
		{"Clamp(-0.5, 0, 1)", FromFloat32(-0.5), FromFloat32(0), FromFloat32(1), FromFloat32(0), Clamp},
		{"Clamp(1.5, 0, 1)", FromFloat32(1.5), FromFloat32(0), FromFloat32(1), FromFloat32(1), Clamp},
		{"Clamp(NaN, 0, 1)", QuietNaN, FromFloat32(0), FromFloat32(1), QuietNaN, Clamp},

		// Lerp
		{"Lerp(0, 1, 0.5)", FromFloat32(0), FromFloat32(1), FromFloat32(0.5), FromFloat32(0.5), Lerp},
		{"Lerp(0, 1, 0)", FromFloat32(0), FromFloat32(1), FromFloat32(0), FromFloat32(0), Lerp},
		{"Lerp(0, 1, 1)", FromFloat32(0), FromFloat32(1), FromFloat32(1), FromFloat32(1), Lerp},

		// Sign
		{"Sign(5)", FromFloat32(5), 0, 0, FromInt(1), Sign},
		{"Sign(-5)", FromFloat32(-5), 0, 0, FromInt(-1), Sign},
		{"Sign(0)", FromFloat32(0), 0, 0, FromFloat32(0), Sign},
		{"Sign(NaN)", QuietNaN, 0, 0, QuietNaN, Sign},

		// Gamma
		{"Gamma(1)", FromFloat32(1), 0, 0, FromInt(1), Gamma},
		{"Gamma(0.5)", FromFloat32(0.5), 0, 0, ToFloat16(float32(math.Gamma(0.5))), Gamma},
		{"Gamma(NaN)", QuietNaN, 0, 0, QuietNaN, Gamma},
		{"Gamma(-Inf)", NegativeInfinity, 0, 0, QuietNaN, Gamma},
		{"Gamma(+Inf)", PositiveInfinity, 0, 0, PositiveInfinity, Gamma},

		// J0
		{"J0(0)", FromFloat32(0), 0, 0, FromInt(1), J0},
		{"J0(1)", FromFloat32(1), 0, 0, ToFloat16(float32(math.J0(1))), J0},
		{"J0(NaN)", QuietNaN, 0, 0, QuietNaN, J0},
		{"J0(Inf)", PositiveInfinity, 0, 0, FromFloat32(0), J0},

		// J1
		{"J1(0)", FromFloat32(0), 0, 0, FromFloat32(0), J1},
		{"J1(1)", FromFloat32(1), 0, 0, ToFloat16(float32(math.J1(1))), J1},
		{"J1(NaN)", QuietNaN, 0, 0, QuietNaN, J1},
		{"J1(Inf)", PositiveInfinity, 0, 0, FromFloat32(0), J1},

		// Erf
		{"Erf(0)", FromFloat32(0), 0, 0, FromFloat32(0), Erf},
		{"Erf(1)", FromFloat32(1), 0, 0, ToFloat16(float32(math.Erf(1))), Erf},
		{"Erf(-1)", FromFloat32(-1), 0, 0, ToFloat16(float32(math.Erf(-1))), Erf},
		{"Erf(NaN)", QuietNaN, 0, 0, QuietNaN, Erf},
		{"Erf(Inf)", PositiveInfinity, 0, 0, FromInt(1), Erf},
		{"Erf(-Inf)", NegativeInfinity, 0, 0, FromInt(-1), Erf},

		// Erfc
		{"Erfc(0)", FromFloat32(0), 0, 0, FromInt(1), Erfc},
		{"Erfc(1)", FromFloat32(1), 0, 0, ToFloat16(float32(math.Erfc(1))), Erfc},
		{"Erfc(NaN)", QuietNaN, 0, 0, QuietNaN, Erfc},
		{"Erfc(Inf)", PositiveInfinity, 0, 0, FromFloat32(0), Erfc},
		{"Erfc(-Inf)", NegativeInfinity, 0, 0, FromInt(2), Erfc},

		// Pow
		{"Pow(2, 3)", FromFloat32(2), FromFloat32(3), 0, FromFloat32(8), Pow},
		{"Pow(0, -1)", FromFloat32(0), FromFloat32(-1), 0, PositiveInfinity, Pow},
		{"Pow(Inf, -1)", PositiveInfinity, FromFloat32(-1), 0, FromFloat32(0), Pow},
		{"Pow(NaN, 1)", QuietNaN, FromFloat32(1), 0, QuietNaN, Pow},

		// Exp
		{"Exp(0)", FromFloat32(0), 0, 0, FromFloat32(1), Exp},
		{"Exp(1)", FromFloat32(1), 0, 0, ToFloat16(float32(math.Exp(1))), Exp},
		{"Exp(Inf)", PositiveInfinity, 0, 0, PositiveInfinity, Exp},
		{"Exp(-Inf)", NegativeInfinity, 0, 0, FromFloat32(0), Exp},
		{"Exp(NaN)", QuietNaN, 0, 0, QuietNaN, Exp},

		// Exp2
		{"Exp2(0)", FromFloat32(0), 0, 0, FromFloat32(1), Exp2},
		{"Exp2(1)", FromFloat32(1), 0, 0, FromFloat32(2), Exp2},
		{"Exp2(Inf)", PositiveInfinity, 0, 0, PositiveInfinity, Exp2},
		{"Exp2(-Inf)", NegativeInfinity, 0, 0, FromFloat32(0), Exp2},
		{"Exp2(NaN)", QuietNaN, 0, 0, QuietNaN, Exp2},

		// Log
		{"Log(1)", FromFloat32(1), 0, 0, FromFloat32(0), Log},
		{"Log(0)", FromFloat32(0), 0, 0, NegativeInfinity, Log},
		{"Log(-1)", FromFloat32(-1), 0, 0, QuietNaN, Log},
		{"Log(Inf)", PositiveInfinity, 0, 0, PositiveInfinity, Log},
		{"Log(NaN)", QuietNaN, 0, 0, QuietNaN, Log},

		// Log2
		{"Log2(1)", FromFloat32(1), 0, 0, FromFloat32(0), Log2},
		{"Log2(0)", FromFloat32(0), 0, 0, NegativeInfinity, Log2},
		{"Log2(-1)", FromFloat32(-1), 0, 0, QuietNaN, Log2},
		{"Log2(Inf)", PositiveInfinity, 0, 0, PositiveInfinity, Log2},
		{"Log2(NaN)", QuietNaN, 0, 0, QuietNaN, Log2},

		// Log10
		{"Log10(1)", FromFloat32(1), 0, 0, FromFloat32(0), Log10},
		{"Log10(0)", FromFloat32(0), 0, 0, NegativeInfinity, Log10},
		{"Log10(-1)", FromFloat32(-1), 0, 0, QuietNaN, Log10},
		{"Log10(Inf)", PositiveInfinity, 0, 0, PositiveInfinity, Log10},
		{"Log10(NaN)", QuietNaN, 0, 0, QuietNaN, Log10},

		// Sin
		{"Sin(0)", FromFloat32(0), 0, 0, FromFloat32(0), Sin},
				{"Sin(Pi/2)", Div(Pi, FromInt(2)), 0, 0, FromFloat32(1), Sin},
		{"Sin(NaN)", QuietNaN, 0, 0, QuietNaN, Sin},

		// Cos
		{"Cos(0)", FromFloat32(0), 0, 0, FromFloat32(1), Cos},
		{"Cos(Pi)", Pi, 0, 0, FromFloat32(-1), Cos},
		{"Cos(NaN)", QuietNaN, 0, 0, QuietNaN, Cos},

		// Tan
		{"Tan(0)", FromFloat32(0), 0, 0, FromFloat32(0), Tan},
		{"Tan(Pi/4)", Pi, 0, 0, FromFloat32(0), Tan},
		{"Tan(NaN)", QuietNaN, 0, 0, QuietNaN, Tan},

		// Floor
		{"Floor(1.5)", FromFloat32(1.5), 0, 0, FromFloat32(1), Floor},
		{"Floor(-1.5)", FromFloat32(-1.5), 0, 0, FromFloat32(-2), Floor},
		{"Floor(NaN)", QuietNaN, 0, 0, QuietNaN, Floor},

		// Ceil
		{"Ceil(1.5)", FromFloat32(1.5), 0, 0, FromFloat32(2), Ceil},
		{"Ceil(-1.5)", FromFloat32(-1.5), 0, 0, FromFloat32(-1), Ceil},
		{"Ceil(NaN)", QuietNaN, 0, 0, QuietNaN, Ceil},

		// Round
		{"Round(1.5)", FromFloat32(1.5), 0, 0, FromFloat32(2), Round},
		{"Round(1.4)", FromFloat32(1.4), 0, 0, FromFloat32(1), Round},
		{"Round(NaN)", QuietNaN, 0, 0, QuietNaN, Round},

		// Trunc
		{"Trunc(1.5)", FromFloat32(1.5), 0, 0, FromFloat32(1), Trunc},
		{"Trunc(-1.5)", FromFloat32(-1.5), 0, 0, FromFloat32(-1), Trunc},
		{"Trunc(NaN)", QuietNaN, 0, 0, QuietNaN, Trunc},

		// Mod
		{"Mod(5, 3)", FromFloat32(5), FromFloat32(3), 0, FromFloat32(2), Mod},
		{"Mod(5, 0)", FromFloat32(5), FromFloat32(0), 0, QuietNaN, Mod},
		{"Mod(NaN, 1)", QuietNaN, FromFloat32(1), 0, QuietNaN, Mod},

		// Hypot
		{"Hypot(3, 4)", FromFloat32(3), FromFloat32(4), 0, FromFloat32(5), Hypot},
		{"Hypot(Inf, 4)", PositiveInfinity, FromFloat32(4), 0, PositiveInfinity, Hypot},
		{"Hypot(NaN, 4)", QuietNaN, FromFloat32(4), 0, QuietNaN, Hypot},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var res Float16
			switch op := tc.op.(type) {
			case func(Float16) Float16:
				res = op(tc.f)
			case func(Float16, Float16) Float16:
				res = op(tc.f, tc.f2)
			case func(Float16, Float16, Float16) Float16:
				res = op(tc.f, tc.f2, tc.f3)
			}
			if tc.expected.IsNaN() {
				if !res.IsNaN() {
					t.Errorf("Expected NaN, got %v", res)
				}
			} else if Abs(Sub(res, tc.expected)).ToFloat32() > 1e-3 {
				t.Errorf("Expected %v, got %v", tc.expected, res)
			}
		})
	}
}


func TestLgamma(t *testing.T) {
	testCases := []struct {
		name         string
		f            Float16
		expectedLgam Float16
		expectedSign int
	}{
		{"Lgamma(1)", FromFloat32(1), FromFloat32(0), 1},
		{"Lgamma(0.5)", FromFloat32(0.5), func() Float16 {
			lg, _ := math.Lgamma(0.5)
			return ToFloat16(float32(lg))
		}(), 1},
		{"Lgamma(NaN)", QuietNaN, QuietNaN, 1},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			lgam, sign := Lgamma(tc.f)
			if lgam.Bits() != tc.expectedLgam.Bits() || sign != tc.expectedSign {
				t.Errorf("Expected (%v, %v), got (%v, %v)", tc.expectedLgam, tc.expectedSign, lgam, sign)
			}
		})
	}
}

func TestY0Y1(t *testing.T) {
	testCases := []struct {
		name     string
		f        Float16
		expected Float16
		op       func(Float16) Float16
	}{
		// Y0
		{"Y0(1)", FromFloat32(1), ToFloat16(float32(math.Y0(1))), Y0},
		{"Y0(0)", FromFloat32(0), NegativeInfinity, Y0},
		{"Y0(-1)", FromFloat32(-1), QuietNaN, Y0},
		{"Y0(NaN)", QuietNaN, QuietNaN, Y0},
		{"Y0(Inf)", PositiveInfinity, FromFloat32(0), Y0},

		// Y1
		{"Y1(1)", FromFloat32(1), ToFloat16(float32(math.Y1(1))), Y1},
		{"Y1(0)", FromFloat32(0), NegativeInfinity, Y1},
		{"Y1(-1)", FromFloat32(-1), QuietNaN, Y1},
		{"Y1(NaN)", QuietNaN, QuietNaN, Y1},
		{"Y1(Inf)", PositiveInfinity, FromFloat32(0), Y1},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res := tc.op(tc.f)
						if res.Bits() != tc.expected.Bits() {
				t.Errorf("Expected %v, got %v", tc.expected, res)
			}
		})
	}
}
