package agent

import (
	"sort"
	"testing"
)

// TestSchemaRequiredPropertiesExist verifica que cada campo listado en "required"
// exista realmente en "properties". Algunos proveedores (Gemini) rechazan el
// request completo cuando un esquema es inconsistente, asi que esto protege el
// contrato de herramientas.
func TestSchemaRequiredPropertiesExist(t *testing.T) {
	problemas := 0
	for agente, herramientas := range builtinToolsByAgent() {
		for _, h := range herramientas {
			esquema := h.InputSchema()
			if esquema == nil {
				continue
			}
			props, _ := esquema["properties"].(map[string]any)
			req, _ := esquema["required"].([]any)
			if len(req) == 0 {
				continue
			}
			nombres := make([]string, 0, len(req))
			for _, r := range req {
				if s, ok := r.(string); ok {
					nombres = append(nombres, s)
				}
			}
			sort.Strings(nombres)
			for _, n := range nombres {
				if _, existe := props[n]; !existe {
					t.Errorf("AGENTE=%s TOOL=%s: required %q no existe en properties (%v)", agente, h.Name(), n, claves(props))
					problemas++
				}
			}
		}
	}
	if problemas == 0 {
		t.Logf("sin inconsistencias en %d herramientas", contar())
	}
}

func claves(m map[string]any) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

func contar() int {
	n := 0
	for _, hs := range builtinToolsByAgent() {
		n += len(hs)
	}
	return n
}
