package llmpool

import (
	"context"
	"testing"

	"github.com/Autumn-27/norma/llm"
)

// A call cancelled before any answer says nothing about the profile: it must not clear an open breaker.
func TestCancelledCallKeepsOpenBreaker(t *testing.T) {
	reg := NewRegistry(nil, nil)
	reg.Trip(3, "invalid key", true)
	prov := &fakeProv{name: "solo", errs: []error{context.Canceled}}
	p := New([]*Member{member(3, "solo", 1, prov)}, reg)
	if _, _, _, err := p.Complete(context.Background(), llm.CompletionRequest{}); err == nil {
		t.Fatal("expected the cancellation error to be returned")
	}
	if prov.calls == 0 {
		t.Fatal("the profile was never called, so the test proves nothing")
	}
	if reg.m[3] == nil || reg.m[3].Trips == 0 {
		t.Fatal("a cancelled call cleared the open breaker")
	}
}
