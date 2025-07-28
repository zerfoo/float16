package float16

import (
	"math"
	"testing"
)

func TestFloat32Multiplication(t *testing.T) {
	// Test 4.0 * 4.0 in float32
	a := float32(4.0)
	b := float32(4.0)
	result := a * b
	
	t.Logf("a = %v, bits: %08b\n", a, math.Float32bits(a))
	t.Logf("b = %v, bits: %08b\n", b, math.Float32bits(b))
	t.Logf("a * b = %v, bits: %08b\n", result, math.Float32bits(result))
	
	// Verify the result is 16.0
	expected := float32(16.0)
	if result != expected {
		t.Errorf("Multiplication failed: expected %v, got %v", expected, result)
	}
}
