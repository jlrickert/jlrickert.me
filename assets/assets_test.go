package assets

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewAssetsManager(t *testing.T) {
	manager := NewAssetManager()
	assert.NotNil(t, manager)
}

func TestGetPost(t *testing.T) {
	manager := NewAssetManager()
	post, err := manager.GetPost(context.Background(), "first-post")

	require.NoError(t, err)
	assert.NotNil(t, post)
	assert.Equal(t, "first-post", post.Slug)
}

func TestGetPostNotFound(t *testing.T) {
	manager := NewAssetManager()
	post, err := manager.GetPost(context.Background(), "nonexistent-post")

	assert.Error(t, err)
	assert.Nil(t, post)
}

func TestGetPostContent(t *testing.T) {
	manager := NewAssetManager()
	post, err := manager.GetPost(context.Background(), "first-post")

	require.NoError(t, err)
	assert.NotEmpty(t, post.Content)
	assert.Contains(t, string(post.Content), "Welcome to My Blog")
}

func TestGetPostMetadata(t *testing.T) {
	manager := NewAssetManager()
	post, err := manager.GetPost(context.Background(), "first-post")

	require.NoError(t, err)
	assert.NotEmpty(t, post.Meta)
	assert.Equal(t, "Welcome to My Blog", post.Meta["title"])
	assert.Equal(t, "2025-11-17", post.Meta["date"])
}

func TestPostMetadata(t *testing.T) {
	manager := NewAssetManager()
	post, err := manager.GetPost(context.Background(), "first-post")

	require.NoError(t, err)
	assert.Equal(t, "Welcome to My Blog", post.Title())
	assert.Equal(t, "2025-11-17", post.Date())
	assert.Equal(t,
		"An introduction to my portfolio and what you can expect here",
		post.Description(),
	)
}

func TestPostTags(t *testing.T) {
	manager := NewAssetManager()
	post, err := manager.GetPost(context.Background(), "first-post")

	require.NoError(t, err)
	tags := post.Tags()
	assert.Equal(t, 2, len(tags))
	assert.Contains(t, tags, "welcome")
	assert.Contains(t, tags, "introduction")
}

func TestGetPostWithContext(t *testing.T) {
	manager := NewAssetManager()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	post, err := manager.GetPost(ctx, "first-post")

	assert.Error(t, err)
	assert.Nil(t, post)
}

func TestGetPostEmptySlug(t *testing.T) {
	manager := NewAssetManager()
	post, err := manager.GetPost(context.Background(), "")

	assert.Error(t, err)
	assert.Nil(t, post)
}
