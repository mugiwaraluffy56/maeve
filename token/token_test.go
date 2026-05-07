package token

import "testing"

func TestCounter_Count(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		content []byte
		want    int
	}{
		{name: "empty", content: nil, want: 0},
		{name: "spaces only", content: []byte(" \n\t "), want: 0},
		{name: "words", content: []byte("alpha beta\ngamma"), want: 3},
	}

	counter := NewCounter()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := counter.Count(tt.content)
			if got != tt.want {
				t.Fatalf("Count() = %d, want %d", got, tt.want)
			}
		})
	}
}
