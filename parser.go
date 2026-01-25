package pin

// Parse converts a PIN string into an Identifier.
//
// The parser uses byte-level validation to ensure security against
// malformed input, Unicode lookalikes, and injection attacks.
//
// Valid formats:
//   - Single letter: "K", "k"
//   - With state modifier: "+K", "-k"
//   - With terminal marker: "K^", "k^"
//   - Combined: "+K^", "-k^"
//
// Returns an error if the string is not valid:
//   - ErrEmptyInput: empty string
//   - ErrInputTooLong: exceeds 3 characters
//   - ErrMustContainOneLetter: no letter found
//   - ErrInvalidStateModifier: invalid prefix character
//   - ErrInvalidTerminalMarker: invalid suffix character
func Parse(s string) (Identifier, error) {
	// Validate input length
	if len(s) == 0 {
		return Identifier{}, ErrEmptyInput
	}
	if len(s) > MaxStringLength {
		return Identifier{}, ErrInputTooLong
	}

	// Convert to bytes for safe parsing
	// This also ensures we reject multi-byte UTF-8 characters
	bytes := []byte(s)

	// Dispatch based on length
	switch len(bytes) {
	case 1:
		return parseLength1(bytes[0])
	case 2:
		return parseLength2(bytes[0], bytes[1])
	case 3:
		return parseLength3(bytes[0], bytes[1], bytes[2])
	default:
		// Should not reach here due to length check above
		return Identifier{}, ErrInputTooLong
	}
}

// parseLength1 handles single-byte input (letter only).
func parseLength1(b byte) (Identifier, error) {
	abbr, side, ok := classifyLetter(b)
	if !ok {
		return Identifier{}, ErrMustContainOneLetter
	}

	return Identifier{
		abbr:     abbr,
		side:     side,
		state:    Normal,
		terminal: false,
	}, nil
}

// parseLength2 handles two-byte input (modifier+letter or letter+terminal).
func parseLength2(first, second byte) (Identifier, error) {
	// Try: modifier + letter
	if state, ok := classifyModifier(first); ok {
		abbr, side, ok := classifyLetter(second)
		if !ok {
			return Identifier{}, ErrMustContainOneLetter
		}
		return Identifier{
			abbr:     abbr,
			side:     side,
			state:    state,
			terminal: false,
		}, nil
	}

	// Try: letter + terminal
	abbr, side, ok := classifyLetter(first)
	if !ok {
		// First byte is not a letter and not a modifier
		if isTerminalMarker(first) {
			return Identifier{}, ErrInvalidStateModifier
		}
		return Identifier{}, ErrInvalidStateModifier
	}

	if !isTerminalMarker(second) {
		return Identifier{}, ErrInvalidTerminalMarker
	}

	return Identifier{
		abbr:     abbr,
		side:     side,
		state:    Normal,
		terminal: true,
	}, nil
}

// parseLength3 handles three-byte input (modifier+letter+terminal).
func parseLength3(first, second, third byte) (Identifier, error) {
	// Must be: modifier + letter + terminal
	state, ok := classifyModifier(first)
	if !ok {
		// First byte is not a valid modifier
		if _, _, isLetter := classifyLetter(first); isLetter {
			return Identifier{}, ErrInvalidTerminalMarker
		}
		return Identifier{}, ErrInvalidStateModifier
	}

	abbr, side, ok := classifyLetter(second)
	if !ok {
		return Identifier{}, ErrMustContainOneLetter
	}

	if !isTerminalMarker(third) {
		return Identifier{}, ErrInvalidTerminalMarker
	}

	return Identifier{
		abbr:     abbr,
		side:     side,
		state:    state,
		terminal: true,
	}, nil
}

// classifyLetter checks if a byte is a valid ASCII letter.
// Returns the uppercase abbreviation, side, and whether it's valid.
func classifyLetter(b byte) (rune, Side, bool) {
	switch {
	case b >= 'A' && b <= 'Z':
		return rune(b), First, true
	case b >= 'a' && b <= 'z':
		return rune(b - 'a' + 'A'), Second, true
	default:
		return 0, 0, false
	}
}

// classifyModifier checks if a byte is a valid state modifier.
// Returns the state and whether it's valid.
func classifyModifier(b byte) (State, bool) {
	switch b {
	case '+':
		return Enhanced, true
	case '-':
		return Diminished, true
	default:
		return Normal, false
	}
}

// isTerminalMarker checks if a byte is the terminal marker.
func isTerminalMarker(b byte) bool {
	return b == '^'
}

// MustParse is like Parse but panics on error.
// Use for constants or trusted input.
func MustParse(s string) Identifier {
	id, err := Parse(s)
	if err != nil {
		panic("pin: MustParse(" + s + "): " + err.Error())
	}
	return id
}

// Validate checks if s is a valid PIN identifier.
// Returns nil if valid, or a descriptive error.
func Validate(s string) error {
	_, err := Parse(s)
	return err
}

// IsValid reports whether s is a valid PIN identifier.
func IsValid(s string) bool {
	_, err := Parse(s)
	return err == nil
}
