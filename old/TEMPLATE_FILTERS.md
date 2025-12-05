# Template Filters Documentation

Custom template filters are available in your Go templates through the `TemplateFuncs()` function. These filters help
format and manipulate data for display.

## Available Filters

### `truncate`

Limits a string to a maximum length and adds "..." if truncated.

**Usage:**

```html
{{.Summary | truncate 100}}
```

**Example:**

```html
<!-- Input: "This is a very long summary text..." -->
<!-- Output: "This is a very long summary text..." (first 100 chars) -->
<p>{{.Summary | truncate 100}}</p>
```

---

### `formatDate`

Converts a YYYY-MM-DD date string to a readable format (Jan 02, 2006).

**Usage:**

```html
{{.StartDate | formatDate}}
```

**Example:**

```html
<!-- Input: "2025-01-15" -->
<!-- Output: "Jan 15, 2025" -->
<span>Started: {{.StartDate | formatDate}}</span>
```

---

### `contains`

Checks if a slice contains a string value. Returns boolean (true/false).

**Usage:**

```html
{{if contains .Tags "go"}}
<span class="badge">Go Developer</span>
{{end}}
```

**Example:**

```html
{{if contains .Skills "Docker"}}
<span>✓ Docker</span>
{{else}}
<span>○ Docker</span>
{{end}}
```

---

### `join`

Concatenates a slice of strings with a separator.

**Usage:**

```html
{{.Tags | join ", "}}
```

**Example:**

```html
<!-- Input: ["go", "php", "rust"] -->
<!-- Output: "go, php, rust" -->
<p>Skills: {{.Skills | join ", "}}</p>
```

---

### `timeAgo`

Calculates human-readable time elapsed since a date (e.g., "2 days ago").

**Usage:**

```html
{{.PublishedDate | timeAgo}}
```

**Example:**

```html
<!-- Input: "2025-11-15" (5 days ago) -->
<!-- Output: "5 days ago" -->
<span>{{.PublishedDate | timeAgo}}</span>
```

---

### `humanize`

Converts technical strings to readable format by replacing underscores/hyphens with spaces and title-casing.

**Usage:**

```html
{{.Category | humanize}}
```

**Example:**

```html
<!-- Input: "go_programming" -->
<!-- Output: "Go Programming" -->
<h3>{{.Category | humanize}}</h3>

<!-- Input: "cloud-devops" -->
<!-- Output: "Cloud Devops" -->
<span>{{.Skill | humanize}}</span>
```

---

## Combining Filters (Pipes)

You can chain filters together using pipes:

```html
<!-- Truncate, then lowercase output -->
{{.Description | truncate 50}}

<!-- Format date, then display with prefix -->
Started: {{.StartDate | formatDate}}

<!-- Join with commas, then uppercase (using built-in upper filter) -->
{{.Tags | join ", " | upper}}
```

---

## Examples in Context

### Skills Display

```html
<div class="skills">
    {{range .Skills.Languages}}
    <span class="tag">
        {{. | humanize}} {{if contains $.SelectedSkills .}}
        <span class="selected">✓</span>
        {{end}}
    </span>
    {{end}}
</div>
```

### Experience Timeline

```html
<div class="experience">
    {{range .Experience}}
    <div class="entry">
        <h4>{{.Title}} @ {{.Company}}</h4>
        <p class="date">{{.StartDate | formatDate}} - {{.EndDate | formatDate}}</p>
        <p class="description">{{.Summary | truncate 200}}</p>
        <p class="tech">Tech: {{.Technologies | join ", "}}</p>
    </div>
    {{end}}
</div>
```

### Blog Post List

```html
<div class="posts">
    {{range .Posts}}
    <article>
        <h2>{{.Title}}</h2>
        <time>{{.PublishedDate | timeAgo}} ({{.PublishedDate | formatDate}})</time>
        <p>{{.Content | truncate 150}}</p>
        <div class="tags">
            {{.Tags | join " "}}
        </div>
    </article>
    {{end}}
</div>
```

---

## Built-in Go Template Functions

In addition to custom filters, Go templates provide these built-in functions:

- `html` - HTML escape
- `js` - JavaScript escape
- `urlquery` - URL encode
- `upper` - Uppercase
- `lower` - Lowercase
- `len` - Length of slice/string
- `index` - Get map/slice value
- `range` - Loop over items
- `if/else` - Conditional logic
- `with` - Set context

**Example combining custom and built-in:**

```html
<p class="description">{{.Summary | truncate 100 | html}}</p>
<a href="?q={{.Title | urlquery}}">Search</a>
<span class="category">{{.Category | humanize | upper}}</span>
```
