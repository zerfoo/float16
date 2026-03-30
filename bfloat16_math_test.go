package float16

import (
	"math"
	"testing"
)

func TestBFloat16Sqrt(t *testing.T) {
	tests := []struct {
		name string
		in   BFloat16
		want float32
	}{
		{"sqrt(4)=2", BFloat16FromFloat32(4), 2},
		{"sqrt(1)=1", BFloat16FromFloat32(1), 1},
		{"sqrt(0)=0", BFloat16PositiveZero, 0},
		{"sqrt(9)=3", BFloat16FromFloat32(9), 3},
		{"sqrt(0.25)=0.5", BFloat16FromFloat32(0.25), 0.5},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BFloat16Sqrt(tt.in).ToFloat32()
			if got != tt.want {
				t.Errorf("BFloat16Sqrt(%v) = %v, want %v", tt.in.ToFloat32(), got, tt.want)
			}
		})
	}
	// Special cases
	if !BFloat16Sqrt(BFloat16QuietNaN).IsNaN() {
		t.Error("Sqrt(NaN) should be NaN")
	}
	if !BFloat16Sqrt(BFloat16PositiveInfinity).IsInf(1) {
		t.Error("Sqrt(+Inf) should be +Inf")
	}
	if !BFloat16Sqrt(BFloat16FromFloat32(-1)).IsNaN() {
		t.Error("Sqrt(-1) should be NaN")
	}
	if !BFloat16Sqrt(BFloat16NegativeZero).IsZero() {
		t.Error("Sqrt(-0) should be zero")
	}
}

func TestBFloat16Exp(t *testing.T) {
	tests := []struct {
		name    string
		in      BFloat16
		wantF64 float64
		tol     float64
	}{
		{"exp(0)=1", BFloat16FromFloat32(0), 1, 0},
		{"exp(1)~=e", BFloat16FromFloat32(1), math.E, 0.05},
		{"exp(-1)~=1/e", BFloat16FromFloat32(-1), 1.0 / math.E, 0.01},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := float64(BFloat16Exp(tt.in).ToFloat32())
			if math.Abs(got-tt.wantF64) > tt.tol {
				t.Errorf("BFloat16Exp(%v) = %v, want ~%v (tol=%v)", tt.in.ToFloat32(), got, tt.wantF64, tt.tol)
			}
		})
	}
	if !BFloat16Exp(BFloat16QuietNaN).IsNaN() {
		t.Error("Exp(NaN) should be NaN")
	}
	if !BFloat16Exp(BFloat16PositiveInfinity).IsInf(1) {
		t.Error("Exp(+Inf) should be +Inf")
	}
	if !BFloat16Exp(BFloat16NegativeInfinity).IsZero() {
		t.Error("Exp(-Inf) should be 0")
	}
}

func TestBFloat16Log(t *testing.T) {
	tests := []struct {
		name    string
		in      BFloat16
		wantF64 float64
		tol     float64
	}{
		{"log(1)=0", BFloat16FromFloat32(1), 0, 0},
		{"log(e)~=1", BFloat16FromFloat32(float32(math.E)), 1, 0.02},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := float64(BFloat16Log(tt.in).ToFloat32())
			if math.Abs(got-tt.wantF64) > tt.tol {
				t.Errorf("BFloat16Log(%v) = %v, want ~%v", tt.in.ToFloat32(), got, tt.wantF64)
			}
		})
	}
	if !BFloat16Log(BFloat16PositiveZero).IsInf(-1) {
		t.Error("Log(0) should be -Inf")
	}
	if !BFloat16Log(BFloat16QuietNaN).IsNaN() {
		t.Error("Log(NaN) should be NaN")
	}
	if !BFloat16Log(BFloat16PositiveInfinity).IsInf(1) {
		t.Error("Log(+Inf) should be +Inf")
	}
	if !BFloat16Log(BFloat16FromFloat32(-1)).IsNaN() {
		t.Error("Log(-1) should be NaN")
	}
}

func TestBFloat16Log2(t *testing.T) {
	tests := []struct {
		name string
		in   BFloat16
		want float32
	}{
		{"log2(1)=0", BFloat16FromFloat32(1), 0},
		{"log2(2)=1", BFloat16FromFloat32(2), 1},
		{"log2(4)=2", BFloat16FromFloat32(4), 2},
		{"log2(8)=3", BFloat16FromFloat32(8), 3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BFloat16Log2(tt.in).ToFloat32()
			if got != tt.want {
				t.Errorf("BFloat16Log2(%v) = %v, want %v", tt.in.ToFloat32(), got, tt.want)
			}
		})
	}
	if !BFloat16Log2(BFloat16FromFloat32(-1)).IsNaN() {
		t.Error("Log2(-1) should be NaN")
	}
}

func TestBFloat16Sin(t *testing.T) {
	tests := []struct {
		name    string
		in      BFloat16
		wantF64 float64
		tol     float64
	}{
		{"sin(0)=0", BFloat16PositiveZero, 0, 0},
		{"sin(pi/2)~=1", BFloat16FromFloat32(float32(math.Pi / 2)), 1, 0.01},
		{"sin(pi)~=0", BFloat16FromFloat32(float32(math.Pi)), 0, 0.01},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := float64(BFloat16Sin(tt.in).ToFloat32())
			if math.Abs(got-tt.wantF64) > tt.tol {
				t.Errorf("BFloat16Sin(%v) = %v, want ~%v", tt.in.ToFloat32(), got, tt.wantF64)
			}
		})
	}
	if !BFloat16Sin(BFloat16QuietNaN).IsNaN() {
		t.Error("Sin(NaN) should be NaN")
	}
	if !BFloat16Sin(BFloat16PositiveInfinity).IsNaN() {
		t.Error("Sin(+Inf) should be NaN")
	}
}

func TestBFloat16Cos(t *testing.T) {
	tests := []struct {
		name    string
		in      BFloat16
		wantF64 float64
		tol     float64
	}{
		{"cos(0)=1", BFloat16PositiveZero, 1, 0},
		{"cos(pi)~=-1", BFloat16FromFloat32(float32(math.Pi)), -1, 0.02},
		{"cos(pi/2)~=0", BFloat16FromFloat32(float32(math.Pi / 2)), 0, 0.01},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := float64(BFloat16Cos(tt.in).ToFloat32())
			if math.Abs(got-tt.wantF64) > tt.tol {
				t.Errorf("BFloat16Cos(%v) = %v, want ~%v", tt.in.ToFloat32(), got, tt.wantF64)
			}
		})
	}
	if !BFloat16Cos(BFloat16QuietNaN).IsNaN() {
		t.Error("Cos(NaN) should be NaN")
	}
	if !BFloat16Cos(BFloat16PositiveInfinity).IsNaN() {
		t.Error("Cos(+Inf) should be NaN")
	}
}

func TestBFloat16Tanh(t *testing.T) {
	tests := []struct {
		name    string
		in      BFloat16
		wantF64 float64
		tol     float64
	}{
		{"tanh(0)=0", BFloat16PositiveZero, 0, 0},
		{"tanh(1)~=0.7616", BFloat16FromFloat32(1), 0.7616, 0.02},
		{"tanh(-1)~=-0.7616", BFloat16FromFloat32(-1), -0.7616, 0.02},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := float64(BFloat16Tanh(tt.in).ToFloat32())
			if math.Abs(got-tt.wantF64) > tt.tol {
				t.Errorf("BFloat16Tanh(%v) = %v, want ~%v", tt.in.ToFloat32(), got, tt.wantF64)
			}
		})
	}
	if got := BFloat16Tanh(BFloat16PositiveInfinity).ToFloat32(); got != 1 {
		t.Errorf("Tanh(+Inf) = %v, want 1", got)
	}
	if got := BFloat16Tanh(BFloat16NegativeInfinity).ToFloat32(); got != -1 {
		t.Errorf("Tanh(-Inf) = %v, want -1", got)
	}
	if !BFloat16Tanh(BFloat16QuietNaN).IsNaN() {
		t.Error("Tanh(NaN) should be NaN")
	}
}

func TestBFloat16Sigmoid(t *testing.T) {
	tests := []struct {
		name    string
		in      BFloat16
		wantF64 float64
		tol     float64
	}{
		{"sigmoid(0)~=0.5", BFloat16PositiveZero, 0.5, 0.01},
		{"sigmoid(large)~=1", BFloat16FromFloat32(10), 1, 0.001},
		{"sigmoid(-large)~=0", BFloat16FromFloat32(-10), 0, 0.001},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := float64(BFloat16Sigmoid(tt.in).ToFloat32())
			if math.Abs(got-tt.wantF64) > tt.tol {
				t.Errorf("BFloat16Sigmoid(%v) = %v, want ~%v", tt.in.ToFloat32(), got, tt.wantF64)
			}
		})
	}
	if got := BFloat16Sigmoid(BFloat16PositiveInfinity).ToFloat32(); got != 1 {
		t.Errorf("Sigmoid(+Inf) = %v, want 1", got)
	}
	if !BFloat16Sigmoid(BFloat16NegativeInfinity).IsZero() {
		t.Error("Sigmoid(-Inf) should be 0")
	}
	if !BFloat16Sigmoid(BFloat16QuietNaN).IsNaN() {
		t.Error("Sigmoid(NaN) should be NaN")
	}
}

func TestBFloat16FastSigmoid(t *testing.T) {
	// FastSigmoid should be within reasonable range of exact sigmoid
	vals := []float32{-5, -2, -1, -0.5, 0, 0.5, 1, 2, 5}
	for _, v := range vals {
		b := BFloat16FromFloat32(v)
		exact := float64(BFloat16Sigmoid(b).ToFloat32())
		fast := float64(BFloat16FastSigmoid(b).ToFloat32())
		// Allow up to 0.1 deviation for the fast approximation
		if math.Abs(exact-fast) > 0.1 {
			t.Errorf("FastSigmoid(%v) = %v, exact = %v, diff = %v", v, fast, exact, math.Abs(exact-fast))
		}
	}
	// Special cases
	if !BFloat16FastSigmoid(BFloat16QuietNaN).IsNaN() {
		t.Error("FastSigmoid(NaN) should be NaN")
	}
	if got := BFloat16FastSigmoid(BFloat16PositiveInfinity).ToFloat32(); got != 1 {
		t.Errorf("FastSigmoid(+Inf) = %v, want 1", got)
	}
}

func TestBFloat16FastTanh(t *testing.T) {
	// FastTanh should be within reasonable range of exact tanh
	vals := []float32{-5, -2, -1, -0.5, 0, 0.5, 1, 2, 5}
	for _, v := range vals {
		b := BFloat16FromFloat32(v)
		exact := float64(BFloat16Tanh(b).ToFloat32())
		fast := float64(BFloat16FastTanh(b).ToFloat32())
		if math.Abs(exact-fast) > 0.2 {
			t.Errorf("FastTanh(%v) = %v, exact = %v, diff = %v", v, fast, exact, math.Abs(exact-fast))
		}
	}
	// Special cases
	if !BFloat16FastTanh(BFloat16QuietNaN).IsNaN() {
		t.Error("FastTanh(NaN) should be NaN")
	}
	if got := BFloat16FastTanh(BFloat16PositiveInfinity).ToFloat32(); got != 1 {
		t.Errorf("FastTanh(+Inf) = %v, want 1", got)
	}
	if got := BFloat16FastTanh(BFloat16NegativeInfinity).ToFloat32(); got != -1 {
		t.Errorf("FastTanh(-Inf) = %v, want -1", got)
	}
}

// TestBFloat16MathFloat64Accuracy verifies all math functions match float64 within BFloat16 precision.
func TestBFloat16MathFloat64Accuracy(t *testing.T) {
	// BFloat16 has 7 mantissa bits, so relative error ~2^-7 ≈ 0.0078
	// We use a generous tolerance since we round-trip through BFloat16.
	const tol = 0.02

	testVals := []float32{0.1, 0.5, 1.0, 1.5, 2.0, 3.0, 4.0, 10.0}

	for _, v := range testVals {
		b := BFloat16FromFloat32(v)
		f64 := float64(v)

		// Sqrt
		gotSqrt := float64(BFloat16Sqrt(b).ToFloat32())
		wantSqrt := math.Sqrt(f64)
		if math.Abs(gotSqrt-wantSqrt)/wantSqrt > tol {
			t.Errorf("Sqrt(%v): got %v, want %v", v, gotSqrt, wantSqrt)
		}

		// Exp
		gotExp := float64(BFloat16Exp(b).ToFloat32())
		wantExp := math.Exp(f64)
		if math.Abs(gotExp-wantExp)/wantExp > tol {
			t.Errorf("Exp(%v): got %v, want %v", v, gotExp, wantExp)
		}

		// Log
		gotLog := float64(BFloat16Log(b).ToFloat32())
		wantLog := math.Log(f64)
		if wantLog != 0 && math.Abs(gotLog-wantLog)/math.Abs(wantLog) > tol {
			t.Errorf("Log(%v): got %v, want %v", v, gotLog, wantLog)
		}

		// Sin
		gotSin := float64(BFloat16Sin(b).ToFloat32())
		wantSin := math.Sin(f64)
		if math.Abs(wantSin) > 0.01 && math.Abs(gotSin-wantSin)/math.Abs(wantSin) > tol {
			t.Errorf("Sin(%v): got %v, want %v", v, gotSin, wantSin)
		}

		// Cos
		gotCos := float64(BFloat16Cos(b).ToFloat32())
		wantCos := math.Cos(f64)
		if math.Abs(wantCos) > 0.01 && math.Abs(gotCos-wantCos)/math.Abs(wantCos) > tol {
			t.Errorf("Cos(%v): got %v, want %v", v, gotCos, wantCos)
		}
	}
}
