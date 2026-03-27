# float16

[![Go Reference](https://pkg.go.dev/badge/github.com/zerfoo/float16.svg)](https://pkg.go.dev/github.com/zerfoo/float16)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

IEEE 754-2008 half-precision (Float16) and BFloat16 arithmetic library for Go.

Part of the [Zerfoo](https://github.com/zerfoo) ML ecosystem.

## Features

- **Full IEEE 754-2008 compliance** for 16-bit floating-point arithmetic
- **BFloat16 support** — Google Brain format for ML training and inference
- **Special value handling** — ±0, ±Inf, NaN (with payload), normalized and subnormal numbers
- **Multiple rounding modes** — nearest-even, toward zero, toward ±Inf, nearest-away
- **Vectorized operations** — batch add, multiply, and dot product
- **Fast math mode** — optional lookup-table acceleration for performance-critical paths
- **Zero dependencies** — pure Go, no CGo

## Installation

```bash
go get github.com/zerfoo/float16
```

Requires Go 1.26+.

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/zerfoo/float16"
)

func main() {
    a := float16.FromFloat32(3.14159)
    b := float16.FromFloat32(2.71828)

    sum := a.Add(b)
    product := a.Mul(b)

    fmt.Printf("Sum: %f\n", sum.ToFloat32())
    fmt.Printf("Product: %f\n", product.ToFloat32())

    // Special values
    inf := float16.Inf(1)
    fmt.Printf("Inf: %v, IsInf: %v\n", inf, inf.IsInf(0))
}
```

## Conversion

```go
// From float32/float64
f16 := float16.FromFloat32(3.14)
f16 := float16.FromFloat64(2.718)

// From bit representation
f16 := float16.FromBits(0x4200) // 3.0

// Back to native types
f32 := f16.ToFloat32()
f64 := f16.ToFloat64()
```

## Rounding Modes

```go
config := float16.GetConfig()
config.DefaultRoundingMode = float16.RoundTowardZero
float16.Configure(config)

// RoundNearestEven (default), RoundTowardZero, RoundTowardPositive,
// RoundTowardNegative, RoundNearestAway
```

## Range and Precision

| Property | Value |
|----------|-------|
| Range | ±65,504 |
| Precision | ~3-4 decimal digits |
| Smallest normal | ~6.10 × 10⁻⁵ |
| Smallest subnormal | ~5.96 × 10⁻⁸ |
| Machine epsilon | ~9.77 × 10⁻⁴ |

## Used By

- [ztensor](https://github.com/zerfoo/ztensor) — GPU-accelerated tensor library

## License

Apache 2.0
