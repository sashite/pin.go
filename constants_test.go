package pin

import "testing"

// ============================================================================
// Side Tests
// ============================================================================

func TestSideConstants(t *testing.T) {
	// Verify iota values
	if First != 0 {
		t.Errorf("First = %d, want 0", First)
	}
	if Second != 1 {
		t.Errorf("Second = %d, want 1", Second)
	}
}

func TestSideString(t *testing.T) {
	tests := []struct {
		side Side
		want string
	}{
		{First, "First"},
		{Second, "Second"},
		{Side(99), "Unknown"},
	}

	for _, tt := range tests {
		got := tt.side.String()
		if got != tt.want {
			t.Errorf("Side(%d).String() = %q, want %q", tt.side, got, tt.want)
		}
	}
}

// ============================================================================
// State Tests
// ============================================================================

func TestStateConstants(t *testing.T) {
	// Verify iota values
	if Normal != 0 {
		t.Errorf("Normal = %d, want 0", Normal)
	}
	if Enhanced != 1 {
		t.Errorf("Enhanced = %d, want 1", Enhanced)
	}
	if Diminished != 2 {
		t.Errorf("Diminished = %d, want 2", Diminished)
	}
}

func TestStateString(t *testing.T) {
	tests := []struct {
		state State
		want  string
	}{
		{Normal, "Normal"},
		{Enhanced, "Enhanced"},
		{Diminished, "Diminished"},
		{State(99), "Unknown"},
	}

	for _, tt := range tests {
		got := tt.state.String()
		if got != tt.want {
			t.Errorf("State(%d).String() = %q, want %q", tt.state, got, tt.want)
		}
	}
}

// ============================================================================
// Limits Tests
// ============================================================================

func TestMaxStringLength(t *testing.T) {
	if MaxStringLength != 3 {
		t.Errorf("MaxStringLength = %d, want 3", MaxStringLength)
	}
}

// ============================================================================
// Validation Helper Tests
// ============================================================================

func TestIsValidSide(t *testing.T) {
	tests := []struct {
		side Side
		want bool
	}{
		{First, true},
		{Second, true},
		{Side(2), false},
		{Side(99), false},
	}

	for _, tt := range tests {
		got := isValidSide(tt.side)
		if got != tt.want {
			t.Errorf("isValidSide(%d) = %v, want %v", tt.side, got, tt.want)
		}
	}
}

func TestIsValidState(t *testing.T) {
	tests := []struct {
		state State
		want  bool
	}{
		{Normal, true},
		{Enhanced, true},
		{Diminished, true},
		{State(3), false},
		{State(99), false},
	}

	for _, tt := range tests {
		got := isValidState(tt.state)
		if got != tt.want {
			t.Errorf("isValidState(%d) = %v, want %v", tt.state, got, tt.want)
		}
	}
}

func TestIsValidType(t *testing.T) {
	// Valid types: A-Z
	for r := 'A'; r <= 'Z'; r++ {
		if !isValidType(r) {
			t.Errorf("isValidType(%q) = false, want true", r)
		}
	}

	// Invalid types
	invalidTypes := []rune{
		'a', 'z', // lowercase
		'0', '9', // digits
		'+', '-', '^', // modifiers
		' ', '\t', '\n', // whitespace
		'@', '[', // adjacent to A-Z
		0, 127, // control characters
	}

	for _, r := range invalidTypes {
		if isValidType(r) {
			t.Errorf("isValidType(%q) = true, want false", r)
		}
	}
}

// ============================================================================
// Formatting Constants Tests
// ============================================================================

func TestFormattingConstants(t *testing.T) {
	if enhancedPrefix != '+' {
		t.Errorf("enhancedPrefix = %q, want '+'", enhancedPrefix)
	}
	if diminishedPrefix != '-' {
		t.Errorf("diminishedPrefix = %q, want '-'", diminishedPrefix)
	}
	if terminalSuffix != '^' {
		t.Errorf("terminalSuffix = %q, want '^'", terminalSuffix)
	}
}
