package session

import "errors"

var (
	ErrSessionNotFound  = errors.New("session not found")
	ErrSnapshotNotFound = errors.New("snapshot not found")
)
