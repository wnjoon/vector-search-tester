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
	assert.NoError(t, err)

	geminiEmbedder := NewGeminiEmbedder(client, "models/embedding-001")

	t.Run("Embed - ShortText", func(t *testing.T) {
		t.Run("today's weather is wonderful", func(t *testing.T) {
			req := model.EmbeddingRequest{
				Text: "today's weather is wonderful",
			}
			resp, err := geminiEmbedder.Embed(ctx, req)
			assert.NoError(t, err)
			t.Logf("Vector Dimension: %d\n", len(resp.Embedding))
		})
	})

	t.Run("Embed - LongText - Eng", func(t *testing.T) {
		t.Run("GopherCon", func(t *testing.T) {
			req := model.EmbeddingRequest{
				Text: LongTextGopherConEng,
			}
			resp, err := geminiEmbedder.Embed(ctx, req)
			assert.NoError(t, err)
			t.Logf("Vector Dimension: %d\n", len(resp.Embedding))
		})
	})

	t.Run("Embed - LongText - Ko", func(t *testing.T) {
		t.Run("GopherCon", func(t *testing.T) {
			req := model.EmbeddingRequest{
				Text: LongTextGopherConKor,
			}
			resp, err := geminiEmbedder.Embed(ctx, req)
			assert.NoError(t, err)
			t.Logf("Vector Dimension: %d\n", len(resp.Embedding))
		})
	})
}

func LoadAPIKey() string {
	if err := godotenv.Load("../../.env"); err != nil {
		return ""
	}
	return os.Getenv("GEMINI_API_KEY")
}
