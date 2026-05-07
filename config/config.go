package config

import "time"

type Config struct {
	Session     SessionConfig
	Compression CompressionConfig
	Scoring     ScoringConfig
	Server      ServerConfig
	Tokenizer   TokenizerConfig
}

type SessionConfig struct {
	AutoInit         bool
	DecayWindowHours int
	DecayFactor      float64
}

type CompressionConfig struct {
	DefaultBudgetTokens int
	DedupThreshold      float64
	TruncationThreshold float64
}

type ScoringConfig struct {
	WeightRecency      float64
	WeightAccessFreq   float64
	WeightSymbolType   float64
	WeightDiffCoverage float64
	WeightAgentRefs    float64
}

type ServerConfig struct {
	Port    int
	MCPPort int
}

type TokenizerConfig struct {
	Model string
}

func Default() Config {
	return Config{
		Session: SessionConfig{
			AutoInit:         true,
			DecayWindowHours: 2,
			DecayFactor:      0.7,
		},
		Compression: CompressionConfig{
			DefaultBudgetTokens: 8000,
			DedupThreshold:      0.85,
			TruncationThreshold: 0.4,
		},
		Scoring: ScoringConfig{
			WeightRecency:      0.30,
			WeightAccessFreq:   0.20,
			WeightSymbolType:   0.25,
			WeightDiffCoverage: 0.15,
			WeightAgentRefs:    0.10,
		},
		Server: ServerConfig{
			Port:    7432,
			MCPPort: 7433,
		},
		Tokenizer: TokenizerConfig{
			Model: "cl100k_base",
		},
	}
}

func (c Config) DecayWindow() time.Duration {
	return time.Duration(c.Session.DecayWindowHours) * time.Hour
}
