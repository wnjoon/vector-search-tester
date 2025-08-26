package typesense

// import (
// 	"context"
// 	"os"
// 	"strconv"
// 	"strings"
// 	"testing"
// 	"time"

// 	"github.com/joho/godotenv"
// 	"github.com/stretchr/testify/assert"
// 	tsapi "github.com/typesense/typesense-go/typesense/api"
// 	"github.com/wnjoon/vector-search-tester/pkg/embedding"
// 	"github.com/wnjoon/vector-search-tester/pkg/model"
// 	"google.golang.org/genai"
// )

// const sentenceBertEmbedderUrl string = "http://localhost:6600"
// const countriesDescriptionCollName string = "countries_description"

// func TestNewTypeSenseClient(t *testing.T) {
// 	ctx := context.Background()
// 	embedder := embedding.NewSentenceBertEmbedder(sentenceBertEmbedderUrl)
// 	ts := New("http://localhost:8108", "xyz", embedder)
// 	assert.NotNil(t, ts)

// 	t.Run("HealthCheck", func(t *testing.T) {
// 		ok, err := ts.client.Health(ctx, 1*time.Second)
// 		assert.True(t, ok)
// 		assert.NoError(t, err)
// 	})
// }

// func TestTypeSenseClientWithSentenceBertEmbedder(t *testing.T) {
// 	ctx := context.Background()
// 	embedder := embedding.NewSentenceBertEmbedder("http://localhost:6600")
// 	ts := New("http://localhost:8108", "xyz", embedder)

// 	t.Run("test english text", func(t *testing.T) {
// 		name := "developer_conferences"
// 		numDim := 384

// 		t.Run("delete collection", func(t *testing.T) {
// 			// ignore error
// 			ts.client.Collection(name).Delete(ctx)
// 		})

// 		t.Run("create collection", func(t *testing.T) {
// 			schema := &tsapi.CollectionSchema{
// 				Name: name,
// 				Fields: []tsapi.Field{
// 					{
// 						Name: "id",
// 						Type: "string",
// 					},
// 					{
// 						Name: "lang",
// 						Type: "string",
// 					},
// 					{
// 						Name: "content",
// 						Type: "string",
// 					},
// 					{
// 						Name:   "embedding",
// 						Type:   "float[]",
// 						NumDim: &numDim,
// 					},
// 				},
// 			}
// 			err := ts.CreateCollection(ctx, schema)
// 			assert.NoError(t, err)
// 		})

// 		t.Run("check collection", func(t *testing.T) {
// 			col, err := ts.client.Collection(name).Retrieve(ctx)
// 			assert.NoError(t, err)
// 			assert.NotNil(t, col)
// 			assert.Equal(t, name, col.Name)
// 		})

// 		// documents to index
// 		docsToIndex := []Document{
// 			{ID: "gophercon-en", Lang: "en", Content: GopherConTextEn},
// 			{ID: "pycon-en", Lang: "en", Content: PyconTextEn},
// 			{ID: "nodecongress-en", Lang: "en", Content: NodeCongressTextEn},
// 			{ID: "devcon-en", Lang: "en", Content: DevconTextEn},
// 		}

// 		t.Run("add documents", func(t *testing.T) {
// 			for _, doc := range docsToIndex {
// 				t.Run("add document: "+doc.ID, func(t *testing.T) {
// 					resp, err := ts.embedder.Embed(ctx, model.EmbeddingRequest{
// 						Text:     doc.Content,
// 						Language: doc.Lang,
// 					})
// 					assert.NoError(t, err)
// 					assert.NotNil(t, resp)
// 					assert.Equal(t, len(resp.Embedding), numDim)
// 					assert.Equal(t, doc.Content, resp.Text)

// 					typesenseDoc := map[string]interface{}{
// 						"id":        doc.ID,
// 						"lang":      doc.Lang,
// 						"content":   doc.Content,
// 						"embedding": resp.Embedding,
// 					}

// 					_, err = ts.client.Collection(name).Documents().Create(ctx, typesenseDoc)
// 					assert.NoError(t, err)
// 				})
// 			}
// 		})

// 		t.Run("search documents", func(t *testing.T) {
// 			quries := []string{
// 				// "a conference for the Go community",
// 				// "what conference's mascot is gopher?",
// 				// "where is the best place to get the newest tech news for javascript?",
// 				// "where is the highest proportion to meet vitalik buterin?",
// 				// "please recommend the best conference for python developer",
// 				"golang", "gopher", "javascript", "ethereum", "python",
// 			}

// 			results := make([]string, 0)

// 			for _, q := range quries {
// 				t.Run("search document: "+q, func(t *testing.T) {
// 					vector, err := ts.embedder.Embed(ctx, model.EmbeddingRequest{
// 						Text:     q,
// 						Language: "en",
// 					})
// 					assert.NoError(t, err)
// 					assert.NotNil(t, vector)
// 					assert.Equal(t, len(vector.Embedding), numDim)
// 					assert.Equal(t, q, vector.Text)

// 					vectorString := "embedding:([" + vectorToString(vector.Embedding) + "])"
// 					includeFields := []string{"id"}
// 					includeFieldsStr := strings.Join(includeFields, ",")

// 					queryBy := "content"

// 					// when using search, an error occurs
// 					// status: 400 response: {"message":"Query string exceeds max allowed length of 4000.
// 					// Use the /multi_search end-point for larger payloads."
// 					t.Run("use multi-search", func(t *testing.T) {
// 						res, err := ts.client.MultiSearch.Perform(
// 							ctx,
// 							&tsapi.MultiSearchParams{
// 								Q:             &q,
// 								QueryBy:       &queryBy,
// 								IncludeFields: &includeFieldsStr,
// 								VectorQuery:   &vectorString,
// 							},
// 							tsapi.MultiSearchSearchesParameter{
// 								Searches: []tsapi.MultiSearchCollectionParameters{
// 									{
// 										Collection: name,
// 									},
// 								},
// 							},
// 						)
// 						assert.NoError(t, err)
// 						assert.NotNil(t, res)

// 						assert.GreaterOrEqual(t, len(res.Results), 1)

// 						hits := *(res.Results[0].Hits)
// 						docs := *(hits[0].Document)
// 						results = append(results, docs["id"].(string))
// 					})
// 				})
// 			}

// 			assert.Equal(t, len(results), len(quries))
// 			for i, result := range results {
// 				// assert.Equal(t, result, docsToIndex[i].ID)
// 				t.Log(i, ":", result)
// 			}
// 		})
// 	})
// }

// // embedding model for gemini
// type embeddingModel struct {
// 	Model  string `json:"model"`
// 	NumDim int    `json:"num_dim"`
// }

// var embeddingModels = []embeddingModel{
// 	{
// 		Model:  "models/embedding-001",
// 		NumDim: 768,
// 	},
// 	{
// 		Model:  "gemini-embedding-001",
// 		NumDim: 3072,
// 	},
// }

// func TestInsertNewDocuments(t *testing.T) {
// 	ctx := context.Background()
// 	gc, err := genai.NewClient(
// 		ctx,
// 		&genai.ClientConfig{
// 			APIKey: LoadAPIKey(),
// 		},
// 	)
// 	assert.NoError(t, err)
// 	embeddingModel := embeddingModels[0]
// 	embedder := embedding.NewGeminiEmbedder(gc, embeddingModel.Model)
// 	ts := New("http://localhost:8108", "xyz", embedder)

// 	collectionName := "short_text_collection"

// 	t.Run("create new collection", func(t *testing.T) {
// 		ts.client.Collection(collectionName).Delete(ctx) // ignore error
// 		schema := &tsapi.CollectionSchema{
// 			Name: collectionName,
// 			Fields: []tsapi.Field{
// 				{
// 					Name: "id",
// 					Type: "string",
// 				},
// 				{
// 					Name: "content",
// 					Type: "string",
// 				},
// 				{
// 					Name:   "embedding",
// 					Type:   "float[]",
// 					NumDim: &embeddingModel.NumDim,
// 				},
// 			},
// 		}
// 		assert.NoError(t, ts.CreateCollection(ctx, schema))
// 	})

// 	t.Run("add documents", func(t *testing.T) {
// 		sampleShortTexts := []string{
// 			"paris is the capital city of france",
// 			"tokyo is the capital city of japan",
// 			"beijing is the capital city of china",
// 			"seoul is the capital city of south korea",
// 			"busan is the city of south korea",
// 			"rome is the capital city of italy",
// 			"moscow is the capital city of russia",
// 		}

// 		for _, text := range sampleShortTexts {
// 			time.Sleep(5 * time.Second)
// 			t.Run("add document: "+text, func(t *testing.T) {
// 				vector, err := ts.embedder.Embed(ctx, model.EmbeddingRequest{
// 					Text:     text,
// 					Language: "en",
// 				})
// 				assert.NoError(t, err)
// 				assert.NotNil(t, vector)
// 				assert.Equal(t, len(vector.Embedding), embeddingModel.NumDim)
// 				t.Log("Vector: ", vector.Embedding)
// 				assert.Equal(t, text, vector.Text)

// 				typesenseDoc := map[string]interface{}{
// 					"id":        text,
// 					"content":   text,
// 					"embedding": vector.Embedding,
// 				}

// 				_, err = ts.client.Collection(collectionName).Documents().Create(ctx, typesenseDoc)
// 				assert.NoError(t, err)
// 			})
// 		}
// 	})
// }

// func TestInsertNewDocumentsIntoTypeSenseCollection(t *testing.T) {
// 	ctx := context.Background()
// 	t.Run("sentence-bert", func(t *testing.T) {
// 		embedder := embedding.NewSentenceBertEmbedder(sentenceBertEmbedderUrl)
// 		ts := New("http://localhost:8108", "xyz", embedder)

// 		t.Run("countries descriptions", func(t *testing.T) {

// 		})
// 	})
// 	gc, err := genai.NewClient(
// 		ctx,
// 		&genai.ClientConfig{
// 			APIKey: LoadAPIKey(),
// 		},
// 	)
// 	assert.NoError(t, err)
// 	embeddingModel := embeddingModels[0]
// 	embedder := embedding.NewGeminiEmbedder(gc, embeddingModel.Model)
// 	ts := New("http://localhost:8108", "xyz", embedder)

// 	collectionName := "short_text_collection"

// 	t.Run("create new collection", func(t *testing.T) {
// 		ts.client.Collection(collectionName).Delete(ctx) // ignore error
// 		schema := &tsapi.CollectionSchema{
// 			Name: collectionName,
// 			Fields: []tsapi.Field{
// 				{
// 					Name: "id",
// 					Type: "string",
// 				},
// 				{
// 					Name: "content",
// 					Type: "string",
// 				},
// 				{
// 					Name:   "embedding",
// 					Type:   "float[]",
// 					NumDim: &embeddingModel.NumDim,
// 				},
// 			},
// 		}
// 		assert.NoError(t, ts.CreateCollection(ctx, schema))
// 	})

// 	t.Run("add documents", func(t *testing.T) {
// 		for _, text := range loadCountriesDescriptions() {
// 			time.Sleep(5 * time.Second)
// 			t.Run("add document: "+text, func(t *testing.T) {
// 				vector, err := ts.embedder.Embed(ctx, model.EmbeddingRequest{
// 					Text:     text,
// 					Language: "en",
// 				})
// 				assert.NoError(t, err)
// 				assert.NotNil(t, vector)
// 				assert.Equal(t, len(vector.Embedding), embeddingModel.NumDim)
// 				t.Log("Vector: ", vector.Embedding)
// 				assert.Equal(t, text, vector.Text)

// 				typesenseDoc := map[string]interface{}{
// 					"id":        text,
// 					"content":   text,
// 					"embedding": vector.Embedding,
// 				}

// 				_, err = ts.client.Collection(collectionName).Documents().Create(ctx, typesenseDoc)
// 				assert.NoError(t, err)
// 			})
// 		}
// 	})

// }

// func TestTypeSenseClientWithGeminiEmbedder(t *testing.T) {
// 	ctx := context.Background()
// 	t.Run("test load api key from .env", func(t *testing.T) {
// 		apiKey := LoadAPIKey()
// 		assert.NotEmpty(t, apiKey)
// 	})

// 	genCli, err := genai.NewClient(
// 		ctx,
// 		&genai.ClientConfig{
// 			APIKey: LoadAPIKey(),
// 		},
// 	)
// 	assert.NoError(t, err)

// 	embeddingModel := struct {
// 		Model  string `json:"model"`
// 		NumDim int    `json:"num_dim"`
// 	}{
// 		Model:  "models/embedding-001",
// 		NumDim: 768,
// 	}

// 	// embeddingModel := struct {
// 	// 	Model  string `json:"model"`
// 	// 	NumDim int    `json:"num_dim"`
// 	// }{
// 	// 	Model:  "gemini-embedding-001",
// 	// 	NumDim: 3072,
// 	// }

// 	embedder := embedding.NewGeminiEmbedder(genCli, embeddingModel.Model)
// 	ts := New("http://localhost:8108", "xyz", embedder)

// 	t.Run("test english text", func(t *testing.T) {
// 		name := "developer_conferences"
// 		numDim := embeddingModel.NumDim

// 		t.Run("delete collection", func(t *testing.T) {
// 			// ignore error
// 			ts.client.Collection(name).Delete(ctx)
// 		})

// 		t.Run("create collection", func(t *testing.T) {
// 			schema := &tsapi.CollectionSchema{
// 				Name: name,
// 				Fields: []tsapi.Field{
// 					{
// 						Name: "id",
// 						Type: "string",
// 					},
// 					{
// 						Name: "lang",
// 						Type: "string",
// 					},
// 					{
// 						Name: "content",
// 						Type: "string",
// 					},
// 					{
// 						Name:   "embedding",
// 						Type:   "float[]",
// 						NumDim: &numDim,
// 					},
// 				},
// 			}
// 			err := ts.CreateCollection(ctx, schema)
// 			assert.NoError(t, err)
// 		})

// 		t.Run("check collection", func(t *testing.T) {
// 			col, err := ts.client.Collection(name).Retrieve(ctx)
// 			assert.NoError(t, err)
// 			assert.NotNil(t, col)
// 			assert.Equal(t, name, col.Name)
// 		})

// 		// documents to index
// 		docsToIndex := []Document{
// 			{ID: "gophercon-en", Lang: "en", Content: GopherConTextEn},
// 			{ID: "pycon-en", Lang: "en", Content: PyconTextEn},
// 			{ID: "nodecongress-en", Lang: "en", Content: NodeCongressTextEn},
// 			{ID: "devcon-en", Lang: "en", Content: DevconTextEn},
// 		}

// 		t.Run("add documents", func(t *testing.T) {
// 			for _, doc := range docsToIndex {
// 				t.Run("add document: "+doc.ID, func(t *testing.T) {
// 					resp, err := ts.embedder.Embed(ctx, model.EmbeddingRequest{
// 						Text:     doc.Content,
// 						Language: doc.Lang,
// 						// TaskType: "RETRIEVAL_DOCUMENT",
// 					})
// 					assert.NoError(t, err)
// 					assert.NotNil(t, resp)
// 					assert.Equal(t, len(resp.Embedding), numDim)
// 					assert.Equal(t, doc.Content, resp.Text)

// 					typesenseDoc := map[string]interface{}{
// 						"id":        doc.ID,
// 						"lang":      doc.Lang,
// 						"content":   doc.Content,
// 						"embedding": resp.Embedding,
// 					}

// 					_, err = ts.client.Collection(name).Documents().Create(ctx, typesenseDoc)
// 					assert.NoError(t, err)
// 				})
// 			}
// 		})

// 		// quries := []string{
// 		// 	"a conference for the Go community",
// 		// 	"what conference's mascot is gopher?",
// 		// 	"where is the best place to get the newest tech news for javascript?",
// 		// 	"where is the highest proportion to meet vitalik buterin?",
// 		// 	"please recommend the best conference for python developer",
// 		// }

// 		quries := []string{
// 			"golang", "gopher", "javascript", "ethereum", "python",
// 		}

// 		t.Run("search documents (hybrid)", func(t *testing.T) {
// 			results := make([]string, 0)
// 			for _, q := range quries {
// 				t.Run("search document: "+q, func(t *testing.T) {
// 					vector, err := ts.embedder.Embed(ctx, model.EmbeddingRequest{
// 						Text:     q,
// 						Language: "en",
// 						// TaskType: "RETRIEVAL_QUERY",
// 					})
// 					assert.NoError(t, err)
// 					assert.NotNil(t, vector)
// 					assert.Equal(t, len(vector.Embedding), numDim)
// 					assert.Equal(t, q, vector.Text)

// 					vectorString := "embedding:([" + vectorToString(vector.Embedding) + "])"
// 					includeFields := []string{"id"}
// 					includeFieldsStr := strings.Join(includeFields, ",")

// 					queryBy := "content"

// 					// when using search, an error occurs
// 					// status: 400 response: {"message":"Query string exceeds max allowed length of 4000.
// 					// Use the /multi_search end-point for larger payloads."
// 					t.Run("use multi-search", func(t *testing.T) {
// 						res, err := ts.client.MultiSearch.Perform(
// 							ctx,
// 							&tsapi.MultiSearchParams{
// 								Q:             &q,
// 								QueryBy:       &queryBy,
// 								IncludeFields: &includeFieldsStr,
// 								VectorQuery:   &vectorString,
// 							},
// 							tsapi.MultiSearchSearchesParameter{
// 								Searches: []tsapi.MultiSearchCollectionParameters{
// 									{
// 										Collection: name,
// 									},
// 								},
// 							},
// 						)
// 						assert.NoError(t, err)
// 						assert.NotNil(t, res)

// 						assert.GreaterOrEqual(t, len(res.Results), 1)

// 						hits := *(res.Results[0].Hits)
// 						docs := *(hits[0].Document)
// 						results = append(results, docs["id"].(string))
// 					})
// 				})
// 			}

// 			assert.Equal(t, len(results), len(quries))
// 			for i, result := range results {
// 				// assert.Equal(t, result, docsToIndex[i].ID)
// 				t.Log(i, ":", result)
// 			}
// 		})

// 		// t.Run("search documents (pure vector)", func(t *testing.T) {
// 		// 	// in pure vector search, we don't use query string
// 		// 	// it uses only performance of model itself
// 		// 	time.Sleep(1 * time.Second)
// 		// 	results := make([]string, 0)
// 		// 	for _, q := range quries {
// 		// 		t.Run("search document: "+q, func(t *testing.T) {
// 		// 			vector, err := ts.embedder.Embed(ctx, model.EmbeddingRequest{
// 		// 				Text:     q,
// 		// 				Language: "en",
// 		// 				TaskType: "RETRIEVAL_QUERY",
// 		// 			})
// 		// 			assert.NoError(t, err)
// 		// 			assert.NotNil(t, vector)
// 		// 			assert.Equal(t, len(vector.Embedding), numDim)
// 		// 			assert.Equal(t, q, vector.Text)

// 		// 			vectorString := "embedding:([" + vectorToString(vector.Embedding) + "])"
// 		// 			includeFields := []string{"id"}
// 		// 			includeFieldsStr := strings.Join(includeFields, ",")

// 		// 			// queryBy := "content"

// 		// 			// when using search, an error occurs
// 		// 			// status: 400 response: {"message":"Query string exceeds max allowed length of 4000.
// 		// 			// Use the /multi_search end-point for larger payloads."
// 		// 			t.Run("use multi-search", func(t *testing.T) {
// 		// 				res, err := ts.client.MultiSearch.Perform(
// 		// 					ctx,
// 		// 					&tsapi.MultiSearchParams{
// 		// 						// Q:             &q,
// 		// 						// QueryBy:       &queryBy,
// 		// 						IncludeFields: &includeFieldsStr,
// 		// 						VectorQuery:   &vectorString,
// 		// 					},
// 		// 					tsapi.MultiSearchSearchesParameter{
// 		// 						Searches: []tsapi.MultiSearchCollectionParameters{
// 		// 							{
// 		// 								Collection: name,
// 		// 							},
// 		// 						},
// 		// 					},
// 		// 				)
// 		// 				assert.NoError(t, err)
// 		// 				assert.NotNil(t, res)

// 		// 				assert.GreaterOrEqual(t, len(res.Results), 1)

// 		// 				firstSearchResult := res.Results[0]

// 		// 				if firstSearchResult.Found != nil {
// 		// 					foundCount := *firstSearchResult.Found
// 		// 					t.Log("found count in query:", q, ":", foundCount)

// 		// 					if foundCount > 0 {
// 		// 						hits := *(firstSearchResult.Hits)
// 		// 						assert.NotEmpty(t, hits)
// 		// 						docs := *(hits[0].Document)
// 		// 						results = append(results, docs["id"].(string))
// 		// 					}
// 		// 				}
// 		// 			})
// 		// 		})
// 		// 	}

// 		// 	// assert.Equal(t, len(results), len(quries))
// 		// 	t.Log("results count:", len(results))
// 		// 	for i, result := range results {
// 		// 		// assert.Equal(t, result, docsToIndex[i].ID)
// 		// 		t.Log(i, ":", result)
// 		// 	}
// 		// })
// 	})
// }

// func vectorToString(v []float32) string {
// 	var sb strings.Builder
// 	for i, f := range v {
// 		if i > 0 {
// 			sb.WriteString(",")
// 		}
// 		sb.WriteString(strconv.FormatFloat(float64(f), 'f', -1, 64))
// 	}
// 	return sb.String()
// }

// func LoadAPIKey() string {
// 	if err := godotenv.Load("../../.env"); err != nil {
// 		return ""
// 	}
// 	return os.Getenv("GEMINI_API_KEY")
// }

// func loadCountriesDescriptions() []string {
// 	return []string{
// 		"paris is the capital city of france",
// 		"tokyo is the capital city of japan",
// 		"beijing is the capital city of china",
// 		"seoul is the capital city of south korea",
// 		"busan is the city of south korea",
// 		"rome is the capital city of italy",
// 		"moscow is the capital city of russia",
// 	}
// }

// func loadDeveloperConferenceDescriptions() []string {
// 	return []string{
// 		GopherConTextEn,
// 		PyconTextEn,
// 		NodeCongressTextEn,
// 		DevconTextEn,
// 	}
// }

// const GopherConTextEn string = `GopherCon is the premier annual global conference for the Go (Golang) programming language, named after its mascot, the gopher. It provides a platform for Go developers to learn, network, and share knowledge about the language's features, tools, and best practices. Organized by a community-driven group, GopherCon fosters an inclusive and supportive Go community by offering state-of-the-art talks and creating opportunities for growth for developers of all experience levels.
// Key Aspects of GopherCon
// Focus on Go:
// The conference is dedicated to the Go programming language, bringing together enthusiasts, contributors, and professionals.
// Community-Driven:
// GopherCon is organized by a community of developers who aim to promote the use of Go and foster an inclusive community.
// Educational Content:
// Attendees can expect talks, presentations, and discussions on Go's latest features, concurrency models, and applications in areas like AI and distributed systems.
// Networking and Collaboration:
// It offers a valuable space for developers to connect with peers, share experiences, and collaborate on new ideas.
// Global Reach:
// While the original conference started in 2014, GopherCon has expanded into a global series of events, with conferences hosted in various countries around the world.
// The "Gopher" Mascot
// The name GopherCon comes from the Go language's official mascot, the gopher. The conference's branding and community are strongly tied to this friendly character, which is a common symbol for the Go programming language community. `

// const PyconTextEn string = `PyCon is the largest annual international conference for the Python programming language community, organized by the Python Software Foundation (PSF). It features talks, tutorials, and other community events to advance Python and foster its global community, keeping registration costs low due to its non-profit nature.
// What is PyCon?
// A Community Gathering:
// PyCon is the premier annual event for the Python programming community, bringing together users and developers of the open-source Python language.
// Organized by the PSF:
// The conference is produced and underwritten by the Python Software Foundation, a non-profit organization dedicated to promoting Python.
// Global and Local:
// While originating in North America, PyCon events are held in many countries worldwide, fostering a global Python community.
// What Happens at PyCon?
// Talks and Tutorials:
// The event includes numerous talks, tutorials, demonstrations, and presentations on a wide variety of Python-related topics.
// Community Events:
// PyCon hosts community events, workshops, and gatherings where attendees can network and share ideas.
// Professional Development:
// It offers opportunities for learning new skills, discovering cutting-edge ideas, and connecting with industry professionals.
// Why is PyCon Important?
// Fosters Community:
// PyCon's main goal is to grow and support the international community of Python programmers.
// Accessibility:
// Being a non-profit event, the PSF keeps registration costs low to ensure broader accessibility for attendees.
// Code of Conduct:
// PyCon is known for having and adhering to a code of conduct, promoting an inclusive and enjoyable environment for everyone.`

// const NodeCongressTextEn string = `NodeCongress is widely recognized as a significant conference focusing on Node.js and JavaScript backends. It covers modern JavaScript backend runtimes, best practices for scaling and maintaining applications, and features talks from prominent speakers in the JavaScript and Node.js communities.
// While NodeCongress is a major event specifically for Node.js, other large JavaScript conferences, such as JSNation and JSWORLD Conference, also dedicate substantial attention to Node.js within their broader JavaScript scope. These conferences attract large audiences and feature content relevant to Node.js developers.`

// const DevconTextEn string = `Devcon is a premier annual conference by the Ethereum Foundation for builders, thinkers, and enthusiasts in the Ethereum ecosystem, aiming to foster education, collaboration, and innovation in decentralized technology. The event serves as an intensive introduction for newcomers and a global gathering for the established Ethereum community to learn about new tools, protocols, and the culture of decentralized systems. Devcon's programming features a range of technical to human-focused content, including talks, workshops, and peer-to-peer interactions, all delivered in a sponsor-free environment to ensure unbiased information. Etheruem is created by vitalik buterin.
// Key Aspects of Devcon:
// Purpose:
// To empower the community to build and use decentralized systems, pushing the boundaries of decentralized protocols, tools, and culture.
// Audience:
// Developers, thinkers, makers, and anyone interested in Ethereum's mission to make decentralized systems accessible to the world.
// Format:
// A four-day event that rotates globally, featuring talks, workshops, and other interactive sessions.
// Programming:
// Content covers a wide spectrum from deep technical discussions on the base protocol to broader topics on the application layer, with sessions tailored to different expertise levels.
// Community:
// It's described as a "global family reunion" for the Ethereum community, providing a unique opportunity to connect with peers and gain a deeper understanding of the ecosystem's culture.
// Structure:
// Devcon is the Ethereum Foundation's principal event, distinguished from Devconnect, which is a separate collaborative initiative of smaller, topic-focused community events.
// Philosophy:
// As a sponsor-free event, Devcon focuses on delivering unbiased overviews of Ethereum's state, values, and future direction, with speakers and talks selected based on merit rather than financial incentives. `
