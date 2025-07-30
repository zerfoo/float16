package float16

import "math"

func fromFloat32New(f32 float32) Float16 {
	bits := math.Float32bits(f32)
	sign := uint16(bits >> 31)
	exp := int32((bits >> 23) & 0xff)
	mant := uint32(bits & 0x7fffff)

	// Handle special cases (infinity and NaN)
	if exp == 0xff {
		if mant == 0 {
			return Float16(sign<<15 | 0x7c00) // infinity
		}
		return Float16(sign<<15 | 0x7e00) // qNaN
	}

	// Handle zero
	if exp == 0 && mant == 0 {
		return Float16(sign << 15)
	}

	// Adjust exponent from float32 bias (127) to float16 bias (15)
	exp -= 127 - 15

	// Handle overflow (exponent too large)
	if exp >= 0x1f {
		return Float16(sign<<15 | 0x7c00) // infinity
	}

	// Handle underflow and subnormal numbers
	if exp <= 0 {
		if exp < -10 {
			return Float16(sign << 15) // zero
		}
		// Convert to subnormal
		mant = (mant | 1<<23) >> uint(1-exp)
		// Round to nearest even
		if mant&0x1fff > 0x1000 || (mant&0x1fff == 0x1000 && mant&0x2000 != 0) {
			mant += 0x2000
		}
		return Float16(uint16(sign<<15) | uint16(mant>>13))
	}

	// Handle normal numbers
	// Add implicit 1
	mant |= 1 << 23

	// For float32 to float16, we need to round the 23-bit mantissa to 10 bits
	// We work with the original 23-bit mantissa and round to get 10 bits

	// Round to nearest even
	// Look at bit 12 (guard), bits 11-0 (round/sticky)
	guard := (mant >> 12) & 1
	sticky := mant & 0xFFF
	lsb := (mant >> 13) & 1

	// Round up if: guard=1 AND (sticky!=0 OR lsb=1)
	if guard != 0 && (sticky != 0 || lsb != 0) {
		mant += 1 << 13
	}

	// Check for mantissa overflow after rounding
	if mant >= 1<<24 {
		// Mantissa overflowed, increment exponent
		exp++
		mant = 0 // Reset mantissa to 0 (implicit 1 will be added by IEEE format)
	}

	// Check for exponent overflow after rounding
	if exp >= 0x1f {
		return Float16(sign<<15 | 0x7c00) // infinity
	}

	// Extract the 10-bit mantissa (bits 22-13 of the original 23-bit mantissa)
	mantissa10 := (mant >> 13) & 0x3FF

	return Float16(uint16(sign<<15) | uint16(exp<<10) | uint16(mantissa10))
}
