package loaders

import (
	"context"
	"io"
	"regexp"
	"strings"

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
	re := regexp.MustCompile(`^\s*\n`)
	re2 := regexp.MustCompile(`^\s*$`)

	docs := []schema.Document{}
	strs := make([]string, 0)

	line := ""
	for _, r := range paragraphs {
		// replace the \u00a0 in the p.Texts
		for i, t := range r.Texts {
			r.Texts[i].Content = strings.ReplaceAll(t.Content, "\u00a0", "")
			r.Texts[i].Content = re.ReplaceAllString(r.Texts[i].Content, "")
			r.Texts[i].Content = re2.ReplaceAllString(r.Texts[i].Content, "")

		}

		//delete empty line in p.Texts
		for i, t := range r.Texts {
			if t.Content == "" {
				r.Texts = append(r.Texts[:i], r.Texts[i+1:]...)
			}
		}

		for _, t := range r.Texts {
			line += t.Content + "\n"
		}
		// }

		for _, t := range r.Hyperlink {
			//remove the end \n in line
			line = strings.TrimSuffix(line, "\n")
			line += t.Content + "\n"
		}
		// for _, r := range r.Runs {
		// 使用页面分隔符来分割段落
		if len(r.LastRenderedPageBreak) > 0 {
			strs = append(strs, line)
			line = ""
		}

		// 如果需要使用空行来分割段落，可以使用以下代码
		// if line != "" && len(p.Texts) == 0 {
		// 	strs = append(strs, line)
		// 	line = ""
		// }

	}
	if line != "" {
		strs = append(strs, line)
		line = ""
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
