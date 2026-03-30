package float16

import (
	"errors"
	"math"
	"testing"
)

func TestBFloat16Error(t *testing.T) {
	t.Run("Error_WithOp", func(t *testing.T) {
		e := &BFloat16Error{Op: "BFloat16FromFloat32", Msg: "overflow in strict mode", Code: ErrOverflow}
		got := e.Error()
		want := "bfloat16 BFloat16FromFloat32: overflow in strict mode"
		if got != want {
			t.Errorf("Error() = %q, want %q", got, want)
		}
	})

	t.Run("Error_WithoutOp", func(t *testing.T) {
		e := &BFloat16Error{Msg: "some error", Code: ErrNaN}
		got := e.Error()
		want := "bfloat16: some error"
		if got != want {
			t.Errorf("Error() = %q, want %q", got, want)
		}
	})

	t.Run("Error_Nil", func(t *testing.T) {
		var e *BFloat16Error
		if e.Error() != "<nil>" {
			t.Errorf("nil BFloat16Error.Error() = %q, want %q", e.Error(), "<nil>")
		}
	})

	t.Run("StrictConversion_ReturnsTypedError", func(t *testing.T) {
		_, err := BFloat16FromFloat32WithMode(math.MaxFloat32, ModeStrict, RoundNearestEven)
		if err == nil {
			t.Fatal("expected error for overflow in strict mode")
		}
		var bfErr *BFloat16Error
		if !errors.As(err, &bfErr) {
			t.Fatalf("expected *BFloat16Error, got %T", err)
		}
		if bfErr.Code != ErrOverflow {
			t.Errorf("Code = %v, want ErrOverflow", bfErr.Code)
		}
	})

	t.Run("StrictConversion_NaN_ReturnsTypedError", func(t *testing.T) {
		_, err := BFloat16FromFloat32WithMode(float32(math.NaN()), ModeStrict, RoundNearestEven)
		if err == nil {
			t.Fatal("expected error for NaN in strict mode")
		}
		var bfErr *BFloat16Error
		if !errors.As(err, &bfErr) {
			t.Fatalf("expected *BFloat16Error, got %T", err)
		}
		if bfErr.Code != ErrNaN {
			t.Errorf("Code = %v, want ErrNaN", bfErr.Code)
		}
	})

	t.Run("CheckedArithmetic_ReturnsTypedError", func(t *testing.T) {
		nan := BFloat16QuietNaN
		_, err := BFloat16AddWithMode(nan, BFloat16One, ModeExactArithmetic, RoundNearestEven)
		if err == nil {
			t.Fatal("expected error for NaN in exact mode")
		}
		var bfErr *BFloat16Error
		if !errors.As(err, &bfErr) {
			t.Fatalf("expected *BFloat16Error, got %T", err)
		}
		if bfErr.Code != ErrNaN {
			t.Errorf("Code = %v, want ErrNaN", bfErr.Code)
		}
	})

	t.Run("CheckedDivByZero_ReturnsTypedError", func(t *testing.T) {
		_, err := BFloat16DivWithMode(BFloat16One, BFloat16PositiveZero, ModeExactArithmetic, RoundNearestEven)
		if err == nil {
			t.Fatal("expected error for division by zero in exact mode")
		}
		var bfErr *BFloat16Error
		if !errors.As(err, &bfErr) {
			t.Fatalf("expected *BFloat16Error, got %T", err)
		}
		if bfErr.Code != ErrDivisionByZero {
			t.Errorf("Code = %v, want ErrDivisionByZero", bfErr.Code)
		}
	})
}

func TestFromInt(t *testing.T) {
	tests := []struct {
		name string
		i    int
		want Float16
	}{
		{"FromInt(0)", 0, PositiveZero},
		{"FromInt(1)", 1, 0x3C00},
		{"FromInt(-1)", -1, 0xBC00},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FromInt(tt.i); got != tt.want {
				t.Errorf("FromInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFromInt32(t *testing.T) {
	tests := []struct {
		name string
		i    int32
		want Float16
	}{
		{"FromInt32(0)", 0, PositiveZero},
		{"FromInt32(1)", 1, 0x3C00},
		{"FromInt32(-1)", -1, 0xBC00},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FromInt32(tt.i); got != tt.want {
				t.Errorf("FromInt32() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFromInt64(t *testing.T) {
	tests := []struct {
		name string
		i    int64
		want Float16
	}{
		{"FromInt64(0)", 0, PositiveZero},
		{"FromInt64(1)", 1, 0x3C00},
		{"FromInt64(-1)", -1, 0xBC00},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FromInt64(tt.i); got != tt.want {
				t.Errorf("FromInt64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToInt(t *testing.T) {
	tests := []struct {
		name string
		f    Float16
		want int
	}{
		{"ToInt(0)", PositiveZero, 0},
		{"ToInt(1.0)", 0x3C00, 1},
		{"ToInt(-1.0)", 0xBC00, -1},
		{"ToInt(1.9)", 0x3F33, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.f.ToInt(); got != tt.want {
				t.Errorf("ToInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToInt32(t *testing.T) {
	tests := []struct {
		name string
		f    Float16
		want int32
	}{
		{"ToInt32(0)", PositiveZero, 0},
		{"ToInt32(1.0)", 0x3C00, 1},
		{"ToInt32(-1.0)", 0xBC00, -1},
		{"ToInt32(1.9)", 0x3F33, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.f.ToInt32(); got != tt.want {
				t.Errorf("ToInt32() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToInt64(t *testing.T) {
	tests := []struct {
		name string
		f    Float16
		want int64
	}{
		{"ToInt64(0)", PositiveZero, 0},
		{"ToInt64(1.0)", 0x3C00, 1},
		{"ToInt64(-1.0)", 0xBC00, -1},
		{"ToInt64(1.9)", 0x3F33, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.f.ToInt64(); got != tt.want {
				t.Errorf("ToInt64() = %v, want %v", got, tt.want)
			}
		})
	}
}
