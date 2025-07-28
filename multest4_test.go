package float16

import (
	"fmt"
	"testing"
)

func TestSpecificMultiplications(t *testing.T) {
	tests := []struct {
		a, b, want Float16
	}{
		{FromBits(0x3C00), FromBits(0x4400), FromBits(0x4400)}, // 1.0 * 4.0 = 4.0
		{FromBits(0x4000), FromBits(0x4400), FromBits(0x4800)}, // 2.0 * 4.0 = 8.0
		{FromBits(0x4400), FromBits(0x4400), FromBits(0x4C00)}, // 4.0 * 4.0 = 16.0 (0x4C00 is the correct bit pattern for 16.0)
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%v*%v", tt.a, tt.b), func(t *testing.T) {
			got, _ := mulIEEE754(tt.a, tt.b, RoundNearestEven)
			if got != tt.want {
				t.Errorf("mulIEEE754(%v, %v) = %v (0x%04X), want %v (0x%04X)",
					tt.a, tt.b, got, got, tt.want, tt.want)
			}
		})
	}
}
