package pin

// Identifier represents a parsed PIN (Piece Identifier Notation) identifier.
//
// An Identifier encodes four attributes of a piece:
//   - Abbr: the piece name abbreviation (A-Z as uppercase rune)
//   - Side: the piece side (First or Second)
//   - State: the piece state (Normal, Enhanced, or Diminished)
//   - Terminal: whether the piece is terminal (true or false)
//
// Identifier is an immutable value type. All transformation methods return
// new Identifier values, leaving the original unchanged.
//
// The zero value is not valid; use NewIdentifier, NewIdentifierWithOptions,
// or Parse to create valid instances.
type Identifier struct {
	abbr     rune
	side     Side
	state    State
	terminal bool
}

// ============================================================================
// Constructors
// ============================================================================

// NewIdentifier creates a new Identifier with the given abbreviation and side.
// State defaults to Normal and terminal defaults to false.
//
// The abbreviation must be an uppercase letter (A-Z). If a lowercase letter is provided,
// it will be converted to uppercase.
//
// Panics if the abbreviation is not a letter A-Z (case-insensitive) or if the side is invalid.
func NewIdentifier(abbr rune, side Side) Identifier {
	return NewIdentifierWithOptions(abbr, side, Normal, false)
}

// NewIdentifierWithOptions creates a new Identifier with all attributes specified.
//
// The abbreviation must be an uppercase letter (A-Z). If a lowercase letter is provided,
// it will be converted to uppercase.
//
// Panics if any attribute is invalid.
func NewIdentifierWithOptions(abbr rune, side Side, state State, terminal bool) Identifier {
	// Normalize to uppercase
	if abbr >= 'a' && abbr <= 'z' {
		abbr = abbr - 'a' + 'A'
	}

	if !isValidAbbr(abbr) {
		panic(ErrInvalidAbbr)
	}
	if !isValidSide(side) {
		panic(ErrInvalidSide)
	}
	if !isValidState(state) {
		panic(ErrInvalidState)
	}

	return Identifier{
		abbr:     abbr,
		side:     side,
		state:    state,
		terminal: terminal,
	}
}

// ============================================================================
// Accessors
// ============================================================================

// Abbr returns the piece name abbreviation as an uppercase rune (A-Z).
func (id Identifier) Abbr() rune {
	return id.abbr
}

// Side returns the piece side.
func (id Identifier) Side() Side {
	return id.side
}

// State returns the piece state.
func (id Identifier) State() State {
	return id.state
}

// IsTerminal returns the terminal status.
func (id Identifier) IsTerminal() bool {
	return id.terminal
}

// ============================================================================
// String Conversion
// ============================================================================

// String returns the PIN string representation.
//
// Examples:
//
//	NewIdentifier('K', First).String()                              // "K"
//	NewIdentifierWithOptions('R', Second, Enhanced, false).String() // "+r"
//	NewIdentifierWithOptions('K', First, Normal, true).String()     // "K^"
func (id Identifier) String() string {
	buf := make([]byte, 0, MaxStringLength)
	return string(id.AppendTo(buf))
}

// AppendTo appends the PIN string representation to dst and returns the result.
// This is the zero-allocation primitive for high-performance serialization.
func (id Identifier) AppendTo(dst []byte) []byte {
	// 1. State prefix
	switch id.state {
	case Enhanced:
		dst = append(dst, enhancedPrefix)
	case Diminished:
		dst = append(dst, diminishedPrefix)
	}

	// 2. Letter (case determined by side)
	dst = append(dst, id.Letter()...)

	// 3. Terminal suffix
	if id.terminal {
		dst = append(dst, terminalSuffix)
	}

	return dst
}

// Letter returns the letter component of the PIN.
// Returns uppercase for First player, lowercase for Second player.
func (id Identifier) Letter() string {
	r := id.abbr
	if id.side == Second {
		r = r - 'A' + 'a'
	}
	return string(r)
}

// Prefix returns the state prefix of the PIN.
// Returns "+" for Enhanced, "-" for Diminished, "" for Normal.
func (id Identifier) Prefix() string {
	switch id.state {
	case Enhanced:
		return "+"
	case Diminished:
		return "-"
	default:
		return ""
	}
}

// Suffix returns the terminal suffix of the PIN.
// Returns "^" if terminal, "" otherwise.
func (id Identifier) Suffix() string {
	if id.terminal {
		return "^"
	}
	return ""
}

// ============================================================================
// State Transformations
// ============================================================================

// Enhance returns a new Identifier with Enhanced state.
func (id Identifier) Enhance() Identifier {
	id.state = Enhanced
	return id
}

// Diminish returns a new Identifier with Diminished state.
func (id Identifier) Diminish() Identifier {
	id.state = Diminished
	return id
}

// Normalize returns a new Identifier with Normal state.
func (id Identifier) Normalize() Identifier {
	id.state = Normal
	return id
}

// ============================================================================
// Side Transformations
// ============================================================================

// Flip returns a new Identifier with the opposite side.
// First becomes Second, Second becomes First.
func (id Identifier) Flip() Identifier {
	if id.side == First {
		id.side = Second
	} else {
		id.side = First
	}
	return id
}

// ============================================================================
// Terminal Transformations
// ============================================================================

// Terminal returns a new Identifier marked as terminal.
func (id Identifier) Terminal() Identifier {
	id.terminal = true
	return id
}

// NonTerminal returns a new Identifier unmarked as terminal.
func (id Identifier) NonTerminal() Identifier {
	id.terminal = false
	return id
}

// ============================================================================
// Attribute Transformations
// ============================================================================

// WithAbbr returns a new Identifier with the specified abbreviation.
// The abbreviation must be an uppercase letter (A-Z). If a lowercase letter is provided,
// it will be converted to uppercase.
//
// Panics if the abbreviation is not a letter A-Z (case-insensitive).
func (id Identifier) WithAbbr(abbr rune) Identifier {
	// Normalize to uppercase
	if abbr >= 'a' && abbr <= 'z' {
		abbr = abbr - 'a' + 'A'
	}

	if !isValidAbbr(abbr) {
		panic(ErrInvalidAbbr)
	}

	id.abbr = abbr
	return id
}

// WithSide returns a new Identifier with the specified side.
//
// Panics if the side is invalid.
func (id Identifier) WithSide(side Side) Identifier {
	if !isValidSide(side) {
		panic(ErrInvalidSide)
	}

	id.side = side
	return id
}

// WithState returns a new Identifier with the specified state.
//
// Panics if the state is invalid.
func (id Identifier) WithState(state State) Identifier {
	if !isValidState(state) {
		panic(ErrInvalidState)
	}

	id.state = state
	return id
}

// WithTerminal returns a new Identifier with the specified terminal status.
func (id Identifier) WithTerminal(terminal bool) Identifier {
	id.terminal = terminal
	return id
}

// ============================================================================
// State Queries
// ============================================================================

// IsNormal reports whether the Identifier has Normal state.
func (id Identifier) IsNormal() bool {
	return id.state == Normal
}

// IsEnhanced reports whether the Identifier has Enhanced state.
func (id Identifier) IsEnhanced() bool {
	return id.state == Enhanced
}

// IsDiminished reports whether the Identifier has Diminished state.
func (id Identifier) IsDiminished() bool {
	return id.state == Diminished
}

// ============================================================================
// Side Queries
// ============================================================================

// IsFirstPlayer reports whether the Identifier belongs to the first player.
func (id Identifier) IsFirstPlayer() bool {
	return id.side == First
}

// IsSecondPlayer reports whether the Identifier belongs to the second player.
func (id Identifier) IsSecondPlayer() bool {
	return id.side == Second
}

// ============================================================================
// Comparison Queries
// ============================================================================

// SameAbbr reports whether two Identifiers have the same abbreviation.
func (id Identifier) SameAbbr(other Identifier) bool {
	return id.abbr == other.abbr
}

// SameSide reports whether two Identifiers have the same side.
func (id Identifier) SameSide(other Identifier) bool {
	return id.side == other.side
}

// SameState reports whether two Identifiers have the same state.
func (id Identifier) SameState(other Identifier) bool {
	return id.state == other.state
}

// SameTerminal reports whether two Identifiers have the same terminal status.
func (id Identifier) SameTerminal(other Identifier) bool {
	return id.terminal == other.terminal
}
