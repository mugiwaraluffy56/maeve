package score

import (
	"math"
	"time"

	"github.com/mugiwaraluffy56/maeve/config"
	"github.com/mugiwaraluffy56/maeve/session"
)

type Scorer struct {
	cfg config.Config
}

func New(cfg config.Config) Scorer {
	return Scorer{cfg: cfg}
}

func (s Scorer) Score(obj session.ContextObject, now time.Time) float64 {
	recency := s.recency(obj, now)
	typeScore := symbolTypeScore(obj.Type)
	base := s.cfg.Scoring.WeightRecency*recency +
		s.cfg.Scoring.WeightSymbolType*typeScore +
		s.cfg.Scoring.WeightAccessFreq*0.5 +
		s.cfg.Scoring.WeightDiffCoverage*diffScore(obj.Type) +
		s.cfg.Scoring.WeightAgentRefs*0.5
	return clamp(base, 0, 1)
}

func (s Scorer) recency(obj session.ContextObject, now time.Time) float64 {
	created := time.Unix(obj.CreatedAt, 0)
	age := now.Sub(created)
	if age <= 0 {
		return 1
	}
	window := s.cfg.DecayWindow()
	if window <= 0 {
		return 1
	}
	return math.Pow(s.cfg.Session.DecayFactor, age.Hours()/window.Hours())
}

func symbolTypeScore(typ session.ContextObjectType) float64 {
	switch typ {
	case session.ObjectSymbol:
		return 0.9
	case session.ObjectDiffBlock:
		return 0.85
	case session.ObjectFileSnapshot:
		return 0.65
	case session.ObjectConversationTurn:
		return 0.55
	case session.ObjectTerminalBlock:
		return 0.45
	default:
		return 0.5
	}
}

func diffScore(typ session.ContextObjectType) float64 {
	if typ == session.ObjectDiffBlock {
		return 1
	}
	return 0
}

func clamp(v, min, max float64) float64 {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}
