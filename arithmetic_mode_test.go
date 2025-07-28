package float16

import "testing"

func TestSubWithMode(t *testing.T) {
	tests := []struct {
		name     string
		a, b     Float16
		mode     ArithmeticMode
		rounding RoundingMode
		want     Float16
		wantErr  bool
	}{
		{
			name:     "exact mode with valid inputs",
			a:        FromBits(0x4000), // 2.0
			b:        FromBits(0x3C00), // 1.0
			mode:     ModeExactArithmetic,
			rounding: RoundNearestEven,
			want:     FromBits(0x3C00), // 1.0
			wantErr:  false,
		},
		{
			name:     "exact mode with NaN",
			a:        NaN(),
			b:        FromBits(0x3C00), // 1.0
			mode:     ModeExactArithmetic,
			rounding: RoundNearestEven,
			want:     0,
			wantErr:  true,
		},
		{
			name:     "fast mode with valid inputs",
			a:        FromBits(0x4200), // 3.0
			b:        FromBits(0x3C00), // 1.0
			mode:     ModeFastArithmetic,
			rounding: RoundNearestEven,
			want:     FromBits(0x4000), // 2.0
			wantErr:  false,
		},
		{
			name:     "IEEE mode with valid inputs",
			a:        FromBits(0x4400), // 4.0
			b:        FromBits(0x3C00), // 1.0
			mode:     ModeIEEEArithmetic,
			rounding: RoundNearestEven,
			want:     FromBits(0x4200), // 3.0
			wantErr:  false,
		},
		{
			name:     "infinity minus infinity",
			a:        PositiveInfinity,
			b:        PositiveInfinity,
			mode:     ModeExactArithmetic,
			rounding: RoundNearestEven,
			want:     0,
			wantErr:  true,
		},
		{
			name:     "IEEE mode with valid inputs",
			a:        FromBits(0x4400), // 4.0
			b:        FromBits(0x3C00), // 1.0
			mode:     ModeIEEEArithmetic,
			rounding: RoundNearestEven,
			want:     FromBits(0x4200), // 3.0
			wantErr:  false,
		},
		{
			name:     "SubWithMode - Negative infinity minus negative infinity",
			a:        NegativeInfinity,
			b:        NegativeInfinity,
			mode:     ModeExactArithmetic, // This operation should only error in exact mode
			rounding: RoundNearestEven,
			want:     QuietNaN,
			wantErr:  true,
		},
		{
			name:     "IEEE mode with negative infinity minus negative infinity",
			a:        NegativeInfinity,
			b:        NegativeInfinity,
			mode:     ModeIEEEArithmetic,
			rounding: RoundNearestEven,
			want:     QuietNaN, // Should return NaN without error in IEEE mode
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SubWithMode(tt.a, tt.b, tt.mode, tt.rounding)
			if (err != nil) != tt.wantErr {
				t.Errorf("SubWithMode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("SubWithMode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMulWithMode(t *testing.T) {
	tests := []struct {
		name     string
		a, b     Float16
		mode     ArithmeticMode
		rounding RoundingMode
		want     Float16
		wantErr  bool
	}{
		{
			name:     "exact mode with valid inputs",
			a:        FromBits(0x4000), // 2.0
			b:        FromBits(0x4200), // 3.0
			mode:     ModeExactArithmetic,
			rounding: RoundNearestEven,
			want:     FromBits(0x4600), // 6.0 (0x4600 is 6.0 in float16)
			wantErr:  false,
		},
		{
			name:     "exact mode with NaN",
			a:        NaN(),
			b:        FromBits(0x3C00), // 1.0
			mode:     ModeExactArithmetic,
			rounding: RoundNearestEven,
			want:     0,
			wantErr:  true,
		},
		{
			name:     "infinity times zero in exact mode",
			a:        Infinity(1),
			b:        FromBits(0x0000), // 0.0
			mode:     ModeExactArithmetic,
			rounding: RoundNearestEven,
			want:     0,
			wantErr:  true,
		},
		{
			name:     "infinity times zero in IEEE mode",
			a:        Infinity(1),
			b:        FromBits(0x0000), // 0.0
			mode:     ModeIEEEArithmetic,
			rounding: RoundNearestEven,
			want:     QuietNaN,
			wantErr:  false,
		},
		{
			name:     "Zero times infinity in exact mode",
			a:        PositiveZero,
			b:        PositiveInfinity,
			mode:     ModeExactArithmetic,
			rounding: RoundNearestEven,
			want:     0,
			wantErr:  true,
		},
		{
			name:     "Zero times infinity in IEEE mode",
			a:        PositiveZero,
			b:        PositiveInfinity,
			mode:     ModeIEEEArithmetic,
			rounding: RoundNearestEven,
			want:     QuietNaN,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := MulWithMode(tt.a, tt.b, tt.mode, tt.rounding)
			if (err != nil) != tt.wantErr {
				t.Errorf("MulWithMode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("MulWithMode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDivWithMode(t *testing.T) {
	tests := []struct {
		name     string
		a, b     Float16
		mode     ArithmeticMode
		rounding RoundingMode
		want     Float16
		wantErr  bool
	}{
		{
			name:     "exact mode with valid inputs",
			a:        FromBits(0x4600), // 6.0 (0x4600 = 6.0, 0x4800 = 8.0)
			b:        FromBits(0x4200), // 3.0
			mode:     ModeExactArithmetic,
			rounding: RoundNearestEven,
			want:     FromBits(0x4000), // 2.0 (0x4000 = 2.0, 0x4400 = 4.0)
			wantErr:  false,
		},
		{
			name:     "exact mode with NaN",
			a:        NaN(),
			b:        FromBits(0x3C00), // 1.0
			mode:     ModeExactArithmetic,
			rounding: RoundNearestEven,
			want:     0,
			wantErr:  true,
		},
		{
			name:     "division by zero",
			a:        FromBits(0x3C00), // 1.0
			b:        FromBits(0x0000), // 0.0
			mode:     ModeExactArithmetic,
			rounding: RoundNearestEven,
			want:     0,
			wantErr:  true,
		},
		{
			name:     "infinity divided by infinity",
			a:        Infinity(1),
			b:        Infinity(1),
			mode:     ModeExactArithmetic,
			rounding: RoundNearestEven,
			want:     0,
			wantErr:  true,
		},
		{
			name:     "DivWithMode - Division by zero",
			a:        FromBits(0x3C00),
			b:        PositiveZero,
			mode:     ModeIEEEArithmetic,
			rounding: RoundNearestEven,
			want:     PositiveInfinity,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DivWithMode(tt.a, tt.b, tt.mode, tt.rounding)
			if (err != nil) != tt.wantErr {
				t.Errorf("DivWithMode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("DivWithMode() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Infinity returns positive or negative infinity based on the sign parameter
func Infinity(sign int) Float16 {
	if sign >= 0 {
		return PositiveInfinity
	}
	return NegativeInfinity
}

func TestSliceOperationsWithMode(t *testing.T) {
	t.Run("SubSlice", func(t *testing.T) {
		a := []Float16{FromBits(0x4400), FromBits(0x4500), FromBits(0x4600)} // [4.0, 5.0, 6.0]
		b := []Float16{FromBits(0x3C00), FromBits(0x4000), FromBits(0x4200)} // [1.0, 2.0, 3.0]
		want := []Float16{FromBits(0x4200), FromBits(0x4200), FromBits(0x4200)} // [3.0, 3.0, 3.0]
		got := SubSlice(a, b)
		if len(got) != len(want) {
			t.Fatalf("SubSlice() length = %d, want %d", len(got), len(want))
		}
		for i := range got {
			if got[i] != want[i] {
				t.Errorf("SubSlice()[%d] = %v, want %v", i, got[i], want[i])
			}
		}
	})

	t.Run("MulSlice", func(t *testing.T) {
		a := []Float16{FromBits(0x3C00), FromBits(0x4000), FromBits(0x4400)} // [1.0, 2.0, 4.0]
		b := []Float16{FromBits(0x4400), FromBits(0x4400), FromBits(0x4400)} // [4.0, 4.0, 4.0]
		want := []Float16{FromBits(0x4400), FromBits(0x4800), FromBits(0x4C00)} // [4.0, 8.0, 16.0]
		got := MulSlice(a, b)
		if len(got) != len(want) {
			t.Fatalf("MulSlice() length = %d, want %d", len(got), len(want))
		}
		for i := range got {
			if got[i] != want[i] {
				t.Errorf("MulSlice()[%d] = %v, want %v", i, got[i], want[i])
			}
		}
	})

	t.Run("DivSlice", func(t *testing.T) {
		a := []Float16{FromBits(0x4400), FromBits(0x4800), FromBits(0x4C00)} // [4.0, 8.0, 16.0]
		b := []Float16{FromBits(0x3C00), FromBits(0x4000), FromBits(0x4400)} // [1.0, 2.0, 4.0]
		want := []Float16{FromBits(0x4400), FromBits(0x4400), FromBits(0x4400)} // [4.0, 4.0, 4.0]
		got := DivSlice(a, b)
		if len(got) != len(want) {
			t.Fatalf("DivSlice() length = %d, want %d", len(got), len(want))
		}
		for i := range got {
			if got[i] != want[i] {
				t.Errorf("DivSlice()[%d] = %v, want %v", i, got[i], want[i])
			}
		}
	})

	t.Run("ScaleSlice", func(t *testing.T) {
		s := []Float16{FromBits(0x3C00), FromBits(0x4000), FromBits(0x4200)}
		scalar := FromBits(0x4000)
		want := []Float16{FromBits(0x4000), FromBits(0x4400), FromBits(0x4600)}
		got := ScaleSlice(s, scalar)
		if len(got) != len(want) {
			t.Fatalf("ScaleSlice() length = %d, want %d", len(got), len(want))
		}
		for i := range got {
			if got[i] != want[i] {
				t.Errorf("ScaleSlice()[%d] = %v, want %v", i, got[i], want[i])
			}
		}
	})

	t.Run("SumSlice", func(t *testing.T) {
		s := []Float16{FromBits(0x3C00), FromBits(0x4000), FromBits(0x4200)}
		want := FromBits(0x4600) // 6.0
		got := SumSlice(s)
		if got != want {
			t.Errorf("SumSlice() = %v, want %v", got, want)
		}
	})

	t.Run("Norm2", func(t *testing.T) {
		s := []Float16{FromBits(0x4200), FromBits(0x4400)} // 3-4-5 right triangle
		want := FromBits(0x4500) // 5.0
		got := Norm2(s)
		if got != want {
			t.Errorf("Norm2() = %v, want %v", got, want)
		}
	})
}
