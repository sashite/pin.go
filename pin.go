package pin

import (
	"errors"
	"fmt"
	"unicode"
)

// Side represents the player owning the piece (First or Second).
type Side int8

const (
	First  Side = iota // First player (represented by Uppercase)
	Second             // Second player (represented by Lowercase)
)

// State represents the mutation state of the piece (Normal, Enhanced, Diminished).
type State int8

const (
	Normal     State = iota // No prefix
	Enhanced                // Prefix '+'
	Diminished              // Prefix '-'
)

// Pin is an immutable representation of a Piece Identifier.
// It holds the core attributes defined in the PIN specification.
type Pin struct {
	// Abbr is the piece name abbreviation.
	// It should ideally be stored as Uppercase (the piece "Type").
	// The visual case (upper/lower) is determined by the Side field during formatting.
	Abbr     rune
	Side     Side
	State    State
	Terminal bool
}

// Common errors
var (
	ErrEmpty       = errors.New("pin: empty string")
	ErrInvalid     = errors.New("pin: invalid format")
	ErrInvalidAbbr = errors.New("pin: invalid abbreviation (must be a letter)")
)

// Parse converts a PIN string into a Pin structure.
// It enforces strict adherence to the spec: [<state>]<letter>[<terminal>].
func Parse(s string) (Pin, error) {
	if len(s) == 0 {
		return Pin{}, ErrEmpty
	}

	var p Pin
	chars := []rune(s)
	cursor := 0
	length := len(chars)

	// 1. Parse State Modifier (Optional Prefix)
	if length > cursor {
		switch chars[cursor] {
		case '+':
			p.State = Enhanced
			cursor++
		case '-':
			p.State = Diminished
			cursor++
		}
	}

	// 2. Parse Abbreviation (Mandatory)
	if cursor >= length {
		return Pin{}, ErrInvalid
	}

	r := chars[cursor]
	if !unicode.IsLetter(r) {
		return Pin{}, ErrInvalidAbbr
	}

	// Determine Side based on case, then normalize Abbr to Uppercase
	if unicode.IsLower(r) {
		p.Side = Second
		p.Abbr = unicode.ToUpper(r)
	} else {
		p.Side = First
		p.Abbr = r
	}
	cursor++

	// 3. Parse Terminal Marker (Optional Suffix)
	if cursor < length {
		if chars[cursor] == '^' {
			p.Terminal = true
			cursor++
		} else {
			// If there is a character after the letter that is not '^', it's invalid
			// (e.g., "K+" or "KK")
			return Pin{}, ErrInvalid
		}
	}

	// 4. Validation: Ensure no extra characters remain
	if cursor != length {
		return Pin{}, ErrInvalid
	}

	return p, nil
}

// MustParse is a helper that panics if parsing fails.
// Useful for initializing global constants or in tests.
func MustParse(s string) Pin {
	p, err := Parse(s)
	if err != nil {
		panic(fmt.Sprintf("pin: MustParse(%q) failed: %v", s, err))
	}
	return p
}

// String implements the fmt.Stringer interface.
// It returns the standard string representation of the PIN.
func (p Pin) String() string {
	// Pre-allocate a buffer of max length 4 (+, A, ^, extra safety)
	// to avoid multiple allocations, though AppendTo handles logic.
	buf := make([]byte, 0, 4)
	return string(p.AppendTo(buf))
}

// AppendTo appends the string representation of the PIN to dst.
// This is the zero-allocation primitive for high-performance serialization.
func (p Pin) AppendTo(dst []byte) []byte {
	// 1. State
	switch p.State {
	case Enhanced:
		dst = append(dst, '+')
	case Diminished:
		dst = append(dst, '-')
	}

	// 2. Abbr (Case determined by Side)
	r := p.Abbr
	if p.Side == Second {
		r = unicode.ToLower(r)
	} else {
		r = unicode.ToUpper(r)
	}
	dst = append(dst, string(r)...)

	// 3. Terminal
	if p.Terminal {
		dst = append(dst, '^')
	}

	return dst
}

// -----------------------------------------------------------------------------
// Transformations (Fluent API)
// These methods return a new Pin, preserving immutability.
// -----------------------------------------------------------------------------

// WithSide returns a new Pin with the specified Side.
func (p Pin) WithSide(s Side) Pin {
	p.Side = s
	return p
}

// WithState returns a new Pin with the specified State.
func (p Pin) WithState(s State) Pin {
	p.State = s
	return p
}

// WithTerminal returns a new Pin with the specified Terminal status.
func (p Pin) WithTerminal(t bool) Pin {
	p.Terminal = t
	return p
}

// -----------------------------------------------------------------------------
// Convenience Helpers
// -----------------------------------------------------------------------------

// Flip returns a new Pin with the Side inverted.
// First -> Second, Second -> First.
func (p Pin) Flip() Pin {
	if p.Side == First {
		p.Side = Second
	} else {
		p.Side = First
	}
	return p
}

// Enhance returns a new Pin promoted to the Enhanced state.
func (p Pin) Enhance() Pin {
	p.State = Enhanced
	return p
}

// Diminish returns a new Pin demoted to the Diminished state.
func (p Pin) Diminish() Pin {
	p.State = Diminished
	return p
}
