package diff

import "testing"

func TestSummarize(t *testing.T) {
	t.Parallel()

	got := Summarize([]byte("diff --git a/a b/a\n--- a/a\n+++ b/a\n-old\n+new\n+line\n"))
	if got.FilesChanged != 1 || got.AddedLines != 2 || got.RemovedLines != 1 {
		t.Fatalf("Summarize() = %+v", got)
	}
}
