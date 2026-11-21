package server

import (
	"log"
	"sync"
	"time"
)

// Failed logins are counted per user name and per client IP, in memory; a restart clears them.
// Nothing is refused and nobody is locked out. Each failure delays the answer by a short pause
// that grows with the recent failures, which slows password guessing, and a warning is logged
// when a count reaches its limit.
const (
	loginMaxFailsUser = 10
	loginMaxFailsIP   = 50
	loginFailWindow   = 15 * time.Minute
	loginPauseStep    = 250 * time.Millisecond
	loginMaxPause     = 3 * time.Second
	loginMaxKeys      = 10000 // above this many tracked keys, expired entries are dropped
)

type loginFails struct {
	count       int
	windowStart time.Time
}

// loginGuard counts failed logins under the keys "ip:<addr>" and "user:<name>".
type loginGuard struct {
	mu    sync.Mutex
	fails map[string]*loginFails
}

var loginFailGuard = &loginGuard{fails: map[string]*loginFails{}}

// failed records a wrong password for user from ip and returns how long the answer must wait.
// Pass an empty user for an unknown name. It never refuses a login.
func (g *loginGuard) failed(user, ip string, now time.Time) time.Duration {
	g.mu.Lock()
	defer g.mu.Unlock()
	if len(g.fails) >= loginMaxKeys {
		g.prune(now)
	}
	n := g.countFail("ip:"+ip, loginMaxFailsIP, now)
	if user != "" {
		if m := g.countFail("user:"+user, loginMaxFailsUser, now); m > n {
			n = m
		}
	}
	return pauseFor(n)
}

// pauseFor grows with the number of recent failures and is capped at loginMaxPause.
func pauseFor(failures int) time.Duration {
	d := time.Duration(failures) * loginPauseStep
	if d > loginMaxPause {
		return loginMaxPause
	}
	return d
}

// succeeded forgets the failures of user after a correct password.
func (g *loginGuard) succeeded(user string) {
	g.mu.Lock()
	defer g.mu.Unlock()
	delete(g.fails, "user:"+user)
}

// countFail adds one failure to key and returns the failures inside the current window.
// Reaching the limit logs a warning; nothing is blocked. The caller holds mu.
func (g *loginGuard) countFail(key string, limit int, now time.Time) int {
	f := g.fails[key]
	if f == nil {
		f = &loginFails{}
		g.fails[key] = f
	}
	if now.Sub(f.windowStart) > loginFailWindow {
		f.windowStart, f.count = now, 0
	}
	f.count++
	if f.count == limit {
		log.Printf("auth: WARNING %d failed logins for %s in %s; possible password guessing. No login is refused.", f.count, key, loginFailWindow)
	}
	return f.count
}

// prune drops keys whose window has ended. The caller holds mu.
func (g *loginGuard) prune(now time.Time) {
	for k, f := range g.fails {
		if now.Sub(f.windowStart) > loginFailWindow {
			delete(g.fails, k)
		}
	}
}
