package server

import (
	"testing"

	"github.com/Autumn-27/artex/agent"
)

func TestLLMKeyMayReuseSoloMismoDestino(t *testing.T) {
	guardada := agent.ConfigFrom("openai", "gpt-4o", "https://api.openai.com/v1/", "", "")
	if !llmKeyMayReuse(agent.ConfigFrom("openai", "gpt-4o", "https://api.openai.com/v1", "", ""), guardada) {
		t.Fatal("misma URL y formato: la clave guardada debe poder reutilizarse")
	}
	if llmKeyMayReuse(agent.ConfigFrom("openai", "gpt-4o", "https://evil.example.com/v1", "", ""), guardada) {
		t.Fatal("otra URL: la clave guardada no debe enviarse")
	}
	if llmKeyMayReuse(agent.ConfigFrom("openai-responses", "gpt-4o", "https://api.openai.com/v1", "", ""), guardada) {
		t.Fatal("otro formato: la clave guardada no debe enviarse")
	}
}
