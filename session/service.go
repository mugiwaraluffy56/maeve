package session

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

type Repository interface {
	CreateSession(context.Context, Session) error
	LatestSession(context.Context, string) (Session, error)
}

type Service struct {
	repo  Repository
	clock Clock
}

func NewService(repo Repository, clock Clock) Service {
	return Service{repo: repo, clock: clock}
}

func (s Service) Init(ctx context.Context, name, repoPath, branch string) (Session, error) {
	now := s.clock.Now().Unix()
	sess := Session{
		ID:        "sess_" + uuid.NewString(),
		Name:      name,
		RepoPath:  repoPath,
		Branch:    branch,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := s.repo.CreateSession(ctx, sess); err != nil {
		return Session{}, fmt.Errorf("init session: %w", err)
	}
	return sess, nil
}

func (s Service) Resolve(ctx context.Context, repoPath string) (Session, error) {
	return s.repo.LatestSession(ctx, repoPath)
}
