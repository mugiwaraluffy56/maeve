package token

import "strings"

type Counter struct{}

func NewCounter() Counter {
	return Counter{}
}

func (Counter) Count(content []byte) int {
	text := strings.TrimSpace(string(content))
	if text == "" {
		return 0
	}
	return len(strings.Fields(text))
}
