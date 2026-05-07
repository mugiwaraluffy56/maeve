package store

import (
	"context"
	"fmt"
)

type Status struct {
	SessionID   string
	ObjectCount int
	TokenCount  int
	MinScore    float64
	MaxScore    float64
	AvgScore    float64
}

func (s *Store) Status(ctx context.Context, sessionID string) (Status, error) {
	var status Status
	status.SessionID = sessionID

	err := s.db.QueryRowContext(ctx, `
SELECT
  COUNT(*),
  COALESCE(SUM(token_count), 0),
  COALESCE(MIN(importance), 0),
  COALESCE(MAX(importance), 0),
  COALESCE(AVG(importance), 0)
FROM context_objects
WHERE session_id = ?`, sessionID).Scan(
		&status.ObjectCount,
		&status.TokenCount,
		&status.MinScore,
		&status.MaxScore,
		&status.AvgScore,
	)
	if err != nil {
		return Status{}, fmt.Errorf("load status: %w", err)
	}
	return status, nil
}
