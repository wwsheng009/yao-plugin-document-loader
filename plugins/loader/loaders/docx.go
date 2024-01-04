package loaders

import (
	"context"
	"fmt"
	"io"

	"loader/schema"
	"loader/textsplitter"

	"github.com/unidoc/unioffice/document"
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

	doc, err := document.Read(d.r, d.s)
	if err != nil {
		return nil, err
	}
	// defer doc.Close()

	pagecontent := ""

	for _, para := range doc.Paragraphs() {
		for _, run := range para.Runs() {
			fmt.Print(run.Text())
			pagecontent += run.Text() + "\n"
		}
	}

	return []schema.Document{
		{
			PageContent: pagecontent,
			Metadata: map[string]any{
				"total_paragraphs": len(doc.Paragraphs()),
			},
		},
	}, nil
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
