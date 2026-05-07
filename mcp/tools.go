package mcp

type Tool struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	InputSchema map[string]any `json:"inputSchema"`
}

func Tools() []Tool {
	return []Tool{
		{
			Name:        "maeve_compress",
			Description: "Compress current session context to a token budget.",
			InputSchema: objectSchema(map[string]any{
				"budget_tokens": map[string]any{"type": "integer"},
				"session_id":    map[string]any{"type": "string"},
			}, []string{"budget_tokens"}),
		},
		{
			Name:        "maeve_snapshot",
			Description: "Save named snapshot of current session state.",
			InputSchema: objectSchema(map[string]any{
				"name":       map[string]any{"type": "string"},
				"session_id": map[string]any{"type": "string"},
			}, []string{"name"}),
		},
		{
			Name:        "maeve_status",
			Description: "Get current token usage and object count.",
			InputSchema: objectSchema(map[string]any{
				"session_id": map[string]any{"type": "string"},
			}, nil),
		},
	}
}

func objectSchema(properties map[string]any, required []string) map[string]any {
	return map[string]any{
		"type":       "object",
		"properties": properties,
		"required":   required,
	}
}
