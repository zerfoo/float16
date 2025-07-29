package float16

import (
	"math"
	"testing"
)

func TestIntConversions(t *testing.T) {
	testCases := []struct {
		name     string
		i32      int32
		i64      int64
		f        Float16
		expected interface{}
		op       interface{}
	}{
		// FromInt32
		{"FromInt32(0)", 0, 0, 0, FromFloat32(0), FromInt32},
		{"FromInt32(1)", 1, 0, 0, FromFloat32(1), FromInt32},
		{"FromInt32(-1)", -1, 0, 0, FromFloat32(-1), FromInt32},

		// FromInt64
		{"FromInt64(0)", 0, 0, 0, FromFloat32(0), FromInt64},
		{"FromInt64(1)", 0, 1, 0, FromFloat32(1), FromInt64},
		{"FromInt64(-1)", 0, -1, 0, FromFloat32(-1), FromInt64},

		// ToInt
		{"ToInt(0)", 0, 0, FromFloat32(0), 0, "ToInt"},
		{"ToInt(1.5)", 0, 0, FromFloat32(1.5), 1, "ToInt"},
		{"ToInt(-1.5)", 0, 0, FromFloat32(-1.5), -1, "ToInt"},

		// ToInt32
		{"ToInt32(0)", 0, 0, FromFloat32(0), int32(0), "ToInt32"},
		{"ToInt32(1.5)", 0, 0, FromFloat32(1.5), int32(1), "ToInt32"},
		{"ToInt32(-1.5)", 0, 0, FromFloat32(-1.5), int32(-1), "ToInt32"},

		// ToInt64
		{"ToInt64(0)", 0, 0, FromFloat32(0), int64(0), "ToInt64"},
		{"ToInt64(1.5)", 0, 0, FromFloat32(1.5), int64(1), "ToInt64"},
		{"ToInt64(-1.5)", 0, 0, FromFloat32(-1.5), int64(-1), "ToInt64"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			switch op := tc.op.(type) {
			case func(int32) Float16:
				res := op(tc.i32)
				if res.Bits() != tc.expected.(Float16).Bits() {
					t.Errorf("Expected %v, got %v", tc.expected, res)
				}
			case func(int64) Float16:
				res := op(tc.i64)
				if res.Bits() != tc.expected.(Float16).Bits() {
					t.Errorf("Expected %v, got %v", tc.expected, res)
				}
			case string:
				switch op {
				case "ToInt":
					res := tc.f.ToInt()
					if res != tc.expected.(int) {
						t.Errorf("Expected %v, got %v", tc.expected, res)
					}
				case "ToInt32":
					res := tc.f.ToInt32()
					if res != tc.expected.(int32) {
						t.Errorf("Expected %v, got %v", tc.expected, res)
					}
				case "ToInt64":
					res := tc.f.ToInt64()
					if res != tc.expected.(int64) {
						t.Errorf("Expected %v, got %v", tc.expected, res)
					}
				}
			}
		})
	}
}

func TestParse(t *testing.T) {
	_, err := Parse("1.0")
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}

func TestFloat32Conversions(t *testing.T) {
	testCases := []struct {
		name     string
		f32      float32
		f16      Float16
		expected interface{}
		op       string
	}{
		// ToFloat16
		{"ToFloat16(0)", 0, 0, FromFloat32(0), "ToFloat16"},
		{"ToFloat16(1)", 1, 0, FromFloat32(1), "ToFloat16"},
		{"ToFloat16(-1)", -1, 0, FromFloat32(-1), "ToFloat16"},
		{"ToFloat16(65504)", 65504, 0, FromFloat32(65504), "ToFloat16"},
		{"ToFloat16(Inf)", float32(math.Inf(1)), 0, PositiveInfinity, "ToFloat16"},
		{"ToFloat16(-Inf)", float32(math.Inf(-1)), 0, NegativeInfinity, "ToFloat16"},
		{"ToFloat16(NaN)", float32(math.NaN()), 0, QuietNaN, "ToFloat16"},

		// FromFloat32
		{"FromFloat32(0)", 0, FromFloat32(0), float32(0), "FromFloat32"},
		{"FromFloat32(1)", 0, FromFloat32(1), float32(1), "FromFloat32"},
		{"FromFloat32(-1)", 0, FromFloat32(-1), float32(-1), "FromFloat32"},
		{"FromFloat32(Inf)", 0, PositiveInfinity, float32(math.Inf(1)), "FromFloat32"},
		{"FromFloat32(-Inf)", 0, NegativeInfinity, float32(math.Inf(-1)), "FromFloat32"},
		{"FromFloat32(NaN)", 0, QuietNaN, float32(math.NaN()), "FromFloat32"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			switch tc.op {
			case "ToFloat16":
				res := ToFloat16(tc.f32)
				if tc.expected.(Float16).IsNaN() {
					if !res.IsNaN() {
						t.Errorf("Expected NaN, got %v", res)
					}
				} else if res.Bits() != tc.expected.(Float16).Bits() {
					t.Errorf("Expected %v, got %v", tc.expected, res)
				}
			case "FromFloat32":
				res := tc.f16.ToFloat32()
				if math.IsNaN(float64(tc.expected.(float32))) {
					if !math.IsNaN(float64(res)) {
						t.Errorf("Expected NaN, got %v", res)
					}
				} else if res != tc.expected.(float32) {
					t.Errorf("Expected %v, got %v", tc.expected, res)
				}
			}
		})
	}
}

func TestFloat64Conversions(t *testing.T) {
	testCases := []struct {
		name     string
		f64      float64
		f16      Float16
		expected interface{}
		op       string
	}{
		// FromFloat64
		{"FromFloat64(0)", 0, 0, FromFloat64(0), "FromFloat64"},
		{"FromFloat64(1)", 1, 0, FromFloat64(1), "FromFloat64"},
		{"FromFloat64(-1)", -1, 0, FromFloat64(-1), "FromFloat64"},
		{"FromFloat64(65504)", 65504, 0, FromFloat64(65504), "FromFloat64"},
		{"FromFloat64(Inf)", math.Inf(1), 0, PositiveInfinity, "FromFloat64"},
		{"FromFloat64(-Inf)", math.Inf(-1), 0, NegativeInfinity, "FromFloat64"},
		{"FromFloat64(NaN)", math.NaN(), 0, QuietNaN, "FromFloat64"},

		// ToFloat64
		{"ToFloat64(0)", 0, FromFloat64(0), float64(0), "ToFloat64"},
		{"ToFloat64(1)", 0, FromFloat64(1), float64(1), "ToFloat64"},
		{"ToFloat64(-1)", 0, FromFloat64(-1), float64(-1), "ToFloat64"},
		{"ToFloat64(Inf)", 0, PositiveInfinity, math.Inf(1), "ToFloat64"},
		{"ToFloat64(-Inf)", 0, NegativeInfinity, math.Inf(-1), "ToFloat64"},
		{"ToFloat64(NaN)", 0, QuietNaN, math.NaN(), "ToFloat64"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			switch tc.op {
			case "FromFloat64":
				res := FromFloat64(tc.f64)
				if tc.expected.(Float16).IsNaN() {
					if !res.IsNaN() {
						t.Errorf("Expected NaN, got %v", res)
					}
				} else if res.Bits() != tc.expected.(Float16).Bits() {
					t.Errorf("Expected %v, got %v", tc.expected, res)
				}
			case "ToFloat64":
				res := tc.f16.ToFloat64()
				if math.IsNaN(tc.expected.(float64)) {
					if !math.IsNaN(res) {
						t.Errorf("Expected NaN, got %v", res)
					}
				} else if res != tc.expected.(float64) {
					t.Errorf("Expected %v, got %v", tc.expected, res)
				}
			}
		})
	}
}

func TestFromFloat64WithModeExtra(t *testing.T) {
	testCases := []struct {
		name      string
		f64       float64
		convMode  ConversionMode
		roundMode RoundingMode
		expected  Float16
		err       bool
	}{
		{"Normal", 1.0, ModeIEEE, RoundNearestEven, FromFloat64(1.0), false},
		{"Strict Inf", math.Inf(1), ModeStrict, RoundNearestEven, 0, true},
		{"Strict NaN", math.NaN(), ModeStrict, RoundNearestEven, 0, true},
		{"Overflow", 70000.0, ModeIEEE, RoundNearestEven, PositiveInfinity, false},
		{"Underflow", 1e-40, ModeIEEE, RoundNearestEven, PositiveZero, false},
		{"Underflow Strict", 1e-40, ModeStrict, RoundNearestEven, 0, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := FromFloat64WithMode(tc.f64, tc.convMode, tc.roundMode)
			if (err != nil) != tc.err {
				t.Errorf("Expected error %v, got %v", tc.err, err)
			}
			if res.Bits() != tc.expected.Bits() {
				t.Errorf("Expected %v, got %v", tc.expected, res)
			}
		})
	}
}

func TestToFloat16WithMode(t *testing.T) {
	testCases := []struct {
		name      string
		f32       float32
		convMode  ConversionMode
		roundMode RoundingMode
		expected  Float16
		err       bool
	}{
		{"Normal", 1.0, ModeIEEE, RoundNearestEven, FromFloat32(1.0), false},
		{"Strict Inf", float32(math.Inf(1)), ModeStrict, RoundNearestEven, 0, true},
		{"Strict NaN", float32(math.NaN()), ModeStrict, RoundNearestEven, 0, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := ToFloat16WithMode(tc.f32, tc.convMode, tc.roundMode)
			if (err != nil) != tc.err {
				t.Errorf("Expected error %v, got %v", tc.err, err)
			}
			if res.Bits() != tc.expected.Bits() {
				t.Errorf("Expected %v, got %v", tc.expected, res)
			}
		})
	}
}

func TestShouldRound(t *testing.T) {
	testCases := []struct {
		name      string
		mantissa  uint32
		shift     int
		mode      RoundingMode
		expected  bool
	}{
		{"NearestEven_NoRound", 0x1234, 4, RoundNearestEven, false},
		{"NearestEven_Round", 0x1238, 4, RoundNearestEven, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res := shouldRound(tc.mantissa, tc.shift, tc.mode)
			if res != tc.expected {
				t.Errorf("Expected %v, got %v", tc.expected, res)
			}
		})
	}
}

func TestSliceConversions(t *testing.T) {
	// Test ToSlice16
	f32s := []float32{0, 1, -1, 65504}
	f16s := ToSlice16(f32s)
	if len(f16s) != len(f32s) {
		t.Errorf("ToSlice16: expected length %d, got %d", len(f32s), len(f16s))
	}
	expectedF16s := []Float16{FromFloat32(0), FromFloat32(1), FromFloat32(-1), FromFloat32(65504)}
	for i := range f16s {
		if f16s[i].Bits() != expectedF16s[i].Bits() {
			t.Errorf("ToSlice16: at index %d, expected %v, got %v", i, expectedF16s[i], f16s[i])
		}
	}

	// Test ToSlice32
	f32sBack := ToSlice32(f16s)
	if len(f32sBack) != len(f16s) {
		t.Errorf("ToSlice32: expected length %d, got %d", len(f16s), len(f32sBack))
	}
	for i := range f32sBack {
		if f32sBack[i] != f32s[i] {
			t.Errorf("ToSlice32: at index %d, expected %v, got %v", i, f32s[i], f32sBack[i])
		}
	}

	// Test ToSlice64 and FromSlice64
	f64s := []float64{0, 1, -1, 65504}
	f16sFrom64 := FromSlice64(f64s)
	if len(f16sFrom64) != len(f64s) {
		t.Errorf("FromSlice64: expected length %d, got %d", len(f64s), len(f16sFrom64))
	}
	expectedF16sFrom64 := []Float16{FromFloat64(0), FromFloat64(1), FromFloat64(-1), FromFloat64(65504)}
	for i := range f16sFrom64 {
		if f16sFrom64[i].Bits() != expectedF16sFrom64[i].Bits() {
			t.Errorf("FromSlice64: at index %d, expected %v, got %v", i, expectedF16sFrom64[i], f16sFrom64[i])
		}
	}

	f64sBack := ToSlice64(f16sFrom64)
	if len(f64sBack) != len(f16sFrom64) {
		t.Errorf("ToSlice64: expected length %d, got %d", len(f16sFrom64), len(f64sBack))
	}
	for i := range f64sBack {
		if f64sBack[i] != f64s[i] {
			t.Errorf("ToSlice64: at index %d, expected %v, got %v", i, f64s[i], f64sBack[i])
		}
	}

	// Test empty slices
	if ToSlice16(nil) != nil {
		t.Error("ToSlice16(nil) should be nil")
	}
	if ToSlice32(nil) != nil {
		t.Error("ToSlice32(nil) should be nil")
	}
	if ToSlice64(nil) != nil {
		t.Error("ToSlice64(nil) should be nil")
	}
	if FromSlice64(nil) != nil {
		t.Error("FromSlice64(nil) should be nil")
	}
}

func TestToSlice16WithModeExtra(t *testing.T) {
	f32s := []float32{1.0, float32(math.Inf(1)), float32(math.NaN())}
	_, errs := ToSlice16WithMode(f32s, ModeStrict, RoundNearestEven)
	if len(errs) != 2 {
		t.Errorf("Expected 2 errors, got %d", len(errs))
	}

	f16s, errs := ToSlice16WithMode(f32s, ModeIEEE, RoundNearestEven)
	if len(errs) != 0 {
		t.Errorf("Expected 0 errors, got %d", len(errs))
	}
	expectedF16s := []Float16{FromFloat32(1.0), PositiveInfinity, QuietNaN}
	for i := range f16s {
		if expectedF16s[i].IsNaN() {
			if !f16s[i].IsNaN() {
				t.Errorf("ToSlice16WithMode: at index %d, expected NaN, got %v", i, f16s[i])
			}
		} else if f16s[i].Bits() != expectedF16s[i].Bits() {
			t.Errorf("ToSlice16WithMode: at index %d, expected %v, got %v", i, expectedF16s[i], f16s[i])
		}
	}

	f16s, errs = ToSlice16WithMode(nil, ModeIEEE, RoundNearestEven)
	if f16s != nil || errs != nil {
		t.Error("ToSlice16WithMode(nil) should be (nil, nil)")
	}
}
