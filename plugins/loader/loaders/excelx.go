package loaders

import (
	"context"
	"io"
	"loader/schema"
	"loader/textsplitter"

	"github.com/xuri/excelize/v2"
)

// HTML loads parses and sanitizes html content from an io.Reader.
type Excelx struct {
	r    io.Reader
	opts []excelize.Options
}

var _ Loader = HTML{}

// NewHTML creates a new html loader with an io.Reader.
func NewExcelx(r io.Reader, opts ...excelize.Options) Excelx {
	return Excelx{r, opts}
}

// Load reads from the io.Reader and returns a single document with the data.
func (e Excelx) Load(_ context.Context) ([]schema.Document, error) {
	f, err := excelize.OpenReader(e.r, e.opts...)
	if err != nil {
		return nil, err
		// log.Fatalf("error opening Excel file: %s", err)
	}
	// defer func() {
	// 	if err := f.Close(); err != nil {
	// 		log.Fatalf("error closing Excel file: %s", err)
	// 	}
	// }()

	// Get all sheet names
	sheets := f.GetSheetList()

	// Iterate over each sheet
	docs := []schema.Document{}

	numSheets := len(sheets)

	for i, sheet := range sheets {
		// fmt.Println("Reading Sheet:", sheet)

		rows, err := f.GetRows(sheet)
		if err != nil {
			return nil, err
			// log.Fatalf("error getting rows from sheet %s: %s", sheet, err)
		}

		pagecontent := ""
		for _, row := range rows {
			for _, cell := range row {
				pagecontent += cell + "\t"
				// fmt.Print(cell, "\t")
			}
			pagecontent += "\n"
		}

		docs = append(docs, schema.Document{
			PageContent: pagecontent,
			Metadata: map[string]any{
				"shee":         i,
				"sheet_name":   sheet,
				"total_sheets": numSheets,
			},
		})
	}

	return docs, nil
}

// LoadAndSplit reads text data from the io.Reader and splits it into multiple
// documents using a text splitter.
func (e Excelx) LoadAndSplit(ctx context.Context, splitter textsplitter.TextSplitter) ([]schema.Document, error) {
	docs, err := e.Load(ctx)
	if err != nil {
		return nil, err
	}
	return textsplitter.SplitDocuments(splitter, docs)
}
