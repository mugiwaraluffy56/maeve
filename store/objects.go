package store

import (
	"context"
	"fmt"

	"github.com/puneethaditya/maeve/session"
)

func (s *Store) AddContextObject(ctx context.Context, obj session.ContextObject) error {
	_, err := s.db.ExecContext(ctx, `
INSERT INTO context_objects
  (id, session_id, type, source_path, content, importance, token_count, created_at, accessed_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		obj.ID,
		obj.SessionID,
		string(obj.Type),
		obj.SourcePath,
		obj.Content,
		obj.Importance,
		obj.TokenCount,
		obj.CreatedAt,
		obj.AccessedAt,
	)
	if err != nil {
		return fmt.Errorf("add context object: %w", err)
	}
	return nil
}

func (s *Store) ListContextObjects(ctx context.Context, sessionID string) ([]session.ContextObject, error) {
	rows, err := s.db.QueryContext(ctx, `
SELECT id, session_id, type, source_path, content, importance, token_count, created_at, accessed_at
FROM context_objects
WHERE session_id = ?
ORDER BY importance DESC, created_at DESC`, sessionID)
	if err != nil {
		return nil, fmt.Errorf("list context objects: %w", err)
	}
	defer rows.Close()

	var objects []session.ContextObject
	for rows.Next() {
		var obj session.ContextObject
		var typ string
		if err := rows.Scan(
			&obj.ID,
			&obj.SessionID,
			&typ,
			&obj.SourcePath,
			&obj.Content,
			&obj.Importance,
			&obj.TokenCount,
			&obj.CreatedAt,
			&obj.AccessedAt,
		); err != nil {
			return nil, fmt.Errorf("scan context object: %w", err)
		}
		obj.Type = session.ContextObjectType(typ)
		objects = append(objects, obj)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate context objects: %w", err)
	}
	return objects, nil
}

func (s *Store) PruneContextObjects(ctx context.Context, sessionID string, threshold float64) (int64, error) {
	result, err := s.db.ExecContext(ctx, `
DELETE FROM context_objects
WHERE session_id = ? AND importance < ?`, sessionID, threshold)
	if err != nil {
		return 0, fmt.Errorf("prune context objects: %w", err)
	}
	count, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("count pruned objects: %w", err)
	}
	return count, nil
}
