package float16

import (
	"testing"
)

func TestBitPatterns(t *testing.T) {
	tests := []struct {
		name string
		bits uint16
	}{
		{"1.0", 0x3C00},
		{"2.0", 0x4000},
		{"4.0", 0x4400},
		{"8.0", 0x4800},
		{"16.0", 0x4C00},
		{"32.0", 0x5000},
		{"0x3136", 0x3136},
		{"0x3332", 0x3332},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := FromBits(tt.bits)
			_ = f // Use the value to prevent unused variable warning
		})
	}
}
