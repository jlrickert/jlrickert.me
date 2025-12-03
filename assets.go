package jlrickert

import "embed"

//go:embed data/posts/** data/data.yaml themes/** example.html
var Assets embed.FS
