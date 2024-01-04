package loaders

import (
	"os"
	"path/filepath"
	"strings"

	"loader/schema"
)

// NotionDirectoryLoader is a document loader that reads content from pages within a Notion Database.
type NotionDirectoryLoader struct {
	filePath string
	encoding string
}

// NewNotionDirectory creates a new NotionDirectoryLoader with the given file path and encoding.
func NewNotionDirectory(filePath string, encoding ...string) *NotionDirectoryLoader {
	defaultEncoding := "utf-8"

	if len(encoding) > 0 {
		return &NotionDirectoryLoader{
			filePath: filePath,
			encoding: encoding[0],
		}
	}

	return &NotionDirectoryLoader{
		filePath: filePath,
		encoding: defaultEncoding,
	}
}

// Load retrieves data from a Notion directory and returns a list of schema.Document objects.
func (n *NotionDirectoryLoader) Load() ([]schema.Document, error) {
	files, err := os.ReadDir(n.filePath)
	if err != nil {
		return nil, err
	}

	documents := make([]schema.Document, 0, len(files))

	for _, file := range files {
		ext := filepath.Ext(strings.ToLower(file.Name()))

		if file.IsDir() || (ext != ".md" && ext != ".mdx") {
			continue
		}

		filePath := filepath.Join(n.filePath, file.Name())
		text, err := os.ReadFile(filePath)
		if err != nil {
			return nil, err
		}

		metadata := map[string]interface{}{"source": filePath}
		documents = append(documents, schema.Document{PageContent: string(text), Metadata: metadata})
	}

	return documents, nil
}
