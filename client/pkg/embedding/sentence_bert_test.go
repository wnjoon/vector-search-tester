package embedding

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wnjoon/vector-search-tester/pkg/model"
)

const URL string = "http://localhost:6600"

func TestSentenceBertEmbedder(t *testing.T) {
	client := NewSentenceBertEmbedder(URL)
	t.Run("Embed - ShortText", func(t *testing.T) {
		req := model.EmbeddingRequest{
			Language: "en",
			Text:     "today's weather is wonderful",
		}
		resp, err := client.Embed(context.Background(), req)
		assert.NoError(t, err)
		t.Logf("Vector Dimension: %d\n", len(resp.Embedding))
	})

	t.Run("Embed - LongText - Eng", func(t *testing.T) {
		req := model.EmbeddingRequest{
			Language: "en",
			Text:     LongTextGopherConEng,
		}
		resp, err := client.Embed(context.Background(), req)
		assert.NoError(t, err)
		t.Logf("Vector Dimension: %d\n", len(resp.Embedding))
	})

	t.Run("Embed - LongText - Ko", func(t *testing.T) {
		req := model.EmbeddingRequest{
			Language: "ko",
			Text:     LongTextGopherConKor,
		}
		resp, err := client.Embed(context.Background(), req)
		assert.NoError(t, err)
		t.Logf("Vector Dimension: %d\n", len(resp.Embedding))
	})
}
