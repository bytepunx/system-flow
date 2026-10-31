package check

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/issues"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// S-0198: an open issue older than issues.story_after that no open story
// links is warned about, on its first_reported line, naming the threshold and
// the command that makes its story.
func TestIssuesWithNoStory(t *testing.T) {
	root := t.TempDir()
	for _, d := range []string{"design/issues", "wip/kanban/epics", "wip/kanban/stories", "wip/kanban/tasks", "wip/agents", "wip/archive"} {
		_ = os.MkdirAll(filepath.Join(root, d), 0o755)
	}
	issue := "---\nid: %s\ntitle: Issue\nclass: defect\nstatus: %s\ncount: 1\nfirst_reported: %s\nlast_reported: 2026-09-01T12:00:00Z\nupdated: 2026-09-01T12:00:00Z\n---\n\n# %s Issue\n"
	for _, is := range []struct{ id, status, first string }{
		{"I-0001", "open", "2026-08-01T12:00:00Z"},   // old, no story
		{"I-0002", "open", "2026-08-01T12:00:00Z"},   // old, linked by an open story
		{"I-0003", "open", "2026-08-01T12:00:00Z"},   // old, linked only by a done story
		{"I-0004", "closed", "2026-08-01T12:00:00Z"}, // closed
		{"I-0005", "open", "2026-08-30T12:00:00Z"},   // two days old
	} {
		body := fmt.Sprintf(issue, is.id, is.status, is.first, is.id)
		_ = os.WriteFile(filepath.Join(root, "design/issues", is.id+"-issue.md"), []byte(body), 0o644)
	}
	manifest := func(extra string) *workitem.Repo {
		t.Helper()
		_ = os.WriteFile(filepath.Join(root, "system-flow.yaml"), []byte("version: 1\nname: t\nkey: t\nlayout:\n  design: design\n  docs: docs\n  wip: wip\n"+extra), 0o644)
		repo, err := workitem.Open(root)
		if err != nil {
			t.Fatal(err)
		}
		return repo
	}
	repo := manifest("")
	// a story with the status, whose body names the issue; only the status
	// matters to the rule, so it is written rather than moved to
	links := func(title, issue, status string) {
		t.Helper()
		s := mustItem(t, repo, workitem.Story, title, "")
		data, _ := os.ReadFile(s.Path)
		text := strings.Replace(string(data), "\nstatus: backlog\n", "\nstatus: "+status+"\n", 1)
		_ = os.WriteFile(s.Path, []byte(text+"This story remediates "+issue+".\n"), 0o644)
	}
	links("Open", "I-0002", workitem.InProgress)
	links("Done", "I-0003", workitem.Done)

	warned := func(repo *workitem.Repo) []string {
		t.Helper()
		res, err := Run(repo, now)
		if err != nil {
			t.Fatal(err)
		}
		var out []string
		for _, f := range res.Findings {
			if f.Rule == "issues.no-story" {
				if f.Level != Warning || f.Line != 7 {
					t.Errorf("not a warning on the first_reported line: %+v", f)
				}
				out = append(out, filepath.Base(f.Path)[:6])
			}
		}
		return out
	}
	if got := warned(repo); !slices.Equal(got, []string{"I-0001", "I-0003"}) {
		t.Errorf("default threshold: warned %v, want I-0001 and I-0003", got)
	}
	want := "I-0001 has been open since 2026-08-01 with no open story linking it (issues.story_after is unset, so 168h): make one with flai issue story I-0001"
	if msg := strings.Join(rules(t, repo, "issues.no-story"), "\n"); !strings.Contains(msg, want) {
		t.Errorf("message %q does not say %q", msg, want)
	}
	repo = manifest("issues:\n  story_after: 24h\n")
	if got := warned(repo); !slices.Equal(got, []string{"I-0001", "I-0003", "I-0005"}) {
		t.Errorf("24h: warned %v, want I-0001, I-0003, and I-0005", got)
	}
	if msg := strings.Join(rules(t, repo, "issues.no-story"), "\n"); !strings.Contains(msg, "(issues.story_after is 24h)") {
		t.Errorf("the message does not name the setting: %s", msg)
	}
	repo = manifest("issues:\n  story_after: \"0\"\n")
	if got := warned(repo); len(got) != 0 {
		t.Errorf("0 turns the warning off: warned %v", got)
	}
	repo = manifest("issues:\n  story_after: 7d\n")
	if got := warned(repo); len(got) != 0 {
		t.Errorf("an invalid story_after turns the warning off: warned %v", got)
	}
	res, err := Run(repo, now)
	if err != nil {
		t.Fatal(err)
	}
	var invalid []Finding
	for _, f := range res.Findings {
		if f.Rule == "manifest.issues" {
			invalid = append(invalid, f)
		}
	}
	if len(invalid) != 1 || invalid[0].Level != Error || invalid[0].Path != "system-flow.yaml" || invalid[0].Line != 8 || !strings.Contains(invalid[0].Message, `issues.story_after "7d"`) {
		t.Errorf("an invalid story_after is one error on its key: %+v", invalid)
	}
}

// issuesProject is a project with the issues given as id and status in
// design/issues, and no summary.md yet.
func issuesProject(t *testing.T, list ...[2]string) *workitem.Repo {
	t.Helper()
	root := t.TempDir()
	for _, d := range []string{"design/issues", "wip/kanban/epics", "wip/kanban/stories", "wip/kanban/tasks", "wip/agents", "wip/archive"} {
		if err := os.MkdirAll(filepath.Join(root, d), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	files := map[string]string{"system-flow.yaml": "version: 1\nname: t\nkey: t\nlayout:\n  design: design\n  docs: docs\n  wip: wip\nissues:\n  story_after: \"0\"\n"}
	for _, is := range list {
		files["design/issues/"+is[0]+"-issue.md"] = fmt.Sprintf("---\nid: %s\ntitle: Issue %s\nclass: defect\nstatus: %s\ncount: 2\ncost: 15m\nfirst_reported: 2026-08-01T12:00:00Z\nlast_reported: 2026-08-02T12:00:00Z\nupdated: 2026-08-02T12:00:00Z\n---\n\n# %s Issue %s\n", is[0], is[0], is[1], is[0], is[0])
	}
	for name, body := range files {
		if err := os.WriteFile(filepath.Join(root, name), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	return repo
}

// S-0550: summary.md is compared whole with what its issue files render,
// its updated line aside; a summary that differs, or is missing, is one
// issues.summary warning naming flai check --fix and flai issue summary.
// A project that has never had an issue keeps the template's summary.
func TestIssuesSummaryIsComparedWithWhatItsFolderRenders(t *testing.T) {
	repo := issuesProject(t, [2]string{"I-0001", "open"}, [2]string{"I-0002", "closed"})
	summary := filepath.Join(repo.Root, "design", "issues", "summary.md")
	want := []string{"summary.md is missing; run flai check --fix or flai issue summary, which write it"}
	if got := rules(t, repo, "issues.summary"); !slices.Equal(got, want) {
		t.Errorf("missing: %q, want %q", got, want)
	}
	if _, err := issues.WriteSummary(repo, now.Add(-48*time.Hour)); err != nil {
		t.Fatal(err)
	}
	if got := rules(t, repo, "issues.summary"); len(got) != 0 {
		t.Errorf("a summary written two days ago, from the same files, differs only in its updated line: %q", got)
	}
	data, err := os.ReadFile(summary)
	if err != nil {
		t.Fatal(err)
	}
	// the open issue linked, as the old rule asked, but its count stale
	stale := strings.Replace(string(data), "| 2 | 15m | 30m |", "| 1 | 15m | 15m |", 1)
	if stale == string(data) {
		t.Fatalf("the fixture's row is not as expected:\n%s", data)
	}
	if err := os.WriteFile(summary, []byte(stale), 0o644); err != nil {
		t.Fatal(err)
	}
	res, err := Run(repo, now)
	if err != nil {
		t.Fatal(err)
	}
	var found []Finding
	for _, f := range res.Findings {
		if f.Rule == "issues.summary" {
			found = append(found, f)
		}
	}
	if len(found) != 1 || found[0].Level != Warning || found[0].Path != "design/issues/summary.md" || found[0].Line != 1 ||
		found[0].Message != "summary.md differs from what its issue files render; run flai check --fix or flai issue summary, which rewrite it" {
		t.Errorf("a stale count is one warning on summary.md: %+v", found)
	}

	template := "---\ntitle: Open issues\nupdated: 2026-08-01\nstatus: active\n---\n\n# Open issues\n\n| ID | Class | Title | Count | Avg cost | Total cost | Last reported |\n|----|-------|-------|-------|----------|------------|---------------|\n"
	empty := issuesProject(t)
	if err := os.WriteFile(filepath.Join(empty.Root, "design", "issues", "summary.md"), []byte(template), 0o644); err != nil {
		t.Fatal(err)
	}
	if got := rules(t, empty, "issues.summary"); len(got) != 0 {
		t.Errorf("a project with no issues keeps the template's summary: %q", got)
	}
}
