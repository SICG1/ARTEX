package server

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCreateUploadSinBucle(t *testing.T) {
	dir := t.TempDir()
	if _, _, err := createUpload(dir, strings.Repeat("a", 300)); err == nil {
		t.Fatal("nombre de 300 bytes: debe devolver error y no quedar en bucle")
	}
	if _, _, err := createUpload(dir, "x\x00y.txt"); err == nil {
		t.Fatal("nombre con byte NUL: debe devolver error")
	}
	if err := os.WriteFile(filepath.Join(dir, "a.txt"), nil, 0o600); err != nil {
		t.Fatal(err)
	}
	f1, got1, err := createUpload(dir, "a.txt")
	if err != nil || filepath.Base(got1) != "a-1.txt" {
		t.Fatalf("createUpload(a.txt) = %q, %v; se esperaba a-1.txt", got1, err)
	}
	f1.Close()
	f2, got2, err := createUpload(dir, "a.txt")
	if err != nil || filepath.Base(got2) != "a-2.txt" {
		t.Fatalf("segunda subida con el mismo nombre = %q, %v; se esperaba a-2.txt", got2, err)
	}
	f2.Close()
}
