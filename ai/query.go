package ai

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/tmc/langchaingo/llms"
)

func (a *TwinkleshineAI) fetchKnowledge(text string) (string, error) {
	matches, err := a.VDB.SimilaritySearch(a.ctx, text, 10)
	if err != nil {
		return "", err
	}

	if len(matches) <= 0 {
		return "", errors.New("no matches found")
	}

	var parsedKnowledge string
	for i, match := range matches {
		parsedMatch := fmt.Sprintf("Source %d.\nText: %s\n", i+1, match.PageContent)
		if file, ok := match.Metadata["file"].(map[string]any); ok {
			if name, ok := file["name"]; ok {
				parsedMatch += fmt.Sprintf("Filename: %s\n", name)
			}
			if url, ok := file["url"]; ok {
				parsedMatch += fmt.Sprintf("URL: %s\n", url)
			}
		}
		parsedKnowledge += parsedMatch + "\n"
	}

	return parsedKnowledge, nil
}

func (a *TwinkleshineAI) Query(text string) (string, error) {
	if len(text) <= a.options.MinMsgLen {
		return "", errors.New("message is too short")
	}

	knowledge, err := a.fetchKnowledge(text)
	if err != nil {
		log.Printf("Cannot fetch knowledge: %v\n", err)
		err = fmt.Errorf("cannot fetch knowledge: %v", err)
		return "", err
	}

	knowledge = fmt.Sprintf("```\n%s\n```", strings.TrimSpace(knowledge))

	rsp, err := a.Model.GenerateContent(
		a.ctx,
		[]llms.MessageContent{
			{
				Role: llms.ChatMessageTypeSystem,
				Parts: []llms.ContentPart{
					llms.TextPart(a.options.SystemPrompt),
					llms.TextPart(strings.ReplaceAll(a.options.RagPrompt, "{RAG_KNOWLEDGE}", knowledge)),
				},
			},
			{
				Role: llms.ChatMessageTypeHuman,
				Parts: []llms.ContentPart{
					llms.TextPart(text),
				},
			},
		},
		llms.WithOptions(a.options.CallOptions),
	)
	if err != nil {
		err = fmt.Errorf("cannot generate response: %v", err)
		return "", err
	}

	response := rsp.Choices[0].Content

	log.Println("Generated response:", response)

	return response, nil
}
