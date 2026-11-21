package report

import (
	"strings"
	"testing"
)

func TestCodeFenceNoSeCierraAntes(t *testing.T) {
	if got := codeFence("texto sin comillas"); got != "```" {
		t.Fatalf("valla sin comillas: %q", got)
	}
	contenido := "antes ``` despues"
	if got := codeFence(contenido); got != "````" || strings.Contains(contenido, got) {
		t.Fatalf("valla con tres comillas dentro: %q", got)
	}
}
