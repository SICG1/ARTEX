package guard

import (
	"context"
	"database/sql"
	"encoding/json"
	"testing"

	"github.com/Autumn-27/artex/db"
	"github.com/Autumn-27/artex/intercept"
	_ "github.com/jackc/pgx/v5/stdlib"
)

// With the rule store unreadable the call is blocked (fail closed, v0.7.0 #124)
// and the block is recorded in the audit ring.
func TestUnreadableRulesBlockCall(t *testing.T) {
	conn, err := sql.Open("pgx", "postgres://test@/test?host=/nonexistent-artex-socket&sslmode=disable&connect_timeout=1")
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	g := NewWithInterceptor(intercept.New(&db.DB{DB: conn}))
	input, _ := json.Marshal(map[string]string{"command": "ls -la"})
	if blocked, _, _ := g.Hooks().PreToolUse(context.Background(), "Bash", input); !blocked {
		t.Fatal("the call went through while the rules cannot be read; it must be blocked")
	}
}
