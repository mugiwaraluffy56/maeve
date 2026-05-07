package parse

import "strings"

type Parser struct{}

func NewParser() Parser {
	return Parser{}
}

func (Parser) Symbols(path string, content []byte) []Symbol {
	lang := Detect(path)
	if lang != LanguageGo {
		return nil
	}

	lines := strings.Split(string(content), "\n")
	var symbols []Symbol
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "func ") {
			symbols = append(symbols, Symbol{
				Name:      symbolName(trimmed, "func "),
				Kind:      "function",
				Language:  lang,
				StartLine: i + 1,
				EndLine:   i + 1,
			})
		}
	}
	return symbols
}

func symbolName(line, prefix string) string {
	name := strings.TrimPrefix(line, prefix)
	if idx := strings.IndexAny(name, " ("); idx >= 0 {
		return name[:idx]
	}
	return name
}
