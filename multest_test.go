package float16

import "testing"

func TestMulIEEE754(t *testing.T) {
	tests := []struct {
		a, b Float16
		want Float16
	}{
		{FromBits(0x4400), FromBits(0x4400), FromBits(0x4C00)}, // 4.0 * 4.0 = 16.0 (0x4C00 is the correct bit pattern for 16.0)
	}

	for _, tt := range tests {
		got, _ := mulIEEE754(tt.a, tt.b, RoundNearestEven)
		if got != tt.want {
			t.Errorf("mulIEEE754(%v, %v) = %v, want %v", tt.a, tt.b, got, tt.want)
		}
	}
}
