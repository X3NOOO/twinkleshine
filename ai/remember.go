package ai

import (
	"errors"
	"net/http"
	"slices"

	"github.com/X3NOOO/llamaparse-go"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/textsplitter"
)

type Chunk struct {
	Content  string
	Metadata map[string]any
}

func parse(body []byte) (string, error) {
	mime := http.DetectContentType(body)

	if mime == "text/plain" {
		return string(body), nil
	} else if slices.Contains(llamaparse.SUPPORTED_MIME_TYPES, mime) {
		return llamaparse.Parse(body, llamaparse.MARKDOWN, nil, nil, nil, nil)
	} else {
		return "", errors.New("unsupported filetype")
	}
}

func parseFile(body []byte, chunkLength int, chunkOverlap int) ([]string, error) {
	parsed, err := parse(body)
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

func (a *TwinkleshineAI) Remember(text string, metadata map[string]any) error {
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
	chunks, err := parseFile(body, a.options.ChunkLength, a.options.ChunkOverlap)
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
