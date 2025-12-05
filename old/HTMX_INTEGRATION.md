# HTMX Integration Guide

This document describes how htmx is integrated into the portfolio theme for progressive enhancement.

## Overview

HTMX is loaded from the CDN in the base template and configured for default behavior. The theme uses htmx attributes to
enable dynamic content loading without full page refreshes while maintaining progressive enhancement (works without
JavaScript).

## Base Template Setup

File: `themes/green-nebula-terminal/default/_base.html`

```html
<script src="https://cdn.jsdelivr.net/npm/htmx.org@2.0.8/dist/htmx.min.js"></script>
<script>
    htmx.config.defaultIndicatorStyle = "spinner";
    htmx.config.timeout = 10000;
    htmx.config.historyCacheSize = 10;
</script>
```

## Endpoints

### Posts Partial

**GET `/api/posts/partial`**

Query Parameters:

- `tag` - Filter posts by tag
- `search` - Search posts by title/description
- `page` - Pagination support (future enhancement)

Returns HTML fragment of filtered posts list.

Usage in template:

```html
<div hx-get="/api/posts/partial?tag=go" hx-target="#posts-list">
    Load Go posts
</div>
```

### Experience Partial

**GET `/api/experience/partial`**

Returns expandable experience section.

Usage in template:

```html
<div hx-get="/api/experience/partial" hx-trigger="load">
    Loading experience...
</div>
```

### Skills Partial

**GET `/api/skills/partial`**

Returns expandable skills section with category organization.

Usage in template:

```html
<div hx-get="/api/skills/partial" hx-trigger="load">
    Loading skills...
</div>
```

### Theme Switcher

**POST `/api/theme`**

Form Data:

- `theme` - Theme name (default, dark, light, etc.)

Sets cookie for persistence and returns success.

Usage in template:

```html
<button hx-post="/api/theme" hx-vals='{"theme":"dark"}'>
    Dark Mode
</button>
```

## HTMX Attributes Used

### Triggers

- `hx-trigger="load"` - Load content when element is rendered
- `hx-trigger="click"` - Load content on click
- `hx-trigger="change"` - Load content on form change

### Content Swapping

- `hx-swap="innerHTML"` - Replace element content
- `hx-swap="outerHTML"` - Replace entire element
- `hx-swap="beforeend"` - Append to element

### Indicators

- `hx-indicator=".htmx-indicator"` - Show spinner during request

## Progressive Enhancement

All htmx-enhanced features work without JavaScript:

1. Regular links/forms still navigate normally
2. Theme switcher sets cookie for future requests
3. Expandable sections use semantic HTML with details/summary (future enhancement)

## TODO Implementations

1. **Post Filtering** - Implement backend logic in `handlePostsPartial()`
   - Parse query parameters for tag and search
   - Filter posts by criteria
   - Render posts-list.html partial

2. **Skills Expansion** - Implement backend logic in `handleSkillsPartial()`
   - Parse data.yaml for skills structure
   - Support click-to-expand behavior
   - Render skills-section.html partial

3. **Experience Details** - Implement `handleExperiencePartial()`
   - Load experience data from data.yaml
   - Support expandable job descriptions
   - Render experience-section.html partial

4. **Theme Validation** - Update `handleThemeSwitch()`
   - Validate theme against available themes
   - Apply theme to next page load
   - Consider OOB swap for immediate visual change

## Performance Considerations

- Cache control headers set on partial responses (10 minutes for posts, 1 hour for data)
- GZIP compression enabled via chi middleware
- 5-second timeout on context operations
- History cache limits (10 entries) to prevent memory bloat

## Browser Compatibility

HTMX works on all modern browsers (ES6+). For older browser support, add htmx-compat.js or use server-side fallbacks.

## Security Notes

- Form parameters validated on server
- Theme cookie is HttpOnly (prevents XSS access)
- CSRF protection should be implemented if forms submit state changes
- All user input sanitized before rendering
