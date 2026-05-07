package cli

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/puneethaditya/maeve/compress"
	"github.com/puneethaditya/maeve/config"
	"github.com/puneethaditya/maeve/ingest"
	"github.com/puneethaditya/maeve/score"
	"github.com/puneethaditya/maeve/session"
	"github.com/puneethaditya/maeve/store"
	"github.com/puneethaditya/maeve/token"
)

type app struct {
	cfg       config.Config
	store     *store.Store
	repoPath  string
	session   session.Service
	ingest    ingest.Service
	compress  compress.Compressor
	snapshots session.SnapshotService
}

func newApp(ctx context.Context, configPath, storePath string) (*app, error) {
	cfg, err := config.Load(configPath)
	if err != nil {
		return nil, err
	}
	repoPath, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("resolve repo path: %w", err)
	}
	if storePath == "" {
		storePath = filepath.Join(repoPath, ".maeve", "maeve.db")
	}

	st, err := store.Open(ctx, storePath)
	if err != nil {
		return nil, err
	}

	clock := session.RealClock{}
	counter := token.NewCounter()
	scorer := score.New(cfg)
	return &app{
		cfg:       cfg,
		store:     st,
		repoPath:  repoPath,
		session:   session.NewService(st, clock),
		ingest:    ingest.New(st, clock, counter, scorer),
		compress:  compress.New(st),
		snapshots: session.NewSnapshotService(st, clock),
	}, nil
}

func (a *app) close() error {
	return a.store.Close()
}

func (a *app) resolveSessionID(ctx context.Context, explicit string) (string, error) {
	if explicit != "" {
		return explicit, nil
	}
	sess, err := a.session.Resolve(ctx, a.repoPath)
	if err != nil {
		return "", err
	}
	return sess.ID, nil
}
