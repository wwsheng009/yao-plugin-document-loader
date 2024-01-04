package loaders

import (
	"context"
	"io"
	"loader/schema"
	"loader/textsplitter"
	"loader/wiz"
)

// HTML loads parses and sanitizes html content from an io.Reader.
type WIZ struct {
	r io.ReaderAt
	s int64
}

var _ Loader = WIZ{}

// NewWIZ creates a new wiz loader with an io.Reader.
func NewWIZ(r io.ReaderAt, size int64) WIZ {
	return WIZ{r, size}
}

// Load reads from the io.Reader and returns a single document with the data.
func (d WIZ) Load(_ context.Context) ([]schema.Document, error) {

	pagecontent, err := wiz.Read(d.r, d.s)
	if err != nil {
		return nil, err
	}
	return []schema.Document{
		{
			PageContent: pagecontent,
			Metadata:    map[string]any{},
		},
	}, nil
}

// LoadAndSplit reads text data from the io.Reader and splits it into multiple
// documents using a text splitter.
func (d WIZ) LoadAndSplit(ctx context.Context, splitter textsplitter.TextSplitter) ([]schema.Document, error) {
	docs, err := d.Load(ctx)
	if err != nil {
		return nil, err
	}
	return textsplitter.SplitDocuments(splitter, docs)
}
