package jlrickert

import "embed"

//go:embed all:content/** all:themes/** data.yaml example.html
var Assets embed.FS
