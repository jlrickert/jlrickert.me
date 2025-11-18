package assets

import (
	"bytes"
	"context"
	"embed"
	"fmt"
	"maps"
	"path/filepath"

	"github.com/jlrickert/jlrickert.me"
	"gopkg.in/yaml.v3"
)

// AssetManager manages embedded assets including posts and data files.
type AssetManager struct {
	assets embed.FS
}

// NewAssetManager creates and returns a new AssetManager instance.
func NewAssetManager() *AssetManager {
	return &AssetManager{assets: jlrickert.Assets}
}

// GetPost retrieves and parses a post by slug, converting its markdown
// content to HTML.
func (m *AssetManager) GetPost(
	ctx context.Context,
	slug string,
) (*Post, error) {
	data, err := m.assets.ReadFile(filepath.Join("posts", slug+".md"))
	if err != nil {
		return nil, fmt.Errorf(
			"slug \"%s\" does not exist: %w",
			slug,
			err,
		)
	}

	meta, content, err := m.parseFrontmatter(data)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	if err := Markdown.Convert(content, &buf); err != nil {
		return nil, err
	}

	return &Post{
		Slug:    slug,
		Content: buf.Bytes(),
		Meta:    meta,
	}, nil
}

// GetData retrieves and parses the data.yaml file into a Data struct.
func (m *AssetManager) GetData(ctx context.Context) (*Data, error) {
	data, err := m.assets.ReadFile("data.yaml")
	if err != nil {
		return nil, fmt.Errorf(
			"failed to read data.yaml: %w",
			err,
		)
	}

	return LoadData(data)
}

// parseFrontmatter extracts YAML frontmatter from markdown content.
//
// It expects frontmatter to be delimited by --- at the start and
// separated from content by \n---\n. Returns metadata map, remaining
// content, and any parsing errors.
func (m *AssetManager) parseFrontmatter(data []byte) (
	map[string]any,
	[]byte,
	error,
) {
	meta := make(map[string]any)
	content := data

	if !bytes.HasPrefix(data, []byte("---")) {
		return meta, content, nil
	}

	parts := bytes.SplitN(data, []byte("\n---\n"), 2)
	if len(parts) < 2 {
		return meta, content, nil
	}

	frontmatterStr := string(bytes.TrimPrefix(
		parts[0],
		[]byte("---\n"),
	))
	content = parts[1]

	yamlMap := make(map[string]any)
	if err := yaml.Unmarshal(
		[]byte(frontmatterStr),
		yamlMap,
	); err != nil {
		return meta, content, fmt.Errorf(
			"failed to parse frontmatter: %w",
			err,
		)
	}

	maps.Copy(meta, yamlMap)

	return meta, content, nil
}
