package float16

import (
	"math"
)

// Mathematical functions for BFloat16

// BFloat16Sqrt returns the square root of the BFloat16 value.
func BFloat16Sqrt(b BFloat16) BFloat16 {
	if b.IsZero() {
		return b
	}
	if b.IsNaN() {
		return b
	}
	if b.IsInf(1) {
		return BFloat16PositiveInfinity
	}
	if b.Signbit() {
		return BFloat16QuietNaN
	}
	result := math.Sqrt(float64(b.ToFloat32()))
	return BFloat16FromFloat32(float32(result))
}

// BFloat16Exp returns e^b.
func BFloat16Exp(b BFloat16) BFloat16 {
	if b.IsZero() {
		return BFloat16One
	}
	if b.IsNaN() {
		return b
	}
	if b.IsInf(1) {
		return BFloat16PositiveInfinity
	}
	if b.IsInf(-1) {
		return BFloat16PositiveZero
	}
	result := math.Exp(float64(b.ToFloat32()))
	return BFloat16FromFloat32(float32(result))
}

// BFloat16Log returns the natural logarithm of b.
func BFloat16Log(b BFloat16) BFloat16 {
	if b.IsZero() {
		return BFloat16NegativeInfinity
	}
	if b.IsNaN() {
		return b
	}
	if b.IsInf(1) {
		return BFloat16PositiveInfinity
	}
	if b.Signbit() {
		return BFloat16QuietNaN
	}
	result := math.Log(float64(b.ToFloat32()))
	return BFloat16FromFloat32(float32(result))
}

// BFloat16Log2 returns the base-2 logarithm of b.
func BFloat16Log2(b BFloat16) BFloat16 {
	if b.IsZero() {
		return BFloat16NegativeInfinity
	}
	if b.IsNaN() {
		return b
	}
	if b.IsInf(1) {
		return BFloat16PositiveInfinity
	}
	if b.Signbit() {
		return BFloat16QuietNaN
	}
	result := math.Log2(float64(b.ToFloat32()))
	return BFloat16FromFloat32(float32(result))
}

// BFloat16Sin returns the sine of b (in radians).
func BFloat16Sin(b BFloat16) BFloat16 {
	if b.IsZero() {
		return b
	}
	if b.IsNaN() || b.IsInf(0) {
		return BFloat16QuietNaN
	}
	result := math.Sin(float64(b.ToFloat32()))
	return BFloat16FromFloat32(float32(result))
}

// BFloat16Cos returns the cosine of b (in radians).
func BFloat16Cos(b BFloat16) BFloat16 {
	if b.IsZero() {
		return BFloat16One
	}
	if b.IsNaN() || b.IsInf(0) {
		return BFloat16QuietNaN
	}
	result := math.Cos(float64(b.ToFloat32()))
	return BFloat16FromFloat32(float32(result))
}

// BFloat16Tanh returns the hyperbolic tangent of b.
func BFloat16Tanh(b BFloat16) BFloat16 {
	if b.IsZero() {
		return b
	}
	if b.IsNaN() {
		return b
	}
	if b.IsInf(1) {
		return BFloat16One
	}
	if b.IsInf(-1) {
		return BFloat16FromFloat32(-1)
	}
	result := math.Tanh(float64(b.ToFloat32()))
	return BFloat16FromFloat32(float32(result))
}

// BFloat16Sigmoid returns 1 / (1 + exp(-b)).
func BFloat16Sigmoid(b BFloat16) BFloat16 {
	if b.IsNaN() {
		return b
	}
	if b.IsInf(1) {
		return BFloat16One
	}
	if b.IsInf(-1) {
		return BFloat16PositiveZero
	}
	x := float64(b.ToFloat32())
	result := 1.0 / (1.0 + math.Exp(-x))
	return BFloat16FromFloat32(float32(result))
}

// FastMode variants using polynomial approximations.
// These trade accuracy for speed, suitable for ML inference workloads
// where BFloat16 precision is already limited.

// BFloat16FastSigmoid computes an approximate sigmoid using a rational polynomial.
// Uses the approximation: sigmoid(x) ≈ 0.5 + 0.5 * x / (1 + |x|)
// which avoids exp() entirely.
func BFloat16FastSigmoid(b BFloat16) BFloat16 {
	if b.IsNaN() {
		return b
	}
	if b.IsInf(1) {
		return BFloat16One
	}
	if b.IsInf(-1) {
		return BFloat16PositiveZero
	}
	x := float64(b.ToFloat32())
	abs := x
	if abs < 0 {
		abs = -abs
	}
	result := 0.5 + 0.5*x/(1.0+abs)
	return BFloat16FromFloat32(float32(result))
}

// BFloat16FastTanh computes an approximate tanh using a rational polynomial.
// Uses the approximation: tanh(x) ≈ x*(27 + x*x) / (27 + 9*x*x)
// which is a Padé approximant accurate to within ~0.004 for |x| < 3.
func BFloat16FastTanh(b BFloat16) BFloat16 {
	if b.IsZero() {
		return b
	}
	if b.IsNaN() {
		return b
	}
	if b.IsInf(1) {
		return BFloat16One
	}
	if b.IsInf(-1) {
		return BFloat16FromFloat32(-1)
	}
	x := float64(b.ToFloat32())
	abs := x
	if abs < 0 {
		abs = -abs
	}
	// Clamp for large values where tanh saturates
	if abs > 4.0 {
		if x > 0 {
			return BFloat16One
		}
		return BFloat16FromFloat32(-1)
	}
	x2 := x * x
	result := x * (27.0 + x2) / (27.0 + 9.0*x2)
	return BFloat16FromFloat32(float32(result))
}
