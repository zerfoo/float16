package float16

import (
	"fmt"
	"math"
	"testing"
)

func modes() []RoundingMode {
	return []RoundingMode{
		RoundNearestEven,
		RoundNearestAway,
		RoundTowardZero,
		RoundTowardPositive,
		RoundTowardNegative,
	}
}

func TestAddWithMode_RoundingMatchesConverter(t *testing.T) {
	cases := [][2]float32{
		{1.0, float32(math.Pow(2, -11))},     // halfway between 1.0 and next
		{1.0, 1e-3},                          // general positive
		{-1.0, float32(math.Pow(2, -11))},    // negative with halfway increment
		{-0.75, 0.125},                       // mixed signs, exact binary fractions
	}

	for _, c := range cases {
		for _, m := range modes() {
			name := func(a, b float32, mode RoundingMode) string {
				return fmt.Sprintf("a=%g b=%g mode=%v", a, b, mode)
			}(c[0], c[1], m)
			t.Run(name, func(t *testing.T) {
				a16 := FromFloat32(c[0])
				b16 := FromFloat32(c[1])
				expected := FromFloat32WithRounding(a16.ToFloat32()+b16.ToFloat32(), m)
				got, err := AddWithMode(a16, b16, ModeIEEEArithmetic, m)
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if got != expected {
					t.Fatalf("AddWithMode mismatch: got=%v want=%v (a=%g b=%g mode=%v)", got, expected, c[0], c[1], m)
				}
			})
		}
	}
}

func TestMulWithMode_RoundingMatchesConverter(t *testing.T) {
	cases := [][2]float32{
		{1.25, 0.2},          // positive * positive
		{-1.25, 0.2},         // negative * positive
		{1.5, -0.75},         // positive * negative
		{-0.5, -0.125},       // negative * negative
		{float32(math.Pow(2, -3)), float32(math.Pow(2, -8))}, // exact powers of two
	}

	for _, c := range cases {
		for _, m := range modes() {
			name := func(a, b float32, mode RoundingMode) string {
				return fmt.Sprintf("a=%g b=%g mode=%v", a, b, mode)
			}(c[0], c[1], m)
			t.Run(name, func(t *testing.T) {
				a16 := FromFloat32(c[0])
				b16 := FromFloat32(c[1])
				expected := FromFloat32WithRounding(a16.ToFloat32()*b16.ToFloat32(), m)
				got, err := MulWithMode(a16, b16, ModeIEEEArithmetic, m)
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if got != expected {
					t.Fatalf("MulWithMode mismatch: got=%v want=%v (a=%g b=%g mode=%v)", got, expected, c[0], c[1], m)
				}
			})
		}
	}
}

func TestDivWithMode_RoundingMatchesConverter(t *testing.T) {
	cases := [][2]float32{
		{1.25, 0.2},          // positive / positive
		{-1.25, 0.2},         // negative / positive
		{1.5, -0.75},         // positive / negative
		{-0.5, -0.125},       // negative / negative
		{7.0, 3.0},           // non-terminating in binary
	}

	for _, c := range cases {
		// avoid division by zero
		if c[1] == 0 {
			continue
		}
		for _, m := range modes() {
			name := func(a, b float32, mode RoundingMode) string {
				return fmt.Sprintf("a=%g b=%g mode=%v", a, b, mode)
			}(c[0], c[1], m)
			t.Run(name, func(t *testing.T) {
				a16 := FromFloat32(c[0])
				b16 := FromFloat32(c[1])
				expected := FromFloat32WithRounding(a16.ToFloat32()/b16.ToFloat32(), m)
				got, err := DivWithMode(a16, b16, ModeIEEEArithmetic, m)
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if got != expected {
					t.Fatalf("DivWithMode mismatch: got=%v want=%v (a=%g b=%g mode=%v)", got, expected, c[0], c[1], m)
				}
			})
		}
	}
}
