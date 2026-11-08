package check

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// itemsProject is a project with a markdownlint configuration, an epic, and
// two stories under it, each with a forecast whose delivery is not a
// timestamp and a line that ends in a space: what Items is given and what it
// is not, both bad the same way.
func itemsProject(t *testing.T) (repo *workitem.Repo, given, other *workitem.Item) {
	t.Helper()
	repo = overlapProject(t)
	_ = os.WriteFile(filepath.Join(repo.Root, ".markdownlint.yaml"), []byte("default: true\nMD022:\n  lines_below: 0\nMD024:\n  siblings_only: true\nMD025:\n  front_matter_title: \"\"\n"), 0o644)
	_ = os.WriteFile(filepath.Join(repo.Root, ".gitignore"), []byte(".flai-cache/\n"), 0o644)
	epic := mustItem(t, repo, workitem.Epic, "E", "")
	spoil := func(it *workitem.Item) {
		data, err := os.ReadFile(it.Path)
		if err != nil {
			t.Fatal(err)
		}
		doc := strings.Replace(string(data), "\n---\n", "\nforecast:\n  duration: 1h\n  delivery: soon\n  by: t\n  at: 2026-09-01T12:00:00Z\n---\n", 1)
		doc = strings.Replace(doc, "## Notes\n", "## Notes\n\ntrailing \n", 1)
		if err := os.WriteFile(it.Path, []byte(doc), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	given = mustItem(t, repo, workitem.Story, "Given", epic.ID)
	other = mustItem(t, repo, workitem.Story, "Other", epic.ID)
	spoil(given)
	spoil(other)
	return repo, given, other
}

// listed is the project's items as Run lists them, and the one with id.
func listed(t *testing.T, repo *workitem.Repo, id string) ([]*workitem.Item, *workitem.Item) {
	t.Helper()
	all, err := repo.List(true)
	if err != nil {
		t.Fatal(err)
	}
	for _, it := range all {
		if it.ID == id {
			return all, it
		}
	}
	t.Fatalf("%s not listed", id)
	return nil, nil
}

// S-0599: Items reports a bad forecast and a markdown finding on the items it
// is given, and nothing on an item it is not given.
func TestItemsChecksOnlyTheItemsGiven(t *testing.T) {
	repo, given, other := itemsProject(t)
	all, it := listed(t, repo, given.ID)
	res := Items(repo, all, []*workitem.Item{it}, now)
	if res.Items != 1 {
		t.Errorf("Items = %d, want 1", res.Items)
	}
	var forecast, lint bool
	for _, f := range res.Findings {
		if strings.Contains(f.Path, other.ID) {
			t.Errorf("a finding on %s, which Items was not given: %+v", other.ID, f)
		}
		if !strings.Contains(f.Path, given.ID) {
			continue
		}
		forecast = forecast || (f.Rule == "item.front-matter" && f.Level == Error && strings.Contains(f.Message, "forecast.delivery"))
		lint = lint || (f.Rule == "markdown.MD009" && f.Level == Warning)
	}
	if !forecast || !lint {
		t.Errorf("want item.front-matter on forecast.delivery and markdown.MD009 on %s, got %+v", given.ID, res.Findings)
	}
	if res.Errors == 0 || res.Warnings == 0 {
		t.Errorf("Errors = %d, Warnings = %d, want both counted", res.Errors, res.Warnings)
	}
}

// S-0599: on the items given, Items finds what Run finds under the per-item
// rules and the markdown lint. Run's item rules that look across items,
// item.duplicate-id, story.component-tag, and story.unaccepted, are not
// Items' to give.
func TestItemsMatchesRunOnTheItemsGiven(t *testing.T) {
	repo, given, _ := itemsProject(t)
	all, it := listed(t, repo, given.ID)
	var epic *workitem.Item
	for _, x := range all {
		if x.ID == it.Parent {
			epic = x
		}
	}
	items := []*workitem.Item{epic, it}
	paths := map[string]bool{}
	for _, x := range items {
		rel, _ := filepath.Rel(repo.Root, x.Path)
		paths[rel] = true
	}
	across := []string{"item.duplicate-id", "story.component-tag", "story.unaccepted"}
	keep := func(fs []Finding) []string {
		var out []string
		for _, f := range fs {
			per := (strings.HasPrefix(f.Rule, "item.") || strings.HasPrefix(f.Rule, "story.") || strings.HasPrefix(f.Rule, "markdown.")) && !slices.Contains(across, f.Rule)
			if per && paths[f.Path] {
				out = append(out, f.Level+" "+f.Rule+" "+f.Path+" "+f.Message)
			}
		}
		slices.Sort(out)
		return out
	}
	run, err := Run(repo, now)
	if err != nil {
		t.Fatal(err)
	}
	want, got := keep(run.Findings), keep(Items(repo, all, items, now).Findings)
	if len(want) == 0 {
		t.Fatal("Run found nothing on the items given; the fixture should give it findings")
	}
	if !slices.Equal(got, want) {
		t.Errorf("Items findings differ from Run's on the items given:\ngot  %q\nwant %q", got, want)
	}
}
