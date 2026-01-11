# pin.go

[![Go Reference](https://pkg.go.dev/badge/github.com/sashite/pin.go.svg)](https://pkg.go.dev/github.com/sashite/pin.go)
[![Go Report Card](https://goreportcard.com/badge/github.com/sashite/pin.go)](https://goreportcard.com/report/github.com/sashite/pin.go)
[![License](https://img.shields.io/github/license/sashite/pin.go)](https://github.com/sashite/pin.go/blob/main/LICENSE)

> Idiomatic Go implementation of the **PIN** (Piece Identifier Notation) specification.

## What is PIN?

PIN (Piece Identifier Notation) is a compact, ASCII-based format for representing pieces in abstract strategy board games. It translates piece attributes from the [Game Protocol](https://sashite.dev/game-protocol/) into a portable notation system.

This library implements the [PIN Specification v1.0.0](https://sashite.dev/specs/pin/1.0.0/). It is designed to be performant (zero-allocation manipulations) and safe (immutable value semantics).

## Installation

```bash
go get github.com/sashite/pin.go

```

## Usage

### 1. Parsing (String to Pin)

Use `Parse` to convert a PIN string into a `Pin` structure.

```go
package main

import (
	"fmt"
	"github.com/sashite/pin.go"
)

func main() {
	// Standard parsing
	p, err := pin.Parse("+K^")
	if err != nil {
		panic(err)
	}
	fmt.Println(p.String()) // Output: "+K^"

	// Strict parsing (panics on error) - Useful for constants
	k := pin.MustParse("k")
	fmt.Println(k.Side) // Output: Second (1)
}

```

### 2. Manual Creation

You can create `Pin` values using struct literals. This is the most idiomatic way to initialize known values.

```go
// Via struct literal (explicit and clear)
p := pin.Pin{
	Abbr:     'n',        // Note: strictly a rune (int32), matches spec <abbr>
	Side:     pin.Second,
	State:    pin.Enhanced,
	Terminal: true,
}
fmt.Println(p.String()) // Output: "+n^"

```

### 3. Immutability and Manipulation

`Pin` objects are **immutable values**. Transformation methods always return a *new* copy, leaving the original unmodified.

The API provides "Fluent Setters" (`WithX`) for absolute control, and "Action Helpers" (`Flip`, `Enhance`) for relative changes.

```go
origin := pin.MustParse("P") // Normal white pawn "P"

// Fluent setters (Consistent with Sashité Ruby API)
promoted := origin.WithState(pin.Enhanced).WithTerminal(true)

// Action helpers
blackPawn := origin.Flip()

fmt.Println(origin.String())    // "P"   (Original remains unchanged)
fmt.Println(promoted.String())  // "+P^" (New value)
fmt.Println(blackPawn.String()) // "p"

```

### 4. High Performance (Zero-Allocation)

For game engines requiring high-throughput serialization (e.g., generating FEN strings or logging moves), use `AppendTo` instead of `String`. This avoids allocating new string objects by appending directly to a byte buffer.

```go
// Pre-allocate a buffer
buf := make([]byte, 0, 16)

// In your game loop:
pieces := []pin.Pin{
    {Abbr: 'R', Side: pin.First},
    {Abbr: 'k', Side: pin.Second, Terminal: true},
}

for _, p := range pieces {
    // Append directly to buffer (no allocation)
    buf = p.AppendTo(buf)
    buf = append(buf, ' ') // Separator
}

fmt.Printf("%s", buf) // Output: "R k^ "

```

### 5. Inspection

The `Pin` struct fields are public, allowing for direct, high-performance access.

```go
p := pin.MustParse("+r^")

// Check attributes directly
if p.Terminal {
	fmt.Println("Game ends if this piece is captured.")
}

if p.State == pin.Enhanced {
	fmt.Println("This piece is in an enhanced state.")
}

// Abbr is a rune, making Unicode comparisons easy
if p.Abbr == 'r' {
	fmt.Println("It is a rook/chariot.")
}

```

## API Reference

### Types

```go
type Pin struct {
	Abbr     rune  // ASCII character (e.g., 'K', 'p')
	Side     Side  // First (0) or Second (1)
	State    State // Normal (0), Enhanced (1), Diminished (2)
	Terminal bool  // true or false
}

```

### Constants

* **Side**: `pin.First` (Uppercase), `pin.Second` (Lowercase).
* **State**: `pin.Normal`, `pin.Enhanced` (`+`), `pin.Diminished` (`-`).

### Methods

```go
// Parsing
func Parse(s string) (Pin, error)
func MustParse(s string) Pin

// Stringer Interface
func (p Pin) String() string

// Zero-Allocation Serialization
func (p Pin) AppendTo(dst []byte) []byte

// Transformations (Return new Pin)
func (p Pin) WithSide(s Side) Pin
func (p Pin) WithState(s State) Pin
func (p Pin) WithTerminal(t bool) Pin

// Convenience Helpers
func (p Pin) Flip() Pin           // Toggles Side
func (p Pin) Enhance() Pin        // Normal -> Enhanced
func (p Pin) Diminish() Pin       // Normal -> Diminished

```

## Related Specifications

* [Game Protocol](https://sashite.dev/game-protocol/) — Conceptual foundation
* [PIN Specification](https://sashite.dev/specs/pin/1.0.0/) — Official specification
* [PIN Examples](https://sashite.dev/specs/pin/1.0.0/examples/) — Game-specific mappings

## License

Available as open source under the [Apache License 2.0](https://opensource.org/licenses/Apache-2.0).
