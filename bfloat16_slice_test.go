package float16

import (
	"testing"
)

func TestBFloat16AddSlice(t *testing.T) {
	a := []BFloat16{BFloat16FromFloat32(1.0), BFloat16FromFloat32(2.0), BFloat16FromFloat32(3.0)}
	b := []BFloat16{BFloat16FromFloat32(4.0), BFloat16FromFloat32(5.0), BFloat16FromFloat32(6.0)}

	result := BFloat16AddSlice(a, b)
	expected := []float32{5.0, 7.0, 9.0}

	for i, r := range result {
		got := r.ToFloat32()
		if got != expected[i] {
			t.Errorf("BFloat16AddSlice[%d] = %v, want %v", i, got, expected[i])
		}
	}
}

func TestBFloat16SubSlice(t *testing.T) {
	a := []BFloat16{BFloat16FromFloat32(10.0), BFloat16FromFloat32(5.0), BFloat16FromFloat32(3.0)}
	b := []BFloat16{BFloat16FromFloat32(4.0), BFloat16FromFloat32(2.0), BFloat16FromFloat32(1.0)}

	result := BFloat16SubSlice(a, b)
	expected := []float32{6.0, 3.0, 2.0}

	for i, r := range result {
		got := r.ToFloat32()
		if got != expected[i] {
			t.Errorf("BFloat16SubSlice[%d] = %v, want %v", i, got, expected[i])
		}
	}
}

func TestBFloat16MulSlice(t *testing.T) {
	a := []BFloat16{BFloat16FromFloat32(2.0), BFloat16FromFloat32(3.0), BFloat16FromFloat32(4.0)}
	b := []BFloat16{BFloat16FromFloat32(5.0), BFloat16FromFloat32(6.0), BFloat16FromFloat32(7.0)}

	result := BFloat16MulSlice(a, b)
	expected := []float32{10.0, 18.0, 28.0}

	for i, r := range result {
		got := r.ToFloat32()
		if got != expected[i] {
			t.Errorf("BFloat16MulSlice[%d] = %v, want %v", i, got, expected[i])
		}
	}
}

func TestBFloat16DivSlice(t *testing.T) {
	a := []BFloat16{BFloat16FromFloat32(10.0), BFloat16FromFloat32(9.0), BFloat16FromFloat32(8.0)}
	b := []BFloat16{BFloat16FromFloat32(2.0), BFloat16FromFloat32(3.0), BFloat16FromFloat32(4.0)}

	result := BFloat16DivSlice(a, b)
	expected := []float32{5.0, 3.0, 2.0}

	for i, r := range result {
		got := r.ToFloat32()
		if got != expected[i] {
			t.Errorf("BFloat16DivSlice[%d] = %v, want %v", i, got, expected[i])
		}
	}
}

func TestBFloat16ScaleSlice(t *testing.T) {
	s := []BFloat16{BFloat16FromFloat32(1.0), BFloat16FromFloat32(2.0), BFloat16FromFloat32(3.0)}
	scalar := BFloat16FromFloat32(2.0)

	result := BFloat16ScaleSlice(s, scalar)
	expected := []float32{2.0, 4.0, 6.0}

	for i, r := range result {
		got := r.ToFloat32()
		if got != expected[i] {
			t.Errorf("BFloat16ScaleSlice[%d] = %v, want %v", i, got, expected[i])
		}
	}
}

func TestBFloat16SumSlice(t *testing.T) {
	s := []BFloat16{BFloat16FromFloat32(1.0), BFloat16FromFloat32(2.0), BFloat16FromFloat32(3.0), BFloat16FromFloat32(4.0)}

	result := BFloat16SumSlice(s)
	got := result.ToFloat32()
	if got != 10.0 {
		t.Errorf("BFloat16SumSlice = %v, want 10.0", got)
	}
}

func TestBFloat16SumSliceEmpty(t *testing.T) {
	result := BFloat16SumSlice([]BFloat16{})
	if !result.IsZero() {
		t.Errorf("BFloat16SumSlice(empty) = %v, want zero", result.ToFloat32())
	}
}

func TestBFloat16SliceNegativeValues(t *testing.T) {
	a := []BFloat16{BFloat16FromFloat32(-1.0), BFloat16FromFloat32(-2.0)}
	b := []BFloat16{BFloat16FromFloat32(3.0), BFloat16FromFloat32(-4.0)}

	add := BFloat16AddSlice(a, b)
	if add[0].ToFloat32() != 2.0 {
		t.Errorf("AddSlice(-1+3) = %v, want 2.0", add[0].ToFloat32())
	}
	if add[1].ToFloat32() != -6.0 {
		t.Errorf("AddSlice(-2+-4) = %v, want -6.0", add[1].ToFloat32())
	}

	sub := BFloat16SubSlice(a, b)
	if sub[0].ToFloat32() != -4.0 {
		t.Errorf("SubSlice(-1-3) = %v, want -4.0", sub[0].ToFloat32())
	}
	if sub[1].ToFloat32() != 2.0 {
		t.Errorf("SubSlice(-2--4) = %v, want 2.0", sub[1].ToFloat32())
	}
}

func TestBFloat16SliceSpecialValues(t *testing.T) {
	inf := BFloat16PositiveInfinity
	nan := BFloat16QuietNaN
	zero := BFloat16PositiveZero
	one := BFloat16FromFloat32(1.0)

	// Inf + finite = Inf
	result := BFloat16AddSlice([]BFloat16{inf}, []BFloat16{one})
	if !result[0].IsInf(0) {
		t.Error("Inf + 1 should be Inf")
	}

	// Anything * 0 = 0
	result = BFloat16MulSlice([]BFloat16{one}, []BFloat16{zero})
	if !result[0].IsZero() {
		t.Error("1 * 0 should be 0")
	}

	// NaN + anything = NaN
	result = BFloat16AddSlice([]BFloat16{nan}, []BFloat16{one})
	if !result[0].IsNaN() {
		t.Error("NaN + 1 should be NaN")
	}

	// Div by zero = Inf
	result = BFloat16DivSlice([]BFloat16{one}, []BFloat16{zero})
	if !result[0].IsInf(0) {
		t.Error("1 / 0 should be Inf")
	}

	// Scale by zero
	result = BFloat16ScaleSlice([]BFloat16{one, BFloat16FromFloat32(2.0)}, zero)
	for i, r := range result {
		if !r.IsZero() {
			t.Errorf("ScaleSlice[%d] by zero = %v, want 0", i, r.ToFloat32())
		}
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

func TestBFloat16ToSlice32(t *testing.T) {
	input := []BFloat16{BFloat16FromFloat32(1.0), BFloat16FromFloat32(2.5), BFloat16FromFloat32(-3.0)}
	result := BFloat16ToSlice32(input)

	expected := []float32{1.0, 2.5, -3.0}
	for i, r := range result {
		if r != expected[i] {
			t.Errorf("BFloat16ToSlice32[%d] = %v, want %v", i, r, expected[i])
		}
	}
}

func TestBFloat16ToSlice32Empty(t *testing.T) {
	result := BFloat16ToSlice32([]BFloat16{})
	if len(result) != 0 {
		t.Errorf("BFloat16ToSlice32(empty) len = %d, want 0", len(result))
	}
}

func TestBFloat16FromSlice32(t *testing.T) {
	input := []float32{1.0, 2.5, -3.0}
	result := BFloat16FromSlice32(input)

	for i, r := range result {
		got := r.ToFloat32()
		if got != input[i] {
			t.Errorf("BFloat16FromSlice32[%d] = %v, want %v", i, got, input[i])
		}
	}
}

func TestBFloat16ToSlice64(t *testing.T) {
	input := []BFloat16{BFloat16FromFloat32(1.0), BFloat16FromFloat32(2.5), BFloat16FromFloat32(-3.0)}
	result := BFloat16ToSlice64(input)

	expected := []float64{1.0, 2.5, -3.0}
	for i, r := range result {
		if r != expected[i] {
			t.Errorf("BFloat16ToSlice64[%d] = %v, want %v", i, r, expected[i])
		}
	}
}

func TestBFloat16FromSlice64(t *testing.T) {
	input := []float64{1.0, 2.5, -3.0}
	result := BFloat16FromSlice64(input)

	for i, r := range result {
		got := r.ToFloat32()
		if got != float32(input[i]) {
			t.Errorf("BFloat16FromSlice64[%d] = %v, want %v", i, got, float32(input[i]))
		}
	}
}

func TestBFloat16SliceRoundTrip(t *testing.T) {
	original := []float32{0.0, 1.0, -1.0, 0.5, 128.0, -256.0}
	bf16 := BFloat16FromSlice32(original)
	roundtrip := BFloat16ToSlice32(bf16)

	for i, r := range roundtrip {
		if r != original[i] {
			t.Errorf("round-trip[%d] = %v, want %v", i, r, original[i])
		}
	}
}

func TestBFloat16SliceRoundTrip64(t *testing.T) {
	original := []float64{0.0, 1.0, -1.0, 0.5}
	bf16 := BFloat16FromSlice64(original)
	roundtrip := BFloat16ToSlice64(bf16)

	for i, r := range roundtrip {
		if r != original[i] {
			t.Errorf("round-trip64[%d] = %v, want %v", i, r, original[i])
		}
	}
}

func TestBFloat16ScaleSliceByOne(t *testing.T) {
	s := []BFloat16{BFloat16FromFloat32(1.0), BFloat16FromFloat32(2.0), BFloat16FromFloat32(3.0)}
	result := BFloat16ScaleSlice(s, BFloat16FromFloat32(1.0))

	for i, r := range result {
		if r != s[i] {
			t.Errorf("ScaleSlice by 1[%d] = %v, want %v", i, r.ToFloat32(), s[i].ToFloat32())
		}
	}
}

func TestBFloat16AddSliceIdentity(t *testing.T) {
	a := []BFloat16{BFloat16FromFloat32(1.0), BFloat16FromFloat32(2.0), BFloat16FromFloat32(3.0)}
	zeros := []BFloat16{BFloat16PositiveZero, BFloat16PositiveZero, BFloat16PositiveZero}

	result := BFloat16AddSlice(a, zeros)
	for i, r := range result {
		if r.ToFloat32() != a[i].ToFloat32() {
			t.Errorf("AddSlice identity[%d] = %v, want %v", i, r.ToFloat32(), a[i].ToFloat32())
		}
	}
}

// Benchmarks

func BenchmarkBFloat16SliceAdd(b *testing.B) {
	size := 1024
	a := make([]BFloat16, size)
	s := make([]BFloat16, size)
	for i := range a {
		a[i] = BFloat16FromFloat32(float32(i))
		s[i] = BFloat16FromFloat32(float32(i + 1))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		BFloat16AddSlice(a, s)
	}
}

func BenchmarkBFloat16SliceSub(b *testing.B) {
	size := 1024
	a := make([]BFloat16, size)
	s := make([]BFloat16, size)
	for i := range a {
		a[i] = BFloat16FromFloat32(float32(i))
		s[i] = BFloat16FromFloat32(float32(i + 1))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		BFloat16SubSlice(a, s)
	}
}

func BenchmarkBFloat16SliceMul(b *testing.B) {
	size := 1024
	a := make([]BFloat16, size)
	s := make([]BFloat16, size)
	for i := range a {
		a[i] = BFloat16FromFloat32(float32(i))
		s[i] = BFloat16FromFloat32(float32(i + 1))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		BFloat16MulSlice(a, s)
	}
}

func BenchmarkBFloat16SliceDiv(b *testing.B) {
	size := 1024
	a := make([]BFloat16, size)
	s := make([]BFloat16, size)
	for i := range a {
		a[i] = BFloat16FromFloat32(float32(i + 1))
		s[i] = BFloat16FromFloat32(float32(i + 1))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		BFloat16DivSlice(a, s)
	}
}

func BenchmarkBFloat16SliceScale(b *testing.B) {
	size := 1024
	s := make([]BFloat16, size)
	for i := range s {
		s[i] = BFloat16FromFloat32(float32(i))
	}
	scalar := BFloat16FromFloat32(2.0)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		BFloat16ScaleSlice(s, scalar)
	}
}

func BenchmarkBFloat16SliceSum(b *testing.B) {
	size := 1024
	s := make([]BFloat16, size)
	for i := range s {
		s[i] = BFloat16FromFloat32(float32(i) * 0.001) // small values to avoid overflow
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		BFloat16SumSlice(s)
	}
}

func BenchmarkBFloat16ToSlice32(b *testing.B) {
	size := 1024
	s := make([]BFloat16, size)
	for i := range s {
		s[i] = BFloat16FromFloat32(float32(i))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		BFloat16ToSlice32(s)
	}
}

func BenchmarkBFloat16FromSlice32(b *testing.B) {
	size := 1024
	s := make([]float32, size)
	for i := range s {
		s[i] = float32(i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		BFloat16FromSlice32(s)
	}
}

