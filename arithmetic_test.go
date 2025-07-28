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
		expect   Float16
		hasError bool
	}{
		{name: "1.0 - 0.5", a: 0x3C00, b: 0x3800, expect: 0x3800, hasError: false},  // 1.0 - 0.5 = 0.5
		{name: "0.5 - 0.25", a: 0x3800, b: 0x3400, expect: 0x3400, hasError: false}, // 0.5 - 0.25 = 0.25
		{name: "1.0 - 1.0", a: 0x3C00, b: 0x3C00, expect: 0x0000, hasError: false},  // 1.0 - 1.0 = 0.0
		{name: "1.0 - -1.0", a: 0x3C00, b: 0xBC00, expect: 0x4000, hasError: false}, // 1.0 - (-1.0) = 2.0
		{name: "-1.0 - 1.0", a: 0xBC00, b: 0x3C00, expect: 0xC000, hasError: false}, // -1.0 - 1.0 = -2.0
		{name: "0.0 - 0.0", a: 0x0000, b: 0x0000, expect: 0x0000, hasError: false},  // 0.0 - 0.0 = 0.0 (or -0.0 is also valid)
		{name: "Inf - 1.0", a: 0x7C00, b: 0x3C00, expect: 0x7C00, hasError: false},  // +Inf - 1.0 = +Inf
		{name: "1.0 - Inf", a: 0x3C00, b: 0x7C00, expect: 0xFC00, hasError: false},  // 1.0 - +Inf = -Inf
		{name: "NaN - 1.0", a: 0x7E00, b: 0x3C00, expect: 0x7E00, hasError: false},  // NaN - 1.0 = NaN
		{name: "1.0 - NaN", a: 0x3C00, b: 0x7E00, expect: 0x7E00, hasError: false},  // 1.0 - NaN = NaN
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := SubWithMode(tt.a, tt.b, ModeIEEEArithmetic, RoundNearestEven)

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
