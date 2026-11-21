package config

import (
	"strings"
	"testing"
)

func TestRedactNoDevuelveDSNSinAnalizar(t *testing.T) {
	if got := Redact("postgres://user:pa%zz@127.0.0.1/db"); strings.Contains(got, "pa%zz") {
		t.Fatalf("Redact dejo la contrasena en claro: %q", got)
	}
}
