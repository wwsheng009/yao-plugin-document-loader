package loaders

import (
	"context"
	"io"

	"loader/docx"
	"loader/schema"
	"loader/textsplitter"
)

// HTML loads parses and sanitizes html content from an io.Reader.
type Docx struct {
	r io.ReaderAt
	s int64
}

var _ Loader = Docx{}

// NewHTML creates a new html loader with an io.Reader.
func NewDocx(r io.ReaderAt, size int64) Docx {
	return Docx{r, size}
}

// Load reads from the io.Reader and returns a single document with the data.
func (d Docx) Load(_ context.Context) ([]schema.Document, error) {

	paragraphs, err := docx.Read(d.r, d.s)
	if err != nil {
		return nil, err
	}

	docs := []schema.Document{}
	strs := make([]string, 0)
	for _, p := range paragraphs {
		line := ""
		for _, t := range p.Texts {
			line += t.Content + "\n"
			// fmt.Print(t.Content)
		}
		if line != "" {
			strs = append(strs, line)
		}

	}
	numPages := len(strs)
	for i, v := range strs {
		docs = append(docs, schema.Document{
			PageContent: v,
			Metadata: map[string]any{
				"paragraph":       i,
				"total_paragraph": numPages,
			},
		})
	}
	return docs, nil

	// return []schema.Document{
	// 	{
	// 		PageContent: pagecontent,
	// 		Metadata: map[string]any{
	// 			"total_paragraphs": len(paragraphs),
	// 		},
	// 	},
	// }, nil
}

// LoadAndSplit reads text data from the io.Reader and splits it into multiple
// documents using a text splitter.
func (d Docx) LoadAndSplit(ctx context.Context, splitter textsplitter.TextSplitter) ([]schema.Document, error) {
	docs, err := d.Load(ctx)
	if err != nil {
		return nil, err
	}
	return textsplitter.SplitDocuments(splitter, docs)
}
