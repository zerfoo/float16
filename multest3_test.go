package float16

import (
	"math"
	"testing"
)

func TestMulPipeline(t *testing.T) {
	// Test the full multiplication pipeline for 4.0 * 4.0
	a := FromFloat32(4.0)
	b := FromFloat32(4.0)
	
	t.Logf("a = %v, bits: %016b\n", a, a)
	t.Logf("b = %v, bits: %016b\n", b, b)
	
	// Convert to float32 and multiply
	a32 := a.ToFloat32()
	b32 := b.ToFloat32()
	result32 := a32 * b32
	
	t.Logf("a.ToFloat32() = %v, bits: %08b\n", a32, math.Float32bits(a32))
	t.Logf("b.ToFloat32() = %v, bits: %08b\n", b32, math.Float32bits(b32))
	t.Logf("a32 * b32 = %v, bits: %08b\n", result32, math.Float32bits(result32))
	
	// Convert back to Float16
	result := FromFloat32(result32)
	
	t.Logf("FromFloat32(result32) = %v, bits: %016b\n", result, result)
	
	// Compare with direct multiplication
	directResult, _ := mulIEEE754(a, b, RoundNearestEven)
	t.Logf("mulIEEE754(a, b) = %v, bits: %016b\n", directResult, directResult)
	
	// Expected result: 16.0 in Float16 (0x5000)
	expected := FromFloat32(16.0)
	t.Logf("Expected = %v, bits: %016b\n", expected, expected)
	
	if directResult != expected {
		t.Errorf("mulIEEE754(4.0, 4.0) = %v, want %v", directResult, expected)
	}
}
