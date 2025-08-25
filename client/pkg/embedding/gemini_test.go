package embedding

import (
	"context"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/wnjoon/vector-search-tester/pkg/model"
	"google.golang.org/genai"
)

func TestGeminiEmbedder(t *testing.T) {
	ctx := context.Background()
	client, err := genai.NewClient(
		ctx,
		&genai.ClientConfig{
			APIKey: LoadAPIKey(),
		},
	)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	geminiEmbedder := NewGeminiEmbedder(client)
	t.Run("Embed", func(t *testing.T) {
		req := model.EmbeddingRequest{
			Text: "today's weather is wonderful",
		}
		resp, err := geminiEmbedder.Embed(ctx, req)
		assert.NoError(t, err)
		t.Logf("Vector Dimension: %d\n", len(resp.Embedding))
	})

}

func LoadAPIKey() string {
	if err := godotenv.Load("../../.env"); err != nil {
		return ""
	}
	return os.Getenv("GEMINI_API_KEY")
}
