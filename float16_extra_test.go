package float16

import (
	"testing"
)

func TestFloat16Extra(t *testing.T) {
	testCases := []struct {
		name     string
		f        Float16
		f2       Float16
		sign     int
		expected interface{}
		op       string
	}{
		// IsInf
		{"IsInf(+Inf, 1)", PositiveInfinity, 0, 1, true, "IsInf"},
		{"IsInf(-Inf, -1)", NegativeInfinity, 0, -1, true, "IsInf"},
		{"IsInf(+Inf, 0)", PositiveInfinity, 0, 0, true, "IsInf"},
		{"IsInf(1, 0)", FromFloat32(1), 0, 0, false, "IsInf"},

		// IsNaN
		{"IsNaN(NaN)", QuietNaN, 0, 0, true, "IsNaN"},
		{"IsNaN(1)", FromFloat32(1), 0, 0, false, "IsNaN"},

		// Signbit
		{"Signbit(-1)", FromFloat32(-1), 0, 0, true, "Signbit"},
		{"Signbit(1)", FromFloat32(1), 0, 0, false, "Signbit"},
		{"Signbit(-0)", NegativeZero, 0, 0, true, "Signbit"},

		// FpClassify
		{"FpClassify(NaN)", QuietNaN, 0, 0, ClassQuietNaN, "FpClassify"},
		{"FpClassify(0)", PositiveZero, 0, 0, ClassPositiveZero, "FpClassify"},
		{"FpClassify(Normal)", FromFloat32(1), 0, 0, ClassPositiveNormal, "FpClassify"},
		{"FpClassify(Subnormal)", FromBits(0x0001), 0, 0, ClassPositiveSubnormal, "FpClassify"},
		{"FpClassify(Inf)", PositiveInfinity, 0, 0, ClassPositiveInfinity, "FpClassify"},

		// NextAfter
		{"NextAfter(1, 2)", FromFloat32(1), FromFloat32(2), 0, FromBits(0x3c01), "NextAfter"},
		{"NextAfter(2, 1)", FromFloat32(2), FromFloat32(1), 0, FromBits(0x3fff), "NextAfter"},
		{"NextAfter(0, 1)", PositiveZero, FromFloat32(1), 0, FromBits(0x0001), "NextAfter"},
		{"NextAfter(0, -1)", PositiveZero, FromFloat32(-1), 0, FromBits(0x8001), "NextAfter"},
		{"NextAfter(NaN, 1)", QuietNaN, FromFloat32(1), 0, QuietNaN, "NextAfter"},
		{"NextAfter(1, NaN)", FromFloat32(1), QuietNaN, 0, QuietNaN, "NextAfter"},
		{"NextAfter(1, 1)", FromFloat32(1), FromFloat32(1), 0, FromFloat32(1), "NextAfter"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			switch tc.op {
			case "IsInf":
				res := IsInf(tc.f, tc.sign)
				if res != tc.expected.(bool) {
					t.Errorf("Expected %v, got %v", tc.expected, res)
				}
			case "IsNaN":
				res := IsNaN(tc.f)
				if res != tc.expected.(bool) {
					t.Errorf("Expected %v, got %v", tc.expected, res)
				}
			case "Signbit":
				res := Signbit(tc.f)
				if res != tc.expected.(bool) {
					t.Errorf("Expected %v, got %v", tc.expected, res)
				}
			case "FpClassify":
				res := FpClassify(tc.f)
				if res != tc.expected.(FloatClass) {
					t.Errorf("Expected %v, got %v", tc.expected, res)
				}
			case "NextAfter":
				res := NextAfter(tc.f, tc.f2)
				if res.Bits() != tc.expected.(Float16).Bits() {
					t.Errorf("Expected %v, got %v", tc.expected, res)
				}
			}
		})
	}
}

func TestMiscFloat16Functions(t *testing.T) {
	if GetVersion() != Version {
		t.Errorf("GetVersion: expected %s, got %s", Version, GetVersion())
	}
	if Zero().Bits() != PositiveZero.Bits() {
		t.Errorf("Zero: expected %v, got %v", PositiveZero, Zero())
	}
	if One().Bits() != FromFloat32(1.0).Bits() {
		t.Errorf("One: expected %v, got %v", FromFloat32(1.0), One())
	}
	if Inf(1).Bits() != PositiveInfinity.Bits() {
		t.Errorf("Inf(1): expected %v, got %v", PositiveInfinity, Inf(1))
	}
	if Inf(-1).Bits() != NegativeInfinity.Bits() {
		t.Errorf("Inf(-1): expected %v, got %v", NegativeInfinity, Inf(-1))
	}
	if !IsFinite(FromFloat32(1)) {
		t.Error("IsFinite(1) should be true")
	}
	if !IsNormal(FromFloat32(1)) {
		t.Error("IsNormal(1) should be true")
	}
	if !IsSubnormal(FromBits(0x0001)) {
		t.Error("IsSubnormal(0x0001) should be true")
	}
	if len(GetBenchmarkOperations()) != 4 {
		t.Error("GetBenchmarkOperations should return 4 operations")
	}
}

func TestFrexpLdexp(t *testing.T) {
	f, exp := Frexp(FromFloat32(2.0))
	if f.ToFloat32() != 0.5 || exp != 2 {
		t.Errorf("Frexp(2.0): expected (0.5, 2), got (%v, %v)", f.ToFloat32(), exp)
	}
	res := Ldexp(FromFloat32(0.5), 2)
	if res.ToFloat32() != 2.0 {
		t.Errorf("Ldexp(0.5, 2): expected 2.0, got %v", res.ToFloat32())
	}
}

func TestModf(t *testing.T) {
	i, f := Modf(FromFloat32(3.14))
	if i.ToFloat32() != 3.0 || f.ToFloat32() > 0.15 || f.ToFloat32() < 0.13 {
		t.Errorf("Modf(3.14): expected (3.0, ~0.14), got (%v, %v)", i.ToFloat32(), f.ToFloat32())
	}
}

func TestComputeSliceStats(t *testing.T) {
	stats := ComputeSliceStats([]Float16{FromFloat32(1), FromFloat32(2), FromFloat32(3)})
	if stats.Min.ToFloat32() != 1 || stats.Max.ToFloat32() != 3 || stats.Sum.ToFloat32() != 6 || stats.Mean.ToFloat32() != 2 {
		t.Errorf("ComputeSliceStats: got %+v", stats)
	}
	stats = ComputeSliceStats(nil)
	if stats.Length != 0 {
		t.Error("ComputeSliceStats(nil) should be empty")
	}
}

func TestFastMath(t *testing.T) {
	a, b := FromFloat32(2), FromFloat32(3)
	if FastAdd(a, b).ToFloat32() != 5 {
		t.Error("FastAdd")
	}
	if FastMul(a, b).ToFloat32() != 6 {
		t.Error("FastMul")
	}
}

func TestVectorMath(t *testing.T) {
	a, b := []Float16{FromFloat32(1)}, []Float16{FromFloat32(2)}
	res := VectorAdd(a, b)
	if len(res) != 1 || res[0].ToFloat32() != 3 {
		t.Error("VectorAdd")
	}
	res = VectorMul(a, b)
	if len(res) != 1 || res[0].ToFloat32() != 2 {
		t.Error("VectorMul")
	}
}

func TestValidateSliceLength(t *testing.T) {
	err := ValidateSliceLength([]Float16{1}, []Float16{1})
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	err = ValidateSliceLength([]Float16{1}, []Float16{1, 2})
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}
