package evidence

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"os"
	"testing"
)

func TestWriteBodyReemplazaBlobDanado(t *testing.T) {
	s := &Store{Dir: t.TempDir()}
	body := []byte("cuerpo de evidencia de prueba")
	sum := sha256.Sum256(body)
	hash := hex.EncodeToString(sum[:])
	if _, err := s.writeBody(bytes.NewReader(body), int64(len(body)), hash); err != nil {
		t.Fatal(err)
	}
	path, err := hashPath(s.Dir, hash)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, bytes.Repeat([]byte("x"), len(body)), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := s.writeBody(bytes.NewReader(body), int64(len(body)), hash); err != nil {
		t.Fatalf("un blob dañado debe reemplazarse, no bloquear la evidencia: %v", err)
	}
	got, err := os.ReadFile(path)
	if err != nil || !bytes.Equal(got, body) {
		t.Fatalf("contenido tras reemplazar: %q, error %v", got, err)
	}
}
