package loaders

import (
	"context"
	"loader/textsplitter"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewDocx(t *testing.T) {
	t.Parallel()
	f, err := os.Open("../yaoapp/data/test.docx")
	require.NoError(t, err)

	finfo, err := f.Stat()
	loader := NewDocx(f, finfo.Size())

	splitter := textsplitter.NewTokenSplitter()
	splitter2 := textsplitter.NewRecursiveCharacter()
	docs2, err := loader.LoadAndSplit(context.Background(),splitter2)

	docs, err := loader.LoadAndSplit(context.Background(),splitter)
	require.NoError(t, err)
	require.Len(t, docs, 1)

	// expectedPageContent := "Foo Bar Baz"
	assert.Equal(t, docs2[0].PageContent, docs[0].PageContent)

	expectedMetadata := map[string]any{}
	assert.Equal(t, expectedMetadata, docs[0].Metadata)
}
