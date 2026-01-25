package pin

import (
	"errors"
	"testing"
)

// ============================================================================
// Parsing Errors Tests
// ============================================================================

func TestParsingErrorsExist(t *testing.T) {
	// Verify all parsing errors are non-nil
	parsingErrors := []error{
		ErrEmptyInput,
		ErrInputTooLong,
		ErrMustContainOneLetter,
		ErrInvalidStateModifier,
		ErrInvalidTerminalMarker,
	}

	for _, err := range parsingErrors {
		if err == nil {
			t.Error("parsing error is nil")
		}
	}
}

func TestParsingErrorMessages(t *testing.T) {
	tests := []struct {
		err  error
		want string
	}{
		{ErrEmptyInput, "pin: empty input"},
		{ErrInputTooLong, "pin: input exceeds 3 characters"},
		{ErrMustContainOneLetter, "pin: must contain exactly one letter"},
		{ErrInvalidStateModifier, "pin: invalid state modifier"},
		{ErrInvalidTerminalMarker, "pin: invalid terminal marker"},
	}

	for _, tt := range tests {
		got := tt.err.Error()
		if got != tt.want {
			t.Errorf("error message = %q, want %q", got, tt.want)
		}
	}
}

func TestParsingErrorsAreDistinct(t *testing.T) {
	parsingErrors := []error{
		ErrEmptyInput,
		ErrInputTooLong,
		ErrMustContainOneLetter,
		ErrInvalidStateModifier,
		ErrInvalidTerminalMarker,
	}

	for i, err1 := range parsingErrors {
		for j, err2 := range parsingErrors {
			if i != j && errors.Is(err1, err2) {
				t.Errorf("errors should be distinct: %v and %v", err1, err2)
			}
		}
	}
}

// ============================================================================
// Validation Errors Tests
// ============================================================================

func TestValidationErrorsExist(t *testing.T) {
	// Verify all validation errors are non-nil
	validationErrors := []error{
		ErrInvalidAbbr,
		ErrInvalidSide,
		ErrInvalidState,
	}

	for _, err := range validationErrors {
		if err == nil {
			t.Error("validation error is nil")
		}
	}
}

func TestValidationErrorMessages(t *testing.T) {
	tests := []struct {
		err  error
		want string
	}{
		{ErrInvalidAbbr, "pin: invalid abbr (must be A-Z)"},
		{ErrInvalidSide, "pin: invalid side"},
		{ErrInvalidState, "pin: invalid state"},
	}

	for _, tt := range tests {
		got := tt.err.Error()
		if got != tt.want {
			t.Errorf("error message = %q, want %q", got, tt.want)
		}
	}
}

func TestValidationErrorsAreDistinct(t *testing.T) {
	validationErrors := []error{
		ErrInvalidAbbr,
		ErrInvalidSide,
		ErrInvalidState,
	}

	for i, err1 := range validationErrors {
		for j, err2 := range validationErrors {
			if i != j && errors.Is(err1, err2) {
				t.Errorf("errors should be distinct: %v and %v", err1, err2)
			}
		}
	}
}

// ============================================================================
// errors.Is Tests
// ============================================================================

func TestErrorsIs(t *testing.T) {
	// Verify errors.Is works correctly with our sentinel errors
	allErrors := []error{
		ErrEmptyInput,
		ErrInputTooLong,
		ErrMustContainOneLetter,
		ErrInvalidStateModifier,
		ErrInvalidTerminalMarker,
		ErrInvalidAbbr,
		ErrInvalidSide,
		ErrInvalidState,
	}

	for _, err := range allErrors {
		if !errors.Is(err, err) {
			t.Errorf("errors.Is(%v, %v) = false, want true", err, err)
		}
	}
}

// ============================================================================
// Error Prefix Tests
// ============================================================================

func TestErrorsHavePinPrefix(t *testing.T) {
	allErrors := []error{
		ErrEmptyInput,
		ErrInputTooLong,
		ErrMustContainOneLetter,
		ErrInvalidStateModifier,
		ErrInvalidTerminalMarker,
		ErrInvalidAbbr,
		ErrInvalidSide,
		ErrInvalidState,
	}

	for _, err := range allErrors {
		msg := err.Error()
		if len(msg) < 4 || msg[:4] != "pin:" {
			t.Errorf("error %q should start with 'pin:'", msg)
		}
	}
}
