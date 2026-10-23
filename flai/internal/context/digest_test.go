package context

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/conventions"
	"github.com/bytepunx/system-flow/flai/internal/digest"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// largeBody is the body of an entry too large for the digest's share of a
// 16 KB budget, 2 KB, though the budget itself has room for it.
var largeBody = strings.Repeat("The digest pack loads entries ranked against the story. ", 60)

// digestRepo is a project with one convention, a design file on cli, and a
// story on cli that touches flai/cmd/prime.go, whose task touches
// flai/internal/context/digest.go. Its digest holds an entry for each way a
// story's pack selects one or leaves it out:
//
//   - D-0001, on svelte, its path a folder holding the story's touch;
//   - D-0002, on cli and all, sharing words with the story;
//   - D-0003, on cli, sharing words with the story, too large for the share;
//   - D-0004, absorbed, on cli, its path the story's touch;
//   - D-0005, compacted, on cli, its title sharing words with the story;
//   - D-0006, on all alone, sharing words with the story;
//   - D-0007, on cli, with no word in common with the story;
//   - D-0008, on svelte, its path a folder holding the task's touch.
func digestRepo(t *testing.T) (*workitem.Repo, string) {
	t.Helper()
	root := t.TempDir()
	write := func(rel, s string) {
		p := filepath.Join(root, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte(s), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	write("system-flow.yaml", "version: 1\nname: t\nkey: t\nlayout:\n  design: design\n  docs: docs\n  wip: wip\n")
	for _, d := range []string{"wip/kanban/epics", "wip/kanban/stories", "wip/kanban/tasks", "wip/agents"} {
		if err := os.MkdirAll(filepath.Join(root, d), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	write("design/conventions/README.md", "# Conventions\n\n- [safety.md](safety.md)\n")
	write("design/conventions/safety.md", "---\ntitle: Safety\nupdated: 2026-10-23\naudience: agent\norder: 10\nstatus: active\ntopics: [all]\n---\n\n# Safety\n\n## Rules\n\n- Safety rule.\n\n<!-- system-flow:end-of-baseline -->\n\n## Project additions\n")
	write("design/system/cli.md", "---\ntitle: CLI\nupdated: 2026-10-23\nstatus: active\ntopics: [cli]\n---\n\n# CLI\n\nThe command line.\n\n## Commands\n\nThe commands prime the pack.\n")
	entry := func(id, slug, title, topics, paths, date, status, body string) {
		absorbed := `""`
		if status == "absorbed" {
			absorbed = "design/system/cli.md"
		}
		text := "---\nid: " + id + "\ntitle: " + title + "\nkind: decision\nsources: [S-0001]\ntopics: [" + topics + "]\npaths: [" + paths + "]\ndate: " + date + "\nstatus: " + status + "\nabsorbed_by: " + absorbed + "\nasked_again: []\ncited: []\n---\n\n# " + id + " " + title + "\n"
		if body != "" {
			text += "\n" + body + "\n"
		}
		write("design/digest/"+id+"-"+slug+".md", text)
	}
	entry("D-0001", "touches", "Prime reads the touches", "svelte", "flai/cmd", "2026-10-01", "active", "Each touch is met against the entry's paths.")
	entry("D-0002", "ranks", "The digest pack ranks entries", "cli, all", "docs", "2026-10-02", "active", "Entries rank against the story goal.")
	entry("D-0003", "large", "Large digest entry", "cli", "", "2026-10-03", "active", largeBody)
	entry("D-0004", "absorbed", "The digest pack absorbed entry", "cli", "flai/cmd/prime.go", "2026-10-04", "absorbed", "Absorbed digest entries never load in the pack.")
	entry("D-0005", "compacted", "Digest entries ranked in the pack", "cli", "", "2026-10-05", "compacted", "")
	entry("D-0006", "everywhere", "The digest pack for everyone", "all", "", "2026-10-06", "active", "Digest entries ranked against every story.")
	entry("D-0007", "unrelated", "Unrelated wombat", "cli", "", "2026-10-07", "active", "Zebra quokka.")
	entry("D-0008", "task", "Unrelated platypus", "svelte", "flai/internal/context", "2026-10-08", "active", "Echidna.")
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	now := time.Date(2026, 10, 23, 9, 0, 0, 0, time.UTC)
	st, err := repo.Create(workitem.NewOptions{Type: workitem.Story, Title: "Prime loads the digest", Owner: "alex", Topics: []string{"cli"}, Touches: []string{"flai/cmd/prime.go"}, Now: now,
		Body: "## Goal\n\nThe pack loads digest entries ranked against the story.\n\n## Acceptance criteria\n\n- [ ] it loads\n"})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := repo.Create(workitem.NewOptions{Type: workitem.Task, Title: "The step", Parent: st.ID, Owner: "alex", Touches: []string{"flai/internal/context/digest.go"}, Now: now}); err != nil {
		t.Fatal(err)
	}
	return repo, st.ID
}

// digestItems is the pack's digest items by ID.
func digestItems(p *Pack) map[string]Item {
	out := map[string]Item{}
	for _, it := range p.Items {
		if it.Step == StepDigest {
			out[it.ID] = it
		}
	}
	return out
}

// tried is the IDs of the digest items a story's pack tries, in order.
func tried(t *testing.T, repo *workitem.Repo, id string, whole bool) []string {
	t.Helper()
	st, err := loadStory(repo, id)
	if err != nil {
		t.Fatal(err)
	}
	var out []string
	for _, it := range digestFor(st.digest, []string{"cli", "all"}, st.claim, Query(st.item.Title, st.item.Body), whole) {
		out = append(out, it.ID)
	}
	return out
}

// ADR-0194, S-0430 criterion 2: a story's pack loads the entries its claim
// selects first, then those its topics select that share words with it,
// within an eighth of the budget, each marked as the digest's with its
// reason and size, after the briefs and before the ranked sections. One too
// large for the share is selected but not loaded; an absorbed entry, one on
// all alone, and one with no word in common with the story never load.
func TestAStorysPackLoadsTheDigestEntriesItsTopicsAndClaimSelect(t *testing.T) {
	repo, id := digestRepo(t)
	if got := tried(t, repo, id, true); len(got) != 5 || got[0] != "D-0001" || got[1] != "D-0008" || !slices.Contains(got, "D-0003") {
		t.Fatalf("tried %v; want D-0001 and D-0008 first, then D-0002, D-0003, and D-0005", got)
	}
	p, err := ForStory(repo, id, "16KB")
	if err != nil {
		t.Fatal(err)
	}
	got := digestItems(p)
	if len(got) != 4 {
		t.Fatalf("digest items %+v; want D-0001, D-0002, D-0005, D-0008", got)
	}
	for _, out := range []string{"D-0003", "D-0004", "D-0006", "D-0007"} {
		if _, ok := got[out]; ok {
			t.Errorf("%s loaded: %+v", out, got[out])
		}
	}
	if r := got["D-0001"].Reason; r != "path: flai/cmd/prime.go; rank 1" {
		t.Errorf("D-0001 reason %q", r)
	}
	if r := got["D-0008"].Reason; r != "path: flai/internal/context/digest.go; rank 2" {
		t.Errorf("D-0008 reason %q", r)
	}
	if r := got["D-0002"].Reason; !strings.HasPrefix(r, "topics: cli; rank ") || strings.HasSuffix(r, "compacted") {
		t.Errorf("D-0002 reason %q", r)
	}
	if r := got["D-0005"].Reason; !strings.HasPrefix(r, "topics: cli; rank ") || !strings.HasSuffix(r, "; compacted") {
		t.Errorf("D-0005 reason %q", r)
	}
	one := got["D-0001"]
	if one.Path != "design/digest/D-0001-touches.md" || one.Title != "Prime reads the touches" || one.Text != "# D-0001 Prime reads the touches\n\nEach touch is met against the entry's paths.\n" || one.Size != len(one.Text) {
		t.Errorf("D-0001 %+v", one)
	}
	if c := got["D-0005"]; c.Text != "" || c.Size != 0 || c.Whole != 0 {
		t.Errorf("a compacted entry loads a body: %+v", c)
	}
	large := digestItem(digest.Entry{ID: "D-0003", Title: "Large digest entry", Path: "design/digest/D-0003-large.md", Status: "active", Body: strings.TrimSpace(largeBody)}, "topics: cli; rank 4", true)
	if n := len(large.Printed()); n <= p.Budget/DigestShare || n >= p.Budget {
		t.Errorf("D-0003 prints %d bytes; the fixture wants it over the share %d and under the budget %d", n, p.Budget/DigestShare, p.Budget)
	}

	used, first, last := 0, -1, -1
	for i, it := range p.Items {
		switch it.Step {
		case StepDigest:
			used += len(it.Printed())
			if first < 0 {
				first = i
			}
			last = i
		case StepBriefed:
			if first >= 0 {
				t.Errorf("brief %s after a digest item", it.Path)
			}
		case StepRanked:
			if first < 0 || i < last {
				t.Errorf("ranked %s before a digest item", it.Path)
			}
		}
	}
	if used > p.Budget/DigestShare || p.Size.Bytes > p.Budget {
		t.Errorf("digest prints %d bytes of a share of %d; pack %d of %d", used, p.Budget/DigestShare, p.Size.Bytes, p.Budget)
	}

	body, header := p.Body(), p.Header()
	for _, want := range []string{groupHead[groupDigest], "\ndesign/digest/D-0001-touches.md\n", "reason: path: flai/cmd/prime.go; rank 1\n\n# D-0001 Prime reads the touches\n", "\n- D-0005 Digest entries ranked in the pack (design/digest/D-0005-compacted.md; " + got["D-0005"].Reason + ")\n"} {
		if !strings.Contains(body, want) {
			t.Errorf("body lacks %q:\n%s", want, body)
		}
	}
	if strings.Count(body, "\ndigest\n======\n") != 1 {
		t.Errorf("the digest heading is not printed once:\n%s", body)
	}
	if want := digestRow(4, 3); !strings.Contains(header, want) || want != "4 digest entries, 3 with their bodies, the rest as one-line briefs" {
		t.Errorf("header lacks %q:\n%s", want, header)
	}

	b, err := json.Marshal(got["D-0001"])
	if err != nil {
		t.Fatal(err)
	}
	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatal(err)
	}
	if m["step"] != "digest" || m["id"] != "D-0001" || m["reason"] != "path: flai/cmd/prime.go; rank 1" || m["size"] != float64(len(one.Text)) {
		t.Errorf("json %s", b)
	}
}

// ADR-0194: an absorbed entry is never selected, though its path is the
// story's touch and its words are the story's, in the story's pack or a
// sub-agent's, whatever the budget.
func TestAnAbsorbedEntryNeverLoads(t *testing.T) {
	repo, id := digestRepo(t)
	if slices.Contains(tried(t, repo, id, true), "D-0004") {
		t.Fatalf("D-0004 tried")
	}
	for _, budget := range []string{"", "16KB", "1MB"} {
		p, err := ForStory(repo, id, budget)
		if err != nil {
			t.Fatal(err)
		}
		r, err := ForRole(repo, id, conventions.RoleExplore, budget)
		if err != nil {
			t.Fatal(err)
		}
		for _, pack := range []*Pack{p, r} {
			if it, ok := digestItems(pack)["D-0004"]; ok || strings.Contains(pack.Body(), "D-0004") {
				t.Errorf("budget %q, role %q: the absorbed entry loads: %+v", budget, pack.Role, it)
			}
		}
	}
}

// ADR-0194: entries are selected by a topic other than all, or by a path
// meeting the claim, the claim's entry being equal to the path, a folder
// holding it, or held by it; the reason names what selected it.
func TestSelectDigestTakesTopicsOtherThanAllAndPathsMeetingTheClaim(t *testing.T) {
	entries := []digest.Entry{
		{ID: "D-0001", Status: "active", Topics: []string{"all"}},
		{ID: "D-0002", Status: "active", Topics: []string{"cli", "all", "go", "cli"}},
		{ID: "D-0003", Status: "absorbed", Topics: []string{"cli"}, Paths: []string{"flai/cmd"}},
		{ID: "D-0004", Status: "compacted", Topics: []string{"svelte"}, Paths: []string{"flai/cmd/prime.go"}},
		{ID: "D-0005", Status: "active", Paths: []string{"flai"}},
		{ID: "D-0006", Status: "active", Paths: []string{"flai/cmdline", "docs"}},
		{ID: "D-0007", Status: "active", Topics: []string{"go"}, Paths: []string{"flai/cmd/prime_test.go"}},
	}
	got := map[string]string{}
	for _, p := range selectDigest(entries, []string{"cli", "go", "all"}, []string{"flai/cmd"}) {
		got[p.entry.ID] = p.reason(1)
	}
	want := map[string]string{
		"D-0002": "topics: cli, go; rank 1",
		"D-0004": "path: flai/cmd; rank 1; compacted",
		"D-0005": "path: flai/cmd; rank 1",
		"D-0007": "path: flai/cmd; rank 1",
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("selected %v, want %v", got, want)
	}
}

// ADR-0194: entries whose paths meet the claim rank first, then by score,
// then the newest date, then by ID; one selected by topics alone with no
// word in common with the query is left out.
func TestRankDigestPutsTheClaimFirstThenScoreThenDate(t *testing.T) {
	p := func(id, path, date string, score float64) pick {
		return pick{entry: digest.Entry{ID: id, Date: date}, path: path, shared: []string{"cli"}, score: score}
	}
	picks := []pick{
		p("D-0001", "", "2026-10-01", 3),
		p("D-0002", "flai", "2026-10-01", 0),
		p("D-0003", "", "2026-10-01", 0),
		p("D-0004", "", "2026-10-09", 3),
		p("D-0005", "flai", "2026-10-01", 1),
		p("D-0006", "", "2026-10-09", 5),
		p("D-0007", "", "2026-10-09", 3),
	}
	var got []string
	for _, r := range rankDigest(picks) {
		got = append(got, r.entry.ID)
	}
	if want := []string{"D-0005", "D-0002", "D-0006", "D-0004", "D-0007", "D-0001"}; !slices.Equal(got, want) {
		t.Errorf("ranked %v, want %v", got, want)
	}
	if len(picks) != 7 || picks[0].entry.ID != "D-0001" {
		t.Errorf("rankDigest changed what it was given: %v", picks)
	}
}

// ADR-0194, ADR-0059: a sub-agent's pack holds the entries the story's
// would try, active or compacted, each as a one-line brief with no body,
// after the design and ADR briefs, with no share of its own; those with no
// room are counted with the briefs left out.
func TestASubAgentsPackHoldsTheDigestAsBriefs(t *testing.T) {
	repo, id := digestRepo(t)
	p, err := ForRole(repo, id, conventions.RoleExplore, "")
	if err != nil {
		t.Fatal(err)
	}
	var ids []string
	seen := false
	for _, it := range p.Items {
		switch it.Step {
		case StepDigest:
			seen = true
			ids = append(ids, it.ID)
			if it.Text != "" || it.Size != 0 {
				t.Errorf("%s loads a body: %+v", it.ID, it)
			}
		case StepBriefed:
			if seen {
				t.Errorf("brief %s after a digest item", it.Path)
			}
		default:
			t.Errorf("a sub-agent's pack loads %s %s", it.Step, it.Path)
		}
	}
	if want := tried(t, repo, id, false); !slices.Equal(ids, want) {
		t.Errorf("digest briefs %v, want %v", ids, want)
	}
	if large := digestItems(p)["D-0003"]; large.Whole != len("# D-0003 Large digest entry\n\n"+strings.TrimSpace(largeBody)+"\n") {
		t.Errorf("D-0003 brief %+v", large)
	}
	body := p.Body()
	if !strings.Contains(body, roleHead[groupDigest]) || strings.Contains(body, "Each touch is met") {
		t.Errorf("body:\n%s", body)
	}
	if !strings.Contains(body, "\n- D-0001 Prime reads the touches (design/digest/D-0001-touches.md, ") {
		t.Errorf("D-0001 is not one line with its size:\n%s", body)
	}

	small, err := ForRole(repo, id, conventions.RoleExplore, strconv.Itoa(p.Size.Bytes-200))
	if err != nil {
		t.Fatal(err)
	}
	if small.BriefsLeftOut == 0 || small.Size.Bytes > small.Budget || len(digestItems(small)) >= len(ids) {
		t.Errorf("a budget 200 bytes under the pack: %d left out, size %d of %d, %d digest briefs", small.BriefsLeftOut, small.Size.Bytes, small.Budget, len(digestItems(small)))
	}
}

// ADR-0194: the digest step never takes the pack over its budget nor its
// items over an eighth of it, and loads nothing when the pack is over its
// budget before it.
func TestTheDigestStepNeverTakesThePackOverItsBudget(t *testing.T) {
	repo, id := digestRepo(t)
	loaded := false
	for budget := 2048; budget <= 40*1024; budget += 512 {
		p, err := ForStory(repo, id, strconv.Itoa(budget))
		if err != nil {
			t.Fatal(err)
		}
		used := 0
		for _, it := range digestItems(p) {
			used += len(it.Printed())
		}
		switch {
		case p.Exceeded != "" && used > 0:
			t.Errorf("budget %d, exceeded %s: digest loads %d bytes", budget, p.Exceeded, used)
		case p.Exceeded == "" && (p.Size.Bytes > budget || used > budget/DigestShare):
			t.Errorf("budget %d: pack %d bytes, digest %d of a share of %d", budget, p.Size.Bytes, used, budget/DigestShare)
		}
		loaded = loaded || used > 0
	}
	if !loaded {
		t.Errorf("no budget loaded a digest entry")
	}
}

// ADR-0194, ADR-0104: digest items go into parts as any item does, and
// each part that holds one prints the digest heading.
func TestDigestItemsGoIntoParts(t *testing.T) {
	repo, id := digestRepo(t)
	p, err := ForStory(repo, id, "16KB")
	if err != nil {
		t.Fatal(err)
	}
	parts := p.Parts(encoded(p) * 2 / 3)
	if len(parts) < 2 {
		t.Fatalf("the pack is in %d part; the test wants more", len(parts))
	}
	var ids []string
	for _, q := range parts {
		holds := false
		for _, it := range q.Items {
			if it.Step == StepDigest {
				holds = true
				ids = append(ids, it.ID)
			}
		}
		if headed := strings.Contains(q.Body(), "\ndigest\n======\n"); headed != holds {
			t.Errorf("part %d holds digest items %v, prints the heading %v", q.PartNumber, holds, headed)
		}
	}
	var want []string
	for _, it := range p.Items {
		if it.Step == StepDigest {
			want = append(want, it.ID)
		}
	}
	if !slices.Equal(ids, want) {
		t.Errorf("parts hold digest items %v, the pack %v", ids, want)
	}
}
