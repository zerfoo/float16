package float16

import (
	"math"
	"testing"
)

func TestBFloat16AddSlice(t *testing.T) {
	a := []BFloat16{BFloat16FromFloat32(1.0), BFloat16FromFloat32(2.0), BFloat16FromFloat32(3.0)}
	b := []BFloat16{BFloat16FromFloat32(4.0), BFloat16FromFloat32(5.0), BFloat16FromFloat32(6.0)}
	result := BFloat16AddSlice(a, b)
	expected := []float32{5.0, 7.0, 9.0}
	for i, v := range result {
		if v.ToFloat32() != expected[i] {
			t.Errorf("BFloat16AddSlice[%d] = %v, want %v", i, v.ToFloat32(), expected[i])
		}
	}
}

func TestBFloat16SubSlice(t *testing.T) {
	a := []BFloat16{BFloat16FromFloat32(5.0), BFloat16FromFloat32(7.0), BFloat16FromFloat32(9.0)}
	b := []BFloat16{BFloat16FromFloat32(1.0), BFloat16FromFloat32(2.0), BFloat16FromFloat32(3.0)}
	result := BFloat16SubSlice(a, b)
	expected := []float32{4.0, 5.0, 6.0}
	for i, v := range result {
		if v.ToFloat32() != expected[i] {
			t.Errorf("BFloat16SubSlice[%d] = %v, want %v", i, v.ToFloat32(), expected[i])
		}
	}
}

func TestBFloat16MulSlice(t *testing.T) {
	a := []BFloat16{BFloat16FromFloat32(2.0), BFloat16FromFloat32(3.0), BFloat16FromFloat32(4.0)}
	b := []BFloat16{BFloat16FromFloat32(5.0), BFloat16FromFloat32(6.0), BFloat16FromFloat32(7.0)}
	result := BFloat16MulSlice(a, b)
	expected := []float32{10.0, 18.0, 28.0}
	for i, v := range result {
		if v.ToFloat32() != expected[i] {
			t.Errorf("BFloat16MulSlice[%d] = %v, want %v", i, v.ToFloat32(), expected[i])
		}
	}
}

func TestBFloat16DivSlice(t *testing.T) {
	a := []BFloat16{BFloat16FromFloat32(10.0), BFloat16FromFloat32(18.0), BFloat16FromFloat32(28.0)}
	b := []BFloat16{BFloat16FromFloat32(2.0), BFloat16FromFloat32(3.0), BFloat16FromFloat32(4.0)}
	result := BFloat16DivSlice(a, b)
	expected := []float32{5.0, 6.0, 7.0}
	for i, v := range result {
		if v.ToFloat32() != expected[i] {
			t.Errorf("BFloat16DivSlice[%d] = %v, want %v", i, v.ToFloat32(), expected[i])
		}
	}
}

func TestBFloat16ScaleSlice(t *testing.T) {
	s := []BFloat16{BFloat16FromFloat32(1.0), BFloat16FromFloat32(2.0), BFloat16FromFloat32(3.0)}
	scalar := BFloat16FromFloat32(10.0)
	result := BFloat16ScaleSlice(s, scalar)
	expected := []float32{10.0, 20.0, 30.0}
	for i, v := range result {
		if v.ToFloat32() != expected[i] {
			t.Errorf("BFloat16ScaleSlice[%d] = %v, want %v", i, v.ToFloat32(), expected[i])
		}
	}
}

func TestBFloat16SumSlice(t *testing.T) {
	s := []BFloat16{BFloat16FromFloat32(1.0), BFloat16FromFloat32(2.0), BFloat16FromFloat32(3.0), BFloat16FromFloat32(4.0)}
	result := BFloat16SumSlice(s)
	if result.ToFloat32() != 10.0 {
		t.Errorf("BFloat16SumSlice = %v, want 10.0", result.ToFloat32())
	}
}

func TestBFloat16SumSliceEmpty(t *testing.T) {
	result := BFloat16SumSlice(nil)
	if !result.IsZero() {
		t.Errorf("BFloat16SumSlice(nil) = %v, want 0", result.ToFloat32())
	}
}

func TestBFloat16DotProduct(t *testing.T) {
	a := []BFloat16{BFloat16FromFloat32(1.0), BFloat16FromFloat32(2.0), BFloat16FromFloat32(3.0)}
	b := []BFloat16{BFloat16FromFloat32(4.0), BFloat16FromFloat32(5.0), BFloat16FromFloat32(6.0)}
	result := BFloat16DotProduct(a, b)
	// 1*4 + 2*5 + 3*6 = 4 + 10 + 18 = 32
	if result.ToFloat32() != 32.0 {
		t.Errorf("BFloat16DotProduct = %v, want 32.0", result.ToFloat32())
	}
}

func TestBFloat16Norm2(t *testing.T) {
	s := []BFloat16{BFloat16FromFloat32(3.0), BFloat16FromFloat32(4.0)}
	result := BFloat16Norm2(s)
	// sqrt(9 + 16) = 5
	if result.ToFloat32() != 5.0 {
		t.Errorf("BFloat16Norm2 = %v, want 5.0", result.ToFloat32())
	}
}

func TestBFloat16SlicePanics(t *testing.T) {
	a := []BFloat16{BFloat16FromFloat32(1.0)}
	b := []BFloat16{BFloat16FromFloat32(1.0), BFloat16FromFloat32(2.0)}

	tests := []struct {
		name string
		fn   func()
	}{
		{"AddSlice", func() { BFloat16AddSlice(a, b) }},
		{"SubSlice", func() { BFloat16SubSlice(a, b) }},
		{"MulSlice", func() { BFloat16MulSlice(a, b) }},
		{"DivSlice", func() { BFloat16DivSlice(a, b) }},
		{"DotProduct", func() { BFloat16DotProduct(a, b) }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r == nil {
					t.Errorf("%s did not panic on mismatched lengths", tt.name)
				}
			}()
			tt.fn()
		})
	}
}

func TestToBFloat16Slice(t *testing.T) {
	input := []float32{1.0, 2.5, -3.0, 0.0}
	result := ToBFloat16Slice(input)
	for i, v := range result {
		if v.ToFloat32() != input[i] {
			t.Errorf("ToBFloat16Slice[%d] = %v, want %v", i, v.ToFloat32(), input[i])
		}
	}
}

func TestToBFloat16SliceEmpty(t *testing.T) {
	result := ToBFloat16Slice(nil)
	if len(result) != 0 {
		t.Errorf("ToBFloat16Slice(nil) len = %d, want 0", len(result))
	}
}

func TestToBFloat16SliceWithMode(t *testing.T) {
	t.Run("IEEE mode", func(t *testing.T) {
		input := []float32{1.0, 2.0, float32(math.Inf(1))}
		result, errs := ToBFloat16SliceWithMode(input, ModeIEEE, RoundNearestEven)
		if len(result) != 3 {
			t.Fatalf("expected 3 results, got %d", len(result))
		}
		for _, e := range errs {
			if e != nil {
				t.Errorf("unexpected error in IEEE mode: %v", e)
			}
		}
		if result[0].ToFloat32() != 1.0 {
			t.Errorf("result[0] = %v, want 1.0", result[0].ToFloat32())
		}
	})

	t.Run("Strict mode overflow", func(t *testing.T) {
		input := []float32{1.0, float32(math.Inf(1))}
		_, errs := ToBFloat16SliceWithMode(input, ModeStrict, RoundNearestEven)
		if errs[0] != nil {
			t.Errorf("expected no error for 1.0, got %v", errs[0])
		}
		if errs[1] == nil {
			t.Error("expected error for Inf in strict mode")
		}
	})
}

func TestBFloat16ToSlice32(t *testing.T) {
	input := []BFloat16{BFloat16FromFloat32(1.0), BFloat16FromFloat32(2.5), BFloat16FromFloat32(-3.0)}
	result := BFloat16ToSlice32(input)
	expected := []float32{1.0, 2.5, -3.0}
	for i, v := range result {
		if v != expected[i] {
			t.Errorf("BFloat16ToSlice32[%d] = %v, want %v", i, v, expected[i])
		}
	}
}

func TestBFloat16ToSlice32Empty(t *testing.T) {
	result := BFloat16ToSlice32(nil)
	if len(result) != 0 {
		t.Errorf("BFloat16ToSlice32(nil) len = %d, want 0", len(result))
	}
}

func TestBFloat16ToSlice64(t *testing.T) {
	input := []BFloat16{BFloat16FromFloat32(1.0), BFloat16FromFloat32(2.0)}
	result := BFloat16ToSlice64(input)
	if result[0] != 1.0 || result[1] != 2.0 {
		t.Errorf("BFloat16ToSlice64 = %v, want [1.0, 2.0]", result)
	}
}

func TestBFloat16FromSlice64(t *testing.T) {
	input := []float64{1.0, 2.5, -3.0}
	result := BFloat16FromSlice64(input)
	for i, v := range result {
		if v.ToFloat32() != float32(input[i]) {
			t.Errorf("BFloat16FromSlice64[%d] = %v, want %v", i, v.ToFloat32(), float32(input[i]))
		}
	}
}

func TestBFloat16RoundTripSlice(t *testing.T) {
	input := []float32{0.0, 1.0, -1.0, 0.5, 100.0}
	bf16 := ToBFloat16Slice(input)
	f32 := BFloat16ToSlice32(bf16)
	for i := range input {
		if f32[i] != input[i] {
			t.Errorf("round trip [%d] = %v, want %v", i, f32[i], input[i])
		}
	}
}

// Benchmarks

func BenchmarkBFloat16AddSlice(b *testing.B) {
	n := 1024
	a := make([]BFloat16, n)
	s := make([]BFloat16, n)
	for i := range a {
		a[i] = BFloat16FromFloat32(float32(i))
		s[i] = BFloat16FromFloat32(float32(i + 1))
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		BFloat16AddSlice(a, s)
	}
}

func BenchmarkBFloat16MulSlice(b *testing.B) {
	n := 1024
	a := make([]BFloat16, n)
	s := make([]BFloat16, n)
	for i := range a {
		a[i] = BFloat16FromFloat32(float32(i))
		s[i] = BFloat16FromFloat32(float32(i + 1))
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		BFloat16MulSlice(a, s)
	}
}

func BenchmarkBFloat16DotProduct(b *testing.B) {
	n := 1024
	a := make([]BFloat16, n)
	s := make([]BFloat16, n)
	for i := range a {
		a[i] = BFloat16FromFloat32(float32(i) * 0.01)
		s[i] = BFloat16FromFloat32(float32(i) * 0.01)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		BFloat16DotProduct(a, s)
	}
}

func BenchmarkToBFloat16Slice(b *testing.B) {
	n := 1024
	s := make([]float32, n)
	for i := range s {
		s[i] = float32(i) * 0.1
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ToBFloat16Slice(s)
	}
}

func BenchmarkBFloat16ToSlice32(b *testing.B) {
	n := 1024
	s := make([]BFloat16, n)
	for i := range s {
		s[i] = BFloat16FromFloat32(float32(i) * 0.1)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		BFloat16ToSlice32(s)
	}
}

func BenchmarkBFloat16ScaleSlice(b *testing.B) {
	n := 1024
	s := make([]BFloat16, n)
	for i := range s {
		s[i] = BFloat16FromFloat32(float32(i))
	}
	scalar := BFloat16FromFloat32(2.0)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		BFloat16ScaleSlice(s, scalar)
	}
}
