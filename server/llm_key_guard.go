package server

import "github.com/Autumn-27/artex/agent"

// llmKeyMayReuse reports whether a key saved for src may be sent to the endpoint cfg names.
// Only the same format and the same base URL qualify, so a saved key never leaves its endpoint.
func llmKeyMayReuse(cfg, src agent.Config) bool {
	return cfg.Format == src.Format && cfg.BaseURL == src.BaseURL
}
