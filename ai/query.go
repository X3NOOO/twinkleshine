package ai

import "github.com/tmc/langchaingo/llms"

func (a *TwinkleshineAI) Query(text string) (string, error) {
	return llms.GenerateFromSinglePrompt(a.ctx, a.Model, text, llms.WithOptions(a.options))
}
