package traffic

import (
	"strings"
	"testing"
)

func TestLogURLQuitaQueryYFragmento(t *testing.T) {
	got := logURL("https://h.example.com/p/a?token=abc123#frag")
	if got != "https://h.example.com/p/a" || strings.Contains(got, "abc123") {
		t.Fatalf("logURL = %q", got)
	}
}
