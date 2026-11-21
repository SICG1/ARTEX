package server

import (
	"os"
	"path/filepath"
	"testing"
)

// wsResolveConfined must not follow a symlink out of the work dir (v0.7.0
// confinement layer): every path through an escaping link is refused, a link
// that stays inside is accepted, and ".." never climbs above the root.
func TestWsResolveConfinedNoSigueEnlacesFuera(t *testing.T) {
	root := t.TempDir()
	outside := t.TempDir()
	if err := os.WriteFile(filepath.Join(outside, "clave.txt"), []byte("x"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(root, "dentro"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(outside, filepath.Join(root, "fuera")); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(filepath.Join(root, "dentro"), filepath.Join(root, "alias")); err != nil {
		t.Fatal(err)
	}
	s := &Server{m: &Manager{dir: root}}
	for _, p := range []string{"fuera", "fuera/clave.txt", "fuera/nuevo.txt", "../x.txt", "dentro/../fuera"} {
		if _, ok := s.wsResolveConfined(p); ok {
			t.Fatalf("wsResolveConfined(%q) acepto una ruta que sale de la raiz", p)
		}
	}
	if _, ok := s.wsResolveConfined("alias/nuevo.txt"); !ok {
		t.Fatal("un enlace que queda dentro de la raiz debe aceptarse")
	}
	if cp, ok := s.wsResolveConfined("dentro/nuevo.txt"); !ok || cp.realPath != filepath.Join(cp.realRoot, "dentro", "nuevo.txt") {
		t.Fatalf("una ruta normal debe resolverse dentro: %q %v", cp.realPath, ok)
	}
}
