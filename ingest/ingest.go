package ingest

import (
	"context"
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/mugiwaraluffy56/maeve/score"
	"github.com/mugiwaraluffy56/maeve/session"
	"github.com/mugiwaraluffy56/maeve/token"
)

type Repository interface {
	AddContextObject(context.Context, session.ContextObject) error
}

type Service struct {
	repo    Repository
	clock   session.Clock
	counter token.Counter
	scorer  score.Scorer
}

func New(repo Repository, clock session.Clock, counter token.Counter, scorer score.Scorer) Service {
	return Service{repo: repo, clock: clock, counter: counter, scorer: scorer}
}

func (s Service) File(ctx context.Context, sessionID, path string) (session.ContextObject, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return session.ContextObject{}, fmt.Errorf("read file: %w", err)
	}
	return s.Object(ctx, sessionID, session.ObjectFileSnapshot, path, content)
}

func (s Service) Diff(ctx context.Context, sessionID string, content []byte) (session.ContextObject, error) {
	return s.Object(ctx, sessionID, session.ObjectDiffBlock, "", content)
}

func (s Service) Object(ctx context.Context, sessionID string, typ session.ContextObjectType, sourcePath string, content []byte) (session.ContextObject, error) {
	now := s.clock.Now()
	obj := session.ContextObject{
		ID:         "obj_" + uuid.NewString(),
		SessionID:  sessionID,
		Type:       typ,
		SourcePath: sourcePath,
		Content:    content,
		TokenCount: s.counter.Count(content),
		CreatedAt:  now.Unix(),
		AccessedAt: now.Unix(),
	}
	obj.Importance = s.scorer.Score(obj, now)
	if err := s.repo.AddContextObject(ctx, obj); err != nil {
		return session.ContextObject{}, fmt.Errorf("ingest object: %w", err)
	}
	return obj, nil
}
