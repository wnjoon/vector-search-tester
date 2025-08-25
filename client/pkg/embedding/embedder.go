package embedding

import (
	"context"

	"github.com/wnjoon/vector-search-tester/pkg/model"
)

type Embedder interface {
	Embed(ctx context.Context, req model.EmbeddingRequest) (*model.EmbeddingResponse, error)
}
