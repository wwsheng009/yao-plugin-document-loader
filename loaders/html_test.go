package loaders

import (
	"context"
	"loader/textsplitter"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHTMLLoader(t *testing.T) {
	t.Parallel()
	file, err := os.Open("./testdata/test.html")
	require.NoError(t, err)

	loader := NewHTML(file)

	docs, err := loader.Load(context.Background())
	require.NoError(t, err)
	require.Len(t, docs, 1)

	content := docs[0].PageContent
	expected := []string{
		"The content",
		"goes here",
		"and here",
	}
	notexpected := []string{
		"console.log(",
		"<title>langchaingo html example",
		"</title>",
		"<footer>",
		"XSS1",
		"onmouseover",
	}

	assert.Contains(t, content, expected[0])
	assert.Contains(t, content, expected[1])
	assert.Contains(t, content, expected[2])
	assert.NotContains(t, content, notexpected[0])
	assert.NotContains(t, content, notexpected[1])
	assert.NotContains(t, content, notexpected[2])
	assert.NotContains(t, content, notexpected[3])
	assert.NotContains(t, content, notexpected[4])
	assert.NotContains(t, content, notexpected[5])

	expectedMetadata := map[string]any{}
	assert.Equal(t, expectedMetadata, docs[0].Metadata)
}

func TestHTML_LoadAndSplit(t *testing.T) {

	t.Parallel()
	file, err := os.Open("./testdata/test.html")
	require.NoError(t, err)

	loader := NewHTML(file)
	docs, err :=loader.LoadAndSplit(context.Background(), textsplitter.NewMarkdownTextSplitter())
	require.NoError(t, err)
	// require.Len(t, docs, 1)
	//print the page content in docs
	for _, doc := range docs {
		println(doc.PageContent)
	}
}
