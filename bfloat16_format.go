package float16

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"strconv"
)

// BFloat16FromString parses a string into a BFloat16 value.
// It handles special values (NaN, Inf) and numeric strings.
func BFloat16FromString(s string) (BFloat16, error) {
	switch s {
	case "NaN":
		return BFloat16QuietNaN, nil
	case "+Inf", "Inf":
		return BFloat16PositiveInfinity, nil
	case "-Inf":
		return BFloat16NegativeInfinity, nil
	case "+0", "0":
		return BFloat16PositiveZero, nil
	case "-0":
		return BFloat16NegativeZero, nil
	}

	f64, err := strconv.ParseFloat(s, 32)
	if err != nil {
		return 0, err
	}
	return BFloat16FromFloat32(float32(f64)), nil
}

// Format implements fmt.Formatter, supporting %e, %f, %g, %E, %F, %G, %v, and %s verbs.
func (b BFloat16) Format(s fmt.State, verb rune) {
	switch verb {
	case 'e', 'E', 'f', 'F', 'g', 'G':
		// Build a format string matching the state flags
		format := "%"
		if s.Flag('+') {
			format += "+"
		}
		if s.Flag('-') {
			format += "-"
		}
		if s.Flag(' ') {
			format += " "
		}
		if s.Flag('0') {
			format += "0"
		}
		if w, ok := s.Width(); ok {
			format += strconv.Itoa(w)
		}
		if p, ok := s.Precision(); ok {
			format += "." + strconv.Itoa(p)
		}
		format += string(verb)
		fmt.Fprintf(s, format, b.ToFloat32())
	case 'v':
		if s.Flag('#') {
			fmt.Fprint(s, b.GoString())
		} else {
			fmt.Fprint(s, b.String())
		}
	case 's':
		fmt.Fprint(s, b.String())
	default:
		fmt.Fprintf(s, "%%!%c(bfloat16=%s)", verb, b.String())
	}
}

// GoString returns a Go syntax representation of the BFloat16 value.
func (b BFloat16) GoString() string {
	return fmt.Sprintf("float16.BFloat16FromBits(0x%04x)", uint16(b))
}

// MarshalJSON implements json.Marshaler.
func (b BFloat16) MarshalJSON() ([]byte, error) {
	if b.IsNaN() {
		return json.Marshal("NaN")
	}
	if b.IsInf(1) {
		return json.Marshal("+Inf")
	}
	if b.IsInf(-1) {
		return json.Marshal("-Inf")
	}
	return json.Marshal(b.ToFloat32())
}

// UnmarshalJSON implements json.Unmarshaler.
func (b *BFloat16) UnmarshalJSON(data []byte) error {
	// Try as string first (for NaN, Inf)
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		v, err := BFloat16FromString(s)
		if err != nil {
			return fmt.Errorf("float16 BFloat16.UnmarshalJSON: invalid string %q", s)
		}
		*b = v
		return nil
	}

	// Try as number
	var f float32
	if err := json.Unmarshal(data, &f); err != nil {
		return fmt.Errorf("float16 BFloat16.UnmarshalJSON: %w", err)
	}
	*b = BFloat16FromFloat32(f)
	return nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
// The encoding is 2 bytes in little-endian order.
func (b BFloat16) MarshalBinary() ([]byte, error) {
	buf := make([]byte, 2)
	binary.LittleEndian.PutUint16(buf, uint16(b))
	return buf, nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
// The encoding is 2 bytes in little-endian order.
func (b *BFloat16) UnmarshalBinary(data []byte) error {
	if len(data) != 2 {
		return fmt.Errorf("float16 BFloat16.UnmarshalBinary: expected 2 bytes, got %d", len(data))
	}
	*b = BFloat16(binary.LittleEndian.Uint16(data))
	return nil
}
