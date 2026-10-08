package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/messages"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// messageLint is the repository's markdown lint configuration, so that a
// conversation the lint refuses fails the test.
const messageLint = "default: true\nMD013: false\nMD022:\n  lines_below: 0\nMD024:\n  siblings_only: true\nMD025:\n  front_matter_title: \"\"\nMD032: false\nMD033: false\nMD041: false\nMD060: false\n"

// messageView is the part of a conversation's --json the tests read.
type messageView struct {
	ID       string   `json:"id"`
	Title    string   `json:"title"`
	From     string   `json:"from"`
	To       string   `json:"to"`
	About    []string `json:"about"`
	Status   string   `json:"status"`
	Closed   bool     `json:"closed"`
	Reason   string   `json:"closed_reason"`
	Awaiting string   `json:"awaiting"`
	Path     string   `json:"path"`
	Entries  []struct {
		At     string `json:"at"`
		Author string `json:"author"`
		Story  string `json:"story"`
		Text   string `json:"text"`
	} `json:"entries"`
}

// messageProject builds a project with the lint configuration, a file to
// talk about, and S-0001 and S-0002 in progress, S-0003 in review, S-0004
// ready, and S-0005 done. The environment names no story and no agent.
func messageProject(t *testing.T) (string, *workitem.Repo) {
	t.Helper()
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	t.Setenv("FLAI_AGENT", "")
	t.Setenv("FLAI_STORY", "")
	root := tempProject(t)
	if err := os.MkdirAll(filepath.Join(root, "design", "system"), 0o755); err != nil {
		t.Fatal(err)
	}
	for path, content := range map[string]string{".markdownlint.yaml": messageLint, "design/system/plan.md": "# Plan\n"} {
		if err := os.WriteFile(filepath.Join(root, path), []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	for _, state := range []string{workitem.InProgress, workitem.InProgress, workitem.Review, workitem.Ready, workitem.Done} {
		s, err := repo.Create(workitem.NewOptions{Type: workitem.Story, Title: "Story " + state, Owner: "t", Now: time.Now()})
		if err != nil {
			t.Fatal(err)
		}
		setStoryStatus(t, repo, s.ID, state)
	}
	return root, repo
}

// setStoryStatus writes a story's status without the rules a move applies.
func setStoryStatus(t *testing.T, repo *workitem.Repo, id, state string) {
	t.Helper()
	s, err := repo.Get(id)
	if err != nil {
		t.Fatal(err)
	}
	s.Status = state
	if err := repo.Save(s); err != nil {
		t.Fatal(err)
	}
}

func decodeMessage(t *testing.T, out string) messageView {
	t.Helper()
	var v messageView
	if err := json.Unmarshal([]byte(out), &v); err != nil {
		t.Fatalf("--json: %v\n%s", err, out)
	}
	return v
}

// S-0330, ADR-0120: flai message send starts a conversation from the --from
// story, else FLAI_STORY, else the story in an agent-S-nnnn FLAI_AGENT, and
// prints it, or its view with --json.
func TestMessageSendStartsAConversation(t *testing.T) {
	root, _ := messageProject(t)

	out, errOut, code := runIn(t, root, "message", "send", "S-0003", "Will you leave root.go to me\n\nUntil I push.", "--from", "S-0001", "--about", "design/system/plan.md", "--by", "agent-S-0001")
	if code != 0 {
		t.Fatalf("send: %d %s", code, errOut)
	}
	if !strings.HasPrefix(out, "MS-0001 Will you leave root.go to me\n  S-0001 → S-0003, about design/system/plan.md\n  wip/messages/MS-0001-") || !strings.HasSuffix(out, ".md\n") {
		t.Errorf("send:\n%s", out)
	}

	t.Setenv("FLAI_STORY", "S-0002")
	out, errOut, code = runIn(t, root, "message", "send", "S-0003", "Who writes the docs row", "--by", "agent-S-0002", "--json")
	if code != 0 {
		t.Fatalf("send under FLAI_STORY: %d %s", code, errOut)
	}
	v := decodeMessage(t, out)
	if v.ID != "MS-0002" || v.From != "S-0002" || v.To != "S-0003" || v.Status != "open" || v.Closed || v.Awaiting != "S-0003" ||
		len(v.About) != 0 || len(v.Entries) != 1 || v.Entries[0].Author != "agent-S-0002" || v.Entries[0].Story != "S-0002" ||
		v.Entries[0].At != "2026-09-15T21:00:00Z" || v.Entries[0].Text != "Who writes the docs row" || !strings.HasPrefix(v.Path, "wip/messages/MS-0002-") {
		t.Errorf("send --json: %+v", v)
	}

	t.Setenv("FLAI_STORY", "")
	t.Setenv("FLAI_AGENT", "agent-S-0003")
	out, errOut, code = runIn(t, root, "message", "send", "S-0001", "And the reference", "--json")
	if code != 0 {
		t.Fatalf("send under FLAI_AGENT: %d %s", code, errOut)
	}
	if v := decodeMessage(t, out); v.From != "S-0003" || v.To != "S-0001" || v.Entries[0].Author != "agent-S-0003" {
		t.Errorf("send under FLAI_AGENT: %+v", v)
	}
}

// A send to or from a story not in progress or in review, or with no sender
// story, is refused with the reason, and nothing is written.
func TestMessageSendRefusals(t *testing.T) {
	root, _ := messageProject(t)
	for _, tc := range []struct {
		name string
		args []string
		want string
	}{
		{"to a ready story", []string{"S-0004", "Hello", "--from", "S-0001"}, "S-0004 is ready; a message goes only between stories in progress or in review"},
		{"from a done story", []string{"S-0001", "Hello", "--from", "S-0005"}, "S-0005 is done; a message goes only between stories in progress or in review"},
		{"to itself", []string{"S-0001", "Hello", "--from", "S-0001"}, "S-0001 cannot message itself"},
		{"from no story", []string{"S-0001", "Hello"}, "no story to send from: give --from S-nnnn"},
	} {
		args := append([]string{"message", "send", "--by", "tester"}, tc.args...)
		if _, errOut, code := runIn(t, root, args...); code == 0 || !strings.Contains(errOut, tc.want) {
			t.Errorf("%s: exit %d, %q, want it refused with %q", tc.name, code, errOut, tc.want)
		}
	}
	if matches, _ := filepath.Glob(filepath.Join(root, "wip", "messages", "*")); len(matches) > 0 {
		t.Errorf("a refused send wrote %v", matches)
	}
}

// flai message reply makes the conversation await the other story; list
// shows the open ones, a story's with --story, and the closed ones with
// --all; show prints every entry; each has --json.
func TestMessageReplyListAndShow(t *testing.T) {
	root, repo := messageProject(t)
	t1 := time.Date(2026, 9, 15, 22, 0, 0, 0, time.UTC)
	for _, args := range [][]string{
		{"message", "send", "S-0003", "Will you leave root.go to me", "--from", "S-0001", "--about", "design/system/plan.md", "--by", "agent-S-0001"},
		{"message", "send", "S-0001", "Who writes the docs row", "--from", "S-0002", "--by", "agent-S-0002"},
	} {
		if _, errOut, code := runIn(t, root, args...); code != 0 {
			t.Fatalf("%v: %d %s", args, code, errOut)
		}
	}

	out, errOut, code := runInAt(t, root, t1, "message", "reply", "ms-1", "Yes, until you push", "--from", "S-0003", "--by", "agent-S-0003")
	if code != 0 || out != "MS-0001 awaits S-0001 (2 entries)\n" {
		t.Fatalf("reply: %d %q %s", code, out, errOut)
	}
	if _, errOut, code := runIn(t, root, "message", "reply", "MS-0001", "Me too", "--from", "S-0002", "--by", "agent-S-0002"); code == 0 || !strings.Contains(errOut, "S-0002 is not in MS-0001") {
		t.Errorf("a reply from a third story: exit %d, %q", code, errOut)
	}
	out, errOut, code = runInAt(t, root, t1, "message", "reply", "MS-0002", "I do", "--from", "S-0001", "--by", "agent-S-0001", "--json")
	if v := decodeMessage(t, out); code != 0 || v.Awaiting != "S-0002" || len(v.Entries) != 2 {
		t.Fatalf("reply --json: %d %+v %s", code, v, errOut)
	}

	out, _, _ = runIn(t, root, "message", "list")
	want := "MS-0001  S-0001 → S-0003  awaiting S-0001  2026-09-15T22:00:00Z  Will you leave root.go to me\n" +
		"MS-0002  S-0002 → S-0001  awaiting S-0002  2026-09-15T22:00:00Z  Who writes the docs row\n"
	if out != want {
		t.Errorf("list:\n%s\nwant:\n%s", out, want)
	}
	if out, _, _ = runIn(t, root, "message", "list", "--story", "s-3"); !strings.HasPrefix(out, "MS-0001 ") || strings.Contains(out, "MS-0002") {
		t.Errorf("list --story s-3:\n%s", out)
	}

	setStoryStatus(t, repo, "S-0002", workitem.Done)
	if out, _, _ = runIn(t, root, "message", "list"); !strings.HasPrefix(out, "MS-0001 ") || strings.Contains(out, "MS-0002") {
		t.Errorf("list leaves the closed out:\n%s", out)
	}
	if out, _, _ = runIn(t, root, "message", "list", "--all"); !strings.Contains(out, "MS-0002  S-0002 → S-0001  closed: S-0002 was accepted  2026-09-15T22:00:00Z  Who writes the docs row\n") {
		t.Errorf("list --all:\n%s", out)
	}
	if out, _, _ = runIn(t, root, "message", "list", "--story", "S-0002"); out != "no open conversations for S-0002\n" {
		t.Errorf("list --story S-0002: %q", out)
	}
	if out, _, _ = runIn(t, root, "message", "list", "--story", "S-0004", "--json"); out != "[]\n" {
		t.Errorf("an empty list --json: %q", out)
	}
	out, _, _ = runIn(t, root, "message", "list", "--all", "--json")
	var rows []messageView
	if err := json.Unmarshal([]byte(out), &rows); err != nil || len(rows) != 2 || rows[0].Awaiting != "S-0001" || !rows[1].Closed || rows[1].Reason != "S-0002 was accepted" || rows[1].Awaiting != "" {
		t.Errorf("list --all --json: %v %+v", err, rows)
	}

	out, _, _ = runIn(t, root, "message", "show", "MS-0001")
	want = "MS-0001 Will you leave root.go to me\n  between S-0001 and S-0003 · open · awaiting S-0001\n  about design/system/plan.md\n" +
		"\n2026-09-15T21:00:00Z agent-S-0001 S-0001\nWill you leave root.go to me\n" +
		"\n2026-09-15T22:00:00Z agent-S-0003 S-0003\nYes, until you push\n"
	if out != want {
		t.Errorf("show:\n%s\nwant:\n%s", out, want)
	}
	if out, _, _ = runIn(t, root, "message", "show", "MS-0002"); !strings.Contains(out, "  between S-0002 and S-0001 · closed: S-0002 was accepted\n") {
		t.Errorf("show of a closed conversation:\n%s", out)
	}
	out, _, _ = runIn(t, root, "message", "show", "ms-1", "--json")
	if v := decodeMessage(t, out); v.ID != "MS-0001" || len(v.Entries) != 2 || v.Entries[1].Story != "S-0003" || v.About[0] != "design/system/plan.md" {
		t.Errorf("show --json: %+v", v)
	}
	if _, errOut, code := runIn(t, root, "message", "show", "MS-0009"); code == 0 || !strings.Contains(errOut, "conversation MS-0009 not found") {
		t.Errorf("show of an unknown ID: exit %d, %q", code, errOut)
	}
}

// S-0332, ADR-0121: flai message escalate opens a thread on the escalating
// story, mirrored into its narrative, records it in the conversation, which
// then awaits the other story, and prints both, or both with --json; it is
// refused from a third story, with an empty reason, and on a closed
// conversation.
func TestMessageEscalate(t *testing.T) {
	root, repo := messageProject(t)
	narrative := filepath.Join(root, "wip", "agents", "S-0003.md")
	if err := os.WriteFile(narrative, []byte("# S-0003\n\n## Open questions\n\nNone.\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	for _, args := range [][]string{
		{"message", "send", "S-0003", "Will you leave plan.md to me", "--from", "S-0001", "--about", "design/system/plan.md", "--by", "agent-S-0001"},
		{"message", "send", "S-0001", "Who writes the docs row", "--from", "S-0002", "--by", "agent-S-0002"},
	} {
		if _, errOut, code := runIn(t, root, args...); code != 0 {
			t.Fatalf("%v: %d %s", args, code, errOut)
		}
	}
	t1 := time.Date(2026, 9, 15, 22, 0, 0, 0, time.UTC)

	out, errOut, code := runInAt(t, root, t1, "message", "escalate", "ms-1", "We both need plan.md this week", "--from", "S-0003", "--by", "agent-S-0003")
	if code != 0 {
		t.Fatalf("escalate: %d %s", code, errOut)
	}
	if !strings.HasPrefix(out, "TH-0001 S-0003 and S-0001 do not agree: Will you leave plan.md to me\n  wip/threads/TH-0001-") || !strings.HasSuffix(out, ".md\nMS-0001 awaits S-0001 (2 entries)\n") {
		t.Errorf("escalate:\n%s", out)
	}
	if data, _ := os.ReadFile(narrative); !strings.Contains(string(data), "- TH-0001 (open, agent-S-0003, 2026-09-15): S-0003 and S-0001 do not agree") {
		t.Errorf("the thread is mirrored into the escalating story's narrative:\n%s", data)
	}
	out, _, _ = runIn(t, root, "message", "show", "MS-0001", "--json")
	if v := decodeMessage(t, out); v.Status != "open" || v.Awaiting != "S-0001" || len(v.Entries) != 2 || v.Entries[1].Story != "S-0003" || !strings.Contains(v.Entries[1].Text, "Asked the operator on TH-0001") {
		t.Errorf("the conversation records the escalation and stays open: %+v", v)
	}

	t.Setenv("FLAI_STORY", "S-0001")
	out, errOut, code = runInAt(t, root, t1, "message", "escalate", "MS-0002", "Neither of us will write it", "--by", "agent-S-0001", "--json")
	if code != 0 {
		t.Fatalf("escalate --json: %d %s", code, errOut)
	}
	var got struct {
		Thread struct {
			ID     string `json:"id"`
			Title  string `json:"title"`
			Status string `json:"status"`
			Story  string `json:"story"`
			Anchor struct {
				Item string `json:"item"`
			} `json:"anchor"`
			Entries []struct {
				Author string `json:"author"`
				Text   string `json:"text"`
			} `json:"entries"`
		} `json:"thread"`
		Conversation messageView `json:"conversation"`
	}
	if err := json.Unmarshal([]byte(out), &got); err != nil {
		t.Fatalf("--json: %v\n%s", err, out)
	}
	th, conv := got.Thread, got.Conversation
	if th.ID != "TH-0002" || th.Title != "S-0001 and S-0002 do not agree: Who writes the docs row" || th.Status != "open" || th.Anchor.Item != "S-0001" || th.Story != "S-0001" || len(th.Entries) != 1 || th.Entries[0].Author != "agent-S-0001" {
		t.Errorf("escalate --json thread: %+v", th)
	}
	if len(th.Entries) == 1 && (!strings.Contains(th.Entries[0].Text, "S-0001 and S-0002 do not agree in their conversation MS-0002, `"+conv.Path+"`") || !strings.Contains(th.Entries[0].Text, "> Neither of us will write it")) {
		t.Errorf("the thread names both stories, links the conversation, and quotes the reason:\n%s", th.Entries[0].Text)
	}
	if conv.ID != "MS-0002" || conv.Awaiting != "S-0002" || conv.Closed || len(conv.Entries) != 2 || !strings.Contains(conv.Entries[1].Text, "Asked the operator on TH-0002") {
		t.Errorf("escalate --json conversation: %+v", conv)
	}

	setStoryStatus(t, repo, "S-0002", workitem.Done)
	for _, tc := range []struct {
		name string
		args []string
		want string
	}{
		{"third story", []string{"MS-0001", "No", "--from", "S-0002"}, "S-0002 is not in MS-0001"},
		{"empty reason", []string{"MS-0001", " ", "--from", "S-0001"}, "an escalation needs a reason"},
		{"closed", []string{"MS-0002", "No", "--from", "S-0001"}, "MS-0002 is closed (S-0002 was accepted) and cannot be escalated"},
	} {
		args := append([]string{"message", "escalate", "--by", "tester"}, tc.args...)
		if _, errOut, code := runIn(t, root, args...); code == 0 || !strings.Contains(errOut, tc.want) {
			t.Errorf("%s: exit %d, %q, want it refused with %q", tc.name, code, errOut, tc.want)
		}
	}
	if matches, _ := filepath.Glob(filepath.Join(root, "wip", "threads", "*.md")); len(matches) != 2 {
		t.Errorf("a refused escalation opens no thread: %v", matches)
	}
}

// shareFixture is messageProject with claims, S-0001 in progress holding
// S-0004, ready, on design/system/plan.md, and flai's ask in MS-0001.
func shareFixture(t *testing.T) (string, *workitem.Repo) {
	t.Helper()
	root, repo := messageProject(t)
	for id, touches := range map[string][]string{"S-0001": {"design/system"}, "S-0002": {"flai"}, "S-0004": {"design/system/plan.md"}} {
		s, err := repo.Get(id)
		if err != nil {
			t.Fatal(err)
		}
		s.Touches = touches
		if err := repo.Save(s); err != nil {
			t.Fatal(err)
		}
	}
	if _, _, err := messages.AskHold(repo, messages.AskOptions{Held: "S-0004", Holder: "S-0001", Paths: []string{"design/system/plan.md"}, Now: time.Now()}); err != nil {
		t.Fatal(err)
	}
	return root, repo
}

// S-0334, ADR-0134: flai message share records the share from the holding
// story's agent, prints it, and the held story is no longer held.
func TestMessageShare(t *testing.T) {
	root, repo := shareFixture(t)
	out, errOut, code := runIn(t, root, "message", "share", "MS-0001", "--paths", "design/system/plan.md", "S-0001 edits the first section; S-0004 adds one at the end.", "--from", "S-0001", "--by", "agent-S-0001")
	if code != 0 {
		t.Fatalf("share: %d %s", code, errOut)
	}
	if out != "MS-0001: S-0001 shares design/system/plan.md with S-0004\nMS-0001 awaits S-0004 (2 entries)\n" {
		t.Errorf("share:\n%s", out)
	}
	items, err := repo.List(false)
	if err != nil {
		t.Fatal(err)
	}
	held, err := repo.Get("S-0004")
	if err != nil {
		t.Fatal(err)
	}
	if h := repo.Holds(items).Of(held); h != nil {
		t.Errorf("S-0004 is still held: %+v", h)
	}

	out, errOut, code = runIn(t, root, "message", "share", "ms-1", "--paths", "design/system/plan.md,design/system", "Apart.", "--from", "S-0001", "--by", "agent-S-0001", "--json")
	if code != 0 {
		t.Fatalf("share --json: %d %s", code, errOut)
	}
	var v struct {
		ID     string           `json:"id"`
		Shares []workitem.Share `json:"shares"`
	}
	if err := json.Unmarshal([]byte(out), &v); err != nil {
		t.Fatalf("--json: %v\n%s", err, out)
	}
	if v.ID != "MS-0001" || len(v.Shares) != 2 || strings.Join(v.Shares[1].Paths, ",") != "design/system/plan.md,design/system" || v.Shares[1].By != "agent-S-0001" || v.Shares[1].Split != "Apart." {
		t.Errorf("share --json: %+v", v)
	}
}

// S-0334, ADR-0134: flai message share is refused from the held story, and
// for a path outside the overlap, writing nothing.
func TestMessageShareRefusals(t *testing.T) {
	root, repo := shareFixture(t)
	for _, tc := range []struct {
		args []string
		want string
	}{
		{[]string{"--paths", "design/system/plan.md", "Mine.", "--from", "S-0004", "--by", "agent-S-0004"}, "agent-S-0004, writing for S-0004, may not share S-0001's paths with S-0004"},
		{[]string{"--paths", "flai", "Mine.", "--from", "S-0001", "--by", "agent-S-0001"}, "--paths flai lies outside the overlap"},
	} {
		args := append([]string{"message", "share", "MS-0001"}, tc.args...)
		if _, errOut, code := runIn(t, root, args...); code == 0 || !strings.Contains(errOut, tc.want) {
			t.Errorf("%v: %d %s, want %q", tc.args, code, errOut, tc.want)
		}
	}
	if shares := repo.Shares(); len(shares) != 0 {
		t.Errorf("a refused share records nothing: %+v", shares)
	}
}
