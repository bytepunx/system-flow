package messages

import (
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// shareProject is project with claims: S-0001, in progress, touches
// flai/internal/workitem and design/system/plan.md; S-0005, ready, touches
// flai/internal/workitem/hold.go and docs, so that S-0001 holds it on
// flai/internal/workitem/hold.go alone; S-0002, in progress, touches a path
// of its own. S-0005 entered ready at t0 and has a goal.
func shareProject(t *testing.T) *workitem.Repo {
	t.Helper()
	r := project(t)
	edit := func(id string, f func(it *workitem.Item)) {
		t.Helper()
		it, err := r.Get(id)
		if err != nil {
			t.Fatal(err)
		}
		f(it)
		if err := r.Save(it); err != nil {
			t.Fatal(err)
		}
	}
	edit("S-0001", func(it *workitem.Item) {
		it.Touches = []string{"flai/internal/workitem", "design/system/plan.md"}
	})
	edit("S-0002", func(it *workitem.Item) { it.Touches = []string{"flai/cmd"} })
	edit("S-0005", func(it *workitem.Item) {
		it.Touches = []string{"flai/internal/workitem/hold.go", "docs"}
		it.Transitions = []workitem.Transition{{To: workitem.Ready, At: t0.Format(workitem.TimeFormat), By: "a"}}
		it.Body = "\n## Goal\n\nA held story starts as soon as its holder shares.\n\n## Acceptance criteria\n\n- [ ] works\n"
	})
	return r
}

// move records a transition of a story to state at the time given.
func move(t *testing.T, r *workitem.Repo, id, state string, at time.Time) {
	t.Helper()
	it, err := r.Get(id)
	if err != nil {
		t.Fatal(err)
	}
	it.Status = state
	it.Transitions = append(it.Transitions, workitem.Transition{To: state, At: at.UTC().Format(workitem.TimeFormat), By: "a"})
	if err := r.Save(it); err != nil {
		t.Fatal(err)
	}
}

// heldStory is why S-0005 is held now, as flai board and the launcher read it.
func heldStory(t *testing.T, r *workitem.Repo) *workitem.Hold {
	t.Helper()
	items, err := r.List(false)
	if err != nil {
		t.Fatal(err)
	}
	it, err := r.Get("S-0005")
	if err != nil {
		t.Fatal(err)
	}
	return r.Holds(items).Of(it)
}

func ask(t *testing.T, r *workitem.Repo, at time.Time) (*Conversation, bool) {
	t.Helper()
	c, wrote, err := AskHold(r, AskOptions{Held: "s-5", Holder: "S-1", Paths: []string{"./flai/internal/workitem/hold.go"}, Now: at})
	if err != nil {
		t.Fatal(err)
	}
	return c, wrote
}

// ADR-0134: flai asks the holding story's agent about the hold, from the held
// story, about the paths, with the held story's title and goal and the three
// answers; the conversation awaits the holder.
func TestAskHoldAsksTheHolder(t *testing.T) {
	r := shareProject(t)
	c, wrote := ask(t, r, t0.Add(time.Hour))
	if !wrote || c.ID != "MS-0001" || c.From != "S-0005" || c.To != "S-0001" || strings.Join(c.About, ",") != "flai/internal/workitem/hold.go" {
		t.Fatalf("ask: wrote %v, %+v", wrote, c)
	}
	es := c.Entries()
	if len(es) != 1 || es[0].Author != AskAuthor || es[0].Story != "S-0005" || c.Awaiting() != "S-0001" {
		t.Fatalf("entries %+v, awaiting %s", es, c.Awaiting())
	}
	if c.Title != "S-0001 holds S-0005 on overlap alone: may S-0005 start beside it?" {
		t.Errorf("title: %q", c.Title)
	}
	for _, want := range []string{
		`S-0005, "Story ready", is ready, and its claim overlaps S-0001's on ` + "`flai/internal/workitem/hold.go`",
		"Its goal:\n\n> A held story starts as soon as its holder shares.",
		"`flai touches S-0001 --remove <path>`",
		"`flai message share MS-0001 --paths <path> \"<split>\"`",
		"`message_share`",
		"`message_reply` saying why the hold stands",
		"S-0005 waits until S-0001 moves to review, is cancelled, or is sent back",
	} {
		if !strings.Contains(es[0].Text, want) {
			t.Errorf("the ask does not say %q:\n%s", want, es[0].Text)
		}
	}
	lintClean(t, r, c)
	if err := c.Validate(); err != nil {
		t.Error(err)
	}
}

// ADR-0134: the ask is made once per pair while the held story stays in
// ready, and again once it has left ready and entered it anew.
func TestAskHoldOncePerReadyStint(t *testing.T) {
	r := shareProject(t)
	first, _ := ask(t, r, t0.Add(time.Hour))
	was := readFile(t, first.Path)
	again, wrote := ask(t, r, t0.Add(2*time.Hour))
	if wrote || again.ID != first.ID || readFile(t, first.Path) != was {
		t.Fatalf("a second ask in the same stint writes nothing: wrote %v, %s", wrote, again.ID)
	}
	move(t, r, "S-0005", workitem.Backlog, t0.Add(3*time.Hour))
	move(t, r, "S-0005", workitem.Ready, t0.Add(4*time.Hour))
	next, wrote := ask(t, r, t0.Add(5*time.Hour))
	if !wrote || next.ID != first.ID || len(next.Entries()) != 2 {
		t.Fatalf("a new stint asks again in the open conversation: wrote %v, %+v", wrote, next.Entries())
	}
	lintClean(t, r, next)
	if _, wrote := ask(t, r, t0.Add(6*time.Hour)); wrote {
		t.Error("the new stint is asked about once")
	}
}

// ADR-0134: a held story with no goal is asked about as having none.
func TestAskHoldWithNoGoal(t *testing.T) {
	r := shareProject(t)
	it, _ := r.Get("S-0005")
	it.Body = "\n## Goal\n\n## Acceptance criteria\n\n- [ ] works\n"
	if err := r.Save(it); err != nil {
		t.Fatal(err)
	}
	c, _ := ask(t, r, t0.Add(time.Hour))
	if text := c.Entries()[0].Text; !strings.Contains(text, "S-0005 has no goal written.") || strings.Contains(text, "Its goal") {
		t.Errorf("no goal:\n%s", text)
	}
	lintClean(t, r, c)
}

// ADR-0134: an ask is refused, writing nothing, for a held story that is not
// ready, a holder not in progress, or no path.
func TestAskHoldRefusals(t *testing.T) {
	r := shareProject(t)
	for _, tc := range []struct {
		name string
		opt  AskOptions
		want string
	}{
		{"held in backlog", AskOptions{Held: "S-0004", Holder: "S-0001", Paths: []string{"docs"}}, "S-0004, the held story, is in the backlog; it must be ready"},
		{"held in progress", AskOptions{Held: "S-0002", Holder: "S-0001", Paths: []string{"docs"}}, "S-0002, the held story, is in progress; it must be ready"},
		{"held archived", AskOptions{Held: "S-0008", Holder: "S-0001", Paths: []string{"docs"}}, "S-0008, the held story, is archived"},
		{"holder in review", AskOptions{Held: "S-0005", Holder: "S-0003", Paths: []string{"docs"}}, "S-0003, the holding story, is in review; it must be in progress"},
		{"holder unknown", AskOptions{Held: "S-0005", Holder: "S-0099", Paths: []string{"docs"}}, "S-0099, the holding story, does not exist"},
		{"no path", AskOptions{Held: "S-0005", Holder: "S-0001"}, "needs the paths S-0001 holds it on"},
		{"bad path", AskOptions{Held: "S-0005", Holder: "S-0001", Paths: []string{"../x"}}, "not a path inside the repository"},
	} {
		tc.opt.Now = t0
		if _, wrote, err := AskHold(r, tc.opt); wrote || err == nil || !strings.Contains(err.Error(), tc.want) {
			t.Errorf("%s: wrote %v, got %v, want %q", tc.name, wrote, err, tc.want)
		}
	}
	if list, _ := List(r); len(list) != 0 {
		t.Errorf("a refused ask writes nothing: %s", ids(list))
	}
}

// ADR-0134: the holder answers the ask; its reply awaits no one, since the
// ready story has no agent, so AwaitingOther leaves it out. The ready story
// itself still cannot write.
func TestTheHolderRepliesToTheAsk(t *testing.T) {
	r := shareProject(t)
	c, _ := ask(t, r, t0.Add(time.Hour))
	got, err := Reply(r, c.ID, "S-0001", "agent-S-0001", "I change only the hold reasons; the hold stands until I push.", t0.Add(2*time.Hour))
	if err != nil {
		t.Fatal(err)
	}
	if got.Awaiting() != "S-0005" {
		t.Errorf("the reply awaits the held story: %s", got.Awaiting())
	}
	if list, err := AwaitingOther(r, "S-0001"); err != nil || len(list) != 0 {
		t.Errorf("a reply to a ready story awaits no one: %s %v", ids(list), err)
	}
	if _, err := Reply(r, c.ID, "S-0005", "agent-S-0005", "Thanks.", t0.Add(3*time.Hour)); err == nil || !strings.Contains(err.Error(), "S-0005 is ready; a message goes only between stories in progress or in review") {
		t.Errorf("a ready story's reply: %v", err)
	}
	move(t, r, "S-0005", workitem.InProgress, t0.Add(4*time.Hour))
	if list, _ := AwaitingOther(r, "S-0001"); ids(list) != c.ID {
		t.Errorf("once the held story starts, the reply awaits it: %s", ids(list))
	}
}

// ADR-0134: the holding story's agent, writing for it, or the operator, its
// owner, shares paths inside the overlap with the held story: the share is in
// the front matter, an entry gives the split, and the held story is no longer
// held on them.
func TestShareLiftsTheHold(t *testing.T) {
	for _, tc := range []struct{ name, story, author string }{
		{"the holder's agent", "s-1", "agent-S-0001"},
		{"the operator", "", "a"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			r := shareProject(t)
			if h := heldStory(t, r); h == nil || h.Code != workitem.HoldOverlap {
				t.Fatalf("before the share S-0005 is held on overlap: %+v", h)
			}
			asked, _ := ask(t, r, t0.Add(time.Hour))
			at := t0.Add(2 * time.Hour)
			c, err := Share(r, ShareOptions{ID: "ms-1", Story: tc.story, Author: tc.author, Paths: []string{"./flai/internal/workitem/hold.go"}, Split: "S-0001 changes holdBy; S-0005 adds a function below it.", Now: at})
			if err != nil {
				t.Fatal(err)
			}
			want := workitem.Share{Holder: "S-0001", Held: "S-0005", Paths: []string{"flai/internal/workitem/hold.go"}, Split: "S-0001 changes holdBy; S-0005 adds a function below it.", By: tc.author, At: "2026-10-07T11:00:00Z"}
			if len(c.Shares) != 1 || !sameShare(c.Shares[0], want) {
				t.Fatalf("shares: %+v", c.Shares)
			}
			es := c.Entries()
			last := es[len(es)-1]
			if last.Author != tc.author || last.Story != "S-0001" || !strings.Contains(last.Text, "S-0001 shares `flai/internal/workitem/hold.go` with S-0005, split so:\n\n> S-0001 changes holdBy") || !strings.Contains(last.Text, "no longer holds (ADR-0134)") {
				t.Errorf("the share's entry: %+v", last)
			}
			if !contains(c.Participants, tc.author) || c.Updated != "2026-10-07T11:00:00Z" {
				t.Errorf("participants %v, updated %s", c.Participants, c.Updated)
			}
			back, err := Read(asked.Path)
			if err != nil {
				t.Fatal(err)
			}
			if len(back.Shares) != 1 || !sameShare(back.Shares[0], want) || back.Marshal() != c.Marshal() || len(back.Unknown) != 0 {
				t.Errorf("the file round-trips: %+v\n%s", back.Shares, readFile(t, asked.Path))
			}
			if err := back.Validate(); err != nil {
				t.Error(err)
			}
			lintClean(t, r, back)
			if v, ok := View(r, back)["shares"].([]workitem.Share); !ok || len(v) != 1 {
				t.Errorf("view shares: %#v", View(r, back)["shares"])
			}
			if got := r.Shares(); len(got) != 1 || !sameShare(got[0], want) {
				t.Errorf("r.Shares: %+v", got)
			}
			if h := heldStory(t, r); h != nil {
				t.Errorf("after the share S-0005 is not held: %+v", h)
			}
		})
	}
}

// ADR-0134, ADR-0121: a share changes no claim, so once the held story has
// started, a trial-merge conflict between the two still reaches them, in the
// conversation the ask began, and the share stays on it.
func TestAConflictAfterAShareJoinsTheirConversation(t *testing.T) {
	r := shareProject(t)
	ask(t, r, t0.Add(time.Hour))
	if _, err := Share(r, ShareOptions{ID: "ms-1", Story: "S-0001", Author: "agent-S-0001", Paths: []string{"flai/internal/workitem/hold.go"}, Split: "S-0001 changes holdBy; S-0005 adds a function below it.", Now: t0.Add(2 * time.Hour)}); err != nil {
		t.Fatal(err)
	}
	setStatus(t, r, "S-0005", workitem.InProgress)
	c, err := Notify(r, SendOptions{From: "S-0005", To: "S-0001", Author: "flai", Text: "story/S-0005 and story/S-0001 conflict when merged.", About: []string{"flai/internal/workitem/hold.go"}, Now: t0.Add(3 * time.Hour)})
	if err != nil {
		t.Fatal(err)
	}
	if c.ID != "MS-0001" || len(c.Shares) != 1 || workitem.CanonicalID(c.Awaiting()) != "S-0001" {
		t.Errorf("the conflict should join MS-0001, keep its share, and await S-0001: %s, %+v, awaiting %s", c.ID, c.Shares, c.Awaiting())
	}
}

func sameShare(a, b workitem.Share) bool {
	return a.Holder == b.Holder && a.Held == b.Held && strings.Join(a.Paths, ",") == strings.Join(b.Paths, ",") && a.Split == b.Split && a.By == b.By && a.At == b.At
}

// ADR-0134: a share ends when either story leaves the open columns, or the
// held story goes back to backlog: the overlap holds again.
func TestAShareEndsWhenAStoryLeaves(t *testing.T) {
	for _, tc := range []struct {
		name  string
		leave func(t *testing.T, r *workitem.Repo)
	}{
		{"the holder goes back to backlog and starts again", func(t *testing.T, r *workitem.Repo) {
			move(t, r, "S-0001", workitem.Backlog, t0.Add(3*time.Hour))
			move(t, r, "S-0001", workitem.InProgress, t0.Add(4*time.Hour))
		}},
		{"the held story goes back to backlog and to ready", func(t *testing.T, r *workitem.Repo) {
			move(t, r, "S-0005", workitem.Backlog, t0.Add(3*time.Hour))
			move(t, r, "S-0005", workitem.Ready, t0.Add(4*time.Hour))
		}},
		{"the holder is cancelled and reopened", func(t *testing.T, r *workitem.Repo) {
			move(t, r, "S-0001", workitem.Cancelled, t0.Add(3*time.Hour))
			move(t, r, "S-0001", workitem.InProgress, t0.Add(4*time.Hour))
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			r := shareProject(t)
			ask(t, r, t0.Add(time.Hour))
			if _, err := Share(r, ShareOptions{ID: "MS-0001", Story: "S-0001", Author: "agent-S-0001", Paths: []string{"flai/internal/workitem/hold.go"}, Split: "Apart.", Now: t0.Add(2 * time.Hour)}); err != nil {
				t.Fatal(err)
			}
			if h := heldStory(t, r); h != nil {
				t.Fatalf("shared: %+v", h)
			}
			tc.leave(t, r)
			if h := heldStory(t, r); h == nil || h.Code != workitem.HoldOverlap {
				t.Errorf("the share has ended, so S-0005 is held again: %+v", h)
			}
		})
	}
}

// ADR-0134: a share is refused, writing nothing, from the held side, a third
// story, an author who is neither the holder's agent nor its owner, on a
// closed conversation or one with no ready story, with no split or path, and
// for a path outside the two claims' overlap.
func TestShareRefusals(t *testing.T) {
	r := shareProject(t)
	c, _ := ask(t, r, t0.Add(time.Hour))
	pair := send(t, r, "S-0001", "S-0002", "Who owns flai/cmd?", t0)
	closed, _, err := AskHold(r, AskOptions{Held: "S-0005", Holder: "S-0002", Paths: []string{"flai/cmd"}, Now: t0})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := Close(r, closed.ID, "a", "settled", t0.Add(time.Minute)); err != nil {
		t.Fatal(err)
	}
	before := map[string]string{}
	for _, x := range []*Conversation{c, pair, closed} {
		before[x.Path] = readFile(t, x.Path)
	}
	path := []string{"flai/internal/workitem/hold.go"}
	for _, tc := range []struct {
		name string
		opt  ShareOptions
		want string
	}{
		{"the held side", ShareOptions{ID: c.ID, Story: "S-0005", Author: "agent-S-0005", Paths: path, Split: "x"}, "agent-S-0005, writing for S-0005, may not share S-0001's paths with S-0005: only S-0001's agent, writing for S-0001, or the operator, its owner a, may (ADR-0134)"},
		{"a third story", ShareOptions{ID: c.ID, Story: "S-0002", Author: "agent-S-0002", Paths: path, Split: "x"}, "agent-S-0002, writing for S-0002, may not share"},
		{"a stranger", ShareOptions{ID: c.ID, Author: "stranger", Paths: path, Split: "x"}, "stranger may not share S-0001's paths with S-0005"},
		{"no author", ShareOptions{ID: c.ID, Story: "S-0001", Author: " ", Paths: path, Split: "x"}, "an author is required"},
		{"empty split", ShareOptions{ID: c.ID, Story: "S-0001", Author: "agent-S-0001", Paths: path, Split: " \n"}, "a share needs the split"},
		{"no path", ShareOptions{ID: c.ID, Story: "S-0001", Author: "agent-S-0001", Split: "x"}, "a share needs at least one path"},
		{"outside the holder's claim", ShareOptions{ID: c.ID, Story: "S-0001", Author: "agent-S-0001", Paths: []string{"docs"}, Split: "x"}, "--paths docs lies outside the overlap of S-0005's and S-0001's claims, which is `flai/internal/workitem/hold.go`"},
		{"outside the held claim", ShareOptions{ID: c.ID, Story: "S-0001", Author: "agent-S-0001", Paths: []string{"flai/internal/workitem/hold.go", "flai/internal/workitem/board.go"}, Split: "x"}, "--paths flai/internal/workitem/board.go lies outside the overlap"},
		{"closed", ShareOptions{ID: closed.ID, Story: "S-0002", Author: "agent-S-0002", Paths: []string{"flai/cmd"}, Split: "x"}, closed.ID + " is closed (settled) and takes no share"},
		{"no ready story", ShareOptions{ID: pair.ID, Story: "S-0001", Author: "agent-S-0001", Paths: path, Split: "x"}, "neither S-0001 (in progress) nor S-0002 (in progress)"},
		{"unknown", ShareOptions{ID: "MS-9", Story: "S-0001", Author: "agent-S-0001", Paths: path, Split: "x"}, "MS-0009 not found"},
		{"lint", ShareOptions{ID: c.ID, Story: "S-0001", Author: "agent-S-0001", Paths: path, Split: "one\ttwo"}, "MD010"},
	} {
		tc.opt.Now = t0.Add(2 * time.Hour)
		if _, err := Share(r, tc.opt); err == nil || !strings.Contains(err.Error(), tc.want) {
			t.Errorf("%s: got %v, want %q", tc.name, err, tc.want)
		}
	}
	for p, was := range before {
		if readFile(t, p) != was {
			t.Errorf("a refused share leaves %s as it was", p)
		}
	}
	if len(r.Shares()) != 0 {
		t.Errorf("no share was recorded: %+v", r.Shares())
	}
}

// A path broader than the overlap but claimed by both stories may be shared.
func TestShareAFolderBothClaim(t *testing.T) {
	r := shareProject(t)
	ask(t, r, t0.Add(time.Hour))
	if _, err := Share(r, ShareOptions{ID: "MS-0001", Story: "S-0001", Author: "agent-S-0001", Paths: []string{"flai/internal/workitem/"}, Split: "Apart."}); err != nil {
		t.Fatal(err)
	}
}

// A share kept in the front matter is checked by Validate.
func TestValidateShares(t *testing.T) {
	c := &Conversation{ID: "MS-0001", Title: "Hi", From: "S-0001", To: "S-0002", Status: StatusOpen, Participants: []string{"a"}, Created: "2026-10-07T09:00:00Z", Updated: "2026-10-07T09:00:00Z",
		Shares: []workitem.Share{{Holder: "T-0001", Held: "S-0002", Split: " ", At: "now"}}}
	err := c.Validate()
	for _, want := range []string{`shares[0]: holder "T-0001" must be a story ID`, "shares[0]: paths must name at least one path", "shares[0]: split is required", "shares[0]: by is required", `shares[0]: at "now" is not a UTC timestamp`} {
		if err == nil || !strings.Contains(err.Error(), want) {
			t.Errorf("got %v, want %q", err, want)
		}
	}
}
