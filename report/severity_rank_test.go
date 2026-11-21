package report

import "testing"

func TestSeverityRankNoSubeDesconocidas(t *testing.T) {
	casos := map[string]int{"critical": 0, "High": 1, " medium ": 2, "low": 3, "info": 4, "": 4, "desconocida": 4}
	for s, want := range casos {
		if got := severityRank(s); got != want {
			t.Fatalf("severityRank(%q) = %d, se esperaba %d", s, got, want)
		}
	}
}
