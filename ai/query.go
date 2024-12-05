package ai

import (
	"errors"

	"github.com/tmc/langchaingo/llms"
)

func (a *TwinkleshineAI) Query(text string) (string, error) {
	if len(text) <= 20 {
		return "", errors.New("message is too short")
	}

	return llms.GenerateFromSinglePrompt(a.ctx, a.Model, text, llms.WithOptions(a.options.CallOptions))
}
