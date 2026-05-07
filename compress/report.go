package compress

type Report struct {
	OriginalTokens   int      `json:"original_tokens"`
	CompressedTokens int      `json:"compressed_tokens"`
	KeptObjects      int      `json:"kept_objects"`
	DroppedObjects   int      `json:"dropped_objects"`
	Reasons          []Reason `json:"reasons"`
}

type Reason struct {
	ObjectID string `json:"object_id"`
	Reason   string `json:"reason"`
}

func (r Report) SavingsPercent() float64 {
	if r.OriginalTokens == 0 {
		return 0
	}
	saved := r.OriginalTokens - r.CompressedTokens
	return float64(saved) / float64(r.OriginalTokens) * 100
}
