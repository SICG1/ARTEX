package report

import "testing"

func TestSummaryKeyCuentaTodo(t *testing.T) {
	casos := map[string]string{"Critical": "critical", " high ": "high", "MEDIUM": "medium", "low": "low", "": "info", "raro": "info"}
	for in, want := range casos {
		if got := summaryKey(in); got != want {
			t.Fatalf("summaryKey(%q) = %q, se esperaba %q", in, got, want)
		}
	}
}
