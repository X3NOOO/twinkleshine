package ai

import (
	"context"
	"errors"
	"net/url"
	"os"
	"strings"

	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/googleai"
	"github.com/tmc/langchaingo/vectorstores"
	"github.com/tmc/langchaingo/vectorstores/qdrant"
	"gopkg.in/yaml.v3"
)

type TwinkleshineAI struct {
	ctx     context.Context
	Options options
	Model   llms.Model
	VDB     vectorstores.VectorStore
}

type Config struct {
	SystemPrompt string `yaml:"system_prompt"`
	Discord      struct {
		Security struct {
			StaffRoleID     string `yaml:"staff_role_id"`
			CooldownSeconds int64  `yaml:"cooldown_seconds"`
		} `yaml:"security"`
	} `yaml:"discord"`
	LLM struct {
		MaxTokens        int     `yaml:"max_tokens"`
		Temperature      float64 `yaml:"temperature"`
		MinMessageLength int     `yaml:"min_message_length"`
	} `yaml:"llm"`
	RAG struct {
		Chunking struct {
			Length  int `yaml:"length"`
			Overlap int `yaml:"overlap"`
		} `yaml:"chunking"`
		Matches struct {
			RootCount int `yaml:"root_count"`
			Count     int `yaml:"count"`
		} `yaml:"matches"`
		RagPrompt string `yaml:"rag_prompt"`
	} `yaml:"rag"`
}

type options struct {
	CallOptions llms.CallOptions
	Config      Config
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
		llm, err = googleai.New(ctx, googleai.WithAPIKey(apiKey))
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

func getOptions() (*options, error) {
	configFilePathEnv, ok := os.LookupEnv("CONFIG_FILE")
	if !ok {
		return nil, errors.New("CONFIG_FILE not set")
	}
	configFilePath := strings.TrimSpace(configFilePathEnv)

	configFile, err := os.ReadFile(configFilePath)
	if err != nil {
		return nil, err
	}

	var cfg Config

	err = yaml.Unmarshal([]byte(configFile), &cfg)
	if err != nil {
		return nil, err
	}

	modelEnv, ok := os.LookupEnv("LLM_MODEL")
	if !ok {
		return nil, errors.New("LLM_MODEL not set")
	}
	model := strings.TrimSpace(modelEnv)

	callOptions := llms.CallOptions{
		Model:          model,
		Temperature:    cfg.LLM.Temperature,
		MaxTokens:      cfg.LLM.MaxTokens,
		CandidateCount: 1,
	}

	options := &options{
		CallOptions: callOptions,
		Config:      cfg,
	}

	return options, nil
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
		Options: *options,
		VDB:     *vdb,
	}, nil
}
