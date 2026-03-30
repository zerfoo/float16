package float16

import (
	"math"
	"testing"
)

func TestBFloat16AddWithMode(t *testing.T) {
	tests := []struct {
		name     string
		a, b     BFloat16
		mode     ArithmeticMode
		rounding RoundingMode
		want     BFloat16
		wantErr  bool
		errCode  ErrorCode
	}{
		{"1+1 IEEE", BFloat16FromFloat32(1), BFloat16FromFloat32(1), ModeIEEEArithmetic, RoundNearestEven, BFloat16FromFloat32(2), false, 0},
		{"1+1 fast", BFloat16FromFloat32(1), BFloat16FromFloat32(1), ModeFastArithmetic, RoundNearestEven, BFloat16FromFloat32(2), false, 0},
		{"a+0", BFloat16FromFloat32(3), BFloat16PositiveZero, ModeIEEEArithmetic, RoundNearestEven, BFloat16FromFloat32(3), false, 0},
		{"0+b", BFloat16PositiveZero, BFloat16FromFloat32(5), ModeIEEEArithmetic, RoundNearestEven, BFloat16FromFloat32(5), false, 0},
		{"+inf+finite", BFloat16PositiveInfinity, BFloat16FromFloat32(1), ModeIEEEArithmetic, RoundNearestEven, BFloat16PositiveInfinity, false, 0},
		{"+inf+-inf IEEE", BFloat16PositiveInfinity, BFloat16NegativeInfinity, ModeIEEEArithmetic, RoundNearestEven, BFloat16QuietNaN, false, 0},
		{"+inf+-inf exact", BFloat16PositiveInfinity, BFloat16NegativeInfinity, ModeExactArithmetic, RoundNearestEven, 0, true, ErrInvalidOperation},
		{"NaN+1 IEEE", BFloat16QuietNaN, BFloat16FromFloat32(1), ModeIEEEArithmetic, RoundNearestEven, BFloat16QuietNaN, false, 0},
		{"NaN+1 exact", BFloat16QuietNaN, BFloat16FromFloat32(1), ModeExactArithmetic, RoundNearestEven, 0, true, ErrNaN},
		{"1+NaN exact", BFloat16FromFloat32(1), BFloat16QuietNaN, ModeExactArithmetic, RoundNearestEven, 0, true, ErrNaN},
		{"negative add", BFloat16FromFloat32(-2.5), BFloat16FromFloat32(1.5), ModeIEEEArithmetic, RoundNearestEven, BFloat16FromFloat32(-1), false, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := BFloat16AddWithMode(tt.a, tt.b, tt.mode, tt.rounding)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				fe, ok := err.(*Float16Error)
				if !ok {
					t.Fatalf("expected *Float16Error, got %T", err)
				}
				if fe.Code != tt.errCode {
					t.Errorf("error code = %d, want %d", fe.Code, tt.errCode)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tt.want.IsNaN() {
				if !got.IsNaN() {
					t.Errorf("got 0x%04X, want NaN", got.Bits())
				}
				return
			}
			if got != tt.want {
				t.Errorf("got 0x%04X (%v), want 0x%04X (%v)", got.Bits(), got, tt.want.Bits(), tt.want)
			}
		})
	}
}

func TestBFloat16SubWithMode(t *testing.T) {
	tests := []struct {
		name     string
		a, b     BFloat16
		mode     ArithmeticMode
		rounding RoundingMode
		want     BFloat16
		wantErr  bool
		errCode  ErrorCode
	}{
		{"3-1", BFloat16FromFloat32(3), BFloat16FromFloat32(1), ModeIEEEArithmetic, RoundNearestEven, BFloat16FromFloat32(2), false, 0},
		{"1-1", BFloat16FromFloat32(1), BFloat16FromFloat32(1), ModeIEEEArithmetic, RoundNearestEven, BFloat16PositiveZero, false, 0},
		{"inf-inf exact", BFloat16PositiveInfinity, BFloat16PositiveInfinity, ModeExactArithmetic, RoundNearestEven, 0, true, ErrInvalidOperation},
		{"NaN-1 exact", BFloat16QuietNaN, BFloat16FromFloat32(1), ModeExactArithmetic, RoundNearestEven, 0, true, ErrNaN},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := BFloat16SubWithMode(tt.a, tt.b, tt.mode, tt.rounding)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				fe, ok := err.(*Float16Error)
				if !ok {
					t.Fatalf("expected *Float16Error, got %T", err)
				}
				if fe.Code != tt.errCode {
					t.Errorf("error code = %d, want %d", fe.Code, tt.errCode)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("got 0x%04X (%v), want 0x%04X (%v)", got.Bits(), got, tt.want.Bits(), tt.want)
			}
		})
	}
}

func TestBFloat16MulWithMode(t *testing.T) {
	tests := []struct {
		name     string
		a, b     BFloat16
		mode     ArithmeticMode
		rounding RoundingMode
		want     BFloat16
		wantErr  bool
		errCode  ErrorCode
	}{
		{"2*3", BFloat16FromFloat32(2), BFloat16FromFloat32(3), ModeIEEEArithmetic, RoundNearestEven, BFloat16FromFloat32(6), false, 0},
		{"2*3 fast", BFloat16FromFloat32(2), BFloat16FromFloat32(3), ModeFastArithmetic, RoundNearestEven, BFloat16FromFloat32(6), false, 0},
		{"0*5", BFloat16PositiveZero, BFloat16FromFloat32(5), ModeIEEEArithmetic, RoundNearestEven, BFloat16PositiveZero, false, 0},
		{"neg*pos", BFloat16FromFloat32(-2), BFloat16FromFloat32(3), ModeIEEEArithmetic, RoundNearestEven, BFloat16FromFloat32(-6), false, 0},
		{"neg*neg", BFloat16FromFloat32(-2), BFloat16FromFloat32(-3), ModeIEEEArithmetic, RoundNearestEven, BFloat16FromFloat32(6), false, 0},
		{"0*inf IEEE", BFloat16PositiveZero, BFloat16PositiveInfinity, ModeIEEEArithmetic, RoundNearestEven, BFloat16QuietNaN, false, 0},
		{"0*inf exact", BFloat16PositiveZero, BFloat16PositiveInfinity, ModeExactArithmetic, RoundNearestEven, 0, true, ErrInvalidOperation},
		{"inf*finite", BFloat16PositiveInfinity, BFloat16FromFloat32(2), ModeIEEEArithmetic, RoundNearestEven, BFloat16PositiveInfinity, false, 0},
		{"-inf*pos", BFloat16NegativeInfinity, BFloat16FromFloat32(2), ModeIEEEArithmetic, RoundNearestEven, BFloat16NegativeInfinity, false, 0},
		{"NaN*1 IEEE", BFloat16QuietNaN, BFloat16FromFloat32(1), ModeIEEEArithmetic, RoundNearestEven, BFloat16QuietNaN, false, 0},
		{"NaN*1 exact", BFloat16QuietNaN, BFloat16FromFloat32(1), ModeExactArithmetic, RoundNearestEven, 0, true, ErrNaN},
		{"-0*pos", BFloat16NegativeZero, BFloat16FromFloat32(1), ModeIEEEArithmetic, RoundNearestEven, BFloat16NegativeZero, false, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := BFloat16MulWithMode(tt.a, tt.b, tt.mode, tt.rounding)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				fe, ok := err.(*Float16Error)
				if !ok {
					t.Fatalf("expected *Float16Error, got %T", err)
				}
				if fe.Code != tt.errCode {
					t.Errorf("error code = %d, want %d", fe.Code, tt.errCode)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tt.want.IsNaN() {
				if !got.IsNaN() {
					t.Errorf("got 0x%04X, want NaN", got.Bits())
				}
				return
			}
			if got != tt.want {
				t.Errorf("got 0x%04X (%v), want 0x%04X (%v)", got.Bits(), got, tt.want.Bits(), tt.want)
			}
		})
	}
}

func TestBFloat16DivWithMode(t *testing.T) {
	tests := []struct {
		name     string
		a, b     BFloat16
		mode     ArithmeticMode
		rounding RoundingMode
		want     BFloat16
		wantErr  bool
		errCode  ErrorCode
	}{
		{"6/2", BFloat16FromFloat32(6), BFloat16FromFloat32(2), ModeIEEEArithmetic, RoundNearestEven, BFloat16FromFloat32(3), false, 0},
		{"6/2 fast", BFloat16FromFloat32(6), BFloat16FromFloat32(2), ModeFastArithmetic, RoundNearestEven, BFloat16FromFloat32(3), false, 0},
		{"0/0 IEEE", BFloat16PositiveZero, BFloat16PositiveZero, ModeIEEEArithmetic, RoundNearestEven, BFloat16QuietNaN, false, 0},
		{"0/0 exact", BFloat16PositiveZero, BFloat16PositiveZero, ModeExactArithmetic, RoundNearestEven, 0, true, ErrInvalidOperation},
		{"1/0 IEEE", BFloat16FromFloat32(1), BFloat16PositiveZero, ModeIEEEArithmetic, RoundNearestEven, BFloat16PositiveInfinity, false, 0},
		{"1/0 exact", BFloat16FromFloat32(1), BFloat16PositiveZero, ModeExactArithmetic, RoundNearestEven, 0, true, ErrDivisionByZero},
		{"-1/0", BFloat16FromFloat32(-1), BFloat16PositiveZero, ModeIEEEArithmetic, RoundNearestEven, BFloat16NegativeInfinity, false, 0},
		{"0/1", BFloat16PositiveZero, BFloat16FromFloat32(1), ModeIEEEArithmetic, RoundNearestEven, BFloat16PositiveZero, false, 0},
		{"-0/1", BFloat16NegativeZero, BFloat16FromFloat32(1), ModeIEEEArithmetic, RoundNearestEven, BFloat16NegativeZero, false, 0},
		{"inf/inf IEEE", BFloat16PositiveInfinity, BFloat16PositiveInfinity, ModeIEEEArithmetic, RoundNearestEven, BFloat16QuietNaN, false, 0},
		{"inf/inf exact", BFloat16PositiveInfinity, BFloat16PositiveInfinity, ModeExactArithmetic, RoundNearestEven, 0, true, ErrInvalidOperation},
		{"inf/1", BFloat16PositiveInfinity, BFloat16FromFloat32(1), ModeIEEEArithmetic, RoundNearestEven, BFloat16PositiveInfinity, false, 0},
		{"1/inf", BFloat16FromFloat32(1), BFloat16PositiveInfinity, ModeIEEEArithmetic, RoundNearestEven, BFloat16PositiveZero, false, 0},
		{"NaN/1 IEEE", BFloat16QuietNaN, BFloat16FromFloat32(1), ModeIEEEArithmetic, RoundNearestEven, BFloat16QuietNaN, false, 0},
		{"NaN/1 exact", BFloat16QuietNaN, BFloat16FromFloat32(1), ModeExactArithmetic, RoundNearestEven, 0, true, ErrNaN},
		{"neg/pos", BFloat16FromFloat32(-6), BFloat16FromFloat32(2), ModeIEEEArithmetic, RoundNearestEven, BFloat16FromFloat32(-3), false, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := BFloat16DivWithMode(tt.a, tt.b, tt.mode, tt.rounding)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				fe, ok := err.(*Float16Error)
				if !ok {
					t.Fatalf("expected *Float16Error, got %T", err)
				}
				if fe.Code != tt.errCode {
					t.Errorf("error code = %d, want %d", fe.Code, tt.errCode)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tt.want.IsNaN() {
				if !got.IsNaN() {
					t.Errorf("got 0x%04X, want NaN", got.Bits())
				}
				return
			}
			if got != tt.want {
				t.Errorf("got 0x%04X (%v), want 0x%04X (%v)", got.Bits(), got, tt.want.Bits(), tt.want)
			}
		})
	}
}

func TestBFloat16FMA(t *testing.T) {
	tests := []struct {
		name    string
		a, b, c BFloat16
		wantNaN bool
		wantF32 float32
	}{
		{"2*3+1", BFloat16FromFloat32(2), BFloat16FromFloat32(3), BFloat16FromFloat32(1), false, 7},
		{"NaN*1+0", BFloat16QuietNaN, BFloat16FromFloat32(1), BFloat16PositiveZero, true, 0},
		{"1*NaN+0", BFloat16FromFloat32(1), BFloat16QuietNaN, BFloat16PositiveZero, true, 0},
		{"1*1+NaN", BFloat16FromFloat32(1), BFloat16FromFloat32(1), BFloat16QuietNaN, true, 0},
		{"-2*3+10", BFloat16FromFloat32(-2), BFloat16FromFloat32(3), BFloat16FromFloat32(10), false, 4},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := BFloat16FMA(tt.a, tt.b, tt.c)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tt.wantNaN {
				if !got.IsNaN() {
					t.Errorf("expected NaN, got %v", got)
				}
				return
			}
			gotF32 := got.ToFloat32()
			if gotF32 != tt.wantF32 {
				t.Errorf("got %v, want %v", gotF32, tt.wantF32)
			}
		})
	}
}

func TestBFloat16NaNPropagation(t *testing.T) {
	nan := BFloat16QuietNaN
	one := BFloat16FromFloat32(1)

	ops := []struct {
		name string
		fn   func() (BFloat16, error)
	}{
		{"add(NaN,1)", func() (BFloat16, error) { return BFloat16AddWithMode(nan, one, ModeIEEEArithmetic, RoundNearestEven) }},
		{"add(1,NaN)", func() (BFloat16, error) { return BFloat16AddWithMode(one, nan, ModeIEEEArithmetic, RoundNearestEven) }},
		{"sub(NaN,1)", func() (BFloat16, error) { return BFloat16SubWithMode(nan, one, ModeIEEEArithmetic, RoundNearestEven) }},
		{"mul(NaN,1)", func() (BFloat16, error) { return BFloat16MulWithMode(nan, one, ModeIEEEArithmetic, RoundNearestEven) }},
		{"div(NaN,1)", func() (BFloat16, error) { return BFloat16DivWithMode(nan, one, ModeIEEEArithmetic, RoundNearestEven) }},
		{"div(1,NaN)", func() (BFloat16, error) { return BFloat16DivWithMode(one, nan, ModeIEEEArithmetic, RoundNearestEven) }},
	}

	for _, op := range ops {
		t.Run(op.name, func(t *testing.T) {
			got, err := op.fn()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !got.IsNaN() {
				t.Errorf("expected NaN, got 0x%04X (%v)", got.Bits(), got)
			}
		})
	}
}

func TestBFloat16GradualUnderflow(t *testing.T) {
	// Multiplying two very small normal numbers should produce a subnormal
	// rather than flushing to zero.
	smallest := BFloat16SmallestPos // smallest positive normal
	half := BFloat16FromFloat32(0.5)

	got, err := BFloat16MulWithMode(smallest, half, ModeIEEEArithmetic, RoundNearestEven)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// The result should be subnormal (half the smallest normal)
	if got.IsZero() {
		t.Error("expected subnormal result, got zero (gradual underflow not working)")
	}

	// Verify the result is approximately half of the smallest normal
	gotF := got.ToFloat32()
	wantF := smallest.ToFloat32() * 0.5
	if math.Abs(float64(gotF-wantF)) > float64(wantF)*0.1 {
		t.Errorf("got %e, want approximately %e", gotF, wantF)
	}
}

func TestBFloat16ArithmeticWithMode(t *testing.T) {
	// Integration test: verifies all WithMode functions work together
	a := BFloat16FromFloat32(10)
	b := BFloat16FromFloat32(3)

	sum, err := BFloat16AddWithMode(a, b, ModeIEEEArithmetic, RoundNearestEven)
	if err != nil {
		t.Fatalf("add: %v", err)
	}
	if sum.ToFloat32() != 13 {
		t.Errorf("add: got %v, want 13", sum.ToFloat32())
	}

	diff, err := BFloat16SubWithMode(a, b, ModeIEEEArithmetic, RoundNearestEven)
	if err != nil {
		t.Fatalf("sub: %v", err)
	}
	if diff.ToFloat32() != 7 {
		t.Errorf("sub: got %v, want 7", diff.ToFloat32())
	}

	prod, err := BFloat16MulWithMode(a, b, ModeIEEEArithmetic, RoundNearestEven)
	if err != nil {
		t.Fatalf("mul: %v", err)
	}
	if prod.ToFloat32() != 30 {
		t.Errorf("mul: got %v, want 30", prod.ToFloat32())
	}

	quot, err := BFloat16DivWithMode(BFloat16FromFloat32(6), b, ModeIEEEArithmetic, RoundNearestEven)
	if err != nil {
		t.Fatalf("div: %v", err)
	}
	if quot.ToFloat32() != 2 {
		t.Errorf("div: got %v, want 2", quot.ToFloat32())
	}
}
