package loaders

import (
	"context"
	"io"
	"loader/pptx"
	"loader/schema"
	"loader/textsplitter"
)

// HTML loads parses and sanitizes html content from an io.Reader.
type PPTX struct {
	r io.ReaderAt
	s int64
}

var _ Loader = Docx{}

// NewHTML creates a new html loader with an io.Reader.
func NewPPTX(r io.ReaderAt, size int64) PPTX {
	return PPTX{r, size}
}

// Load reads from the io.Reader and returns a single document with the data.
func (d PPTX) Load(_ context.Context) ([]schema.Document, error) {

	alltexts, err := pptx.Read(d.r, d.s)
	if err != nil {
		return nil, err
	}

	docs := []schema.Document{}
	slides := make([]string, 0)
	for _, p := range alltexts {
		line := ""
		for _, t := range p {
			line += t + ""
		}
		if line != "" {
			slides = append(slides, line)
		}

	}
	numPages := len(slides)
	for i, v := range slides {
		docs = append(docs, schema.Document{
			PageContent: v,
			Metadata: map[string]any{
				"slide":        i,
				"total_slides": numPages,
			},
		})
	}
	return docs, nil
}

// LoadAndSplit reads text data from the io.Reader and splits it into multiple
// documents using a text splitter.
func (d PPTX) LoadAndSplit(ctx context.Context, splitter textsplitter.TextSplitter) ([]schema.Document, error) {
	docs, err := d.Load(ctx)
	if err != nil {
		return nil, err
	}
	return textsplitter.SplitDocuments(splitter, docs)
}
