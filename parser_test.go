package pin

import (
	"errors"
	"testing"
)

// ============================================================================
// Valid Inputs - Simple Letters
// ============================================================================

func TestParseUppercaseLetter(t *testing.T) {
	id, err := Parse("K")
	if err != nil {
		t.Fatalf("Parse(\"K\") error = %v", err)
	}

	if id.Type() != 'K' {
		t.Errorf("Type() = %q, want 'K'", id.Type())
	}
	if id.Side() != First {
		t.Errorf("Side() = %v, want First", id.Side())
	}
	if id.State() != Normal {
		t.Errorf("State() = %v, want Normal", id.State())
	}
	if id.Terminal() != false {
		t.Errorf("Terminal() = %v, want false", id.Terminal())
	}
}

func TestParseLowercaseLetter(t *testing.T) {
	id, err := Parse("k")
	if err != nil {
		t.Fatalf("Parse(\"k\") error = %v", err)
	}

	if id.Type() != 'K' {
		t.Errorf("Type() = %q, want 'K' (normalized)", id.Type())
	}
	if id.Side() != Second {
		t.Errorf("Side() = %v, want Second", id.Side())
	}
}

func TestParseAllUppercaseLetters(t *testing.T) {
	for r := 'A'; r <= 'Z'; r++ {
		input := string(r)
		id, err := Parse(input)
		if err != nil {
			t.Errorf("Parse(%q) error = %v", input, err)
			continue
		}
		if id.Type() != r {
			t.Errorf("Parse(%q).Type() = %q, want %q", input, id.Type(), r)
		}
		if id.Side() != First {
			t.Errorf("Parse(%q).Side() = %v, want First", input, id.Side())
		}
	}
}

func TestParseAllLowercaseLetters(t *testing.T) {
	for r := 'a'; r <= 'z'; r++ {
		input := string(r)
		id, err := Parse(input)
		if err != nil {
			t.Errorf("Parse(%q) error = %v", input, err)
			continue
		}
		expectedType := r - 'a' + 'A'
		if id.Type() != expectedType {
			t.Errorf("Parse(%q).Type() = %q, want %q", input, id.Type(), expectedType)
		}
		if id.Side() != Second {
			t.Errorf("Parse(%q).Side() = %v, want Second", input, id.Side())
		}
	}
}

// ============================================================================
// Valid Inputs - State Modifiers
// ============================================================================

func TestParseEnhancedUppercase(t *testing.T) {
	id, err := Parse("+R")
	if err != nil {
		t.Fatalf("Parse(\"+R\") error = %v", err)
	}

	if id.Type() != 'R' {
		t.Errorf("Type() = %q, want 'R'", id.Type())
	}
	if id.Side() != First {
		t.Errorf("Side() = %v, want First", id.Side())
	}
	if id.State() != Enhanced {
		t.Errorf("State() = %v, want Enhanced", id.State())
	}
	if id.Terminal() != false {
		t.Errorf("Terminal() = %v, want false", id.Terminal())
	}
}

func TestParseEnhancedLowercase(t *testing.T) {
	id, err := Parse("+r")
	if err != nil {
		t.Fatalf("Parse(\"+r\") error = %v", err)
	}

	if id.State() != Enhanced {
		t.Errorf("State() = %v, want Enhanced", id.State())
	}
	if id.Side() != Second {
		t.Errorf("Side() = %v, want Second", id.Side())
	}
}

func TestParseDiminishedUppercase(t *testing.T) {
	id, err := Parse("-P")
	if err != nil {
		t.Fatalf("Parse(\"-P\") error = %v", err)
	}

	if id.Type() != 'P' {
		t.Errorf("Type() = %q, want 'P'", id.Type())
	}
	if id.State() != Diminished {
		t.Errorf("State() = %v, want Diminished", id.State())
	}
}

func TestParseDiminishedLowercase(t *testing.T) {
	id, err := Parse("-p")
	if err != nil {
		t.Fatalf("Parse(\"-p\") error = %v", err)
	}

	if id.State() != Diminished {
		t.Errorf("State() = %v, want Diminished", id.State())
	}
	if id.Side() != Second {
		t.Errorf("Side() = %v, want Second", id.Side())
	}
}

// ============================================================================
// Valid Inputs - Terminal Marker
// ============================================================================

func TestParseTerminalUppercase(t *testing.T) {
	id, err := Parse("K^")
	if err != nil {
		t.Fatalf("Parse(\"K^\") error = %v", err)
	}

	if id.Type() != 'K' {
		t.Errorf("Type() = %q, want 'K'", id.Type())
	}
	if id.Side() != First {
		t.Errorf("Side() = %v, want First", id.Side())
	}
	if id.State() != Normal {
		t.Errorf("State() = %v, want Normal", id.State())
	}
	if id.Terminal() != true {
		t.Errorf("Terminal() = %v, want true", id.Terminal())
	}
}

func TestParseTerminalLowercase(t *testing.T) {
	id, err := Parse("k^")
	if err != nil {
		t.Fatalf("Parse(\"k^\") error = %v", err)
	}

	if id.Side() != Second {
		t.Errorf("Side() = %v, want Second", id.Side())
	}
	if id.Terminal() != true {
		t.Errorf("Terminal() = %v, want true", id.Terminal())
	}
}

// ============================================================================
// Valid Inputs - Combined
// ============================================================================

func TestParseEnhancedTerminal(t *testing.T) {
	id, err := Parse("+K^")
	if err != nil {
		t.Fatalf("Parse(\"+K^\") error = %v", err)
	}

	if id.Type() != 'K' {
		t.Errorf("Type() = %q, want 'K'", id.Type())
	}
	if id.Side() != First {
		t.Errorf("Side() = %v, want First", id.Side())
	}
	if id.State() != Enhanced {
		t.Errorf("State() = %v, want Enhanced", id.State())
	}
	if id.Terminal() != true {
		t.Errorf("Terminal() = %v, want true", id.Terminal())
	}
}

func TestParseDiminishedTerminal(t *testing.T) {
	id, err := Parse("-k^")
	if err != nil {
		t.Fatalf("Parse(\"-k^\") error = %v", err)
	}

	if id.Type() != 'K' {
		t.Errorf("Type() = %q, want 'K'", id.Type())
	}
	if id.Side() != Second {
		t.Errorf("Side() = %v, want Second", id.Side())
	}
	if id.State() != Diminished {
		t.Errorf("State() = %v, want Diminished", id.State())
	}
	if id.Terminal() != true {
		t.Errorf("Terminal() = %v, want true", id.Terminal())
	}
}

// ============================================================================
// MustParse Tests
// ============================================================================

func TestMustParseValid(t *testing.T) {
	id := MustParse("+K^")

	if id.Type() != 'K' {
		t.Errorf("Type() = %q, want 'K'", id.Type())
	}
	if id.State() != Enhanced {
		t.Errorf("State() = %v, want Enhanced", id.State())
	}
	if id.Terminal() != true {
		t.Errorf("Terminal() = %v, want true", id.Terminal())
	}
}

func TestMustParsePanicsOnInvalid(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for invalid input")
		}
	}()

	MustParse("invalid")
}

func TestMustParsePanicsOnEmpty(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for empty input")
		}
	}()

	MustParse("")
}

// ============================================================================
// Validate Tests
// ============================================================================

func TestValidateValid(t *testing.T) {
	validInputs := []string{"K", "k", "+R", "-p", "K^", "+K^", "-k^"}

	for _, input := range validInputs {
		if err := Validate(input); err != nil {
			t.Errorf("Validate(%q) = %v, want nil", input, err)
		}
	}
}

func TestValidateInvalid(t *testing.T) {
	invalidInputs := []string{"", "invalid", "+", "KK", "*K"}

	for _, input := range invalidInputs {
		if err := Validate(input); err == nil {
			t.Errorf("Validate(%q) = nil, want error", input)
		}
	}
}

// ============================================================================
// IsValid Tests
// ============================================================================

func TestIsValidTrue(t *testing.T) {
	validInputs := []string{"K", "k", "+R", "-p", "K^", "+K^", "-k^"}

	for _, input := range validInputs {
		if !IsValid(input) {
			t.Errorf("IsValid(%q) = false, want true", input)
		}
	}
}

func TestIsValidFalse(t *testing.T) {
	invalidInputs := []string{"", "invalid", "+", "KK", "*K", "1", "^"}

	for _, input := range invalidInputs {
		if IsValid(input) {
			t.Errorf("IsValid(%q) = true, want false", input)
		}
	}
}

// ============================================================================
// Error Cases - Empty Input
// ============================================================================

func TestParseEmptyInput(t *testing.T) {
	_, err := Parse("")
	if !errors.Is(err, ErrEmptyInput) {
		t.Errorf("Parse(\"\") error = %v, want ErrEmptyInput", err)
	}
}

// ============================================================================
// Error Cases - Input Too Long
// ============================================================================

func TestParseInputTooLong(t *testing.T) {
	inputs := []string{"+K^X", "invalid", "KKKK", "abcd"}

	for _, input := range inputs {
		_, err := Parse(input)
		if !errors.Is(err, ErrInputTooLong) {
			t.Errorf("Parse(%q) error = %v, want ErrInputTooLong", input, err)
		}
	}
}

// ============================================================================
// Error Cases - Must Contain One Letter
// ============================================================================

func TestParseMustContainOneLetter(t *testing.T) {
	inputs := []string{"+", "-", "1", "++", "+-", "+1"}

	for _, input := range inputs {
		_, err := Parse(input)
		if !errors.Is(err, ErrMustContainOneLetter) {
			t.Errorf("Parse(%q) error = %v, want ErrMustContainOneLetter", input, err)
		}
	}
}

// ============================================================================
// Error Cases - Invalid State Modifier
// ============================================================================

func TestParseInvalidStateModifier(t *testing.T) {
	inputs := []string{"*K", "!R", "1K", "^K", "^K^", "*K^", "1K^"}

	for _, input := range inputs {
		_, err := Parse(input)
		if !errors.Is(err, ErrInvalidStateModifier) {
			t.Errorf("Parse(%q) error = %v, want ErrInvalidStateModifier", input, err)
		}
	}
}

// ============================================================================
// Error Cases - Invalid Terminal Marker
// ============================================================================

func TestParseInvalidTerminalMarker(t *testing.T) {
	inputs := []string{"KQ", "K!", "K1", "KK", "K+", "K-", "+K!", "-R1", "K1^", "KQR"}

	for _, input := range inputs {
		_, err := Parse(input)
		if !errors.Is(err, ErrInvalidTerminalMarker) {
			t.Errorf("Parse(%q) error = %v, want ErrInvalidTerminalMarker", input, err)
		}
	}
}

// ============================================================================
// Security - Null Byte Injection
// ============================================================================

func TestParseRejectsNullBytes(t *testing.T) {
	inputs := []string{"K\x00", "\x00K", "+\x00K"}

	for _, input := range inputs {
		if IsValid(input) {
			t.Errorf("IsValid(%q) = true, want false (null byte)", input)
		}
	}
}

// ============================================================================
// Security - Control Characters
// ============================================================================

func TestParseRejectsControlCharacters(t *testing.T) {
	inputs := []string{
		"K\n", "\nK", // newline
		"K\r", "\rK", // carriage return
		"K\t", "\tK", // tab
		"K\x01", // SOH
		"K\x1b", // ESC
		"K\x7f", // DEL
	}

	for _, input := range inputs {
		if IsValid(input) {
			t.Errorf("IsValid(%q) = true, want false (control char)", input)
		}
	}
}

// ============================================================================
// Security - Unicode Lookalikes
// ============================================================================

func TestParseRejectsUnicodeLookalikes(t *testing.T) {
	inputs := []string{
		"\xD0\x9A",     // Cyrillic 'К' (U+041A) looks like Latin 'K'
		"\xD0\xB0",     // Cyrillic 'а' (U+0430) looks like Latin 'a'
		"\xCE\x91",     // Greek 'Α' (U+0391) looks like Latin 'A'
		"\xEF\xBC\xAB", // Full-width 'K' (U+FF2B)
		"\xEF\xBD\x8B", // Full-width 'k' (U+FF4B)
	}

	for _, input := range inputs {
		if IsValid(input) {
			t.Errorf("IsValid(%q) = true, want false (unicode lookalike)", input)
		}
	}
}

// ============================================================================
// Security - Combining Characters
// ============================================================================

func TestParseRejectsCombiningCharacters(t *testing.T) {
	inputs := []string{
		"K\xCC\x81", // 'K' + combining acute accent (U+0301)
		"K\xCC\x88", // 'K' + combining diaeresis (U+0308)
	}

	for _, input := range inputs {
		if IsValid(input) {
			t.Errorf("IsValid(%q) = true, want false (combining char)", input)
		}
	}
}

// ============================================================================
// Security - Zero-Width Characters
// ============================================================================

func TestParseRejectsZeroWidthCharacters(t *testing.T) {
	inputs := []string{
		"K\xE2\x80\x8B", // Zero-width space (U+200B)
		"K\xE2\x80\x8C", // Zero-width non-joiner (U+200C)
		"\xEF\xBB\xBFK", // BOM (U+FEFF)
	}

	for _, input := range inputs {
		if IsValid(input) {
			t.Errorf("IsValid(%q) = true, want false (zero-width char)", input)
		}
	}
}

// ============================================================================
// Round-Trip Tests
// ============================================================================

func TestRoundTripSimpleLetters(t *testing.T) {
	inputs := []string{"K", "k", "Q", "q", "R", "r"}

	for _, input := range inputs {
		id, err := Parse(input)
		if err != nil {
			t.Errorf("Parse(%q) error = %v", input, err)
			continue
		}
		if id.String() != input {
			t.Errorf("Parse(%q).String() = %q, want %q", input, id.String(), input)
		}
	}
}

func TestRoundTripWithStateModifiers(t *testing.T) {
	inputs := []string{"+K", "+k", "-P", "-p", "+R", "-r"}

	for _, input := range inputs {
		id, err := Parse(input)
		if err != nil {
			t.Errorf("Parse(%q) error = %v", input, err)
			continue
		}
		if id.String() != input {
			t.Errorf("Parse(%q).String() = %q, want %q", input, id.String(), input)
		}
	}
}

func TestRoundTripWithTerminalMarker(t *testing.T) {
	inputs := []string{"K^", "k^", "Q^", "q^"}

	for _, input := range inputs {
		id, err := Parse(input)
		if err != nil {
			t.Errorf("Parse(%q) error = %v", input, err)
			continue
		}
		if id.String() != input {
			t.Errorf("Parse(%q).String() = %q, want %q", input, id.String(), input)
		}
	}
}

func TestRoundTripCombined(t *testing.T) {
	inputs := []string{"+K^", "-k^", "+Q^", "-q^"}

	for _, input := range inputs {
		id, err := Parse(input)
		if err != nil {
			t.Errorf("Parse(%q) error = %v", input, err)
			continue
		}
		if id.String() != input {
			t.Errorf("Parse(%q).String() = %q, want %q", input, id.String(), input)
		}
	}
}

// ============================================================================
// Edge Cases
// ============================================================================

func TestParseExactlyMaxLength(t *testing.T) {
	// 3 characters is the max valid length
	inputs := []string{"+K^", "-k^", "+R^", "-p^"}

	for _, input := range inputs {
		if len(input) != MaxStringLength {
			t.Errorf("test input %q should be exactly %d characters", input, MaxStringLength)
		}
		if _, err := Parse(input); err != nil {
			t.Errorf("Parse(%q) error = %v, want nil", input, err)
		}
	}
}

func TestParseOneOverMaxLength(t *testing.T) {
	// 4 characters should fail
	input := "+K^X"
	if len(input) != MaxStringLength+1 {
		t.Errorf("test input %q should be exactly %d characters", input, MaxStringLength+1)
	}
	_, err := Parse(input)
	if !errors.Is(err, ErrInputTooLong) {
		t.Errorf("Parse(%q) error = %v, want ErrInputTooLong", input, err)
	}
}
