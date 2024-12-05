package ai

import (
	"errors"
	"fmt"
	"log"
	"slices"
	"strings"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/vectorstores"
)

/*
Retrieval process explained:
1. Fetch the root matches from the vector database.
2. Extract the sources from the root matches.
3. Fetch the matches only from the sources.
*/
func (a *TwinkleshineAI) fetchKnowledge(text string) (string, error) {
	rootMatches, err := a.VDB.SimilaritySearch(a.ctx, text, a.options.RagRootMatchesCount)
	if err != nil {
		return "", err
	}

	if len(rootMatches) <= 0 {
		return "", errors.New("no root matches found")
	}

	var sources []string
	for _, match := range rootMatches {
		if rootMatch, ok := match.Metadata["file"].(map[string]any); ok {
			if source, ok := rootMatch["name"].(string); ok {
				if !slices.Contains(sources, source) {
					sources = append(sources, source)
				}
			}
		}
	}

	// FIXME: This filter is qdrant-specific.
	filter := map[string]any{
		"must": map[string]any{
			"key": "file.name",
			"match": map[string]any{
				"any": sources,
			},
		},
	}

	matches, err := a.VDB.SimilaritySearch(a.ctx, text, a.options.RagMatchesCount, vectorstores.WithFilters(filter))
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

	return response, nil
}
