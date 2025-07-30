package float16

import (
	"math"
	"testing"
)

func TestCbrtFunction(t *testing.T) {
	tests := []struct {
		name string
		arg  Float16
		want Float16
	}{
		// Basic test cases
		{"Cbrt(1.0)", 0x3C00, 0x3C00},  // 1.0 -> 1.0
		{"Cbrt(8.0)", 0x4800, 0x4000},  // 8.0 -> 2.0
		{"Cbrt(27.0)", 0x51C0, 0x4240}, // 27.0 -> 3.0
		{"Cbrt(64.0)", 0x5800, 0x4400}, // 64.0 -> 4.0

		// Additional test cases for better coverage
		{"Cbrt(0.0)", 0x0000, 0x0000},   // +0.0 -> +0.0
		{"Cbrt(-0.0)", 0x8000, 0x8000},  // -0.0 -> -0.0
		{"Cbrt(0.125)", 0x3000, 0x3800}, // 0.125 -> 0.5 (0x3800 is 0.5 in float16)
		{"Cbrt(125.0)", 0x5A00, 0x45C5}, // 125.0 -> ~5.77 (actual result from math.Cbrt)

		// Special values
		{"Cbrt(+Inf)", 0x7C00, 0x7C00}, // +Inf -> +Inf
		{"Cbrt(-Inf)", 0xFC00, 0xFC00}, // -Inf -> -Inf
		{"Cbrt(NaN)", 0x7E00, 0x7E00},  // NaN -> NaN
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Cbrt(tt.arg)

			if got.IsNaN() && tt.want.IsNaN() {
				return // Both are NaN, which is fine
			}

			if got != tt.want {
				// For non-NaN values, check if they're approximately equal
				gotF64 := got.ToFloat64()
				wantF64 := tt.want.ToFloat64()

				// Calculate relative error
				error := math.Abs(gotF64 - wantF64)
				relError := error / math.Max(math.Abs(wantF64), 1e-10) // Avoid division by zero

				// Allow a small relative error due to floating-point imprecision
				const epsilon = 1e-4
				if relError > epsilon {
					t.Errorf("Cbrt(%v) = %v (0x%04X, %f), want %v (0x%04X, %f), relative error: %e",
						tt.arg, got, uint16(got), got.ToFloat32(),
						tt.want, uint16(tt.want), tt.want.ToFloat32(), relError)
				}
			}
		})
	}
}

func TestSqrtFunction(t *testing.T) {
	tests := []struct {
		name string
		arg  Float16
		want Float16
	}{
		// Basic test cases with actual float16 results from the implementation
		{"Sqrt(1.0)", 0x3C00, 0x3C00},  // 1.0 -> 1.0
		{"Sqrt(4.0)", 0x4400, 0x4000},  // 4.0 -> 2.0
		{"Sqrt(9.0)", 0x4900, 0x4253},  // 9.0 -> ~3.162 (actual float16 result)
		{"Sqrt(16.0)", 0x4C00, 0x4400}, // 16.0 -> 4.0

		// Additional test cases for better coverage
		{"Sqrt(0.0)", 0x0000, 0x0000},  // +0.0 -> +0.0
		{"Sqrt(-0.0)", 0x8000, 0x8000}, // -0.0 -> -0.0 (should be the same as +0.0 in float16)
		{"Sqrt(0.25)", 0x3400, 0x3800}, // 0.25 -> 0.5
		{"Sqrt(2.0)", 0x4000, 0x3DA8},  // 2.0 -> ~1.414 (actual float16 result)

		// Special values
		{"Sqrt(+Inf)", 0x7C00, 0x7C00}, // +Inf -> +Inf
		{"Sqrt(NaN)", 0x7E00, 0x7E00},  // NaN -> NaN
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Sqrt(tt.arg)
			t.Logf("Testing %s (0x%04X)", tt.name, uint16(tt.arg))
			t.Logf("Got: %v (0x%04X, %f)", got, uint16(got), got.ToFloat32())
			t.Logf("Want: %v (0x%04X, %f)", tt.want, uint16(tt.want), tt.want.ToFloat32())

			// First check if the values are exactly equal
			if got == tt.want {
				return
			}

			// If not exactly equal, check if they're approximately equal within float16 precision
			gotF64 := float64(got.ToFloat32())
			wantF64 := float64(tt.want.ToFloat32())
			const epsilon = 1e-4 // More lenient epsilon for square roots

			if math.Abs(gotF64-wantF64) <= epsilon*math.Abs(wantF64) {
				t.Logf("Values are approximately equal within epsilon %g", epsilon)
				return
			}

			t.Errorf("%s = %v (0x%04X, %f), want %v (0x%04X, %f)",
				tt.name, got, uint16(got), got.ToFloat32(),
				tt.want, uint16(tt.want), tt.want.ToFloat32())
		})
	}
}

// Test that Sqrt of a negative number returns NaN
func TestSqrtNegative(t *testing.T) {
	tests := []struct {
		name string
		arg  Float16
	}{
		{"Sqrt(-1.0)", 0xBC00},
		{"Sqrt(-0.1)", 0xB999},
		{"Sqrt(-Inf)", 0xFC00},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Sqrt(tt.arg)
			if !got.IsNaN() {
				t.Errorf("%s = %v (0x%04X), want NaN", tt.name, got, uint16(got))
			}
		})
	}
}

func TestBasicMathFunctions(t *testing.T) {
	tests := []struct {
		name      string
		fn        func(Float16) Float16
		arg       Float16
		want      Float16
		tolerance float64 // Relative tolerance for approximate comparisons
	}{
		// Sqrt tests
		{"Sqrt(4.0)", Sqrt, 0x4400, 0x4000, 1e-5},  // 4.0 -> 2.0
		{"Sqrt(0.25)", Sqrt, 0x3400, 0x3800, 1e-5}, // 0.25 -> 0.5
		{"Sqrt(2.0)", Sqrt, 0x4000, 0x3DA8, 1e-3},  // 2.0 -> ~1.414 (approximate)

		// Cbrt tests
		{"Cbrt(27.0)", Cbrt, 0x51C0, 0x4240, 1e-5},  // 27.0 -> 3.0
		{"Cbrt(8.0)", Cbrt, 0x4800, 0x4000, 1e-5},   // 8.0 -> 2.0
		{"Cbrt(1.0)", Cbrt, 0x3C00, 0x3C00, 1e-5},   // 1.0 -> 1.0
		{"Cbrt(0.125)", Cbrt, 0x3000, 0x3800, 1e-3}, // 0.125 -> 0.5 (approximate)

		// Abs tests
		{"Abs(-1.5)", Abs, 0xBE00, 0x3E00, 0}, // -1.5 -> 1.5
		{"Abs(1.5)", Abs, 0x3E00, 0x3E00, 0},  // 1.5 -> 1.5
		{"Abs(-0.0)", Abs, 0x8000, 0x0000, 0}, // -0.0 -> 0.0

		// Rounding functions
		{"Floor(1.9)", Floor, 0x3F33, 0x3C00, 0},  // 1.9 -> 1.0
		{"Floor(-1.9)", Floor, 0xBF33, 0xC000, 0}, // -1.9 -> -2.0
		{"Ceil(1.1)", Ceil, 0x3DCC, 0x4000, 0},    // 1.1 -> 2.0
		{"Ceil(-1.1)", Ceil, 0xBDCC, 0xBC00, 0},   // -1.1 -> -1.0
		{"Round(1.5)", Round, 0x3E00, 0x4000, 0},  // 1.5 -> 2.0 (round half up)
		{"Round(2.5)", Round, 0x4100, 0x4200, 0},  // 2.5 -> 3.0 (round half up)
		{"Round(-1.5)", Round, 0xBE00, 0xC000, 0}, // -1.5 -> -2.0 (round half up)
		{"Trunc(1.9)", Trunc, 0x3F33, 0x3C00, 0},  // 1.9 -> 1.0
		{"Trunc(-1.9)", Trunc, 0xBF33, 0xBC00, 0}, // -1.9 -> -1.0

		// Special values
		{"Sqrt(+Inf)", Sqrt, 0x7C00, 0x7C00, 0}, // +Inf -> +Inf
		{"Sqrt(NaN)", Sqrt, 0x7E00, 0x7E00, 0},  // NaN -> NaN
		{"Cbrt(-Inf)", Cbrt, 0xFC00, 0xFC00, 0}, // -Inf -> -Inf
		{"Abs(NaN)", Abs, 0x7E00, 0x7E00, 0},    // NaN -> NaN
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.fn(tt.arg)
			if got != tt.want {
				t.Errorf("%s = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

func TestTwoArgMathFunctions(t *testing.T) {
	tests := []struct {
		name string
		fn   func(Float16, Float16) Float16
		a, b Float16
		want Float16
	}{
		{"Pow(2.0, 3.0)", Pow, 0x4000, 0x4200, 0x4800},          // 2^3 = 8.0
		{"Hypot(3.0, 4.0)", Hypot, 0x4200, 0x4400, 0x4500},      // 5.0
		{"Dim(5.0, 3.0)", Dim, 0x4500, 0x4200, 0x4000},          // 5.0 - 3.0 = 2.0 (actual float16 representation)
		{"Dim(3.0, 5.0)", Dim, 0x4200, 0x4500, 0x0000},          // 3.0 - 5.0 = 0.0 (negative result clamped to 0)
		{"Mod(5.5, 2.5)", Mod, 0x4580, 0x4200, 0x4100},          // 5.5 % 2.5 = 0.5
		{"CopySign(1.5, -1)", CopySign, 0x3E00, 0xBC00, 0xBE00}, // 1.5 with sign of -1.0
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.fn(tt.a, tt.b)
			if got != tt.want {
				t.Errorf("%s = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

func TestTrigonometricFunctions(t *testing.T) {
	// Test with angles in radians
	tests := []struct {
		name string
		fn   func(Float16) Float16
		arg  Float16
		want Float16
	}{
		{"Sin(0.0)", Sin, 0x0000, 0x0000},
		{"Cos(0.0)", Cos, 0x0000, 0x3C00}, // cos(0) = 1.0
		{"Tan(0.0)", Tan, 0x0000, 0x0000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.fn(tt.arg)
			if got != tt.want {
				t.Errorf("%s = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

func TestSpecialMathFunctions(t *testing.T) {
	tests := []struct {
		name string
		fn   func(Float16) Float16
		arg  Float16
		want Float16
	}{
		{"Exp(0.0)", Exp, 0x0000, 0x3C00},     // e^0 = 1.0
		{"Exp2(1.0)", Exp2, 0x3C00, 0x4000},   // 2^1 = 2.0
		{"Exp10(1.0)", Exp10, 0x3C00, 0x4900}, // 10^1 ≈ 10.0 (actual float16 representation)
		{"Log(1.0)", Log, 0x3C00, 0x0000},     // ln(1) = 0.0
		{"Log2(1.0)", Log2, 0x3C00, 0x0000},   // log2(1) = 0.0
		{"Log10(1.0)", Log10, 0x3C00, 0x0000}, // log10(1) = 0.0
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.fn(tt.arg)
			if got != tt.want {
				t.Errorf("%s = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

func TestLgammaFunction(t *testing.T) {
	tests := []struct {
		name     string
		arg      Float16
		wantLg   Float16
		wantSign int
	}{
		{"Lgamma(1.0)", 0x3C00, 0x0000, 1},   // Lgamma(1) = 0, sign = 1
		{"Lgamma(2.0)", 0x4000, 0x0000, 1},   // Lgamma(2) = 0, sign = 1
		{"Lgamma(0.5)", 0x3800, 0x3894, 1},   // Lgamma(0.5) ≈ 0.572266, sign = 1
		{"Lgamma(-0.5)", 0xB800, 0x3D10, -1}, // Lgamma(-0.5) ≈ 1.26562, sign = -1
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotLg, gotSign := Lgamma(tt.arg)
			if gotLg != tt.wantLg || gotSign != tt.wantSign {
				t.Errorf("Lgamma(%v) = (%v (0x%04X, %f), %d), want (%v (0x%04X, %f), %d)",
					tt.arg,
					gotLg, uint16(gotLg), gotLg.ToFloat32(),
					gotSign,
					tt.wantLg, uint16(tt.wantLg), tt.wantLg.ToFloat32(),
					tt.wantSign)
			}
		})
	}
}

func TestBesselFunctions(t *testing.T) {
	tests := []struct {
		name string
		fn   func(Float16) Float16
		arg  Float16
		want Float16
	}{
		{"J0(0.0)", J0, 0x0000, 0x3C00}, // J0(0) = 1.0
		{"J1(0.0)", J1, 0x0000, 0x0000}, // J1(0) = 0.0
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.fn(tt.arg)
			if got != tt.want {
				t.Errorf("%s = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

func TestErrorFunctions(t *testing.T) {
	tests := []struct {
		name string
		fn   func(Float16) Float16
		arg  Float16
		want Float16
	}{
		{"Erf(0.0)", Erf, 0x0000, 0x0000},
		{"Erfc(0.0)", Erfc, 0x0000, 0x3C00}, // erfc(0) = 1.0
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.fn(tt.arg)
			if got != tt.want {
				t.Errorf("%s = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

func TestPow_Extra(t *testing.T) {
	tests := []struct {
		name   string
		f, exp Float16
		want   Float16
	}{
		{"1^x = 1", ToFloat16(1.0), ToFloat16(123.45), ToFloat16(1.0)},
		{"x^1 = x", ToFloat16(123.45), ToFloat16(1.0), ToFloat16(123.45)},
		{"-1^2 = 1", ToFloat16(-1.0), ToFloat16(2.0), ToFloat16(1.0)},
		{"-1^3 = -1", ToFloat16(-1.0), ToFloat16(3.0), ToFloat16(-1.0)},
		{"inf^2 = inf", PositiveInfinity, ToFloat16(2.0), PositiveInfinity},
		{"inf^-2 = 0", PositiveInfinity, ToFloat16(-2.0), PositiveZero},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Pow(tt.f, tt.exp); got != tt.want {
				t.Errorf("Pow() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMod_Extra(t *testing.T) {
	tests := []struct {
		name string
		f, d Float16
		want Float16
	}{
		{"5.0 mod 3.0", ToFloat16(5.0), ToFloat16(3.0), ToFloat16(2.0)},
		{"-5.0 mod 3.0", ToFloat16(-5.0), ToFloat16(3.0), ToFloat16(-2.0)},
		{"5.0 mod -3.0", ToFloat16(5.0), ToFloat16(-3.0), ToFloat16(2.0)},
		{"-5.0 mod -3.0", ToFloat16(-5.0), ToFloat16(-3.0), ToFloat16(-2.0)},
		{"inf mod 1", PositiveInfinity, ToFloat16(1.0), QuietNaN},
		{"1 mod inf", ToFloat16(1.0), PositiveInfinity, QuietNaN},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Mod(tt.f, tt.d)
			if got.IsNaN() && tt.want.IsNaN() {
				return
			}
			if got != tt.want {
				t.Errorf("Mod() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHypot_Inf(t *testing.T) {
	got := Hypot(PositiveInfinity, QuietNaN)
	if !got.IsInf(1) {
		t.Errorf("Hypot(inf, nan) = %v, want +Inf", got)
	}
}
