package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/mugiwaraluffy56/maeve/session"
)

func (s *Store) CreateSnapshot(ctx context.Context, snap session.Snapshot) error {
	_, err := s.db.ExecContext(ctx, `
INSERT INTO snapshots
  (id, session_id, name, compressed_payload, original_tokens, compressed_tokens, created_at)
VALUES (?, ?, ?, ?, ?, ?, ?)`,
		snap.ID,
		snap.SessionID,
		snap.Name,
		snap.CompressedPayload,
		snap.OriginalTokens,
		snap.CompressedTokens,
		snap.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("create snapshot: %w", err)
	}
	return nil
}

func (s *Store) GetSnapshot(ctx context.Context, id string) (session.Snapshot, error) {
	var snap session.Snapshot
	err := s.db.QueryRowContext(ctx, `
SELECT id, session_id, name, compressed_payload, original_tokens, compressed_tokens, created_at
FROM snapshots
WHERE id = ?`, id).Scan(
		&snap.ID,
		&snap.SessionID,
		&snap.Name,
		&snap.CompressedPayload,
		&snap.OriginalTokens,
		&snap.CompressedTokens,
		&snap.CreatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return session.Snapshot{}, session.ErrSnapshotNotFound
	}
	if err != nil {
		return session.Snapshot{}, fmt.Errorf("get snapshot: %w", err)
	}
	return snap, nil
}

func (s *Store) ListSnapshots(ctx context.Context, sessionID string) ([]session.Snapshot, error) {
	rows, err := s.db.QueryContext(ctx, `
SELECT id, session_id, name, compressed_payload, original_tokens, compressed_tokens, created_at
FROM snapshots
WHERE session_id = ?
ORDER BY created_at DESC`, sessionID)
	if err != nil {
		return nil, fmt.Errorf("list snapshots: %w", err)
	}
	defer rows.Close()

	var snapshots []session.Snapshot
	for rows.Next() {
		var snap session.Snapshot
		if err := rows.Scan(
			&snap.ID,
			&snap.SessionID,
			&snap.Name,
			&snap.CompressedPayload,
			&snap.OriginalTokens,
			&snap.CompressedTokens,
			&snap.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan snapshot: %w", err)
		}
		snapshots = append(snapshots, snap)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate snapshots: %w", err)
	}
	return snapshots, nil
}
