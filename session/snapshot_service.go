package session

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

type SnapshotRepository interface {
	CreateSnapshot(context.Context, Snapshot) error
	ListSnapshots(context.Context, string) ([]Snapshot, error)
}

type SnapshotService struct {
	repo  SnapshotRepository
	clock Clock
}

func NewSnapshotService(repo SnapshotRepository, clock Clock) SnapshotService {
	return SnapshotService{repo: repo, clock: clock}
}

func (s SnapshotService) Save(ctx context.Context, sessionID, name string, payload []byte, originalTokens, compressedTokens int) (Snapshot, error) {
	now := s.clock.Now().Unix()
	snap := Snapshot{
		ID:                "snap_" + uuid.NewString(),
		SessionID:         sessionID,
		Name:              name,
		CompressedPayload: payload,
		OriginalTokens:    originalTokens,
		CompressedTokens:  compressedTokens,
		CreatedAt:         now,
	}
	if err := s.repo.CreateSnapshot(ctx, snap); err != nil {
		return Snapshot{}, fmt.Errorf("save snapshot: %w", err)
	}
	return snap, nil
}

func (s SnapshotService) List(ctx context.Context, sessionID string) ([]Snapshot, error) {
	return s.repo.ListSnapshots(ctx, sessionID)
}
