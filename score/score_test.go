package score

import (
	"testing"
	"time"

	"github.com/mugiwaraluffy56/maeve/config"
	"github.com/mugiwaraluffy56/maeve/session"
)

func TestScorer_ScoreDiffAboveTerminal(t *testing.T) {
	t.Parallel()

	now := time.Unix(100, 0)
	scorer := New(config.Default())
	diff := session.ContextObject{Type: session.ObjectDiffBlock, CreatedAt: now.Unix()}
	terminal := session.ContextObject{Type: session.ObjectTerminalBlock, CreatedAt: now.Unix()}

	if scorer.Score(diff, now) <= scorer.Score(terminal, now) {
		t.Fatal("diff score should be higher than terminal score")
	}
}
