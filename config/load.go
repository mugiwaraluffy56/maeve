package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

func Load(path string) (Config, error) {
	cfg := Default()

	v := viper.New()
	v.SetConfigType("toml")
	v.SetEnvPrefix("MAEVE")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	setDefaults(v, cfg)
	if path != "" {
		v.SetConfigFile(path)
		if err := v.ReadInConfig(); err != nil {
			return Config{}, fmt.Errorf("read config: %w", err)
		}
	}

	if err := v.Unmarshal(&cfg); err != nil {
		return Config{}, fmt.Errorf("decode config: %w", err)
	}
	return cfg, nil
}

func setDefaults(v *viper.Viper, cfg Config) {
	v.SetDefault("session.auto_init", cfg.Session.AutoInit)
	v.SetDefault("session.decay_window_hours", cfg.Session.DecayWindowHours)
	v.SetDefault("session.decay_factor", cfg.Session.DecayFactor)
	v.SetDefault("compression.default_budget_tokens", cfg.Compression.DefaultBudgetTokens)
	v.SetDefault("compression.dedup_threshold", cfg.Compression.DedupThreshold)
	v.SetDefault("compression.truncation_threshold", cfg.Compression.TruncationThreshold)
	v.SetDefault("scoring.weight_recency", cfg.Scoring.WeightRecency)
	v.SetDefault("scoring.weight_access_freq", cfg.Scoring.WeightAccessFreq)
	v.SetDefault("scoring.weight_symbol_type", cfg.Scoring.WeightSymbolType)
	v.SetDefault("scoring.weight_diff_coverage", cfg.Scoring.WeightDiffCoverage)
	v.SetDefault("scoring.weight_agent_refs", cfg.Scoring.WeightAgentRefs)
	v.SetDefault("server.port", cfg.Server.Port)
	v.SetDefault("server.mcp_port", cfg.Server.MCPPort)
	v.SetDefault("tokenizer.model", cfg.Tokenizer.Model)
}
