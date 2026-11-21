package server

import (
	"testing"
	"time"
)

// Failed logins never lock anyone out: the pause is capped and a later correct login is not refused.
func TestLoginGuardNeverRefuses(t *testing.T) {
	g := &loginGuard{fails: map[string]*loginFails{}}
	now := time.Now()
	for i := 0; i < 3*loginMaxFailsIP; i++ {
		g.failed("ARTEX", "192.0.2.1", now)
	}
	if d := g.failed("ARTEX", "192.0.2.2", now); d != loginMaxPause {
		t.Fatalf("pause after many failures = %v, want the cap %v", d, loginMaxPause)
	}
	g.succeeded("ARTEX")
	if d := g.failed("ARTEX", "192.0.2.3", now); d != loginPauseStep {
		t.Fatalf("pause after a correct login = %v, want the first step %v", d, loginPauseStep)
	}
}

// The pause grows with consecutive failures and stops at the cap.
func TestLoginGuardPauseGrowsAndCaps(t *testing.T) {
	g := &loginGuard{fails: map[string]*loginFails{}}
	now := time.Now()
	first := g.failed("", "192.0.2.9", now)
	second := g.failed("", "192.0.2.9", now)
	if !(second > first) {
		t.Fatalf("second pause %v must be longer than the first %v", second, first)
	}
	var last time.Duration
	for i := 0; i < 50; i++ {
		last = g.failed("", "192.0.2.9", now)
	}
	if last != loginMaxPause {
		t.Fatalf("pause = %v, want the cap %v", last, loginMaxPause)
	}
}

// After the failure window ends the count starts again.
func TestLoginGuardWindowResets(t *testing.T) {
	g := &loginGuard{fails: map[string]*loginFails{}}
	now := time.Now()
	for i := 0; i < 20; i++ {
		g.failed("", "192.0.2.10", now)
	}
	later := now.Add(loginFailWindow + time.Second)
	if d := g.failed("", "192.0.2.10", later); d != loginPauseStep {
		t.Fatalf("pause after the window = %v, want %v", d, loginPauseStep)
	}
}
