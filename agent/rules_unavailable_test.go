package agent

import (
	"context"
	"database/sql"
	"encoding/json"
	"strings"
	"testing"

	"github.com/Autumn-27/artex/db"
	_ "github.com/jackc/pgx/v5/stdlib"
)

// deadDB returns a DB whose queries always fail. The socket directory does not
// exist, so the driver fails at once and opens no network connection.
func deadDB(t *testing.T) *db.DB {
	t.Helper()
	conn, err := sql.Open("pgx", "postgres://test@/test?host=/nonexistent-artex-socket&sslmode=disable&connect_timeout=1")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = conn.Close() })
	return &db.DB{DB: conn}
}

// When the asset rules cannot be read, the insert goes on and the result carries a warning
// (it is not refused).
func TestInsertAssetsWarnsWhenRulesUnreadable(t *testing.T) {
	d := deadDB(t)
	ts := NewToolSet(nil, "")
	ts.SetAssetStore(d.Assets(), d.Companies())
	raw, _ := json.Marshal(map[string]any{"assets": []any{map[string]any{"type": "root_domain", "domain": "warn-only.example"}}})
	res, err := ts.insertAssets().Call(context.Background(), raw, nil)
	if err != nil {
		t.Fatalf("Call error: %v", err)
	}
	text := res.Flatten()
	if !strings.Contains(text, "could not be read") {
		t.Fatalf("no warning about unreadable rules in the result: %s", text)
	}
	if strings.Contains(text, "insertion blocked") {
		t.Fatalf("insert was refused for unreadable rules; it must go on with a warning: %s", text)
	}
}

// When the operator constraints cannot be read, the prompt must say so instead of omitting them.
func TestConstraintBlockSaysWhenUnreadable(t *testing.T) {
	if got := constraintBlock(deadDB(t).Exploration(1)); !strings.Contains(got, "could not be loaded") {
		t.Fatalf("constraintBlock = %q, want a notice that the constraints could not be loaded", got)
	}
	if got := constraintBlock(nil); got != "" {
		t.Fatalf("constraintBlock(nil) = %q, want empty", got)
	}
}
