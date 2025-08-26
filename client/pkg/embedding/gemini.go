package embedding

import (
	"context"
	"fmt"

	"github.com/wnjoon/vector-search-tester/pkg/model"
	"google.golang.org/genai"
)

type geminiEmbedder struct {
	model  string
	client *genai.Client
}

func NewGeminiEmbedder(client *genai.Client, model string) Embedder {
	return &geminiEmbedder{
		client: client,
		model:  model,
	}
}

func (g *geminiEmbedder) Embed(ctx context.Context, req model.EmbeddingRequest) (*model.EmbeddingResponse, error) {
	contents := []*genai.Content{
		genai.NewContentFromText(req.Text, genai.RoleUser),
	}
	resp, err := g.client.Models.EmbedContent(
		ctx,
		g.model,
		contents,
		nil,
		// &genai.EmbedContentConfig{
		// 	TaskType: req.TaskType,
		// },
	)
	if err != nil {
		return nil, fmt.Errorf("failed to embed: %w", err)
	}

	if len(resp.Embeddings) == 0 {
		return nil, fmt.Errorf("no embeddings in response")
	}

	return &model.EmbeddingResponse{
		Text:      req.Text,
		Embedding: resp.Embeddings[0].Values,
	}, nil
}
