package traffic

import (
	"crypto/sha256"
	"os"
	"path/filepath"
	"testing"
)

func TestBlobIntactDetectaParcial(t *testing.T) {
	p := filepath.Join(t.TempDir(), "x.bin")
	body := []byte("cuerpo completo")
	sum := sha256.Sum256(body)
	if blobIntact(p, sum) {
		t.Fatal("archivo inexistente: no debe darse por bueno")
	}
	if err := os.WriteFile(p, body[:4], 0o600); err != nil {
		t.Fatal(err)
	}
	if blobIntact(p, sum) {
		t.Fatal("archivo parcial: no debe darse por bueno")
	}
	if err := os.WriteFile(p, body, 0o600); err != nil {
		t.Fatal(err)
	}
	if !blobIntact(p, sum) {
		t.Fatal("archivo completo: debe darse por bueno")
	}
}
