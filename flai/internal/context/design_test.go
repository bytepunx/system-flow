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

// briefed maps each brief to its reason, with " (also ...)" when later
// reasons found it too, and " [sections]" naming the sections of a design
// file its topics selected when they are not all of it.
func briefed(s *Selection) map[string]string {
	out := map[string]string{}
	for _, b := range s.briefs {
		v := b.reason
		if len(b.also) > 0 {
			v += " (also " + strings.Join(b.also, "; ") + ")"
		}
		if partial(b) {
			var secs []string
			for i := range b.doc.Sections {
				if b.secs[i] {
					secs = append(secs, b.doc.Label(i))
				}
			}
			v += " [" + strings.Join(secs, ", ") + "]"
		}
		out[b.doc.Path] = v
	}
	return out
}

func partial(b *brief) bool {
	for i, sec := range b.doc.Sections {
		if sec.Level > 0 && !b.secs[i] {
			return len(b.secs) > 0
		}
	}
	return false
}

func TestByTopicsBriefsDesignFilesAndADRs(t *testing.T) {
	s := fixture(t)
	s.ByTopics([]string{"cli", "go", "code", "all"})
	want := map[string]string{
		"design/system/cli.md": "topics: cli [The CLI, Rules, Prime]",
		"design/tech/go.md":    "topics: go",
	}
	if got := briefed(s); !reflect.DeepEqual(got, want) {
		t.Errorf("got %v", got)
	}
	if got := chosen(s); len(got) != 0 {
		t.Errorf("topics loaded %v", got)
	}

	s = fixture(t)
	s.ByTopics([]string{"dashboard", "all"})
	want = map[string]string{
		"design/system/cli.md":       "topics: dashboard [Rules › Board view]",
		"design/system/dashboard.md": "topics: dashboard [The dashboard, Charts, Charts › Colours]",
		"design/adrs/0005-fifth.md":  "topics: dashboard",
	}
	if got := briefed(s); !reflect.DeepEqual(got, want) {
		t.Errorf("got %v", got)
	}
}

func TestByTopicsDoesNotBriefWhatIsNamedWhole(t *testing.T) {
	s := fixture(t)
	s.Linked([]Source{{ID: "S-0001", Path: "wip/kanban/stories/S-0001-x.md", Body: "design/tech/go.md and ADR-0005"}})
	s.ByTopics([]string{"go", "dashboard"})
	if got := chosen(s); got["design/tech/go.md"] != "linked from S-0001 (also topics: go)" || got["design/adrs/0005-fifth.md"] != "linked from S-0001 (also topics: dashboard)" {
		t.Errorf("chosen %v", got)
	}
	if s.Briefed(s.byPath["design/tech/go.md"]) || s.Briefed(s.byID["ADR-0005"]) {
		t.Errorf("briefed %v", briefed(s))
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

func TestStepBriefsTheADRsOneStepFromWhatIsNamedOrBriefed(t *testing.T) {
	s := fixture(t)
	s.ByTopics([]string{"cli"})
	s.Step()
	got := briefed(s)
	if got["design/adrs/0002-second.md"] != "linked from design/system/cli.md § Rules" {
		t.Errorf("section link: %v", got)
	}
	if _, ok := got["design/adrs/0004-fourth.md"]; ok {
		t.Errorf("followed a second step: %v", got)
	}

	s = fixture(t)
	s.Linked([]Source{{ID: "S-0001", Path: "wip/kanban/stories/S-0001-x.md", Body: "ADR-0002"}})
	s.Step()
	if got := briefed(s)["design/adrs/0004-fourth.md"]; got != "refined by ADR-0002 (also linked from design/adrs/0002-second.md § Decision)" {
		t.Errorf("refines: %q", got)
	}
	if !s.Whole(s.byID["ADR-0002"]) || s.Loaded(s.byID["ADR-0004"]) {
		t.Errorf("chosen %v", chosen(s))
	}

	s = fixture(t)
	s.Linked([]Source{{ID: "S-0001", Path: "wip/kanban/stories/S-0001-x.md", Body: "ADR-0004"}})
	s.Step()
	if got := briefed(s); !reflect.DeepEqual(got, map[string]string{"design/adrs/0002-second.md": "refines ADR-0004"}) {
		t.Errorf("refined by: %v", got)
	}
}

func TestStepFollowsOnlyTheSectionsTopicsSelected(t *testing.T) {
	s := fixture(t)
	s.ByTopics([]string{"dashboard"})
	s.Step()
	if s.Briefed(s.byID["ADR-0002"]) {
		t.Errorf("followed a section the topics did not select: %v", briefed(s))
	}
}

func TestStepListsItsReasonOnAnADRLoadedWhole(t *testing.T) {
	s := fixture(t)
	s.Linked([]Source{{ID: "S-0001", Path: "wip/kanban/stories/S-0001-x.md", Body: "ADR-0002"}})
	s.ByTopics([]string{"cli"})
	s.Step()
	if got := chosen(s)["design/adrs/0002-second.md"]; got != "linked from S-0001 (also linked from design/system/cli.md § Rules)" {
		t.Errorf("got %q", got)
	}
	if s.Briefed(s.byID["ADR-0002"]) {
		t.Error("briefed an ADR loaded whole")
	}
}

func TestNothingIsChosenTwiceAndLaterReasonsAreListed(t *testing.T) {
	s := fixture(t)
	s.Linked([]Source{
		{ID: "S-0001", Path: "wip/kanban/stories/S-0001-x.md", Body: "[go](../../../design/tech/go.md) and `design/tech/go.md`"},
		{ID: "T-0001", Path: "wip/kanban/tasks/T-0001-x.md", Body: "design/tech/go.md"},
	})
	s.ByTopics([]string{"go"})
	want := map[string]string{"design/tech/go.md": "linked from S-0001 (also linked from T-0001; topics: go)"}
	if got := chosen(s); !reflect.DeepEqual(got, want) {
		t.Errorf("got %v", got)
	}
	if len(s.choices) != 1 || len(s.briefs) != 0 {
		t.Errorf("%d choices, %d briefs", len(s.choices), len(s.briefs))
	}
}

func TestSlug(t *testing.T) {
	for in, want := range map[string]string{"Charts": "charts", "`flai prime --story`": "flai-prime---story", "What a story says, today": "what-a-story-says-today"} {
		if got := Slug(in); got != want {
			t.Errorf("Slug(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestCandidatesRankSectionsNothingLoadedAndEachADROnce(t *testing.T) {
	s := fixture(t)
	s.Linked([]Source{{ID: "S-0001", Path: "wip/kanban/stories/S-0001-x.md", Body: "ADR-0007"}})
	s.ByTopics([]string{"cli"})
	cands := s.Candidates("herons on the estuary at the turn of the tide")
	var got []string
	seen := map[string]bool{}
	for _, c := range cands {
		k := c.doc.Path + " § " + c.doc.Label(c.sec)
		if seen[c.doc.Path] && c.doc.Kind == KindADR {
			t.Errorf("an ADR twice: %s", k)
		}
		seen[c.doc.Path] = true
		got = append(got, k)
	}
	if len(got) == 0 || got[0] != "design/system/plain.md § Heron migration" {
		t.Fatalf("candidates %v", got)
	}
	for _, id := range []string{"ADR-0006", "ADR-0008", "ADR-0009", "ADR-0010", "ADR-0011"} {
		if !seen[s.byID[id].Path] {
			t.Errorf("%s not a candidate: %v", id, got)
		}
	}
	if seen[s.byID["ADR-0007"].Path] || seen[s.byID["ADR-0012"].Path] {
		t.Errorf("a loaded or unrelated ADR is a candidate: %v", got)
	}
}

func TestCandidatesLeaveOutSupersededADRs(t *testing.T) {
	s := fixture(t)
	for _, c := range s.Candidates("the first way") {
		if c.doc.ID == "ADR-0001" {
			t.Errorf("superseded ADR-0001 is a candidate")
		}
	}
}

func TestRankLoadsACandidateAndUndoes(t *testing.T) {
	s := fixture(t)
	var adr Candidate
	for _, c := range s.Candidates("herons estuary") {
		if c.doc.Kind == KindADR {
			adr = c
			break
		}
	}
	undo := s.Rank(adr, 1, true)
	if !s.Whole(adr.doc) || chosen(s)[adr.doc.Path] != "rank 1" {
		t.Errorf("whole: %v", chosen(s))
	}
	undo()
	if s.Loaded(adr.doc) || len(s.choices) != 0 {
		t.Errorf("undo left %v", chosen(s))
	}
	s.Rank(adr, 2, false)
	if s.Whole(adr.doc) || !s.Loaded(adr.doc) {
		t.Errorf("section: %v", chosen(s))
	}
	if adr.RankSize(true) != len(adr.doc.Raw) || adr.RankSize(false) != len(adr.doc.Sections[adr.sec].Text) {
		t.Error("RankSize")
	}
}

func TestQueryIsTitleGoalAndCriteria(t *testing.T) {
	body := "# S-0001 T\n\n## Goal\n\nHerons.\n\n## Acceptance criteria\n- [ ] Tides.\n\n## Notes\n\nCompilers.\n"
	q := Query("Title words", body)
	for _, want := range []string{"Title words", "Herons.", "Tides."} {
		if !strings.Contains(q, want) {
			t.Errorf("missing %q in %q", want, q)
		}
	}
	if strings.Contains(q, "Compilers") {
		t.Errorf("notes in %q", q)
	}
	if q := Query("T", "just text"); !strings.Contains(q, "just text") {
		t.Errorf("no sections: %q", q)
	}
}

func TestItemsAreNamedThenBriefsThenDecisionsThenRanked(t *testing.T) {
	s := fixture(t)
	s.Linked([]Source{{ID: "S-0001", Path: "wip/kanban/stories/S-0001-x.md", Body: "design/tech/go.md and [charts](../../../design/system/dashboard.md#charts)"}})
	s.ByTopics([]string{"cli", "go", "dashboard"})
	s.Step()
	for _, c := range s.Candidates("herons") {
		if c.doc.Path == "design/system/plain.md" {
			s.Rank(c, 1, false)
			break
		}
	}
	items := s.Items()
	var got []string
	for _, it := range items {
		got = append(got, it.Step+" "+it.Path+" § "+it.Label+" = "+it.Reason+" "+strings.Join(it.Also, ";"))
		if it.Size != len(it.Text) {
			t.Errorf("size %d of %d", it.Size, len(it.Text))
		}
	}
	want := []string{
		"named design/system/dashboard.md § Charts = linked from S-0001 ",
		"named design/tech/go.md §  = linked from S-0001 topics: go",
		"briefed design/system/cli.md §  = topics: cli, dashboard ",
		"briefed design/system/dashboard.md §  = topics: dashboard ",
		"briefed design/adrs/0005-fifth.md §  = topics: dashboard ",
		"briefed design/adrs/0002-second.md §  = linked from design/system/cli.md § Rules ",
		"ranked design/system/plain.md § Heron migration = rank 1 ",
	}
	if !slices.Equal(got, want) {
		t.Errorf("got %q", got)
	}
	if !strings.Contains(items[0].Text, "### Colours") || !reflect.DeepEqual(items[0].Heading, []string{"The dashboard", "Charts"}) {
		t.Errorf("named section %+v", items[0])
	}
	if !strings.HasPrefix(items[1].Text, "---\ntitle: Go") || items[1].Heading != nil {
		t.Errorf("whole document %+v", items[1])
	}
	if want := "What the command line does.\n\n- Rules\n  - Board view\n- Prime\n"; items[2].Text != want || items[2].Whole != len(s.byPath["design/system/cli.md"].Raw) {
		t.Errorf("design brief %q", items[2].Text)
	}
	if want := "Cumulative flow and cycle time.\n\n- Charts (selected, loaded)\n  - Colours (selected, loaded)\n- Login\n"; items[3].Text != want {
		t.Errorf("design brief marks %q", items[3].Text)
	}
	if items[5].Text != "Commands are verbs. This refines ADR-0004." || items[5].ID != "ADR-0002" {
		t.Errorf("ADR brief %+v", items[5])
	}
}

func TestPrintedBriefsNameTheirFirstReasonAndCountTheRest(t *testing.T) {
	adr := Item{Path: "design/adrs/0002-second.md", ID: "ADR-0002", Title: "Second", Step: StepBriefed, Reason: "linked from design/system/cli.md § Rules", Also: []string{"refined by ADR-0003", "topics: cli"}, Text: "Commands are verbs.", Whole: 99}
	if got := adr.Printed(); got != "- ADR-0002 Second (design/adrs/0002-second.md; linked from design/system/cli.md, and 2 more): Commands are verbs.\n" {
		t.Errorf("ADR brief %q", got)
	}
	design := Item{Path: "design/system/cli.md", Title: "The CLI", Step: StepBriefed, Reason: "topics: cli", Text: "- Rules\n", Whole: 120}
	if got := design.Printed(); got != "\ndesign/system/cli.md: The CLI (120 bytes; topics: cli)\n\n- Rules\n" {
		t.Errorf("design brief %q", got)
	}
	named := Item{Path: "design/system/cli.md", Label: "Rules", Step: StepNamed, Reason: "linked from S-0001", Also: []string{"topics: cli"}, Text: "## Rules\n"}
	if got := named.Printed(); got != "\ndesign/system/cli.md § Rules\n============================\nreason: linked from S-0001; also topics: cli\n\n## Rules\n" {
		t.Errorf("named %q", got)
	}
}

func TestCatalogListsWhatIsNeitherLoadedNorBriefedAndOutlinesWhatIsInPart(t *testing.T) {
	s := fixture(t)
	s.Linked([]Source{{ID: "S-0001", Path: "wip/kanban/stories/S-0001-x.md", Body: "[charts](../../../design/system/dashboard.md#charts)"}})
	s.ByTopics([]string{"cli"})
	notLoaded, inPart := s.Catalog()
	if len(inPart) != 1 || inPart[0].Path != "design/system/dashboard.md" {
		t.Fatalf("in part %+v", inPart)
	}
	want := []Outline{{1, "The dashboard", false}, {2, "Charts", true}, {3, "Colours", true}, {2, "Login", false}}
	if !reflect.DeepEqual(inPart[0].Outline, want) {
		t.Errorf("outline %+v", inPart[0].Outline)
	}
	if len(notLoaded) != 15 {
		t.Errorf("%d not loaded: %+v", len(notLoaded), notLoaded)
	}
	for _, e := range notLoaded {
		if e.Path == "design/system/cli.md" {
			t.Errorf("a briefed document is in the catalog")
		}
		if e.Path == "design/adrs/0001-first.md" && (!reflect.DeepEqual(e.SupersededBy, []string{"ADR-0003"}) || e.Title != "First") {
			t.Errorf("superseded entry %+v", e)
		}
		if e.Outline != nil {
			t.Errorf("outline on a document not loaded: %+v", e)
		}
	}
}
