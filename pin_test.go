package pin

import (
	"bytes"
	"testing"
)

// -----------------------------------------------------------------------------
// 1. Table-Driven Tests (Parsing & Logic)
// -----------------------------------------------------------------------------

func TestParseAndRoundTrip(t *testing.T) {
	tests := []struct {
		input    string
		wantAbbr rune
		wantSide Side
		wantStat State
		wantTerm bool
		wantErr  error
	}{
		// --- Valid Cases ---
		// Simple
		{"P", 'P', First, Normal, false, nil},
		{"p", 'P', Second, Normal, false, nil}, // 'p' -> Abbr 'P', Side Second

		// Side / Abbr
		{"K", 'K', First, Normal, false, nil},
		{"n", 'N', Second, Normal, false, nil},

		// State Modifiers
		{"+R", 'R', First, Enhanced, false, nil},
		{"-r", 'R', Second, Diminished, false, nil},

		// Terminal Markers
		{"K^", 'K', First, Normal, true, nil},
		{"k^", 'K', Second, Normal, true, nil},

		// Complex Combinations (All fields)
		{"+S^", 'S', First, Enhanced, true, nil},
		{"-c^", 'C', Second, Diminished, true, nil},

		// --- Error Cases ---
		{"", 0, 0, 0, false, ErrEmpty},
		{"1", 0, 0, 0, false, ErrInvalidAbbr},   // Not a letter
		{"@", 0, 0, 0, false, ErrInvalidAbbr},   // Not a letter
		{"KK", 0, 0, 0, false, ErrInvalid},      // Too long
		{"K ", 0, 0, 0, false, ErrInvalid},      // Trailing space
		{" K", 0, 0, 0, false, ErrInvalidAbbr},  // Leading space
		{"++K", 0, 0, 0, false, ErrInvalidAbbr}, // Double state
		{"K^^", 0, 0, 0, false, ErrInvalid},     // Double terminal
		{"^K", 0, 0, 0, false, ErrInvalidAbbr},  // Wrong order (starts with terminal)
		{"K+", 0, 0, 0, false, ErrInvalid},      // Wrong order (state at end)
		{"+", 0, 0, 0, false, ErrInvalid},       // Missing letter
		{"^", 0, 0, 0, false, ErrInvalidAbbr},   // Missing letter
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := Parse(tt.input)

			// 1. Check Error
			if err != tt.wantErr {
				t.Errorf("Parse(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}

			if err != nil {
				return
			}

			// 2. Check Fields
			if got.Abbr != tt.wantAbbr {
				t.Errorf("Abbr = %c, want %c", got.Abbr, tt.wantAbbr)
			}
			if got.Side != tt.wantSide {
				t.Errorf("Side = %v, want %v", got.Side, tt.wantSide)
			}
			if got.State != tt.wantStat {
				t.Errorf("State = %v, want %v", got.State, tt.wantStat)
			}
			if got.Terminal != tt.wantTerm {
				t.Errorf("Terminal = %v, want %v", got.Terminal, tt.wantTerm)
			}

			// 3. Round-Trip Check (String() matches Input)
			// Note: String() normalizes case. "p" input becomes "p" output.
			if got.String() != tt.input {
				t.Errorf("RoundTrip mismatch: got.String() = %q, want %q", got.String(), tt.input)
			}
		})
	}
}

// -----------------------------------------------------------------------------
// 2. Immutability & API Logic
// -----------------------------------------------------------------------------

func TestImmutability(t *testing.T) {
	origin := MustParse("P") // First, Normal

	// Fluent chain
	promoted := origin.WithSide(Second).WithState(Enhanced).WithTerminal(true)

	// Verify Origin is untouched
	if origin.Side != First {
		t.Errorf("Immutability violation: Origin side changed to %v", origin.Side)
	}
	if origin.State != Normal {
		t.Errorf("Immutability violation: Origin state changed to %v", origin.State)
	}

	// Verify Promoted is correct
	if promoted.Side != Second {
		t.Errorf("Fluent API failed: Side is %v", promoted.Side)
	}
	if promoted.State != Enhanced {
		t.Errorf("Fluent API failed: State is %v", promoted.State)
	}
	if !promoted.Terminal {
		t.Errorf("Fluent API failed: Terminal is false")
	}
}

func TestHelpers(t *testing.T) {
	p := MustParse("P") // First

	// Flip
	p2 := p.Flip()
	if p2.Side != Second {
		t.Errorf("Flip() failed: expected Second, got %v", p2.Side)
	}
	if p.Side != First { // Re-check origin
		t.Error("Flip() mutated original!")
	}

	// Enhance/Diminish
	if p.Enhance().State != Enhanced {
		t.Error("Enhance() failed")
	}
	if p.Diminish().State != Diminished {
		t.Error("Diminish() failed")
	}
}

func TestAppendTo(t *testing.T) {
	// Setup: A PIN that uses all features
	p := Pin{Abbr: 'K', Side: Second, State: Enhanced, Terminal: true}
	// Expected output: +k^ (State + Lowercase + Terminal)
	expected := []byte("+k^")

	// 1. Append to empty
	dst := p.AppendTo(nil)
	if !bytes.Equal(dst, expected) {
		t.Errorf("AppendTo(nil) = %s, want %s", dst, expected)
	}

	// 2. Append to existing buffer
	prefix := []byte("Prefix: ")
	dst = append(make([]byte, 0, 20), prefix...)
	dst = p.AppendTo(dst)

	if string(dst) != "Prefix: +k^" {
		t.Errorf("AppendTo(buf) result %q incorrect", dst)
	}
}

func TestMustParse(t *testing.T) {
	// Should work
	_ = MustParse("P")

	// Should panic
	defer func() {
		if r := recover(); r == nil {
			t.Error("MustParse did not panic on invalid input")
		}
	}()
	_ = MustParse("Invalid")
}

// -----------------------------------------------------------------------------
// 3. Benchmarks (Performance)
// -----------------------------------------------------------------------------

// BenchmarkParse measures the cost of creating a Pin from a string.
// Expectation: 0 Allocs/op (Stack allocation).
func BenchmarkParse(b *testing.B) {
	input := "+K^"
	for i := 0; i < b.N; i++ {
		_, _ = Parse(input)
	}
}

// BenchmarkString measures standard formatting.
// Expectation: 1 Alloc/op (Creating the result string).
func BenchmarkString(b *testing.B) {
	p := Pin{Abbr: 'K', Side: First, State: Enhanced, Terminal: true}
	for i := 0; i < b.N; i++ {
		_ = p.String()
	}
}

// BenchmarkAppendTo measures the zero-allocation formatting path.
// Expectation: 0 Allocs/op (Reusing the buffer).
func BenchmarkAppendTo(b *testing.B) {
	p := Pin{Abbr: 'K', Side: First, State: Enhanced, Terminal: true}
	buf := make([]byte, 0, 8) // Small buffer on stack/heap

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Reset buffer length (keep capacity) -> No allocation
		buf = buf[:0]
		_ = p.AppendTo(buf)
	}
}
