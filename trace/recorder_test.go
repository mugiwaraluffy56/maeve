package trace

import (
	"testing"
	"time"
)

func TestRecorder_EventsReturnsCopy(t *testing.T) {
	t.Parallel()

	rec := NewRecorder()
	rec.Record(Event{Time: time.Unix(1, 0), Name: "ingest"})

	events := rec.Events()
	events[0].Name = "mutated"

	if got := rec.Events()[0].Name; got != "ingest" {
		t.Fatalf("stored event = %q, want ingest", got)
	}
}
