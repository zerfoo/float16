package float16

import (
	"math"
	"testing"
)

func TestFromFloat64WithMode_Extra(t *testing.T) {
	tests := []struct {
		name      string
		input     float64
		convMode  ConversionMode
		roundMode RoundingMode
		expected  Float16
		hasError  bool
		errCode   ErrorCode
	}{
		{
			name:      "Strict mode overflow",
			input:     70000.0,
			convMode:  ModeStrict,
			roundMode: RoundNearestEven,
			hasError:  true,
			errCode:   ErrOverflow,
		},
		{
			name:      "Strict mode underflow",
			input:     1e-6,
			convMode:  ModeStrict,
			roundMode: RoundNearestEven,
			hasError:  true,
			errCode:   ErrUnderflow,
		},
		{
			name:      "Strict mode NaN",
			input:     math.NaN(),
			convMode:  ModeStrict,
			roundMode: RoundNearestEven,
			hasError:  true,
			errCode:   ErrNaN,
		},
		{
			name:      "Strict mode Inf",
			input:     math.Inf(1),
			convMode:  ModeStrict,
			roundMode: RoundNearestEven,
			hasError:  true,
			errCode:   ErrInfinity,
		},
		{
			name:      "IEEE mode overflow",
			input:     70000.0,
			convMode:  ModeIEEE,
			roundMode: RoundNearestEven,
			expected:  PositiveInfinity,
		},
		{
			name:      "IEEE mode underflow",
			input:     1e-6,
			convMode:  ModeIEEE,
			roundMode: RoundNearestEven,
			expected:  0x0011,
		},
		{
			name:      "IEEE mode out of range",
			input:     80000.0,
			convMode:  ModeIEEE,
			roundMode: RoundNearestEven,
			expected:  PositiveInfinity,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := FromFloat64WithMode(tt.input, tt.convMode, tt.roundMode)

			if tt.hasError {
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
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if result != tt.expected {
				t.Errorf("FromFloat64WithMode() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestShouldRound(t *testing.T) {
	tests := []struct {
		name        string
		mantissa    uint32
		shift       int
		mode        RoundingMode
		sign        uint16
		shouldRound bool
	}{
		{"NearestEven, no round", 0b1010_0000, 4, RoundNearestEven, 0, false},
		{"NearestEven, round up", 0b1011_1000, 4, RoundNearestEven, 0, true},
		{"NearestEven, halfway, even", 0b1010_1000, 4, RoundNearestEven, 0, false},
		{"NearestEven, halfway, odd", 0b1011_1000, 4, RoundNearestEven, 0, true},
		{"NearestAway, round up", 0b1010_1000, 4, RoundNearestAway, 0, true},
		{"TowardZero, no round", 0b1010_1111, 4, RoundTowardZero, 0, false},
		{"TowardPositive, round up", 0b1010_0001, 4, RoundTowardPositive, 0, true},
		{"TowardNegative, no round", 0b1010_0001, 4, RoundTowardNegative, 0, false},
		{"TowardPositive, no round (negative)", 0b1010_0001, 4, RoundTowardPositive, 0x8000, false},
		{"TowardNegative, round up (negative)", 0b1010_0001, 4, RoundTowardNegative, 0x8000, true},

		// Default case
		{"Invalid rounding mode", 0b1010_0001, 4, 99, 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shouldRound(tt.mantissa, tt.shift, tt.sign)
			if got != tt.shouldRound {
				t.Errorf("shouldRound(%d, %d, %d) = %v, want %v", tt.mantissa, tt.shift, tt.sign, got, tt.shouldRound)
			}
		})
	}
}

func TestParse(t *testing.T) {
	_, err := Parse("1.0")
	if err == nil {
		t.Error("Expected error, got nil")
	}
}
