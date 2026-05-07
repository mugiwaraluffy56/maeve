package mcp

import (
	"encoding/json"
	"io"
)

type Server struct{}

func NewServer() Server {
	return Server{}
}

func (Server) ListTools(w io.Writer) error {
	return json.NewEncoder(w).Encode(map[string]any{
		"tools": Tools(),
	})
}
