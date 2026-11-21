package traffic

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCopyFileToCopiaConModo0600(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "origen.bin")
	if err := os.WriteFile(src, []byte("cuerpo"), 0o644); err != nil {
		t.Fatal(err)
	}
	dst := filepath.Join(dir, "destino.bin")
	if err := copyFileTo(src, dst); err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(dst)
	if err != nil || string(b) != "cuerpo" {
		t.Fatalf("contenido copiado %q, error %v", b, err)
	}
	fi, err := os.Stat(dst)
	if err != nil {
		t.Fatal(err)
	}
	if fi.Mode().Perm() != 0o600 {
		t.Fatalf("modo de la copia %v; se esperaba 0600", fi.Mode().Perm())
	}
}

// A missing legacy head is logged and left empty, so archiving goes on (as in the base).
func TestLegacyHeadsSigueSiFaltaUnArchivo(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "ex1"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "ex1", "request.http"), []byte("GET /"), 0o600); err != nil {
		t.Fatal(err)
	}
	req, resp, err := legacyHeads(root, "ex1")
	if err != nil || req != "GET /" || resp != "" {
		t.Fatalf("falta response.http: cabeceras %q %q, error %v; se esperaba la solicitud y una respuesta vacia sin error", req, resp, err)
	}
	if err := os.WriteFile(filepath.Join(root, "ex1", "response.http"), []byte("200 OK"), 0o600); err != nil {
		t.Fatal(err)
	}
	req, resp, err = legacyHeads(root, "ex1")
	if err != nil || req != "GET /" || resp != "200 OK" {
		t.Fatalf("cabeceras %q %q, error %v", req, resp, err)
	}
}
