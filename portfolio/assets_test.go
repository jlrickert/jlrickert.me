package portfolio

import (
	"bytes"
	"context"
	"testing"
	"time"

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
}

func TestPostMetadata(t *testing.T) {
	manager := NewAssetManager()
	post, err := manager.GetPost(context.Background(), "first-post")

	require.NoError(t, err)
	assert.Equal(t, "Welcome to My Blog", post.Title())
	assert.Equal(t, time.Date(2025, 11, 17, 0, 0, 0, 0, time.Local),
		post.Date())
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
	assert.Contains(t, tags, "introduction")
	assert.Contains(t, tags, "welcome")
}

func TestGetPostEmptySlug(t *testing.T) {
	manager := NewAssetManager()
	post, err := manager.GetPost(context.Background(), "")

	assert.Error(t, err)
	assert.Nil(t, post)
}

func TestContentIsRenderedHTML(t *testing.T) {
	manager := NewAssetManager()
	post, err := manager.GetPost(context.Background(), "first-post")

	require.NoError(t, err)
	require.NotNil(t, post)

	content := string(post.Content)
	assert.Contains(t, content, "<h1")
	assert.Contains(t, content, "<h2")
	assert.Contains(t, content, "<p")
	assert.Contains(t, content, "<ul")
	assert.Contains(t, content, "<li")
	assert.NotContains(t, content, "<br")
	assert.True(t, bytes.Contains([]byte(content), []byte("<")),
		"Content should be HTML with opening tags")
	assert.True(t, bytes.Contains([]byte(content), []byte(">")),
		"Content should be HTML with closing brackets")
}

func TestGetData(t *testing.T) {
	manager := NewAssetManager()
	data, err := manager.GetData(context.Background())

	require.NoError(t, err)
	assert.NotNil(t, data)
}

func TestGetDataNotNil(t *testing.T) {
	manager := NewAssetManager()
	data, err := manager.GetData(context.Background())

	require.NoError(t, err)
	assert.NotNil(t, data.Name)
	assert.NotNil(t, data.Title)
	assert.NotNil(t, data.Email)
}

func TestGetDataBasicInfo(t *testing.T) {
	manager := NewAssetManager()
	data, err := manager.GetData(context.Background())

	require.NoError(t, err)
	assert.Equal(t, "Jared Rickert", data.Name)
	assert.Equal(t, "Minneapolis, MN", data.Location)
	assert.Equal(t, "jaredrickert52@gmail.com", data.Email)
	assert.Equal(t, "(320) 360-8538", data.Phone)
}

func TestGetDataExperience(t *testing.T) {
	manager := NewAssetManager()
	data, err := manager.GetData(context.Background())

	require.NoError(t, err)
	assert.NotEmpty(t, data.Experience)
	assert.True(t, len(data.Experience) > 0)

	firstExp := data.Experience[0]
	assert.NotEmpty(t, firstExp.Title)
	assert.NotEmpty(t, firstExp.Company)
	assert.NotEmpty(t, firstExp.StartDate)
}

func TestGetDataSkills(t *testing.T) {
	manager := NewAssetManager()
	data, err := manager.GetData(context.Background())

	require.NoError(t, err)
	assert.NotEmpty(t, data.Skills)
	assert.NotEmpty(t, data.Skills.Languages)
	assert.NotEmpty(t, data.Skills.Frontend)
	assert.NotEmpty(t, data.Skills.Backend)
	assert.NotEmpty(t, data.Skills.CloudDevOps)
	assert.NotEmpty(t, data.Skills.Databases)
	assert.NotEmpty(t, data.Skills.Tools)
}

func TestGetDataEducation(t *testing.T) {
	manager := NewAssetManager()
	data, err := manager.GetData(context.Background())

	require.NoError(t, err)
	assert.NotEmpty(t, data.Education)
	assert.True(t, len(data.Education) > 0)

	firstEdu := data.Education[0]
	assert.NotEmpty(t, firstEdu.School)
	assert.NotEmpty(t, firstEdu.Degree)
}

func TestGetDataCertifications(t *testing.T) {
	manager := NewAssetManager()
	data, err := manager.GetData(context.Background())

	require.NoError(t, err)
	assert.NotEmpty(t, data.Certifications)
	assert.True(t, len(data.Certifications) > 0)

	firstCert := data.Certifications[0]
	assert.NotEmpty(t, firstCert.Name)
	assert.NotEmpty(t, firstCert.Issued)
	assert.NotEmpty(t, firstCert.CredentialID)
}

func TestGetDataSummary(t *testing.T) {
	manager := NewAssetManager()
	data, err := manager.GetData(context.Background())

	require.NoError(t, err)
	assert.NotEmpty(t, data.Summary)
	assert.Contains(t, data.Summary, "Software Engineer")
}

func TestGetDataSocialLinks(t *testing.T) {
	manager := NewAssetManager()
	data, err := manager.GetData(context.Background())

	require.NoError(t, err)
	assert.NotEmpty(t, data.LinkedIn)
	assert.NotEmpty(t, data.Portfolio)
	assert.Contains(t, data.LinkedIn, "linkedin.com")
	assert.Contains(t, data.Portfolio, "jlrickert.me")
}
