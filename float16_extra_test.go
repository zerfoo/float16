package float16

import (
	"math"
	"testing"
)

func TestVersionFunctions(t *testing.T) {
	if got := GetVersion(); got != Version {
		t.Errorf("GetVersion() = %s, want %s", got, Version)
	}
}

func TestZeroAndOne(t *testing.T) {
	if got := Zero(); got != PositiveZero {
		t.Errorf("Zero() = %v, want %v", got, PositiveZero)
	}
	if got := One(); got != 0x3C00 {
		t.Errorf("One() = %v, want 0x3C00", got)
	}
}

func TestInf(t *testing.T) {
	if got := Inf(1); got != PositiveInfinity {
		t.Errorf("Inf(1) = %v, want %v", got, PositiveInfinity)
	}
	if got := Inf(-1); got != NegativeInfinity {
		t.Errorf("Inf(-1) = %v, want %v", got, NegativeInfinity)
	}
}

func TestNextAfter(t *testing.T) {
	tests := []struct {
		a, b, want Float16
	}{
		{0x3C00, 0x4000, 0x3C01}, // 1.0 -> 2.0 = 1.0009766
		{0x3C00, 0x0000, 0x3BFF}, // 1.0 -> 0.0 = 0.9995117
		{0x3C00, 0x3C00, 0x3C00}, // Equal values
	}

	for _, tt := range tests {
		if got := NextAfter(tt.a, tt.b); got != tt.want {
			t.Errorf("NextAfter(%v, %v) = %v, want %v", tt.a, tt.b, got, tt.want)
		}
	}
}

func TestFrexpLdexp(t *testing.T) {
	tests := []struct {
		f        Float16
		wantFrac Float16
		wantExp  int
		desc     string
	}{
		{0x3C00, 0x3C00, 0, "1.0 = 1.0 * 2^0"},
		{0x4000, 0x3C00, 1, "2.0 = 1.0 * 2^1"},
		{0x3800, 0x3C00, -1, "0.5 = 1.0 * 2^-1"},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			frac, exp := Frexp(tt.f)

			// Convert to float64 for comparison to handle rounding
			gotFrac := frac.ToFloat64()
			wantFrac := tt.wantFrac.ToFloat64()

			// Allow small floating point differences
			if math.Abs(gotFrac*math.Pow(2, float64(exp))-wantFrac*math.Pow(2, float64(tt.wantExp))) > 1e-5 {
				t.Errorf("Frexp(%v) = (%v, %d), want (%v, %d)",
					tt.f, gotFrac, exp, wantFrac, tt.wantExp)
			}

			// Test roundtrip
			if got := Ldexp(frac, exp); got != tt.f {
				t.Errorf("Ldexp(%v, %d) = %v, want %v",
					frac, exp, got, tt.f)
			}
		})
	}
}

func TestModf(t *testing.T) {
	tests := []struct {
		f    Float16
		desc string
	}{
		{0x3E00, "1.5"},
		{0xBE00, "-1.5"},
		{0x3C00, "1.0"},
		{0x0000, "0.0"},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			intPart, fracPart := Modf(tt.f)

			// The sum should equal the original value (within float16 precision)
			sum := ToFloat16(intPart.ToFloat32() + fracPart.ToFloat32())
			if sum != tt.f && !(math.IsNaN(float64(sum.ToFloat32())) && math.IsNaN(float64(tt.f.ToFloat32()))) {
				t.Errorf("Modf(%v) sum = %v + %v = %v, want %v (diff: %g)",
					tt.f, intPart, fracPart, sum, tt.f,
					math.Abs(float64(sum.ToFloat32()-tt.f.ToFloat32())))
			}

			// The fractional part should have the same sign as the input
			if !tt.f.IsZero() && !fracPart.IsZero() {
				if fracPart.Signbit() != tt.f.Signbit() {
					t.Errorf("Modf(%v) fractional part has wrong sign: got %v",
						tt.f, fracPart)
				}
			}
		})
	}
}

func TestFloatClassification(t *testing.T) {
	if !IsFinite(0x3C00) {
		t.Error("IsFinite(1.0) = false, want true")
	}
	if !IsNormal(0x3C00) {
		t.Error("IsNormal(1.0) = false, want true")
	}
	if !IsSubnormal(0x0001) {
		t.Error("IsSubnormal(smallest subnormal) = false, want true")
	}
}

func TestGetDebugInfo(t *testing.T) {
	info := DebugInfo()
	if info == nil {
		t.Fatal("DebugInfo() returned nil")
	}
}

func TestGetBenchmarkOperations(t *testing.T) {
	if ops := GetBenchmarkOperations(); len(ops) == 0 {
		t.Error("GetBenchmarkOperations() returned empty map")
	}
}

func TestVectorOperations(t *testing.T) {
	a := []Float16{0x3C00, 0x4000} // [1.0, 2.0]
	b := []Float16{0x3C00, 0x3C00} // [1.0, 1.0]

	// Test VectorAdd
	result := VectorAdd(a, b)
	if len(result) != 2 || result[0] != 0x4000 || result[1] != 0x4200 {
		t.Errorf("VectorAdd() = %v, want [0x4000, 0x4200]", result)
	}

	// Test VectorMul
	result = VectorMul(a, b)
	if len(result) != 2 || result[0] != 0x3C00 || result[1] != 0x4000 {
		t.Errorf("VectorMul() = %v, want [0x3C00, 0x4000]", result)
	}
}

func TestFastMath(t *testing.T) {
	a, b := Float16(0x3C00), Float16(0x4000) // 1.0, 2.0

	// Test FastAdd
	if got := FastAdd(a, b); got != 0x4200 { // 3.0
		t.Errorf("FastAdd() = %v, want 0x4200", got)
	}

	// Test FastMul
	if got := FastMul(a, b); got != 0x4000 { // 2.0
		t.Errorf("FastMul() = %v, want 0x4000", got)
	}
}

func TestSliceUtils(t *testing.T) {
	s := []Float16{0x3C00, 0x4000, 0x4200} // [1.0, 2.0, 3.0]

	// Test ComputeSliceStats
	stats := ComputeSliceStats(s)
	if stats.Min != 0x3C00 || stats.Max != 0x4200 || stats.Length != 3 {
		t.Errorf("ComputeSliceStats() = %+v, want {Min:0x3C00, Max:0x4200, Length:3}", stats)
	}
}
