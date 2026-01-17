# pin.go

[![Go Reference](https://pkg.go.dev/badge/github.com/sashite/pin.go.svg)](https://pkg.go.dev/github.com/sashite/pin.go)
[![Go Report Card](https://goreportcard.com/badge/github.com/sashite/pin.go)](https://goreportcard.com/report/github.com/sashite/pin.go)
[![License](https://img.shields.io/github/license/sashite/pin.go)](https://github.com/sashite/pin.go/blob/main/LICENSE)

> **PIN** (Piece Identifier Notation) implementation for Go.

## Overview

This library implements the [PIN Specification v1.0.0](https://sashite.dev/specs/pin/1.0.0/).

## Installation

```bash
go get github.com/sashite/pin.go/v2
```

## Usage

### Parsing (String → Identifier)

Convert a PIN string into an `Identifier` struct.

```go
package main

import (
	"fmt"
	"github.com/sashite/pin.go/v2"
)

func main() {
	// Standard parsing (returns error)
	id, err := pin.Parse("K")
	if err != nil {
		panic(err)
	}
	fmt.Println(id.Type())     // K
	fmt.Println(id.Side())     // First
	fmt.Println(id.State())    // Normal
	fmt.Println(id.Terminal()) // false

	// With state modifier
	id, _ = pin.Parse("+R")
	fmt.Println(id.State()) // Enhanced

	// With terminal marker
	id, _ = pin.Parse("K^")
	fmt.Println(id.Terminal()) // true

	// Panic on error (for constants or trusted input)
	k := pin.MustParse("+K^")
	fmt.Println(k.String()) // +K^
}
```

### Formatting (Identifier → String)

Convert an `Identifier` back to a PIN string.

```go
// From Identifier
id := pin.NewIdentifier('K', pin.First)
fmt.Println(id.String()) // "K"

// With attributes
id = pin.NewIdentifierWithOptions('R', pin.Second, pin.Enhanced, false)
fmt.Println(id.String()) // "+r"

id = pin.NewIdentifierWithOptions('K', pin.First, pin.Normal, true)
fmt.Println(id.String()) // "K^"
```

### Validation

```go
// Boolean check
if pin.IsValid("K") {
	// valid identifier
}

// Detailed error
if err := pin.Validate("+K^"); err != nil {
	fmt.Println(err)
}
```

### Accessing Identifier Data

```go
id := pin.MustParse("+K^")

// Get attributes
fmt.Println(id.Type())     // K (rune)
fmt.Println(id.Side())     // First
fmt.Println(id.State())    // Enhanced
fmt.Println(id.Terminal()) // true

// Get string components
fmt.Println(id.Letter()) // "K"
fmt.Println(id.Prefix()) // "+"
fmt.Println(id.Suffix()) // "^"
```

### Transformations

All transformations return new immutable values.

```go
id := pin.MustParse("K")

// State transformations
fmt.Println(id.Enhance().String())   // "+K"
fmt.Println(id.Diminish().String())  // "-K"
fmt.Println(id.Normalize().String()) // "K"

// Side transformation
fmt.Println(id.Flip().String()) // "k"

// Terminal transformations
fmt.Println(id.MarkTerminal().String())   // "K^"
fmt.Println(id.UnmarkTerminal().String()) // "K"

// Attribute changes
fmt.Println(id.WithType('Q').String())           // "Q"
fmt.Println(id.WithSide(pin.Second).String())    // "k"
fmt.Println(id.WithState(pin.Enhanced).String()) // "+K"
fmt.Println(id.WithTerminal(true).String())      // "K^"
```

### Queries

```go
id := pin.MustParse("+K^")

// State queries
fmt.Println(id.IsNormal())     // false
fmt.Println(id.IsEnhanced())   // true
fmt.Println(id.IsDiminished()) // false

// Side queries
fmt.Println(id.IsFirstPlayer())  // true
fmt.Println(id.IsSecondPlayer()) // false

// Terminal query
fmt.Println(id.Terminal()) // true

// Comparison queries
other := pin.MustParse("k")
fmt.Println(id.SameType(other))     // true
fmt.Println(id.SameSide(other))     // false
fmt.Println(id.SameState(other))    // false
fmt.Println(id.SameTerminal(other)) // false
```

### Zero-Allocation Serialization

For high-performance scenarios, use `AppendTo` to avoid allocations.

```go
buf := make([]byte, 0, 32)
id := pin.MustParse("+K^")
buf = id.AppendTo(buf)
fmt.Printf("%s\n", buf) // "+K^"
```

## API Reference

### Types

```go
// Identifier represents a parsed PIN identifier with all attributes.
// Zero value is not valid; use NewIdentifier or Parse to create.
type Identifier struct {
	// contains unexported fields
}

// Side represents the player side.
type Side uint8

// State represents the piece state.
type State uint8

// NewIdentifier creates an Identifier with default state (Normal) and terminal (false).
func NewIdentifier(typ rune, side Side) Identifier

// NewIdentifierWithOptions creates an Identifier with all attributes specified.
func NewIdentifierWithOptions(typ rune, side Side, state State, terminal bool) Identifier

// Type returns the piece type as uppercase rune (A-Z).
func (id Identifier) Type() rune

// Side returns the player side.
func (id Identifier) Side() Side

// State returns the piece state.
func (id Identifier) State() State

// Terminal returns the terminal status.
func (id Identifier) Terminal() bool

// String returns the PIN string representation.
func (id Identifier) String() string

// AppendTo appends the PIN string to dst without allocation.
func (id Identifier) AppendTo(dst []byte) []byte
```

### Constants

```go
const (
	First  Side = iota // Uppercase letter
	Second             // Lowercase letter
)

const (
	Normal     State = iota // No prefix
	Enhanced                // Prefix '+'
	Diminished              // Prefix '-'
)

const MaxStringLength = 3
```

### Parsing

```go
// Parse converts a PIN string to an Identifier.
// Returns an error if the string is not valid.
func Parse(s string) (Identifier, error)

// MustParse is like Parse but panics on error.
// Use for constants or trusted input.
func MustParse(s string) Identifier
```

### Validation

```go
// Validate checks if s is a valid PIN identifier.
// Returns nil if valid, or a descriptive error.
func Validate(s string) error

// IsValid reports whether s is a valid PIN identifier.
func IsValid(s string) bool
```

### Transformations

```go
// State transformations
func (id Identifier) Enhance() Identifier
func (id Identifier) Diminish() Identifier
func (id Identifier) Normalize() Identifier

// Side transformation
func (id Identifier) Flip() Identifier

// Terminal transformations
func (id Identifier) MarkTerminal() Identifier
func (id Identifier) UnmarkTerminal() Identifier

// Attribute changes
func (id Identifier) WithType(typ rune) Identifier
func (id Identifier) WithSide(side Side) Identifier
func (id Identifier) WithState(state State) Identifier
func (id Identifier) WithTerminal(terminal bool) Identifier
```

### Queries

```go
// State queries
func (id Identifier) IsNormal() bool
func (id Identifier) IsEnhanced() bool
func (id Identifier) IsDiminished() bool

// Side queries
func (id Identifier) IsFirstPlayer() bool
func (id Identifier) IsSecondPlayer() bool

// Comparison queries
func (id Identifier) SameType(other Identifier) bool
func (id Identifier) SameSide(other Identifier) bool
func (id Identifier) SameState(other Identifier) bool
func (id Identifier) SameTerminal(other Identifier) bool
```

### Errors

```go
var (
	ErrEmptyInput           = errors.New("pin: empty input")
	ErrInputTooLong         = errors.New("pin: input exceeds 3 characters")
	ErrMustContainOneLetter = errors.New("pin: must contain exactly one letter")
	ErrInvalidStateModifier = errors.New("pin: invalid state modifier")
	ErrInvalidTerminalMarker = errors.New("pin: invalid terminal marker")
)
```

## Design Principles

- **Bounded types**: Fixed-size struct, no heap allocation in hot path
- **Value semantics**: `Identifier` is a value type, safe to copy
- **Sentinel errors**: Standard Go error handling with `errors.Is()`
- **strconv-style API**: Familiar `Parse`, `MustParse`, `String()` patterns
- **Zero-allocation option**: `AppendTo` for high-performance serialization
- **Security-hardened**: Byte-level parsing, rejects Unicode lookalikes
- **No dependencies**: Pure Go standard library only

## Related Specifications

- [Game Protocol](https://sashite.dev/game-protocol/) — Conceptual foundation
- [PIN Specification](https://sashite.dev/specs/pin/1.0.0/) — Official specification
- [PIN Examples](https://sashite.dev/specs/pin/1.0.0/examples/) — Usage examples

## License

Available as open source under the [Apache License 2.0](https://opensource.org/licenses/Apache-2.0).
