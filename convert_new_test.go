package float16

import (
	"math"
	"testing"
)

func TestFromFloat32New(t *testing.T) {
	tests := []struct {
		name     string
		input    float32
		expected Float16
	}{
		{"zero", 0.0, 0x0000},
		{"-zero", float32(math.Copysign(0, -1)), 0x8000},
		{"one", 1.0, 0x3C00},
		{"-one", -1.0, 0xBC00},
		{"two", 2.0, 0x4000},
		{"half", 0.5, 0x3800},
		{"+inf", float32(math.Inf(1)), 0x7C00},
		{"-inf", float32(math.Inf(-1)), 0xFC00},
		{"max", 65504.0, 0x7BFF},
		{"-max", -65504.0, 0xFBFF},
		{"smallest normal", 6.103515625e-05, 0x0400},
		{"-smallest normal", -6.103515625e-05, 0x8400},
		{"smallest subnormal", 5.960464477539063e-08, 0x0001},
		{"-smallest subnormal", -5.960464477539063e-08, 0x8001},
		{"1.2", 1.2, 0x3CCD},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := fromFloat32New(tt.input)
			if result != tt.expected {
				t.Errorf("fromFloat32New(%g) = 0x%04x, expected 0x%04x",
					tt.input, result.Bits(), tt.expected.Bits())
			}
		})
	}
}
