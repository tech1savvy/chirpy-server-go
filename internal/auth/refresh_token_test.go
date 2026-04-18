package auth

import (
	"testing"
	"unicode"
)

func TestMakeRefreshToken(t *testing.T) {
	got := MakeRefreshToken()
	if len(got) != 64 {
		t.Errorf("expected 64 chars, got %d", len(got))
	}
	for _, r := range got {
		if !unicode.Is(unicode.Hex_Digit, r) {
			t.Errorf("invalid hex character: %c", r)
		}
	}
}
