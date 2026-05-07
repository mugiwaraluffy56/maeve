package diff

import "strings"

type Summary struct {
	FilesChanged int
	AddedLines   int
	RemovedLines int
}

func Summarize(patch []byte) Summary {
	var summary Summary
	for _, line := range strings.Split(string(patch), "\n") {
		switch {
		case strings.HasPrefix(line, "diff --git "):
			summary.FilesChanged++
		case strings.HasPrefix(line, "+++") || strings.HasPrefix(line, "---"):
			continue
		case strings.HasPrefix(line, "+"):
			summary.AddedLines++
		case strings.HasPrefix(line, "-"):
			summary.RemovedLines++
		}
	}
	return summary
}
