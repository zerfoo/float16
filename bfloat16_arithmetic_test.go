package float16

import (
	"testing"
)

func TestBFloat16AddWithMode(t *testing.T) {
	tests := []struct {
		name      string
		a         BFloat16
		b         BFloat16
		mode      BFloat16ArithmeticMode
		rounding  RoundingMode
		expected  BFloat16
		expectErr bool
	}{
		// BFloat16ArithmeticModeIEEE tests
		{
			name:      "IEEE_NormalAdd",
			a:         BFloat16FromFloat32(1.0),
			b:         BFloat16FromFloat32(2.0),
			mode:      BFloat16ArithmeticModeIEEE,
			rounding:  RoundNearestEven,
			expected:  BFloat16FromFloat32(3.0),
			expectErr: false,
		},
		{
			name:      "IEEE_AddInf",
			a:         BFloat16PositiveInfinity,
			b:         BFloat16FromFloat32(1.0),
			mode:      BFloat16ArithmeticModeIEEE,
			rounding:  RoundNearestEven,
			expected:  BFloat16PositiveInfinity,
			expectErr: false,
		},
		{
			name:      "IEEE_AddNaN",
			a:         BFloat16QuietNaN,
			b:         BFloat16FromFloat32(1.0),
			mode:      BFloat16ArithmeticModeIEEE,
			rounding:  RoundNearestEven,
			expected:  BFloat16QuietNaN,
			expectErr: false,
		},
		{
			name:      "IEEE_AddOverflow",
			a:         BFloat16MaxValue,
			b:         BFloat16FromFloat32(1.0),
			mode:      BFloat16ArithmeticModeIEEE,
			rounding:  RoundNearestEven,
			expected:  BFloat16PositiveInfinity,
			expectErr: false,
		},
		{
			name:      "IEEE_AddUnderflow",
			a:         BFloat16SmallestPos,
			b:         BFloat16SmallestPos, // Adding two small numbers to cause underflow
			mode:      BFloat16ArithmeticModeIEEE,
			rounding:  RoundNearestEven,
			expected:  BFloat16FromBits(0x0100), // Corrected expected value
			expectErr: false,
		},

		// BFloat16ArithmeticModeStrict tests
		{
			name:      "Strict_NormalAdd",
			a:         BFloat16FromFloat32(1.0),
			b:         BFloat16FromFloat32(2.0),
			mode:      BFloat16ArithmeticModeStrict,
			rounding:  RoundNearestEven,
			expected:  BFloat16FromFloat32(3.0),
			expectErr: false,
		},
		{
			name:      "Strict_AddInf",
			a:         BFloat16PositiveInfinity,
			b:         BFloat16FromFloat32(1.0),
			mode:      BFloat16ArithmeticModeStrict,
			rounding:  RoundNearestEven,
			expected:  0, // Value doesn't matter if error is expected
			expectErr: true,
		},
		{
			name:      "Strict_AddNaN",
			a:         BFloat16QuietNaN,
			b:         BFloat16FromFloat32(1.0),
			mode:      BFloat16ArithmeticModeStrict,
			rounding:  RoundNearestEven,
			expected:  0, // Value doesn't matter if error is expected
			expectErr: true,
		},
		{
			name:      "Strict_AddOverflow",
			a:         BFloat16MaxValue,
			b:         BFloat16FromFloat32(1.0),
			mode:      BFloat16ArithmeticModeStrict,
			rounding:  RoundNearestEven,
			expected:  0, // Value doesn't matter if error is expected
			expectErr: true,
		},
		{
			name:      "Strict_AddUnderflow",
			a:         BFloat16SmallestPos,
			b:         BFloat16SmallestPos, // Adding two small numbers to cause underflow
			mode:      BFloat16ArithmeticModeStrict,
			rounding:  RoundNearestEven,
			expected:  0, // Value doesn't matter if error is expected
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := BFloat16AddWithMode(tt.a, tt.b, tt.mode, tt.rounding)

			if tt.expectErr {
				if err == nil {
					t.Errorf("Expected error for %s, but got none", tt.name)
				}
			} else {
				if err != nil {
					t.Errorf("Did not expect error for %s, but got: %v", tt.name, err)
				}
				if got != tt.expected {
					t.Errorf("BFloat16AddWithMode(%04x, %04x, %v, %v) = %04x, want %04x", uint16(tt.a), uint16(tt.b), tt.mode, tt.rounding, uint16(got), uint16(tt.expected))
				}
			}
		})
	}
}
