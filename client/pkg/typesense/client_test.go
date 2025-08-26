package typesense

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/typesense/typesense-go/typesense"
	tsapi "github.com/typesense/typesense-go/typesense/api"
	"github.com/typesense/typesense-go/typesense/api/pointer"
	"github.com/wnjoon/vector-search-tester/pkg/embedding"
	"github.com/wnjoon/vector-search-tester/pkg/model"
	"google.golang.org/genai"
)

const typesenseApiKey string = "xyz"
const typesenseUrl string = "http://localhost:8108"
const sentenceBertEmbedderUrl string = "http://localhost:6600"

const (
	collCountriesDescriptionGeminiModelsEmbedding001 string = "countries_description_gemini_models_embedding_001"
	collCountriesDescriptionGeminiGeminiEmbedding001 string = "countries_description_gemini_gemini_embedding_001"
	collCountriesDescriptionSentenceBert             string = "countries_description_sentence_bert"
	collDeveloperConferenceGeminiModelsEmbedding001  string = "developer_conferences_gemini_models_embedding_001"
	collDeveloperConferenceGeminiGeminiEmbedding001  string = "developer_conferences_gemini_gemini_embedding_001"
	collDeveloperConferenceSentenceBert              string = "developer_conferences_sentence_bert"
)

// embedding model for gemini
type embeddingModel struct {
	Model  string `json:"model"`
	NumDim int    `json:"num_dim"`
}

var embeddingModelMap = map[string]embeddingModel{
	"models/embedding-001": {
		Model:  "models/embedding-001",
		NumDim: 768,
	},
	"gemini-embedding-001": {
		Model:  "gemini-embedding-001",
		NumDim: 3072,
	},
	"sentence-bert": {
		Model:  "sentence-bert",
		NumDim: 384,
	},
}

// create typesence schema by collection name and number of dimensions
func schema(collName string, numDim int) *tsapi.CollectionSchema {
	return &tsapi.CollectionSchema{
		Name: collName,
		Fields: []tsapi.Field{
			{Name: "id", Type: "string"},
			{Name: "content", Type: "string"},
			{Name: "embedding", Type: "float[]", NumDim: pointer.Int(numDim)},
		},
	}
}

func TestNewTypeSenseClient(t *testing.T) {
	ctx := context.Background()
	embedder := embedding.NewSentenceBertEmbedder(sentenceBertEmbedderUrl)
	ts := New("http://localhost:8108", "xyz", embedder)
	assert.NotNil(t, ts)

	t.Run("HealthCheck", func(t *testing.T) {
		ok, err := ts.client.Health(ctx, 1*time.Second)
		assert.True(t, ok)
		assert.NoError(t, err)
	})
}

/*
 * [TestCreateCollectionsForCountriesDescription]
 * create collections for countries description
 * this test creates collections for countries description with different embedding models
 * sentence-bert, models/embedding-001, gemini-embedding-001
 */
func TestCreateCollectionsForCountriesDescription(t *testing.T) {
	ctx := context.Background()
	ts := New(typesenseUrl, typesenseApiKey, nil)
	assert.NotNil(t, ts)

	t.Run("create sentence-bert embedding collection", func(t *testing.T) {
		em := embeddingModelMap["sentence-bert"]
		embedder := embedding.NewSentenceBertEmbedder(sentenceBertEmbedderUrl)
		ts.SetEmbedder(embedder)
		ts.DeleteCollection(ctx, collCountriesDescriptionSentenceBert)
		assert.NoError(t, ts.CreateCollection(ctx, schema(collCountriesDescriptionSentenceBert, em.NumDim)))
	})

	geminiClient, _ := genai.NewClient(ctx, &genai.ClientConfig{APIKey: LoadAPIKey()})

	t.Run("create models/embedding-001 embedding collection", func(t *testing.T) {
		em := embeddingModelMap["models/embedding-001"]
		embedder := embedding.NewGeminiEmbedder(geminiClient, em.Model)
		ts.SetEmbedder(embedder)
		ts.DeleteCollection(ctx, collCountriesDescriptionGeminiModelsEmbedding001)
		assert.NoError(t, ts.CreateCollection(ctx, schema(collCountriesDescriptionGeminiModelsEmbedding001, em.NumDim)))
	})

	t.Run("create gemini-embedding-001 embedding collection", func(t *testing.T) {
		em := embeddingModelMap["gemini-embedding-001"]
		embedder := embedding.NewGeminiEmbedder(geminiClient, em.Model)
		ts.SetEmbedder(embedder)
		ts.DeleteCollection(ctx, collCountriesDescriptionGeminiGeminiEmbedding001)
		assert.NoError(t, ts.CreateCollection(ctx, schema(collCountriesDescriptionGeminiGeminiEmbedding001, em.NumDim)))
	})
}

/*
 * [TestCreateCollectionsForDeveloperConferences]
 * create collections for developer conferences
 * this test creates collections for developer conferences with different embedding models
 * sentence-bert, models/embedding-001, gemini-embedding-001
 */
func TestCreateCollectionsForDeveloperConferences(t *testing.T) {
	ctx := context.Background()
	ts := New(typesenseUrl, typesenseApiKey, nil)
	assert.NotNil(t, ts)

	t.Run("create sentence-bert embedding collection", func(t *testing.T) {
		em := embeddingModelMap["sentence-bert"]
		embedder := embedding.NewSentenceBertEmbedder(sentenceBertEmbedderUrl)
		ts.SetEmbedder(embedder)
		ts.DeleteCollection(ctx, collDeveloperConferenceSentenceBert)
		assert.NoError(t, ts.CreateCollection(ctx, schema(collDeveloperConferenceSentenceBert, em.NumDim)))
	})

	geminiClient, _ := genai.NewClient(ctx, &genai.ClientConfig{APIKey: LoadAPIKey()})

	t.Run("create models/embedding-001 embedding collection", func(t *testing.T) {
		em := embeddingModelMap["models/embedding-001"]
		embedder := embedding.NewGeminiEmbedder(geminiClient, em.Model)
		ts.SetEmbedder(embedder)
		ts.DeleteCollection(ctx, collDeveloperConferenceGeminiModelsEmbedding001)
		assert.NoError(t, ts.CreateCollection(ctx, schema(collDeveloperConferenceGeminiModelsEmbedding001, em.NumDim)))
	})

	t.Run("create gemini-embedding-001 embedding collection", func(t *testing.T) {
		em := embeddingModelMap["gemini-embedding-001"]
		embedder := embedding.NewGeminiEmbedder(geminiClient, em.Model)
		ts.SetEmbedder(embedder)
		ts.DeleteCollection(ctx, collDeveloperConferenceGeminiGeminiEmbedding001)
		assert.NoError(t, ts.CreateCollection(ctx, schema(collDeveloperConferenceGeminiGeminiEmbedding001, em.NumDim)))
	})
}

// texts for countries description
// key is country name, value is description
// this document is inserted into typesense collections
var countriesDescriptionTexts map[string]string = map[string]string{
	"tokyo": "tokyo is the capital city of japan and there is a disney land. sushi is a popular food in japan.",
	"seoul": "seoul is the capital city of south korea and there is a lotte world. gangnam style is a popular k-pop song.",
	"rome":  "rome is the capital city of italy and there is a colosseum and Vatican City. the Vatican City is the smallest country in the world.",
}

/*
 * [TestAddCountriesDescriptionDocuments_SentenceBert]
 * add documents to countries description collections with sentence-bert embedding model
 */
func TestAddCountriesDescriptionDocuments_SentenceBert(t *testing.T) {
	ctx := context.Background()
	ts := New(typesenseUrl, typesenseApiKey, embedding.NewSentenceBertEmbedder(sentenceBertEmbedderUrl))
	assert.NotNil(t, ts)

	t.Run("add countries descriptions with sentence-bert", func(t *testing.T) {
		err := addDocumentToCollection(ctx, ts, embeddingModelMap["sentence-bert"], countriesDescriptionTexts, collCountriesDescriptionSentenceBert, "en")
		assert.NoError(t, err)
	})
}

/*
 * [TestAddCountriesDescriptionDocuments_ModelEmbedding001]
 * add documents to countries description collections with models/embedding-001 embedding model
 */
func TestAddCountriesDescriptionDocuments_ModelEmbedding001(t *testing.T) {
	ctx := context.Background()
	geminiClient, _ := genai.NewClient(ctx, &genai.ClientConfig{APIKey: LoadAPIKey()})
	em := embeddingModelMap["models/embedding-001"]
	ts := New(typesenseUrl, typesenseApiKey, embedding.NewGeminiEmbedder(geminiClient, em.Model))
	assert.NotNil(t, ts)

	t.Run("add countries descriptions with models/embedding-001", func(t *testing.T) {
		err := addDocumentToCollection(ctx, ts, em, countriesDescriptionTexts, collCountriesDescriptionGeminiModelsEmbedding001, "")
		assert.NoError(t, err)
	})
}

/*
 * [TestAddCountriesDescriptionDocuments_GeminiEmbedding001]
 * add documents to countries description collections with gemini-embedding-001 embedding model
 */
func TestAddCountriesDescriptionDocuments_GeminiEmbedding001(t *testing.T) {
	ctx := context.Background()
	geminiClient, _ := genai.NewClient(ctx, &genai.ClientConfig{APIKey: LoadAPIKey()})
	em := embeddingModelMap["gemini-embedding-001"]
	ts := New(typesenseUrl, typesenseApiKey, embedding.NewGeminiEmbedder(geminiClient, em.Model))
	assert.NotNil(t, ts)

	t.Run("add countries descriptions with models/embedding-001", func(t *testing.T) {
		err := addDocumentToCollection(ctx, ts, em, countriesDescriptionTexts, collCountriesDescriptionGeminiGeminiEmbedding001, "")
		assert.NoError(t, err)
	})
}

// texts for developer conferences
// key is conference name, value is description
// this document is inserted into typesense collections
var developerConferencesTexts map[string]string = map[string]string{
	"gophercon":    GopherConTextEn,
	"pycon":        PyconTextEn,
	"nodecongress": NodeCongressTextEn,
	"devcon":       DevconTextEn,
}

/*
 * [TestAddDeveloperConferencesDocuments_SentenceBert]
 * add documents to developer conferences collections with sentence-bert embedding model
 */
func TestAddDeveloperConferencesDocuments_SentenceBert(t *testing.T) {
	ctx := context.Background()
	ts := New(typesenseUrl, typesenseApiKey, embedding.NewSentenceBertEmbedder(sentenceBertEmbedderUrl))
	assert.NotNil(t, ts)

	t.Run("add developer conferences with sentence-bert", func(t *testing.T) {
		err := addDocumentToCollection(ctx, ts, embeddingModelMap["sentence-bert"], developerConferencesTexts, collDeveloperConferenceSentenceBert, "en")
		assert.NoError(t, err)
	})
}

/*
 * [TestAddDeveloperConferencesDocuments_ModelEmbedding001]
 * add documents to developer conferences collections with models/embedding-001 embedding model
 */
func TestAddDeveloperConferencesDocuments_ModelEmbedding001(t *testing.T) {
	ctx := context.Background()
	geminiClient, _ := genai.NewClient(ctx, &genai.ClientConfig{APIKey: LoadAPIKey()})
	em := embeddingModelMap["models/embedding-001"]
	ts := New(typesenseUrl, typesenseApiKey, embedding.NewGeminiEmbedder(geminiClient, em.Model))
	assert.NotNil(t, ts)

	t.Run("add developer conferences with models/embedding-001", func(t *testing.T) {
		err := addDocumentToCollection(ctx, ts, em, developerConferencesTexts, collDeveloperConferenceGeminiModelsEmbedding001, "")
		assert.NoError(t, err)
	})
}

/*
 * [TestAddDeveloperConferencesDocuments_GeminiEmbedding001]
 * add documents to developer conferences collections with gemini-embedding-001 embedding model
 */
func TestAddDeveloperConferencesDocuments_GeminiEmbedding001(t *testing.T) {
	ctx := context.Background()
	geminiClient, _ := genai.NewClient(ctx, &genai.ClientConfig{APIKey: LoadAPIKey()})
	em := embeddingModelMap["gemini-embedding-001"]
	ts := New(typesenseUrl, typesenseApiKey, embedding.NewGeminiEmbedder(geminiClient, em.Model))
	assert.NotNil(t, ts)

	t.Run("add developer conferences with gemini-embedding-001", func(t *testing.T) {
		err := addDocumentToCollection(ctx, ts, em, developerConferencesTexts, collDeveloperConferenceGeminiGeminiEmbedding001, "")
		assert.NoError(t, err)
	})
}

// texts for countries description
// we use this texts to search documents as query parameters
var vectorSearchRequestsCountriesDescription []string = []string{
	"seoul",
	"tokyo",
	"rome",
	"south korea",
	"japan",
	"italy",
	"lotte world",
	"sushi",
	"disney land",
	"colosseum",
	"smallest country",
}

// texts for countries description
// we use this texts as expected results
var vectorSearchExpectedCountriesDescription []string = []string{
	"seoul",
	"tokyo",
	"rome",
	"seoul",
	"tokyo",
	"rome",
	"seoul",
	"tokyo",
	"tokyo",
	"rome",
	"rome",
}

// create vectorized string from query
// result will be "embedding:([vector])"
func vectorString(ctx context.Context, embedder embedding.Embedder, q, lang string) (string, error) {
	vector, err := embedder.Embed(ctx, model.EmbeddingRequest{
		Text:     q,
		Language: lang,
	})
	if err != nil {
		return "", err
	}
	if vector == nil {
		return "", fmt.Errorf("vector is nil")
	}
	if vector.Text != q {
		return "", fmt.Errorf("vector text is not equal to query")
	}
	return "embedding:([" + vectorToString(vector.Embedding) + "])", nil
}

// perform multi search
// fields(QueryBy, IncludeFields) is static
func doMultiSearch(ctx context.Context, cli *typesense.Client, collName, q, vectorstr string) (string, error) {
	res, err := cli.MultiSearch.Perform(
		ctx,
		&tsapi.MultiSearchParams{
			Q:             &q,
			QueryBy:       pointer.String("content"),
			IncludeFields: pointer.String("id,content"),
			VectorQuery:   &vectorstr,
		},
		tsapi.MultiSearchSearchesParameter{
			Searches: []tsapi.MultiSearchCollectionParameters{
				{
					Collection: collName,
				},
			},
		},
	)
	if err != nil {
		return "", err
	}
	if res == nil {
		return "", fmt.Errorf("multi search result is nil")
	}
	if len(res.Results) == 0 {
		return "", fmt.Errorf("multi search result is empty")
	}
	hits := *(res.Results[0].Hits)
	docs := *(hits[0].Document)
	return docs["id"].(string), nil
}

// compare expected and actual results
func isMatched(expected, actual []string) bool {
	for i := range expected {
		if expected[i] != actual[i] {
			return false
		}
	}
	return true
}

/*
 * [TestVectorSearchCountriesDescription_SentenceBert]
 * vector search countries description with sentence-bert embedding model
 */
func TestVectorSearchCountriesDescription_SentenceBert(t *testing.T) {
	ctx := context.Background()
	ts := New(typesenseUrl, typesenseApiKey, embedding.NewSentenceBertEmbedder(sentenceBertEmbedderUrl))

	result := make([]string, 0)
	for _, q := range vectorSearchRequestsCountriesDescription {
		vectorstr, err := vectorString(ctx, ts.embedder, q, "en")
		assert.NoError(t, err)

		res, err := doMultiSearch(ctx, ts.client, collCountriesDescriptionSentenceBert, q, vectorstr)
		assert.NoError(t, err)
		result = append(result, res)
	}
	for i := range vectorSearchExpectedCountriesDescription {
		assert.Equal(t, result[i], vectorSearchExpectedCountriesDescription[i])
	}
	assert.True(t, isMatched(result, vectorSearchExpectedCountriesDescription))
}

/*
 * [TestVectorSearchCountriesDescription_ModelEmbedding001]
 * vector search countries description with models/embedding-001 embedding model
 */
func TestVectorSearchCountriesDescription_ModelEmbedding001(t *testing.T) {
	ctx := context.Background()
	geminiClient, _ := genai.NewClient(ctx, &genai.ClientConfig{APIKey: LoadAPIKey()})
	em := embeddingModelMap["models/embedding-001"]
	ts := New(typesenseUrl, typesenseApiKey, embedding.NewGeminiEmbedder(geminiClient, em.Model))

	result := make([]string, 0)
	for _, q := range vectorSearchRequestsCountriesDescription {
		vectorstr, err := vectorString(ctx, ts.embedder, q, "")
		assert.NoError(t, err)

		res, err := doMultiSearch(ctx, ts.client, collCountriesDescriptionGeminiModelsEmbedding001, q, vectorstr)
		assert.NoError(t, err)
		result = append(result, res)
	}
	for i := range vectorSearchExpectedCountriesDescription {
		assert.Equal(t, result[i], vectorSearchExpectedCountriesDescription[i])
	}
	assert.True(t, isMatched(result, vectorSearchExpectedCountriesDescription))
}

/*
 * [TestVectorSearchCountriesDescription_GeminiEmbedding001]
 * vector search countries description with gemini-embedding-001 embedding model
 */
func TestVectorSearchCountriesDescription_GeminiEmbedding001(t *testing.T) {
	ctx := context.Background()
	geminiClient, _ := genai.NewClient(ctx, &genai.ClientConfig{APIKey: LoadAPIKey()})
	em := embeddingModelMap["gemini-embedding-001"]
	ts := New(typesenseUrl, typesenseApiKey, embedding.NewGeminiEmbedder(geminiClient, em.Model))

	result := make([]string, 0)
	for _, q := range vectorSearchRequestsCountriesDescription {
		vectorstr, err := vectorString(ctx, ts.embedder, q, "")
		assert.NoError(t, err)

		res, err := doMultiSearch(ctx, ts.client, collCountriesDescriptionGeminiGeminiEmbedding001, q, vectorstr)
		assert.NoError(t, err)
		result = append(result, res)
	}
	for i := range vectorSearchExpectedCountriesDescription {
		assert.Equal(t, result[i], vectorSearchExpectedCountriesDescription[i])
	}
	assert.True(t, isMatched(result, vectorSearchExpectedCountriesDescription))
}

// texts for developer conferences
// we use this texts to search documents as query parameters
var vectorSearchRequestsDeveloperConferences []string = []string{
	"golang",
	"nodejs",
	"ethereum",
	"python",
}

// texts for developer conferences
// we use this texts as expected results
var vectorSearchExpectedDeveloperConferences []string = []string{
	"gophercon",
	"nodecongress",
	"devcon",
	"pycon",
}

/*
 * [TestVectorSearchDeveloperConferences_SentenceBert]
 * vector search developer conferences with sentence-bert embedding model
 */
func TestVectorSearchDeveloperConferences_SentenceBert(t *testing.T) {
	ctx := context.Background()
	ts := New(typesenseUrl, typesenseApiKey, embedding.NewSentenceBertEmbedder(sentenceBertEmbedderUrl))

	result := make([]string, 0)
	for _, q := range vectorSearchRequestsDeveloperConferences {
		vectorstr, err := vectorString(ctx, ts.embedder, q, "en")
		assert.NoError(t, err)

		res, err := doMultiSearch(ctx, ts.client, collDeveloperConferenceSentenceBert, q, vectorstr)
		assert.NoError(t, err)
		result = append(result, res)
	}
	for i := range vectorSearchExpectedDeveloperConferences {
		assert.Equal(t, result[i], vectorSearchExpectedDeveloperConferences[i])
	}
	assert.True(t, isMatched(result, vectorSearchExpectedDeveloperConferences))
}

/*
 * [TestVectorSearchDeveloperConferences_ModelEmbedding001]
 * vector search developer conferences with models/embedding-001 embedding model
 */
func TestVectorSearchDeveloperConferences_ModelEmbedding001(t *testing.T) {
	ctx := context.Background()
	geminiClient, _ := genai.NewClient(ctx, &genai.ClientConfig{APIKey: LoadAPIKey()})
	em := embeddingModelMap["models/embedding-001"]
	ts := New(typesenseUrl, typesenseApiKey, embedding.NewGeminiEmbedder(geminiClient, em.Model))

	result := make([]string, 0)
	for _, q := range vectorSearchRequestsDeveloperConferences {
		vectorstr, err := vectorString(ctx, ts.embedder, q, "")
		assert.NoError(t, err)

		res, err := doMultiSearch(ctx, ts.client, collDeveloperConferenceGeminiModelsEmbedding001, q, vectorstr)
		assert.NoError(t, err)
		result = append(result, res)
	}
	for i := range vectorSearchExpectedDeveloperConferences {
		assert.Equal(t, result[i], vectorSearchExpectedDeveloperConferences[i])
	}
	assert.True(t, isMatched(result, vectorSearchExpectedDeveloperConferences))
}

/*
 * [TestVectorSearchDeveloperConferences_GeminiEmbedding001]
 * vector search countries description with gemini-embedding-001 embedding model
 */
func TestVectorSearchDeveloperConferences_GeminiEmbedding001(t *testing.T) {
	ctx := context.Background()
	geminiClient, _ := genai.NewClient(ctx, &genai.ClientConfig{APIKey: LoadAPIKey()})
	em := embeddingModelMap["gemini-embedding-001"]
	ts := New(typesenseUrl, typesenseApiKey, embedding.NewGeminiEmbedder(geminiClient, em.Model))

	result := make([]string, 0)
	for _, q := range vectorSearchRequestsDeveloperConferences {
		vectorstr, err := vectorString(ctx, ts.embedder, q, "")
		assert.NoError(t, err)

		res, err := doMultiSearch(ctx, ts.client, collDeveloperConferenceGeminiGeminiEmbedding001, q, vectorstr)
		assert.NoError(t, err)
		result = append(result, res)
	}
	for i := range vectorSearchExpectedDeveloperConferences {
		assert.Equal(t, result[i], vectorSearchExpectedDeveloperConferences[i])
	}
	assert.True(t, isMatched(result, vectorSearchExpectedDeveloperConferences))
}

// add documents to collection
func addDocumentToCollection(ctx context.Context, ts *Client, m embeddingModel, info map[string]string, collectionName, lang string) error {
	for k, v := range info {
		vector, err := ts.embedder.Embed(ctx, model.EmbeddingRequest{
			Text:     v,
			Language: lang,
		})
		if err != nil {
			return err
		}
		if vector == nil {
			return fmt.Errorf("vector is nil")
		}
		if len(vector.Embedding) != m.NumDim {
			return fmt.Errorf("vector dimension is not %d", m.NumDim)
		}
		if vector.Text != v {
			return fmt.Errorf("vector text is not %s", v)
		}

		typesenseDoc := map[string]interface{}{
			"id":        k,
			"content":   v,
			"embedding": vector.Embedding,
		}

		_, err = ts.client.Collection(collectionName).Documents().Create(ctx, typesenseDoc)
		if err != nil {
			return err
		}
	}
	return nil
}

func vectorToString(v []float32) string {
	var sb strings.Builder
	for i, f := range v {
		if i > 0 {
			sb.WriteString(",")
		}
		sb.WriteString(strconv.FormatFloat(float64(f), 'f', -1, 64))
	}
	return sb.String()
}

func LoadAPIKey() string {
	if err := godotenv.Load("../../.env"); err != nil {
		return ""
	}
	return os.Getenv("GEMINI_API_KEY")
}

const GopherConTextEn string = `GopherCon is the premier annual global conference for the Go (Golang) programming language, named after its mascot, the gopher. It provides a platform for Go developers to learn, network, and share knowledge about the language's features, tools, and best practices. Organized by a community-driven group, GopherCon fosters an inclusive and supportive Go community by offering state-of-the-art talks and creating opportunities for growth for developers of all experience levels. 
Key Aspects of GopherCon
Focus on Go:
The conference is dedicated to the Go programming language, bringing together enthusiasts, contributors, and professionals. 
Community-Driven:
GopherCon is organized by a community of developers who aim to promote the use of Go and foster an inclusive community. 
Educational Content:
Attendees can expect talks, presentations, and discussions on Go's latest features, concurrency models, and applications in areas like AI and distributed systems. 
Networking and Collaboration:
It offers a valuable space for developers to connect with peers, share experiences, and collaborate on new ideas. 
Global Reach:
While the original conference started in 2014, GopherCon has expanded into a global series of events, with conferences hosted in various countries around the world. 
The "Gopher" Mascot
The name GopherCon comes from the Go language's official mascot, the gopher. The conference's branding and community are strongly tied to this friendly character, which is a common symbol for the Go programming language community. `

const PyconTextEn string = `PyCon is the largest annual international conference for the Python programming language community, organized by the Python Software Foundation (PSF). It features talks, tutorials, and other community events to advance Python and foster its global community, keeping registration costs low due to its non-profit nature. 
What is PyCon?
A Community Gathering:
PyCon is the premier annual event for the Python programming community, bringing together users and developers of the open-source Python language. 
Organized by the PSF:
The conference is produced and underwritten by the Python Software Foundation, a non-profit organization dedicated to promoting Python. 
Global and Local:
While originating in North America, PyCon events are held in many countries worldwide, fostering a global Python community. 
What Happens at PyCon?
Talks and Tutorials:
The event includes numerous talks, tutorials, demonstrations, and presentations on a wide variety of Python-related topics. 
Community Events:
PyCon hosts community events, workshops, and gatherings where attendees can network and share ideas. 
Professional Development:
It offers opportunities for learning new skills, discovering cutting-edge ideas, and connecting with industry professionals. 
Why is PyCon Important?
Fosters Community:
PyCon's main goal is to grow and support the international community of Python programmers. 
Accessibility:
Being a non-profit event, the PSF keeps registration costs low to ensure broader accessibility for attendees. 
Code of Conduct:
PyCon is known for having and adhering to a code of conduct, promoting an inclusive and enjoyable environment for everyone.`

const NodeCongressTextEn string = `NodeCongress is widely recognized as a significant conference focusing on Node.js and JavaScript backends. It covers modern JavaScript backend runtimes, best practices for scaling and maintaining applications, and features talks from prominent speakers in the JavaScript and Node.js communities.
While NodeCongress is a major event specifically for Node.js, other large JavaScript conferences, such as JSNation and JSWORLD Conference, also dedicate substantial attention to Node.js within their broader JavaScript scope. These conferences attract large audiences and feature content relevant to Node.js developers.`

const DevconTextEn string = `Devcon is a premier annual conference by the Ethereum Foundation for builders, thinkers, and enthusiasts in the Ethereum ecosystem, aiming to foster education, collaboration, and innovation in decentralized technology. The event serves as an intensive introduction for newcomers and a global gathering for the established Ethereum community to learn about new tools, protocols, and the culture of decentralized systems. Devcon's programming features a range of technical to human-focused content, including talks, workshops, and peer-to-peer interactions, all delivered in a sponsor-free environment to ensure unbiased information. Etheruem is created by vitalik buterin. 
Key Aspects of Devcon:
Purpose:
To empower the community to build and use decentralized systems, pushing the boundaries of decentralized protocols, tools, and culture. 
Audience:
Developers, thinkers, makers, and anyone interested in Ethereum's mission to make decentralized systems accessible to the world. 
Format:
A four-day event that rotates globally, featuring talks, workshops, and other interactive sessions. 
Programming:
Content covers a wide spectrum from deep technical discussions on the base protocol to broader topics on the application layer, with sessions tailored to different expertise levels. 
Community:
It's described as a "global family reunion" for the Ethereum community, providing a unique opportunity to connect with peers and gain a deeper understanding of the ecosystem's culture. 
Structure:
Devcon is the Ethereum Foundation's principal event, distinguished from Devconnect, which is a separate collaborative initiative of smaller, topic-focused community events. 
Philosophy:
As a sponsor-free event, Devcon focuses on delivering unbiased overviews of Ethereum's state, values, and future direction, with speakers and talks selected based on merit rather than financial incentives. `
