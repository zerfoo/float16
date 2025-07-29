package float16

import (
	"testing"
)

func TestTypesExtra(t *testing.T) {
	// Test Error method
	err := &Float16Error{Op: "test", Msg: "test error", Value: 1.0}
	if err.Error() != "float16.test: test error (value: 1)" {
		t.Errorf("Error method: expected 'float16.test: test error (value: 1)', got '%s'", err.Error())
	}
	err = &Float16Error{Op: "test", Msg: "test error"}
	if err.Error() != "float16.test: test error" {
		t.Errorf("Error method: expected 'float16.test: test error', got '%s'", err.Error())
	}

	// Test String method
	if PositiveZero.String() != "0" {
		t.Errorf("String method: expected '0', got '%s'", PositiveZero.String())
	}
	if NegativeZero.String() != "-0" {
		t.Errorf("String method: expected '-0', got '%s'", NegativeZero.String())
	}
	if PositiveInfinity.String() != "+Inf" {
		t.Errorf("String method: expected '+Inf', got '%s'", PositiveInfinity.String())
	}
	if NegativeInfinity.String() != "-Inf" {
		t.Errorf("String method: expected '-Inf', got '%s'", NegativeInfinity.String())
	}
	if QuietNaN.String() != "NaN" {
		t.Errorf("String method: expected 'NaN', got '%s'", QuietNaN.String())
	}
	if NegativeQNaN.String() != "-NaN" {
		t.Errorf("String method: expected '-NaN', got '%s'", NegativeQNaN.String())
	}

	// Test GoString method
	if PositiveZero.GoString() != "float16.FromBits(0x0000)" {
		t.Errorf("GoString method: expected 'float16.FromBits(0x0000)', got '%s'", PositiveZero.GoString())
	}

	// Test Class method
	if SignalingNaN.Class() != ClassSignalingNaN {
		t.Errorf("Class method: expected ClassSignalingNaN, got %v", SignalingNaN.Class())
	}
	if NegativeInfinity.Class() != ClassNegativeInfinity {
		t.Errorf("Class method: expected ClassNegativeInfinity, got %v", NegativeInfinity.Class())
	}
	if NegativeZero.Class() != ClassNegativeZero {
		t.Errorf("Class method: expected ClassNegativeZero, got %v", NegativeZero.Class())
	}
	if FromFloat32(-1).Class() != ClassNegativeNormal {
		t.Errorf("Class method: expected ClassNegativeNormal, got %v", FromFloat32(-1).Class())
	}
	if FromBits(0x8001).Class() != ClassNegativeSubnormal {
		t.Errorf("Class method: expected ClassNegativeSubnormal, got %v", FromBits(0x8001).Class())
	}

	// Test leadingZeros10
	if leadingZeros10(0) != 10 {
		t.Errorf("leadingZeros10(0): expected 10, got %d", leadingZeros10(0))
	}
	if leadingZeros10(1) != 3 {
		t.Errorf("leadingZeros10(1): expected 3, got %d", leadingZeros10(1))
	}
}
