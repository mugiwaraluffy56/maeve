package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/puneethaditya/maeve/session"
)

func (s *Store) CreateSession(ctx context.Context, sess session.Session) error {
	_, err := s.db.ExecContext(ctx, `
INSERT INTO sessions (id, name, repo_path, branch, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?)`,
		sess.ID, sess.Name, sess.RepoPath, sess.Branch, sess.CreatedAt, sess.UpdatedAt)
	if err != nil {
		return fmt.Errorf("create session: %w", err)
	}
	return nil
}

func (s *Store) GetSession(ctx context.Context, id string) (session.Session, error) {
	var sess session.Session
	err := s.db.QueryRowContext(ctx, `
SELECT id, name, repo_path, branch, created_at, updated_at
FROM sessions
WHERE id = ?`, id).Scan(
		&sess.ID,
		&sess.Name,
		&sess.RepoPath,
		&sess.Branch,
		&sess.CreatedAt,
		&sess.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return session.Session{}, session.ErrSessionNotFound
	}
	if err != nil {
		return session.Session{}, fmt.Errorf("get session: %w", err)
	}
	return sess, nil
}

func (s *Store) LatestSession(ctx context.Context, repoPath string) (session.Session, error) {
	var sess session.Session
	err := s.db.QueryRowContext(ctx, `
SELECT id, name, repo_path, branch, created_at, updated_at
FROM sessions
WHERE repo_path = ?
ORDER BY updated_at DESC
LIMIT 1`, repoPath).Scan(
		&sess.ID,
		&sess.Name,
		&sess.RepoPath,
		&sess.Branch,
		&sess.CreatedAt,
		&sess.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return session.Session{}, session.ErrSessionNotFound
	}
	if err != nil {
		return session.Session{}, fmt.Errorf("latest session: %w", err)
	}
	return sess, nil
}
