package traffic

import "testing"

func TestHeaderValueIgnoraMayusculas(t *testing.T) {
	h := map[string]string{"Content-Type": "text/html"}
	if got := headerValue(h, "content-type"); got != "text/html" {
		t.Fatalf("headerValue = %q", got)
	}
	if got := headerValue(h, "accept"); got != "" {
		t.Fatalf("cabecera ausente debe dar vacio: %q", got)
	}
}
