package parse

import "strings"

type Language string

const (
	LanguageGo         Language = "go"
	LanguageMarkdown   Language = "markdown"
	LanguageTypeScript Language = "typescript"
	LanguageUnknown    Language = "unknown"
)

func Detect(path string) Language {
	switch {
	case strings.HasSuffix(path, ".go"):
		return LanguageGo
	case strings.HasSuffix(path, ".md"):
		return LanguageMarkdown
	case strings.HasSuffix(path, ".ts"), strings.HasSuffix(path, ".tsx"):
		return LanguageTypeScript
	default:
		return LanguageUnknown
	}
}
