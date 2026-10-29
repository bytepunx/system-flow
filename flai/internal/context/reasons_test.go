package context

import (
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/conventions"
	"github.com/bytepunx/system-flow/flai/internal/threads"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// reasonsConvention is the orchestrator's convention in reasonsRepo, its
// sections tagged as ADR-0217 §2 tags them: one for every pack, one each
// for reviews, threads, and the role, one for the analyzer beside all, and
// one without all for a topic the orchestrator lacks.
const reasonsConvention = "---\ntitle: Orchestration\nupdated: 2026-10-29\naudience: agent\norder: 50\nstatus: active\nroles: [orchestrate]\n---\n\n# Orchestration\n\n" +
	"## Every run\n\n- Every run rule.\n\n" +
	"## Reviews <!-- topics: all, review -->\n\n- Review rule.\n\n" +
	"## Threads <!-- topics: all, thread -->\n\n- Thread rule.\n\n" +
	"## Orchestrating <!-- topics: all, orchestration -->\n\n- Orchestrating rule.\n\n" +
	"## Analysis <!-- topics: all, analysis -->\n\n- Analysis rule.\n\n" +
	"## Svelte <!-- topics: svelte -->\n\n- Svelte rule.\n\n" +
	"<!-- system-flow:end-of-baseline -->\n\n## Project additions\n"

// reasonsRepo is strategicRepo with the orchestrator's convention tagged by
// reason, a second story under the epic and a cancelled third, a thread
// opened by a story's agent on the first story, one on a design document,
// the planner's open thread on the epic, its resolved one, and an agent's
// open one on the epic.
func reasonsRepo(t *testing.T) (*workitem.Repo, *workitem.Item, *workitem.Item) {
	t.Helper()
	repo, epic, story := strategicRepo(t)
	writeRepo(t, repo, "design/conventions/orchestration.md", reasonsConvention)
	now := time.Date(2026, 10, 29, 9, 0, 0, 0, time.UTC)
	if _, err := repo.Create(workitem.NewOptions{Type: workitem.Story, Title: "Second story", Parent: epic.ID, Owner: "alex", Now: now}); err != nil {
		t.Fatal(err)
	}
	gone, err := repo.Create(workitem.NewOptions{Type: workitem.Story, Title: "Cancelled story", Parent: epic.ID, Owner: "alex", Now: now})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := repo.Transition(gone, workitem.Cancelled, "alex", "not needed", now); err != nil {
		t.Fatal(err)
	}
	open := func(title, on, by string) *threads.Thread {
		th, err := threads.New(repo, threads.NewOptions{Title: title, On: on, Author: by, Text: title + " Asked in full.", Now: now})
		if err != nil {
			t.Fatal(err)
		}
		return th
	}
	// TH-0001 to TH-0005, in this order.
	open("Which way round?", story.ID, "agent-S-0001")
	open("Is the design right?", "design/system/orchestration.md", "agent-S-0001")
	open("Split the epic?", epic.ID, "planner-E-0001")
	resolved := open("Rename the epic?", epic.ID, "planner-E-0001")
	if _, err := threads.Resolve(repo, resolved.ID, "planner-E-0001", "settled", now); err != nil {
		t.Fatal(err)
	}
	open("Who owns the epic?", epic.ID, "agent-S-0002")
	return repo, epic, story
}

// reasonItems is the reason items of the pack by path, each with its
// reason.
func reasonItems(p *Pack) map[string]Item {
	out := map[string]Item{}
	for _, it := range p.Items {
		if it.Step == StepReason {
			out[it.Path] = it
		}
	}
	return out
}

// keptHeadings is the headings of the sections kept in the pack's
// convention named name.
func keptHeadings(p *Pack, name string) []string {
	var out []string
	for _, c := range p.Conventions {
		if c.Name != name {
			continue
		}
		for _, k := range c.Kept {
			if k.Heading != "" {
				out = append(out, k.Heading)
			}
		}
	}
	return out
}

// mustReasons is ForReasons for reasons with the project's budget, which
// the test fails on refusing.
func mustReasons(t *testing.T, repo *workitem.Repo, reasons ...string) *Pack {
	t.Helper()
	p, err := ForReasons(repo, reasons, "")
	if err != nil {
		t.Fatal(err)
	}
	return p
}

// fixturePath is the path of a file of the fixture as a pack names it.
func fixturePath(t *testing.T, repo *workitem.Repo, path string) string {
	t.Helper()
	r, err := filepath.Rel(repo.MainRoot, path)
	if err != nil {
		t.Fatal(err)
	}
	return filepath.ToSlash(r)
}

// ADR-0217 §1: each kind of reason parses in any padding and any case of
// its ID, a reason given twice is kept once, and every refusal names the
// known kinds in their forms.
func TestParseReasons(t *testing.T) {
	got, err := ParseReasons([]string{"thread:th-7", "thread:TH-0007", "review:S-12", "planned:e-3", "planned:S-0004", "draft:s-5", "release", "ready", "ready"})
	if err != nil {
		t.Fatal(err)
	}
	var each []string
	for _, r := range got {
		each = append(each, r.String())
	}
	want := []string{"thread:TH-0007", "review:S-0012", "planned:E-0003", "planned:S-0004", "draft:S-0005", "release", "ready"}
	if !slices.Equal(each, want) {
		t.Errorf("parsed %v; want %v", each, want)
	}
	forms := "thread:TH-nnnn, review:S-nnnn, planned:E-nnnn or planned:S-nnnn, draft:S-nnnn, release, or ready"
	if ReasonForms() != forms {
		t.Errorf("forms %q; want %q", ReasonForms(), forms)
	}
	for in, says := range map[string]string{
		"bogus:S-0001":  `no such kind of reason "bogus"`,
		"":              `no such kind of reason ""`,
		"thread":        "thread needs the ID of a thread (TH-nnnn)",
		"review:":       "review needs the ID of a story (S-nnnn)",
		"release:S-1":   "release takes no ID",
		"ready:x":       "ready takes no ID",
		"draft:E-0001":  "E-0001 is not the ID of a story (S-nnnn)",
		"thread:S-0001": "S-0001 is not the ID of a thread (TH-nnnn)",
		"planned:T-1":   "T-1 is not the ID of an epic (E-nnnn) or a story (S-nnnn)",
		"review:S-x":    "S-x is not the ID of a story (S-nnnn)",
	} {
		_, err := ParseReasons([]string{"ready", in})
		if err == nil {
			t.Errorf("%q parsed", in)
			continue
		}
		if !strings.Contains(err.Error(), says) || !strings.HasSuffix(err.Error(), "; a reason is "+forms) {
			t.Errorf("%q: %v; want it to say %q and name the kinds", in, err, says)
		}
	}
	repo, _, _ := reasonsRepo(t)
	if _, err := ForReasons(repo, nil, ""); err == nil || !strings.Contains(err.Error(), forms) {
		t.Errorf("no reasons: %v; want a refusal naming the kinds", err)
	}
	if _, err := ForReasons(repo, []string{"thread:TH-0099"}, ""); err == nil || !strings.Contains(err.Error(), "TH-0099") {
		t.Errorf("a thread that does not exist: %v", err)
	}
}

// ADR-0217 §2: the words are the kinds of reason and the built-in strategic
// roles' topics, each once.
func TestSectionWords(t *testing.T) {
	want := []string{"thread", "review", "planned", "draft", "release", "ready", "planning", "orchestration", "analysis", "digest"}
	if got := SectionWords(); !slices.Equal(got, want) {
		t.Errorf("section words %v; want %v", got, want)
	}
}

// ADR-0217 §2, S-0537: a reason pack keeps a section tagged all alone; one
// tagged all and review only for a review; one tagged all and orchestration
// for any reason; one tagged all and analysis for none, though the general
// orchestrate pack keeps it; and one without all as the role's topics have
// it.
func TestAReasonPackKeepsTheSectionsItsReasonsName(t *testing.T) {
	repo, epic, story := reasonsRepo(t)
	general, err := ForStrategic(repo, conventions.RoleOrchestrate, "", "", "")
	if err != nil {
		t.Fatal(err)
	}
	all := []string{"Orchestration", "Every run", "Reviews", "Threads", "Orchestrating", "Analysis", "Project additions"}
	if got := keptHeadings(general, "orchestration.md"); !slices.Equal(got, all) {
		t.Errorf("general pack keeps %v; want %v", got, all)
	}
	for reason, want := range map[string][]string{
		"review:" + story.ID:  {"Orchestration", "Every run", "Reviews", "Orchestrating", "Project additions"},
		"thread:TH-0001":      {"Orchestration", "Every run", "Threads", "Orchestrating", "Project additions"},
		"draft:" + story.ID:   {"Orchestration", "Every run", "Orchestrating", "Project additions"},
		"planned:" + epic.ID:  {"Orchestration", "Every run", "Orchestrating", "Project additions"},
		"release":             {"Orchestration", "Every run", "Orchestrating", "Project additions"},
		"ready":               {"Orchestration", "Every run", "Orchestrating", "Project additions"},
		"thread:TH-0001 both": {"Orchestration", "Every run", "Reviews", "Threads", "Orchestrating", "Project additions"},
	} {
		reasons := []string{reason}
		if r, ok := strings.CutSuffix(reason, " both"); ok {
			reasons = []string{r, "review:" + story.ID}
		}
		p := mustReasons(t, repo, reasons...)
		if got := keptHeadings(p, "orchestration.md"); !slices.Equal(got, want) {
			t.Errorf("%v keeps %v; want %v", reasons, got, want)
		}
		text := p.Conventions[len(p.Conventions)-1].Text
		if strings.Contains(text, "Analysis rule") || strings.Contains(text, "Svelte rule") {
			t.Errorf("%v: the convention holds what no reason needs:\n%s", reasons, text)
		}
		if !slices.Contains(p.Omitted, "orchestration.md § Analysis (all, analysis)") {
			t.Errorf("%v: left out %v; want the analysis section listed", reasons, p.Omitted)
		}
	}
}

// ADR-0217 §3, S-0537 criterion 1: a thread loads the thread with its
// entries and the work item it is on, whole; a thread on a document loads
// the thread alone. The pack is the orchestrator's conventions without the
// README, and holds no briefs, no ranked sections, no issues, and no
// catalog.
func TestAThreadReasonLoadsTheThreadAndItsItem(t *testing.T) {
	repo, _, story := reasonsRepo(t)
	p := mustReasons(t, repo, "thread:th-1")
	th, err := threads.Get(repo, "TH-0001")
	if err != nil {
		t.Fatal(err)
	}
	thPath, storyPath := fixturePath(t, repo, th.Path), fixturePath(t, repo, story.Path)
	got := reasonItems(p)
	raw, _ := os.ReadFile(th.Path)
	if it := got[thPath]; it.Text != string(raw) || it.ID != "TH-0001" || it.Reason != "thread:TH-0001" || it.Size != len(raw) || !strings.Contains(it.Text, "Which way round? Asked in full.") {
		t.Errorf("thread item %+v", it)
	}
	raw, _ = os.ReadFile(story.Path)
	if it := got[storyPath]; it.Text != string(raw) || it.ID != story.ID || it.Reason != "thread:TH-0001" {
		t.Errorf("story item %+v", it)
	}
	if len(got) != 2 || len(p.Items) != 2 || p.Items[0].Path != thPath {
		t.Errorf("items %+v; want the thread then its story", p.Items)
	}
	conv, _ := paths(p)
	for i, c := range conv {
		conv[i] = filepath.Base(c)
	}
	if !slices.Equal(conv, []string{"safety.md", "strategy.md", "orchestration.md"}) || p.README != nil {
		t.Errorf("conventions %v, readme %v", conv, p.README != nil)
	}
	if p.Issues != "" || len(p.Catalog.NotLoaded)+len(p.Catalog.InPart) != 0 || p.Role != conventions.RoleOrchestrate || !slices.Equal(p.Reasons, []string{"thread:TH-0001"}) {
		t.Errorf("issues %q, catalog %+v, role %q, reasons %v", p.Issues, p.Catalog, p.Role, p.Reasons)
	}
	header, body := p.Header(), p.Body()
	for _, want := range []string{"primed for thread:TH-0001 (ADR-0217)", "doc_search and doc_get find the rest", "reason      " + thPath} {
		if !strings.Contains(header, want) {
			t.Errorf("header lacks %q:\n%s", want, header)
		}
	}
	if !strings.Contains(body, "\n"+thPath+"\n"+strings.Repeat("=", len(thPath))+"\nreason: thread:TH-0001\n\n---\n") || strings.Contains(body, "\ncatalog\n") || strings.Contains(body, "\nbriefs\n") {
		t.Errorf("body:\n%s", body)
	}
	doc := mustReasons(t, repo, "thread:TH-0002")
	if got := reasonItems(doc); len(got) != 1 {
		t.Errorf("a thread on a document loads %v; want the thread alone", got)
	}
}

// ADR-0217 §3, S-0537 criterion 1: a draft loads the draft story whole.
func TestADraftReasonLoadsTheStory(t *testing.T) {
	repo, _, story := reasonsRepo(t)
	p := mustReasons(t, repo, "draft:S-1")
	raw, _ := os.ReadFile(story.Path)
	if len(p.Items) != 1 || p.Items[0].Text != string(raw) || p.Items[0].Reason != "draft:"+story.ID || p.Items[0].Step != StepReason {
		t.Errorf("items %+v; want the story whole", p.Items)
	}
}

// ADR-0217 §3, S-0537 criterion 1: a plan of an epic loads the epic, its
// stories but the cancelled one, and the planner's open thread on it, not
// its resolved one nor an agent's; a plan of a story loads the story and
// its tasks.
func TestAPlannedReasonLoadsTheItemWhatThePlannerWroteAndItsThread(t *testing.T) {
	repo, epic, story := reasonsRepo(t)
	p := mustReasons(t, repo, "planned:"+epic.ID)
	var ids []string
	for _, it := range p.Items {
		if it.Step == StepReason && it.Reason == "planned:"+epic.ID {
			ids = append(ids, it.ID)
		}
	}
	if want := []string{epic.ID, story.ID, "S-0002", "TH-0003"}; !slices.Equal(ids, want) {
		t.Errorf("planned %s loads %v; want %v", epic.ID, ids, want)
	}
	p = mustReasons(t, repo, "planned:"+story.ID)
	ids = nil
	for _, it := range p.Items {
		ids = append(ids, it.ID)
	}
	if want := []string{story.ID, "T-0001"}; !slices.Equal(ids, want) {
		t.Errorf("planned %s loads %v; want %v", story.ID, ids, want)
	}
}

// ADR-0217 §3: an item two reasons name loads once, with the second reason
// beside the first, and the conventions load once.
func TestAnItemTwoReasonsNameLoadsOnce(t *testing.T) {
	repo, epic, story := reasonsRepo(t)
	p := mustReasons(t, repo, "thread:TH-0001", "draft:"+story.ID, "planned:"+epic.ID, "thread:TH-1")
	count := map[string]int{}
	for _, it := range p.Items {
		count[it.ID]++
	}
	for id, n := range count {
		if n != 1 {
			t.Errorf("%s loads %d times", id, n)
		}
	}
	var st Item
	for _, it := range p.Items {
		if it.ID == story.ID {
			st = it
		}
	}
	if st.Reason != "thread:TH-0001" || !slices.Equal(st.Also, []string{"draft:" + story.ID, "planned:" + epic.ID}) {
		t.Errorf("story item reason %q, also %v", st.Reason, st.Also)
	}
	conv, _ := paths(p)
	if len(conv) != 3 || !slices.Equal(p.Reasons, []string{"thread:TH-0001", "draft:" + story.ID, "planned:" + epic.ID}) {
		t.Errorf("conventions %v, reasons %v", conv, p.Reasons)
	}
	limit := encoded(p) / 3
	parts := p.Parts(limit)
	if len(parts) < 2 {
		t.Fatalf("the pack is in %d part; the test wants more", len(parts))
	}
	for _, q := range parts {
		if n := size(t, q); n > limit {
			t.Errorf("part %d encodes to %d bytes, over %d", q.PartNumber, n, limit)
		}
	}
	if !slices.Equal(parts[0].Reasons, p.Reasons) {
		t.Errorf("part 1 reasons %v; want %v", parts[0].Reasons, p.Reasons)
	}
	convs, items := rejoin(t, parts)
	if !reflect.DeepEqual(convs, p.Conventions) || !reflect.DeepEqual(items, p.Items) {
		t.Errorf("the parts' conventions and items are not the pack's")
	}
	if q, err := p.Part(len(parts), limit); err != nil || q.PartNumber != len(parts) {
		t.Errorf("last part: %v", err)
	}
}

// ADR-0217 §4: the digest's answers are ranked against the reasons'
// threads, not against every thread awaiting the operator.
func TestAReasonPackRanksTheDigestAgainstItsThreads(t *testing.T) {
	repo := answersRepo(t)
	p := mustReasons(t, repo, "thread:TH-0001")
	got := digestItems(p)
	if one := got["D-0001"]; one.Reason != "thread TH-0001; rank 1" {
		t.Errorf("D-0001 %+v; want it first for TH-0001", one)
	}
	if _, ok := got["D-0002"]; ok {
		t.Errorf("digest %v holds TH-0002's answer", got)
	}
	if !strings.Contains(p.Body(), reasonsHead[groupDigest]) {
		t.Errorf("body lacks the reasons' digest heading")
	}
	if p := mustReasons(t, repo, "ready"); len(digestItems(p)) != 0 {
		t.Errorf("a pack with no thread loads digest %v", digestItems(p))
	}
}

// S-0537 criterion 2: a pack primed for one thread is under a quarter of
// the general orchestrate pack over this repository's conventions and
// design. The general pack is built over this repository as it stands; the
// thread pack over a copy of its manifest and design, with a story and a
// thread on it, so that nothing is written into this repository.
func TestAThreadPackIsUnderAQuarterOfTheGeneralPack(t *testing.T) {
	if testing.Short() {
		t.Skip("integration: reads this repository's design")
	}
	here, err := workitem.Open(".")
	if err != nil {
		t.Fatal(err)
	}
	general, err := ForStrategic(here, conventions.RoleOrchestrate, "", "", "")
	if err != nil {
		t.Fatal(err)
	}
	root := t.TempDir()
	copyTree(t, filepath.Join(here.Root, "system-flow.yaml"), filepath.Join(root, "system-flow.yaml"))
	copyTree(t, filepath.Join(here.Root, "design"), filepath.Join(root, "design"))
	for _, d := range []string{"wip/kanban/epics", "wip/kanban/stories", "wip/kanban/tasks", "wip/agents", "wip/threads"} {
		if err := os.MkdirAll(filepath.Join(root, d), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	now := time.Date(2026, 10, 29, 9, 0, 0, 0, time.UTC)
	epic, err := repo.Create(workitem.NewOptions{Type: workitem.Epic, Title: "Orchestration", Owner: "alex", Now: now})
	if err != nil {
		t.Fatal(err)
	}
	story, err := repo.Create(workitem.NewOptions{Type: workitem.Story, Title: "The orchestrator primes for its reasons", Parent: epic.ID, Owner: "alex", Now: now,
		Body: "## Goal\n\nThe orchestrator primes for the reasons of its run.\n\n## Acceptance criteria\n\n- [ ] A pack for one thread is under a quarter of the general pack.\n"})
	if err != nil {
		t.Fatal(err)
	}
	th, err := threads.New(repo, threads.NewOptions{Title: "Which pack does a thread run read?", On: story.ID, Author: "agent-" + story.ID, Now: now,
		Text: "Should a run for one thread read the general pack or one primed for the thread?"})
	if err != nil {
		t.Fatal(err)
	}
	p, err := ForReasons(repo, []string{"thread:" + th.ID}, "")
	if err != nil {
		t.Fatal(err)
	}
	conventionBytes := 0
	for _, c := range p.Conventions {
		conventionBytes += c.Size
	}
	t.Logf("general orchestrate pack %d bytes; pack for thread:%s %d bytes (%d of conventions, %d items), %.1f%% of it", general.Size.Bytes, th.ID, p.Size.Bytes, conventionBytes, len(p.Items), 100*float64(p.Size.Bytes)/float64(general.Size.Bytes))
	if p.Size.Bytes*4 >= general.Size.Bytes {
		t.Errorf("the pack for one thread is %d bytes, %d of them conventions; want under a quarter of the general pack's %d, %d", p.Size.Bytes, conventionBytes, general.Size.Bytes, general.Size.Bytes/4)
	}
}

// copyTree copies the file or folder at from to to.
func copyTree(t *testing.T, from, to string) {
	t.Helper()
	err := filepath.WalkDir(from, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		r, err := filepath.Rel(from, path)
		if err != nil {
			return err
		}
		dst := filepath.Join(to, r)
		if d.IsDir() {
			return os.MkdirAll(dst, 0o755)
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
			return err
		}
		return os.WriteFile(dst, data, 0o644)
	})
	if err != nil {
		t.Fatal(err)
	}
}
