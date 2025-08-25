package embedding

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/wnjoon/vector-search-tester/pkg/model"
)

type sentenceBertEmbedder struct {
	url string
}

func NewSentenceBertEmbedder(url string) Embedder {
	return &sentenceBertEmbedder{
		url: url,
	}
}

func (s *sentenceBertEmbedder) Embed(ctx context.Context, req model.EmbeddingRequest) (*model.EmbeddingResponse, error) {
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON: %w", err)
	}

	if req.Language != "eng" && req.Language != "ko" {
		return nil, fmt.Errorf("invalid language: %s", req.Language)
	}

	url := s.url + "/embed/" + req.Language

	resp, err := http.Post(
		url,
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to make HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API server returned non-200 status: %d %s", resp.StatusCode, string(body))
	}

	var embeddingResp model.EmbeddingResponse
	if err := json.NewDecoder(resp.Body).Decode(&embeddingResp); err != nil {
		return nil, fmt.Errorf("failed to decode JSON response: %w", err)
	}

	return &embeddingResp, nil
}
