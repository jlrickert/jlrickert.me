package portfolio

import (
	"bytes"
	"context"
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"log/slog"
	"maps"
	"path/filepath"

	"github.com/jlrickert/jlrickert.me"
	"gopkg.in/yaml.v3"
)

const DefaultTheme = "green-nebula-terminal"

// AssetManager manages embedded assets including posts and data files.
type AssetManager struct {
	Assets embed.FS
	Logger *slog.Logger
}

// NewAssetManager creates and returns a new AssetManager instance.
func NewAssetManager(theme string) *AssetManager {
	return &AssetManager{Assets: jlrickert.Assets}
}

// GetPage retrieves and parses a post by slug, converting its markdown
// content to HTML.
func (m *AssetManager) GetPage(ctx context.Context, path string) (*Page, error) {
	fp := filepath.Join("content", path)

	data, err := fs.ReadFile(m.Assets, fp)
	if err != nil {
		return nil, fmt.Errorf(
			"slug \"%s\" does not exist: %w",
			path,
			err,
		)
	}

	ext := filepath.Ext(fp)

	switch ext {
	case "md":
		meta, content, err := m.parseFrontmatter(data)
		if err != nil {
			return nil, err
		}
		var buf bytes.Buffer
		if err := Markdown.Convert(content, &buf); err != nil {
			return nil, err
		}

		return &Page{
			Path:    path,
			Type:    "markdown",
			Content: buf.Bytes(),
			Meta:    meta,
		}, nil
	case "html":
		return &Page{
			Path:    path,
			Type:    "html",
			Content: data,
			Meta:    map[string]any{},
		}, nil
	}
	return nil, fmt.Errorf("%s is unsupported", path)
}

// GetData retrieves and parses the data.yaml file into a Data struct.
func (m *AssetManager) GetData(ctx context.Context) (*Data, error) {
	data, err := m.Assets.ReadFile(filepath.Join("data", "data.yaml"))
	if err != nil {
		return nil, fmt.Errorf(
			"failed to read data.yaml: %w",
			err,
		)
	}

	return LoadData(data)
}

// GetTemplateContent retrieves raw template file contents without parsing
func (m *AssetManager) GetTemplateContent(theme, name string) ([]byte, error) {
	path := fmt.Sprintf("themes/%s/default/%s.html", theme, name)
	content, err := m.Assets.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read template file %s: %w", path, err)
	}
	return content, nil
}

func (m *AssetManager) GetTemplate(theme, name string) (*template.Template, error) {
	path := fmt.Sprintf("themes/%s/default/%s.html", theme, name)
	content, err := m.Assets.ReadFile(path)
	if err != nil {
		// Fallback to default theme
		path = fmt.Sprintf("themes/%s/default/%s.html", theme, name)
		content, err = m.Assets.ReadFile(path)
	}

	tmpl, err := template.New(name).Funcs(TemplateFuncs()).Parse(string(content))
	if err != nil {
		m.Logger.Error("failed to parse base template", "template", "_base", "error", err)
		return nil, err
	}

	return tmpl, err
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
