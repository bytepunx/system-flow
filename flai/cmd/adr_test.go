package cmd

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

const adrTemplate = "---\nid: ADR-0000\ntitle: Template\nstatus: template\ndate: 2026-09-15\nsupersedes: []\nsuperseded_by: []\n---\n\n# ADR-0000 Title\n\n## Context\n\nWhy.\n\n## Decision\n\nWhat.\n\n## Consequences\n\nSo.\n\n## Alternatives considered\n\nOthers.\n"

func adrFile(n, title, status string) string {
	return "---\nid: ADR-" + n + "\ntitle: " + title + "\nstatus: " + status + "\ndate: 2026-09-01\nsupersedes: []\nsuperseded_by: []\n---\n\n# ADR-" + n + " " + title + "\n\n## Context\n\nc\n\n## Decision\n\nd\n\n## Consequences\n\ne\n\n## Alternatives considered\n\nf\n"
}

// adrProject is a git project with ADRs 0001, 0002, and 0007 (a gap), a
// template, and an index, all committed.
func adrProject(t *testing.T) string {
	t.Helper()
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	root := tempProject(t)
	dir := filepath.Join(root, "design", "adrs")
	_ = os.MkdirAll(dir, 0o755)
	w := func(name, body string) { _ = os.WriteFile(filepath.Join(dir, name), []byte(body), 0o644) }
	w("0000-template.md", adrTemplate)
	w("0001-first.md", adrFile("0001", "First", "accepted"))
	w("0002-second.md", adrFile("0002", "Second", "accepted"))
	w("0007-seventh.md", adrFile("0007", "Seventh", "accepted"))
	w("README.md", "---\ntitle: ADRs\nupdated: 2026-09-01\nstatus: active\n---\n\n# Architecture Decision Records\n\n| ADR | Title | Status |\n|-----|-------|--------|\n| [0001](0001-first.md) | First | accepted |\n| [0002](0002-second.md) | Second | accepted |\n| [0007](0007-seventh.md) | Seventh | accepted |\n\nText after the table stays after it.\n")
	_ = os.WriteFile(filepath.Join(root, ".gitignore"), []byte(".flai-cache/\n"), 0o644)
	gitIn(t, root, "init", "-q", "-b", "main")
	gitIn(t, root, "config", "user.email", "olive@example.invalid")
	gitIn(t, root, "config", "user.name", "Olive")
	gitIn(t, root, "add", "-A")
	gitIn(t, root, "commit", "-q", "-m", "init")
	return root
}

const decision = "## Context\n\nA reason.\n\n## Decision\n\nWe do the thing.\n\n## Consequences\n\nIt is done.\n\n## Alternatives considered\n\nNot doing it: lost.\n"

// S-0060: the author writes the decision; the number, file, front matter,
// index row, and superseded_by are supplied.
func TestAdrNew(t *testing.T) {
	root := adrProject(t)
	out, errOut, code := runStdin(t, root, decision, "adr", "new", "Tokens: one per project, never shared", "--status", "accepted", "--supersedes", "2", "--refines", "ADR-0001", "--body-stdin", "--autocommit", "--trailer", "Co-Authored-By: flaiover <flaiover@localhost>", "--json")
	if code != 0 {
		t.Fatalf("adr new: %d %s", code, errOut)
	}
	var res struct {
		ID, Title, Status, Path string
		Number                  int
		Changed                 []string
		Committed               bool
	}
	if err := json.Unmarshal([]byte(out), &res); err != nil {
		t.Fatalf("json: %v\n%s", err, out)
	}
	if res.ID != "ADR-0008" || res.Number != 8 || res.Status != "accepted" || res.Path != "design/adrs/0008-tokens-one-per-project-never-shared.md" || !res.Committed {
		t.Errorf("the number follows the highest file, gaps are not filled: %+v", res)
	}
	file, _ := os.ReadFile(filepath.Join(root, res.Path))
	text := string(file)
	for _, want := range []string{
		"---\nid: ADR-0008\ntitle: \"Tokens: one per project, never shared\"\nstatus: accepted\ndate: 2026-09-19\nsupersedes: [ADR-0002]\nsuperseded_by: []\nrefines: [ADR-0001]\n---\n",
		"\n# ADR-0008 Tokens: one per project, never shared\n\n## Context\n\nA reason.\n",
		"Not doing it: lost.\n",
	} {
		if !strings.Contains(text, want) {
			t.Errorf("missing %q in:\n%s", want, text)
		}
	}
	index, _ := os.ReadFile(filepath.Join(root, "design/adrs/README.md"))
	if !strings.Contains(string(index), "| [0007](0007-seventh.md) | Seventh | accepted |\n| [0008](0008-tokens-one-per-project-never-shared.md) | Tokens: one per project, never shared | accepted, supersedes 0002, refines 0001 |\n\nText after the table stays after it.\n") {
		t.Errorf("the row goes after the last row, with the index's notes:\n%s", index)
	}
	if !strings.Contains(string(index), "| [0002](0002-second.md) | Second | accepted, superseded by 0008 |") {
		t.Errorf("the superseded ADR's row says so:\n%s", index)
	}
	old, _ := os.ReadFile(filepath.Join(root, "design/adrs/0002-second.md"))
	if !strings.Contains(string(old), "superseded_by: [ADR-0008]\n") || !strings.Contains(string(old), "## Decision\n\nd\n") {
		t.Errorf("superseded_by is set and nothing else changes:\n%s", old)
	}
	show := gitIn(t, root, "show", "--stat", "--format=%an|%s|%b", "HEAD")
	if !strings.HasPrefix(show, "Olive|docs: ADR-0008 Tokens: one per project, never shared|Co-Authored-By: flaiover") || !strings.Contains(show, "3 files changed") {
		t.Errorf("one docs commit of the file, the index, and the superseded ADR: %s", show)
	}
	if st := strings.TrimSpace(gitIn(t, root, "status", "--porcelain")); st != "" {
		t.Errorf("clean tree: %q", st)
	}
	// the bare test project has findings of its own; none may be about the ADRs
	if out, _, _ := runIn(t, root, "check"); strings.Contains(out, "adrs/") || strings.Contains(out, "adr.") {
		t.Errorf("flai check has nothing to say about the ADRs:\n%s", out)
	}
}

func TestAdrNewFromTheTemplateThenAccept(t *testing.T) {
	root := adrProject(t)
	out, errOut, code := runIn(t, root, "adr", "new", "A proposal")
	if code != 0 || !strings.HasPrefix(out, "ADR-0008 A proposal (proposed)") {
		t.Fatalf("adr new: %d %s %s", code, out, errOut)
	}
	file := filepath.Join(root, "design/adrs/0008-a-proposal.md")
	data, _ := os.ReadFile(file)
	if !strings.Contains(string(data), "status: proposed\n") || !strings.Contains(string(data), "# ADR-0008 A proposal\n\n## Context\n\nWhy.\n") || strings.Contains(string(data), "refines:") {
		t.Errorf("proposed by default, the body from the project's template:\n%s", data)
	}
	body, _, _ := runIn(t, root, "adr", "new", "--print-body")
	if !strings.HasPrefix(body, "## Context\n\nWhy.\n") || strings.Contains(body, "ADR-0000") {
		t.Errorf("--print-body is the template's sections:\n%s", body)
	}
	gitIn(t, root, "add", "-A")
	gitIn(t, root, "commit", "-q", "-m", "docs: a proposal")

	out, errOut, code = runIn(t, root, "adr", "accept", "ADR-0008", "--autocommit")
	if code != 0 || !strings.Contains(out, "accepted and committed") {
		t.Fatalf("accept: %d %s %s", code, out, errOut)
	}
	data, _ = os.ReadFile(file)
	if !strings.Contains(string(data), "status: accepted\ndate: 2026-09-15\n") {
		t.Errorf("accepted, dated the day it was accepted:\n%s", data)
	}
	index, _ := os.ReadFile(filepath.Join(root, "design/adrs/README.md"))
	if !strings.Contains(string(index), "| A proposal | accepted |") {
		t.Errorf("the index row follows:\n%s", index)
	}
	if _, errOut, code := runIn(t, root, "adr", "accept", "8"); code == 0 || !strings.Contains(errOut, "only a proposed ADR can be accepted") {
		t.Errorf("accepting twice is refused: %d %s", code, errOut)
	}
}

func TestAdrNewRefusals(t *testing.T) {
	root := adrProject(t)
	head := strings.TrimSpace(gitIn(t, root, "rev-parse", "HEAD"))
	for _, c := range []struct {
		args []string
		want string
	}{
		{[]string{"adr", "new", "   "}, "needs a title"},
		{[]string{"adr", "new", "X", "--status", "final"}, "proposed or accepted"},
		{[]string{"adr", "new", "X", "--supersedes", "5"}, "no ADR-0005"},
		{[]string{"adr", "new", "X", "--refines", "seven"}, "not an ADR number"},
		{[]string{"adr", "new", "X", "--supersedes", "1", "--refines", "1"}, "both superseded and refined"},
		{[]string{"adr", "accept", "99"}, "no ADR-0099"},
	} {
		_, errOut, code := runIn(t, root, c.args...)
		if code == 0 || !strings.Contains(errOut, c.want) {
			t.Errorf("flai %v: %d %s, want %q", c.args, code, errOut, c.want)
		}
	}
	if _, errOut, code := runStdin(t, root, " \n", "adr", "new", "X", "--body-stdin"); code == 0 || !strings.Contains(errOut, "standard input is empty") {
		t.Errorf("empty body: %d %s", code, errOut)
	}
	if st := strings.TrimSpace(gitIn(t, root, "status", "--porcelain")); st != "" || strings.TrimSpace(gitIn(t, root, "rev-parse", "HEAD")) != head {
		t.Errorf("refusals leave nothing behind: %q", st)
	}
}

// A failure after files were written puts every one of them back.
func TestAdrNewUndoesEverythingWhenItCannotFinish(t *testing.T) {
	root := adrProject(t)
	// an old ADR with no superseded_by line cannot be marked superseded
	legacy := filepath.Join(root, "design/adrs/0007-seventh.md")
	data, _ := os.ReadFile(legacy)
	_ = os.WriteFile(legacy, []byte(strings.Replace(string(data), "superseded_by: []\n", "", 1)), 0o644)
	gitIn(t, root, "commit", "-q", "-am", "legacy front matter")
	_, errOut, code := runStdin(t, root, decision, "adr", "new", "Replaces two", "--supersedes", "2", "--supersedes", "7", "--body-stdin", "--autocommit")
	if code == 0 || !strings.Contains(errOut, "ADR-0007 has no superseded_by line") {
		t.Fatalf("must fail: %d %s", code, errOut)
	}
	if st := strings.TrimSpace(gitIn(t, root, "status", "--porcelain")); st != "" {
		t.Errorf("the new file is gone, and the index and ADR-0002 are as they were: %q", st)
	}
	// the number was not spent
	if out, _, _ := runIn(t, root, "adr", "new", "Next"); !strings.HasPrefix(out, "ADR-0008 ") {
		t.Errorf("next: %s", out)
	}
}

// flai check keeps the index honest for ADRs made by hand.
func TestCheckAdrIndex(t *testing.T) {
	root := adrProject(t)
	if out, _, _ := runIn(t, root, "check"); strings.Contains(out, "adr.index") {
		t.Fatalf("a consistent index has no finding:\n%s", out)
	}
	_ = os.WriteFile(filepath.Join(root, "design/adrs/0009-by-hand.md"), []byte(adrFile("0009", "By hand", "proposed")), 0o644)
	_ = os.Remove(filepath.Join(root, "design/adrs/0001-first.md"))
	out, _, _ := runIn(t, root, "check")
	if !strings.Contains(out, "design/adrs/0009-by-hand.md:1: warning: adr.index: no row in design/adrs/README.md") {
		t.Errorf("a file with no row:\n%s", out)
	}
	if !strings.Contains(out, "design/adrs/README.md:11: warning: adr.index: the row links to 0001-first.md, which is not there") {
		t.Errorf("a row with no file, at its line:\n%s", out)
	}
	if strings.Contains(out, "0000-template.md") {
		t.Errorf("the template needs no row:\n%s", out)
	}
}

// An accepted ADR is immutable in the editor too; a proposed one is a draft.
func TestAcceptedAdrCannotBeEdited(t *testing.T) {
	root := adrProject(t)
	var doc struct {
		Mode, Reason, Hash, Content string
	}
	out, _, _ := runIn(t, root, "doc", "show", "design/adrs/0001-first.md", "--json")
	_ = json.Unmarshal([]byte(out), &doc)
	if doc.Mode != "none" || !strings.Contains(doc.Reason, "immutable") || !strings.Contains(doc.Reason, "supersedes") {
		t.Errorf("accepted: %+v", doc)
	}
	_, errOut, code := runStdin(t, root, strings.Replace(doc.Content, "## Decision\n\nd\n", "## Decision\n\nrewritten\n", 1), "doc", "save", "design/adrs/0001-first.md", "--hash", doc.Hash)
	if code != 4 || !strings.Contains(errOut, "immutable") {
		t.Errorf("a save is refused: %d %s", code, errOut)
	}
	if _, errOut, code := runIn(t, root, "adr", "new", "A draft"); code != 0 {
		t.Fatal(errOut)
	}
	out, _, _ = runIn(t, root, "doc", "show", "design/adrs/0008-a-draft.md", "--json")
	_ = json.Unmarshal([]byte(out), &doc)
	if doc.Mode != "full" {
		t.Errorf("a proposed ADR is edited like any document: %+v", doc)
	}
	out, _, _ = runIn(t, root, "doc", "show", "design/adrs/0000-template.md", "--json")
	_ = json.Unmarshal([]byte(out), &doc)
	if doc.Mode != "full" {
		t.Errorf("the template is not an accepted ADR: %+v", doc)
	}
}

func TestAdrSlugIsCutAtAWord(t *testing.T) {
	root := adrProject(t)
	out, errOut, code := runIn(t, root, "adr", "new", "The dashboard container cannot write git hooks, git config, or git info: they are mounted read-only over the clone")
	if code != 0 {
		t.Fatal(errOut)
	}
	if !strings.Contains(out, "design/adrs/0008-the-dashboard-container-cannot-write-git-hooks-git-config-or-git-info-they-are.md") {
		t.Errorf("a long title gives a file name that ends on a whole word:\n%s", out)
	}
}
