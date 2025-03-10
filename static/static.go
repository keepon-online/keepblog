package static

import "embed"

//go:embed *
var Static embed.FS

//go:embed favicon.ico
var Favicon embed.FS

//go:embed robots.txt
var Robots embed.FS
