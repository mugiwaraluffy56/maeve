package watch

import "strings"

func ShouldIngest(path string) bool {
	if strings.Contains(path, "/.git/") || strings.Contains(path, "/.maeve/") {
		return false
	}
	return strings.HasSuffix(path, ".go") ||
		strings.HasSuffix(path, ".md") ||
		strings.HasSuffix(path, ".toml") ||
		strings.HasSuffix(path, ".diff") ||
		strings.HasSuffix(path, ".patch")
}
