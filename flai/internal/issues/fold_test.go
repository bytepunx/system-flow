package issues

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// opened opens an issue in r, failing the test when it is refused.
func opened(t *testing.T, r *workitem.Repo, title, cost, story string, at time.Time) *Issue {
	t.Helper()
	is, err := New(r, NewOptions{Title: title, Class: "efficiency", Cost: cost, Story: story, Note: "Found by " + story + ".", Now: at})
	if err != nil {
		t.Fatal(err)
	}
	return is
}

// rel is the issue's file relative to r's root, as Fold names it.
func rel(t *testing.T, r *workitem.Repo, is *Issue) string {
	t.Helper()
	p, err := filepath.Rel(r.Root, is.Path)
	if err != nil {
		t.Fatal(err)
	}
	return filepath.ToSlash(p)
}

// files is the text of each issue file in r, by base name.
func files(t *testing.T, r *workitem.Repo) map[string]string {
	t.Helper()
	matches, _ := filepath.Glob(filepath.Join(Dir(r), "*.md"))
	out := map[string]string{}
	for _, m := range matches {
		data, err := os.ReadFile(m)
		if err != nil {
			t.Fatal(err)
		}
		out[filepath.Base(m)] = string(data)
	}
	return out
}

// I-0112, first instance: the branch opened an issue under the title of one
// the main branch opened meanwhile; the fold moves the branch's into it.
func TestFoldMovesTheBranchsIssueIntoTheMainBranchsOfTheSameTitle(t *testing.T) {
	r := repo(t)
	title := "flai check finds `x` outside the story at close-out"
	onMain := opened(t, r, title, "10m", "S-0001", t0)
	branch := opened(t, r, title, "30m", "S-0002", t0.Add(time.Hour))
	if err := Bump(branch, "S-0002", "20m", "Found again.", t0.Add(2*time.Hour)); err != nil {
		t.Fatal(err)
	}
	folded, paths, err := Fold(r, []string{filepath.Base(onMain.Path)}, t0.Add(3*time.Hour))
	if err != nil {
		t.Fatal(err)
	}
	if want := []Folded{{From: "I-0002", Into: "I-0001"}}; !reflect.DeepEqual(folded, want) {
		t.Errorf("folded %+v, want %+v", folded, want)
	}
	if want := []string{rel(t, r, onMain), rel(t, r, branch)}; !reflect.DeepEqual(paths, want) {
		t.Errorf("paths %q, want the main branch's issue written and the branch's deleted, %q", paths, want)
	}
	if _, err := os.Stat(branch.Path); !os.IsNotExist(err) {
		t.Errorf("the branch's issue file is still there: %v", err)
	}
	got, err := Read(onMain.Path)
	if err != nil {
		t.Fatal(err)
	}
	if err := got.Validate(); err != nil {
		t.Errorf("the folded issue is not valid: %v", err)
	}
	if got.Count != 3 || got.Cost != "20m" || got.Status != "open" {
		t.Errorf("count 3 at (10m+30m+20m)/3 = 20m, open: %+v", got)
	}
	if got.FirstReported != "2026-09-16T10:00:00Z" || got.LastReported != "2026-09-16T12:00:00Z" || got.Updated != "2026-09-16T13:00:00Z" {
		t.Errorf("the earlier first_reported, the later last_reported, updated at the fold: %+v", got)
	}
	parts, _ := splitInstances(got.Body)
	var headings []string
	for _, in := range parts.blocks {
		headings = append(headings, in.heading)
	}
	if want := []string{"### 2026-09-16T10:00:00Z", "### 2026-09-16T11:00:00Z", "### 2026-09-16T12:00:00Z"}; !reflect.DeepEqual(headings, want) {
		t.Errorf("instances %q, want both issues' in timestamp order %q", headings, want)
	}
	if left := files(t, r); len(left) != 1 {
		t.Errorf("one issue file is left, got %d", len(left))
	}
}

// Two issues the branch added under one title fold into the same issue, one
// after the other, and its file is named once.
func TestFoldFoldsSeveralBranchIssuesIntoOne(t *testing.T) {
	r := repo(t)
	onMain := opened(t, r, "Flaky sync", "10m", "S-0001", t0)
	first := opened(t, r, "Flaky sync", "10m", "S-0002", t0.Add(time.Hour))
	second := opened(t, r, "Flaky sync", "10m", "S-0002", t0.Add(2*time.Hour))
	folded, paths, err := Fold(r, []string{filepath.Base(onMain.Path)}, t0.Add(3*time.Hour))
	if err != nil {
		t.Fatal(err)
	}
	if want := []Folded{{From: "I-0002", Into: "I-0001"}, {From: "I-0003", Into: "I-0001"}}; !reflect.DeepEqual(folded, want) {
		t.Errorf("folded %+v, want %+v", folded, want)
	}
	if want := []string{rel(t, r, onMain), rel(t, r, first), rel(t, r, second)}; !reflect.DeepEqual(paths, want) {
		t.Errorf("paths %q, want %q", paths, want)
	}
	got, err := Read(onMain.Path)
	if err != nil {
		t.Fatal(err)
	}
	if got.Count != 3 || got.Cost != "10m" {
		t.Errorf("count 3 at 10m: %+v", got)
	}
}

// A fold leaves alone what it does not fold, and writes nothing.
func TestFoldLeavesAlone(t *testing.T) {
	for _, tc := range []struct {
		name string
		// make opens the issues and returns the base names of the files the
		// main branch has.
		make func(t *testing.T, r *workitem.Repo) []string
	}{
		{"an issue the branch added under a title the main branch has no open issue for", func(t *testing.T, r *workitem.Repo) []string {
			onMain := opened(t, r, "Flaky sync", "10m", "S-0001", t0)
			opened(t, r, "Slow lint", "10m", "S-0002", t0.Add(time.Hour))
			return []string{filepath.Base(onMain.Path)}
		}},
		{"an issue the main branch has too", func(t *testing.T, r *workitem.Repo) []string {
			a := opened(t, r, "Flaky sync", "10m", "S-0001", t0)
			b := opened(t, r, "Flaky sync", "10m", "S-0002", t0.Add(time.Hour))
			return []string{filepath.Base(a.Path), filepath.Base(b.Path)}
		}},
		{"an issue of the title of a closed issue on the main branch", func(t *testing.T, r *workitem.Repo) []string {
			onMain := opened(t, r, "Flaky sync", "10m", "S-0001", t0)
			if err := Close(onMain, "fixed", t0.Add(time.Minute)); err != nil {
				t.Fatal(err)
			}
			opened(t, r, "Flaky sync", "10m", "S-0002", t0.Add(time.Hour))
			return []string{filepath.Base(onMain.Path)}
		}},
		{"a closed issue the branch added", func(t *testing.T, r *workitem.Repo) []string {
			onMain := opened(t, r, "Flaky sync", "10m", "S-0001", t0)
			branch := opened(t, r, "Flaky sync", "10m", "S-0002", t0.Add(time.Hour))
			if err := Close(branch, "a duplicate", t0.Add(2*time.Hour)); err != nil {
				t.Fatal(err)
			}
			return []string{filepath.Base(onMain.Path)}
		}},
		{"nothing in the folder", func(t *testing.T, r *workitem.Repo) []string { return nil }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			r := repo(t)
			onMain := tc.make(t, r)
			before := files(t, r)
			folded, paths, err := Fold(r, onMain, t0.Add(3*time.Hour))
			if err != nil {
				t.Fatal(err)
			}
			if len(folded) != 0 || len(paths) != 0 {
				t.Errorf("folded %+v and changed %q, want nothing", folded, paths)
			}
			if after := files(t, r); !reflect.DeepEqual(after, before) {
				t.Errorf("the issue files changed:\nbefore %v\nafter %v", before, after)
			}
		})
	}
}
