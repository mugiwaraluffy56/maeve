package mcp

import "testing"

func TestTools(t *testing.T) {
	t.Parallel()

	tools := Tools()
	if len(tools) != 3 {
		t.Fatalf("len(Tools()) = %d, want 3", len(tools))
	}
	if tools[0].Name != "maeve_compress" {
		t.Fatalf("first tool = %q, want maeve_compress", tools[0].Name)
	}
}
