package context

import (
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"testing"
)

// fixture loads the documents of testdata/repo.
func fixture(t *testing.T) *Selection {
	t.Helper()
	root := filepath.Join("testdata", "repo")
	docs, err := LoadDocs(root, filepath.Join(root, "design"))
	if err != nil {
		t.Fatal(err)
	}
	return NewSelection(docs)
}

// chosen maps what a selection loads to its reason: a document's path when
// one reason owns all of it, "path § label" for each section otherwise, and
// " (also ...)" when later reasons found it too.
func chosen(s *Selection) map[string]string {
	out := map[string]string{}
	for _, d := range s.Docs {
		own := s.owner[d]
		if len(own) > 0 && own[0] >= 0 && !slices.ContainsFunc(own, func(o int) bool { return o != own[0] }) {
			out[d.Path] = describe(s.choices[own[0]])
			continue
		}
		for i, o := range own {
			if o >= 0 {
				out[d.Path+" § "+d.Label(i)] = describe(s.choices[o])
			}
		}
	}
	return out
}

func describe(c *choice) string {
	if len(c.also) == 0 {
		return c.reason
	}
	return c.reason + " (also " + strings.Join(c.also, "; ") + ")"
}

func TestLoadDocsReadsSystemTechAndADRsWithoutIndexesOrTheTemplate(t *testing.T) {
	s := fixture(t)
	var paths []string
	for _, d := range s.Docs {
		paths = append(paths, d.Path)
	}
	want := []string{"design/system/cli.md", "design/system/dashboard.md", "design/system/plain.md", "design/tech/go.md", "design/tech/node.md"}
	if len(paths) != len(want)+12 || !slices.Equal(paths[:5], want) {
		t.Fatalf("paths %v", paths)
	}
	for _, p := range paths[5:] {
		if !strings.HasPrefix(p, "design/adrs/00") || strings.HasPrefix(p, "design/adrs/0000") {
			t.Errorf("unexpected %s", p)
		}
	}
	first, second, cli := s.byID["ADR-0001"], s.byID["ADR-0002"], s.byPath["design/system/cli.md"]
	if first == nil || second == nil || cli == nil {
		t.Fatal("missing documents")
	}
	if first.Kind != KindADR || !reflect.DeepEqual(first.SupersededBy, []string{"ADR-0003"}) || !reflect.DeepEqual(second.Refines, []string{"ADR-0004"}) {
		t.Errorf("ADR metadata %+v %+v", first, second)
	}
	if cli.Kind != KindSystem || cli.Title != "The CLI" || !reflect.DeepEqual(cli.Topics, []string{"cli"}) || cli.Label(2) != "Rules › Board view" {
		t.Errorf("cli %+v %q", cli, cli.Label(2))
	}
	if node := s.byPath["design/tech/node.md"]; node.Kind != KindTech || len(node.Topics) != 0 {
		t.Errorf("node %+v", node)
	}
}

func TestByTopicsChoosesMatchingSectionsAndWholeADRs(t *testing.T) {
	s := fixture(t)
	s.ByTopics([]string{"cli", "go", "code", "all"})
	want := map[string]string{
		"design/system/cli.md § The CLI": "topics: cli",
		"design/system/cli.md § Rules":   "topics: cli",
		"design/system/cli.md § Prime":   "topics: cli",
		"design/tech/go.md":              "topics: go",
	}
	if got := chosen(s); !reflect.DeepEqual(got, want) {
		t.Errorf("got %v", got)
	}

	s = fixture(t)
	s.ByTopics([]string{"dashboard", "all"})
	want = map[string]string{
		"design/system/cli.md § Rules › Board view":     "topics: dashboard",
		"design/system/dashboard.md § The dashboard":    "topics: dashboard",
		"design/system/dashboard.md § Charts":           "topics: dashboard",
		"design/system/dashboard.md § Charts › Colours": "topics: dashboard",
		"design/adrs/0005-fifth.md":                     "topics: dashboard",
	}
	if got := chosen(s); !reflect.DeepEqual(got, want) {
		t.Errorf("got %v", got)
	}
}

func TestByTopicsLoadsADocumentWholeWhenNoHeadingNarrowsIt(t *testing.T) {
	s := fixture(t)
	s.ByTopics([]string{"go"})
	if !s.Whole(s.byPath["design/tech/go.md"]) || s.Loaded(s.byPath["design/tech/node.md"]) {
		t.Errorf("got %v", chosen(s))
	}
}

func TestLinkedFollowsLinksPathsAndIDsFromStoryEpicAndTasks(t *testing.T) {
	s := fixture(t)
	s.Linked([]Source{
		{ID: "S-0001", Path: "wip/kanban/stories/S-0001-x.md", Body: "See [charts](../../../design/system/dashboard.md#charts) and [the web](https://example.com/design/system/plain.md)."},
		{ID: "E-0001", Path: "wip/kanban/epics/E-0001-x.md", Body: "Built with `design/tech/node.md`."},
		{ID: "T-0001", Path: "wip/kanban/tasks/T-0001-x.md", Body: "Keep to ADR-4, and ADR-0099 which does not exist."},
	})
	want := map[string]string{
		"design/system/dashboard.md § Charts":           "linked from S-0001",
		"design/system/dashboard.md § Charts › Colours": "linked from S-0001",
		"design/tech/node.md":                           "linked from E-0001",
		"design/adrs/0004-fourth.md":                    "linked from T-0001",
	}
	if got := chosen(s); !reflect.DeepEqual(got, want) {
		t.Errorf("got %v", got)
	}
}

func TestLinkedResolvesTheLinksOfAnArchivedItem(t *testing.T) {
	s := fixture(t)
	s.Linked([]Source{{ID: "S-0001", Path: "wip/archive/kanban/stories/S-0001-x.md", Body: "[plain](../../../design/system/plain.md)"}})
	if got := chosen(s); !reflect.DeepEqual(got, map[string]string{"design/system/plain.md": "linked from S-0001"}) {
		t.Errorf("got %v", got)
	}
}

func TestASupersededADRGivesWayToWhatSupersedesIt(t *testing.T) {
	s := fixture(t)
	s.Linked([]Source{{ID: "S-0001", Path: "wip/kanban/stories/S-0001-x.md", Body: "As ADR-0001 says."}})
	if got := chosen(s); !reflect.DeepEqual(got, map[string]string{"design/adrs/0003-third.md": "supersedes ADR-0001"}) {
		t.Errorf("got %v", got)
	}
	if !s.Superseded(s.byID["ADR-0001"]) || s.Superseded(s.byID["ADR-0003"]) {
		t.Error("Superseded")
	}
}

func TestStepFollowsSectionLinksAndRefinesOnce(t *testing.T) {
	s := fixture(t)
	s.ByTopics([]string{"cli"})
	s.Step(StepTopics, StepLinked)
	got := chosen(s)
	if got["design/adrs/0002-second.md"] != "linked from design/system/cli.md § Rules" {
		t.Errorf("section link: %v", got)
	}
	if _, ok := got["design/adrs/0004-fourth.md"]; ok {
		t.Errorf("followed a second step: %v", got)
	}

	s = fixture(t)
	s.Linked([]Source{{ID: "S-0001", Path: "wip/kanban/stories/S-0001-x.md", Body: "ADR-0002"}})
	s.Step(StepTopics, StepLinked)
	if got := chosen(s)["design/adrs/0004-fourth.md"]; got != "refined by ADR-0002 (also linked from design/adrs/0002-second.md § Decision)" {
		t.Errorf("refines: %q", got)
	}
}

func TestStepFollowsOnlyTheStepsItIsGiven(t *testing.T) {
	s := fixture(t)
	s.ByTopics([]string{"cli"})
	s.Step(StepLinked)
	if s.Loaded(s.byID["ADR-0002"]) {
		t.Errorf("followed a topic's links: %v", chosen(s))
	}
}

func TestNothingIsChosenTwiceAndLaterReasonsAreListed(t *testing.T) {
	s := fixture(t)
	s.ByTopics([]string{"go"})
	s.Linked([]Source{
		{ID: "S-0001", Path: "wip/kanban/stories/S-0001-x.md", Body: "[go](../../../design/tech/go.md) and `design/tech/go.md`"},
		{ID: "T-0001", Path: "wip/kanban/tasks/T-0001-x.md", Body: "design/tech/go.md"},
	})
	want := map[string]string{"design/tech/go.md": "topics: go (also linked from S-0001; linked from T-0001)"}
	if got := chosen(s); !reflect.DeepEqual(got, want) {
		t.Errorf("got %v", got)
	}
	if len(s.choices) != 1 {
		t.Errorf("%d choices", len(s.choices))
	}
}

func TestSlug(t *testing.T) {
	for in, want := range map[string]string{"Charts": "charts", "`flai prime --story`": "flai-prime---story", "What a story says, today": "what-a-story-says-today"} {
		if got := Slug(in); got != want {
			t.Errorf("Slug(%q) = %q, want %q", in, got, want)
		}
	}
}
