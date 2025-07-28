package float16

import (
	"math"
	"testing"
)

func TestFloat32Conversion(t *testing.T) {
	// Test conversion of 4.0
	f32 := float32(4.0)
	f16 := FromFloat32(f32)
	backToF32 := f16.ToFloat32()

	t.Logf("float32(4.0) = %v, bits: %08b\n", f32, math.Float32bits(f32))
	t.Logf("FromFloat32(4.0) = %v, bits: %016b\n", f16, f16)
	t.Logf("back to float32 = %v, bits: %08b\n", backToF32, math.Float32bits(backToF32))

	// Verify the conversion is correct
	if backToF32 != f32 {
		t.Errorf("Conversion failed: expected %v, got %v", f32, backToF32)
	}
}
