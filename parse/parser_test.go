package parse

import "testing"

func TestParser_SymbolsFindsGoFunctions(t *testing.T) {
	t.Parallel()

	symbols := NewParser().Symbols("main.go", []byte("package main\n\nfunc Run() {}\n"))
	if len(symbols) != 1 {
		t.Fatalf("len(Symbols()) = %d, want 1", len(symbols))
	}
	if symbols[0].Name != "Run" {
		t.Fatalf("symbol name = %q, want Run", symbols[0].Name)
	}
}
