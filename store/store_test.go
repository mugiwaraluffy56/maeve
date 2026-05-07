package store

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/mugiwaraluffy56/maeve/session"
)

func TestStore_SessionLifecycle(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	st := openTestStore(t)

	sess := session.Session{
		ID:        "sess_1",
		Name:      "maeve",
		RepoPath:  "/tmp/maeve",
		Branch:    "main",
		CreatedAt: 1,
		UpdatedAt: 2,
	}
	if err := st.CreateSession(ctx, sess); err != nil {
		t.Fatalf("CreateSession() error = %v", err)
	}

	got, err := st.GetSession(ctx, sess.ID)
	if err != nil {
		t.Fatalf("GetSession() error = %v", err)
	}
	if got.ID != sess.ID || got.RepoPath != sess.RepoPath {
		t.Fatalf("GetSession() = %+v, want %+v", got, sess)
	}
}

func TestStore_ContextStatusAndPrune(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	st := openTestStore(t)
	createTestSession(t, st)

	objects := []session.ContextObject{
		{ID: "obj_1", SessionID: "sess_1", Type: session.ObjectFileSnapshot, Content: []byte("alpha"), Importance: 0.8, TokenCount: 1, CreatedAt: 1, AccessedAt: 1},
		{ID: "obj_2", SessionID: "sess_1", Type: session.ObjectTerminalBlock, Content: []byte("beta gamma"), Importance: 0.2, TokenCount: 2, CreatedAt: 2, AccessedAt: 2},
	}
	for _, obj := range objects {
		if err := st.AddContextObject(ctx, obj); err != nil {
			t.Fatalf("AddContextObject() error = %v", err)
		}
	}

	status, err := st.Status(ctx, "sess_1")
	if err != nil {
		t.Fatalf("Status() error = %v", err)
	}
	if status.ObjectCount != 2 || status.TokenCount != 3 {
		t.Fatalf("Status() = %+v, want 2 objects and 3 tokens", status)
	}

	pruned, err := st.PruneContextObjects(ctx, "sess_1", 0.5)
	if err != nil {
		t.Fatalf("PruneContextObjects() error = %v", err)
	}
	if pruned != 1 {
		t.Fatalf("PruneContextObjects() = %d, want 1", pruned)
	}
}

func TestStore_SnapshotLifecycle(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	st := openTestStore(t)
	createTestSession(t, st)

	snap := session.Snapshot{
		ID:                "snap_1",
		SessionID:         "sess_1",
		Name:              "before refactor",
		CompressedPayload: []byte("payload"),
		OriginalTokens:    10,
		CompressedTokens:  3,
		CreatedAt:         5,
	}
	if err := st.CreateSnapshot(ctx, snap); err != nil {
		t.Fatalf("CreateSnapshot() error = %v", err)
	}

	got, err := st.GetSnapshot(ctx, snap.ID)
	if err != nil {
		t.Fatalf("GetSnapshot() error = %v", err)
	}
	if got.Name != snap.Name || string(got.CompressedPayload) != "payload" {
		t.Fatalf("GetSnapshot() = %+v, want %+v", got, snap)
	}
}

func openTestStore(t *testing.T) *Store {
	t.Helper()

	st, err := Open(context.Background(), filepath.Join(t.TempDir(), "maeve.db"))
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	t.Cleanup(func() {
		if err := st.Close(); err != nil {
			t.Fatalf("Close() error = %v", err)
		}
	})
	return st
}

func createTestSession(t *testing.T, st *Store) {
	t.Helper()

	err := st.CreateSession(context.Background(), session.Session{
		ID:        "sess_1",
		Name:      "maeve",
		RepoPath:  "/tmp/maeve",
		Branch:    "main",
		CreatedAt: 1,
		UpdatedAt: 1,
	})
	if err != nil {
		t.Fatalf("CreateSession() error = %v", err)
	}
}
