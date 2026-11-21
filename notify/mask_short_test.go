package notify

import (
	"strings"
	"testing"
)

// A short secret, such as an SMTP password, must not show any of its characters.
func TestMaskedValueHidesShortSecretsCompletely(t *testing.T) {
	short := strings.Repeat("x", 13) + "abcdef" // 19 characters
	if got := MaskedValue(short); got != MaskedPrefix {
		t.Fatalf("19-character secret leaked a hint: %q", got)
	}
	// From 20 characters on, the tail hint is kept as before.
	if got := MaskedValue(short + "g"); !strings.HasSuffix(got, "bcdefg") {
		t.Fatalf("20-character secret lost its tail hint: %q", got)
	}
}
