package portfolio

import "html/template"

type ThemeManager struct {
	defaultTheme string
	validThemes  map[string]bool

	templates map[string]*template.Template // cache
}
