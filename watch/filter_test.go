package watch

import "testing"

func TestShouldIngest(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		path string
		want bool
	}{
		{name: "go file", path: "/repo/main.go", want: true},
		{name: "git file", path: "/repo/.git/config", want: false},
		{name: "database file", path: "/repo/.maeve/maeve.db", want: false},
		{name: "image", path: "/repo/logo.png", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := ShouldIngest(tt.path); got != tt.want {
				t.Fatalf("ShouldIngest(%q) = %v, want %v", tt.path, got, tt.want)
			}
		})
	}
}
