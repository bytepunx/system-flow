package license

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// The embedded copy is the repository's LICENSE.md, byte for byte: the binary
// must carry the terms it is distributed under, not an older draft of them.
func TestEmbeddedLicenseMatchesRepository(t *testing.T) {
	root := filepath.Join("..", "..", "..", "LICENSE.md")
	want, err := os.ReadFile(root)
	if err != nil {
		t.Fatalf("read %s: %v", root, err)
	}
	if string(want) != Text() {
		t.Fatalf("flai/internal/license/LICENSE.md differs from LICENSE.md at the repository root; copy the root file over it")
	}
}

func TestNameIsTheFirstHeading(t *testing.T) {
	if Name() != "Bytepunx Shield License 1.0" {
		t.Fatalf("name %q", Name())
	}
	if !strings.HasPrefix(Text(), "# "+Name()+"\n") {
		t.Fatalf("text does not open with the name as its heading")
	}
}
