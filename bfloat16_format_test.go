package float16

import (
	"encoding/json"
	"fmt"
	"math"
	"testing"
)

func TestBFloat16FromString(t *testing.T) {
	tests := []struct {
		input   string
		want    BFloat16
		wantErr bool
	}{
		{"NaN", BFloat16QuietNaN, false},
		{"+Inf", BFloat16PositiveInfinity, false},
		{"-Inf", BFloat16NegativeInfinity, false},
		{"Inf", BFloat16PositiveInfinity, false},
		{"0", BFloat16PositiveZero, false},
		{"+0", BFloat16PositiveZero, false},
		{"-0", BFloat16NegativeZero, false},
		{"1.0", BFloat16FromFloat32(1.0), false},
		{"-1.0", BFloat16FromFloat32(-1.0), false},
		{"3.14", BFloat16FromFloat32(3.14), false},
		{"1e10", BFloat16FromFloat32(1e10), false},
		{"not_a_number", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := BFloat16FromString(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("BFloat16FromString(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
			if !tt.wantErr {
				if tt.want.IsNaN() {
					if !got.IsNaN() {
						t.Errorf("BFloat16FromString(%q) = %v, want NaN", tt.input, got)
					}
				} else if got != tt.want {
					t.Errorf("BFloat16FromString(%q) = 0x%04x, want 0x%04x", tt.input, got, tt.want)
				}
			}
		})
	}
}

func TestBFloat16FromStringRoundTrip(t *testing.T) {
	values := []float32{0, 1, -1, 0.5, 3.14, 100, -256, 1e-5, 1e10}
	for _, f := range values {
		b := BFloat16FromFloat32(f)
		s := b.String()
		got, err := BFloat16FromString(s)
		if err != nil {
			t.Fatalf("BFloat16FromString(%q) error: %v (original float32: %v)", s, err, f)
		}
		if got != b {
			t.Errorf("round-trip failed for %v: String()=%q, parsed=0x%04x, want=0x%04x", f, s, got, b)
		}
	}
}

func TestBFloat16FormatVerbs(t *testing.T) {
	b := BFloat16FromFloat32(3.140625) // BFloat16 representable value close to pi

	// Compare format output against float32 reference
	f32 := b.ToFloat32()

	verbs := []string{"%e", "%E", "%f", "%F", "%g", "%G"}
	for _, verb := range verbs {
		got := fmt.Sprintf(verb, b)
		want := fmt.Sprintf(verb, f32)
		if got != want {
			t.Errorf("Sprintf(%q, bfloat16) = %q, want %q (float32 reference)", verb, got, want)
		}
	}

	// Test with precision
	got := fmt.Sprintf("%.2f", b)
	want := fmt.Sprintf("%.2f", f32)
	if got != want {
		t.Errorf("Sprintf(%%.2f) = %q, want %q", got, want)
	}

	// Test with width
	got = fmt.Sprintf("%10.3e", b)
	want = fmt.Sprintf("%10.3e", f32)
	if got != want {
		t.Errorf("Sprintf(%%10.3e) = %q, want %q", got, want)
	}

	// Test %v and %s use String()
	got = fmt.Sprintf("%v", b)
	want = b.String()
	if got != want {
		t.Errorf("Sprintf(%%v) = %q, want %q", got, want)
	}

	got = fmt.Sprintf("%s", b)
	want = b.String()
	if got != want {
		t.Errorf("Sprintf(%%s) = %q, want %q", got, want)
	}
}

func TestBFloat16FormatSpecialValues(t *testing.T) {
	tests := []struct {
		name string
		val  BFloat16
	}{
		{"zero", BFloat16PositiveZero},
		{"neg_zero", BFloat16NegativeZero},
		{"inf", BFloat16PositiveInfinity},
		{"neg_inf", BFloat16NegativeInfinity},
		{"nan", BFloat16QuietNaN},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f32 := tt.val.ToFloat32()
			for _, verb := range []string{"%e", "%f", "%g"} {
				got := fmt.Sprintf(verb, tt.val)
				want := fmt.Sprintf(verb, f32)
				if got != want {
					t.Errorf("Sprintf(%q, %s) = %q, want %q", verb, tt.name, got, want)
				}
			}
		})
	}
}

func TestBFloat16GoString(t *testing.T) {
	b := BFloat16FromFloat32(1.0)
	got := fmt.Sprintf("%#v", b)
	want := fmt.Sprintf("float16.BFloat16FromBits(0x%04x)", uint16(b))
	if got != want {
		t.Errorf("GoString() = %q, want %q", got, want)
	}
}

func TestBFloat16MarshalJSON(t *testing.T) {
	tests := []struct {
		name string
		val  BFloat16
		want string
	}{
		{"zero", BFloat16PositiveZero, "0"},
		{"one", BFloat16FromFloat32(1.0), "1"},
		{"pi", BFloat16FromFloat32(3.140625), "3.140625"},
		{"negative", BFloat16FromFloat32(-2.0), "-2"},
		{"nan", BFloat16QuietNaN, `"NaN"`},
		{"+inf", BFloat16PositiveInfinity, `"+Inf"`},
		{"-inf", BFloat16NegativeInfinity, `"-Inf"`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.val)
			if err != nil {
				t.Fatalf("MarshalJSON error: %v", err)
			}
			if string(data) != tt.want {
				t.Errorf("MarshalJSON() = %s, want %s", data, tt.want)
			}
		})
	}
}

func TestBFloat16UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    BFloat16
		isNaN   bool
		wantErr bool
	}{
		{"zero", "0", BFloat16PositiveZero, false, false},
		{"one", "1", BFloat16FromFloat32(1.0), false, false},
		{"negative", "-2", BFloat16FromFloat32(-2.0), false, false},
		{"float", "3.14", BFloat16FromFloat32(3.14), false, false},
		{"nan_string", `"NaN"`, BFloat16QuietNaN, true, false},
		{"inf_string", `"+Inf"`, BFloat16PositiveInfinity, false, false},
		{"neg_inf_string", `"-Inf"`, BFloat16NegativeInfinity, false, false},
		{"invalid_string", `"hello"`, 0, false, true},
		{"invalid_json", `{bad}`, 0, false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got BFloat16
			err := json.Unmarshal([]byte(tt.input), &got)
			if (err != nil) != tt.wantErr {
				t.Fatalf("UnmarshalJSON(%s) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
			if !tt.wantErr {
				if tt.isNaN {
					if !got.IsNaN() {
						t.Errorf("UnmarshalJSON(%s) = %v, want NaN", tt.input, got)
					}
				} else if got != tt.want {
					t.Errorf("UnmarshalJSON(%s) = 0x%04x, want 0x%04x", tt.input, got, tt.want)
				}
			}
		})
	}
}

func TestBFloat16JSONRoundTrip(t *testing.T) {
	values := []BFloat16{
		BFloat16PositiveZero,
		BFloat16NegativeZero,
		BFloat16FromFloat32(1.0),
		BFloat16FromFloat32(-1.0),
		BFloat16FromFloat32(3.14),
		BFloat16FromFloat32(0.001),
		BFloat16MaxValue,
		BFloat16PositiveInfinity,
		BFloat16NegativeInfinity,
		BFloat16QuietNaN,
	}

	for _, orig := range values {
		data, err := json.Marshal(orig)
		if err != nil {
			t.Fatalf("MarshalJSON(0x%04x) error: %v", orig, err)
		}

		var got BFloat16
		if err := json.Unmarshal(data, &got); err != nil {
			t.Fatalf("UnmarshalJSON(%s) error: %v", data, err)
		}

		if orig.IsNaN() {
			if !got.IsNaN() {
				t.Errorf("JSON round-trip NaN: got %v, want NaN", got)
			}
		} else if orig != got {
			t.Errorf("JSON round-trip 0x%04x: marshal=%s, got=0x%04x", orig, data, got)
		}
	}
}

func TestBFloat16MarshalBinary(t *testing.T) {
	tests := []struct {
		name string
		val  BFloat16
	}{
		{"zero", BFloat16PositiveZero},
		{"neg_zero", BFloat16NegativeZero},
		{"one", BFloat16FromFloat32(1.0)},
		{"pi", BFloat16FromFloat32(float32(math.Pi))},
		{"max", BFloat16MaxValue},
		{"+inf", BFloat16PositiveInfinity},
		{"-inf", BFloat16NegativeInfinity},
		{"nan", BFloat16QuietNaN},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := tt.val.MarshalBinary()
			if err != nil {
				t.Fatalf("MarshalBinary error: %v", err)
			}
			if len(data) != 2 {
				t.Fatalf("MarshalBinary len = %d, want 2", len(data))
			}

			var got BFloat16
			if err := got.UnmarshalBinary(data); err != nil {
				t.Fatalf("UnmarshalBinary error: %v", err)
			}
			if got != tt.val {
				t.Errorf("binary round-trip: got 0x%04x, want 0x%04x", got, tt.val)
			}
		})
	}
}

func TestBFloat16UnmarshalBinaryErrors(t *testing.T) {
	var b BFloat16

	// Wrong length
	if err := b.UnmarshalBinary([]byte{0x00}); err == nil {
		t.Error("UnmarshalBinary(1 byte) should return error")
	}
	if err := b.UnmarshalBinary([]byte{0x00, 0x00, 0x00}); err == nil {
		t.Error("UnmarshalBinary(3 bytes) should return error")
	}
	if err := b.UnmarshalBinary(nil); err == nil {
		t.Error("UnmarshalBinary(nil) should return error")
	}
}

func TestBFloat16JSONInStruct(t *testing.T) {
	type Weights struct {
		W1 BFloat16 `json:"w1"`
		W2 BFloat16 `json:"w2"`
	}

	orig := Weights{
		W1: BFloat16FromFloat32(0.5),
		W2: BFloat16FromFloat32(-1.25),
	}

	data, err := json.Marshal(orig)
	if err != nil {
		t.Fatalf("Marshal struct error: %v", err)
	}

	var got Weights
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("Unmarshal struct error: %v", err)
	}

	if got.W1 != orig.W1 || got.W2 != orig.W2 {
		t.Errorf("struct round-trip: got %+v, want %+v", got, orig)
	}
}
