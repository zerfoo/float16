package float16

import (
	"math"
)

// Mathematical functions for Float16

// Sqrt returns the square root of the Float16 value
func Sqrt(f Float16) Float16 {
	// Handle special cases
	if f.IsZero() {
		return f // Preserve sign of zero
	}
	if f.IsNaN() {
		return f
	}
	if f.IsInf(1) {
		return PositiveInfinity
	}
	if f.Signbit() {
		// Square root of negative number
		return QuietNaN
	}
	
	// Use float32 for computation and convert back
	f32 := f.ToFloat32()
	result := float32(math.Sqrt(float64(f32)))
	return ToFloat16(result)
}

// Cbrt returns the cube root of the Float16 value
func Cbrt(f Float16) Float16 {
	if f.IsZero() || f.IsNaN() {
		return f
	}
	if f.IsInf(0) {
		return f
	}
	
	f32 := f.ToFloat32()
	result := float32(math.Cbrt(float64(f32)))
	return ToFloat16(result)
}

// Pow returns f raised to the power of exp
func Pow(f, exp Float16) Float16 {
	// Handle special cases according to IEEE 754
	if exp.IsZero() {
		return FromInt(1) // x^0 = 1 for any x (including NaN)
	}
	if f.IsZero() {
		if exp.Signbit() {
			return PositiveInfinity // 0^(-y) = +∞
		}
		return PositiveZero // 0^y = 0 for positive y
	}
	if f.IsNaN() || exp.IsNaN() {
		return QuietNaN
	}
	if f.IsInf(0) {
		if exp.Signbit() {
			return PositiveZero // ∞^(-y) = 0
		}
		return PositiveInfinity // ∞^y = ∞
	}
	
	f32 := f.ToFloat32()
	exp32 := exp.ToFloat32()
	result := float32(math.Pow(float64(f32), float64(exp32)))
	return ToFloat16(result)
}

// Exp returns e^f
func Exp(f Float16) Float16 {
	if f.IsZero() {
		return FromInt(1)
	}
	if f.IsNaN() {
		return f
	}
	if f.IsInf(1) {
		return PositiveInfinity
	}
	if f.IsInf(-1) {
		return PositiveZero
	}
	
	f32 := f.ToFloat32()
	result := float32(math.Exp(float64(f32)))
	return ToFloat16(result)
}

// Exp2 returns 2^f
func Exp2(f Float16) Float16 {
	if f.IsZero() {
		return FromInt(1)
	}
	if f.IsNaN() {
		return f
	}
	if f.IsInf(1) {
		return PositiveInfinity
	}
	if f.IsInf(-1) {
		return PositiveZero
	}
	
	f32 := f.ToFloat32()
	result := float32(math.Exp2(float64(f32)))
	return ToFloat16(result)
}

// Exp10 returns 10^f
func Exp10(f Float16) Float16 {
	return Pow(FromInt(10), f)
}

// Log returns the natural logarithm of f
func Log(f Float16) Float16 {
	if f.IsZero() {
		return NegativeInfinity
	}
	if f.IsNaN() {
		return f
	}
	if f.IsInf(1) {
		return PositiveInfinity
	}
	if f.Signbit() {
		return QuietNaN // log of negative number
	}
	
	f32 := f.ToFloat32()
	result := float32(math.Log(float64(f32)))
	return ToFloat16(result)
}

// Log2 returns the base-2 logarithm of f
func Log2(f Float16) Float16 {
	if f.IsZero() {
		return NegativeInfinity
	}
	if f.IsNaN() {
		return f
	}
	if f.IsInf(1) {
		return PositiveInfinity
	}
	if f.Signbit() {
		return QuietNaN
	}
	
	f32 := f.ToFloat32()
	result := float32(math.Log2(float64(f32)))
	return ToFloat16(result)
}

// Log10 returns the base-10 logarithm of f
func Log10(f Float16) Float16 {
	if f.IsZero() {
		return NegativeInfinity
	}
	if f.IsNaN() {
		return f
	}
	if f.IsInf(1) {
		return PositiveInfinity
	}
	if f.Signbit() {
		return QuietNaN
	}
	
	f32 := f.ToFloat32()
	result := float32(math.Log10(float64(f32)))
	return ToFloat16(result)
}

// Trigonometric functions

// Sin returns the sine of f (in radians)
func Sin(f Float16) Float16 {
	if f.IsZero() {
		return f // Preserve sign of zero
	}
	if f.IsNaN() || f.IsInf(0) {
		return QuietNaN
	}
	
	f32 := f.ToFloat32()
	result := float32(math.Sin(float64(f32)))
	return ToFloat16(result)
}

// Cos returns the cosine of f (in radians)
func Cos(f Float16) Float16 {
	if f.IsZero() {
		return FromInt(1)
	}
	if f.IsNaN() || f.IsInf(0) {
		return QuietNaN
	}
	
	f32 := f.ToFloat32()
	result := float32(math.Cos(float64(f32)))
	return ToFloat16(result)
}

// Tan returns the tangent of f (in radians)
func Tan(f Float16) Float16 {
	if f.IsZero() {
		return f // Preserve sign of zero
	}
	if f.IsNaN() || f.IsInf(0) {
		return QuietNaN
	}
	
	f32 := f.ToFloat32()
	result := float32(math.Tan(float64(f32)))
	return ToFloat16(result)
}

// Asin returns the arcsine of f
func Asin(f Float16) Float16 {
	if f.IsZero() {
		return f
	}
	if f.IsNaN() {
		return f
	}
	
	// Check domain: [-1, 1]
	if f.Abs().ToFloat32() > 1.0 {
		return QuietNaN
	}
	
	f32 := f.ToFloat32()
	result := float32(math.Asin(float64(f32)))
	return ToFloat16(result)
}

// Acos returns the arccosine of f
func Acos(f Float16) Float16 {
	if f.IsNaN() {
		return f
	}
	
	// Check domain: [-1, 1]
	if f.Abs().ToFloat32() > 1.0 {
		return QuietNaN
	}
	
	f32 := f.ToFloat32()
	result := float32(math.Acos(float64(f32)))
	return ToFloat16(result)
}

// Atan returns the arctangent of f
func Atan(f Float16) Float16 {
	if f.IsZero() {
		return f
	}
	if f.IsNaN() {
		return f
	}
	if f.IsInf(1) {
		return Div(Pi, FromInt(2))
	}
	if f.IsInf(-1) {
		return Div(Pi, FromInt(2)).Neg()
	}
	
	f32 := f.ToFloat32()
	result := float32(math.Atan(float64(f32)))
	return ToFloat16(result)
}

// Atan2 returns the arctangent of y/x
func Atan2(y, x Float16) Float16 {
	if y.IsNaN() || x.IsNaN() {
		return QuietNaN
	}
	
	y32 := y.ToFloat32()
	x32 := x.ToFloat32()
	result := float32(math.Atan2(float64(y32), float64(x32)))
	return ToFloat16(result)
}

// Hyperbolic functions

// Sinh returns the hyperbolic sine of f
func Sinh(f Float16) Float16 {
	if f.IsZero() {
		return f
	}
	if f.IsNaN() {
		return f
	}
	if f.IsInf(0) {
		return f
	}
	
	f32 := f.ToFloat32()
	result := float32(math.Sinh(float64(f32)))
	return ToFloat16(result)
}

// Cosh returns the hyperbolic cosine of f
func Cosh(f Float16) Float16 {
	if f.IsZero() {
		return FromInt(1)
	}
	if f.IsNaN() {
		return f
	}
	if f.IsInf(0) {
		return PositiveInfinity
	}
	
	f32 := f.ToFloat32()
	result := float32(math.Cosh(float64(f32)))
	return ToFloat16(result)
}

// Tanh returns the hyperbolic tangent of f
func Tanh(f Float16) Float16 {
	if f.IsZero() {
		return f
	}
	if f.IsNaN() {
		return f
	}
	if f.IsInf(1) {
		return FromInt(1)
	}
	if f.IsInf(-1) {
		return FromInt(-1)
	}
	
	f32 := f.ToFloat32()
	result := float32(math.Tanh(float64(f32)))
	return ToFloat16(result)
}

// Rounding and truncation functions

// Floor returns the largest integer value less than or equal to f
func Floor(f Float16) Float16 {
	if f.IsZero() || f.IsNaN() || f.IsInf(0) {
		return f
	}
	
	f32 := f.ToFloat32()
	result := float32(math.Floor(float64(f32)))
	return ToFloat16(result)
}

// Ceil returns the smallest integer value greater than or equal to f
func Ceil(f Float16) Float16 {
	if f.IsZero() || f.IsNaN() || f.IsInf(0) {
		return f
	}
	
	f32 := f.ToFloat32()
	result := float32(math.Ceil(float64(f32)))
	return ToFloat16(result)
}

// Round returns the nearest integer value to f
func Round(f Float16) Float16 {
	if f.IsZero() || f.IsNaN() || f.IsInf(0) {
		return f
	}
	
	f32 := f.ToFloat32()
	result := float32(math.Round(float64(f32)))
	return ToFloat16(result)
}

// RoundToEven returns the nearest integer value to f, rounding ties to even
func RoundToEven(f Float16) Float16 {
	if f.IsZero() || f.IsNaN() || f.IsInf(0) {
		return f
	}
	
	f32 := f.ToFloat32()
	result := float32(math.RoundToEven(float64(f32)))
	return ToFloat16(result)
}

// Trunc returns the integer part of f (truncated towards zero)
func Trunc(f Float16) Float16 {
	if f.IsZero() || f.IsNaN() || f.IsInf(0) {
		return f
	}
	
	f32 := f.ToFloat32()
	result := float32(math.Trunc(float64(f32)))
	return ToFloat16(result)
}

// Mod returns the floating-point remainder of f/divisor
func Mod(f, divisor Float16) Float16 {
	if divisor.IsZero() {
		return QuietNaN
	}
	if f.IsZero() {
		return f
	}
	if f.IsNaN() || divisor.IsNaN() {
		return QuietNaN
	}
	if f.IsInf(0) || divisor.IsInf(0) {
		return QuietNaN
	}
	
	f32 := f.ToFloat32()
	div32 := divisor.ToFloat32()
	result := float32(math.Mod(float64(f32), float64(div32)))
	return ToFloat16(result)
}

// Remainder returns the IEEE 754 floating-point remainder of f/divisor
func Remainder(f, divisor Float16) Float16 {
	if divisor.IsZero() {
		return QuietNaN
	}
	if f.IsZero() {
		return f
	}
	if f.IsNaN() || divisor.IsNaN() {
		return QuietNaN
	}
	if f.IsInf(0) {
		return QuietNaN
	}
	if divisor.IsInf(0) {
		return f
	}
	
	f32 := f.ToFloat32()
	div32 := divisor.ToFloat32()
	result := float32(math.Remainder(float64(f32), float64(div32)))
	return ToFloat16(result)
}

// Mathematical constants as Float16 values
var (
	E          = ToFloat16(float32(math.E))          // Euler's number
	Pi         = ToFloat16(float32(math.Pi))         // Pi
	Phi        = ToFloat16(float32(math.Phi))        // Golden ratio
	Sqrt2      = ToFloat16(float32(math.Sqrt2))      // Square root of 2
	SqrtE      = ToFloat16(float32(math.SqrtE))      // Square root of E
	SqrtPi     = ToFloat16(float32(math.SqrtPi))     // Square root of Pi
	SqrtPhi    = ToFloat16(float32(math.SqrtPhi))    // Square root of Phi
	Ln2        = ToFloat16(float32(math.Ln2))        // Natural logarithm of 2
	Log2E      = ToFloat16(float32(math.Log2E))      // Base-2 logarithm of E
	Ln10       = ToFloat16(float32(math.Ln10))       // Natural logarithm of 10
	Log10E     = ToFloat16(float32(math.Log10E))     // Base-10 logarithm of E
)

// Utility functions

// Abs returns the absolute value of f
func Abs(f Float16) Float16 {
	return f.Abs()
}

// Clamp restricts f to the range [min, max]
func Clamp(f, min, max Float16) Float16 {
	if f.IsNaN() {
		return f
	}
	if Less(f, min) {
		return min
	}
	if Greater(f, max) {
		return max
	}
	return f
}

// Lerp performs linear interpolation between a and b by factor t
func Lerp(a, b, t Float16) Float16 {
	// lerp(a, b, t) = a + t * (b - a) = a * (1 - t) + b * t
	if t.IsZero() {
		return a
	}
	if Equal(t, FromInt(1)) {
		return b
	}
	
	diff := Sub(b, a)
	scaled := Mul(t, diff)
	return Add(a, scaled)
}

// Sign returns -1, 0, or 1 depending on the sign of f
func Sign(f Float16) Float16 {
	if f.IsNaN() {
		return f
	}
	if f.IsZero() {
		return PositiveZero
	}
	if f.Signbit() {
		return FromInt(-1)
	}
	return FromInt(1)
}

// CopySign returns a Float16 with the magnitude of f and the sign of sign
func CopySign(f, sign Float16) Float16 {
	return f.CopySign(sign)
}

// Dim returns the positive difference between f and g: max(f-g, 0)
func Dim(f, g Float16) Float16 {
	diff := Sub(f, g)
	if Less(diff, PositiveZero) {
		return PositiveZero
	}
	return diff
}

// Hypot returns sqrt(f*f + g*g), taking care to avoid overflow and underflow
func Hypot(f, g Float16) Float16 {
	if f.IsInf(0) || g.IsInf(0) {
		return PositiveInfinity
	}
	if f.IsNaN() || g.IsNaN() {
		return QuietNaN
	}
	
	f32 := f.ToFloat32()
	g32 := g.ToFloat32()
	result := float32(math.Hypot(float64(f32), float64(g32)))
	return ToFloat16(result)
}

// Gamma returns the Gamma function of f
func Gamma(f Float16) Float16 {
	if f.IsNaN() {
		return f
	}
	if f.IsInf(-1) {
		return QuietNaN
	}
	if f.IsInf(1) {
		return PositiveInfinity
	}
	
	f32 := f.ToFloat32()
	result := float32(math.Gamma(float64(f32)))
	return ToFloat16(result)
}

// Lgamma returns the natural logarithm and sign of Gamma(f)
func Lgamma(f Float16) (Float16, int) {
	if f.IsNaN() {
		return f, 1
	}
	
	f32 := f.ToFloat32()
	lgamma, sign := math.Lgamma(float64(f32))
	return ToFloat16(float32(lgamma)), sign
}

// J0 returns the order-zero Bessel function of the first kind
func J0(f Float16) Float16 {
	if f.IsNaN() {
		return f
	}
	if f.IsInf(0) {
		return PositiveZero
	}
	
	f32 := f.ToFloat32()
	result := float32(math.J0(float64(f32)))
	return ToFloat16(result)
}

// J1 returns the order-one Bessel function of the first kind
func J1(f Float16) Float16 {
	if f.IsNaN() {
		return f
	}
	if f.IsInf(0) {
		return PositiveZero
	}
	
	f32 := f.ToFloat32()
	result := float32(math.J1(float64(f32)))
	return ToFloat16(result)
}

// Y0 returns the order-zero Bessel function of the second kind
func Y0(f Float16) Float16 {
	if f.IsNaN() || f.Signbit() {
		return QuietNaN
	}
	if f.IsZero() {
		return NegativeInfinity
	}
	if f.IsInf(1) {
		return PositiveZero
	}
	
	f32 := f.ToFloat32()
	result := float32(math.Y0(float64(f32)))
	return ToFloat16(result)
}

// Y1 returns the order-one Bessel function of the second kind
func Y1(f Float16) Float16 {
	if f.IsNaN() || f.Signbit() {
		return QuietNaN
	}
	if f.IsZero() {
		return NegativeInfinity
	}
	if f.IsInf(1) {
		return PositiveZero
	}
	
	f32 := f.ToFloat32()
	result := float32(math.Y1(float64(f32)))
	return ToFloat16(result)
}

// Erf returns the error function of f
func Erf(f Float16) Float16 {
	if f.IsZero() {
		return f
	}
	if f.IsNaN() {
		return f
	}
	if f.IsInf(1) {
		return FromInt(1)
	}
	if f.IsInf(-1) {
		return FromInt(-1)
	}
	
	f32 := f.ToFloat32()
	result := float32(math.Erf(float64(f32)))
	return ToFloat16(result)
}

// Erfc returns the complementary error function of f
func Erfc(f Float16) Float16 {
	if f.IsNaN() {
		return f
	}
	if f.IsInf(1) {
		return PositiveZero
	}
	if f.IsInf(-1) {
		return FromInt(2)
	}
	
	f32 := f.ToFloat32()
	result := float32(math.Erfc(float64(f32)))
	return ToFloat16(result)
}
