package assets

import (
	"bytes"
	"time"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
	"gopkg.in/yaml.v3"
)

type Post struct {
	Slug    string
	Content []byte
	Meta    map[string]any
}

// Title returns the post title from metadata or the first h1 from
// content.
func (p *Post) Title() string {
	if title, ok := p.Meta["title"].(string); ok && title != "" {
		return title
	}
	return p.extractFirstHeading()
}

// Date returns the post date from metadata or current date.
func (p *Post) Date() time.Time {
	if dateStr, ok := p.Meta["date"].(string); ok && dateStr != "" {
		if parsedDate, err := time.Parse("2006-01-02", dateStr); err == nil {
			return time.Date(
				parsedDate.Year(),
				parsedDate.Month(),
				parsedDate.Day(),
				0, 0, 0, 0,
				time.Local,
			)
		}
	} else if date, ok := p.Meta["date"].(time.Time); ok {
		return time.Date(
			date.Year(),
			date.Month(),
			date.Day(),
			0, 0, 0, 0,
			time.Local,
		)
	}
	return time.Now()
}

// Description returns the post description from metadata or lead
// paragraph.
func (p *Post) Description() string {
	if desc, ok := p.Meta["description"].(string); ok && desc != "" {
		return desc
	}
	return p.extractLeadParagraph()
}

// Tags returns the tags from metadata as a slice of strings.
func (p *Post) Tags() []string {
	// If already a slice of strings, return it
	if tags, ok := p.Meta["tags"].([]string); ok {
		return tags
	}

	// Handle []any from YAML unmarshaling
	if tagsInterface, ok := p.Meta["tags"].([]any); ok {
		tags := make([]string, 0, len(tagsInterface))
		for _, tag := range tagsInterface {
			if tagStr, ok := tag.(string); ok {
				tags = append(tags, tagStr)
			}
		}
		if len(tags) > 0 {
			return tags
		}
	}

	// Try as string and unmarshal
	if tagsStr, ok := p.Meta["tags"].(string); ok {
		if tagsStr == "" {
			return []string{}
		}
		var tags []string
		if err := yaml.Unmarshal([]byte(tagsStr), &tags); err != nil {
			return []string{}
		}
		return tags
	}

	return []string{}
}

// extractFirstHeading finds the first h1 heading in the content and
// returns its text value.
func (p *Post) extractFirstHeading() string {
	doc := Markdown.Parser().Parse(text.NewReader(p.Content))

	var heading string
	ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}

		if h, ok := n.(*ast.Heading); ok && h.Level == 1 {
			heading = p.extractHeadingText(h)
			return ast.WalkStop, nil
		}

		return ast.WalkContinue, nil
	})

	return heading
}

// extractHeadingText extracts text from a heading node.
func (p *Post) extractHeadingText(heading *ast.Heading) string {
	buf := new(bytes.Buffer)
	for child := heading.FirstChild(); child != nil; child = child.NextSibling() {
		if textNode, ok := child.(*ast.Text); ok {
			buf.Write(textNode.Segment.Value(p.Content))
		}
	}
	return buf.String()
}

// extractLeadParagraph finds the first paragraph in the content and
// returns its text value.
func (p *Post) extractLeadParagraph() string {
	md := goldmark.New()
	doc := md.Parser().Parse(text.NewReader(p.Content))

	var paragraph string
	ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}

		if para, ok := n.(*ast.Paragraph); ok && paragraph == "" {
			buf := new(bytes.Buffer)
			for child := para.FirstChild(); child != nil; child =
				child.NextSibling() {
				if textNode, ok := child.(*ast.Text); ok {
					buf.Write(textNode.Segment.Value(p.Content))
				}
			}
			paragraph = buf.String()
			return ast.WalkStop, nil
		}

		return ast.WalkContinue, nil
	})

	return paragraph
}
