package traffic

import "testing"

func TestLikeContainsEscapaComodines(t *testing.T) {
	if got := likeContains(`a%_b\c`); got != `%a\%\_b\\c%` {
		t.Fatalf("likeContains = %q", got)
	}
}
