package ai

import (
	"context"
	"errors"
	"os"
	"strconv"
	"strings"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/googleai"
)

type TwinkleshineAI struct {
	ctx     context.Context
	options llms.CallOptions
	Model   llms.Model
}

func NewAI() (*TwinkleshineAI, error) {
	ctx := context.Background()

	llmProvider, ok := os.LookupEnv("LLM_PROVIDER")
	if !ok {
		return nil, errors.New("LLM_PROVIDER not set")
	}

	apiKey, ok := os.LookupEnv("LLM_API_KEY")
	if !ok {
		return nil, errors.New("LLM_API_KEY not set")
	}
	apiKey = strings.TrimSpace(apiKey)

	var llmModel llms.Model
	var err error = nil

	switch strings.ToLower(strings.TrimSpace(llmProvider)) {
	case "google":
		llmModel, err = googleai.New(ctx, googleai.WithAPIKey(apiKey))
	default:
		err = errors.New("unknown LLM provider")
	}
	if err != nil {
		return nil, err
	}

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

	return &TwinkleshineAI{
		ctx:     ctx,
		Model:   llmModel,
		options: options,
	}, nil
}
