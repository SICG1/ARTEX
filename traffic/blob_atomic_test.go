package traffic

import (
	"os"
	"path/filepath"
	"testing"
)

func TestWriteFileAtomicSinParciales(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "abc.bin")
	if err := writeFileAtomic(p, []byte("cuerpo"), 0o644); err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(p)
	if err != nil || string(b) != "cuerpo" {
		t.Fatalf("contenido %q, error %v", b, err)
	}
	entries, _ := os.ReadDir(dir)
	if len(entries) != 1 {
		t.Fatalf("quedan archivos temporales: %d entradas", len(entries))
	}
}
