package pin

import "testing"

// ============================================================================
// Constructor Tests
// ============================================================================

func TestNewIdentifier(t *testing.T) {
	id := NewIdentifier('K', First)

	if id.Abbr() != 'K' {
		t.Errorf("Abbr() = %q, want 'K'", id.Abbr())
	}
	if id.Side() != First {
		t.Errorf("Side() = %v, want First", id.Side())
	}
	if id.State() != Normal {
		t.Errorf("State() = %v, want Normal", id.State())
	}
	if id.IsTerminal() != false {
		t.Errorf("IsTerminal() = %v, want false", id.IsTerminal())
	}
}

func TestNewIdentifierWithOptions(t *testing.T) {
	id := NewIdentifierWithOptions('R', Second, Enhanced, true)

	if id.Abbr() != 'R' {
		t.Errorf("Abbr() = %q, want 'R'", id.Abbr())
	}
	if id.Side() != Second {
		t.Errorf("Side() = %v, want Second", id.Side())
	}
	if id.State() != Enhanced {
		t.Errorf("State() = %v, want Enhanced", id.State())
	}
	if id.IsTerminal() != true {
		t.Errorf("IsTerminal() = %v, want true", id.IsTerminal())
	}
}

func TestNewIdentifierNormalizesLowercase(t *testing.T) {
	id := NewIdentifier('k', First)

	if id.Abbr() != 'K' {
		t.Errorf("Abbr() = %q, want 'K' (normalized from 'k')", id.Abbr())
	}
}

func TestNewIdentifierAllAbbrs(t *testing.T) {
	for r := 'A'; r <= 'Z'; r++ {
		id := NewIdentifier(r, First)
		if id.Abbr() != r {
			t.Errorf("NewIdentifier(%q, First).Abbr() = %q, want %q", r, id.Abbr(), r)
		}
	}
}

func TestNewIdentifierPanicsOnInvalidAbbr(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for invalid abbr")
		}
	}()

	NewIdentifier('1', First)
}

func TestNewIdentifierPanicsOnInvalidSide(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for invalid side")
		}
	}()

	NewIdentifier('K', Side(99))
}

func TestNewIdentifierWithOptionsPanicsOnInvalidState(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for invalid state")
		}
	}()

	NewIdentifierWithOptions('K', First, State(99), false)
}

// ============================================================================
// String Conversion Tests
// ============================================================================

func TestIdentifierString(t *testing.T) {
	tests := []struct {
		id   Identifier
		want string
	}{
		{NewIdentifier('K', First), "K"},
		{NewIdentifier('K', Second), "k"},
		{NewIdentifierWithOptions('R', First, Enhanced, false), "+R"},
		{NewIdentifierWithOptions('R', Second, Enhanced, false), "+r"},
		{NewIdentifierWithOptions('P', First, Diminished, false), "-P"},
		{NewIdentifierWithOptions('P', Second, Diminished, false), "-p"},
		{NewIdentifierWithOptions('K', First, Normal, true), "K^"},
		{NewIdentifierWithOptions('K', Second, Normal, true), "k^"},
		{NewIdentifierWithOptions('K', First, Enhanced, true), "+K^"},
		{NewIdentifierWithOptions('K', Second, Diminished, true), "-k^"},
	}

	for _, tt := range tests {
		got := tt.id.String()
		if got != tt.want {
			t.Errorf("String() = %q, want %q", got, tt.want)
		}
	}
}

func TestIdentifierLetter(t *testing.T) {
	tests := []struct {
		id   Identifier
		want string
	}{
		{NewIdentifier('K', First), "K"},
		{NewIdentifier('K', Second), "k"},
		{NewIdentifier('R', First), "R"},
		{NewIdentifier('R', Second), "r"},
	}

	for _, tt := range tests {
		got := tt.id.Letter()
		if got != tt.want {
			t.Errorf("Letter() = %q, want %q", got, tt.want)
		}
	}
}

func TestIdentifierPrefix(t *testing.T) {
	tests := []struct {
		state State
		want  string
	}{
		{Normal, ""},
		{Enhanced, "+"},
		{Diminished, "-"},
	}

	for _, tt := range tests {
		id := NewIdentifierWithOptions('K', First, tt.state, false)
		got := id.Prefix()
		if got != tt.want {
			t.Errorf("Prefix() for state %v = %q, want %q", tt.state, got, tt.want)
		}
	}
}

func TestIdentifierSuffix(t *testing.T) {
	tests := []struct {
		terminal bool
		want     string
	}{
		{false, ""},
		{true, "^"},
	}

	for _, tt := range tests {
		id := NewIdentifierWithOptions('K', First, Normal, tt.terminal)
		got := id.Suffix()
		if got != tt.want {
			t.Errorf("Suffix() for terminal=%v = %q, want %q", tt.terminal, got, tt.want)
		}
	}
}

// ============================================================================
// AppendTo Tests
// ============================================================================

func TestIdentifierAppendTo(t *testing.T) {
	buf := make([]byte, 0, 16)

	id := NewIdentifierWithOptions('K', First, Enhanced, true)
	buf = id.AppendTo(buf)

	if string(buf) != "+K^" {
		t.Errorf("AppendTo result = %q, want \"+K^\"", string(buf))
	}
}

func TestIdentifierAppendToMultiple(t *testing.T) {
	buf := make([]byte, 0, 32)

	ids := []Identifier{
		NewIdentifier('K', First),
		NewIdentifier('Q', Second),
		NewIdentifierWithOptions('R', First, Enhanced, true),
	}

	for i, id := range ids {
		if i > 0 {
			buf = append(buf, ' ')
		}
		buf = id.AppendTo(buf)
	}

	want := "K q +R^"
	if string(buf) != want {
		t.Errorf("AppendTo result = %q, want %q", string(buf), want)
	}
}

// ============================================================================
// State Transformation Tests
// ============================================================================

func TestIdentifierEnhance(t *testing.T) {
	id := NewIdentifier('K', First)
	enhanced := id.Enhance()

	// Original unchanged
	if id.State() != Normal {
		t.Errorf("original State() = %v, want Normal", id.State())
	}

	// New is enhanced
	if enhanced.State() != Enhanced {
		t.Errorf("enhanced State() = %v, want Enhanced", enhanced.State())
	}
	if enhanced.String() != "+K" {
		t.Errorf("enhanced String() = %q, want \"+K\"", enhanced.String())
	}
}

func TestIdentifierDiminish(t *testing.T) {
	id := NewIdentifier('K', First)
	diminished := id.Diminish()

	// Original unchanged
	if id.State() != Normal {
		t.Errorf("original State() = %v, want Normal", id.State())
	}

	// New is diminished
	if diminished.State() != Diminished {
		t.Errorf("diminished State() = %v, want Diminished", diminished.State())
	}
	if diminished.String() != "-K" {
		t.Errorf("diminished String() = %q, want \"-K\"", diminished.String())
	}
}

func TestIdentifierNormalize(t *testing.T) {
	id := NewIdentifierWithOptions('K', First, Enhanced, false)
	normalized := id.Normalize()

	// Original unchanged
	if id.State() != Enhanced {
		t.Errorf("original State() = %v, want Enhanced", id.State())
	}

	// New is normal
	if normalized.State() != Normal {
		t.Errorf("normalized State() = %v, want Normal", normalized.State())
	}
	if normalized.String() != "K" {
		t.Errorf("normalized String() = %q, want \"K\"", normalized.String())
	}
}

// ============================================================================
// Side Transformation Tests
// ============================================================================

func TestIdentifierFlip(t *testing.T) {
	tests := []struct {
		initial Side
		want    Side
	}{
		{First, Second},
		{Second, First},
	}

	for _, tt := range tests {
		id := NewIdentifier('K', tt.initial)
		flipped := id.Flip()

		// Original unchanged
		if id.Side() != tt.initial {
			t.Errorf("original Side() = %v, want %v", id.Side(), tt.initial)
		}

		// New is flipped
		if flipped.Side() != tt.want {
			t.Errorf("flipped Side() = %v, want %v", flipped.Side(), tt.want)
		}
	}
}

func TestIdentifierFlipString(t *testing.T) {
	id := NewIdentifier('K', First)
	if id.Flip().String() != "k" {
		t.Errorf("Flip().String() = %q, want \"k\"", id.Flip().String())
	}

	id = NewIdentifier('K', Second)
	if id.Flip().String() != "K" {
		t.Errorf("Flip().String() = %q, want \"K\"", id.Flip().String())
	}
}

// ============================================================================
// Terminal Transformation Tests
// ============================================================================

func TestIdentifierTerminal(t *testing.T) {
	id := NewIdentifier('K', First)
	marked := id.Terminal()

	// Original unchanged
	if id.IsTerminal() != false {
		t.Errorf("original IsTerminal() = %v, want false", id.IsTerminal())
	}

	// New is terminal
	if marked.IsTerminal() != true {
		t.Errorf("marked IsTerminal() = %v, want true", marked.IsTerminal())
	}
	if marked.String() != "K^" {
		t.Errorf("marked String() = %q, want \"K^\"", marked.String())
	}
}

func TestIdentifierNonTerminal(t *testing.T) {
	id := NewIdentifierWithOptions('K', First, Normal, true)
	unmarked := id.NonTerminal()

	// Original unchanged
	if id.IsTerminal() != true {
		t.Errorf("original IsTerminal() = %v, want true", id.IsTerminal())
	}

	// New is not terminal
	if unmarked.IsTerminal() != false {
		t.Errorf("unmarked IsTerminal() = %v, want false", unmarked.IsTerminal())
	}
	if unmarked.String() != "K" {
		t.Errorf("unmarked String() = %q, want \"K\"", unmarked.String())
	}
}

// ============================================================================
// WithX Transformation Tests
// ============================================================================

func TestIdentifierWithAbbr(t *testing.T) {
	id := NewIdentifier('K', First)
	changed := id.WithAbbr('Q')

	// Original unchanged
	if id.Abbr() != 'K' {
		t.Errorf("original Abbr() = %q, want 'K'", id.Abbr())
	}

	// New has new abbr
	if changed.Abbr() != 'Q' {
		t.Errorf("changed Abbr() = %q, want 'Q'", changed.Abbr())
	}
}

func TestIdentifierWithAbbrNormalizesLowercase(t *testing.T) {
	id := NewIdentifier('K', First)
	changed := id.WithAbbr('q')

	if changed.Abbr() != 'Q' {
		t.Errorf("WithAbbr('q').Abbr() = %q, want 'Q'", changed.Abbr())
	}
}

func TestIdentifierWithAbbrPanicsOnInvalid(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for invalid abbr")
		}
	}()

	id := NewIdentifier('K', First)
	id.WithAbbr('1')
}

func TestIdentifierWithSide(t *testing.T) {
	id := NewIdentifier('K', First)
	changed := id.WithSide(Second)

	// Original unchanged
	if id.Side() != First {
		t.Errorf("original Side() = %v, want First", id.Side())
	}

	// New has new side
	if changed.Side() != Second {
		t.Errorf("changed Side() = %v, want Second", changed.Side())
	}
}

func TestIdentifierWithSidePanicsOnInvalid(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for invalid side")
		}
	}()

	id := NewIdentifier('K', First)
	id.WithSide(Side(99))
}

func TestIdentifierWithState(t *testing.T) {
	id := NewIdentifier('K', First)
	changed := id.WithState(Enhanced)

	// Original unchanged
	if id.State() != Normal {
		t.Errorf("original State() = %v, want Normal", id.State())
	}

	// New has new state
	if changed.State() != Enhanced {
		t.Errorf("changed State() = %v, want Enhanced", changed.State())
	}
}

func TestIdentifierWithStatePanicsOnInvalid(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for invalid state")
		}
	}()

	id := NewIdentifier('K', First)
	id.WithState(State(99))
}

func TestIdentifierWithTerminal(t *testing.T) {
	id := NewIdentifier('K', First)
	changed := id.WithTerminal(true)

	// Original unchanged
	if id.IsTerminal() != false {
		t.Errorf("original IsTerminal() = %v, want false", id.IsTerminal())
	}

	// New has new terminal
	if changed.IsTerminal() != true {
		t.Errorf("changed IsTerminal() = %v, want true", changed.IsTerminal())
	}
}

// ============================================================================
// State Query Tests
// ============================================================================

func TestIdentifierIsNormal(t *testing.T) {
	tests := []struct {
		state State
		want  bool
	}{
		{Normal, true},
		{Enhanced, false},
		{Diminished, false},
	}

	for _, tt := range tests {
		id := NewIdentifierWithOptions('K', First, tt.state, false)
		if got := id.IsNormal(); got != tt.want {
			t.Errorf("IsNormal() for state %v = %v, want %v", tt.state, got, tt.want)
		}
	}
}

func TestIdentifierIsEnhanced(t *testing.T) {
	tests := []struct {
		state State
		want  bool
	}{
		{Normal, false},
		{Enhanced, true},
		{Diminished, false},
	}

	for _, tt := range tests {
		id := NewIdentifierWithOptions('K', First, tt.state, false)
		if got := id.IsEnhanced(); got != tt.want {
			t.Errorf("IsEnhanced() for state %v = %v, want %v", tt.state, got, tt.want)
		}
	}
}

func TestIdentifierIsDiminished(t *testing.T) {
	tests := []struct {
		state State
		want  bool
	}{
		{Normal, false},
		{Enhanced, false},
		{Diminished, true},
	}

	for _, tt := range tests {
		id := NewIdentifierWithOptions('K', First, tt.state, false)
		if got := id.IsDiminished(); got != tt.want {
			t.Errorf("IsDiminished() for state %v = %v, want %v", tt.state, got, tt.want)
		}
	}
}

// ============================================================================
// Side Query Tests
// ============================================================================

func TestIdentifierIsFirstPlayer(t *testing.T) {
	tests := []struct {
		side Side
		want bool
	}{
		{First, true},
		{Second, false},
	}

	for _, tt := range tests {
		id := NewIdentifier('K', tt.side)
		if got := id.IsFirstPlayer(); got != tt.want {
			t.Errorf("IsFirstPlayer() for side %v = %v, want %v", tt.side, got, tt.want)
		}
	}
}

func TestIdentifierIsSecondPlayer(t *testing.T) {
	tests := []struct {
		side Side
		want bool
	}{
		{First, false},
		{Second, true},
	}

	for _, tt := range tests {
		id := NewIdentifier('K', tt.side)
		if got := id.IsSecondPlayer(); got != tt.want {
			t.Errorf("IsSecondPlayer() for side %v = %v, want %v", tt.side, got, tt.want)
		}
	}
}

// ============================================================================
// Comparison Query Tests
// ============================================================================

func TestIdentifierSameAbbr(t *testing.T) {
	id1 := NewIdentifier('K', First)
	id2 := NewIdentifier('K', Second)
	id3 := NewIdentifier('Q', First)

	if !id1.SameAbbr(id2) {
		t.Error("SameAbbr(K, K) = false, want true")
	}
	if id1.SameAbbr(id3) {
		t.Error("SameAbbr(K, Q) = true, want false")
	}
}

func TestIdentifierSameSide(t *testing.T) {
	id1 := NewIdentifier('K', First)
	id2 := NewIdentifier('Q', First)
	id3 := NewIdentifier('K', Second)

	if !id1.SameSide(id2) {
		t.Error("SameSide(First, First) = false, want true")
	}
	if id1.SameSide(id3) {
		t.Error("SameSide(First, Second) = true, want false")
	}
}

func TestIdentifierSameState(t *testing.T) {
	id1 := NewIdentifier('K', First)
	id2 := NewIdentifier('Q', Second)
	id3 := NewIdentifierWithOptions('K', First, Enhanced, false)

	if !id1.SameState(id2) {
		t.Error("SameState(Normal, Normal) = false, want true")
	}
	if id1.SameState(id3) {
		t.Error("SameState(Normal, Enhanced) = true, want false")
	}
}

func TestIdentifierSameTerminal(t *testing.T) {
	id1 := NewIdentifier('K', First)
	id2 := NewIdentifier('Q', Second)
	id3 := NewIdentifierWithOptions('K', First, Normal, true)

	if !id1.SameTerminal(id2) {
		t.Error("SameTerminal(false, false) = false, want true")
	}
	if id1.SameTerminal(id3) {
		t.Error("SameTerminal(false, true) = true, want false")
	}
}

// ============================================================================
// Value Semantics Tests
// ============================================================================

func TestIdentifierIsValueType(t *testing.T) {
	// Verify that Identifier behaves as a value type
	id1 := NewIdentifier('K', First)
	id2 := id1 // Copy

	// Modify id2 (through transformation)
	id2 = id2.Enhance()

	// id1 should be unchanged
	if id1.State() != Normal {
		t.Error("modifying copy affected original")
	}
	if id2.State() != Enhanced {
		t.Error("copy was not modified")
	}
}

func TestIdentifierEquality(t *testing.T) {
	id1 := NewIdentifier('K', First)
	id2 := NewIdentifier('K', First)
	id3 := NewIdentifier('K', Second)

	if id1 != id2 {
		t.Error("identical identifiers should be equal")
	}
	if id1 == id3 {
		t.Error("different identifiers should not be equal")
	}
}

// ============================================================================
// Chaining Tests
// ============================================================================

func TestIdentifierChaining(t *testing.T) {
	id := NewIdentifier('K', First).
		Enhance().
		Terminal().
		Flip()

	want := "+k^"
	if id.String() != want {
		t.Errorf("chained result = %q, want %q", id.String(), want)
	}
}
