package hostapi

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/channel"
	"github.com/bytepunx/system-flow/flai/internal/designsystem"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
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

// worktree gives story a worktree as git links one, its manifest the
// harbour's, and answers its root.
func worktree(t *testing.T, p channel.Project, story string) string {
	t.Helper()
	repo, err := workitem.Open(p.Root)
	if err != nil {
		t.Fatal(err)
	}
	wt := repo.WorktreePath(story)
	manifest, err := os.ReadFile(filepath.Join(p.Root, "system-flow.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	writeFile(t, wt, ".git", []byte("gitdir: "+filepath.Join(p.Root, ".git", "worktrees", story)+"\n"))
	writeFile(t, wt, "system-flow.yaml", manifest)
	return wt
}

// S-0384: with a story, design.get reads the design system in its
// worktree, not the main checkout's; no params, or null, still read main's.
func TestDesignGetWithAStoryReadsItsWorktree(t *testing.T) {
	p := harbour(t)
	tokens, err := os.ReadFile("../designsystem/testdata/full.tokens.json")
	if err != nil {
		t.Fatal(err)
	}
	wt := worktree(t, p, "S-0001")
	writeFile(t, wt, tokenFile, tokens)
	writeFile(t, wt, guidanceFile, []byte("# Design system\n\n## Proposed\n"))

	var got designsystem.Shown
	if err := call(t, p, "design.get", `{"story":"S-1"}`, &got); err != nil {
		t.Fatal(err)
	}
	if !got.Exists || got.Path != tokenFile || len(got.Themes) != 2 || strings.Join(got.Headings, "|") != "Design system|Proposed" {
		t.Errorf("design.get of S-0001 = %+v, want its worktree's token file and guidance", got)
	}
	for _, params := range []string{``, `null`, `{}`} {
		var main designsystem.Shown
		if err := call(t, p, "design.get", params, &main); err != nil {
			t.Fatalf("design.get with %q: %+v", params, err)
		}
		if main.Exists {
			t.Errorf("design.get with %q = %+v, want the main checkout's, which has no token file", params, main)
		}
	}
}

// S-0384: a story with no design system of its own to propose is refused,
// saying why, with no path on the host.
func TestDesignGetRefusesAStoryWithNothingToPropose(t *testing.T) {
	p := harbour(t)
	worktree(t, p, "S-0003")
	cases := []struct {
		name, params string
		code         int
		want         string
	}{
		{"an ID no item has", `{"story":"S-0099"}`, NotFound, "S-0099"},
		{"what is not an ID", `{"story":"../S-0001"}`, channel.CodeInvalidParams, "not an item ID"},
		{"a story that is not a string", `{"story":1}`, channel.CodeInvalidParams, "params"},
		{"a task", `{"story":"T-0001"}`, channel.CodeInvalidParams, "it is a task"},
		{"a done story", `{"story":"S-0003"}`, channel.CodeInvalidParams, "S-0003: it is done"},
		{"a story without a worktree", `{"story":"S-0002"}`, channel.CodeInvalidParams, "S-0002: it has no worktree at .flai-cache/worktrees/S-0002"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var got designsystem.Shown
			err := call(t, p, "design.get", c.params, &got)
			if err == nil || err.Code != c.code || !strings.Contains(err.Message, c.want) {
				t.Fatalf("design.get %s = %+v, %+v; want code %d saying %q", c.params, got, err, c.code, c.want)
			}
			if strings.Contains(err.Message, p.Root) {
				t.Errorf("the refusal names a path on the host: %s", err.Message)
			}
		})
	}
}

// S-0384: a project with no story changing its design system proposes
// nothing, an empty list on the wire.
func TestDesignProposalsNoneIsAnEmptyList(t *testing.T) {
	p := harbour(t)
	worktree(t, p, "S-0001")
	var raw json.RawMessage
	if err := call(t, p, "design.proposals", `{}`, &raw); err != nil {
		t.Fatal(err)
	}
	if string(raw) != "[]" {
		t.Errorf("design.proposals = %s, want []", raw)
	}
}
