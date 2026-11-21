package traffic

import (
	"os"
	"testing"
)

func TestLegacyWalkErrorSoloOmiteInexistentes(t *testing.T) {
	if err := legacyWalkError("x", os.ErrNotExist); err != nil {
		t.Fatal("una entrada que ya no existe debe omitirse")
	}
	if err := legacyWalkError("x", os.ErrPermission); err == nil {
		t.Fatal("un error de lectura debe abortar el GC")
	}
}
