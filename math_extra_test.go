package float16

import (
	"testing"
)

func TestAsin(t *testing.T) {
	tests := []struct {
		name string
		arg  Float16
		want Float16
	}{
		{"Asin(0)", PositiveZero, PositiveZero},
		{"Asin(1)", One(), 0x3E48},
		{"Asin(-1)", One().Neg(), 0xBE48},
		{"Asin(NaN)", QuietNaN, QuietNaN},
		{"Asin(2)", FromInt(2), QuietNaN},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Asin(tt.arg)
			if !Equal(got, tt.want) && !got.IsNaN() {
				t.Errorf("Asin() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAcos(t *testing.T) {
	tests := []struct {
		name string
		arg  Float16
		want Float16
	}{
		{"Acos(1)", One(), PositiveZero},
		{"Acos(NaN)", QuietNaN, QuietNaN},
		{"Acos(2)", FromInt(2), QuietNaN},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Acos(tt.arg)
			if !Equal(got, tt.want) && !got.IsNaN() {
				t.Errorf("Acos() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAtan(t *testing.T) {
	tests := []struct {
		name string
		arg  Float16
		want Float16
	}{
		{"Atan(0)", PositiveZero, PositiveZero},
		{"Atan(inf)", PositiveInfinity, 0x3E48},
		{"Atan(-inf)", NegativeInfinity, 0xBE48},
		{"Atan(NaN)", QuietNaN, QuietNaN},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Atan(tt.arg)
			if !Equal(got, tt.want) && !got.IsNaN() {
				t.Errorf("Atan() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAtan2(t *testing.T) {
	tests := []struct {
		name string
		y, x Float16
		want Float16
	}{
		{"Atan2(0, 0)", PositiveZero, PositiveZero, PositiveZero},
		{"Atan2(NaN, 1)", QuietNaN, One(), QuietNaN},
		{"Atan2(1, NaN)", One(), QuietNaN, QuietNaN},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Atan2(tt.y, tt.x)
			if !Equal(got, tt.want) && !got.IsNaN() {
				t.Errorf("Atan2() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSinh(t *testing.T) {
	tests := []struct {
		name string
		arg  Float16
		want Float16
	}{
		{"Sinh(0)", PositiveZero, PositiveZero},
		{"Sinh(inf)", PositiveInfinity, PositiveInfinity},
		{"Sinh(-inf)", NegativeInfinity, NegativeInfinity},
		{"Sinh(NaN)", QuietNaN, QuietNaN},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Sinh(tt.arg)
			if !Equal(got, tt.want) && !got.IsNaN() {
				t.Errorf("Sinh() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCosh(t *testing.T) {
	tests := []struct {
		name string
		arg  Float16
		want Float16
	}{
		{"Cosh(0)", PositiveZero, One()},
		{"Cosh(inf)", PositiveInfinity, PositiveInfinity},
		{"Cosh(-inf)", NegativeInfinity, PositiveInfinity},
		{"Cosh(NaN)", QuietNaN, QuietNaN},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Cosh(tt.arg)
			if !Equal(got, tt.want) && !got.IsNaN() {
				t.Errorf("Cosh() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTanh(t *testing.T) {
	tests := []struct {
		name string
		arg  Float16
		want Float16
	}{
		{"Tanh(0)", PositiveZero, PositiveZero},
		{"Tanh(inf)", PositiveInfinity, One()},
		{"Tanh(-inf)", NegativeInfinity, One().Neg()},
		{"Tanh(NaN)", QuietNaN, QuietNaN},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Tanh(tt.arg)
			if !Equal(got, tt.want) && !got.IsNaN() {
				t.Errorf("Tanh() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRoundToEven(t *testing.T) {
	tests := []struct {
		name string
		arg  Float16
		want Float16
	}{
		{"RoundToEven(2.5)", FromFloat64(2.5), FromFloat64(2)},
		{"RoundToEven(3.5)", FromFloat64(3.5), FromFloat64(4)},
		{"RoundToEven(NaN)", QuietNaN, QuietNaN},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RoundToEven(tt.arg)
			if !Equal(got, tt.want) && !got.IsNaN() {
				t.Errorf("RoundToEven() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRemainder(t *testing.T) {
	tests := []struct {
		name string
		x, y Float16
		want Float16
	}{
		{"Remainder(5, 3)", FromInt(5), FromInt(3), FromInt(-1)},
		{"Remainder(0, 1)", PositiveZero, One(), PositiveZero},
		{"Remainder(1, 0)", One(), PositiveZero, QuietNaN},
		{"Remainder(NaN, 1)", QuietNaN, One(), QuietNaN},
		{"Remainder(1, NaN)", One(), QuietNaN, QuietNaN},
		{"Remainder(inf, 1)", PositiveInfinity, One(), QuietNaN},
		{"Remainder(1, inf)", One(), PositiveInfinity, One()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Remainder(tt.x, tt.y)
			if !Equal(got, tt.want) && !got.IsNaN() {
				t.Errorf("Remainder() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClamp(t *testing.T) {
	tests := []struct {
		name        string
		x, min, max Float16
		want        Float16
	}{
		{"Clamp(1, 0, 2)", One(), PositiveZero, FromInt(2), One()},
		{"Clamp(-1, 0, 2)", One().Neg(), PositiveZero, FromInt(2), PositiveZero},
		{"Clamp(3, 0, 2)", FromInt(3), PositiveZero, FromInt(2), FromInt(2)},
		{"Clamp(NaN, 0, 2)", QuietNaN, PositiveZero, FromInt(2), QuietNaN},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Clamp(tt.x, tt.min, tt.max)
			if !Equal(got, tt.want) && !got.IsNaN() {
				t.Errorf("Clamp() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLerp(t *testing.T) {
	tests := []struct {
		name    string
		a, b, t Float16
		want    Float16
	}{
		{"Lerp(0, 1, 0.5)", PositiveZero, One(), FromFloat64(0.5), FromFloat64(0.5)},
		{"Lerp(0, 1, 0)", PositiveZero, One(), PositiveZero, PositiveZero},
		{"Lerp(0, 1, 1)", PositiveZero, One(), One(), One()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Lerp(tt.a, tt.b, tt.t)
			if !Equal(got, tt.want) {
				t.Errorf("Lerp() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSign(t *testing.T) {
	tests := []struct {
		name string
		arg  Float16
		want Float16
	}{
		{"Sign(1)", One(), One()},
		{"Sign(-1)", One().Neg(), One().Neg()},
		{"Sign(0)", PositiveZero, PositiveZero},
		{"Sign(NaN)", QuietNaN, QuietNaN},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Sign(tt.arg)
			if !Equal(got, tt.want) && !got.IsNaN() {
				t.Errorf("Sign() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGamma(t *testing.T) {
	tests := []struct {
		name string
		arg  Float16
		want Float16
	}{
		{"Gamma(NaN)", QuietNaN, QuietNaN},
		{"Gamma(-inf)", NegativeInfinity, QuietNaN},
		{"Gamma(+inf)", PositiveInfinity, PositiveInfinity},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Gamma(tt.arg)
			if !Equal(got, tt.want) && !got.IsNaN() {
				t.Errorf("Gamma() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestY0(t *testing.T) {
	tests := []struct {
		name string
		arg  Float16
		want Float16
	}{
		{"Y0(NaN)", QuietNaN, QuietNaN},
		{"Y0(-1)", One().Neg(), QuietNaN},
		{"Y0(0)", PositiveZero, NegativeInfinity},
		{"Y0(inf)", PositiveInfinity, PositiveZero},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Y0(tt.arg)
			if !Equal(got, tt.want) && !got.IsNaN() {
				t.Errorf("Y0() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestY1(t *testing.T) {
	tests := []struct {
		name string
		arg  Float16
		want Float16
	}{
		{"Y1(NaN)", QuietNaN, QuietNaN},
		{"Y1(-1)", One().Neg(), QuietNaN},
		{"Y1(0)", PositiveZero, NegativeInfinity},
		{"Y1(inf)", PositiveInfinity, PositiveZero},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Y1(tt.arg)
			if !Equal(got, tt.want) && !got.IsNaN() {
				t.Errorf("Y1() = %v, want %v", got, tt.want)
			}
		})
	}
}
