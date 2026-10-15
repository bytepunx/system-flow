package hostapi

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/channel"
	"github.com/bytepunx/system-flow/flai/internal/designsystem"
)

// The token file and guidance where the harbour's manifest, with no design
// block, puts them.
const (
	tokenFile    = "design/system/design-system.tokens.json"
	guidanceFile = "design/system/design-system.md"
)

func writeFile(t *testing.T, root, rel string, data []byte) {
	t.Helper()
	full := filepath.Join(root, filepath.FromSlash(rel))
	if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(full, data, 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestDesignGetAnswersWhatDesignShowPrints(t *testing.T) {
	p := harbour(t)
	tokens, err := os.ReadFile("../designsystem/testdata/full.tokens.json")
	if err != nil {
		t.Fatal(err)
	}
	writeFile(t, p.Root, tokenFile, tokens)
	writeFile(t, p.Root, guidanceFile, []byte("---\ntitle: Design system\n---\n\n# Design system\n\n## Colour\n\ntext\n"))

	var got designsystem.Shown
	if err := call(t, p, "design.get", `{}`, &got); err != nil {
		t.Fatal(err)
	}
	if !got.Exists || got.Path != tokenFile || got.Guidance != guidanceFile || len(got.Problems) != 0 ||
		len(got.Themes) != 2 || got.Themes[0].Name != "light" || got.Themes[1].Name != "dark" ||
		!strings.Contains(got.CSS, "--color-primary") || !strings.Contains(got.Tailwind, "@theme") || len(got.Contrast) != 2 {
		t.Errorf("design.get = %+v, want the full token file shown in two themes with no problems", got)
	}
	if strings.Join(got.Headings, "|") != "Design system|Colour" {
		t.Errorf("headings = %q", got.Headings)
	}

	// exactly what flai design show --json prints: designsystem.Show's answer
	want, err := designsystem.Show(os.DirFS(p.Root), tokenFile, guidanceFile)
	if err != nil {
		t.Fatal(err)
	}
	wantJSON, err := json.Marshal(want)
	if err != nil {
		t.Fatal(err)
	}
	var raw json.RawMessage
	_ = call(t, p, "design.get", `{}`, &raw)
	if string(raw) != string(wantJSON) {
		t.Errorf("design.get differs from designsystem.Show:\n%s\nwant\n%s", raw, wantJSON)
	}
}

func TestDesignGetWithoutATokenFileSaysSo(t *testing.T) {
	p := harbour(t)
	var got designsystem.Shown
	if err := call(t, p, "design.get", `{}`, &got); err != nil {
		t.Fatalf("no token file is no error: %+v", err)
	}
	if got.Exists || got.Path != tokenFile || !strings.Contains(got.NoTokenFile, "no token file at "+tokenFile) ||
		!strings.Contains(got.NoTokenFile, "flai design init") || got.CSS != "" || len(got.Themes) != 0 {
		t.Errorf("design.get with no token file = %+v, want exists false and the message", got)
	}
	var raw map[string]any
	_ = call(t, p, "design.get", `{}`, &raw)
	if exists, ok := raw["exists"].(bool); !ok || exists {
		t.Errorf("exists = %v, want false on the wire", raw["exists"])
	}
}

func TestDesignGetListsATokenThatDoesNotParse(t *testing.T) {
	p := harbour(t)
	writeFile(t, p.Root, tokenFile, []byte(`{
  "color": {
    "$type": "color",
    "text": { "$value": { "colorSpace": "oklch", "components": [0.21, 0.02, 265] } },
    "link": { "$value": "{color.nowhere}" }
  }
}
`))
	var got designsystem.Shown
	if err := call(t, p, "design.get", `{}`, &got); err != nil {
		t.Fatalf("a token that does not parse is a problem, not an error: %+v", err)
	}
	if !got.Exists || len(got.Problems) == 0 {
		t.Fatalf("design.get = %+v, want the file's problems", got)
	}
	found := false
	for _, pr := range got.Problems {
		if pr.Path == "color.link" && pr.Line == 5 && strings.Contains(pr.Message, "color.nowhere") {
			found = true
		}
	}
	if !found {
		t.Errorf("problems = %+v, want color.link at line 5 naming color.nowhere", got.Problems)
	}
}

func TestDesignGetAnswersAnErrorWhenTheTokenFileCannotBeRead(t *testing.T) {
	p := harbour(t)
	if err := os.MkdirAll(filepath.Join(p.Root, filepath.FromSlash(tokenFile)), 0o755); err != nil {
		t.Fatal(err)
	}
	var got designsystem.Shown
	if err := call(t, p, "design.get", `{}`, &got); err == nil || err.Code != channel.CodeInternal {
		t.Errorf("design.get on a token file that is a folder = %+v, %+v; want an internal error", got, err)
	}
}
