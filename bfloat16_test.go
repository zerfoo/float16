package float16

import (
	"math"
	"testing"
)

func TestBFloat16FromFloat32(t *testing.T) {
	tests := []struct {
		input    float32
		expected uint16
		desc     string
	}{
		{0.0, 0x0000, "positive zero"},
		{float32(math.Copysign(0.0, -1.0)), 0x8000, "negative zero"},
		{1.0, 0x3F80, "one"},
		{-1.0, 0xBF80, "negative one"},
		{float32(math.Inf(1)), 0x7F80, "positive infinity"},
		{float32(math.Inf(-1)), 0xFF80, "negative infinity"},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			result := BFloat16FromFloat32(test.input)
			if result.Bits() != test.expected {
				t.Errorf("BFloat16FromFloat32(%v) = 0x%04X, expected 0x%04X",
					test.input, result.Bits(), test.expected)
			}
		})
	}
}

func TestBFloat16ToFloat32(t *testing.T) {
	tests := []struct {
		input    uint16
		expected float32
		desc     string
	}{
		{0x0000, 0.0, "positive zero"},
		{0x8000, float32(math.Copysign(0.0, -1.0)), "negative zero"},
		{0x3F80, 1.0, "one"},
		{0xBF80, -1.0, "negative one"},
		{0x7F80, float32(math.Inf(1)), "positive infinity"},
		{0xFF80, float32(math.Inf(-1)), "negative infinity"},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			bf := BFloat16FromBits(test.input)
			result := bf.ToFloat32()

			// Handle special case of negative zero
			if test.input == 0x8000 {
				if math.Signbit(float64(result)) == false {
					t.Errorf("BFloat16(0x%04X).ToFloat32() should be negative zero", test.input)
				}
				return
			}

			if result != test.expected {
				t.Errorf("BFloat16(0x%04X).ToFloat32() = %v, expected %v",
					test.input, result, test.expected)
			}
		})
	}
}

func TestBFloat16Arithmetic(t *testing.T) {
	a := BFloat16FromFloat32(2.0)
	b := BFloat16FromFloat32(3.0)

	// Test addition
	sum := BFloat16Add(a, b)
	if !BFloat16Equal(sum, BFloat16FromFloat32(5.0)) {
		t.Errorf("2.0 + 3.0 should equal 5.0, got %v", sum.ToFloat32())
	}

	// Test multiplication
	prod := BFloat16Mul(a, b)
	if !BFloat16Equal(prod, BFloat16FromFloat32(6.0)) {
		t.Errorf("2.0 * 3.0 should equal 6.0, got %v", prod.ToFloat32())
	}
}

func TestBFloat16Classification(t *testing.T) {
	// Test zero
	zero := BFloat16FromFloat32(0.0)
	if !zero.IsZero() {
		t.Error("0.0 should be identified as zero")
	}
	if !zero.IsFinite() {
		t.Error("0.0 should be finite")
	}

	// Test infinity
	inf := BFloat16PositiveInfinity
	if !inf.IsInf(0) {
		t.Error("positive infinity should be identified as infinity")
	}
	if inf.IsFinite() {
		t.Error("infinity should not be finite")
	}

	// Test normal number
	one := BFloat16FromFloat32(1.0)
	if !one.IsNormal() {
		t.Error("1.0 should be a normal number")
	}
	if !one.IsFinite() {
		t.Error("1.0 should be finite")
	}
}

func TestFloat16BFloat16Conversion(t *testing.T) {
	// Test round-trip conversion
	original := FromFloat32(3.14159)
	asBFloat := original.ToBFloat16()
	backToFloat16 := asBFloat.ToFloat16()

	// Due to different precision, we expect some loss
	// Just verify the conversion functions work without panicking
	if backToFloat16.IsNaN() {
		t.Error("Round-trip conversion should not produce NaN for normal values")
	}

	// Test the reverse direction
	originalBF := BFloat16FromFloat32(2.718)
	asFloat16 := originalBF.ToFloat16()
	backToBFloat := asFloat16.ToBFloat16()

	if backToBFloat.IsNaN() {
		t.Error("Round-trip conversion should not produce NaN for normal values")
	}
}

func TestBFloat16String(t *testing.T) {
	tests := []struct {
		value BFloat16
		desc  string
	}{
		{BFloat16FromFloat32(1.0), "one"},
		{BFloat16FromFloat32(-1.0), "negative one"},
		{BFloat16PositiveInfinity, "positive infinity"},
		{BFloat16NegativeInfinity, "negative infinity"},
		{BFloat16QuietNaN, "quiet NaN"},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			str := test.value.String()
			if str == "" {
				t.Errorf("String representation should not be empty for %s", test.desc)
			}
		})
	}
}

func TestBFloat16Class(t *testing.T) {
	tests := []struct {
		name     string
		input    BFloat16
		expected FloatClass
	}{
		{"PositiveZero", BFloat16PositiveZero, ClassPositiveZero},
		{"NegativeZero", BFloat16NegativeZero, ClassNegativeZero},
		{"PositiveInfinity", BFloat16PositiveInfinity, ClassPositiveInfinity},
		{"NegativeInfinity", BFloat16NegativeInfinity, ClassNegativeInfinity},
		{"QuietNaN", BFloat16QuietNaN, ClassQuietNaN},
		{"SignalingNaN", BFloat16SignalingNaN, ClassSignalingNaN},
		{"PositiveNormal", BFloat16FromFloat32(1.0), ClassPositiveNormal},
		{"NegativeNormal", BFloat16FromFloat32(-1.0), ClassNegativeNormal},
		{"PositiveSubnormal", BFloat16SmallestPosSubnormal, ClassPositiveSubnormal},
		{"NegativeSubnormal", BFloat16SmallestNegSubnormal, ClassNegativeSubnormal},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.input.Class()
			if got != tt.expected {
				t.Errorf("Class() for %s (0x%04x) = %v, want %v", tt.name, uint16(tt.input), got, tt.expected)
			}
		})
	}
}

func TestBFloat16CopySign(t *testing.T) {
	tests := []struct {
		name     string
		f        BFloat16
		s        BFloat16
		expected BFloat16
	}{
		{"PositiveMagnitudePositiveSign", BFloat16FromFloat32(1.0), BFloat16FromFloat32(2.0), BFloat16FromFloat32(1.0)},
		{"PositiveMagnitudeNegativeSign", BFloat16FromFloat32(1.0), BFloat16FromFloat32(-2.0), BFloat16FromFloat32(-1.0)},
		{"NegativeMagnitudePositiveSign", BFloat16FromFloat32(-1.0), BFloat16FromFloat32(2.0), BFloat16FromFloat32(1.0)},
		{"NegativeMagnitudeNegativeSign", BFloat16FromFloat32(-1.0), BFloat16FromFloat32(-2.0), BFloat16FromFloat32(-1.0)},
		{"ZeroMagnitudePositiveSign", BFloat16PositiveZero, BFloat16FromFloat32(2.0), BFloat16PositiveZero},
		{"ZeroMagnitudeNegativeSign", BFloat16PositiveZero, BFloat16FromFloat32(-2.0), BFloat16NegativeZero},
		{"InfMagnitudePositiveSign", BFloat16PositiveInfinity, BFloat16FromFloat32(2.0), BFloat16PositiveInfinity},
		{"InfMagnitudeNegativeSign", BFloat16PositiveInfinity, BFloat16FromFloat32(-2.0), BFloat16NegativeInfinity},
		{"NaNMagnitudePositiveSign", BFloat16QuietNaN, BFloat16FromFloat32(2.0), BFloat16QuietNaN},
		{"NaNMagnitudeNegativeSign", BFloat16QuietNaN, BFloat16FromFloat32(-2.0), BFloat16QuietNaN | BFloat16SignMask},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.f.CopySign(tt.s)
			if got != tt.expected {
				t.Errorf("CopySign(%04x, %04x) = %04x, want %04x", uint16(tt.f), uint16(tt.s), uint16(got), uint16(tt.expected))
			}
		})
	}
}
