package assets

import (
	"bytes"
	"context"
	"embed"
	"fmt"
	"path/filepath"

	"github.com/jlrickert/jlrickert.me"
	"gopkg.in/yaml.v3"
)

type AssetManager struct {
	assets embed.FS
}

func NewAssetManager() *AssetManager {
	return &AssetManager{assets: jlrickert.Assets}
}

func (m *AssetManager) parseFrontmatter(data []byte) (map[string]string, []byte, error) {
	meta := map[string]string{}
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

	for key, val := range yamlMap {
		meta[key] = fmt.Sprint(val)
	}

	return meta, content, nil
}

func (m *AssetManager) GetPost(ctx context.Context, slug string) (*Post, error) {
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
