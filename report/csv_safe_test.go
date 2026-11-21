package report

import "testing"

func TestCSVSafeNeutralizaFormulas(t *testing.T) {
	casos := map[string]string{
		"=cmd|' /C calc'!A0": "'=cmd|' /C calc'!A0",
		"+1+1":               "'+1+1",
		"-2":                 "'-2",
		"@SUM(A1)":           "'@SUM(A1)",
		"\tx":                "'\tx",
		"\rx":                "'\rx",
		"Normal":             "Normal",
		"":                   "",
	}
	for in, want := range casos {
		if got := csvSafe(in); got != want {
			t.Fatalf("csvSafe(%q) = %q, se esperaba %q", in, got, want)
		}
	}
}
