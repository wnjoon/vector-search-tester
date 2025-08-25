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
	t.Run("Embed", func(t *testing.T) {
		req := model.EmbeddingRequest{
			Text: "today's weather is wonderful",
		}
		resp, err := client.Embed(context.Background(), req)
		assert.NoError(t, err)
		t.Logf("Vector Dimension: %d\n", len(resp.Embedding))
	})
}
