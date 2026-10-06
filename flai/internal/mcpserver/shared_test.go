package mcpserver

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/guard"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// sharedManifest is the fixture's manifest with a component and three shared
// patterns, the last not valid, and a comment in the list.
const sharedManifest = "version: 1\nname: t\nkey: t\nlayout:\n  design: design\n  docs: docs\n  wip: wip\nprojects:\n  - name: cli\n    path: flai\n    kind: go\nclaims:\n  shared:\n    # many stories change these\n    - design/adrs\n    - docs/users/*.md\n    - /abs\n"

// sharedFixture is setup in a session the operator started, with
// sharedManifest written after the server opened the project, so that the
// tools must read the list as it is now.
func sharedFixture(t *testing.T, served string) *fixture {
	t.Helper()
	for _, v := range []string{"FLAI_ROLE", "FLAI_STORY"} {
		t.Setenv(v, "")
	}
	t.Setenv(guard.StartedByEnv, served)
	f := setup(t)
	if err := os.WriteFile(f.manifest(), []byte(sharedManifest), 0o644); err != nil {
		t.Fatal(err)
	}
	return f
}

func (f *fixture) manifest() string { return filepath.Join(f.repo.Root, manifest.File) }

// sharedOut decodes a shared_paths answer.
func sharedOut(t *testing.T, out map[string]any) SharedPathsOut {
	t.Helper()
	var got SharedPathsOut
	data, _ := json.Marshal(out)
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatal(err)
	}
	return got
}

// entries are a shared_paths answer's entries, one "entry path shared
// pattern story" line each.
func entries(es []SharedEntry) string {
	var lines []string
	for _, e := range es {
		lines = append(lines, strings.Join([]string{e.Entry, e.Path, map[bool]string{true: "shared", false: "-"}[e.Shared], e.Pattern, e.Story}, " "))
	}
	return strings.Join(lines, "\n")
}

// S-0295, ADR-0096: shared_paths lists the patterns, the invalid ones with
// why, and checks paths and touches entries, a component read as its path.
func TestSharedPathsListsAndChecksPaths(t *testing.T) {
	f := sharedFixture(t, "")
	out, failed := f.call(t, "shared_paths", map[string]any{})
	if failed != "" {
		t.Fatal(failed)
	}
	got := sharedOut(t, out)
	if strings.Join(got.Shared, ",") != "design/adrs,docs/users/*.md,/abs" || len(got.Entries) != 0 {
		t.Errorf("list: %+v", got)
	}
	if len(got.Invalid) != 1 || got.Invalid[0].Pattern != "/abs" || !strings.Contains(got.Invalid[0].Reason, "absolute") {
		t.Errorf("invalid: %+v", got.Invalid)
	}

	out, failed = f.call(t, "shared_paths", map[string]any{"paths": []string{"design/adrs/0096-x.md", "./design/adrs/", "docs/users/flai.md", "docs/users", "docs", "cli"}})
	if failed != "" {
		t.Fatal(failed)
	}
	want := strings.Join([]string{
		"design/adrs/0096-x.md design/adrs/0096-x.md shared design/adrs ",
		"./design/adrs/ design/adrs shared design/adrs ",
		"docs/users/flai.md docs/users/flai.md shared docs/users/*.md ",
		"docs/users docs/users -  ", // a folder may hold more than *.md
		"docs docs -  ",
		"cli flai -  ",
	}, "\n")
	if e := entries(sharedOut(t, out).Entries); e != want {
		t.Errorf("entries:\n%s\nwant:\n%s", e, want)
	}

	if _, failed := f.call(t, "shared_paths", map[string]any{"paths": []string{" "}}); !strings.Contains(failed, "names no path") {
		t.Errorf("an empty entry: %q", failed)
	}
}

// shared_paths checks a story's claim entry by entry, as its hold reads it:
// its touches, a component as its path, and its open tasks' touches.
func TestSharedPathsChecksAStorysClaim(t *testing.T) {
	f := sharedFixture(t, "")
	story, err := f.repo.Create(workitem.NewOptions{Type: workitem.Story, Title: "Shared", Owner: "alex", Touches: []string{"design/adrs", "cli"}, Now: t0})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.repo.Create(workitem.NewOptions{Type: workitem.Task, Title: "Docs", Parent: story.ID, Owner: "alex", Touches: []string{"docs/users/flai.md"}, Now: t0}); err != nil {
		t.Fatal(err)
	}
	padded := strings.Replace(story.ID, "S-0", "s-", 1)
	out, failed := f.call(t, "shared_paths", map[string]any{"story": padded, "paths": []string{"docs"}})
	if failed != "" {
		t.Fatal(failed)
	}
	want := strings.Join([]string{
		"docs docs -  ",
		"design/adrs design/adrs shared design/adrs " + story.ID,
		"flai flai -  " + story.ID,
		"docs/users/flai.md docs/users/flai.md shared docs/users/*.md " + story.ID,
	}, "\n")
	if e := entries(sharedOut(t, out).Entries); e != want {
		t.Errorf("entries:\n%s\nwant:\n%s", e, want)
	}
	for id, says := range map[string]string{"E-0001": "is not a story", "S-0999": "S-0999"} {
		if _, failed := f.call(t, "shared_paths", map[string]any{"story": id}); !strings.Contains(failed, says) {
			t.Errorf("story %s: %q", id, failed)
		}
	}
}

// shared_paths_edit adds and removes patterns, rewriting the list alone and
// keeping its comment, and returns what changed and the list after.
func TestSharedPathsEditAddsAndRemoves(t *testing.T) {
	f := sharedFixture(t, "")
	out, failed := f.call(t, "shared_paths_edit", map[string]any{"add": []string{"template/CHANGELOG.md", "design/**/README.md"}})
	if failed != "" {
		t.Fatal(failed)
	}
	if a, s := strings.Join(toStrings(out["added"]), ","), strings.Join(toStrings(out["shared"]), ","); a != "template/CHANGELOG.md,design/**/README.md" || s != "design/adrs,docs/users/*.md,/abs,template/CHANGELOG.md,design/**/README.md" || out["removed"] != nil {
		t.Errorf("add: %v", out)
	}

	out, failed = f.call(t, "shared_paths_edit", map[string]any{"remove": []string{"/abs", "design/adrs"}, "add": []string{"design/adrs/**"}})
	if failed != "" {
		t.Fatal(failed)
	}
	if r, a, s := strings.Join(toStrings(out["removed"]), ","), strings.Join(toStrings(out["added"]), ","), strings.Join(toStrings(out["shared"]), ","); r != "/abs,design/adrs" || a != "design/adrs/**" || s != "docs/users/*.md,template/CHANGELOG.md,design/**/README.md,design/adrs/**" {
		t.Errorf("remove and add: %v", out)
	}

	data, _ := os.ReadFile(f.manifest())
	want := strings.Replace(sharedManifest, "    - design/adrs\n    - docs/users/*.md\n    - /abs\n", "    - docs/users/*.md\n    - template/CHANGELOG.md\n    - \"design/**/README.md\"\n    - \"design/adrs/**\"\n", 1)
	if string(data) != want {
		t.Errorf("the manifest:\n%s\nwant:\n%s", data, want)
	}
	listed, _ := f.call(t, "shared_paths", map[string]any{})
	if s := strings.Join(sharedOut(t, listed).Shared, ","); s != "docs/users/*.md,template/CHANGELOG.md,design/**/README.md,design/adrs/**" {
		t.Errorf("shared_paths after the edit: %s", s)
	}
}

// A pattern that is not valid, one already listed, or one to remove that is
// not listed is refused with flai's reason, and nothing is written, not even
// what the call's other half would have changed.
func TestSharedPathsEditRefusesWhatIsNotValid(t *testing.T) {
	f := sharedFixture(t, "")
	for _, c := range []struct {
		args map[string]any
		says string
	}{
		{map[string]any{"add": []string{"../outside"}}, ".. segment"},
		{map[string]any{"add": []string{"docs/[x"}}, "not a valid glob"},
		{map[string]any{"add": []string{""}}, "is empty"},
		{map[string]any{"add": []string{"design/adrs"}}, "in the list already"},
		{map[string]any{"remove": []string{"design/adrs"}, "add": []string{"/root"}}, "absolute path"},
		{map[string]any{"remove": []string{"docs/users"}}, "not in the list"},
		{map[string]any{}, "give add, remove, or both"},
	} {
		_, failed := f.call(t, "shared_paths_edit", c.args)
		if !strings.Contains(failed, c.says) {
			t.Errorf("%v: %q, want it to say %q", c.args, failed, c.says)
		}
		if data, _ := os.ReadFile(f.manifest()); string(data) != sharedManifest {
			t.Errorf("%v changed the manifest:\n%s", c.args, data)
		}
	}
}

// In a session flai serve started, shared_paths_edit refuses itself, should
// flai guard not run, and says to ask the operator; shared_paths still reads.
func TestSharedPathsEditIsTheOperatorsOnly(t *testing.T) {
	f := sharedFixture(t, guard.StartedByServe)
	_, failed := f.call(t, "shared_paths_edit", map[string]any{"add": []string{"template"}})
	if !strings.Contains(failed, "operator's (ADR-0096)") || !strings.Contains(failed, "thread_open") {
		t.Errorf("refusal: %q", failed)
	}
	if data, _ := os.ReadFile(f.manifest()); string(data) != sharedManifest {
		t.Errorf("the manifest changed:\n%s", data)
	}
	if _, failed := f.call(t, "shared_paths", map[string]any{"paths": []string{"design/adrs"}}); failed != "" {
		t.Errorf("shared_paths: %s", failed)
	}
}
