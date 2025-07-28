package float16

import "testing"

func TestMul(t *testing.T) {
	tests := []struct {
		a, b, want Float16
	}{
		{FromBits(0x3C00), FromBits(0x4400), FromBits(0x4400)}, // 1.0 * 4.0 = 4.0
		{FromBits(0x4000), FromBits(0x4400), FromBits(0x4800)}, // 2.0 * 4.0 = 8.0
		{FromBits(0x4400), FromBits(0x4400), FromBits(0x4C00)}, // 4.0 * 4.0 = 16.0 (0x4C00 is the correct bit pattern for 16.0)
	}

	for _, tt := range tests {
		got := Mul(tt.a, tt.b)
		if got != tt.want {
			t.Errorf("Mul(%v, %v) = %v, want %v", tt.a, tt.b, got, tt.want)
		}
	}
}
