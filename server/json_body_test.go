package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// jsonBody caps a JSON body at maxJSONBody: one request cannot make the server decode
// an unbounded body, while an ordinary body still decodes.
func TestJSONBodyLimitaCuerpo(t *testing.T) {
	var v map[string]any
	grande := `{"x":"` + strings.Repeat("a", maxJSONBody) + `"}`
	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(grande))
	if err := json.NewDecoder(jsonBody(r)).Decode(&v); err == nil {
		t.Fatal("un cuerpo por encima del tope debe fallar al decodificar")
	}
	r = httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"x":"ok"}`))
	if err := json.NewDecoder(jsonBody(r)).Decode(&v); err != nil {
		t.Fatalf("un cuerpo normal debe decodificarse: %v", err)
	}
	// decode (the shared helper) enforces its own, larger cap.
	if err := decode(httptest.NewRecorder(), httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"x":"ok"}`)), &v); err != nil {
		t.Fatalf("decode con un cuerpo normal: %v", err)
	}
}
