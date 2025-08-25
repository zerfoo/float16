package float16

import (
	"math"
	"testing"
)

func TestBFloat16FromFloat32WithRounding(t *testing.T) {
	tests := []struct {
		name     string
		input    float32
		mode     RoundingMode
		expected BFloat16
	}{
		// RoundNearestEven tests
		{
			name:     "RNE_Exact",
			input:    1.0,
			mode:     RoundNearestEven,
			expected: BFloat16FromBits(0x3F80), // 1.0
		},
		{
			name:     "RNE_TieToEven_RoundDown",
			input:    float32(math.Float32frombits(0x3F808000)), // 1.0 + 0.5 ULP of BFloat16, LSB of BFloat16 is 0 (even)
			mode:     RoundNearestEven,
			expected: BFloat16FromBits(0x3F80), // Should round down to 1.0
		},
		{
			name:     "RNE_TieToEven_RoundUp",
			input:    float32(math.Float32frombits(0x3F818000)), // Value between 0x3F81 and 0x3F82, LSB of BFloat16 is 1 (odd)
			mode:     RoundNearestEven,
			expected: BFloat16FromBits(0x3F82), // Should round up
		},
		{
			name:     "RNE_RoundUp",
			input:    0.12345679,
			mode:     RoundNearestEven,
			expected: BFloat16FromBits(0x3DFD), // Corrected expected value
		},
		{
			name:     "RNE_RoundDown",
			input:    0.9876543,
			mode:     RoundNearestEven,
			expected: BFloat16FromBits(0x3F7D), // Corrected expected value
		},

		// RoundTowardZero tests
		{
			name:     "RTZ_Positive",
			input:    1.2345679,
			mode:     RoundTowardZero,
			expected: BFloat16FromBits(0x3F9E), // Truncated value
		},
		{
			name:     "RTZ_Negative",
			input:    -1.2345679,
			mode:     RoundTowardZero,
			expected: BFloat16FromBits(0xBF9E), // Truncated value
		},

		// RoundTowardPositive tests
		{
			name:     "RTP_Positive_RoundUp",
			input:    1.00000001,
			mode:     RoundTowardPositive,
			expected: BFloat16FromBits(0x3F80), // Should round up to 1.0
		},
		{
			name:     "RTP_Positive_NoRound",
			input:    1.0,
			mode:     RoundTowardPositive,
			expected: BFloat16FromBits(0x3F80), // No rounding needed
		},
		{
			name:     "RTP_Negative_NoRound",
			input:    -1.00000001,
			mode:     RoundTowardPositive,
			expected: BFloat16FromBits(0xBF80), // Should round towards zero (up)
		},

		// RoundTowardNegative tests
		{
			name:     "RTN_Positive_NoRound",
			input:    1.00000001,
			mode:     RoundTowardNegative,
			expected: BFloat16FromBits(0x3F80), // Should round towards zero (down)
		},
		{
			name:     "RTN_Negative_RoundDown",
			input:    float32(math.Float32frombits(0xBF800001)), // -1.00000001
			mode:     RoundTowardNegative,
			expected: BFloat16FromBits(0xBF81), // Should round down
		},
		{
			name:     "RTN_Negative_NoRound",
			input:    -1.0,
			mode:     RoundTowardNegative,
			expected: BFloat16FromBits(0xBF80), // No rounding needed
		},

		// RoundNearestAway tests
		{
			name:     "RNA_Positive_TieAway",
			input:    float32(math.Float32frombits(0x3F808000)), // 1.0 + 0.5 ULP of BFloat16
			mode:     RoundNearestAway,
			expected: BFloat16FromBits(0x3F81), // Should round up
		},
		{
			name:     "RNA_Negative_TieAway",
			input:    float32(math.Float32frombits(0xBF808000)), // -1.0 - 0.5 ULP of BFloat16
			mode:     RoundNearestAway,
			expected: BFloat16FromBits(0xBF81), // Should round down (away from zero)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BFloat16FromFloat32WithRounding(tt.input, tt.mode)
			if got != tt.expected {
				t.Errorf("BFloat16FromFloat32WithRounding(%f, %v) = %04x, want %04x", tt.input, tt.mode, uint16(got), uint16(tt.expected))
			}
		})
	}
}

func TestBFloat16FromFloat64WithRounding(t *testing.T) {
	tests := []struct {
		name     string
		input    float64
		mode     RoundingMode
		expected BFloat16
	}{
		{
			name:     "RNE_Exact",
			input:    1.0,
			mode:     RoundNearestEven,
			expected: BFloat16FromBits(0x3F80), // 1.0
		},
		{
			name:     "RNE_TieToEven_RoundDown",
			input:    float64(math.Float32frombits(0x3F808000)), // 1.0 + 0.5 ULP of BFloat16, LSB of BFloat16 is 0 (even)
			mode:     RoundNearestEven,
			expected: BFloat16FromBits(0x3F80), // Should round down to 1.0
		},
		{
			name:     "RNE_TieToEven_RoundUp",
			input:    float64(math.Float32frombits(0x3F818000)), // Value between 0x3F81 and 0x3F82, LSB of BFloat16 is 1 (odd)
			mode:     RoundNearestEven,
			expected: BFloat16FromBits(0x3F82), // Should round up
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BFloat16FromFloat64WithRounding(tt.input, tt.mode)
			if got != tt.expected {
				t.Errorf("BFloat16FromFloat64WithRounding(%f, %v) = %04x, want %04x", tt.input, tt.mode, uint16(got), uint16(tt.expected))
			}
		})
	}
}

func TestBFloat16FromFloat32WithMode(t *testing.T) {
	tests := []struct {
		name      string
		input     float32
		convMode  ConversionMode
		roundMode RoundingMode
		expected  BFloat16
		expectErr bool
	}{
		// ModeIEEE tests
		{
			name:      "IEEE_Normal",
			input:     1.0,
			convMode:  ModeIEEE,
			roundMode: RoundNearestEven,
			expected:  BFloat16FromBits(0x3F80),
			expectErr: false,
		},
		{
			name:      "IEEE_Overflow",
			input:     math.MaxFloat32, // A large float32 that overflows BFloat16
			convMode:  ModeIEEE,
			roundMode: RoundNearestEven,
			expected:  BFloat16PositiveInfinity,
			expectErr: false,
		},
		{
			name:      "IEEE_Underflow",
			input:     math.SmallestNonzeroFloat32, // A small float32 that underflows BFloat16
			convMode:  ModeIEEE,
			roundMode: RoundNearestEven,
			expected:  BFloat16PositiveZero,
			expectErr: false,
		},
		{
			name:      "IEEE_NaN",
			input:     float32(math.NaN()),
			convMode:  ModeIEEE,
			roundMode: RoundNearestEven,
			expected:  BFloat16QuietNaN,
			expectErr: false,
		},
		{
			name:      "IEEE_Inf",
			input:     float32(math.Inf(1)),
			convMode:  ModeIEEE,
			roundMode: RoundNearestEven,
			expected:  BFloat16PositiveInfinity,
			expectErr: false,
		},

		// ModeStrict tests
		{
			name:      "Strict_Normal",
			input:     1.0,
			convMode:  ModeStrict,
			roundMode: RoundNearestEven,
			expected:  BFloat16FromBits(0x3F80),
			expectErr: false,
		},
		{
			name:      "Strict_Overflow",
			input:     math.MaxFloat32,
			convMode:  ModeStrict,
			roundMode: RoundNearestEven,
			expected:  0, // Value doesn't matter if error is expected
			expectErr: true,
		},
		{
			name:      "Strict_Underflow",
			input:     math.SmallestNonzeroFloat32,
			convMode:  ModeStrict,
			roundMode: RoundNearestEven,
			expected:  0, // Value doesn't matter if error is expected
			expectErr: true,
		},
		{
			name:      "Strict_NaN",
			input:     float32(math.NaN()),
			convMode:  ModeStrict,
			roundMode: RoundNearestEven,
			expected:  0, // Value doesn't matter if error is expected
			expectErr: true,
		},
		{
			name:      "Strict_Inf",
			input:     float32(math.Inf(1)),
			convMode:  ModeStrict,
			roundMode: RoundNearestEven,
			expected:  0, // Value doesn't matter if error is expected
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := BFloat16FromFloat32WithMode(tt.input, tt.convMode, tt.roundMode)

			if tt.expectErr {
				if err == nil {
					t.Errorf("Expected error for %s, but got none", tt.name)
				}
			} else {
				if err != nil {
					t.Errorf("Did not expect error for %s, but got: %v", tt.name, err)
				}
				if got != tt.expected {
					t.Errorf("BFloat16FromFloat32WithMode(%f, %v, %v) = %04x, want %04x", tt.input, tt.convMode, tt.roundMode, uint16(got), uint16(tt.expected))
				}
			}
		})
	}
}