# Contributing to float16

Thank you for your interest in contributing to the float16 Go package! This document provides guidelines and information for contributors.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [Contributing Guidelines](#contributing-guidelines)
- [Testing](#testing)
- [Code Style](#code-style)
- [Submitting Changes](#submitting-changes)
- [Issue Reporting](#issue-reporting)
- [Performance Considerations](#performance-considerations)
- [Documentation](#documentation)

## Code of Conduct

This project adheres to a code of conduct that promotes a welcoming and inclusive environment. Please be respectful and professional in all interactions.

## Getting Started

### Prerequisites

- Go 1.24 or later
- Git
- Basic understanding of IEEE 754 floating-point arithmetic (helpful but not required)

### Development Setup

1. **Fork and clone the repository:**
   ```bash
   git clone https://github.com/YOUR_USERNAME/float16.git
   cd float16
   ```

2. **Verify your setup:**
   ```bash
   go mod tidy
   go test ./...
   ```

3. **Run benchmarks to establish baseline:**
   ```bash
   go test -bench=. -benchmem
   ```

## Contributing Guidelines

### Types of Contributions

We welcome the following types of contributions:

- **Bug fixes**: Fixes for incorrect behavior or edge cases
- **Performance improvements**: Optimizations that maintain correctness
- **New features**: Additional functionality that aligns with the package goals
- **Documentation**: Improvements to code comments, README, or examples
- **Tests**: Additional test cases, especially for edge cases
- **Benchmarks**: Performance tests for critical operations

### Before You Start

1. **Check existing issues** to see if your contribution is already being worked on
2. **Open an issue** for significant changes to discuss the approach
3. **Start small** - consider beginning with documentation or test improvements
4. **Understand the scope** - this package focuses on IEEE 754-2008 compliance

## Development Setup

### Project Structure

```
float16/
├── README.md              # Package documentation
├── CONTRIBUTING.md        # This file
├── LICENSE               # Apache 2.0 license
├── go.mod               # Go module definition
├── float16.go           # Main package file with utilities
├── types.go             # Core Float16 type and constants
├── convert.go           # Conversion functions
├── arithmetic.go        # Arithmetic operations
├── math.go             # Mathematical functions
├── *_test.go           # Test files
└── coverage.html       # Test coverage report
```

### Key Components

- **`Float16`**: The main 16-bit floating-point type
- **Conversion functions**: Between float16 and other types
- **Arithmetic operations**: Add, subtract, multiply, divide
- **Mathematical functions**: Sqrt, abs, sign operations
- **Special value handling**: NaN, infinity, zero detection
- **Configuration system**: Rounding modes and conversion modes

## Testing

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Generate detailed coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html

# Run specific test patterns
go test -run TestConversion
go test -run TestArithmetic
```

### Test Categories

1. **Unit tests**: Test individual functions and methods
2. **Conversion tests**: Verify accuracy of type conversions
3. **Arithmetic tests**: Test mathematical operations
4. **Edge case tests**: Special values, overflow, underflow
5. **Roundtrip tests**: Ensure conversion consistency
6. **Benchmark tests**: Performance measurements

### Writing Tests

When adding tests, follow these guidelines:

```go
func TestNewFeature(t *testing.T) {
    tests := []struct {
        name     string
        input    float32
        expected Float16
        wantErr  bool
    }{
        {"positive normal", 1.0, 0x3C00, false},
        {"negative normal", -1.0, 0xBC00, false},
        {"zero", 0.0, 0x0000, false},
        {"infinity", math.Inf(1), PositiveInfinity, false},
        // Add edge cases
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result, err := SomeFunction(tt.input)
            
            if tt.wantErr {
                if err == nil {
                    t.Errorf("expected error but got none")
                }
                return
            }
            
            if err != nil {
                t.Errorf("unexpected error: %v", err)
                return
            }
            
            if result != tt.expected {
                t.Errorf("got %v, want %v", result, tt.expected)
            }
        })
    }
}
```

### Test Requirements

- **All new code must have tests** with at least 80% coverage
- **Include edge cases**: NaN, infinity, zero, subnormal values
- **Test error conditions** when applicable
- **Verify IEEE 754 compliance** for critical operations
- **Add benchmarks** for performance-sensitive code

## Code Style

### Go Standards

Follow standard Go conventions:

- Use `gofmt` for formatting
- Follow effective Go guidelines
- Use meaningful variable and function names
- Add comprehensive documentation comments

### Documentation Comments

```go
// FunctionName performs a specific operation on Float16 values.
// It returns the result and an error if the operation cannot be completed.
//
// Special cases are:
//   - FunctionName(NaN) = NaN
//   - FunctionName(±Inf) = ±Inf
//   - FunctionName(±0) = ±0
func FunctionName(f Float16) (Float16, error) {
    // Implementation
}
```

### Naming Conventions

- **Types**: PascalCase (`Float16`, `ConversionMode`)
- **Functions**: PascalCase for exported, camelCase for internal
- **Constants**: PascalCase for exported (`PositiveInfinity`)
- **Variables**: camelCase
- **Test functions**: `TestFunctionName`
- **Benchmark functions**: `BenchmarkFunctionName`

### Error Handling

```go
// Use custom error types for package-specific errors
func SomeOperation(f Float16) (Float16, error) {
    if f.IsNaN() {
        return NaN(), &Float16Error{
            Op:    "SomeOperation",
            Value: f,
            Msg:   "operation not defined for NaN",
            Code:  ErrInvalidOperation,
        }
    }
    // ... rest of implementation
}
```

## Performance Considerations

### Optimization Guidelines

1. **Measure first**: Use benchmarks to identify bottlenecks
2. **Maintain correctness**: Never sacrifice IEEE 754 compliance for speed
3. **Consider memory allocation**: Minimize heap allocations in hot paths
4. **Use bit operations**: Leverage efficient bit manipulation
5. **Profile regularly**: Use `go tool pprof` for detailed analysis

### Benchmark Requirements

Add benchmarks for new performance-critical code:

```go
func BenchmarkNewOperation(b *testing.B) {
    f1 := FromFloat32(3.14159)
    f2 := FromFloat32(2.71828)
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _ = f1.NewOperation(f2)
    }
}

func BenchmarkNewOperationParallel(b *testing.B) {
    f1 := FromFloat32(3.14159)
    f2 := FromFloat32(2.71828)
    
    b.RunParallel(func(pb *testing.PB) {
        for pb.Next() {
            _ = f1.NewOperation(f2)
        }
    })
}
```

## Submitting Changes

### Pull Request Process

1. **Create a feature branch:**
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. **Make your changes** following the guidelines above

3. **Add tests** for new functionality

4. **Run the full test suite:**
   ```bash
   go test ./...
   go test -race ./...
   go vet ./...
   ```

5. **Update documentation** if needed

6. **Commit with clear messages:**
   ```bash
   git commit -m "Add support for new rounding mode
   
   - Implement RoundNearestAway mode
   - Add comprehensive tests
   - Update documentation
   - Maintain IEEE 754 compliance"
   ```

7. **Push and create pull request:**
   ```bash
   git push origin feature/your-feature-name
   ```

### Pull Request Guidelines

- **Clear title and description** explaining the change
- **Reference related issues** using "Fixes #123" or "Closes #123"
- **Include test results** showing all tests pass
- **Update CHANGELOG** if applicable
- **Keep changes focused** - one feature/fix per PR
- **Respond to feedback** promptly and professionally

### Review Process

1. **Automated checks** must pass (tests, linting, formatting)
2. **Code review** by maintainers
3. **Performance review** for optimization changes
4. **Documentation review** for user-facing changes
5. **Final approval** and merge

## Issue Reporting

### Bug Reports

Include the following information:

```markdown
**Describe the bug**
A clear description of what the bug is.

**To Reproduce**
Steps to reproduce the behavior:
1. Create Float16 value '...'
2. Call method '...'
3. See error

**Expected behavior**
What you expected to happen.

**Actual behavior**
What actually happened.

**Environment**
- Go version: [e.g. 1.24]
- OS: [e.g. macOS, Linux, Windows]
- Architecture: [e.g. amd64, arm64]

**Additional context**
Any other context about the problem.
```

### Feature Requests

```markdown
**Is your feature request related to a problem?**
A clear description of what the problem is.

**Describe the solution you'd like**
A clear description of what you want to happen.

**Describe alternatives you've considered**
Any alternative solutions or features you've considered.

**Additional context**
Any other context or screenshots about the feature request.
```

## Documentation

### Code Documentation

- **All exported functions** must have documentation comments
- **Include examples** for complex operations
- **Document special cases** and edge behaviors
- **Explain IEEE 754 compliance** where relevant

### README Updates

When adding new features:

1. Update the feature list
2. Add usage examples
3. Update the table of contents
4. Consider adding to quick start guide

### Example Documentation

```go
// Add performs IEEE 754 compliant addition of two Float16 values.
//
// Special cases are:
//   - Add(x, ±0) = x for any x
//   - Add(±0, ±0) = +0
//   - Add(x, -x) = +0
//   - Add(±Inf, ±Inf) = ±Inf
//   - Add(±Inf, ∓Inf) = NaN
//   - Add(x, NaN) = NaN for any x
//
// Example:
//   a := FromFloat32(1.5)
//   b := FromFloat32(2.5)
//   result := a.Add(b) // result represents 4.0
func (f Float16) Add(other Float16) Float16 {
    // Implementation
}
```

## Questions?

If you have questions about contributing:

1. Check existing issues and discussions
2. Open a new issue with the "question" label
3. Join our community discussions (if applicable)

Thank you for contributing to float16!
