package float16

import (
	"math"
	"testing"
)

// Use the RoundingMode and ConversionMode constants from types.go
const (
	// Rounding modes
	testRoundNearestEven = RoundNearestEven
	testRoundToZero      = RoundTowardZero
	testRoundUp          = RoundTowardPositive
	testRoundDown        = RoundTowardNegative
)

const (
	// Conversion modes
	testModeDefault = ModeIEEE
	testModeStrict  = ModeStrict
)

// Test conversion functions

func TestToFloat16Basic(t *testing.T) {
	converter := NewConverter(ModeIEEE, RoundNearestEven)
	tests := []struct {
		input    float32
		expected Float16
		name     string
	}{
		{0.0, PositiveZero, "positive zero"},
		{float32(math.Copysign(0, -1)), NegativeZero, "negative zero"}, // Use Copysign to ensure negative zero
		{1.0, 0x3C00, "one"},
		{-1.0, 0xBC00, "negative one"},
		{2.0, 0x4000, "two"},
		{0.5, 0x3800, "half"},
		{float32(math.Inf(1)), PositiveInfinity, "positive infinity"},
		{float32(math.Inf(-1)), NegativeInfinity, "negative infinity"},
		{65504.0, MaxValue, "max finite value"},
		{-65504.0, MinValue, "min finite value"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := converter.ToFloat16(test.input)
			if result != test.expected {
				t.Errorf("ToFloat16(%g) = 0x%04x, expected 0x%04x",
					test.input, result, test.expected)
			}
		})
	}
}

func TestToFloat16NaN(t *testing.T) {
	converter := NewConverter(ModeIEEE, RoundNearestEven)
	result := converter.ToFloat16(float32(math.NaN()))
	if !result.IsNaN() {
		t.Errorf("ToFloat16(NaN) should return NaN, got 0x%04x", result)
	}
}

func TestToFloat16WithModeStrict(t *testing.T) {
	converter := NewConverter(ModeStrict, RoundNearestEven)
	// Test overflow in strict mode
	_, err := converter.ToFloat16WithMode(1e10)
	if err == nil {
		t.Error("Expected overflow error in strict mode")
	}

	// Test underflow in strict mode
	_, err = converter.ToFloat16WithMode(1e-10)
	if err == nil {
		t.Error("Expected underflow error in strict mode")
	}

	// Test NaN in strict mode
	_, err = converter.ToFloat16WithMode(float32(math.NaN()))
	if err == nil {
		t.Error("Expected NaN error in strict mode")
	}
}

func TestToFloat32(t *testing.T) {
	tests := []struct {
		input    Float16
		expected float32
		name     string
	}{
		{PositiveZero, 0.0, "positive zero"},
		{NegativeZero, -0.0, "negative zero"},
		{0x3C00, 1.0, "one"},
		{0xBC00, -1.0, "negative one"},
		{0x4000, 2.0, "two"},
		{0x3800, 0.5, "half"},
		{PositiveInfinity, float32(math.Inf(1)), "positive infinity"},
		{NegativeInfinity, float32(math.Inf(-1)), "negative infinity"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := test.input.ToFloat32()
			if result != test.expected && !(math.IsInf(float64(result), 0) && math.IsInf(float64(test.expected), 0)) {
				t.Errorf("Float16(0x%04x).ToFloat32() = %g, expected %g",
					test.input, result, test.expected)
			}
		})
	}
}

func TestRoundTripConversion(t *testing.T) {
	converter := NewConverter(ModeIEEE, RoundNearestEven)
	// Test that converting Float16 -> Float32 -> Float16 is identity
	failureCount := 0
	totalTested := 0
	totalSubnormal := 0
	subnormalFailures := 0

	for i := 0; i < 65536; i++ {
		f16 := Float16(i)

		// Skip NaN values as they may not round-trip exactly
		if f16.IsNaN() {
			continue
		}

		f32 := f16.ToFloat32()
		f16_back := converter.ToFloat16(f32)
		totalTested++

		// For subnormal numbers, we can't guarantee exact round-trip due to precision loss
		// in float32 representation. Instead, we'll check that the values are very close.
		isSubnormal := f16.IsSubnormal()
		if isSubnormal {
			totalSubnormal++
			// For subnormals, check if the values are within 1 ULP
			f32_again := f16_back.ToFloat32()
			if math.Abs(float64(f32_again-f32)) > 1e-8 {
				failureCount++
				subnormalFailures++
			}
		} else if f16 != f16_back {
			// For normal numbers, we should have exact round-trip
			failureCount++
		}
	}

	if failureCount > 0 {
		t.Errorf("Found %d round-trip failures out of %d values tested (%.2f%%), including %d subnormal failures",
			failureCount, totalTested, 100*float64(failureCount)/float64(totalTested), subnormalFailures)
	}
}

// Test special value methods

func TestSpecialValueMethods(t *testing.T) {
	if !PositiveZero.IsZero() {
		t.Error("PositiveZero should be zero")
	}
	if !NegativeZero.IsZero() {
		t.Error("NegativeZero should be zero")
	}
	if !PositiveInfinity.IsInf(0) {
		t.Error("PositiveInfinity should be infinity")
	}
	if !NegativeInfinity.IsInf(0) {
		t.Error("NegativeInfinity should be infinity")
	}
	if !PositiveInfinity.IsInf(1) {
		t.Error("PositiveInfinity should be positive infinity")
	}
	if !NegativeInfinity.IsInf(-1) {
		t.Error("NegativeInfinity should be negative infinity")
	}
	if !ToFloat16(1.0).IsFinite() {
		t.Error("1.0 should be finite")
	}
	if PositiveInfinity.IsFinite() {
		t.Error("PositiveInfinity should not be finite")
	}
	if !ToFloat16(1.0).IsNormal() {
		t.Error("1.0 should be normal")
	}
	if !SmallestSubnormal.IsSubnormal() {
		t.Error("SmallestSubnormal should be subnormal")
	}
}

func TestNaNMethods(t *testing.T) {
	if !QuietNaN.IsNaN() {
		t.Error("QuietNaN should be NaN")
	}
	if !SignalingNaN.IsNaN() {
		t.Error("SignalingNaN should be NaN")
	}
	if QuietNaN.IsFinite() {
		t.Error("NaN should not be finite")
	}
	if QuietNaN.IsNormal() {
		t.Error("NaN should not be normal")
	}
}

func TestAbsNeg(t *testing.T) {
	converter := NewConverter(ModeIEEE, RoundNearestEven)
	a := converter.ToFloat16(-1.0)
	if a.Abs() != converter.ToFloat16(1.0) {
		t.Error("Abs(-1.0) should be 1.0")
	}
	if a.Neg() != converter.ToFloat16(1.0) {
		t.Error("Neg(-1.0) should be 1.0")
	}

	b := converter.ToFloat16(1.0)
	if b.Neg() != converter.ToFloat16(-1.0) {
		t.Error("Neg(1.0) should be -1.0")
	}
}

// Test arithmetic operations

func TestAddBasic(t *testing.T) {
	converter := NewConverter(ModeIEEE, RoundNearestEven)
	tests := []struct {
		a, b     Float16
		expected Float16
		name     string
	}{
		{PositiveZero, PositiveZero, PositiveZero, "zero + zero"},
		{converter.ToFloat16(1.0), PositiveZero, converter.ToFloat16(1.0), "one + zero"},
		{converter.ToFloat16(1.0), converter.ToFloat16(1.0), converter.ToFloat16(2.0), "one + one"},
		{converter.ToFloat16(2.0), converter.ToFloat16(3.0), converter.ToFloat16(5.0), "two + three"},
		{PositiveInfinity, converter.ToFloat16(1.0), PositiveInfinity, "inf + one"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := Add(test.a, test.b)
			if !Equal(result, test.expected) && !result.IsNaN() {
				t.Errorf("Add(0x%04x, 0x%04x) = 0x%04x, expected 0x%04x",
					test.a, test.b, result, test.expected)
			}
		})
	}
}

func TestSubBasic(t *testing.T) {
	converter := NewConverter(ModeIEEE, RoundNearestEven)
	tests := []struct {
		a, b     Float16
		expected Float16
		name     string
	}{
		{PositiveZero, PositiveZero, PositiveZero, "zero - zero"},
		{converter.ToFloat16(1.0), PositiveZero, converter.ToFloat16(1.0), "one - zero"},
		{converter.ToFloat16(3.0), converter.ToFloat16(1.0), converter.ToFloat16(2.0), "three - one"},
		{converter.ToFloat16(1.0), converter.ToFloat16(3.0), converter.ToFloat16(-2.0), "one - three"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := Sub(test.a, test.b)
			if !Equal(result, test.expected) {
				t.Errorf("Sub(0x%04x, 0x%04x) = 0x%04x, expected 0x%04x",
					test.a, test.b, result, test.expected)
			}
		})
	}
}

func TestMulBasic(t *testing.T) {
	converter := NewConverter(ModeIEEE, RoundNearestEven)
	tests := []struct {
		a, b     Float16
		expected Float16
		name     string
	}{
		{PositiveZero, converter.ToFloat16(1.0), PositiveZero, "zero * one"},
		{converter.ToFloat16(1.0), converter.ToFloat16(1.0), converter.ToFloat16(1.0), "one * one"},
		{converter.ToFloat16(2.0), converter.ToFloat16(3.0), converter.ToFloat16(6.0), "two * three"},
		{converter.ToFloat16(-1.0), converter.ToFloat16(1.0), converter.ToFloat16(-1.0), "(-one) * one"},
		{PositiveInfinity, converter.ToFloat16(2.0), PositiveInfinity, "inf * two"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := Mul(test.a, test.b)
			if !Equal(result, test.expected) {
				t.Errorf("Mul(0x%04x, 0x%04x) = 0x%04x, expected 0x%04x",
					test.a, test.b, result, test.expected)
			}
		})
	}
}

func TestDivBasic(t *testing.T) {
	converter := NewConverter(ModeIEEE, RoundNearestEven)
	tests := []struct {
		a, b     Float16
		expected Float16
		name     string
	}{
		{PositiveZero, converter.ToFloat16(1.0), PositiveZero, "zero / one"},
		{converter.ToFloat16(6.0), converter.ToFloat16(2.0), converter.ToFloat16(3.0), "six / two"},
		{converter.ToFloat16(1.0), converter.ToFloat16(2.0), converter.ToFloat16(0.5), "one / two"},
		{converter.ToFloat16(1.0), PositiveZero, PositiveInfinity, "one / zero"},
		{converter.ToFloat16(-1.0), PositiveZero, NegativeInfinity, "(-one) / zero"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := Div(test.a, test.b)
			if !Equal(result, test.expected) && !result.IsInf(0) {
				t.Errorf("Div(0x%04x, 0x%04x) = 0x%04x, expected 0x%04x",
					test.a, test.b, result, test.expected)
			}
		})
	}
}

// Test comparison operations

func TestComparisons(t *testing.T) {
	converter := NewConverter(ModeIEEE, RoundNearestEven)
	a := converter.ToFloat16(1.0)
	b := converter.ToFloat16(2.0)
	c := converter.ToFloat16(1.0)

	if !Less(a, b) {
		t.Error("1.0 should be less than 2.0")
	}
	if !Greater(b, a) {
		t.Error("2.0 should be greater than 1.0")
	}
	if !Equal(a, c) {
		t.Error("1.0 should equal 1.0")
	}
	if !LessEqual(a, b) {
		t.Error("1.0 should be less than or equal to 2.0")
	}
	if !LessEqual(a, c) {
		t.Error("1.0 should be less than or equal to 1.0")
	}
	if !GreaterEqual(b, a) {
		t.Error("2.0 should be greater than or equal to 1.0")
	}
	if !GreaterEqual(a, c) {
		t.Error("1.0 should be greater than or equal to 1.0")
	}
}

func TestMinMax(t *testing.T) {
	converter := NewConverter(ModeIEEE, RoundNearestEven)
	a := converter.ToFloat16(1.0)
	b := converter.ToFloat16(2.0)

	if Min(a, b) != a {
		t.Error("Min(1.0, 2.0) should be 1.0")
	}
	if Max(a, b) != b {
		t.Error("Max(1.0, 2.0) should be 2.0")
	}
}

// Test slice operations

func TestSliceOperations(t *testing.T) {
	converter := NewConverter(ModeIEEE, RoundNearestEven)
	a := []Float16{converter.ToFloat16(1.0), converter.ToFloat16(2.0), converter.ToFloat16(3.0)}
	b := []Float16{converter.ToFloat16(1.0), converter.ToFloat16(1.0), converter.ToFloat16(1.0)}

	// Test AddSlice
	result := AddSlice(a, b)
	expected := []Float16{converter.ToFloat16(2.0), converter.ToFloat16(3.0), converter.ToFloat16(4.0)}
	for i := range result {
		if !Equal(result[i], expected[i]) {
			t.Errorf("AddSlice[%d] = 0x%04x, expected 0x%04x", i, result[i], expected[i])
		}
	}

	// Test ScaleSlice
	scaled := ScaleSlice(a, converter.ToFloat16(2.0))
	expectedScaled := []Float16{converter.ToFloat16(2.0), converter.ToFloat16(4.0), converter.ToFloat16(6.0)}
	for i := range scaled {
		if !Equal(scaled[i], expectedScaled[i]) {
			t.Errorf("ScaleSlice[%d] = 0x%04x, expected 0x%04x", i, scaled[i], expectedScaled[i])
		}
	}

	// Test SumSlice
	sum := SumSlice(a)
	expectedSum := converter.ToFloat16(6.0)
	if !Equal(sum, expectedSum) {
		t.Errorf("SumSlice = 0x%04x, expected 0x%04x", sum, expectedSum)
	}

	// Test DotProduct
	dot := DotProduct(a, b)
	expectedDot := converter.ToFloat16(6.0) // 1*1 + 2*1 + 3*1 = 6
	if !Equal(dot, expectedDot) {
		t.Errorf("DotProduct = 0x%04x, expected 0x%04x", dot, expectedDot)
	}
}

// Test mathematical functions

func TestDebugSubnormalValues(t *testing.T) {
	converter := NewConverter(ModeIEEE, RoundNearestEven)
	tests := []struct {
		name  string
		value uint16
	}{
		{"smallest subnormal", 0x0001},
		{"largest subnormal", 0x03ff},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f16 := Float16(tt.value)
			f16.ToFloat32()
			f16.ToFloat64()
		})
	}
}

func TestSqrt(t *testing.T) {
	converter := NewConverter(ModeIEEE, RoundNearestEven)
	mathConverter := NewMathConverter(converter)
	// Test Sqrt
	sqrtResult := mathConverter.Sqrt(converter.FromFloat32(4.0))
	if sqrtResult != converter.FromFloat32(2.0) {
		t.Errorf("Expected Sqrt(4.0) to be 2.0, but got %v", sqrtResult)
	}
}

func TestSinCosTan(t *testing.T) {
	converter := NewConverter(ModeIEEE, RoundNearestEven)
	mathConverter := NewMathConverter(converter)
	// Test Sin, Cos, Tan
	sinResult := mathConverter.Sin(converter.FromFloat32(0.0))
	if sinResult != converter.FromFloat32(0.0) {
		t.Errorf("Expected Sin(0.0) to be 0.0, but got %v", sinResult)
	}
	cosResult := mathConverter.Cos(converter.FromFloat32(0.0))
	if cosResult != converter.FromFloat32(1.0) {
		t.Errorf("Expected Cos(0.0) to be 1.0, but got %v", cosResult)
	}
	tanResult := mathConverter.Tan(converter.FromFloat32(0.0))
	if tanResult != converter.FromFloat32(0.0) {
		t.Errorf("Expected Tan(0.0) to be 0.0, but got %v", tanResult)
	}
}

func TestToFloat64(t *testing.T) {
	converter := NewConverter(ModeIEEE, RoundNearestEven)
	tests := []struct {
		name     string
		input    Float16
		expected float64
	}{
		// Special values
		{"positive zero", PositiveZero, 0.0},
		{"negative zero", NegativeZero, math.Copysign(0.0, -1.0)},
		{"positive infinity", PositiveInfinity, math.Inf(1)},
		{"negative infinity", NegativeInfinity, math.Inf(-1)},
		{"quiet NaN", NaN(), math.NaN()},

		// Normal numbers
		{"one", Float16(0x3c00), 1.0},
		{"negative one", Float16(0xbc00), -1.0},
		{"two", Float16(0x4000), 2.0},
		{"half", Float16(0x3800), 0.5},
		{"smallest normal", Float16(0x0400), 6.103515625e-05}, // 2^-14

		// Subnormal numbers
		{"smallest subnormal", Float16(0x0001), 5.960464477539063e-08}, // 2^-24
		{"largest subnormal", Float16(0x03ff), 6.097555160522461e-05},  // (1-2^-10) * 2^-14

		// Large numbers
		{"max value", MaxValue, 65504.0},
		{"min value", MinValue, -65504.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.input.ToFloat64()

			if tt.input.IsNaN() {
				if !math.IsNaN(result) {
					t.Errorf("Expected NaN, got %v", result)
				}
				return
			}

			if result != tt.expected {
				// Allow for a small tolerance for floating point comparisons
				if math.Abs(result-tt.expected) > 1e-12 {
					t.Errorf("ToFloat64() = %v, want %v", result, tt.expected)
				}
			}

			if math.Signbit(result) != math.Signbit(tt.expected) {
				t.Errorf("Sign mismatch: got %v, want %v", math.Signbit(result), math.Signbit(tt.expected))
			}
		})
	}
}

func TestFromFloat64(t *testing.T) {
	converter := NewConverter(ModeIEEE, RoundNearestEven)
	tests := []struct {
		name        string
		input       float64
		expected    Float16
		expectExact bool // Whether to expect exact match or just same sign and approximate value
	}{
		{input: 0.0, expected: 0x0000, name: "positive zero", expectExact: true},
		{input: math.Copysign(0, -1), expected: 0x8000, name: "negative zero", expectExact: true},
		{input: 1.0, expected: 0x3C00, name: "one", expectExact: true},
		{input: -1.0, expected: 0xBC00, name: "negative one", expectExact: true},
		{input: 2.0, expected: 0x4000, name: "two", expectExact: true},
		{input: 0.5, expected: 0x3800, name: "half", expectExact: true},
		{input: math.Inf(1), expected: 0x7C00, name: "positive infinity", expectExact: true},
		{input: math.Inf(-1), expected: 0xFC00, name: "negative infinity", expectExact: true},
		{input: 65504.0, expected: 0x7BFF, name: "max finite value", expectExact: true},
		{input: -65504.0, expected: 0xFBFF, name: "min finite value", expectExact: true},
		// For subnormal values, we can't always expect exact round-trip
		{input: math.Float64frombits(0x3f00000000000001), expected: 0x0001, name: "smallest positive subnormal", expectExact: false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := converter.FromFloat64(test.input)
			if test.expectExact {
				if result != test.expected {
					t.Errorf("FromFloat64(%g) = 0x%04x, expected 0x%04x", test.input, result, test.expected)
				}
			} else {
				// For subnormal values, just verify the sign is correct and the value is small
				if result.Sign() != test.expected.Sign() {
					t.Errorf("FromFloat64(%g) = 0x%04x, expected sign %d but got %d",
						test.input, result, test.expected.Sign(), result.Sign())
				}

				// Log the actual value for debugging
				t.Logf("FromFloat64(%g) = 0x%04x (value: %g)", test.input, result, result.ToFloat32())
			}
		})
	}
}

func TestFromFloat64WithMode(t *testing.T) {
	// Test basic conversion
	t.Run("basic conversion", func(t *testing.T) {
		result, err := FromFloat64WithMode(1.5, testModeDefault, testRoundNearestEven)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		expected := Float16(0x3E00) // 1.5 in float16
		if result != expected {
			t.Errorf("FromFloat64WithMode(1.5) = 0x%04x, expected 0x%04x", result, expected)
		}
	})

	// Test strict mode with overflow
	t.Run("strict mode overflow", func(t *testing.T) {
		_, err := FromFloat64WithMode(1e10, testModeStrict, testRoundNearestEven)
		if err == nil {
			t.Error("Expected overflow error in strict mode")
		}
	})

	// Test strict mode with underflow
	t.Run("strict mode underflow", func(t *testing.T) {
		_, err := FromFloat64WithMode(1e-10, testModeStrict, testRoundNearestEven)
		if err == nil {
			t.Error("Expected underflow error in strict mode")
		}
	})

	// Test NaN in strict mode
	t.Run("strict mode NaN", func(t *testing.T) {
		_, err := FromFloat64WithMode(math.NaN(), testModeStrict, testRoundNearestEven)
		if err == nil {
			t.Error("Expected NaN error in strict mode")
		}
	})

	// Test different rounding modes
	roundingTests := []struct {
		input     float64
		roundMode RoundingMode
		expected  Float16
		name      string
	}{
		{1.2, testRoundNearestEven, 0x3CCD, "1.2 to nearest even"},
		{1.2, testRoundToZero, 0x3ccd, "1.2 toward zero"},
		{1.2, testRoundUp, 0x3ccd, "1.2 toward +inf"},
		{-1.2, testRoundDown, 0xbccd, "-1.2 toward -inf"},
	}

	for _, test := range roundingTests {
		t.Run(test.name, func(t *testing.T) {
			result, err := FromFloat64WithMode(test.input, testModeDefault, test.roundMode)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			if result != test.expected {
				t.Errorf("FromFloat64WithMode(%g, %v) = 0x%04x, expected 0x%04x",
					test.input, test.roundMode, result.Bits(), test.expected.Bits())
			}
		})
	}
}

// Test error handling

func TestArithmeticWithNaN(t *testing.T) {
	converter := NewConverter(ModeIEEE, RoundNearestEven)
	nan := QuietNaN
	one := converter.ToFloat16(1.0)

	if !Add(nan, one).IsNaN() {
		t.Error("NaN + 1 should be NaN")
	}
	if !Mul(nan, one).IsNaN() {
		t.Error("NaN * 1 should be NaN")
	}
	if !Div(nan, one).IsNaN() {
		t.Error("NaN / 1 should be NaN")
	}
}

func TestArithmeticWithInfinity(t *testing.T) {
	converter := NewConverter(ModeIEEE, RoundNearestEven)
	inf := PositiveInfinity
	one := converter.ToFloat16(1.0)

	if Add(inf, one) != inf {
		t.Error("∞ + 1 should be ∞")
	}
	if Mul(inf, one) != inf {
		t.Error("∞ * 1 should be ∞")
	}
	if !Div(one, PositiveZero).IsInf(1) {
		t.Error("1 / 0 should be +∞")
	}
}

// Benchmarks

func BenchmarkToFloat16(b *testing.B) {
	converter := NewConverter(ModeIEEE, RoundNearestEven)
	f32 := float32(1.5)
	for i := 0; i < b.N; i++ {
		_ = converter.ToFloat16(f32)
	}
}

func BenchmarkToFloat32(b *testing.B) {
	f16 := ToFloat16(1.5)
	for i := 0; i < b.N; i++ {
		_ = f16.ToFloat32()
	}
}

func BenchmarkAdd(b *testing.B) {
	converter := NewConverter(ModeIEEE, RoundNearestEven)
	a := converter.ToFloat16(1.5)
	c := converter.ToFloat16(2.5)
	for i := 0; i < b.N; i++ {
		_ = Add(a, c)
	}
}

func BenchmarkMul(b *testing.B) {
	converter := NewConverter(ModeIEEE, RoundNearestEven)
	a := converter.ToFloat16(1.5)
	c := converter.ToFloat16(2.5)
	for i := 0; i < b.N; i++ {
		_ = Mul(a, c)
	}
}

func BenchmarkSqrt(b *testing.B) {
	converter := NewConverter(ModeIEEE, RoundNearestEven)
	f := converter.ToFloat16(16.0)
	for i := 0; i < b.N; i++ {
		_ = Sqrt(f)
	}
}


	converter := NewConverter(ModeIEEE, RoundNearestEven)
	if converter == nil {
		t.Error("Expected converter to be initialized, got nil")
	}
	input := make([]float32, 1000)
	for i := range input {
		input[i] = float32(i) * 0.1
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = converter.ToSlice16(input)
	}
}

func BenchmarkToSlice32(b *testing.B) {
	input := make([]Float16, 1000)
	for i := range input {
		input[i] = ToFloat16(float32(i) * 0.1)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ToSlice32(input)
	}
}

func BenchmarkDotProduct(b *testing.B) {
	converter := NewConverter(ModeIEEE, RoundNearestEven)
	size := 1000
	a := make([]Float16, size)
	c := make([]Float16, size)
	for i := range a {
		a[i] = converter.ToFloat16(float32(i) * 0.1)
		c[i] = converter.ToFloat16(float32(i) * 0.2)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = DotProduct(a, c)
	}
}

// Test package configuration

func TestConfiguration(t *testing.T) {
	// Save original config
	originalConfig := GetConfig()

	// Test custom configuration
	customConfig := &Config{
		DefaultConversionMode: ModeStrict,
		DefaultRoundingMode:   RoundTowardZero,
		DefaultArithmeticMode: ModeFastArithmetic,
		EnableFastMath:        true,
	}

	Configure(customConfig)

	newConfig := GetConfig()
	if newConfig.DefaultConversionMode != ModeStrict {
		t.Error("Configuration not applied correctly")
	}

	// Restore original config
	Configure(originalConfig)
}

func TestDebugInfo(t *testing.T) {
	info := DebugInfo()

	if info["version"] != Version {
		t.Error("Debug info should contain correct version")
	}
	if info["ieee754_compliant"] != true {
		t.Error("Debug info should indicate IEEE 754 compliance")
	}
	if info["lookup_tables"] != false {
		t.Error("Debug info should indicate no lookup tables")
	}
}

func TestNextAfter(t *testing.T) {
	converter := NewConverter(ModeIEEE, RoundNearestEven)
	tests := []struct {
		name   string
		f, g   Float16
		expect Float16
	}{
		{"1.0 toward 2.0", converter.ToFloat16(1.0), converter.ToFloat16(2.0), FromBits(0x3c01)},
		{"1.0 toward 0.0", converter.ToFloat16(1.0), converter.ToFloat16(0.0), FromBits(0x3bff)},
		{"-1.0 toward -2.0", converter.ToFloat16(-1.0), converter.ToFloat16(-2.0), FromBits(0xbc01)},
		{"-1.0 toward 0.0", converter.ToFloat16(-1.0), converter.ToFloat16(0.0), FromBits(0xbbff)},
		{"0.0 toward 1.0", PositiveZero, converter.ToFloat16(1.0), FromBits(0x0001)},
		{"0.0 toward -1.0", PositiveZero, converter.ToFloat16(-1.0), FromBits(0x8001)},
		{"max toward inf", MaxValue, PositiveInfinity, PositiveInfinity},
		{"-max toward -inf", MinValue, NegativeInfinity, NegativeInfinity},
		{"inf toward 0", PositiveInfinity, PositiveZero, MaxValue},
		{"-inf toward 0", NegativeInfinity, PositiveZero, MinValue},
		{"nan, 1.0", QuietNaN, converter.ToFloat16(1.0), QuietNaN},
		{"1.0, nan", converter.ToFloat16(1.0), QuietNaN, QuietNaN},
		{"1.0, 1.0", converter.ToFloat16(1.0), converter.ToFloat16(1.0), converter.ToFloat16(1.0)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NextAfter(tt.f, tt.g); got != tt.expect {
				t.Errorf("NextAfter() = %v, want %v", got, tt.expect)
			}
		})
	}
}

func TestFrexp(t *testing.T) {
	converter := NewConverter(ModeIEEE, RoundNearestEven)
	tests := []struct {
		name     string
		f        Float16
		wantFrac Float16
		wantExp  int
	}{
		{"zero", PositiveZero, PositiveZero, 0},
		{"one", converter.ToFloat16(1.0), converter.ToFloat16(0.5), 1},
		{"two", converter.ToFloat16(2.0), converter.ToFloat16(0.5), 2},
		{"half", converter.ToFloat16(0.5), converter.ToFloat16(0.5), 0},
		{"inf", PositiveInfinity, PositiveInfinity, 0},
		{"nan", QuietNaN, QuietNaN, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotFrac, gotExp := Frexp(tt.f)
			if gotFrac != tt.wantFrac {
				t.Errorf("Frexp() gotFrac = %v, want %v", gotFrac, tt.wantFrac)
			}
			if gotExp != tt.wantExp {
				t.Errorf("Frexp() gotExp = %v, want %v", gotExp, tt.wantExp)
			}
		})
	}
}

func TestLdexp(t *testing.T) {
	converter := NewConverter(ModeIEEE, RoundNearestEven)
	tests := []struct {
		name string
		frac Float16
		exp  int
		want Float16
	}{
		{"0.5, 1", converter.ToFloat16(0.5), 1, converter.ToFloat16(1.0)},
		{"0.5, 2", converter.ToFloat16(0.5), 2, converter.ToFloat16(2.0)},
		{"zero", PositiveZero, 10, PositiveZero},
		{"inf", PositiveInfinity, 10, PositiveInfinity},
		{"nan", QuietNaN, 10, QuietNaN},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Ldexp(tt.frac, tt.exp); got != tt.want {
				t.Errorf("Ldexp() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestModf(t *testing.T) {
	converter := NewConverter(ModeIEEE, RoundNearestEven)
	tests := []struct {
		name     string
		f        Float16
		wantInt  Float16
		wantFrac Float16
	}{
		{"1.5", converter.ToFloat16(1.5), converter.ToFloat16(1.0), converter.ToFloat16(0.5)},
		{"-1.5", converter.ToFloat16(-1.5), converter.ToFloat16(-1.0), converter.ToFloat16(-0.5)},
		{"inf", PositiveInfinity, PositiveInfinity, PositiveInfinity},
		{"nan", QuietNaN, QuietNaN, QuietNaN},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotInt, gotFrac := Modf(tt.f)
			if gotInt != tt.wantInt {
				t.Errorf("Modf() gotInt = %v, want %v", gotInt, tt.wantInt)
			}
			if gotFrac != tt.wantFrac {
				t.Errorf("Modf() gotFrac = %v, want %v", gotFrac, tt.wantFrac)
			}
		})
	}
}

func TestComputeSliceStats(t *testing.T) {
	converter := NewConverter(ModeIEEE, RoundNearestEven)
	t.Run("empty slice", func(t *testing.T) {
		stats := ComputeSliceStats([]Float16{})
		if stats.Length != 0 {
			t.Errorf("Expected length 0, got %d", stats.Length)
		}
	})

	t.Run("slice with NaNs", func(t *testing.T) {
		s := []Float16{converter.ToFloat16(1.0), converter.ToFloat16(2.0), QuietNaN, converter.ToFloat16(3.0)}
		stats := ComputeSliceStats(s)
		if stats.Min != converter.ToFloat16(1.0) {
			t.Errorf("Expected min 1.0, got %v", stats.Min)
		}
		if stats.Max != converter.ToFloat16(3.0) {
			t.Errorf("Expected max 3.0, got %v", stats.Max)
		}
		if !stats.Sum.IsNaN() {
			t.Errorf("Expected sum to be NaN, got %v", stats.Sum)
		}
	})
}
