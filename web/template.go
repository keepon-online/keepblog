package templates

import "embed"

//go:embed **/*.html
var Fs embed.FS
