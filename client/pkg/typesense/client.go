package typesense

import (
	"context"
	"fmt"

	ts "github.com/typesense/typesense-go/typesense"
	tsapi "github.com/typesense/typesense-go/typesense/api"
	"github.com/wnjoon/vector-search-tester/pkg/embedding"
)

type Client struct {
	client   *ts.Client
	embedder embedding.Embedder
}

func New(serverUrl, apiKey string, embedder embedding.Embedder) *Client {
	client := ts.NewClient(
		ts.WithServer(serverUrl),
		ts.WithAPIKey(apiKey),
	)
	return &Client{
		client:   client,
		embedder: embedder,
	}
}

func (c *Client) CreateCollection(ctx context.Context, schema *tsapi.CollectionSchema) error {
	if schema == nil {
		return fmt.Errorf("collection schema is nil")
	}

	if schema.Name == "" {
		return fmt.Errorf("collection name is empty")
	}

	_, err := c.client.Collections().Create(
		ctx,
		schema,
	)
	if err != nil {
		return fmt.Errorf("failed to create collection %s: %w", schema.Name, err)
	}
	return nil
}
