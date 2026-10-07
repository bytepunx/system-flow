package messages

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/mdlint"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

var t0 = time.Date(2026, 10, 7, 9, 0, 0, 0, time.UTC)

// lintConfig is the repository's markdown lint configuration, so that a
// write the lint refuses fails the test.
const lintConfig = "default: true\nMD013: false\nMD022:\n  lines_below: 0\nMD024:\n  siblings_only: true\nMD025:\n  front_matter_title: \"\"\nMD032: false\nMD033: false\nMD041: false\nMD060: false\n"

// project builds a repository with the lint configuration, a file to talk
// about, and one story in each state: S-0001 and S-0002 in progress, S-0003
// in review, S-0004 in the backlog, S-0005 ready, S-0006 done, S-0007
// cancelled, and S-0008 done and archived.
func project(t *testing.T) *workitem.Repo {
	t.Helper()
	root := t.TempDir()
	write(t, filepath.Join(root, "system-flow.yaml"), "version: 1\nname: t\nkey: t\nlayout:\n  design: design\n  docs: docs\n  wip: wip\n")
	for _, d := range []string{"design/system", "wip/kanban/epics", "wip/kanban/stories", "wip/kanban/tasks", "wip/agents", "wip/archive"} {
		if err := os.MkdirAll(filepath.Join(root, d), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	write(t, filepath.Join(root, ".markdownlint.yaml"), lintConfig)
	write(t, filepath.Join(root, "design/system/plan.md"), "# Plan\n")
	r, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	for _, state := range []string{workitem.InProgress, workitem.InProgress, workitem.Review, workitem.Backlog, workitem.Ready, workitem.Done, workitem.Cancelled, workitem.Done} {
		it, err := r.Create(workitem.NewOptions{Type: workitem.Story, Title: "Story " + state, Owner: "a", Now: t0})
		if err != nil {
			t.Fatal(err)
		}
		setStatus(t, r, it.ID, state)
	}
	archive(t, r, "S-0008")
	return r
}

func write(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

// setStatus writes a story's status as a move leaves it, without the rules a
// move applies, which these tests do not exercise.
func setStatus(t *testing.T, r *workitem.Repo, id, state string) {
	t.Helper()
	it, err := r.Get(id)
	if err != nil {
		t.Fatal(err)
	}
	it.Status = state
	if err := r.Save(it); err != nil {
		t.Fatal(err)
	}
}

// archive moves a story's file into wip/archive, as flai archive does.
func archive(t *testing.T, r *workitem.Repo, id string) {
	t.Helper()
	it, err := r.Get(id)
	if err != nil {
		t.Fatal(err)
	}
	dest := filepath.Join(r.ItemDir(workitem.Story, true), filepath.Base(it.Path))
	if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Rename(it.Path, dest); err != nil {
		t.Fatal(err)
	}
}

func readFile(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}

// lintClean fails the test when the conversation's file has a lint finding.
func lintClean(t *testing.T, r *workitem.Repo, c *Conversation) {
	t.Helper()
	cfg, err := mdlint.Load(r.Root)
	if err != nil {
		t.Fatal(err)
	}
	if f := cfg.Lint(readFile(t, c.Path)); len(f) > 0 {
		t.Fatalf("lint findings %+v in:\n%s", f, readFile(t, c.Path))
	}
}

func ids(list []*Conversation) string {
	var out []string
	for _, c := range list {
		out = append(out, c.ID)
	}
	return strings.Join(out, " ")
}

func send(t *testing.T, r *workitem.Repo, from, to, text string, at time.Time) *Conversation {
	t.Helper()
	c, err := Send(r, SendOptions{From: from, To: to, Author: "agent-" + from, Text: text, Now: at})
	if err != nil {
		t.Fatal(err)
	}
	return c
}

// ADR-0120: a message between two stories in progress or in review starts a
// conversation in its own file, numbered one more than the highest.
func TestSendWritesTheConversation(t *testing.T) {
	r := project(t)
	c, err := Send(r, SendOptions{From: "s-1", To: "S-003", Author: "agent-S-0001", Text: "I am changing the board's columns. Will you leave flai/internal/workitem/board.go to me until Friday, please?\n\nIt is a small change.", About: []string{"./design/system/plan.md", "design/system/", "design/system/plan.md"}, Now: t0})
	if err != nil {
		t.Fatal(err)
	}
	want := "---\nid: MS-0001\ntitle: I am changing the board's columns. Will you leave\nfrom: S-0001\nto: S-0003\nabout: [design/system/plan.md, design/system]\nstatus: open\nparticipants: [agent-S-0001]\ncreated: 2026-10-07T09:00:00Z\nupdated: 2026-10-07T09:00:00Z\n---\n\n# MS-0001 I am changing the board's columns. Will you leave\n\nBetween S-0001 and S-0003, about `design/system/plan.md`, `design/system`.\n\n## Entries\n\n### 2026-10-07T09:00:00Z agent-S-0001 S-0001\nI am changing the board's columns. Will you leave flai/internal/workitem/board.go to me until Friday, please?\n\nIt is a small change.\n"
	if got := readFile(t, c.Path); got != want {
		t.Fatalf("file:\n%s\nwant:\n%s", got, want)
	}
	if filepath.Base(c.Path) != "MS-0001-i-am-changing-the-board-s-columns-will-you-leave.md" || filepath.Dir(c.Path) != Dir(r) {
		t.Errorf("path %s", c.Path)
	}
	lintClean(t, r, c)
	if err := c.Validate(); err != nil {
		t.Error(err)
	}
	if c.Awaiting() != "S-0003" {
		t.Errorf("a new conversation awaits the addressee: %s", c.Awaiting())
	}
	if closed, why := c.Closed(r); closed {
		t.Errorf("open conversation reads as closed: %s", why)
	}
	second := send(t, r, "S-0003", "S-0002", "Who owns the cache?", t0)
	if second.ID != "MS-0002" || second.Title != "Who owns the cache?" || second.About != nil || strings.Contains(readFile(t, second.Path), "about:") {
		t.Errorf("second: %+v", second)
	}
	lintClean(t, r, second)
}

// ADR-0120: a send is refused, writing nothing, with the reason, when a side
// is not a story in progress or in review, is the other side, or names a
// path that does not exist, and when the text or author is missing.
func TestSendRefusals(t *testing.T) {
	r := project(t)
	for _, tc := range []struct {
		name string
		opt  SendOptions
		want string
	}{
		{"empty text", SendOptions{From: "S-0001", To: "S-0002", Author: "a", Text: " \n"}, "needs text"},
		{"empty author", SendOptions{From: "S-0001", To: "S-0002", Author: " ", Text: "Hi"}, "an author is required (--by, FLAI_AGENT, or config author)"},
		{"same story", SendOptions{From: "S-0001", To: "s-1", Author: "a", Text: "Hi"}, "S-0001 cannot message itself"},
		{"no sender", SendOptions{To: "S-0002", Author: "a", Text: "Hi"}, "sender's story is required"},
		{"unknown story", SendOptions{From: "S-0001", To: "S-0099", Author: "a", Text: "Hi"}, "S-0099, the addressee, does not exist"},
		{"not a story", SendOptions{From: "T-0001", To: "S-0002", Author: "a", Text: "Hi"}, `"T-0001", the sender, is not a story ID`},
		{"not an ID", SendOptions{From: "S-0001", To: "S-x", Author: "a", Text: "Hi"}, `"S-x", the addressee, is not a story ID`},
		{"backlog", SendOptions{From: "S-0001", To: "S-0004", Author: "a", Text: "Hi"}, "S-0004 is in the backlog; a message goes only between stories in progress or in review"},
		{"ready", SendOptions{From: "S-0005", To: "S-0001", Author: "a", Text: "Hi"}, "S-0005 is ready; a message goes only between stories in progress or in review"},
		{"done", SendOptions{From: "S-0001", To: "S-0006", Author: "a", Text: "Hi"}, "S-0006 is done; a message goes only between stories in progress or in review"},
		{"cancelled", SendOptions{From: "S-0001", To: "S-0007", Author: "a", Text: "Hi"}, "S-0007 is cancelled; "},
		{"archived", SendOptions{From: "S-0001", To: "S-0008", Author: "a", Text: "Hi"}, "S-0008 is archived; "},
		{"missing about", SendOptions{From: "S-0001", To: "S-0002", Author: "a", Text: "Hi", About: []string{"design/nope.md"}}, "--about design/nope.md does not exist in the repository"},
		{"about outside", SendOptions{From: "S-0001", To: "S-0002", Author: "a", Text: "Hi", About: []string{"../x"}}, "not a path inside the repository"},
		{"lint", SendOptions{From: "S-0001", To: "S-0002", Author: "a", Text: "Hi\n\n**Bold as a heading**"}, "MD036"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			tc.opt.Now = t0
			if _, err := Send(r, tc.opt); err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Errorf("got %v, want an error containing %q", err, tc.want)
			}
		})
	}
	if list, err := List(r); err != nil || len(list) != 0 {
		t.Errorf("a refused send writes nothing, and a missing folder lists none: %v %v", ids(list), err)
	}
}

// ADR-0120: a reply from one of the two stories appends an entry, names its
// author among the participants, and awaits the other story; one in the same
// second by the same author for the same story joins the entry before it.
func TestReplyAppendsAndAwaitsTheOther(t *testing.T) {
	r := project(t)
	send(t, r, "S-0001", "S-0003", "May I take board.go?", t0)
	at := t0.Add(time.Minute)
	c, err := Reply(r, "ms-1", "S-3", "agent-S-0003", "Yes, after Friday.", at)
	if err != nil {
		t.Fatal(err)
	}
	if c.Awaiting() != "S-0001" || !contains(c.Participants, "agent-S-0003") || c.Updated != "2026-10-07T09:01:00Z" || c.Status != StatusOpen {
		t.Errorf("after the reply: awaiting %s, %+v", c.Awaiting(), c)
	}
	c, err = Reply(r, "MS-0001", "S-0003", "agent-S-0003", "Or now, if you prefer.", at)
	if err != nil {
		t.Fatal(err)
	}
	if n := strings.Count(c.Body, "### 2026-10-07T09:01:00Z agent-S-0003 S-0003"); n != 1 {
		t.Errorf("same heading joins: %d\n%s", n, c.Body)
	}
	c, err = Reply(r, "7", "S-0001", "agent-S-0001", "x", at)
	if err == nil || !strings.Contains(err.Error(), "MS-0007 not found") {
		t.Errorf("an unknown conversation: %v %v", c, err)
	}
	c, err = Reply(r, "MS-0001", "S-0001", "agent-S-0001", "Now, thanks.", at)
	if err != nil {
		t.Fatal(err)
	}
	want := []Entry{
		{At: "2026-10-07T09:00:00Z", Author: "agent-S-0001", Story: "S-0001", Text: "May I take board.go?"},
		{At: "2026-10-07T09:01:00Z", Author: "agent-S-0003", Story: "S-0003", Text: "Yes, after Friday.\n\nOr now, if you prefer."},
		{At: "2026-10-07T09:01:00Z", Author: "agent-S-0001", Story: "S-0001", Text: "Now, thanks."},
	}
	back, err := Read(c.Path)
	if err != nil {
		t.Fatal(err)
	}
	if got := back.Entries(); len(got) != len(want) || got[0] != want[0] || got[1] != want[1] || got[2] != want[2] {
		t.Fatalf("entries: %+v", got)
	}
	if back.Awaiting() != "S-0003" || back.Marshal() != c.Marshal() {
		t.Errorf("awaiting %s; the file round-trips:\n%s", back.Awaiting(), back.Marshal())
	}
	lintClean(t, r, c)

	was := readFile(t, c.Path)
	for _, tc := range []struct{ story, author, text, want string }{
		{"S-0002", "agent-S-0002", "Me too.", "S-0002 is not in MS-0001, which is between S-0001 and S-0003"},
		{"S-0001", "", "Hi", "an author is required"},
		{"S-0001", "a", "", "a reply needs text"},
		{"S-0001", "a", "Hi\n\n\n\nthere", "MD012"},
	} {
		if _, err := Reply(r, "MS-0001", tc.story, tc.author, tc.text, at); err == nil || !strings.Contains(err.Error(), tc.want) {
			t.Errorf("reply %+v: got %v", tc, err)
		}
	}
	var le *mdlint.Error
	if _, err := Reply(r, "MS-0001", "S-0001", "a", "Hi\n\n\n\nthere", at); !errors.As(err, &le) {
		t.Errorf("a reply the lint rejects is an *mdlint.Error: %v", err)
	}
	if readFile(t, c.Path) != was {
		t.Error("a refused reply leaves the conversation as it was")
	}
}

// ADR-0120: a conversation reads as closed, with the reason, once either
// story is done, cancelled, or archived, and takes no reply; a reply from a
// story moved back to ready is refused with its state.
func TestAConversationReadsAsClosedWhenAStoryLeaves(t *testing.T) {
	for _, tc := range []struct {
		name  string
		leave func(t *testing.T, r *workitem.Repo)
		why   string
	}{
		{"done", func(t *testing.T, r *workitem.Repo) { setStatus(t, r, "S-0003", workitem.Done) }, "S-0003 was accepted"},
		{"cancelled", func(t *testing.T, r *workitem.Repo) { setStatus(t, r, "S-0003", workitem.Cancelled) }, "S-0003 was cancelled"},
		{"archived", func(t *testing.T, r *workitem.Repo) { archive(t, r, "S-0003") }, "S-0003 is archived"},
		{"gone", func(t *testing.T, r *workitem.Repo) {
			it, _ := r.Get("S-0003")
			if err := os.Remove(it.Path); err != nil {
				t.Fatal(err)
			}
		}, "S-0003 does not exist"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			r := project(t)
			c := send(t, r, "S-0001", "S-0003", "Hello", t0)
			tc.leave(t, r)
			closed, why := c.Closed(r)
			if !closed || why != tc.why {
				t.Fatalf("reads as closed %v (%s), want %s", closed, why, tc.why)
			}
			v := View(r, c)
			if v["closed"] != true || v["closed_reason"] != tc.why || v["awaiting"] != "" || v["status"] != StatusOpen {
				t.Errorf("view: %+v", v)
			}
			if _, err := Reply(r, c.ID, "S-0001", "agent-S-0001", "Still there?", t0.Add(time.Minute)); err == nil || !strings.Contains(err.Error(), "MS-0001 is closed ("+tc.why+") and takes no reply") {
				t.Errorf("reply to a closed conversation: %v", err)
			}
		})
	}
	r := project(t)
	c := send(t, r, "S-0001", "S-0003", "Hello", t0)
	setStatus(t, r, "S-0001", workitem.Ready)
	if closed, _ := c.Closed(r); closed {
		t.Error("a story moved back to ready leaves the conversation open")
	}
	if _, err := Reply(r, c.ID, "S-0001", "agent-S-0001", "Back.", t0.Add(time.Minute)); err == nil || !strings.Contains(err.Error(), "S-0001 is ready; a message goes only between stories in progress or in review") {
		t.Errorf("a reply from a story no longer in progress: %v", err)
	}
}

// ADR-0120: CloseOn closes each conversation still stored open of the
// stories with a dated entry giving the reason, leaves the others byte for
// byte, and a closed conversation is not closed again or reopened.
func TestCloseOnClosesTheConversationsOfTheStories(t *testing.T) {
	r := project(t)
	send(t, r, "S-0001", "S-0003", "One", t0)
	send(t, r, "S-0002", "S-0001", "Two", t0)
	send(t, r, "S-0002", "S-0003", "Three", t0)
	send(t, r, "S-0003", "S-0001", "Four", t0)
	if _, err := Close(r, "MS-0004", "alex", "settled by hand", t0.Add(time.Minute)); err != nil {
		t.Fatal(err)
	}
	untouched := map[string]string{}
	for _, id := range []string{"MS-0003", "MS-0004"} {
		c, _ := Get(r, id)
		untouched[id] = readFile(t, c.Path)
	}
	at := t0.Add(time.Hour)
	got, err := CloseOn(r, []string{"s-1"}, "flai", "accepted", at)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Join(got, " ") != "MS-0001 MS-0002" {
		t.Fatalf("CloseOn closed %v", got)
	}
	for _, id := range got {
		c, err := Get(r, id)
		if err != nil {
			t.Fatal(err)
		}
		es := c.Entries()
		last := es[len(es)-1]
		if c.Status != StatusClosed || last != (Entry{At: "2026-10-07T10:00:00Z", Author: "flai", Text: "Closed: S-0001 was accepted"}) || !contains(c.Participants, "flai") {
			t.Errorf("%s: status %s, last entry %+v", id, c.Status, last)
		}
		if closed, why := c.Closed(r); !closed || why != "S-0001 was accepted" {
			t.Errorf("%s reads as closed %v (%s)", id, closed, why)
		}
		if c.Awaiting() != c.To {
			t.Errorf("a closing entry has no story and does not change who is awaited: %s", c.Awaiting())
		}
		if err := c.Validate(); err != nil {
			t.Error(err)
		}
		lintClean(t, r, c)
	}
	for id, was := range untouched {
		c, _ := Get(r, id)
		if readFile(t, c.Path) != was {
			t.Errorf("%s must be left byte for byte", id)
		}
	}
	byHand, err := Get(r, "MS-0004")
	if err != nil {
		t.Fatal(err)
	}
	if closed, why := byHand.Closed(r); !closed || why != "settled by hand" {
		t.Errorf("a closed status reads its reason from the closing entry: %v %s", closed, why)
	}
	if _, err := Close(r, "MS-0004", "alex", "again", at); err == nil || !strings.Contains(err.Error(), "closed already") {
		t.Errorf("closing twice: %v", err)
	}
	if got, err := CloseOn(r, []string{"S-0002", "S-0003"}, "flai", "archived", at); err != nil || strings.Join(got, " ") != "MS-0003" {
		t.Fatalf("CloseOn both: %v %v", got, err)
	}
	if c, _ := Get(r, "MS-0003"); !strings.HasSuffix(c.Body, "### 2026-10-07T10:00:00Z flai\nClosed: S-0002 and S-0003 were archived\n") {
		t.Errorf("both stories leaving:\n%s", c.Body)
	}
	if got, err := CloseOn(r, []string{"S-0099"}, "flai", "archived", at); err != nil || len(got) != 0 {
		t.Errorf("no conversation of the stories is no work: %v %v", got, err)
	}
}

// List orders by number, For finds a story's conversations in any padding,
// and Get takes an ID in any padding.
func TestListForAndGet(t *testing.T) {
	r := project(t)
	for i := 0; i < 10; i++ {
		send(t, r, "S-0001", "S-0002", "Message", t0.Add(time.Duration(i)*time.Second))
	}
	send(t, r, "S-0003", "S-0002", "Other", t0)
	all, err := List(r)
	if err != nil || len(all) != 11 || all[9].ID != "MS-0010" || all[10].ID != "MS-0011" {
		t.Fatalf("list: %s %v", ids(all), err)
	}
	if list, _ := For(r, "s-3"); ids(list) != "MS-0011" {
		t.Errorf("For S-0003: %s", ids(list))
	}
	if list, _ := For(r, "S-002"); len(list) != 11 {
		t.Errorf("For S-0002 finds sent and received: %s", ids(list))
	}
	if list, _ := For(r, "S-0004"); len(list) != 0 {
		t.Errorf("For S-0004: %s", ids(list))
	}
	for _, id := range []string{"ms-11", "MS-011", "11", "MS-0011"} {
		if c, err := Get(r, id); err != nil || c.ID != "MS-0011" {
			t.Errorf("Get %s: %v", id, err)
		}
	}
	if CanonicalID("TH-0001") != "TH-0001" || CanonicalID("ms-0") != "MS-0000" {
		t.Errorf("CanonicalID leaves other IDs alone")
	}
	if NextID(r) != "MS-0012" {
		t.Errorf("next: %s", NextID(r))
	}
}

// View is what --json prints: the front matter, the path from the root, who
// is awaited, and the entries.
func TestView(t *testing.T) {
	r := project(t)
	c, err := Send(r, SendOptions{From: "S-0001", To: "S-0002", Author: "agent-S-0001", Text: "Hi", About: []string{"design/system/plan.md"}, Now: t0})
	if err != nil {
		t.Fatal(err)
	}
	v := View(r, c)
	if v["path"] != "wip/messages/MS-0001-hi.md" || v["awaiting"] != "S-0002" || v["closed"] != false || v["closed_reason"] != "" || v["from"] != "S-0001" || v["to"] != "S-0002" {
		t.Errorf("view: %+v", v)
	}
	if es, ok := v["entries"].([]Entry); !ok || len(es) != 1 || es[0].Story != "S-0001" {
		t.Errorf("entries: %+v", v["entries"])
	}
	if about, ok := v["about"].([]string); !ok || len(about) != 1 {
		t.Errorf("about: %+v", v["about"])
	}
	c2 := send(t, r, "S-0002", "S-0001", "Hello", t0)
	if about, ok := View(r, c2)["about"].([]string); !ok || about == nil || len(about) != 0 {
		t.Errorf("no about is an empty list, not null: %#v", View(r, c2)["about"])
	}
}

// A title is the first line of the first message, cleaned and cut on a word.
func TestTitleOf(t *testing.T) {
	for _, tc := range []struct{ text, want string }{
		{"\n\n  Hello there.  \nMore", "Hello there"},
		{"Who owns it?", "Who owns it?"},
		{"...", "From S-0001"},
		{strings.Repeat("word ", 20), "word word word word word word word word word word word word word word"},
		{strings.Repeat("x", 80), strings.Repeat("x", 72)},
		{strings.Repeat("a", 60) + " bcdef, ghijklmnop", strings.Repeat("a", 60) + " bcdef"},
	} {
		if got := titleOf(tc.text, "S-0001"); got != tc.want {
			t.Errorf("titleOf(%q) = %q, want %q", tc.text, got, tc.want)
		}
	}
}

// Parse keeps fields this flai does not know, and Marshal writes them back.
func TestUnknownFieldsRoundTrip(t *testing.T) {
	r := project(t)
	c := send(t, r, "S-0001", "S-0002", "Hi", t0)
	doc := strings.Replace(readFile(t, c.Path), "\n---\n", "\npriority: high\n---\n", 1)
	write(t, c.Path, doc)
	list, err := List(r)
	if err != nil || len(list) != 1 || len(list[0].Unknown) != 1 || list[0].Unknown[0].Name != "priority" || list[0].Marshal() != doc {
		t.Fatalf("listing: %v %+v", err, list)
	}
	if _, err := Reply(r, c.ID, "S-0002", "agent-S-0002", "Yes.", t0.Add(time.Minute)); err != nil {
		t.Fatal(err)
	}
	if after := readFile(t, c.Path); !strings.Contains(after, "priority: high\n---\n") || !strings.Contains(after, "Yes.") {
		t.Errorf("a reply dropped the unknown field:\n%s", after)
	}
}

func TestValidate(t *testing.T) {
	good := func() *Conversation {
		return &Conversation{ID: "MS-0001", Title: "Hi", From: "S-0001", To: "S-0002", About: []string{"design"}, Status: StatusOpen, Participants: []string{"a"}, Created: "2026-10-07T09:00:00Z", Updated: "2026-10-07T09:00:00Z"}
	}
	if err := good().Validate(); err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		edit func(c *Conversation)
		want string
	}{
		{func(c *Conversation) { c.ID = "TH-0001" }, `id "TH-0001" must look like MS-0001`},
		{func(c *Conversation) { c.Title = " " }, "title is required"},
		{func(c *Conversation) { c.From = "T-0001" }, `from "T-0001" must be a story ID`},
		{func(c *Conversation) { c.To = "" }, `to "" must be a story ID`},
		{func(c *Conversation) { c.To = "S-001" }, "from and to are both S-0001"},
		{func(c *Conversation) { c.About = []string{"../out"} }, "about: ../out is not a path inside the repository"},
		{func(c *Conversation) { c.Status = "resolved" }, `status "resolved" must be one of open, closed`},
		{func(c *Conversation) { c.Participants = nil }, "participants must name at least the sender"},
		{func(c *Conversation) { c.Updated = "yesterday" }, `updated "yesterday" is not a UTC timestamp`},
	} {
		c := good()
		tc.edit(c)
		if err := c.Validate(); err == nil || !strings.Contains(err.Error(), tc.want) {
			t.Errorf("got %v, want %q", err, tc.want)
		}
	}
}
