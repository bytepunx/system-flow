package workitem

import (
	"encoding/json"
	"fmt"
	"reflect"
	"slices"
	"strings"
	"testing"
	"time"
)

// S-0087, found in T-0318's end-to-end review: acceptance archives a story
// the moment it merges, and the done column used to drop every archived
// item without exception, so a story could never be seen there before it
// was published, which is exactly the state the operator asked to be able
// to see. A done, archived story named in pendingPublish stays on the done
// column; one not named there, or archived under any other status, does not.
func TestArchivedDoneStoryStaysUntilPublished(t *testing.T) {
	items := []*Item{
		{ID: "S-0001", Type: Story, Title: "Waiting", Status: Done, Archived: true},
		{ID: "S-0002", Type: Story, Title: "Published", Status: Done, Archived: true},
		{ID: "S-0003", Type: Story, Title: "Cancelled but named anyway", Status: Cancelled, Archived: true},
		{ID: "S-0004", Type: Story, Title: "Not done", Status: InProgress},
	}
	board := &Board{}
	pending := map[string]bool{"S-0001": true, "S-0003": true}
	now := time.Date(2026, 9, 21, 12, 0, 0, 0, time.UTC)

	v := NewBoardView(items, board, now, false, pending, nil)
	got := ids(v.Columns[Done])
	if len(got) != 1 || got[0] != "S-0001" {
		t.Errorf("done column: %v", got)
	}
	if !v.Columns[Done][0].Archived {
		t.Error("the card says it is archived, so the dashboard can tell it is there only while pending (S-0174)")
	}
	if len(v.Columns[Cancelled]) != 0 {
		t.Error("pendingPublish names a cancelled item too, but only a done one may stay")
	}

	// nothing named as pending: the old behaviour, nothing archived shows
	if v := NewBoardView(items, board, now, false, nil, nil); len(v.Columns[Done]) != 0 || v.Unpublished != nil {
		t.Errorf("no pending map: %v, unpublished %v", ids(v.Columns[Done]), v.Unpublished)
	}
}

// S-0195, ADR-0067: the view lists what is accepted and not yet published,
// every ID pendingPublish names, in order, for whoever publishes.
func TestBoardViewListsWhatIsUnpublished(t *testing.T) {
	items := []*Item{{ID: "S-0002", Type: Story, Title: "Accepted", Status: Done, Archived: true}}
	pending := map[string]bool{"S-0009": true, "S-0002": true, "S-0005": false}
	v := NewBoardView(items, &Board{}, time.Now(), false, pending, nil)
	if got := v.Unpublished; len(got) != 2 || got[0] != "S-0002" || got[1] != "S-0009" {
		t.Errorf("unpublished: %v", got)
	}
}

// S-0201: a story an agent wrote and the operator has not finalized carries
// draft on its card, so the dashboard can mark it; any other card leaves the
// key out, and only a story is ever a draft.
func TestADraftStorysCardSaysSo(t *testing.T) {
	items := []*Item{
		{ID: "S-0001", Type: Story, Title: "Drafted", Status: Backlog, Draft: true},
		{ID: "S-0002", Type: Story, Title: "Finalized", Status: Backlog},
		{ID: "E-0001", Type: Epic, Title: "Not a story", Status: Backlog, Draft: true},
	}
	v := NewBoardView(items, &Board{}, time.Now(), true, nil, nil)
	cards := map[string]BoardCard{}
	for _, c := range v.Columns[Backlog] {
		cards[c.ID] = c
	}
	for _, tc := range []struct {
		id    string
		draft bool
	}{
		{"S-0001", true},
		{"S-0002", false},
		{"E-0001", false},
	} {
		c, ok := cards[tc.id]
		if !ok {
			t.Fatalf("%s: no card in %v", tc.id, ids(v.Columns[Backlog]))
		}
		if c.Draft != tc.draft {
			t.Errorf("%s: draft %v, want %v", tc.id, c.Draft, tc.draft)
		}
		data, err := json.Marshal(c)
		if err != nil {
			t.Fatalf("%s: marshal: %v", tc.id, err)
		}
		if got := strings.Contains(string(data), `"draft":true`); got != tc.draft {
			t.Errorf("%s: JSON %s", tc.id, data)
		}
		if !tc.draft && strings.Contains(string(data), `"draft"`) {
			t.Errorf("%s: a card that is not a draft has a draft key: %s", tc.id, data)
		}
	}
}

// S-0338: flai board --json gives a ready card held on overlap alone the
// conversation in which flai asked its holding story's agent about the hold,
// and leaves asked out of a held card no one was asked about and of one asked
// only before it last entered ready, or in a closed conversation; it gives a
// card started on a share the shared paths, the split, and the conversation,
// and none once the share has ended.
func TestBoardCardsNameTheAskAndTheShare(t *testing.T) {
	r := newProject(t)
	readyAt := func(id, at string, touches ...string) *Item {
		s := claimed(id, Ready, touches...)
		s.Transitions = []Transition{{To: Ready, At: at, By: "test"}}
		return s
	}
	items := []*Item{
		claimed("S-0001", InProgress, "flai/cmd"),
		claimed("S-0002", InProgress, "docs"),
		claimed("S-0003", InProgress, "flaiover/src/lib", "flaiover/src/routes"),
		readyAt("S-0007", "2026-10-08T09:00:00Z", "flai/cmd/serve"),        // asked
		readyAt("S-0008", "2026-10-08T09:00:00Z", "docs/users"),            // asked before it last entered ready
		readyAt("S-0009", "2026-10-08T09:00:00Z", "flai/cmd/board.go"),     // not asked
		readyAt("S-0006", "2026-10-08T09:00:00Z", "flaiover/src/lib/x.ts"), // asked in a closed conversation
		claimed("S-0005", InProgress, "flaiover/src/routes"),               // started on a share
	}
	writeConversationBody(t, r, "MS-0001-ask.md", "id: MS-0001\nfrom: S-0007\nto: S-0001\nstatus: open\n",
		askEntry("2026-10-08T08:00:00Z", "S-0007")+askEntry("2026-10-08T10:00:00Z", "S-0007"))
	writeConversationBody(t, r, "MS-0002-ask.md", "id: MS-0002\nfrom: S-0002\nto: S-0008\nstatus: open\n",
		askEntry("2026-10-08T08:00:00Z", "S-0008"))
	writeConversationBody(t, r, "MS-0003-closed.md", "id: MS-0003\nfrom: S-0006\nto: S-0003\nstatus: closed\n",
		askEntry("2026-10-08T10:00:00Z", "S-0006"))
	shares := "shares:\n  - holder: S-0003\n    held: S-0005\n    paths: [flaiover/src/routes]\n    split: \"S-0005 adds the page; S-0003 keeps the layout\"\n    by: agent\n    at: " + shareAt + "\n"
	writeConversationBody(t, r, "MS-0004-share.md", "id: MS-0004\nfrom: S-0005\nto: S-0003\nstatus: open\n"+shares,
		askEntry("2026-10-08T09:30:00Z", "S-0005"))

	cards := func() map[string]BoardCard {
		v := NewBoardView(items, &Board{}, time.Now(), false, nil, r.Holds(items))
		out := map[string]BoardCard{}
		for _, col := range v.Columns {
			for _, c := range col {
				out[c.ID] = c
			}
		}
		return out
	}
	asJSON := func(v any) string {
		data, err := json.Marshal(v)
		if err != nil {
			t.Fatal(err)
		}
		return string(data)
	}
	got := cards()
	for id, want := range map[string]string{
		"S-0007": `[{"by":"S-0001","conversation":"MS-0001","at":"2026-10-08T10:00:00Z"}]`,
		"S-0008": "null",
		"S-0009": "null",
		"S-0006": "null",
	} {
		c := got[id]
		if c.Held == nil || c.Held.Code != HoldOverlap {
			t.Fatalf("%s: held = %+v", id, c.Held)
		}
		if a := asJSON(c.Held.Asked); a != want {
			t.Errorf("%s: asked = %s\nwant     %s", id, a, want)
		}
		if want == "null" && strings.Contains(asJSON(c), `"asked"`) {
			t.Errorf("%s: a card no one was asked about since it entered ready has an asked key: %s", id, asJSON(c))
		}
	}
	want := `[{"holder":"S-0003","held":"S-0005","paths":["flaiover/src/routes"],"split":"S-0005 adds the page; S-0003 keeps the layout","by":"agent","at":"` + shareAt + `","conversation":"MS-0004"}]`
	if s := asJSON(got["S-0005"].Shared); s != want {
		t.Errorf("shared = %s\nwant     %s", s, want)
	}
	for id, c := range got {
		if id != "S-0005" && strings.Contains(asJSON(c), `"shared"`) {
			t.Errorf("%s: shared on a card no share was made for: %s", id, asJSON(c))
		}
	}

	// the share ends when the story it was made for goes back to backlog
	// after it, and the card returned to ready names no share
	items[7].Status = Ready
	items[7].Transitions = []Transition{
		{To: Backlog, At: "2026-10-08T11:00:00Z", By: "test"},
		{To: Ready, At: "2026-10-08T11:30:00Z", By: "test"},
	}
	if c := cards()["S-0005"]; c.Shared != nil || strings.Contains(asJSON(c), `"shared"`) {
		t.Errorf("a share that ended is shown: %s", asJSON(c))
	}
}

func ids(cards []BoardCard) []string {
	out := make([]string, len(cards))
	for i, c := range cards {
		out[i] = c.ID
	}
	return out
}

// boardFlow is a workflow the board's tests run a case under: the board
// LoadBoard reads for a new project running it, with its default limits, and
// its name for each of the baseline's states.
type boardFlow struct {
	name  string
	board *Board
	state func(string) string
}

// boardFlows are the two workflows the board's tests run each case under:
// the baseline, and the fixture with in-progress renamed doing and testing
// added after it (S-0739).
func boardFlows(t *testing.T) []boardFlow {
	t.Helper()
	load := func(r *Repo) *Board {
		b, err := r.LoadBoard()
		if err != nil {
			t.Fatal(err)
		}
		return b
	}
	same := func(s string) string { return s }
	return []boardFlow{
		{"baseline", load(newProject(t)), same},
		{"fixture", load(withTestFlow(t, newProject(t))), inTestFlow},
	}
}

// storiesIn are n stories in state, numbered from first.
func storiesIn(first, n int, state string) []*Item {
	var out []*Item
	for i := range n {
		out = append(out, &Item{ID: fmt.Sprintf("S-%04d", first+i), Type: Story, Status: state})
	}
	return out
}

// S-0739, ADR-0301 §2: the view gives its columns in the workflow's order,
// each with its kind, its WIP limit, and whether it takes a limit and reads
// in the pull order, so a reader needs no list of states of its own; under
// the fixture the renamed and the added state are columns in their places.
func TestBoardViewLanesAreTheWorkflowsColumns(t *testing.T) {
	baseline := []Lane{
		{Name: Backlog, Title: Backlog, Kind: "queue", Ordered: true},
		{Name: Ready, Title: Ready, Kind: "queue", WIPLimit: 5, TakesLimit: true, Ordered: true},
		{Name: InProgress, Title: InProgress, Kind: "active", WIPLimit: 2, TakesLimit: true},
		{Name: Review, Title: Review, Kind: "waiting", WIPLimit: 3, TakesLimit: true},
		{Name: Done, Title: Done, Kind: "accepted"},
		{Name: Cancelled, Title: Cancelled, Kind: "cancelled"},
	}
	fixture := slices.Clone(baseline)
	fixture[2] = Lane{Name: fixtureDoing, Title: fixtureDoing, Kind: "active", WIPLimit: 2, TakesLimit: true}
	fixture = slices.Insert(fixture, 3, Lane{Name: fixtureTesting, Title: fixtureTesting, Kind: "active", WIPLimit: 2, TakesLimit: true})
	want := map[string][]Lane{"baseline": baseline, "fixture": fixture}
	for _, f := range boardFlows(t) {
		v := NewBoardView(nil, f.board, time.Now(), true, nil, nil)
		if !reflect.DeepEqual(v.Lanes, want[f.name]) {
			t.Errorf("%s: lanes\n%+v\nwant\n%+v", f.name, v.Lanes, want[f.name])
		}
	}

	// a column with no limit set shows none, and the JSON leaves out what a
	// column does not have
	b := &Board{WIPLimits: map[string]int{}}
	data, err := json.Marshal(NewBoardView(nil, b, time.Now(), true, nil, nil))
	if err != nil {
		t.Fatal(err)
	}
	if lanes := `"lanes":[{"name":"backlog","title":"backlog","kind":"queue","ordered":true},{"name":"ready","title":"ready","kind":"queue","takes_limit":true,"ordered":true},{"name":"in-progress","title":"in-progress","kind":"active","takes_limit":true},`; !strings.Contains(string(data), lanes) {
		t.Errorf("JSON:\n%s\nwant it to hold\n%s", data, lanes)
	}
}

// S-0739: the pull holds, review holds, and breaches read limited states by
// kind, so each case gives the same answer under the fixture as under the
// baseline, with the fixture's name for in-progress.
func TestBoardViewHoldsByKindUnderEitherWorkflow(t *testing.T) {
	for _, tc := range []struct {
		name            string
		active, waiting int
		pull, review    string
		breaches        []string
	}{
		{"room", 1, 2, "", "", nil},
		{"active full", 2, 0, "the in-progress limit is full (2 of 2)", "", nil},
		{"review full", 1, 3, "review is full (3 of 3): accept or send back a story", "review is full (3 of 3): accept or send back a story", nil},
		{"both full, active first", 3, 4, "the in-progress limit is full (3 of 2)", "review is full (4 of 3): accept or send back a story",
			[]string{"in-progress has 3 stories, limit 2", "review has 4 stories, limit 3"}},
	} {
		for _, f := range boardFlows(t) {
			items := append(storiesIn(1, tc.active, f.state(InProgress)), storiesIn(100, tc.waiting, f.state(Review))...)
			v := NewBoardView(items, f.board, time.Now(), false, nil, nil)
			rename := func(s string) string {
				return strings.ReplaceAll(s, InProgress, f.state(InProgress))
			}
			var breaches []string
			for _, b := range tc.breaches {
				breaches = append(breaches, rename(b))
			}
			if got := v.PullHold(); got != rename(tc.pull) {
				t.Errorf("%s, %s: pull hold %q, want %q", tc.name, f.name, got, rename(tc.pull))
			}
			if got := v.ReviewHold(); got != tc.review {
				t.Errorf("%s, %s: review hold %q, want %q", tc.name, f.name, got, tc.review)
			}
			if v.CanPull() != (tc.pull == "") {
				t.Errorf("%s, %s: can pull %v", tc.name, f.name, v.CanPull())
			}
			if !slices.Equal(v.Breaches, breaches) {
				t.Errorf("%s, %s: breaches %v, want %v", tc.name, f.name, v.Breaches, breaches)
			}
		}
	}

	// under the fixture, testing counts against its own limit: full, it
	// holds the pull though doing has room, and over it, it is a breach
	fixture := boardFlows(t)[1].board
	items := append(storiesIn(1, 1, fixtureDoing), storiesIn(10, 2, fixtureTesting)...)
	v := NewBoardView(items, fixture, time.Now(), false, nil, nil)
	if got, want := v.PullHold(), "the testing limit is full (2 of 2)"; got != want {
		t.Errorf("testing full: pull hold %q, want %q", got, want)
	}
	if v.ReviewHold() != "" || v.Breaches != nil || v.Counts[fixtureDoing] != 1 {
		t.Errorf("testing full: review hold %q, breaches %v, doing %d", v.ReviewHold(), v.Breaches, v.Counts[fixtureDoing])
	}
	items = append(items, storiesIn(20, 1, fixtureTesting)...)
	v = NewBoardView(items, fixture, time.Now(), false, nil, nil)
	if want := []string{"testing has 3 stories, limit 2"}; !slices.Equal(v.Breaches, want) {
		t.Errorf("testing over: breaches %v, want %v", v.Breaches, want)
	}
}

// S-0739: the pull's source reads in pull order, the stories held there are
// marked held, and the shares in force show on a story in the pull's source
// or pulled and not closed, the same under the fixture as under the
// baseline, where the fixture's testing is a pulled state like doing.
func TestBoardViewPullSourceAndSharesByKind(t *testing.T) {
	for _, f := range boardFlows(t) {
		items := []*Item{
			{ID: "S-0003", Type: Story, Status: Ready},
			{ID: "S-0001", Type: Story, Status: Ready, After: []string{"S-0009"}},
			{ID: "S-0002", Type: Story, Status: Ready},
			{ID: "E-0001", Type: Epic, Status: Ready},
			{ID: "S-0009", Type: Story, Status: Backlog},
			{ID: "S-0008", Type: Story, Status: Backlog, After: []string{"S-0009"}},
			{ID: "S-0010", Type: Story, Status: f.state(InProgress)},
			{ID: "S-0011", Type: Story, Status: Review},
			{ID: "S-0012", Type: Story, Status: Done},
		}
		shared := []string{"S-0003", "S-0010", "S-0011"}
		if f.name == "fixture" {
			items = append(items, &Item{ID: "S-0013", Type: Story, Status: fixtureTesting})
			shared = append(shared, "S-0013")
		}
		board := *f.board
		board.Order = []string{"S-0002", "S-0003"}
		for _, it := range items {
			it.Touches = []string{"src/" + it.ID} // no overlap holds
		}
		holds := NewHolds(items, nil)
		for _, it := range items {
			holds.shares = append(holds.shares, Share{Holder: "S-0099", Held: it.ID, At: shareAt})
		}
		v := NewBoardView(items, &board, time.Now(), true, nil, holds)

		if got, want := ids(v.ReadyInPullOrder()), []string{"S-0002", "S-0003", "S-0001"}; !slices.Equal(got, want) {
			t.Errorf("%s: ready in pull order %v, want %v", f.name, got, want)
		}
		if c := v.FirstPullable(); c == nil || c.ID != "S-0002" {
			t.Errorf("%s: first pullable %+v", f.name, c)
		}
		var held, sharedOn []string
		for _, col := range v.Columns {
			for _, c := range col {
				if c.Held != nil {
					held = append(held, c.ID)
				}
				if c.Shared != nil {
					sharedOn = append(sharedOn, c.ID)
				}
			}
		}
		slices.Sort(held)
		slices.Sort(sharedOn)
		if want := []string{"S-0001"}; !slices.Equal(held, want) {
			t.Errorf("%s: held %v, want %v: only a story in the pull's source is judged", f.name, held, want)
		}
		if want := append([]string{"S-0001", "S-0002"}, shared...); !slices.Equal(sharedOn, slices.Sorted(slices.Values(want))) {
			t.Errorf("%s: shared on %v, want %v", f.name, sharedOn, want)
		}
	}
}

// S-0739: LoadBoard gives a project running the fixture the default limits
// of its limited states by kind, doing's and testing's among them, and the
// board takes a limit for testing and none for backlog or for a state the
// workflow does not have; Save writes the limit of every limited state.
func TestBoardLimitsFollowTheWorkflow(t *testing.T) {
	if want := map[string]int{Ready: 5, InProgress: 2, Review: 3}; !reflect.DeepEqual(DefaultWIPLimits, want) {
		t.Errorf("default limits %v, want %v", DefaultWIPLimits, want)
	}
	if want := []string{Ready, InProgress, Review}; !slices.Equal(LimitedColumns, want) {
		t.Errorf("limited columns %v, want %v", LimitedColumns, want)
	}
	r := withTestFlow(t, newProject(t))
	b, err := r.LoadBoard()
	if err != nil {
		t.Fatal(err)
	}
	if want := map[string]int{Ready: 5, fixtureDoing: 2, fixtureTesting: 2, Review: 3}; !reflect.DeepEqual(b.WIPLimits, want) {
		t.Errorf("limits %v, want %v", b.WIPLimits, want)
	}
	if err := b.SetLimit(fixtureTesting, 4); err != nil {
		t.Errorf("set testing: %v", err)
	}
	for _, st := range []string{Backlog, InProgress} {
		err := b.SetLimit(st, 1)
		if want := st + " has no WIP limit; only ready, doing, testing, review do"; err == nil || err.Error() != want {
			t.Errorf("set %s: %v, want %q", st, err, want)
		}
	}
	if err := b.Save("2026-11-14"); err != nil {
		t.Fatal(err)
	}
	back, err := r.LoadBoard()
	if err != nil {
		t.Fatal(err)
	}
	if want := map[string]int{Ready: 5, fixtureDoing: 2, fixtureTesting: 4, Review: 3}; !reflect.DeepEqual(back.WIPLimits, want) {
		t.Errorf("limits read back %v, want %v", back.WIPLimits, want)
	}
	if back.Flow().Kind(fixtureTesting) != "active" {
		t.Error("the board read back does not run the repository's workflow")
	}
}
