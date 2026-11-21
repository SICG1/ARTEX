package server

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestComposeAgentMessageSoloAdjuntosEnUploads(t *testing.T) {
	base := t.TempDir()
	atts := []chatAttachment{
		{Name: "ok.txt", Path: "uploads/ok.txt", Size: 3},
		{Name: "x", Path: "uploads/../../secreto.txt", Size: 3},
		{Name: "y", Path: "/etc/passwd", Size: 3},
	}
	out := composeAgentMessage("hola", atts, base)
	if !strings.Contains(out, filepath.Join(base, "uploads", "ok.txt")) {
		t.Fatalf("el adjunto legitimo debe aparecer en el mensaje: %q", out)
	}
	if strings.Contains(out, "secreto") || strings.Contains(out, "/etc/passwd") {
		t.Fatalf("una ruta fuera de uploads no debe llegar al agente: %q", out)
	}
}
