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

func TestBFloat16NaNPropagationAllModes(t *testing.T) {
	nan := BFloat16QuietNaN
	one := BFloat16FromFloat32(1)

	type opFunc func(a, b BFloat16, mode ArithmeticMode, rounding RoundingMode) (BFloat16, error)

	ops := []struct {
		name string
		fn   opFunc
		a, b BFloat16
	}{
		{"add(NaN,1)", BFloat16AddWithMode, nan, one},
		{"add(1,NaN)", BFloat16AddWithMode, one, nan},
		{"sub(NaN,1)", BFloat16SubWithMode, nan, one},
		{"sub(1,NaN)", BFloat16SubWithMode, one, nan},
		{"mul(NaN,1)", BFloat16MulWithMode, nan, one},
		{"mul(1,NaN)", BFloat16MulWithMode, one, nan},
		{"div(NaN,1)", BFloat16DivWithMode, nan, one},
		{"div(1,NaN)", BFloat16DivWithMode, one, nan},
	}

	modes := []struct {
		name    string
		mode    ArithmeticMode
		wantErr bool
	}{
		{"IEEE", ModeIEEEArithmetic, false},
		{"fast", ModeFastArithmetic, false},
		{"exact", ModeExactArithmetic, true},
	}

	for _, op := range ops {
		for _, m := range modes {
			t.Run(op.name+"/"+m.name, func(t *testing.T) {
				got, err := op.fn(op.a, op.b, m.mode, RoundNearestEven)
				if m.wantErr {
					if err == nil {
						t.Fatal("expected error for NaN in exact mode, got nil")
					}
					fe, ok := err.(*Float16Error)
					if !ok {
						t.Fatalf("expected *Float16Error, got %T", err)
					}
					if fe.Code != ErrNaN {
						t.Errorf("error code = %d, want %d (ErrNaN)", fe.Code, ErrNaN)
					}
					return
				}
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if !got.IsNaN() {
					t.Errorf("expected NaN, got 0x%04X (%v)", got.Bits(), got)
				}
			})
		}
	}
}

func TestBFloat16GradualUnderflow(t *testing.T) {
	smallest := BFloat16SmallestPos // smallest positive normal
	half := BFloat16FromFloat32(0.5)

	t.Run("mul/smallest*0.5", func(t *testing.T) {
		got, err := BFloat16MulWithMode(smallest, half, ModeIEEEArithmetic, RoundNearestEven)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got.IsZero() {
			t.Fatal("expected subnormal result, got zero")
		}
		gotF := got.ToFloat32()
		wantF := smallest.ToFloat32() * 0.5
		if math.Abs(float64(gotF-wantF)) > float64(wantF)*0.1 {
			t.Errorf("got %e, want approximately %e", gotF, wantF)
		}
	})

	t.Run("mul/neg_underflow", func(t *testing.T) {
		neg := BFloat16Neg(smallest)
		got, err := BFloat16MulWithMode(neg, half, ModeIEEEArithmetic, RoundNearestEven)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got.IsZero() {
			t.Fatal("expected negative subnormal, got zero")
		}
		if got.ToFloat32() >= 0 {
			t.Errorf("expected negative result, got %e", got.ToFloat32())
		}
	})

	t.Run("add/near_subnormal_boundary", func(t *testing.T) {
		// Adding two values that sum to something below the smallest normal
		// should produce a subnormal, not zero.
		sub := BFloat16SmallestPosSubnormal
		got, err := BFloat16AddWithMode(sub, sub, ModeIEEEArithmetic, RoundNearestEven)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got.IsZero() {
			t.Fatal("expected non-zero subnormal sum, got zero")
		}
		wantF := BFloat16SmallestPosSubnormal.ToFloat32() * 2
		gotF := got.ToFloat32()
		if math.Abs(float64(gotF-wantF)) > float64(wantF)*0.1 {
			t.Errorf("got %e, want approximately %e", gotF, wantF)
		}
	})

	t.Run("div/smallest/2", func(t *testing.T) {
		two := BFloat16FromFloat32(2)
		got, err := BFloat16DivWithMode(smallest, two, ModeIEEEArithmetic, RoundNearestEven)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got.IsZero() {
			t.Fatal("expected subnormal result, got zero")
		}
		gotF := got.ToFloat32()
		wantF := smallest.ToFloat32() / 2
		if math.Abs(float64(gotF-wantF)) > float64(wantF)*0.1 {
			t.Errorf("got %e, want approximately %e", gotF, wantF)
		}
	})

	t.Run("sub/subnormal_boundary", func(t *testing.T) {
		// Subtracting values that are very close should yield a subnormal.
		a := BFloat16FromFloat32(smallest.ToFloat32() * 1.5)
		got, err := BFloat16SubWithMode(a, smallest, ModeIEEEArithmetic, RoundNearestEven)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		// Result should be approximately 0.5 * smallest normal = subnormal
		if got.IsZero() {
			t.Fatal("expected subnormal result, got zero")
		}
	})
}

func TestBFloat16FMACorrectness(t *testing.T) {
	tests := []struct {
		name    string
		a, b, c BFloat16
		wantNaN bool
		wantF32 float32
	}{
		{"2*3+1=7", BFloat16FromFloat32(2), BFloat16FromFloat32(3), BFloat16FromFloat32(1), false, 7},
		{"-2*3+10=4", BFloat16FromFloat32(-2), BFloat16FromFloat32(3), BFloat16FromFloat32(10), false, 4},
		{"0*5+3=3", BFloat16PositiveZero, BFloat16FromFloat32(5), BFloat16FromFloat32(3), false, 3},
		{"5*0+3=3", BFloat16FromFloat32(5), BFloat16PositiveZero, BFloat16FromFloat32(3), false, 3},
		{"1*1+0=1", BFloat16FromFloat32(1), BFloat16FromFloat32(1), BFloat16PositiveZero, false, 1},
		{"-1*-1+0=1", BFloat16FromFloat32(-1), BFloat16FromFloat32(-1), BFloat16PositiveZero, false, 1},
		{"4*0.5+-2=0", BFloat16FromFloat32(4), BFloat16FromFloat32(0.5), BFloat16FromFloat32(-2), false, 0},
		// NaN in each operand position
		{"NaN*1+0", BFloat16QuietNaN, BFloat16FromFloat32(1), BFloat16PositiveZero, true, 0},
		{"1*NaN+0", BFloat16FromFloat32(1), BFloat16QuietNaN, BFloat16PositiveZero, true, 0},
		{"1*1+NaN", BFloat16FromFloat32(1), BFloat16FromFloat32(1), BFloat16QuietNaN, true, 0},
		{"NaN*NaN+NaN", BFloat16QuietNaN, BFloat16QuietNaN, BFloat16QuietNaN, true, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := BFloat16FMA(tt.a, tt.b, tt.c)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tt.wantNaN {
				if !got.IsNaN() {
					t.Errorf("expected NaN, got %v (0x%04X)", got, got.Bits())
				}
				return
			}
			gotF32 := got.ToFloat32()
			if gotF32 != tt.wantF32 {
				t.Errorf("got %v, want %v", gotF32, tt.wantF32)
			}
		})
	}

	// FMA precision test: verify fused multiply-add avoids intermediate rounding.
	// For values where a*b overflows float16 range but a*b+c is representable,
	// FMA via float64 should give a more accurate result.
	t.Run("precision/no_intermediate_rounding", func(t *testing.T) {
		a := BFloat16FromFloat32(100)
		b := BFloat16FromFloat32(100)
		c := BFloat16FromFloat32(-9984) // 100*100 = 10000; 10000 - 9984 = 16
		got, err := BFloat16FMA(a, b, c)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		gotF := got.ToFloat32()
		wantF := float32(math.FMA(float64(a.ToFloat32()), float64(b.ToFloat32()), float64(c.ToFloat32())))
		if gotF != wantF {
			t.Errorf("got %v, want %v", gotF, wantF)
		}
	})
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
