package intercept

import (
	"database/sql"
	"testing"

	"github.com/Autumn-27/artex/db"
	_ "github.com/jackc/pgx/v5/stdlib"
)

// deadStore returns a DB whose queries always fail. The socket directory does not
// exist, so the driver fails at once and opens no network connection.
func deadStore(t *testing.T) *db.DB {
	t.Helper()
	conn, err := sql.Open("pgx", "postgres://test@/test?host=/nonexistent-artex-socket&sslmode=disable&connect_timeout=1")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = conn.Close() })
	return &db.DB{DB: conn}
}

// With the rule store down, the call is denied: the gate keeps the tool in the rule
// path and Match returns a deny, so a call is never let through unchecked
// (fail closed, v0.7.0 #124).
func TestRuleLoadFailureFailsClosed(t *testing.T) {
	i := New(deadStore(t))
	if !i.IsToolEnabled("Bash") {
		t.Fatal("IsToolEnabled = false while rules cannot be loaded; the tool would skip the rule path")
	}
	dec, matched := i.Match("Bash", []byte(`{"command":"ls"}`))
	if !matched || dec.Action != "deny" {
		t.Fatalf("Match = (%+v, %v), want a deny while rules cannot be read", dec, matched)
	}
}
