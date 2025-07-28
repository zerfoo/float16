# float16

[![Go Reference](https://pkg.go.dev/badge/github.com/zerfoo/float16.svg)](https://pkg.go.dev/github.com/zerfoo/float16)
[![Go Report Card](https://goreportcard.com/badge/github.com/zerfoo/float16)](https://goreportcard.com/report/github.com/zerfoo/float16)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

A comprehensive Go implementation of IEEE 754-2008 16-bit floating-point (half-precision) arithmetic with full support for special values, multiple rounding modes, and high-performance operations.

## Features

- **Full IEEE 754-2008 compliance** for 16-bit floating-point arithmetic
- **Complete special value support**: ±0, ±∞, NaN (with payload), normalized and subnormal numbers
- **Multiple rounding modes**: nearest-even, toward zero, toward ±∞, nearest-away
- **Flexible conversion modes**: IEEE standard, strict error handling, fast approximations
- **High-performance operations** with optional fast math optimizations
- **Comprehensive test suite** with extensive edge case coverage
- **Zero dependencies** - pure Go implementation

## Installation

```bash
go get github.com/zerfoo/float16
```

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/zerfoo/float16"
)

func main() {
    // Create float16 values
    a := float16.FromFloat32(3.14159)
    b := float16.FromFloat64(2.71828)
    
    // Basic arithmetic
    sum := a.Add(b)
    product := a.Mul(b)
    
    // Convert back to other types
    fmt.Printf("Sum: %v (float32: %f)\n", sum, sum.ToFloat32())
    fmt.Printf("Product: %v (float64: %f)\n", product, product.ToFloat64())
    
    // Work with special values
    inf := float16.Inf(1)  // positive infinity
    nan := float16.NaN()   // quiet NaN
    zero := float16.Zero() // positive zero
    
    fmt.Printf("Infinity: %v\n", inf)
    fmt.Printf("NaN: %v\n", nan)
    fmt.Printf("Zero: %v\n", zero)
}
```

## Core Types and Constants

### Float16 Type

The `Float16` type represents a 16-bit IEEE 754 half-precision floating-point value:

```go
type Float16 uint16
```

### Special Values

```go
const (
    PositiveZero     Float16 = 0x0000 // +0.0
    NegativeZero     Float16 = 0x8000 // -0.0
    PositiveInfinity Float16 = 0x7C00 // +∞
    NegativeInfinity Float16 = 0xFC00 // -∞
    MaxValue         Float16 = 0x7BFF // ~65504
    MinValue         Float16 = 0xFBFF // ~-65504
)
```

## Conversion Functions

### From Other Types

```go
// From float32/float64
f16 := float16.FromFloat32(3.14159)
f16 := float16.FromFloat64(2.71828)

// From bit representation
f16 := float16.FromBits(0x4200) // 3.0

// From string
f16, err := float16.ParseFloat("3.14159", 32)
```

### To Other Types

```go
f32 := f16.ToFloat32()
f64 := f16.ToFloat64()
bits := f16.Bits()
str := f16.String()
```

## Arithmetic Operations

```go
a := float16.FromFloat32(5.0)
b := float16.FromFloat32(3.0)

// Basic arithmetic
sum := a.Add(b)        // 8.0
diff := a.Sub(b)       // 2.0
product := a.Mul(b)    // 15.0
quotient := a.Div(b)   // 1.666...

// Mathematical functions
sqrt := a.Sqrt()       // √5
abs := a.Abs()         // |a|
neg := a.Neg()         // -a
```

## Rounding Modes

Configure rounding behavior for conversions:

```go
import "github.com/zerfoo/float16"

// Set global rounding mode
config := float16.GetConfig()
config.DefaultRoundingMode = float16.RoundTowardZero
float16.Configure(config)

// Available rounding modes:
// - RoundNearestEven (default)
// - RoundTowardZero
// - RoundTowardPositive  
// - RoundTowardNegative
// - RoundNearestAway
```

## Conversion Modes

Control conversion behavior and error handling:

```go
config := float16.GetConfig()
config.DefaultConversionMode = float16.ModeStrict
float16.Configure(config)

// Available modes:
// - ModeIEEE: Standard IEEE 754 behavior
// - ModeStrict: Returns errors for overflow/underflow
// - ModeFast: Optimized for performance
```

## Special Value Handling

```go
f := float16.FromFloat32(math.Inf(1))

// Check value types
if f.IsInf(0) {
    fmt.Println("Value is infinity")
}
if f.IsNaN() {
    fmt.Println("Value is NaN")
}
if f.IsFinite() {
    fmt.Println("Value is finite")
}
if f.IsNormal() {
    fmt.Println("Value is normalized")
}
if f.IsSubnormal() {
    fmt.Println("Value is subnormal")
}

// IEEE 754 classification
class := f.Class()
switch class {
case float16.ClassPositiveInfinity:
    fmt.Println("Positive infinity")
case float16.ClassQuietNaN:
    fmt.Println("Quiet NaN")
// ... other classes
}
```

## Performance Features

### Fast Math Operations

```go
// Enable fast math for better performance (may sacrifice precision)
config := float16.GetConfig()
config.EnableFastMath = true
float16.Configure(config)

// Use fast operations
result := float16.FastAdd(a, b)
result := float16.FastMul(a, b)
```

### Vectorized Operations

```go
// Vectorized operations (optimized for SIMD when available)
a := []float16.Float16{...}
b := []float16.Float16{...}

sum := float16.VectorAdd(a, b)
product := float16.VectorMul(a, b)
```

## Error Handling

```go
// Strict mode returns errors for exceptional conditions
config := float16.GetConfig()
config.DefaultConversionMode = float16.ModeStrict
float16.Configure(config)

f16, err := float16.FromFloat32WithMode(1e10, float16.ModeStrict)
if err != nil {
    if float16Err, ok := err.(*float16.Float16Error); ok {
        switch float16Err.Code {
        case float16.ErrOverflow:
            fmt.Println("Value too large for float16")
        case float16.ErrUnderflow:
            fmt.Println("Value too small for float16")
        }
    }
}
```

## Utilities

### Statistics for Slices

```go
values := []float16.Float16{
    float16.FromFloat32(1.0),
    float16.FromFloat32(2.0),
    float16.FromFloat32(3.0),
}

stats := float16.ComputeSliceStats(values)
fmt.Printf("Min: %v, Max: %v, Mean: %v\n", stats.Min, stats.Max, stats.Mean)
```

### Debugging and Monitoring

```go
// Get memory usage
usage := float16.GetMemoryUsage()
fmt.Printf("Memory usage: %d bytes\n", usage)

// Get debug information
debug := float16.DebugInfo()
fmt.Printf("Debug info: %+v\n", debug)
```

## Benchmarking

The package includes built-in benchmarking utilities:

```go
ops := float16.GetBenchmarkOperations()
for name, op := range ops {
    // Benchmark operation
    fmt.Printf("Benchmarking %s\n", name)
}
```

## Range and Precision

Float16 has the following characteristics:

- **Range**: ±6.55×10⁴ (approximately ±65,504)
- **Precision**: ~3-4 decimal digits
- **Smallest positive normal**: ~6.10×10⁻⁵
- **Smallest positive subnormal**: ~5.96×10⁻⁸
- **Machine epsilon**: ~9.77×10⁻⁴

## Use Cases

Float16 is ideal for:

- **Machine Learning**: Reduced memory usage and faster training
- **Graphics Programming**: Color values, texture coordinates
- **Scientific Computing**: Large datasets where precision can be traded for memory
- **Embedded Systems**: Memory-constrained environments
- **Data Compression**: Storing floating-point data more efficiently

## Performance Considerations

- Conversions between float16 and float32/float64 have computational overhead
- Native float16 arithmetic is generally faster than conversion-based approaches
- Enable fast math mode for performance-critical applications where precision can be sacrificed
- Use vectorized operations for bulk processing

## Contributing

We welcome contributions! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.

## References

- [IEEE 754-2008 Standard](https://ieeexplore.ieee.org/document/4610935)
- [Half-precision floating-point format](https://en.wikipedia.org/wiki/Half-precision_floating-point_format)
