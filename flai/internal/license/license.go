// Package license carries the text of the license system-flow is distributed
// under, so that every flai binary holds a copy, as the license itself requires
// (S-0231). LICENSE.md here is a copy of the one at the repository root: Go
// embeds nothing above the module, and the test keeps the two the same.
package license

import (
	_ "embed"
	"strings"
)

//go:embed LICENSE.md
var text string

// Name is the license's title: the first heading of LICENSE.md.
func Name() string {
	first, _, _ := strings.Cut(text, "\n")
	return strings.TrimSpace(strings.TrimPrefix(first, "#"))
}

// Text is the whole license, as LICENSE.md at the repository root reads.
func Text() string { return text }
