package loaders

import (
	"context"
	"loader/textsplitter"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWIZ_LoadAndSplit(t *testing.T) {
	t.Parallel()
	file, err := os.Open("../yaoapp/data/test2.ziw")
	require.NoError(t, err)
	finfo, err := file.Stat()
	loader := NewWIZ(file,finfo.Size())

	splitter2 := textsplitter.NewRecursiveCharacter()
	docs, err := loader.LoadAndSplit(context.Background(),splitter2)

	
	require.NoError(t, err)
	require.Len(t, docs, 1)

	expectedPageContent := "Foo Bar Baz"
	assert.Equal(t, expectedPageContent, docs[0].PageContent)

	expectedMetadata := map[string]any{}
	assert.Equal(t, expectedMetadata, docs[0].Metadata)
}
