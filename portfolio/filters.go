package portfolio

import (
	"html/template"
	"slices"
	"strconv"
	"strings"
	"time"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// TemplateFuncs returns a FuncMap of custom template functions
func TemplateFuncs() template.FuncMap {
	return template.FuncMap{
		"truncate":   truncate,
		"formatDate": formatDate,
		"contains":   contains,
		"join":       strings.Join,
		"timeAgo":    timeAgo,
		"humanize":   humanize,
	}
}

// truncate limits a string to a maximum length and adds ellipsis if truncated
func truncate(s string, length int) string {
	if len(s) <= length {
		return s
	}
	return s[:length] + "..."
}

// formatDate converts a YYYY-MM-DD date string to a more readable format
func formatDate(dateStr string) string {
	if dateStr == "" {
		return ""
	}

	// Parse the date string
	t, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return dateStr
	}

	// Format as "Jan 02, 2006"
	return t.Format("Jan 02, 2006")
}

// contains checks if a slice contains a string value
func contains(slice []string, item string) bool {
	return slices.Contains(slice, item)
}

// timeAgo calculates the time elapsed since a given date string
func timeAgo(dateStr string) string {
	if dateStr == "" {
		return ""
	}

	t, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return dateStr
	}

	now := time.Now()
	duration := now.Sub(t)

	days := int(duration.Hours() / 24)
	months := days / 30
	years := months / 12

	switch {
	case years > 0:
		if years == 1 {
			return "1 year ago"
		}
		return strconv.Itoa(years) + " years ago"
	case months > 0:
		if months == 1 {
			return "1 month ago"
		}
		return strconv.Itoa(months) + " months ago"
	case days > 0:
		if days == 1 {
			return "1 day ago"
		}
		return strconv.Itoa(days) + " days ago"
	default:
		return "Today"
	}
}

// humanize converts a technical string to a more readable format
// e.g., "go_programming" -> "Go Programming"
func humanize(s string) string {
	// Replace underscores and hyphens with spaces
	s = strings.NewReplacer(
		"_", " ",
		"-", " ",
	).Replace(s)

	// Title case each word
	return cases.Title(language.AmericanEnglish).String(strings.ToLower(s))
}
