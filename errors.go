package pin

import "errors"

// Parsing errors.
var (
	// ErrEmptyInput is returned when the input string is empty.
	ErrEmptyInput = errors.New("pin: empty input")

	// ErrInputTooLong is returned when the input exceeds MaxStringLength.
	ErrInputTooLong = errors.New("pin: input exceeds 3 characters")

	// ErrMustContainOneLetter is returned when the input does not contain exactly one letter.
	ErrMustContainOneLetter = errors.New("pin: must contain exactly one letter")

	// ErrInvalidStateModifier is returned when the state modifier is invalid.
	ErrInvalidStateModifier = errors.New("pin: invalid state modifier")

	// ErrInvalidTerminalMarker is returned when the terminal marker is invalid.
	ErrInvalidTerminalMarker = errors.New("pin: invalid terminal marker")
)

// Validation errors (for constructors).
var (
	// ErrInvalidType is returned when the piece type is not A-Z.
	ErrInvalidType = errors.New("pin: invalid type (must be A-Z)")

	// ErrInvalidSide is returned when the side is not First or Second.
	ErrInvalidSide = errors.New("pin: invalid side")

	// ErrInvalidState is returned when the state is not Normal, Enhanced, or Diminished.
	ErrInvalidState = errors.New("pin: invalid state")
)
