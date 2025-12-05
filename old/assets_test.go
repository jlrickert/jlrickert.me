package jlrickert

import (
	"io/fs"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAssetsEmbedded(t *testing.T) {
	// Verify that Assets FS is not nil
	assert.NotNil(t, Assets)
}

func TestAssetsContainsDataYAML(t *testing.T) {
	// Verify data.yaml exists in embedded assets
	_, err := fs.Stat(Assets, "data.yaml")
	require.NoError(t, err, "data.yaml should exist in embedded assets")
}

func TestAssetsContainsPosts(t *testing.T) {
	// Verify posts directory exists
	entries, err := fs.ReadDir(Assets, "content/blog")
	require.NoError(t, err, "posts directory should exist in embedded assets")

	// Verify at least one post file exists
	assert.NotEmpty(t, entries, "posts directory should not be empty")

	// Verify first-post.md exists
	_, err = fs.Stat(Assets, "content/blog/first-post.md")
	require.NoError(t, err, "first-post.md should exist in embedded assets")
}

func TestAssetsDataYAMLIsReadable(t *testing.T) {
	// Verify we can read the data.yaml file
	data, err := fs.ReadFile(Assets, "data.yaml")
	require.NoError(t, err, "should be able to read data.yaml")
	assert.NotEmpty(t, data, "data.yaml should contain data")
}

func TestAssetsFirstPostIsReadable(t *testing.T) {
	// Verify we can read the first post
	postContent, err := fs.ReadFile(Assets, "content/blog/first-post.md")
	require.NoError(t, err, "should be able to read blog/first-post.md")
	assert.NotEmpty(t, postContent, "first-post.md should contain content")
	assert.Contains(t, string(postContent), "---", "post should have frontmatter delimiter")
}
