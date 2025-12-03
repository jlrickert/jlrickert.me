package jlrickert

import "embed"

//go:embed data/posts/** data/data.yaml
var Assets embed.FS
