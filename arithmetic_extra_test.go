package float16

import (
	"testing"
)

func TestArithmeticExtra(t *testing.T) {
	testCases := []struct {
		name      string
		a         Float16
		b         Float16
		mode      ArithmeticMode
		rounding  RoundingMode
		expected  Float16
		err       bool
		op        string
	}{
		// AddWithMode
		{"AddWithMode(1, 2)", FromFloat32(1), FromFloat32(2), ModeIEEEArithmetic, RoundNearestEven, FromFloat32(3), false, "AddWithMode"},
		{"AddWithMode(Inf, -Inf)", PositiveInfinity, NegativeInfinity, ModeIEEEArithmetic, RoundNearestEven, QuietNaN, false, "AddWithMode"},
		{"AddWithMode(NaN, 1)", QuietNaN, FromFloat32(1), ModeExactArithmetic, RoundNearestEven, 0, true, "AddWithMode"},
		{"AddWithMode(Fast)", FromFloat32(1), FromFloat32(2), ModeFastArithmetic, RoundNearestEven, FromFloat32(3), false, "AddWithMode"},

		// MulWithMode
		{"MulWithMode(2, 3)", FromFloat32(2), FromFloat32(3), ModeIEEEArithmetic, RoundNearestEven, FromFloat32(6), false, "MulWithMode"},
		{"MulWithMode(0, Inf)", PositiveZero, PositiveInfinity, ModeIEEEArithmetic, RoundNearestEven, QuietNaN, false, "MulWithMode"},
		{"MulWithMode(NaN, 1)", QuietNaN, FromFloat32(1), ModeExactArithmetic, RoundNearestEven, 0, true, "MulWithMode"},
		{"MulWithMode(Fast)", FromFloat32(2), FromFloat32(3), ModeFastArithmetic, RoundNearestEven, FromFloat32(6), false, "MulWithMode"},

		// DivWithMode
		{"DivWithMode(6, 3)", FromFloat32(6), FromFloat32(3), ModeIEEEArithmetic, RoundNearestEven, FromFloat32(2), false, "DivWithMode"},
		{"DivWithMode(0, 0)", PositiveZero, PositiveZero, ModeIEEEArithmetic, RoundNearestEven, QuietNaN, false, "DivWithMode"},
		{"DivWithMode(1, 0)", FromFloat32(1), PositiveZero, ModeExactArithmetic, RoundNearestEven, 0, true, "DivWithMode"},
		{"DivWithMode(Inf, Inf)", PositiveInfinity, PositiveInfinity, ModeIEEEArithmetic, RoundNearestEven, QuietNaN, false, "DivWithMode"},
		{"DivWithMode(NaN, 1)", QuietNaN, FromFloat32(1), ModeExactArithmetic, RoundNearestEven, 0, true, "DivWithMode"},
		{"DivWithMode(Fast)", FromFloat32(6), FromFloat32(3), ModeFastArithmetic, RoundNearestEven, FromFloat32(2), false, "DivWithMode"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var res Float16
			var err error
			switch tc.op {
			case "AddWithMode":
				res, err = AddWithMode(tc.a, tc.b, tc.mode, tc.rounding)
			case "MulWithMode":
				res, err = MulWithMode(tc.a, tc.b, tc.mode, tc.rounding)
			case "DivWithMode":
				res, err = DivWithMode(tc.a, tc.b, tc.mode, tc.rounding)
			}

			if (err != nil) != tc.err {
				t.Errorf("Expected error %v, got %v", tc.err, err)
			}
			if tc.expected.IsNaN() {
				if !res.IsNaN() {
					t.Errorf("Expected NaN, got %v", res)
				}
			} else if res.Bits() != tc.expected.Bits() {
				t.Errorf("Expected %v, got %v", tc.expected, res)
			}
		})
	}
}

func TestComparisonExtra(t *testing.T) {
	if !Equal(FromFloat32(1), FromFloat32(1)) {
		t.Error("Equal(1, 1) should be true")
	}
	if Equal(QuietNaN, QuietNaN) {
		t.Error("Equal(NaN, NaN) should be false")
	}
	if !Less(FromFloat32(1), FromFloat32(2)) {
		t.Error("Less(1, 2) should be true")
	}
	if Less(QuietNaN, FromFloat32(1)) {
		t.Error("Less(NaN, 1) should be false")
	}
	if Min(FromFloat32(1), QuietNaN).Bits() != FromFloat32(1).Bits() {
		t.Error("Min(1, NaN) should be 1")
	}
	if Max(FromFloat32(1), QuietNaN).Bits() != FromFloat32(1).Bits() {
		t.Error("Max(1, NaN) should be 1")
	}
}

func TestSliceOpsExtra(t *testing.T) {
	a := []Float16{FromFloat32(1), FromFloat32(2)}
	b := []Float16{FromFloat32(3), FromFloat32(4)}
	add := AddSlice(a, b)
	if add[0].ToFloat32() != 4 || add[1].ToFloat32() != 6 {
		t.Error("AddSlice")
	}
	sub := SubSlice(b, a)
	if sub[0].ToFloat32() != 2 || sub[1].ToFloat32() != 2 {
		t.Error("SubSlice")
	}
	mul := MulSlice(a, b)
	if mul[0].ToFloat32() != 3 || mul[1].ToFloat32() != 8 {
		t.Error("MulSlice")
	}
	div := DivSlice(b, a)
	if div[0].ToFloat32() != 3 || div[1].ToFloat32() != 2 {
		t.Error("DivSlice")
	}
	dot := DotProduct(a, b)
	if dot.ToFloat32() != 11 {
		t.Error("DotProduct")
	}
}
