// Package prompts embeds the app's prompts in the binary.
package prompts

import _ "embed"

//go:embed support.md
var Support string
