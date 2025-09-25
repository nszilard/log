# Log

[![CI](https://github.com/nszilard/log/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/nszilard/log/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/nszilard/log)](https://goreportcard.com/report/github.com/nszilard/log)
[![GoDoc](https://godoc.org/github.com/nszilard/log?status.svg)](https://godoc.org/github.com/nszilard/log)

A fast, simple, and structured logger for Go applications with minimal overhead and maximum flexibility.

## Installation

```shell
go get github.com/nszilard/log
```

## Features

- **Simple API**: Three variants per level (basic, formatted, structured)
- **Structured Logging**: Key-value pairs for better log analysis
- **Level-based Filtering**: Only log what you need
- **High Performance**: Optimized for speed with minimal allocations
- **Thread Safe**: Concurrent logging without data races
- **Testing Support**: Noop logger for performance-critical paths

## Log Levels

```go
const (
    PanicLevel Level = iota  // Fatal errors and panics
    FatalLevel              // Fatal errors, doesn't panic
    ErrorLevel              // Error conditions
    WarnLevel               // Warning conditions
    InfoLevel               // Informational messages
    DebugLevel              // Debug information
)
```

## Quick Start

### Basic Usage

```go
package main

import (
    "os"
    "github.com/nszilard/log"
)

func main() {
    // Use package-level logger (writes to stdout by default)
    log.Info("Application started")
    log.Error("Something went wrong")

    // Create a custom logger
    logger := log.New(log.InfoLevel, os.Stdout)
    logger.Info("Custom logger message")
}
```

### Three Logging Variants

Each level supports three variants:

```go
// Basic: Log any values
log.Info("User", "john", "logged in")

// Formatted: Printf-style formatting
log.Infof("User %s logged in at %s", "john", time.Now())

// Structured: Key-value pairs for structured logs
log.InfoS(
    log.WithString("user", "john"),
    log.WithTime("timestamp", time.Now()),
    log.WithBool("success", true),
)
```

### Structured Logging

```go
log.InfoS(
    log.WithInt("user_id", 12345),
    log.WithString("action", "login"),
    log.WithInt("duration_ms", 245),
    log.WithString("ip_address", "192.168.1.100"),
)
// Output: {"timestamp":"2025-09-25T13:20:18.524Z","level":"INFO","caller":"main.go:17","user_id":12345,"action":"login","duration_ms":245,"ip_address":"192.168.1.100"}
```

### Level Configuration

```go
logger := log.New(log.WarnLevel, os.Stdout)

logger.Debug("Won't appear") // Filtered out
logger.Info("Won't appear")  // Filtered out
logger.Warn("Will appear")   // Logged
logger.Error("Will appear")  // Logged

// Change level at runtime
logger.SetLevel(log.DebugLevel)
logger.Debug("Now this appears") // Logged
```

## Performance

Benchmarks on Apple M2 Pro:

```
BenchmarkStructuredLogging/WithFileInfo-10         	 1018492	        1116 ns/op	     720 B/op	      12 allocs/op
BenchmarkStructuredLogging/WithoutFileInfo-10      	 1662369	       720.0 ns/op	     472 B/op	      10 allocs/op
BenchmarkFormattedLogging/WithFileInfo-10          	 2055763	       586.0 ns/op	     416 B/op	       6 allocs/op
BenchmarkFormattedLogging/WithoutFileInfo-10       	 4146776	       290.4 ns/op	     168 B/op	       4 allocs/op
BenchmarkFieldTypes/String-10                      	 5654694	       214.2 ns/op	      24 B/op	       1 allocs/op
BenchmarkFieldTypes/Int-10                         	 5544013	       215.6 ns/op	      24 B/op	       1 allocs/op
BenchmarkFieldTypes/Float-10                       	 4701655	       252.8 ns/op	      24 B/op	       1 allocs/op
BenchmarkFieldTypes/Bool-10                        	 5687622	       210.2 ns/op	      24 B/op	       1 allocs/op
BenchmarkFieldTypes/Error-10                       	 5465964	       220.3 ns/op	      24 B/op	       1 allocs/op
BenchmarkFieldTypes/Duration-10                    	 5369445	       223.6 ns/op	      24 B/op	       1 allocs/op
BenchmarkFiltered-10                               	615557968	       1.920 ns/op	       0 B/op	       0 allocs/op
BenchmarkMixedFields-10                            	 3519652	       343.9 ns/op	      24 B/op	       1 allocs/op
BenchmarkLargePayload-10                           	 1943682	       610.9 ns/op	      24 B/op	       1 allocs/op
BenchmarkConcurrent-10                             	 6142543	       196.0 ns/op	      24 B/op	       1 allocs/op
BenchmarkAllocations-10                            	 3739988	       321.2 ns/op	      24 B/op	       1 allocs/op
```

## Thread Safety

All logger methods are thread-safe and can be called concurrently from multiple goroutines without additional synchronization.

## Architecture

- **Simplified Design**: No separate models package - everything in main package
- **Clean Separation**: Internal package contains only core logging logic
- **Interface-based**: Easy to swap implementations (real, mock, noop)
- **Minimal Dependencies**: Only standard library dependencies

## License

MIT License - see LICENSE file for details.
