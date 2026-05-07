package parse

type Symbol struct {
	Name      string
	Kind      string
	Language  Language
	StartLine int
	EndLine   int
}
