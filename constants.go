// Package pin implements the PIN (Piece Identifier Notation) specification.
//
// PIN provides a compact, ASCII-based format for representing pieces in
// abstract strategy board games.
//
// See https://sashite.dev/specs/pin/1.0.0/ for the specification.
package pin

// Side represents the player side (First or Second).
type Side uint8

const (
	// First represents the first player (uppercase letters).
	First Side = iota
	// Second represents the second player (lowercase letters).
	Second
)

// State represents the piece state (Normal, Enhanced, or Diminished).
type State uint8

const (
	// Normal represents a piece with no state modifier.
	Normal State = iota
	// Enhanced represents a piece with the '+' prefix.
	Enhanced
	// Diminished represents a piece with the '-' prefix.
	Diminished
)

// Limits
const (
	// MaxStringLength is the maximum length of a valid PIN string.
	MaxStringLength = 3
)

// Formatting constants (internal use)
const (
	enhancedPrefix   = '+'
	diminishedPrefix = '-'
	terminalSuffix   = '^'
)

// String returns the string representation of the Side.
func (s Side) String() string {
	switch s {
	case First:
		return "First"
	case Second:
		return "Second"
	default:
		return "Unknown"
	}
}

// String returns the string representation of the State.
func (s State) String() string {
	switch s {
	case Normal:
		return "Normal"
	case Enhanced:
		return "Enhanced"
	case Diminished:
		return "Diminished"
	default:
		return "Unknown"
	}
}

// isValidSide reports whether s is a valid Side value.
func isValidSide(s Side) bool {
	return s == First || s == Second
}

// isValidState reports whether s is a valid State value.
func isValidState(s State) bool {
	return s == Normal || s == Enhanced || s == Diminished
}

// isValidType reports whether r is a valid piece type (A-Z).
func isValidType(r rune) bool {
	return r >= 'A' && r <= 'Z'
}
