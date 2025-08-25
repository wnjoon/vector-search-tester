package model

// EmbeddingRequest is a struct for embedding request
type EmbeddingRequest struct {
	Text string `json:"text"`
}

// EmbeddingResponse is a struct for embedding response
type EmbeddingResponse struct {
	Text      string    `json:"text"`
	Embedding []float32 `json:"embedding"`
}
