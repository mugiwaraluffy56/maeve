package compress

import (
	"context"
	"strings"
	"testing"

	"github.com/puneethaditya/maeve/session"
)

func TestCompressor_CompressPacksByImportance(t *testing.T) {
	t.Parallel()

	repo := fakeRepo{objects: []session.ContextObject{
		{ID: "low", Type: session.ObjectFileSnapshot, Content: []byte("low value"), Importance: 0.1, TokenCount: 2},
		{ID: "high", Type: session.ObjectDiffBlock, Content: []byte("high value"), Importance: 0.9, TokenCount: 2},
	}}
	payload, report, err := New(repo).Compress(context.Background(), "sess_1", 2)
	if err != nil {
		t.Fatalf("Compress() error = %v", err)
	}

	got := string(payload)
	if !strings.Contains(got, "high value") {
		t.Fatalf("payload = %q, want high value object", got)
	}
	if strings.Contains(got, "low value") {
		t.Fatalf("payload = %q, did not expect low value object", got)
	}
	if report.KeptObjects != 1 || report.DroppedObjects != 1 {
		t.Fatalf("report = %+v, want one kept and one dropped", report)
	}
}

func TestCompressor_CompressDeduplicatesContent(t *testing.T) {
	t.Parallel()

	repo := fakeRepo{objects: []session.ContextObject{
		{ID: "first", Type: session.ObjectFileSnapshot, Content: []byte("same"), Importance: 0.2, TokenCount: 1},
		{ID: "second", Type: session.ObjectFileSnapshot, Content: []byte("same"), Importance: 0.8, TokenCount: 1},
	}}
	_, report, err := New(repo).Compress(context.Background(), "sess_1", 10)
	if err != nil {
		t.Fatalf("Compress() error = %v", err)
	}

	if report.KeptObjects != 1 {
		t.Fatalf("KeptObjects = %d, want 1", report.KeptObjects)
	}
	if len(report.Reasons) != 1 || report.Reasons[0].Reason != "duplicate_content" {
		t.Fatalf("Reasons = %+v, want duplicate_content", report.Reasons)
	}
}

type fakeRepo struct {
	objects []session.ContextObject
}

func (f fakeRepo) ListContextObjects(context.Context, string) ([]session.ContextObject, error) {
	return f.objects, nil
}
