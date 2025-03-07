package ai

import (
	"fmt"
	"net/http"
	"slices"
	"strings"

	"github.com/X3NOOO/llamaparse-go"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/textsplitter"
	"github.com/tmc/langchaingo/vectorstores"
)

type Chunk struct {
	Content  string
	Metadata map[string]any
}

func parse(body []byte, timeout int) (string, error) {
	mime := http.DetectContentType(body)

	mime = strings.Split(mime, ";")[0]

	if mime == "text/plain" {
		return string(body), nil
	} else if slices.Contains(llamaparse.SUPPORTED_MIME_TYPES, mime) {
		return llamaparse.Parse(body, llamaparse.MARKDOWN, nil, nil, &timeout, nil)
	} else {
		return "", fmt.Errorf("unsupported MIME type: %s", mime)
	}
}

func parseFile(body []byte, timeout int, chunkLength int, chunkOverlap int) ([]string, error) {
	parsed, err := parse(body, timeout)
	if err != nil {
		return nil, err
	}

	slitter := textsplitter.NewRecursiveCharacter(
		textsplitter.WithChunkSize(chunkLength),
		textsplitter.WithChunkOverlap(chunkOverlap),
		textsplitter.WithSeparators([]string{"\n\n", "\n", ". ", "! ", "? ", ".\n", "!\n", "?\n"}),
	)

	chunks, err := slitter.SplitText(parsed)
	if err != nil {
		return nil, err
	}

	return chunks, nil
}

func (a *TwinkleshineAI) Exists(key string, values []any) (bool, error) {
	filter := map[string]any{
		"must": map[string]any{
			"key":   key,
			"match": map[string]any{"any": values},
		},
	}

	match, err := a.VDB.SimilaritySearch(a.ctx, "", 1, vectorstores.WithFilters(filter))
	if err != nil {
		return false, err
	}

	return len(match) > 0, nil
}

func (a *TwinkleshineAI) Remember(text string, metadata map[string]any) error {
	a.log.Println("Remembering:", text)
	_, err := a.VDB.AddDocuments(
		a.ctx,
		[]schema.Document{
			{
				PageContent: text,
				Metadata:    metadata,
			},
		},
	)

	return err
}

func (a *TwinkleshineAI) RememberFile(body []byte, metadata map[string]any) error {
	a.log.Println("Remembering file:", metadata)
	chunks, err := parseFile(body, a.Options.Config.RAG.ParseTimeoutSeconds, a.Options.Config.RAG.Chunking.Length, a.Options.Config.RAG.Chunking.Overlap)
	if err != nil {
		return err
	}

	errs := make(chan error, len(chunks))
	for _, chunk := range chunks {
		go func(c string) {
			_, err := a.VDB.AddDocuments(
				a.ctx,
				[]schema.Document{
					{
						PageContent: c,
						Metadata:    metadata,
					},
				},
			)
			errs <- err
		}(chunk)
	}

	for i := 0; i < len(chunks); i++ {
		if err := <-errs; err != nil {
			return err
		}
	}

	return nil

}
