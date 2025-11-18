package jlrickert

import "embed"

//go:embed posts/**.md data.yaml
var Assets embed.FS
