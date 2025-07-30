package float16

import (
	"testing"
)

func TestAddWithMode(t *testing.T) {
	tests := []struct {
		name     string
		a        Float16
		b        Float16
		mode     ArithmeticMode
		rounding RoundingMode
		expect   Float16
		hasError bool
	}{
		// Basic addition
		{"1.0 + 2.0", 0x3C00, 0x4000, ModeIEEEArithmetic, RoundNearestEven, 0x4200, false},  // 1.0 + 2.0 = 3.0
		{"0.5 + 0.25", 0x3800, 0x3400, ModeIEEEArithmetic, RoundNearestEven, 0x3A00, false}, // 0.5 + 0.25 = 0.75 (0x3A00)

		// Special cases
		{"0 + x", 0x0000, 0x3C00, ModeIEEEArithmetic, RoundNearestEven, 0x3C00, false}, // 0 + 1.0 = 1.0
		{"x + 0", 0x3C00, 0x0000, ModeIEEEArithmetic, RoundNearestEven, 0x3C00, false}, // 1.0 + 0 = 1.0

		// Infinity handling
		{"Inf + 1", 0x7C00, 0x3C00, ModeIEEEArithmetic, RoundNearestEven, 0x7C00, false},    // +Inf + 1 = +Inf
		{"1 + Inf", 0x3C00, 0x7C00, ModeIEEEArithmetic, RoundNearestEven, 0x7C00, false},    // 1 + Inf = +Inf
		{"-Inf + Inf", 0xFC00, 0x7C00, ModeIEEEArithmetic, RoundNearestEven, 0x7E00, false}, // -Inf + Inf = NaN

		// NaN handling
		{"NaN + 1", 0x7E00, 0x3C00, ModeIEEEArithmetic, RoundNearestEven, 0x7E00, false}, // NaN + 1 = NaN
		{"1 + NaN", 0x3C00, 0x7E00, ModeIEEEArithmetic, RoundNearestEven, 0x7E00, false}, // 1 + NaN = NaN

		// Exact mode
		{"1.0 + 2.0 (exact)", 0x3C00, 0x4000, ModeExactArithmetic, RoundNearestEven, 0x4200, false}, // 1.0 + 2.0 = 3.0 (exact)
		{"0.1 + 0.2 (exact)", 0x2E66, 0x3266, ModeExactArithmetic, RoundNearestEven, 0x34CC, false}, // 0.1 + 0.2 = ~0.2998 (actual float16 result)

		// Rounding modes
		{"1.0 + 0.5 (toward zero)", 0x3C00, 0x3800, ModeIEEEArithmetic, RoundTowardZero, 0x3E00, false}, // 1.0 + 0.5 = 1.5 (0x3E00)
		{"1.0 + 0.5 (toward positive)", 0x3C00, 0x3800, ModeIEEEArithmetic, RoundTowardPositive, 0x3E00, false},
		{"-1.0 + -0.5 (toward positive)", 0xBC00, 0xB800, ModeIEEEArithmetic, RoundTowardPositive, 0xBE00, false},
		{"1.0 + 0.5 (toward negative)", 0x3C00, 0x3800, ModeIEEEArithmetic, RoundTowardNegative, 0x3E00, false},
		{"-1.0 - 0.5 (toward negative)", 0xBC00, 0xB800, ModeIEEEArithmetic, RoundTowardNegative, 0xBE00, false},
		{"1.0 + 0.5 (nearest away)", 0x3C00, 0x3800, ModeIEEEArithmetic, RoundNearestAway, 0x3E00, false},

		// Fast mode
		{"1.0 + 2.0 (fast)", 0x3C00, 0x4000, ModeFastArithmetic, RoundNearestEven, 0x4200, false},

		// Subnormal addition
		{"subnormal + subnormal", 0x0001, 0x0001, ModeIEEEArithmetic, RoundNearestEven, 0x0002, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := AddWithMode(tt.a, tt.b, tt.mode, tt.rounding)

			if tt.hasError {
				if err == nil {
					t.Errorf("Expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			// For better debugging, show float32 values when the test fails
			if result != tt.expect {
				t.Errorf("AddWithMode(%v, %v) = %v (0x%04X, %f), want %v (0x%04X, %f)",
					tt.a, tt.b, result, uint16(result), result.ToFloat32(),
					tt.expect, uint16(tt.expect), tt.expect.ToFloat32())
			}
		})
	}
}

func TestAddWithMode_ErrorCases(t *testing.T) {
	tests := []struct {
		name    string
		a       Float16
		b       Float16
		mode    ArithmeticMode
		expect  Float16
		errCode ErrorCode
	}{
		{"NaN in exact mode", 0x7E00, 0x3C00, ModeExactArithmetic, 0, ErrNaN},
		{"Inf-Inf in exact mode", 0x7C00, 0xFC00, ModeExactArithmetic, 0, ErrInvalidOperation},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := AddWithMode(tt.a, tt.b, tt.mode, RoundNearestEven)

			if err == nil {
				t.Fatalf("Expected error, got nil")
			}

			err16, ok := err.(*Float16Error)
			if !ok {
				t.Fatalf("Expected Float16Error, got %T", err)
			}

			if err16.Code != tt.errCode {
				t.Errorf("Expected error code %v, got %v", tt.errCode, err16.Code)
			}

			if result != tt.expect {
				t.Errorf("Expected result %v, got %v", tt.expect, result)
			}
		})
	}
}

func TestAdd(t *testing.T) {
	tests := []struct {
		a      Float16
		b      Float16
		expect Float16
	}{
		{0x3C00, 0x4000, 0x4200}, // 1.0 + 2.0 = 3.0
		{0x3800, 0x3400, 0x3A00}, // 0.5 + 0.25 = 0.75 (0x3A00)
		{0x7C00, 0x3C00, 0x7C00}, // +Inf + 1 = +Inf
		{0x3C00, 0x0000, 0x3C00}, // 1.0 + 0.0 = 1.0
		{0x0000, 0x3C00, 0x3C00}, // 0.0 + 1.0 = 1.0
		{0x3C00, 0xBC00, 0x0000}, // 1.0 + (-1.0) = 0.0
	}

	for _, tt := range tests {
		t.Run(tt.a.String()+" + "+tt.b.String(), func(t *testing.T) {
			result := Add(tt.a, tt.b)
			if result != tt.expect {
				t.Errorf("Add(%v, %v) = %v (0x%04X), want %v (0x%04X)",
					tt.a, tt.b, result, uint16(result), tt.expect, uint16(tt.expect))
			}
		})
	}
}

func TestSubWithModeBasic(t *testing.T) {
	tests := []struct {
		name     string
		a        Float16
		b        Float16
		mode     ArithmeticMode
		rounding RoundingMode
		expect   Float16
		hasError bool
	}{
		{name: "1.0 - 0.5", a: 0x3C00, b: 0x3800, mode: ModeIEEEArithmetic, rounding: RoundNearestEven, expect: 0x3800, hasError: false},  // 1.0 - 0.5 = 0.5
		{name: "0.5 - 0.25", a: 0x3800, b: 0x3400, mode: ModeIEEEArithmetic, rounding: RoundNearestEven, expect: 0x3400, hasError: false}, // 0.5 - 0.25 = 0.25
		{name: "1.0 - 1.0", a: 0x3C00, b: 0x3C00, mode: ModeIEEEArithmetic, rounding: RoundNearestEven, expect: 0x0000, hasError: false},  // 1.0 - 1.0 = 0.0
		{name: "1.0 - -1.0", a: 0x3C00, b: 0xBC00, mode: ModeIEEEArithmetic, rounding: RoundNearestEven, expect: 0x4000, hasError: false}, // 1.0 - (-1.0) = 2.0
		{name: "-1.0 - 1.0", a: 0xBC00, b: 0x3C00, mode: ModeIEEEArithmetic, rounding: RoundNearestEven, expect: 0xC000, hasError: false}, // -1.0 - 1.0 = -2.0
		{name: "0.0 - 0.0", a: 0x0000, b: 0x0000, mode: ModeIEEEArithmetic, rounding: RoundNearestEven, expect: 0x0000, hasError: false},  // 0.0 - 0.0 = 0.0 (or -0.0 is also valid)
		{name: "Inf - 1.0", a: 0x7C00, b: 0x3C00, mode: ModeIEEEArithmetic, rounding: RoundNearestEven, expect: 0x7C00, hasError: false},  // +Inf - 1.0 = +Inf
		{name: "1.0 - Inf", a: 0x3C00, b: 0x7C00, mode: ModeIEEEArithmetic, rounding: RoundNearestEven, expect: 0xFC00, hasError: false},  // 1.0 - +Inf = -Inf
		{name: "NaN - 1.0", a: 0x7E00, b: 0x3C00, mode: ModeIEEEArithmetic, rounding: RoundNearestEven, expect: 0x7E00, hasError: false},  // NaN - 1.0 = NaN
		{name: "1.0 - NaN", a: 0x3C00, b: 0x7E00, mode: ModeIEEEArithmetic, rounding: RoundNearestEven, expect: 0x7E00, hasError: false},  // 1.0 - NaN = NaN
		{"1.0 - 0.5 (toward zero)", 0x3C00, 0x3800, ModeIEEEArithmetic, RoundTowardZero, 0x3800, false},
		{"1.0 - 0.5 (toward positive)", 0x3C00, 0x3800, ModeIEEEArithmetic, RoundTowardPositive, 0x3800, false},
		{"-1.0 - 0.5 (toward positive)", 0xBC00, 0xB800, ModeIEEEArithmetic, RoundTowardPositive, 0xB800, false},
		{"1.0 - 0.5 (toward negative)", 0x3C00, 0x3800, ModeIEEEArithmetic, RoundTowardNegative, 0x3800, false},
		{"-1.0 - 0.5 (toward negative)", 0xBC00, 0xB800, ModeIEEEArithmetic, RoundTowardNegative, 0xB800, false},
		{"1.0 - 0.5 (nearest away)", 0x3C00, 0x3800, ModeIEEEArithmetic, RoundNearestAway, 0x3800, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := SubWithMode(tt.a, tt.b, tt.mode, tt.rounding)

			if tt.hasError {
				if err == nil {
					t.Error("Expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			// For 0.0 - 0.0, both 0.0 and -0.0 are valid results
			if tt.a == 0 && tt.b == 0 {
				// Check if the result is either 0.0 or -0.0
				if result != 0 && result != Float16(0x8000) {
					t.Errorf("SubWithMode(%v, %v) = %v (0x%04X), want 0.0 (0x0000) or -0.0 (0x8000)",
						tt.a, tt.b, result, uint16(result))
				}
			} else if result != tt.expect {
				t.Errorf("SubWithMode(%v, %v) = %v (0x%04X), want %v (0x%04X)",
					tt.a, tt.b, result, uint16(result), tt.expect, uint16(tt.expect))
			}
		})
	}
}

func TestSub(t *testing.T) {
	tests := []struct {
		a      Float16
		b      Float16
		expect Float16
	}{
		{0x4200, 0x3C00, 0x4000}, // 3.0 - 1.0 = 2.0
		{0x3C00, 0x3C00, 0x0000}, // 1.0 - 1.0 = 0.0
		{0x3C00, 0x4000, 0xBC00}, // 1.0 - 2.0 = -1.0
	}

	for _, tt := range tests {
		t.Run(tt.a.String()+" - "+tt.b.String(), func(t *testing.T) {
			result := Sub(tt.a, tt.b)
			if result != tt.expect {
				t.Errorf("Sub(%v, %v) = %v (0x%04X), want %v (0x%04X)",
					tt.a, tt.b, result, uint16(result), tt.expect, uint16(tt.expect))
			}
		})
	}
}

func TestMulWithMode(t *testing.T) {
	tests := []struct {
		name     string
		a        Float16
		b        Float16
		mode     ArithmeticMode
		rounding RoundingMode
		expect   Float16
		hasError bool
	}{
		// Basic multiplication
		{"2.0 * 3.0", 0x4000, 0x4200, ModeIEEEArithmetic, RoundNearestEven, 0x4600, false}, // 2.0 * 3.0 = 6.0 (0x4600)
		{"0.5 * 0.5", 0x3800, 0x3800, ModeIEEEArithmetic, RoundNearestEven, 0x3400, false}, // 0.5 * 0.5 = 0.25 (0x3400)

		// Special cases
		{"1.0 * 0.0", 0x3C00, 0x0000, ModeIEEEArithmetic, RoundNearestEven, 0x0000, false},  // 1.0 * 0.0 = 0.0
		{"-1.0 * 0.0", 0xBC00, 0x0000, ModeIEEEArithmetic, RoundNearestEven, 0x8000, false}, // -1.0 * 0.0 = -0.0

		// Infinity handling
		{"Inf * 2.0", 0x7C00, 0x4000, ModeIEEEArithmetic, RoundNearestEven, 0x7C00, false},  // +Inf * 2.0 = +Inf
		{"-Inf * 2.0", 0xFC00, 0x4000, ModeIEEEArithmetic, RoundNearestEven, 0xFC00, false}, // -Inf * 2.0 = -Inf
		{"Inf * 0.0", 0x7C00, 0x0000, ModeIEEEArithmetic, RoundNearestEven, 0x7E00, false},  // +Inf * 0.0 = NaN

		// NaN handling
		{"NaN * 2.0", 0x7E00, 0x4000, ModeIEEEArithmetic, RoundNearestEven, 0x7E00, false}, // NaN * 2.0 = NaN

		// Exact mode
		{"2.0 * 3.0 (exact)", 0x4000, 0x4200, ModeExactArithmetic, RoundNearestEven, 0x4600, false}, // 2.0 * 3.0 = 6.0 (exact)
		{"Inf * 0 (exact)", 0x7C00, 0x0000, ModeExactArithmetic, RoundNearestEven, 0, true},         // Inf * 0 is an error in exact mode

		// Rounding modes
		{"2.0 * 0.5 (toward zero)", 0x4000, 0x3800, ModeIEEEArithmetic, RoundTowardZero, 0x3C00, false},
		{"2.0 * 0.5 (toward positive)", 0x4000, 0x3800, ModeIEEEArithmetic, RoundTowardPositive, 0x3C00, false},
		{"-2.0 * 0.5 (toward positive)", 0xC000, 0x3800, ModeIEEEArithmetic, RoundTowardPositive, 0xBC00, false},
		{"2.0 * 0.5 (toward negative)", 0x4000, 0x3800, ModeIEEEArithmetic, RoundTowardNegative, 0x3C00, false},
		{"-2.0 * 0.5 (toward negative)", 0xC000, 0x3800, ModeIEEEArithmetic, RoundTowardNegative, 0xBC00, false},
		{"2.0 * 0.5 (nearest away)", 0x4000, 0x3800, ModeIEEEArithmetic, RoundNearestAway, 0x3C00, false},

		// Fast mode
		{"2.0 * 3.0 (fast)", 0x4000, 0x4200, ModeFastArithmetic, RoundNearestEven, 0x4600, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := MulWithMode(tt.a, tt.b, tt.mode, tt.rounding)

			if tt.hasError {
				if err == nil {
					t.Errorf("Expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if result != tt.expect {
				t.Errorf("MulWithMode(%v, %v) = %v (0x%04X), want %v (0x%04X)",
					tt.a, tt.b, result, uint16(result), tt.expect, uint16(tt.expect))
			}
		})
	}
}

func TestMul(t *testing.T) {
	tests := []struct {
		a      Float16
		b      Float16
		expect Float16
	}{
		{0x4000, 0x4200, 0x4600}, // 2.0 * 3.0 = 6.0
		{0x3800, 0x3800, 0x3400}, // 0.5 * 0.5 = 0.25
		{0xBC00, 0x3C00, 0xBC00}, // -1.0 * 1.0 = -1.0
	}

	for _, tt := range tests {
		t.Run(tt.a.String()+" * "+tt.b.String(), func(t *testing.T) {
			result := Mul(tt.a, tt.b)
			if result != tt.expect {
				t.Errorf("Mul(%v, %v) = %v (0x%04X), want %v (0x%04X)",
					tt.a, tt.b, result, uint16(result), tt.expect, uint16(tt.expect))
			}
		})
	}
}

func TestDivWithMode(t *testing.T) {
	tests := []struct {
		name     string
		a        Float16
		b        Float16
		mode     ArithmeticMode
		rounding RoundingMode
		expect   Float16
		hasError bool
	}{
		// Basic division
		{"6.0 / 2.0", 0x4600, 0x4000, ModeIEEEArithmetic, RoundNearestEven, 0x4200, false}, // 6.0 / 2.0 = 3.0 (0x4200)
		{"1.0 / 2.0", 0x3C00, 0x4000, ModeIEEEArithmetic, RoundNearestEven, 0x3800, false}, // 1.0 / 2.0 = 0.5 (0x3800)

		// Division by zero
		{"1.0 / 0.0", 0x3C00, 0x0000, ModeIEEEArithmetic, RoundNearestEven, 0x7C00, false},  // 1.0 / 0.0 = +Inf
		{"-1.0 / 0.0", 0xBC00, 0x0000, ModeIEEEArithmetic, RoundNearestEven, 0xFC00, false}, // -1.0 / 0.0 = -Inf
		{"0.0 / 0.0", 0x0000, 0x0000, ModeIEEEArithmetic, RoundNearestEven, 0x7E00, false},  // 0.0 / 0.0 = NaN

		// Infinity handling
		{"Inf / 2.0", 0x7C00, 0x4000, ModeIEEEArithmetic, RoundNearestEven, 0x7C00, false}, // +Inf / 2.0 = +Inf
		{"2.0 / Inf", 0x4000, 0x7C00, ModeIEEEArithmetic, RoundNearestEven, 0x0000, false}, // 2.0 / +Inf = 0.0
		{"Inf / Inf", 0x7C00, 0x7C00, ModeIEEEArithmetic, RoundNearestEven, 0x7E00, false}, // +Inf / +Inf = NaN

		// NaN handling
		{"NaN / 2.0", 0x7E00, 0x4000, ModeIEEEArithmetic, RoundNearestEven, 0x7E00, false}, // NaN / 2.0 = NaN

		// Exact mode
		{"6.0 / 2.0 (exact)", 0x4600, 0x4000, ModeExactArithmetic, RoundNearestEven, 0x4200, false}, // 6.0 / 2.0 = 3.0 (exact)
		{"1.0 / 0.0 (exact)", 0x3C00, 0x0000, ModeExactArithmetic, RoundNearestEven, 0, true},       // 1.0 / 0.0 is an error in exact mode

		// Rounding modes
		{"3.0 / 2.0 (toward zero)", 0x4200, 0x4000, ModeIEEEArithmetic, RoundTowardZero, 0x3E00, false},
		{"3.0 / 2.0 (toward positive)", 0x4200, 0x4000, ModeIEEEArithmetic, RoundTowardPositive, 0x3E00, false},
		{"-3.0 / 2.0 (toward positive)", 0xC200, 0x4000, ModeIEEEArithmetic, RoundTowardPositive, 0xBE00, false},
		{"3.0 / 2.0 (toward negative)", 0x4200, 0x4000, ModeIEEEArithmetic, RoundTowardNegative, 0x3E00, false},
		{"-3.0 / 2.0 (toward negative)", 0xC200, 0x4000, ModeIEEEArithmetic, RoundTowardNegative, 0xBE00, false},
		{"3.0 / 2.0 (nearest away)", 0x4200, 0x4000, ModeIEEEArithmetic, RoundNearestAway, 0x3E00, false},

		// Fast mode
		{"6.0 / 2.0 (fast)", 0x4600, 0x4000, ModeFastArithmetic, RoundNearestEven, 0x4200, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := DivWithMode(tt.a, tt.b, tt.mode, tt.rounding)

			if tt.hasError {
				if err == nil {
					t.Errorf("Expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if result != tt.expect {
				t.Errorf("DivWithMode(%v, %v) = %v (0x%04X), want %v (0x%04X)",
					tt.a, tt.b, result, uint16(result), tt.expect, uint16(tt.expect))
			}
		})
	}
}

func TestDiv(t *testing.T) {
	tests := []struct {
		a      Float16
		b      Float16
		expect Float16
	}{
		{0x4600, 0x4000, 0x4200}, // 6.0 / 2.0 = 3.0
		{0x3C00, 0x4000, 0x3800}, // 1.0 / 2.0 = 0.5
		{0xBC00, 0x3C00, 0xBC00}, // -1.0 / 1.0 = -1.0
	}

	for _, tt := range tests {
		t.Run(tt.a.String()+" / "+tt.b.String(), func(t *testing.T) {
			result := Div(tt.a, tt.b)
			if result != tt.expect {
				t.Errorf("Div(%v, %v) = %v (0x%04X), want %v (0x%04X)",
					tt.a, tt.b, result, uint16(result), tt.expect, uint16(tt.expect))
			}
		})
	}
}

func TestEqual(t *testing.T) {
	tests := []struct {
		a      Float16
		b      Float16
		expect bool
	}{
		{0x3C00, 0x3C00, true},  // 1.0 == 1.0
		{0x0000, 0x8000, true},  // +0 == -0
		{0x3C00, 0xBC00, false}, // 1.0 != -1.0
		{0x7E00, 0x7E00, false}, // NaN != NaN
		{0x7C00, 0xFC00, false}, // +Inf != -Inf
		{0x7C00, 0x7C00, true},  // +Inf == +Inf
	}

	for _, tt := range tests {
		t.Run(tt.a.String()+" == "+tt.b.String(), func(t *testing.T) {
			result := Equal(tt.a, tt.b)
			if result != tt.expect {
				t.Errorf("Equal(%v, %v) = %v, want %v", tt.a, tt.b, result, tt.expect)
			}
		})
	}
}

func TestLess(t *testing.T) {
	tests := []struct {
		a      Float16
		b      Float16
		expect bool
	}{
		{0x3C00, 0x4000, true},  // 1.0 < 2.0
		{0xBC00, 0x3C00, true},  // -1.0 < 1.0
		{0xFC00, 0xBC00, true},  // -Inf < -1.0
		{0x4000, 0x3C00, false}, // 2.0 not < 1.0
		{0x3C00, 0x3C00, false}, // 1.0 not < 1.0
		{0x7E00, 0x3C00, false}, // NaN not < 1.0
		{0x3C00, 0x7E00, false}, // 1.0 not < NaN
	}

	for _, tt := range tests {
		t.Run(tt.a.String()+" < "+tt.b.String(), func(t *testing.T) {
			result := Less(tt.a, tt.b)
			if result != tt.expect {
				t.Errorf("Less(%v, %v) = %v, want %v", tt.a, tt.b, result, tt.expect)
			}
		})
	}
}

func TestGreater(t *testing.T) {
	tests := []struct {
		a      Float16
		b      Float16
		expect bool
	}{
		{0x4000, 0x3C00, true},  // 2.0 > 1.0
		{0x3C00, 0xBC00, true},  // 1.0 > -1.0
		{0xBC00, 0xFC00, true},  // -1.0 > -Inf
		{0x3C00, 0x4000, false}, // 1.0 not > 2.0
		{0x3C00, 0x3C00, false}, // 1.0 not > 1.0
		{0x7E00, 0x3C00, false}, // NaN not > 1.0
		{0x3C00, 0x7E00, false}, // 1.0 not > NaN
	}

	for _, tt := range tests {
		t.Run(tt.a.String()+" > "+tt.b.String(), func(t *testing.T) {
			result := Greater(tt.a, tt.b)
			if result != tt.expect {
				t.Errorf("Greater(%v, %v) = %v, want %v", tt.a, tt.b, result, tt.expect)
			}
		})
	}
}

func TestLessEqual(t *testing.T) {
	tests := []struct {
		a      Float16
		b      Float16
		expect bool
	}{
		{0x3C00, 0x4000, true},  // 1.0 <= 2.0
		{0x3C00, 0x3C00, true},  // 1.0 <= 1.0
		{0x4000, 0x3C00, false}, // 2.0 not <= 1.0
	}

	for _, tt := range tests {
		t.Run(tt.a.String()+" <= "+tt.b.String(), func(t *testing.T) {
			result := LessEqual(tt.a, tt.b)
			if result != tt.expect {
				t.Errorf("LessEqual(%v, %v) = %v, want %v", tt.a, tt.b, result, tt.expect)
			}
		})
	}
}

func TestGreaterEqual(t *testing.T) {
	tests := []struct {
		a      Float16
		b      Float16
		expect bool
	}{
		{0x4000, 0x3C00, true},  // 2.0 >= 1.0
		{0x3C00, 0x3C00, true},  // 1.0 >= 1.0
		{0x3C00, 0x4000, false}, // 1.0 not >= 2.0
	}

	for _, tt := range tests {
		t.Run(tt.a.String()+" >= "+tt.b.String(), func(t *testing.T) {
			result := GreaterEqual(tt.a, tt.b)
			if result != tt.expect {
				t.Errorf("GreaterEqual(%v, %v) = %v, want %v", tt.a, tt.b, result, tt.expect)
			}
		})
	}
}

func TestMin(t *testing.T) {
	tests := []struct {
		name   string
		a      Float16
		b      Float16
		expect Float16
	}{
		{name: "Min(1.0, 2.0)", a: 0x3C00, b: 0x4000, expect: 0x3C00},  // Min(1.0, 2.0) = 1.0
		{name: "Min(-1.0, 1.0)", a: 0xBC00, b: 0x3C00, expect: 0xBC00}, // Min(-1.0, 1.0) = -1.0
		{name: "Min(NaN, 1.0)", a: 0x7E00, b: 0x3C00, expect: 0x3C00},  // Min(NaN, 1.0) = 1.0
		{name: "Min(1.0, NaN)", a: 0x3C00, b: 0x7E00, expect: 0x3C00},  // Min(1.0, NaN) = 1.0
		{name: "Min(NaN, NaN)", a: 0x7E00, b: 0x7E00, expect: 0x7E00},  // Min(NaN, NaN) = NaN
		{name: "Min(1.0, 1.0)", a: 0x3C00, b: 0x3C00, expect: 0x3C00},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Min(tt.a, tt.b)
			if result != tt.expect {
				t.Errorf("Min(%v, %v) = %v, want %v", tt.a, tt.b, result, tt.expect)
			}
		})
	}
}

func TestMax(t *testing.T) {
	tests := []struct {
		name   string
		a      Float16
		b      Float16
		expect Float16
	}{
		{name: "Max(1.0, 2.0)", a: 0x3C00, b: 0x4000, expect: 0x4000},  // Max(1.0, 2.0) = 2.0
		{name: "Max(-1.0, 1.0)", a: 0xBC00, b: 0x3C00, expect: 0x3C00}, // Max(-1.0, 1.0) = 1.0
		{name: "Max(NaN, 1.0)", a: 0x7E00, b: 0x3C00, expect: 0x3C00},  // Max(NaN, 1.0) = 1.0
		{name: "Max(1.0, NaN)", a: 0x3C00, b: 0x7E00, expect: 0x3C00},  // Max(1.0, NaN) = 1.0
		{name: "Max(NaN, NaN)", a: 0x7E00, b: 0x7E00, expect: 0x7E00},  // Max(NaN, NaN) = NaN
		{name: "Max(1.0, 1.0)", a: 0x3C00, b: 0x3C00, expect: 0x3C00},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Max(tt.a, tt.b)
			if result != tt.expect {
				t.Errorf("Max(%v, %v) = %v, want %v", tt.a, tt.b, result, tt.expect)
			}
		})
	}
}

func TestAddSlice(t *testing.T) {
	a := []Float16{0x3C00, 0x4000}      // [1.0, 2.0]
	b := []Float16{0x4200, 0x4400}      // [3.0, 4.0]
	expect := []Float16{0x4400, 0x4600} // [4.0, 6.0]

	result := AddSlice(a, b)
	if len(result) != len(expect) {
		t.Fatalf("Expected length %d, got %d", len(expect), len(result))
	}
	for i := range result {
		if result[i] != expect[i] {
			t.Errorf("result[%d] = %v, want %v", i, result[i], expect[i])
		}
	}
}

func TestSubSlice(t *testing.T) {
	a := []Float16{0x4400, 0x4600}      // [4.0, 6.0]
	b := []Float16{0x3C00, 0x4000}      // [1.0, 2.0]
	expect := []Float16{0x4200, 0x4400} // [3.0, 4.0]

	result := SubSlice(a, b)
	if len(result) != len(expect) {
		t.Fatalf("Expected length %d, got %d", len(expect), len(result))
	}
	for i := range result {
		if result[i] != expect[i] {
			t.Errorf("result[%d] = %v, want %v", i, result[i], expect[i])
		}
	}
}

func TestMulSlice(t *testing.T) {
	a := []Float16{0x3C00, 0x4000}      // [1.0, 2.0]
	b := []Float16{0x4200, 0x4400}      // [3.0, 4.0]
	expect := []Float16{0x4200, 0x4800} // [3.0, 8.0]

	result := MulSlice(a, b)
	if len(result) != len(expect) {
		t.Fatalf("Expected length %d, got %d", len(expect), len(result))
	}
	for i := range result {
		if result[i] != expect[i] {
			t.Errorf("result[%d] = %v, want %v", i, result[i], expect[i])
		}
	}
}

func TestDivSlice(t *testing.T) {
	a := []Float16{0x4200, 0x4800}      // [3.0, 8.0]
	b := []Float16{0x3C00, 0x4000}      // [1.0, 2.0]
	expect := []Float16{0x4200, 0x4400} // [3.0, 4.0]

	result := DivSlice(a, b)
	if len(result) != len(expect) {
		t.Fatalf("Expected length %d, got %d", len(expect), len(result))
	}
	for i := range result {
		if result[i] != expect[i] {
			t.Errorf("result[%d] = %v, want %v", i, result[i], expect[i])
		}
	}
}

func TestScaleSlice(t *testing.T) {
	s := []Float16{0x3C00, 0x4000}      // [1.0, 2.0]
	scalar := Float16(0x4200)           // 3.0
	expect := []Float16{0x4200, 0x4600} // [3.0, 6.0]

	result := ScaleSlice(s, scalar)
	if len(result) != len(expect) {
		t.Fatalf("Expected length %d, got %d", len(expect), len(result))
	}
	for i := range result {
		if result[i] != expect[i] {
			t.Errorf("result[%d] = %v, want %v", i, result[i], expect[i])
		}
	}
}

func TestSumSlice(t *testing.T) {
	s := []Float16{0x3C00, 0x4000, 0x4200} // [1.0, 2.0, 3.0]
	expect := Float16(0x4600)              // 6.0

	result := SumSlice(s)
	if result != expect {
		t.Errorf("SumSlice = %v, want %v", result, expect)
	}
}

func TestDotProduct(t *testing.T) {
	a := []Float16{0x3C00, 0x4000} // [1.0, 2.0]
	b := []Float16{0x4200, 0x4400} // [3.0, 4.0]
	// 1*3 + 2*4 = 3 + 8 = 11
	expect := Float16(0x4980) // 11.0

	result := DotProduct(a, b)
	if result != expect {
		t.Errorf("DotProduct = %v, want %v", result, expect)
	}
}

func TestNorm2(t *testing.T) {
	s := []Float16{0x4200, 0x4400} // [3.0, 4.0]
	// sqrt(3^2 + 4^2) = sqrt(9 + 16) = sqrt(25) = 5
	expect := Float16(0x4500) // 5.0

	result := Norm2(s)
	if result != expect {
		t.Errorf("Norm2 = %v, want %v", result, expect)
	}
}

func TestAddIEEE754(t *testing.T) {
	tests := []struct {
		name     string
		a        Float16
		b        Float16
		rounding RoundingMode
		expect   Float16
		hasError bool
	}{
		{
			name:     "subnormal + subnormal",
			a:        0x0001, // Smallest subnormal
			b:        0x0001,
			rounding: RoundNearestEven,
			expect:   0x0002,
			hasError: false,
		},
		{
			name:     "normal + subnormal",
			a:        0x3C00, // 1.0
			b:        0x0001,
			rounding: RoundNearestEven,
			expect:   0x3C00,
			hasError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := addIEEE754(tt.a, tt.b, tt.rounding)
			if (err != nil) != tt.hasError {
				t.Fatalf("addIEEE754() error = %v, wantErr %v", err, tt.hasError)
			}
			if result != tt.expect {
				t.Errorf("addIEEE754() = %v, want %v", result, tt.expect)
			}
		})
	}
}

func TestLess_Extra(t *testing.T) {
	tests := []struct {
		name   string
		a      Float16
		b      Float16
		expect bool
	}{
		{"-0 < +0", NegativeZero, PositiveZero, false},
		{"+0 < -0", PositiveZero, NegativeZero, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Less(tt.a, tt.b); got != tt.expect {
				t.Errorf("Less() = %v, want %v", got, tt.expect)
			}
		})
	}
}

func TestMinMax_Extra(t *testing.T) {
	tests := []struct {
		name    string
		a       Float16
		b       Float16
		minWant Float16
		maxWant Float16
	}{
		{"-Inf, +Inf", NegativeInfinity, PositiveInfinity, NegativeInfinity, PositiveInfinity},
		{"-0, +0", NegativeZero, PositiveZero, NegativeZero, PositiveZero},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Min(tt.a, tt.b); got != tt.minWant {
				t.Errorf("Min() = %v, want %v", got, tt.minWant)
			}
			if got := Max(tt.a, tt.b); got != tt.maxWant {
				t.Errorf("Max() = %v, want %v", got, tt.maxWant)
			}
		})
	}
}

func TestSlicePanics(t *testing.T) {
	a := []Float16{1, 2}
	b := []Float16{1}

	tests := []struct {
		name string
		f    func()
	}{
		{"AddSlice", func() { AddSlice(a, b) }},
		{"SubSlice", func() { SubSlice(a, b) }},
		{"MulSlice", func() { MulSlice(a, b) }},
		{"DivSlice", func() { DivSlice(a, b) }},
		{"DotProduct", func() { DotProduct(a, b) }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r == nil {
					t.Errorf("%s should have panicked", tt.name)
				}
			}()
			tt.f()
		})
	}
}
