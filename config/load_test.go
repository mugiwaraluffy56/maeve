package config

import "testing"

func TestLoad_Defaults(t *testing.T) {
	t.Parallel()

	cfg, err := Load("")
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.Compression.DefaultBudgetTokens != 8000 {
		t.Fatalf("DefaultBudgetTokens = %d, want 8000", cfg.Compression.DefaultBudgetTokens)
	}
	if cfg.Server.Port != 7432 {
		t.Fatalf("Server.Port = %d, want 7432", cfg.Server.Port)
	}
}
