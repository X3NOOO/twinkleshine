package ai

import (
	"context"
	"errors"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/googleai"
	"github.com/tmc/langchaingo/vectorstores"
	"github.com/tmc/langchaingo/vectorstores/qdrant"
)

type TwinkleshineAI struct {
	ctx     context.Context
	options llms.CallOptions
	Model   llms.Model
	VDB     vectorstores.VectorStore
}

func getLLM(ctx context.Context) (llms.Model, *embeddings.EmbedderImpl, error) {
	providerEnv, ok := os.LookupEnv("LLM_PROVIDER")
	if !ok {
		return nil, nil, errors.New("LLM_PROVIDER not set")
	}

	apiKey, ok := os.LookupEnv("LLM_API_KEY")
	if !ok {
		return nil, nil, errors.New("LLM_API_KEY not set")
	}
	apiKey = strings.TrimSpace(apiKey)

	// GoogleAI specifically implements creating embeddings, i dont know about other providers, might need a rework
	var llm interface {
		llms.Model
		embeddings.EmbedderClient
	}
	var embedder *embeddings.EmbedderImpl
	var err error = nil

	switch strings.ToLower(strings.TrimSpace(providerEnv)) {
	case "google":
		llm, err := googleai.New(ctx, googleai.WithAPIKey(apiKey))
		if err != nil {
			return nil, nil, err
		}
		embedder, err = embeddings.NewEmbedder(llm)
		if err != nil {
			return nil, nil, err
		}
	default:
		err = errors.New("unknown LLM provider")
	}

	return llm, embedder, err
}

func getOptions() (*llms.CallOptions, error) {
	modelEnv, ok := os.LookupEnv("LLM_MODEL")
	if !ok {
		return nil, errors.New("LLM_MODEL not set")
	}
	model := strings.TrimSpace(modelEnv)

	temperatureEnv, ok := os.LookupEnv("LLM_TEMPERATURE")
	if !ok {
		temperatureEnv = "0.6"
	}
	temperature, err := strconv.ParseFloat(temperatureEnv, 64)
	if err != nil {
		return nil, err
	}

	maxTokensEnv, ok := os.LookupEnv("LLM_MAX_TOKENS")
	if !ok {
		return nil, errors.New("LLM_MAX_TOKENS not set")
	}
	maxTokens, err := strconv.Atoi(maxTokensEnv)
	if err != nil {
		return nil, err
	}

	options := llms.CallOptions{
		Model:          model,
		Temperature:    temperature,
		MaxTokens:      maxTokens,
		CandidateCount: 1,
	}

	return &options, nil
}

func getVDB(embedderImpl embeddings.EmbedderImpl) (*vectorstores.VectorStore, error) {
	providerEnv, ok := os.LookupEnv("VDB_PROVIDER")
	if !ok {
		return nil, errors.New("VDB_PROVIDER not set")
	}
	provider := strings.TrimSpace(providerEnv)

	apiKeyEnv, ok := os.LookupEnv("VDB_API_KEY")
	if !ok {
		return nil, errors.New("VDB_API_KEY not set")
	}
	apiKey := strings.TrimSpace(apiKeyEnv)

	hostEnv, ok := os.LookupEnv("VDB_HOST")
	if !ok {
		return nil, errors.New("VDB_HOST not set")
	}
	host, err := url.Parse(hostEnv)
	if err != nil {
		return nil, err
	}

	collectionName, ok := os.LookupEnv("VDB_COLLECTION_NAME")
	if !ok {
		return nil, errors.New("VDB_COLLECTION_NAME not set")
	}
	collectionName = strings.TrimSpace(collectionName)

	var vdb vectorstores.VectorStore
	switch strings.ToLower(strings.TrimSpace(provider)) {
	case "qdrant":
		vdb, err = qdrant.New(qdrant.WithURL(*host), qdrant.WithAPIKey(apiKey), qdrant.WithCollectionName(collectionName), qdrant.WithEmbedder(&embedderImpl))
	default:
		err = errors.New("unknown VDB provider")
	}

	return &vdb, err
}

func NewAI() (*TwinkleshineAI, error) {
	ctx := context.Background()

	llm, embedder, err := getLLM(ctx)
	if err != nil {
		return nil, err
	}

	options, err := getOptions()
	if err != nil {
		return nil, err
	}

	vdb, err := getVDB(*embedder)
	if err != nil {
		return nil, err
	}

	return &TwinkleshineAI{
		ctx:     ctx,
		Model:   llm,
		options: *options,
		VDB:     *vdb,
	}, nil
}