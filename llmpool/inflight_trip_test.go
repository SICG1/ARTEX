package llmpool

import "testing"

// A failure that arrives while the breaker is already open must not extend the cooldown.
func TestTripWhileOpenDoesNotExtendCooldown(t *testing.T) {
	r := NewRegistry(nil, nil)
	if !r.Trip(7, "invalid key", true) {
		t.Fatal("a hard failure should open the breaker")
	}
	before, trips := r.m[7].OpenUntil, r.m[7].Trips
	if r.Trip(7, "late 429", false) {
		t.Fatal("a failure while the breaker is open must not count as a new trip")
	}
	if r.m[7].Trips != trips || !r.m[7].OpenUntil.Equal(before) {
		t.Fatalf("cooldown changed: trips %d -> %d, open until %v -> %v", trips, r.m[7].Trips, before, r.m[7].OpenUntil)
	}
}
