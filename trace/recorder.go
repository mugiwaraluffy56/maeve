package trace

import "sync"

type Recorder struct {
	mu     sync.Mutex
	events []Event
}

func NewRecorder() *Recorder {
	return &Recorder{}
}

func (r *Recorder) Record(event Event) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.events = append(r.events, event)
}

func (r *Recorder) Events() []Event {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := make([]Event, len(r.events))
	copy(out, r.events)
	return out
}
