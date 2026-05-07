package session

type ContextObjectType string

const (
	ObjectFileSnapshot     ContextObjectType = "file_snapshot"
	ObjectTerminalBlock    ContextObjectType = "terminal_block"
	ObjectConversationTurn ContextObjectType = "conversation_turn"
	ObjectDiffBlock        ContextObjectType = "diff_block"
	ObjectSymbol           ContextObjectType = "symbol"
)

type Session struct {
	ID        string
	Name      string
	RepoPath  string
	Branch    string
	CreatedAt int64
	UpdatedAt int64
}

type ContextObject struct {
	ID         string
	SessionID  string
	Type       ContextObjectType
	SourcePath string
	Content    []byte
	Importance float64
	TokenCount int
	CreatedAt  int64
	AccessedAt int64
}

type Snapshot struct {
	ID                string
	SessionID         string
	Name              string
	CompressedPayload []byte
	OriginalTokens    int
	CompressedTokens  int
	CreatedAt         int64
}
