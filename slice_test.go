package float16

import (
	"math"
	"testing"
)

func TestToSlice64(t *testing.T) {
	tests := []struct {
		name  string
		input []Float16
		want  []float64
	}{
		{
			"Empty slice",
			[]Float16{},
			[]float64{},
		},
		{
			"Single element",
			[]Float16{0x3C00}, // 1.0
			[]float64{1.0},
		},
		{
			"Multiple elements",
			[]Float16{0x3C00, 0x4000, 0x4400}, // 1.0, 2.0, 4.0
			[]float64{1.0, 2.0, 4.0},
		},
		{
			"Special values",
			[]Float16{PositiveZero, NegativeZero, PositiveInfinity, NegativeInfinity, NaN()},
			[]float64{0.0, -0.0, math.Inf(1), math.Inf(-1), math.NaN()},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ToSlice64(tt.input)

			if len(got) != len(tt.want) {
				t.Fatalf("ToSlice64() length = %d, want %d", len(got), len(tt.want))
			}

			for i := range got {
				// Special handling for NaN
				if math.IsNaN(tt.want[i]) {
					if !math.IsNaN(got[i]) {
						t.Errorf("ToSlice64()[%d] = %v, want NaN", i, got[i])
					}
					continue
				}

				// For infinity, check sign
				if math.IsInf(tt.want[i], 0) {
					if !math.IsInf(got[i], 0) || math.Signbit(got[i]) != math.Signbit(tt.want[i]) {
						t.Errorf("ToSlice64()[%d] = %v, want %v", i, got[i], tt.want[i])
					}
					continue
				}

				// For zero, check sign
				if tt.want[i] == 0.0 || tt.want[i] == -0.0 {
					if got[i] != 0.0 && got[i] != -0.0 {
						t.Errorf("ToSlice64()[%d] = %v, want %v", i, got[i], tt.want[i])
					}
					continue
				}

				// For other values, allow small floating point differences
				const epsilon = 1e-10
				diff := math.Abs(got[i] - tt.want[i])
				if diff > epsilon {
					t.Errorf("ToSlice64()[%d] = %v, want %v (diff: %e)", i, got[i], tt.want[i], diff)
				}
			}
		})
	}
}

func TestFromSlice64(t *testing.T) {
	tests := []struct {
		name  string
		input []float64
		want  []Float16
	}{
		{
			"Empty slice",
			[]float64{},
			[]Float16{},
		},
		{
			"Single element",
			[]float64{1.0},
			[]Float16{0x3C00}, // 1.0
		},
		{
			"Multiple elements",
			[]float64{1.0, 2.0, 4.0},
			[]Float16{0x3C00, 0x4000, 0x4400}, // 1.0, 2.0, 4.0
		},
		{
			"Special values",
			[]float64{0.0, -0.0, math.Inf(1), math.Inf(-1), math.NaN()},
			[]Float16{PositiveZero, NegativeZero, PositiveInfinity, NegativeInfinity, NaN()},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FromSlice64(tt.input)

			if len(got) != len(tt.want) {
				t.Fatalf("FromSlice64() length = %d, want %d", len(got), len(tt.want))
			}

			for i := range got {
				// Special handling for NaN
				if tt.want[i].IsNaN() {
					if !got[i].IsNaN() {
						t.Errorf("FromSlice64()[%d] = %v (0x%04X), want NaN", i, got[i], uint16(got[i]))
					}
					continue
				}

				// For infinity, check sign
				if tt.want[i].IsInf(0) {
					if !got[i].IsInf(0) || got[i].Signbit() != tt.want[i].Signbit() {
						t.Errorf("FromSlice64()[%d] = %v (0x%04X), want %v (0x%04X)",
							i, got[i], uint16(got[i]), tt.want[i], uint16(tt.want[i]))
					}
					continue
				}

				// For zero, check sign
				if tt.want[i] == 0 || tt.want[i] == 0x8000 {
					if got[i] != 0 && got[i] != 0x8000 {
						t.Errorf("FromSlice64()[%d] = %v (0x%04X), want 0.0 or -0.0",
							i, got[i], uint16(got[i]))
					}
					continue
				}

				// For other values, check exact match
				if got[i] != tt.want[i] {
					t.Errorf("FromSlice64()[%d] = %v (0x%04X), want %v (0x%04X)",
						i, got[i], uint16(got[i]), tt.want[i], uint16(tt.want[i]))
				}
			}
		})
	}
}

func TestToSlice16WithMode(t *testing.T) {
	tests := []struct {
		name      string
		input     []float32
		convMode  ConversionMode
		roundMode RoundingMode
		want      []Float16
		hasError  bool
	}{
		{
			"Empty slice",
			[]float32{},
			ModeIEEE,
			RoundNearestEven,
			[]Float16{},
			false,
		},
		{
			"Single element",
			[]float32{1.0},
			ModeIEEE,
			RoundNearestEven,
			[]Float16{0x3C00}, // 1.0
			false,
		},
		{
			"Multiple elements",
			[]float32{1.0, 2.0, 4.0},
			ModeIEEE,
			RoundNearestEven,
			[]Float16{0x3C00, 0x4000, 0x4400}, // 1.0, 2.0, 4.0
			false,
		},
		{
			"Special values",
			[]float32{0.0, float32(math.Copysign(0, -1)), float32(math.Inf(1)), float32(math.Inf(-1)), float32(math.NaN())},
			ModeIEEE,
			RoundNearestEven,
			[]Float16{PositiveZero, NegativeZero, PositiveInfinity, NegativeInfinity, NaN()},
			false,
		},
		{
			"Strict mode with overflow",
			[]float32{1e10},
			ModeStrict,
			RoundNearestEven,
			nil,
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, errs := ToSlice16WithMode(tt.input, tt.convMode, tt.roundMode)

			if tt.hasError {
				if len(errs) == 0 || errs[0] == nil {
					t.Error("Expected error, got none")
				}
				return
			}

			if len(errs) > 0 && errs[0] != nil {
				t.Fatalf("Unexpected error: %v", errs[0])
			}

			if len(result) != len(tt.want) {
				t.Fatalf("ToSlice16WithMode() length = %d, want %d", len(result), len(tt.want))
			}

			for i := range result {
				// Special handling for NaN
				if tt.want[i].IsNaN() {
					if !result[i].IsNaN() {
						t.Errorf("ToSlice16WithMode()[%d] = %v (0x%04X), want NaN",
							i, result[i], uint16(result[i]))
					}
					continue
				}

				// For infinity, check sign
				if tt.want[i].IsInf(0) {
					if !result[i].IsInf(0) || result[i].Signbit() != tt.want[i].Signbit() {
						t.Errorf("ToSlice16WithMode()[%d] = %v (0x%04X), want %v (0x%04X)",
							i, result[i], uint16(result[i]), tt.want[i], uint16(tt.want[i]))
					}
					continue
				}

				// For zero, check sign
				if tt.want[i] == 0 || tt.want[i] == 0x8000 {
					if result[i] != 0 && result[i] != 0x8000 {
						t.Errorf("ToSlice16WithMode()[%d] = %v (0x%04X), want 0.0 or -0.0",
							i, result[i], uint16(result[i]))
					}
					continue
				}

				// For other values, check exact match
				if result[i] != tt.want[i] {
					t.Errorf("ToSlice16WithMode()[%d] = %v (0x%04X), want %v (0x%04X)",
						i, result[i], uint16(result[i]), tt.want[i], uint16(tt.want[i]))
				}
			}
		})
	}
}
