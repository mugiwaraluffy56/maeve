package compress

import (
	"bytes"
	"context"
	"crypto/sha256"
	"fmt"
	"sort"

	"github.com/mugiwaraluffy56/maeve/session"
)

type Repository interface {
	ListContextObjects(context.Context, string) ([]session.ContextObject, error)
}

type Compressor struct {
	repo Repository
}

func New(repo Repository) Compressor {
	return Compressor{repo: repo}
}

func (c Compressor) Compress(ctx context.Context, sessionID string, budget int) ([]byte, Report, error) {
	if budget < 0 {
		return nil, Report{}, fmt.Errorf("budget must be non-negative")
	}

	objects, err := c.repo.ListContextObjects(ctx, sessionID)
	if err != nil {
		return nil, Report{}, fmt.Errorf("load objects: %w", err)
	}

	originalTokens := totalTokens(objects)
	objects, reasons := dedup(objects)
	sort.SliceStable(objects, func(i, j int) bool {
		if objects[i].Importance == objects[j].Importance {
			return objects[i].CreatedAt > objects[j].CreatedAt
		}
		return objects[i].Importance > objects[j].Importance
	})

	var payload bytes.Buffer
	report := Report{OriginalTokens: originalTokens, Reasons: reasons}
	for _, obj := range objects {
		if report.CompressedTokens+obj.TokenCount > budget {
			report.DroppedObjects++
			report.Reasons = append(report.Reasons, Reason{ObjectID: obj.ID, Reason: "budget_exhausted"})
			continue
		}
		writeObject(&payload, obj)
		report.CompressedTokens += obj.TokenCount
		report.KeptObjects++
	}

	return payload.Bytes(), report, nil
}

func dedup(objects []session.ContextObject) ([]session.ContextObject, []Reason) {
	seen := map[[32]byte]session.ContextObject{}
	var reasons []Reason
	for _, obj := range objects {
		hash := sha256.Sum256(obj.Content)
		existing, ok := seen[hash]
		if !ok || obj.Importance > existing.Importance {
			if ok {
				reasons = append(reasons, Reason{ObjectID: existing.ID, Reason: "duplicate_content"})
			}
			seen[hash] = obj
			continue
		}
		reasons = append(reasons, Reason{ObjectID: obj.ID, Reason: "duplicate_content"})
	}

	deduped := make([]session.ContextObject, 0, len(seen))
	for _, obj := range seen {
		deduped = append(deduped, obj)
	}
	return deduped, reasons
}

func totalTokens(objects []session.ContextObject) int {
	total := 0
	for _, obj := range objects {
		total += obj.TokenCount
	}
	return total
}

func writeObject(buf *bytes.Buffer, obj session.ContextObject) {
	if buf.Len() > 0 {
		buf.WriteString("\n\n")
	}
	if obj.SourcePath != "" {
		fmt.Fprintf(buf, "### %s %s\n", obj.Type, obj.SourcePath)
	} else {
		fmt.Fprintf(buf, "### %s\n", obj.Type)
	}
	buf.Write(obj.Content)
}
