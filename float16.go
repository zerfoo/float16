// Package float16 implements the 16-bit floating point data type (IEEE 754-2008).
//
// This implementation provides conversion between float16 and other floating-point types
// (float32 and float64) with support for various rounding modes and error handling.
//
// # Special Values
//
// The float16 type supports all IEEE 754-2008 special values:
//   - Positive and negative zero
//   - Positive and negative infinity
//   - Not-a-Number (NaN) values with payload
//   - Normalized numbers
//   - Subnormal (denormal) numbers
//
// # Subnormal Numbers
//
// When converting to higher-precision types (float32/float64), subnormal float16 values
// are preserved. However, when converting back from higher-precision types to float16,
// subnormal values may be rounded to the nearest representable normal float16 value.
// This behavior is consistent with many hardware implementations that handle subnormals
// in a similar way for performance reasons.
//
// # Rounding Modes
//
// The following rounding modes are supported for conversions:
//   - RoundNearestEven: Round to nearest, ties to even (default)
//   - RoundTowardZero: Round toward zero (truncate)
//   - RoundTowardPositive: Round toward positive infinity
//   - RoundTowardNegative: Round toward negative infinity
//   - RoundNearestAway: Round to nearest, ties away from zero
//
// # Error Handling
//
// Conversion functions with a ConversionMode parameter can return errors for:
//   - Overflow: When a value is too large to be represented
//   - Underflow: When a value is too small to be represented (in strict mode)
//   - Inexact: When rounding occurs (in strict mode)
//
// See: http://en.wikipedia.org/wiki/Half-precision_floating-point_format
package float16

import (
	"math"
	"sync"
)

// Package version information
const (
	Version      = "1.0.0"
	VersionMajor = 1
	VersionMinor = 0
	VersionPatch = 0
)

// Package configuration
type Config struct {
	DefaultConversionMode ConversionMode
	DefaultRoundingMode   RoundingMode
	DefaultArithmeticMode ArithmeticMode
	EnableFastMath        bool // Package float16 implements the 16-bit floating point data type (IEEE 754-2008).
	// This implementation provides conversion between float16 and other floating-point types
	// (float32 and float64) with support for various rounding modes and error handling.
	//
	// # Special Values
	//
	// The float16 type supports all IEEE 754-2008 special values:
	//   - Positive and negative zero
	//   - Positive and negative infinity
	//   - Not-a-Number (NaN) values with payload
	//   - Normalized numbers
	//   - Subnormal (denormal) numbers
	//
	// # Subnormal Numbers
	//
	// When converting to higher-precision types (float32/float64), subnormal float16 values
	// are preserved. However, when converting back from higher-precision types to float16,
	// subnormal values may be rounded to the nearest representable normal float16 value.
	// This behavior is consistent with many hardware implementations that handle subnormals
	// in a similar way for performance reasons.
	//
	// # Rounding Modes
	//
	// The following rounding modes are supported for conversions:
	//   - RoundNearestEven: Round to nearest, ties to even (default)
	//   - RoundTowardZero: Round toward zero (truncate)
	//   - RoundTowardPositive: Round toward positive infinity
	//   - RoundTowardNegative: Round toward negative infinity
	//   - RoundNearestAway: Round to nearest, ties away from zero
	//
	// # Error Handling
	//
	// Conversion functions with a ConversionMode parameter can return errors for:
	//   - Overflow: When a value is too large to be represented
	//   - Underflow: When a value is too small to be represented (in strict mode)
	//   - Inexact: When rounding occurs (in strict mode)
	//
	// See: http://en.wikipedia.org/wiki/Half-precision_floating-point_format
}

// DefaultConfig returns the default package configuration
func DefaultConfig() *Config {
	return &Config{
		DefaultConversionMode: DefaultConversionMode,
		DefaultRoundingMode:   DefaultRoundingMode,
		DefaultArithmeticMode: ModeIEEEArithmetic,
		EnableFastMath:        false,
	}
}

var (
	configMutex sync.RWMutex
	config      = DefaultConfig()
)

// Configure applies the given configuration to the package
func Configure(cfg *Config) {
	configMutex.Lock()
	defer configMutex.Unlock()

	config = cfg
	DefaultConversionMode = cfg.DefaultConversionMode
	DefaultRoundingMode = cfg.DefaultRoundingMode
	DefaultArithmeticMode = cfg.DefaultArithmeticMode
}

// GetConfig returns the current package configuration
func GetConfig() *Config {
	configMutex.RLock()
	defer configMutex.RUnlock()

	// Return a copy to prevent external modification
	return &Config{
		DefaultConversionMode: config.DefaultConversionMode,
		DefaultRoundingMode:   config.DefaultRoundingMode,
		DefaultArithmeticMode: config.DefaultArithmeticMode,
		EnableFastMath:        config.EnableFastMath,
	}
}

// GetVersion returns the package version string
func GetVersion() string {
	return Version
}

// Convenience functions for common operations

// Zero returns a Float16 zero value
func Zero() Float16 {
	return PositiveZero
}

// One returns a Float16 value representing 1.0
func One() Float16 {
	converter := NewConverter(DefaultConversionMode, DefaultRoundingMode)
	return converter.ToFloat16(1.0)
}

// NaN returns a Float16 quiet NaN value
func NaN() Float16 {
	return QuietNaN
}

// Inf returns a Float16 infinity value
// If sign >= 0, returns positive infinity
// If sign < 0, returns negative infinity
func Inf(sign int) Float16 {
	if sign >= 0 {
		return PositiveInfinity
	}
	return NegativeInfinity
}

// IsInf reports whether f is an infinity, according to sign
// If sign > 0, IsInf reports whether f is positive infinity
// If sign < 0, IsInf reports whether f is negative infinity
// If sign == 0, IsInf reports whether f is either infinity
func IsInf(f Float16, sign int) bool {
	return f.IsInf(sign)
}

// IsNaN reports whether f is an IEEE 754 "not-a-number" value
func IsNaN(f Float16) bool {
	return f.IsNaN()
}

// Signbit reports whether f is negative or negative zero
func Signbit(f Float16) bool {
	return f.Signbit()
}

// Utility functions for working with Float16 values

// NextAfter returns the next representable Float16 value after f in the direction of g
func NextAfter(f, g Float16) Float16 {
	if f.IsNaN() || g.IsNaN() {
		return QuietNaN
	}

	if Equal(f, g) {
		return g
	}

	if f.IsZero() {
		if g.Signbit() {
			return FromBits(0x8001) // Smallest negative subnormal
		}
		return FromBits(0x0001) // Smallest positive subnormal
	}

	bits := f.Bits()
	if (f.ToFloat32() < g.ToFloat32()) == !f.Signbit() {
		bits++
	} else {
		bits--
	}

	return FromBits(bits)
}

// Frexp breaks f into a normalized fraction and an integral power of two
// It returns frac and exp satisfying f == frac × 2^exp, with the absolute
// value of frac in the interval [0.5, 1) or zero
func Frexp(f Float16) (frac Float16, exp int) {
	if f.IsZero() || f.IsNaN() || f.IsInf(0) {
		return f, 0
	}

	f32 := f.ToFloat32()
	frac32, exp := math.Frexp(float64(f32))
	converter := NewConverter(DefaultConversionMode, DefaultRoundingMode)
	return converter.ToFloat16(float32(frac32)), exp
}

// Ldexp returns frac × 2^exp
func Ldexp(frac Float16, exp int) Float16 {
	if frac.IsZero() || frac.IsNaN() || frac.IsInf(0) {
		return frac
	}

	frac32 := frac.ToFloat32()
	result := math.Ldexp(float64(frac32), exp)
	converter := NewConverter(DefaultConversionMode, DefaultRoundingMode)
	return converter.ToFloat16(float32(result))
}

// Modf returns integer and fractional floating-point numbers that sum to f
// Both values have the same sign as f
func Modf(f Float16) (integer, frac Float16) {
	if f.IsNaN() || f.IsInf(0) {
		return f, f
	}

	f32 := f.ToFloat32()
	int32, frac32 := math.Modf(float64(f32))
	converter := NewConverter(DefaultConversionMode, DefaultRoundingMode)
	return converter.ToFloat16(float32(int32)), converter.ToFloat16(float32(frac32))
}

// Validation and classification functions

// IsFinite reports whether f is neither infinite nor NaN
func IsFinite(f Float16) bool {
	return f.IsFinite()
}

// IsNormal reports whether f is a normal number (not zero, subnormal, infinite, or NaN)
func IsNormal(f Float16) bool {
	return f.IsNormal()
}

// IsSubnormal reports whether f is a subnormal number
func IsSubnormal(f Float16) bool {
	return f.IsSubnormal()
}

// FpClassify returns the IEEE 754 class of f
func FpClassify(f Float16) FloatClass {
	return f.Class()
}

// Performance monitoring and debugging

// GetMemoryUsage returns the current memory usage of the package in bytes
func GetMemoryUsage() int {
	// Float16 package uses minimal memory (no lookup tables)
	// Only constants and code, estimated at ~8KB
	return 8192
}

// DebugInfo returns debugging information about the package state
func DebugInfo() map[string]interface{} {
	cfg := GetConfig()
	return map[string]interface{}{
		"version":                 Version,
		"memory_usage_bytes":      GetMemoryUsage(),
		"default_conversion_mode": cfg.DefaultConversionMode,
		"default_rounding_mode":   cfg.DefaultRoundingMode,
		"default_arithmetic_mode": cfg.DefaultArithmeticMode,
		"fast_math_enabled":       cfg.EnableFastMath,
		"ieee754_compliant":       true,
		"supports_subnormals":     true,
		"lookup_tables":           false,
	}
}

// Benchmark helpers for performance testing

// BenchmarkOperation represents a benchmarkable operation
type BenchmarkOperation func(Float16, Float16) Float16

// GetBenchmarkOperations returns a map of operations suitable for benchmarking
func GetBenchmarkOperations() map[string]BenchmarkOperation {
	return map[string]BenchmarkOperation{
		"Add": Add,
		"Sub": Sub,
		"Mul": Mul,
		"Div": Div,
	}
}

// Constants for common values
var (
	// Common integer values
	Zero16  = PositiveZero
	One16   = NewConverter(DefaultConversionMode, DefaultRoundingMode).ToFloat16(1.0)
	Two16   = NewConverter(DefaultConversionMode, DefaultRoundingMode).ToFloat16(2.0)
	Three16 = NewConverter(DefaultConversionMode, DefaultRoundingMode).ToFloat16(3.0)
	Four16  = NewConverter(DefaultConversionMode, DefaultRoundingMode).ToFloat16(4.0)
	Five16  = NewConverter(DefaultConversionMode, DefaultRoundingMode).ToFloat16(5.0)
	Ten16   = NewConverter(DefaultConversionMode, DefaultRoundingMode).ToFloat16(10.0)

	// Common fractional values
	Half16    = NewConverter(DefaultConversionMode, DefaultRoundingMode).ToFloat16(0.5)
	Quarter16 = NewConverter(DefaultConversionMode, DefaultRoundingMode).ToFloat16(0.25)
	Third16   = NewConverter(DefaultConversionMode, DefaultRoundingMode).ToFloat16(1.0 / 3.0)

	// Special mathematical values
	NaN16  = QuietNaN
	PosInf = PositiveInfinity
	NegInf = NegativeInfinity

	// Commonly used constants
	Deg2Rad = NewConverter(DefaultConversionMode, DefaultRoundingMode).ToFloat16(float32(math.Pi / 180.0)) // Degrees to radians
	Rad2Deg = NewConverter(DefaultConversionMode, DefaultRoundingMode).ToFloat16(float32(180.0 / math.Pi)) // Radians to degrees
)

// Helper functions for slice operations with error handling

// ValidateSliceLength checks if two slices have the same length
func ValidateSliceLength(a, b []Float16) error {
	if len(a) != len(b) {
		return &Float16Error{
			Op:   "slice_operation",
			Msg:  "slice length mismatch",
			Code: ErrInvalidOperation,
		}
	}
	return nil
}

// SliceStats computes basic statistics for a Float16 slice
type SliceStats struct {
	Min    Float16
	Max    Float16
	Sum    Float16
	Mean   Float16
	Length int
}

// ComputeSliceStats calculates statistics for a Float16 slice
func ComputeSliceStats(s []Float16) SliceStats {
	if len(s) == 0 {
		return SliceStats{}
	}

	stats := SliceStats{
		Min:    s[0],
		Max:    s[0],
		Sum:    PositiveZero,
		Length: len(s),
	}

	for _, v := range s {
		if !v.IsNaN() {
			if Less(v, stats.Min) {
				stats.Min = v
			}
			if Greater(v, stats.Max) {
				stats.Max = v
			}
		}
		stats.Sum = Add(stats.Sum, v)
	}

	if stats.Length > 0 {
		stats.Mean = Div(stats.Sum, NewConverter(DefaultConversionMode, DefaultRoundingMode).FromInt(stats.Length))
	}

	return stats
}

// Experimental features (may change in future versions)

// FastAdd performs addition optimized for speed (may sacrifice precision)
func FastAdd(a, b Float16) Float16 {
	converter := NewConverter(DefaultConversionMode, DefaultRoundingMode)
	return converter.ToFloat16(a.ToFloat32() + b.ToFloat32())
}

// FastMul performs multiplication optimized for speed (may sacrifice precision)
func FastMul(a, b Float16) Float16 {
	converter := NewConverter(DefaultConversionMode, DefaultRoundingMode)
	return converter.ToFloat16(a.ToFloat32() * b.ToFloat32())
}

// VectorAdd performs vectorized addition (placeholder for future SIMD implementation)
func VectorAdd(a, b []Float16) []Float16 {
	// Currently just calls the regular slice operation
	// Future versions may implement SIMD optimizations
	return AddSlice(a, b)
}

// VectorMul performs vectorized multiplication (placeholder for future SIMD implementation)
func VectorMul(a, b []Float16) []Float16 {
	// Currently just calls the regular slice operation
	// Future versions may implement SIMD optimizations
	return MulSlice(a, b)
}
