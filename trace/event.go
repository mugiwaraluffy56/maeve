package trace

import "time"

type Event struct {
	Time   time.Time      `json:"time"`
	Name   string         `json:"name"`
	Fields map[string]any `json:"fields,omitempty"`
}
