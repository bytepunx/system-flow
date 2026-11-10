package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// streamStateProject is a project with the template's markdown lint and
// S-0001 in progress with a narrative; S-0002 is in progress with none, and
// S-0003 is ready.
func streamStateProject(t *testing.T) (root, narrative string) {
	t.Helper()
	useConfig(t, filepath.Join(t.TempDir(), "cfg.json"))
	t.Setenv("FLAI_AGENT", "tester")
	t.Setenv("FLAI_SESSION", "sess-1")
	root = tempProject(t)
	_ = os.WriteFile(filepath.Join(root, ".markdownlint.yaml"), []byte("default: true\nMD013: false\nMD022:\n  lines_below: 0\nMD024:\n  siblings_only: true\nMD025:\n  front_matter_title: \"\"\nMD032: false\nMD041: false\n"), 0o644)
	steps := [][]string{{"epic", "new", "Epic"}}
	for _, title := range []string{"One", "Two", "Three"} {
		steps = append(steps, []string{"story", "new", title, "--epic", "E-0001"})
	}
	steps = append(steps, []string{"move", "S-0001", "ready"}, []string{"move", "S-0001", "in-progress"}, []string{"stream", "open", "S-0001", "--no-branch"},
		[]string{"move", "S-0002", "ready"}, []string{"move", "S-0002", "in-progress"}, []string{"move", "S-0003", "ready"})
	for _, args := range steps {
		if len(args) > 1 && args[0] == "move" && args[2] == "ready" {
			matches, _ := filepath.Glob(filepath.Join(root, "wip/kanban/stories", args[1]+"-*.md"))
			s := read(t, matches[0])
			_ = os.WriteFile(matches[0], []byte(strings.Replace(s, "## Acceptance criteria\n- [ ]\n", "## Acceptance criteria\n- [ ] it works\n", 1)), 0o644)
		}
		if _, errOut, code := runIn(t, root, args...); code != 0 {
			t.Fatalf("flai %v: %s", args, errOut)
		}
	}
	return root, filepath.Join(root, "wip/agents/S-0001.md")
}

// section is what is under a narrative's ## name, up to the next ## heading.
func section(narrative, name string) string {
	_, after, _ := strings.Cut(narrative, "\n## "+name+"\n")
	if i := strings.Index(after, "\n## "); i >= 0 {
		after = after[:i]
	}
	return strings.TrimSpace(after)
}

// S-0271: flai stream state replaces both sections, says so in a line,
// appends nothing to the log, and writes the same text again as nothing.
func TestStreamStateWritesBothSections(t *testing.T) {
	root, path := streamStateProject(t)
	before := read(t, path)
	out, errOut, code := runIn(t, root, "stream", "state", "S-0001", "--current", "T-1 is done.", "--next", "1. Do T-2.")
	if code != 0 || out != "S-0001: wrote Current state and Next steps in wip/agents/S-0001.md at 2026-09-15T21:00:00Z\n" {
		t.Fatalf("state: %d %q %s", code, out, errOut)
	}
	after := read(t, path)
	if section(after, "Current state") != "T-1 is done." || section(after, "Next steps") != "1. Do T-2." {
		t.Errorf("sections:\n%s", after)
	}
	for _, name := range []string{"Context", "Decisions", "Open questions", "Log"} {
		if section(after, name) != section(before, name) {
			t.Errorf("## %s changed:\n%s", name, after)
		}
	}
	if !strings.Contains(after, "agent: tester") || !strings.Contains(after, "session: sess-1") {
		t.Errorf("front matter does not record the writer:\n%s", after)
	}
	out, errOut, code = runIn(t, root, "stream", "state", "S-0001", "--current", "T-1 is done.", "--next", "1. Do T-2.")
	if code != 0 || out != "S-0001: unchanged; wip/agents/S-0001.md already holds the Current state and Next steps given\n" {
		t.Errorf("the same text again: %d %q %s", code, out, errOut)
	}
}

// Either flag alone writes its section and leaves the other; --json prints
// the result.
func TestStreamStateEachFlagAloneAndJSON(t *testing.T) {
	root, path := streamStateProject(t)
	if _, errOut, code := runIn(t, root, "stream", "state", "S-0001", "--current", "Old.", "--next", "1. Old."); code != 0 {
		t.Fatal(errOut)
	}
	out, errOut, code := runIn(t, root, "stream", "state", "S-0001", "--next", "1. New.", "--json")
	var v workitem.StreamState
	if code != 0 || json.Unmarshal([]byte(out), &v) != nil {
		t.Fatalf("--next --json: %d %s %s", code, out, errOut)
	}
	if v.Stream != "S-0001" || v.Path != "wip/agents/S-0001.md" || strings.Join(v.Written, "|") != "Next steps" || !v.Changed || v.Updated != "2026-09-15T21:00:00Z" {
		t.Errorf("--json: %+v", v)
	}
	if s := read(t, path); section(s, "Current state") != "Old." || section(s, "Next steps") != "1. New." {
		t.Errorf("--next alone:\n%s", s)
	}
	out, errOut, code = runIn(t, root, "stream", "state", "S-0001", "--current", "New.")
	if code != 0 || out != "S-0001: wrote Current state in wip/agents/S-0001.md at 2026-09-15T21:00:00Z\n" {
		t.Fatalf("--current alone: %d %q %s", code, out, errOut)
	}
	if s := read(t, path); section(s, "Current state") != "New." || section(s, "Next steps") != "1. New." {
		t.Errorf("--current alone:\n%s", s)
	}
	out, _, _ = runIn(t, root, "stream", "state", "S-0001", "--current", "New.", "--json")
	if !strings.Contains(out, `"changed": false`) || !strings.Contains(out, `"written": [`) {
		t.Errorf("unchanged --json: %s", out)
	}
}

// A value of - reads that text from standard input; only one flag can be.
func TestStreamStateReadsStandardInput(t *testing.T) {
	root, path := streamStateProject(t)
	out, errOut, code := runStdin(t, root, "1. First.\n2. Second.\n\nThen review.\n", "stream", "state", "S-0001", "--current", "Working.", "--next", "-")
	if code != 0 || !strings.HasPrefix(out, "S-0001: wrote Current state and Next steps in ") {
		t.Fatalf("--next -: %d %q %s", code, out, errOut)
	}
	if s := read(t, path); section(s, "Next steps") != "1. First.\n2. Second.\n\nThen review." || section(s, "Current state") != "Working." {
		t.Errorf("standard input:\n%s", s)
	}
	before := read(t, path)
	if _, errOut, code := runStdin(t, root, "x\n", "stream", "state", "S-0001", "--current", "-", "--next", "-"); code != 1 || !strings.Contains(errOut, "only one of --current, --next, and --context") {
		t.Errorf("both -: %d %s", code, errOut)
	}
	if _, errOut, code := runStdin(t, root, "x\n", "stream", "state", "S-0001", "--context", "-", "--next", "-"); code != 1 || !strings.Contains(errOut, "only one of --current, --next, and --context") {
		t.Errorf("--context - and --next -: %d %s", code, errOut)
	}
	if read(t, path) != before {
		t.Error("a refused write changed the narrative")
	}
	if _, errOut, code := runStdin(t, root, "Why this story.\n\nWhat it builds on.\n", "stream", "state", "S-0001", "--context", "-"); code != 0 {
		t.Fatalf("--context -: %d %s", code, errOut)
	}
	if s := read(t, path); section(s, "Context") != "Why this story.\n\nWhat it builds on." {
		t.Errorf("--context from standard input:\n%s", s)
	}
}

// None of the five flags is a usage error; a story not in progress or review, one with
// no narrative, and text the lint rejects are refused with exit 4 and write
// nothing.
func TestStreamStateRefuses(t *testing.T) {
	root, path := streamStateProject(t)
	before := read(t, path)
	if _, errOut, code := runIn(t, root, "stream", "state", "S-0001"); code != 1 || !strings.Contains(errOut, "give --current, --next, --context, --decision, or --question") {
		t.Errorf("neither flag: %d %s", code, errOut)
	}
	type refusal struct {
		Refused struct {
			ID       string `json:"id"`
			Status   string `json:"status"`
			Reason   string `json:"reason"`
			Path     string `json:"path"`
			Findings []struct {
				Rule string `json:"rule"`
			} `json:"findings"`
		} `json:"refused"`
	}
	if _, errOut, code := runIn(t, root, "stream", "state", "S-0003", "--current", "Now."); code != exitDocRefused || !strings.Contains(errOut, "S-0003 is ready") {
		t.Errorf("a ready story: %d %s", code, errOut)
	}
	out, _, code := runIn(t, root, "stream", "state", "S-0003", "--current", "Now.", "--json")
	var r refusal
	if code != exitDocRefused || json.Unmarshal([]byte(out), &r) != nil || r.Refused.ID != "S-0003" || r.Refused.Status != "ready" || !strings.Contains(r.Refused.Reason, "in progress or in review") {
		t.Errorf("a ready story --json: %d %s", code, out)
	}
	out, errOut, code := runIn(t, root, "stream", "state", "S-0002", "--next", "1. Then.", "--json")
	r = refusal{}
	if code != exitDocRefused || json.Unmarshal([]byte(out), &r) != nil || r.Refused.ID != "S-0002" || r.Refused.Status != "in-progress" || !strings.Contains(errOut, "flai stream open S-0002") {
		t.Errorf("no narrative: %d %s %s", code, out, errOut)
	}
	out, errOut, code = runIn(t, root, "stream", "state", "S-0001", "--next", "1. Do it.\n\n**Looks like a heading**", "--json")
	r = refusal{}
	if code != exitDocRefused || json.Unmarshal([]byte(out), &r) != nil || r.Refused.Path != "wip/agents/S-0001.md" || len(r.Refused.Findings) == 0 || r.Refused.Findings[0].Rule != "MD036" || !strings.Contains(errOut, "markdown lint") {
		t.Errorf("text the lint rejects: %d %s %s", code, out, errOut)
	}
	if _, errOut, code := runIn(t, root, "stream", "state", "S-0001", "--current", "Fine.\n\n## Decisions"); code != 1 || !strings.Contains(errOut, "use ### or deeper") {
		t.Errorf("a heading in the text: %d %s", code, errOut)
	}
	if read(t, path) != before {
		t.Error("a refused write changed the narrative")
	}
}

// S-0668: --context replaces ## Context, each --decision appends a bullet
// under ## Decisions, and each --question a dated bullet under ## Open
// questions, above the threads block; each alone leaves every other section
// as it was, and --json's written names the sections in template order.
func TestStreamStateContextDecisionsAndQuestions(t *testing.T) {
	root, path := streamStateProject(t)
	n := read(t, path)
	n = strings.Replace(n, "## Open questions\n", "## Open questions\n\n<!-- threads:start -->\n- TH-0001 Which port? (open)\n<!-- threads:end -->\n", 1)
	if err := os.WriteFile(path, []byte(n), 0o644); err != nil {
		t.Fatal(err)
	}
	block := "<!-- threads:start -->\n- TH-0001 Which port? (open)\n<!-- threads:end -->"
	others := func(before, after string, except string) {
		t.Helper()
		for _, name := range []string{"Context", "Current state", "Next steps", "Decisions", "Open questions", "Log"} {
			if name != except && section(after, name) != section(before, name) {
				t.Errorf("## %s changed:\n%s", name, after)
			}
		}
	}

	before := read(t, path)
	out, errOut, code := runIn(t, root, "stream", "state", "S-0001", "--context", "Built on S-0271.")
	if code != 0 || out != "S-0001: wrote Context in wip/agents/S-0001.md at 2026-09-15T21:00:00Z\n" {
		t.Fatalf("--context alone: %d %q %s", code, out, errOut)
	}
	after := read(t, path)
	if section(after, "Context") != "Built on S-0271." {
		t.Errorf("--context:\n%s", after)
	}
	others(before, after, "Context")

	before = after
	if _, errOut, code := runIn(t, root, "stream", "state", "S-0001", "--decision", "Kept the flag, as two callers use it.", "--decision", "T-3 waits for T-2."); code != 0 {
		t.Fatalf("--decision twice: %d %s", code, errOut)
	}
	after = read(t, path)
	if got := section(after, "Decisions"); got != "- Kept the flag, as two callers use it.\n- T-3 waits for T-2." {
		t.Errorf("--decision twice: %q", got)
	}
	others(before, after, "Decisions")
	if _, errOut, code := runIn(t, root, "stream", "state", "S-0001", "--decision", "T-3 waits for T-2."); code != 0 {
		t.Fatalf("--decision again: %d %s", code, errOut)
	}
	if got := section(read(t, path), "Decisions"); strings.Count(got, "- T-3 waits for T-2.") != 2 {
		t.Errorf("the same decision again is appended again: %q", got)
	}

	before = read(t, path)
	out, errOut, code = runIn(t, root, "stream", "state", "S-0001", "--question", "Should the port be configurable, or fixed?", "--json")
	var v workitem.StreamState
	if code != 0 || json.Unmarshal([]byte(out), &v) != nil {
		t.Fatalf("--question --json: %d %s %s", code, out, errOut)
	}
	if strings.Join(v.Written, "|") != "Open questions" || !v.Changed {
		t.Errorf("--question --json: %+v", v)
	}
	after = read(t, path)
	if got, want := section(after, "Open questions"), "- 2026-09-15: Should the port be configurable, or fixed?\n\n"+block; got != want {
		t.Errorf("--question:\n got %q\nwant %q", got, want)
	}
	others(before, after, "Open questions")

	out, errOut, code = runIn(t, root, "stream", "state", "S-0001", "--context", "Built on S-0271.", "--next", "1. Ask.", "--decision", "One more.", "--question", "Another?", "--json")
	v = workitem.StreamState{}
	if code != 0 || json.Unmarshal([]byte(out), &v) != nil {
		t.Fatalf("every flag --json: %d %s %s", code, out, errOut)
	}
	if got := strings.Join(v.Written, "|"); got != "Context|Next steps|Decisions|Open questions" {
		t.Errorf("written in template order: %q", got)
	}
}
