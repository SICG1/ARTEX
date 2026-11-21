package config

import (
	"strings"
	"testing"
)

func TestRedactEnmascaraPasswordClaveValor(t *testing.T) {
	for _, dsn := range []string{
		"host=127.0.0.1 port=5432 user=app password=secreto dbname=x",
		"host=127.0.0.1 password='se creto' dbname=x",
		"postgres://app@127.0.0.1/db?password=secreto",
	} {
		got := Redact(dsn)
		if strings.Contains(got, "secreto") || strings.Contains(got, "se creto") {
			t.Fatalf("Redact(%q) dejo la contrasena en claro: %q", dsn, got)
		}
	}
}
