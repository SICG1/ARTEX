package config

import "testing"

func TestCheckDBPortRango(t *testing.T) {
	for _, p := range []int{0, 1, 5432, 65535} {
		if err := checkDBPort(p); err != nil {
			t.Fatalf("puerto %d debe ser valido: %v", p, err)
		}
	}
	for _, p := range []int{-1, 65536} {
		if err := checkDBPort(p); err == nil {
			t.Fatalf("puerto %d debe rechazarse", p)
		}
	}
}
