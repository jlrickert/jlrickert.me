package portfolio

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTruncate(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		length   int
		expected string
	}{
		{
			name:     "truncates long string",
			input:    "This is a very long string that needs truncation",
			length:   10,
			expected: "This is a ...",
		},
		{
			name:     "no truncation needed",
			input:    "Short",
			length:   20,
			expected: "Short",
		},
		{
			name:     "exact length",
			input:    "Exact",
			length:   5,
			expected: "Exact",
		},
		{
			name:     "empty string",
			input:    "",
			length:   10,
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := truncate(tt.input, tt.length)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFormatDate(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "valid date",
			input:    "2025-01-15",
			expected: "Jan 15, 2025",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "invalid date format",
			input:    "15/01/2025",
			expected: "15/01/2025",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatDate(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestContains(t *testing.T) {
	tests := []struct {
		name     string
		slice    []string
		item     string
		expected bool
	}{
		{
			name:     "item exists",
			slice:    []string{"go", "php", "rust"},
			item:     "go",
			expected: true,
		},
		{
			name:     "item not exists",
			slice:    []string{"go", "php", "rust"},
			item:     "python",
			expected: false,
		},
		{
			name:     "empty slice",
			slice:    []string{},
			item:     "go",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := contains(tt.slice, tt.item)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestJoin(t *testing.T) {
	tests := []struct {
		name     string
		slice    []string
		sep      string
		expected string
	}{
		{
			name:     "join with comma",
			slice:    []string{"go", "php", "rust"},
			sep:      ", ",
			expected: "go, php, rust",
		},
		{
			name:     "join with pipe",
			slice:    []string{"a", "b", "c"},
			sep:      " | ",
			expected: "a | b | c",
		},
		{
			name:     "empty slice",
			slice:    []string{},
			sep:      ",",
			expected: "",
		},
		{
			name:     "single item",
			slice:    []string{"only"},
			sep:      ",",
			expected: "only",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := strings.Join(tt.slice, tt.sep)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestHumanize(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "replace underscores",
			input:    "go_programming",
			expected: "Go Programming",
		},
		{
			name:     "replace hyphens",
			input:    "cloud-devops",
			expected: "Cloud Devops",
		},
		{
			name:     "already formatted",
			input:    "Go",
			expected: "Go",
		},
		{
			name:     "multiple underscores",
			input:    "backend_web_development",
			expected: "Backend Web Development",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := humanize(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestTemplateFuncs(t *testing.T) {
	// Verify that TemplateFuncs returns all expected functions
	funcs := TemplateFuncs()

	expectedFuncs := []string{
		"truncate",
		"formatDate",
		"contains",
		"join",
		"timeAgo",
		"humanize",
	}

	for _, name := range expectedFuncs {
		assert.Contains(t, funcs, name, "function %s should be in FuncMap", name)
	}
}
