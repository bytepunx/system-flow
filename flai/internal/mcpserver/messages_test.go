package mcpserver

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/bytepunx/system-flow/flai/internal/messages"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// messagesFixture is setup with a second story in progress, other, beside
// the fixture's own, and FLAI_STORY cleared, so that a test names the
// session's story itself.
func messagesFixture(t *testing.T) (*fixture, *workitem.Item) {
	t.Helper()
	t.Setenv("FLAI_STORY", "")
	f := setup(t)
	epic, err := f.repo.Get(f.story.Parent)
	if err != nil {
		t.Fatal(err)
	}
	return f, storyOfEpic(t, f.repo, epic, "Other", nil, workitem.Ready, workitem.InProgress)
}

// sessionAs connects another agent's session to f's project.
func sessionAs(t *testing.T, f *fixture, agent string) *fixture {
	t.Helper()
	clock := *f.clock
	cs := connectServer(t, Options{Repo: f.repo, Agent: agent, Version: "test", Now: func() time.Time { return clock }})
	return &fixture{repo: f.repo, cs: cs, story: f.story, task: f.task, clock: &clock}
}

// conversation decodes a message tool's answer.
func conversation(t *testing.T, out map[string]any) ConversationOut {
	t.Helper()
	data, err := json.Marshal(out)
	if err != nil {
		t.Fatal(err)
	}
	var c ConversationOut
	if err := json.Unmarshal(data, &c); err != nil {
		t.Fatal(err)
	}
	return c
}

// conversationFiles lists what wip/messages holds.
func conversationFiles(t *testing.T, repo *workitem.Repo) []string {
	t.Helper()
	matches, err := filepath.Glob(filepath.Join(messages.Dir(repo), "*.md"))
	if err != nil {
		t.Fatal(err)
	}
	return matches
}

// S-0331: message_send starts a conversation from FLAI_STORY, the other
// story's agent, named agent-S-nnnn, answers it with message_reply, and
// message_get reads it as flai message show --json prints it.
func TestMessageToolsWorkAConversationBetweenTwoStories(t *testing.T) {
	f, other := messagesFixture(t)
	t.Setenv("FLAI_STORY", f.story.ID)

	out, failed := f.call(t, "message_send", map[string]any{"to": other.ID, "text": "I am changing plan.md.\nWill you leave it to me?", "about": []string{"design/system/plan.md"}})
	if failed != "" {
		t.Fatal(failed)
	}
	sent := conversation(t, out)
	if sent.ID != "MS-0001" || sent.Title != "I am changing plan.md" || sent.From != f.story.ID || sent.To != other.ID || sent.Status != messages.StatusOpen || sent.Closed || sent.Awaiting != other.ID {
		t.Errorf("message_send should start an open conversation from %s to %s that awaits %s: %+v", f.story.ID, other.ID, other.ID, sent)
	}
	if len(sent.About) != 1 || sent.About[0] != "design/system/plan.md" || !strings.HasPrefix(sent.Path, "wip/messages/MS-0001-") {
		t.Errorf("message_send should keep about and answer the path from the project root: %+v", sent)
	}
	if len(sent.Entries) != 1 || sent.Entries[0].Author != "claude" || sent.Entries[0].Story != f.story.ID || sent.Entries[0].Text != "I am changing plan.md.\nWill you leave it to me?" {
		t.Errorf("the first entry should be the agent's, for its story: %+v", sent.Entries)
	}

	t.Setenv("FLAI_STORY", "")
	theirs := sessionAs(t, f, "agent-"+other.ID)
	out, failed = theirs.call(t, "message_reply", map[string]any{"id": "ms-1", "text": "Yes, it is yours until you push."})
	if failed != "" {
		t.Fatal(failed)
	}
	replied := conversation(t, out)
	if replied.Awaiting != f.story.ID || len(replied.Entries) != 2 || replied.Entries[1].Author != "agent-"+other.ID || replied.Entries[1].Story != other.ID {
		t.Errorf("a reply from %s, named in the agent's name, should await %s: %+v", other.ID, f.story.ID, replied)
	}
	if len(replied.Participants) != 2 || replied.Participants[1] != "agent-"+other.ID {
		t.Errorf("the replying agent should join the participants: %v", replied.Participants)
	}

	out, failed = f.call(t, "message_get", map[string]any{"id": "MS-1"})
	if failed != "" {
		t.Fatal(failed)
	}
	c, err := messages.Get(f.repo, "MS-0001")
	if err != nil {
		t.Fatal(err)
	}
	// both as JSON objects, whose keys json.Marshal writes sorted
	want, _ := json.Marshal(messages.View(f.repo, c))
	got, _ := json.Marshal(out)
	if string(want) != string(got) {
		t.Errorf("message_get should answer what flai message show --json prints:\nwant %s\ngot  %s", want, got)
	}

	if _, failed := f.call(t, "message_get", map[string]any{"id": "MS-9"}); !strings.Contains(failed, "MS-0009 not found") {
		t.Errorf("an unknown conversation should be refused: %q", failed)
	}
}

// S-0331: a session with no story of its own, no FLAI_STORY, an agent name
// not of the form agent-S-nnnn, and no story branch checked out, can neither
// send nor reply, and nothing is written; it may still read a conversation.
func TestAMessageNeedsASessionWithAStory(t *testing.T) {
	f, other := messagesFixture(t)
	if _, failed := f.call(t, "message_send", map[string]any{"to": other.ID, "text": "Hello"}); !strings.Contains(failed, "message_send writes as this session's own story, and it has none") || !strings.Contains(failed, `"claude"`) || !strings.Contains(failed, "nothing was written") {
		t.Errorf("a send with no story should be refused, saying why and what to do: %q", failed)
	}
	if files := conversationFiles(t, f.repo); len(files) != 0 {
		t.Fatalf("a refused send should write nothing: %v", files)
	}

	c, err := messages.Send(f.repo, messages.SendOptions{From: other.ID, To: f.story.ID, Author: "agent-" + other.ID, Text: "Who adds the row?", Now: t0})
	if err != nil {
		t.Fatal(err)
	}
	before, err := os.ReadFile(c.Path)
	if err != nil {
		t.Fatal(err)
	}
	if _, failed := f.call(t, "message_reply", map[string]any{"id": c.ID, "text": "I do."}); !strings.Contains(failed, "message_reply writes as this session's own story, and it has none") {
		t.Errorf("a reply with no story should be refused: %q", failed)
	}
	if after, _ := os.ReadFile(c.Path); string(after) != string(before) {
		t.Errorf("a refused reply should write nothing:\n%s", after)
	}

	out, failed := f.call(t, "message_get", map[string]any{"id": c.ID})
	if failed != "" {
		t.Fatal(failed)
	}
	if got := conversation(t, out); got.ID != c.ID || got.Awaiting != f.story.ID {
		t.Errorf("message_get should read a conversation without a story of its own: %+v", got)
	}
}

// S-0331: the sender's story must be in progress or in review, as the
// package holds it, and so must the addressee; a reply comes only from one of
// the conversation's two stories.
func TestAMessageGoesOnlyBetweenOpenStories(t *testing.T) {
	f, other := messagesFixture(t)
	epic, err := f.repo.Get(f.story.Parent)
	if err != nil {
		t.Fatal(err)
	}
	ready := storyOfEpic(t, f.repo, epic, "Waiting", nil, workitem.Ready)
	t.Setenv("FLAI_STORY", f.story.ID)
	if _, failed := f.call(t, "message_send", map[string]any{"to": ready.ID, "text": "Hello"}); !strings.Contains(failed, ready.ID+" is ready; a message goes only between stories in progress or in review") {
		t.Errorf("a send to a story not in progress should be refused with its state: %q", failed)
	}
	if _, failed := f.call(t, "message_send", map[string]any{"to": f.story.ID, "text": "Hello"}); !strings.Contains(failed, "cannot message itself") {
		t.Errorf("a send to the session's own story should be refused: %q", failed)
	}
	t.Setenv("FLAI_STORY", ready.ID)
	if _, failed := f.call(t, "message_send", map[string]any{"to": other.ID, "text": "Hello"}); !strings.Contains(failed, ready.ID+" is ready") {
		t.Errorf("a send from a story not in progress should be refused with its state: %q", failed)
	}
	if files := conversationFiles(t, f.repo); len(files) != 0 {
		t.Fatalf("a refused send should write nothing: %v", files)
	}

	t.Setenv("FLAI_STORY", other.ID)
	if _, failed := f.call(t, "message_send", map[string]any{"to": f.story.ID, "text": "Hello"}); failed != "" {
		t.Fatal(failed)
	}
	t.Setenv("FLAI_STORY", ready.ID)
	if _, failed := f.call(t, "message_reply", map[string]any{"id": "MS-0001", "text": "Me too"}); !strings.Contains(failed, ready.ID+" is not in MS-0001") {
		t.Errorf("a reply from a story not in the conversation should be refused: %q", failed)
	}
}

// S-0331: the message tools take no story to write as, their descriptions
// say messages are apart from the operator's threads (ADR-0120), and the
// instructions send an agent whose paths another story comes to claim to
// that story's agent, not to a thread, in a project and in a folder.
func TestTheMessageToolsAndInstructions(t *testing.T) {
	f := setup(t)
	res, err := f.cs.ListTools(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}
	seen := map[string]bool{}
	for _, tool := range res.Tools {
		if !strings.HasPrefix(tool.Name, "message_") {
			continue
		}
		seen[tool.Name] = true
		for _, want := range []string{"ADR-0120", "apart from the operator's threads", "flai message"} {
			if !strings.Contains(tool.Description, want) {
				t.Errorf("%s's description does not say %q: %s", tool.Name, want, tool.Description)
			}
		}
		schema, _ := json.Marshal(tool.InputSchema)
		var props struct {
			Properties map[string]any `json:"properties"`
		}
		_ = json.Unmarshal(schema, &props)
		for _, never := range []string{"from", "story"} {
			if _, ok := props.Properties[never]; ok {
				t.Errorf("%s takes %s, so an agent could write as another story: %s", tool.Name, never, schema)
			}
		}
	}
	if len(seen) != 5 {
		t.Errorf("message_send, message_reply, message_escalate, message_get, and message_share should be listed: %v", seen)
	}
	root := t.TempDir()
	makeProject(t, filepath.Join(root, "alpha"), "alpha")
	for name, in := range map[string]string{"project": f.cs.InitializeResult().Instructions, "folder": folderSetup(t, root).cs.InitializeResult().Instructions} {
		for _, want := range []string{"the story whose claim grew has messaged the other in the two stories' conversation", "agree there with message_reply which of you changes them first", "as a message from flai in that conversation", "Answer a message to your story with message_reply before you go on", "message_get", "ask the operator with message_escalate only when the two of you do not agree", "inbox lists your story's open conversations under messages, which awaiting_you does not count", "wait_for_events wakes on a message to your story, an event of kind message"} {
			if !strings.Contains(in, want) {
				t.Errorf("the %s instructions do not say %q: %s", name, want, in)
			}
		}
		if strings.Contains(in, "on a thread before you change them") {
			t.Errorf("the %s instructions still send an overlap to a thread: %s", name, in)
		}
	}
	for name, cs := range map[string]*mcp.ClientSession{"project": f.cs, "folder": folderSetup(t, root).cs} {
		res, err := cs.ListTools(context.Background(), nil)
		if err != nil {
			t.Fatal(err)
		}
		for _, tool := range res.Tools {
			want := map[string]string{"inbox": "messages lists the open conversations of this session's own story", "wait_for_events": "wakes it as an event of kind message"}[tool.Name]
			if want != "" && !strings.Contains(tool.Description, want) {
				t.Errorf("the %s %s description does not say %q: %s", name, tool.Name, want, tool.Description)
			}
		}
	}
}

// S-0331: the tools route by project in a folder of projects, as the other
// per-project tools do.
func TestTheMessageToolsRouteByProject(t *testing.T) {
	t.Setenv("FLAI_STORY", "")
	root := t.TempDir()
	makeProject(t, filepath.Join(root, "alpha"), "alpha")
	makeProject(t, filepath.Join(root, "beta"), "beta")
	f := folderSetup(t, root)
	res, err := f.cs.CallTool(context.Background(), &mcp.CallToolParams{Name: "message_get", Arguments: map[string]any{"id": "MS-1"}})
	if err != nil {
		t.Fatal(err)
	}
	if !res.IsError || !strings.Contains(res.Content[0].(*mcp.TextContent).Text, "name one with project") {
		t.Errorf("message_get with no project should ask for one: %+v", res.Content)
	}
	res, err = f.cs.CallTool(context.Background(), &mcp.CallToolParams{Name: "message_get", Arguments: map[string]any{"project": "beta", "id": "MS-1"}})
	if err != nil {
		t.Fatal(err)
	}
	if !res.IsError || !strings.Contains(res.Content[0].(*mcp.TextContent).Text, filepath.Join(root, "beta")) {
		t.Errorf("message_get should look in the project named: %+v", res.Content)
	}
}

// S-0331: inbox lists the open conversations of the session's own story
// under messages, which awaiting_you does not count, and a session with no
// story of its own has none.
func TestTheInboxToolListsTheStorysConversations(t *testing.T) {
	f, other := messagesFixture(t)
	if _, err := messages.Send(f.repo, messages.SendOptions{From: other.ID, To: f.story.ID, Author: "agent-" + other.ID, Text: "Will you leave plan.md to me?", About: []string{"design/system/plan.md"}, Now: t0}); err != nil {
		t.Fatal(err)
	}

	none, failed := f.call(t, "inbox", map[string]any{})
	if failed != "" {
		t.Fatal(failed)
	}
	if _, ok := none["messages"]; ok {
		t.Errorf("a session with no story of its own has no messages: %v", none["messages"])
	}

	t.Setenv("FLAI_STORY", f.story.ID)
	out, failed := f.call(t, "inbox", map[string]any{})
	if failed != "" {
		t.Fatal(failed)
	}
	list, _ := out["messages"].([]any)
	if len(list) != 1 {
		t.Fatalf("one conversation: %v", out["messages"])
	}
	m := list[0].(map[string]any)
	last, _ := m["last"].(map[string]any)
	if m["id"] != "MS-0001" || m["with"] != other.ID || m["awaiting"] != "you" || m["about"].([]any)[0] != "design/system/plan.md" || last["story"] != other.ID || last["text"] != "Will you leave plan.md to me?" {
		t.Errorf("the conversation as the story's agent sees it: %v", m)
	}
	if out["awaiting_you"].(float64) != 0 {
		t.Errorf("awaiting_you counts threads only: %v", out["awaiting_you"])
	}
}

// messageEvents are the events of kind message in a wait's answer.
func messageEvents(out map[string]any) []map[string]any {
	var got []map[string]any
	events, _ := out["events"].([]any)
	for _, e := range events {
		if e := e.(map[string]any); e["kind"] == "message" {
			got = append(got, e)
		}
	}
	return got
}

// S-0331: a held wait_for_events answers within a poll of a message to the
// session's story, as an event of kind message naming the conversation and
// the sender's story, told once; the story's own reply is not told back.
func TestAMessageWakesAWaitingAgent(t *testing.T) {
	f, other := messagesFixture(t)
	t.Setenv("FLAI_STORY", f.story.ID)
	if _, failed := f.call(t, "inbox", map[string]any{}); failed != "" {
		t.Fatal(failed)
	}

	type waited struct {
		out     map[string]any
		elapsed time.Duration
	}
	done := make(chan waited, 1)
	start := time.Now()
	go func() {
		out, _ := f.call(t, "wait_for_events", map[string]any{"timeout_seconds": 3})
		done <- waited{out, time.Since(start)}
	}()
	time.Sleep(150 * time.Millisecond)
	c, err := messages.Send(f.repo, messages.SendOptions{From: other.ID, To: f.story.ID, Author: "agent-" + other.ID, Text: "Will you leave plan.md to me?", About: []string{"design/system/plan.md"}, Now: *f.clock})
	if err != nil {
		t.Fatal(err)
	}
	select {
	case w := <-done:
		got := messageEvents(w.out)
		if w.out["timed_out"] == true || len(got) != 1 || got[0]["id"] != c.ID || got[0]["title"] != c.Title || got[0]["cause"] != other.ID || got[0]["by"] != "agent-"+other.ID || got[0]["to"] != "design/system/plan.md" {
			t.Fatalf("one message event naming the conversation and the sender's story: %v", w.out)
		}
		if w.elapsed > 1500*time.Millisecond {
			t.Errorf("the agent should hear within a poll, took %v", w.elapsed)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("wait_for_events never returned")
	}

	*f.clock = f.clock.Add(time.Minute)
	if _, failed := f.call(t, "message_reply", map[string]any{"id": c.ID, "text": "Yes."}); failed != "" {
		t.Fatal(failed)
	}
	quiet, _ := f.call(t, "wait_for_events", map[string]any{"timeout_seconds": 1})
	if got := messageEvents(quiet); quiet["timed_out"] != true || len(got) != 0 {
		t.Errorf("neither the message told nor the story's own reply is told again: %v", quiet)
	}
}

// S-0332, ADR-0121: message_escalate, from either story of a conversation,
// opens a thread on the session's story naming both stories, giving the
// conversation's path, and quoting the reason, and records it in the
// conversation, which then awaits the other story; a third story, an empty
// reason, and a closed conversation are refused, and nothing is written.
func TestMessageEscalateAsksTheOperator(t *testing.T) {
	f, other := messagesFixture(t)
	c, err := messages.Send(f.repo, messages.SendOptions{From: f.story.ID, To: other.ID, Author: "claude", Text: "May I take plan.md?", About: []string{"design/system/plan.md"}, Now: t0})
	if err != nil {
		t.Fatal(err)
	}

	t.Setenv("FLAI_STORY", other.ID)
	out, failed := f.call(t, "message_escalate", map[string]any{"id": "ms-1", "reason": "We both need plan.md this week."})
	if failed != "" {
		t.Fatal(failed)
	}
	data, _ := json.Marshal(out)
	var got EscalateOut
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatal(err)
	}
	th, conv := got.Thread, got.Conversation
	if th.ID != "TH-0001" || th.Title != other.ID+" and "+f.story.ID+" do not agree: May I take plan.md?" || len(th.EntryList) != 1 || th.EntryList[0].Author != "claude" {
		t.Errorf("message_escalate should open a thread titled for both stories: %+v", th)
	}
	if len(th.EntryList) == 1 {
		for _, want := range []string{other.ID + " and " + f.story.ID + " do not agree in their conversation MS-0001, `" + conv.Path + "`", "> We both need plan.md this week."} {
			if !strings.Contains(th.EntryList[0].Text, want) {
				t.Errorf("the thread's entry does not say %q:\n%s", want, th.EntryList[0].Text)
			}
		}
	}
	if conv.ID != c.ID || conv.Status != messages.StatusOpen || conv.Awaiting != f.story.ID || len(conv.Entries) != 2 || conv.Entries[1].Story != other.ID || !strings.Contains(conv.Entries[1].Text, "Asked the operator on TH-0001") {
		t.Errorf("the conversation should record the thread and await %s: %+v", f.story.ID, conv)
	}

	t.Setenv("FLAI_STORY", f.story.ID)
	if out, failed = f.call(t, "message_escalate", map[string]any{"id": c.ID, "reason": "And the docs."}); failed != "" {
		t.Fatal(failed)
	}
	data, _ = json.Marshal(out)
	got = EscalateOut{}
	_ = json.Unmarshal(data, &got)
	if got.Thread.ID != "TH-0002" || !strings.HasPrefix(got.Thread.Title, f.story.ID+" and "+other.ID+" do not agree") || got.Conversation.Awaiting != other.ID {
		t.Errorf("the sender's story escalates too, and the conversation then awaits %s: %+v", other.ID, got)
	}

	epic, err := f.repo.Get(f.story.Parent)
	if err != nil {
		t.Fatal(err)
	}
	third := storyOfEpic(t, f.repo, epic, "Third", nil, workitem.Ready, workitem.InProgress)
	if _, err := messages.Close(f.repo, c.ID, "operator", "settled", t0); err != nil {
		t.Fatal(err)
	}
	before, err := os.ReadFile(c.Path)
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct{ story, reason, want string }{
		{third.ID, "No.", third.ID + " is not in MS-0001"},
		{f.story.ID, " ", "an escalation needs a reason"},
		{f.story.ID, "No.", "MS-0001 is closed (settled) and cannot be escalated"},
		{"", "No.", "message_escalate writes as this session's own story, and it has none"},
	} {
		t.Setenv("FLAI_STORY", tc.story)
		if _, failed := f.call(t, "message_escalate", map[string]any{"id": c.ID, "reason": tc.reason}); !strings.Contains(failed, tc.want) {
			t.Errorf("an escalation from %q with %q should be refused with %q: %q", tc.story, tc.reason, tc.want, failed)
		}
	}
	if after, _ := os.ReadFile(c.Path); string(after) != string(before) {
		t.Errorf("a refused escalation should write nothing:\n%s", after)
	}
	if matches, _ := filepath.Glob(filepath.Join(f.repo.WipDir(), "threads", "*.md")); len(matches) != 2 {
		t.Errorf("a refused escalation should open no thread: %v", matches)
	}
}

// S-0334, ADR-0134: message_share, from the holding story's session, records
// the share flai's ask invited and answers the conversation with it; from the
// held story's session it is refused, and nothing is written.
func TestMessageShare(t *testing.T) {
	f, _ := messagesFixture(t)
	epic, err := f.repo.Get(f.story.Parent)
	if err != nil {
		t.Fatal(err)
	}
	held := storyOfEpic(t, f.repo, epic, "Held", []string{"flai/internal/mcpserver/messages.go"}, workitem.Ready)
	asked, _, err := messages.AskHold(f.repo, messages.AskOptions{Held: held.ID, Holder: f.story.ID, Paths: []string{"flai/internal/mcpserver/messages.go"}, Now: t0})
	if err != nil {
		t.Fatal(err)
	}
	clock := t0.Add(time.Hour)
	as := func(agent string) *server {
		return newServer(Options{Repo: f.repo, Agent: agent, Version: "test", Now: func() time.Time { return clock }}, f.repo)
	}
	in := MessageShareIn{ID: "ms-1", Paths: []string{"flai/internal/mcpserver/messages.go"}, Split: "The holder changes the descriptions; the held story adds a tool below them."}

	was := readFile(t, asked.Path)
	if _, _, err := as("agent-"+held.ID).messageShare(context.Background(), nil, in); err == nil || !strings.Contains(err.Error(), "agent-"+held.ID+", writing for "+held.ID+", may not share "+f.story.ID+"'s paths with "+held.ID) {
		t.Errorf("a share from the held story's session should be refused: %v", err)
	}
	if readFile(t, asked.Path) != was {
		t.Error("a refused share should write nothing")
	}

	_, out, err := as("agent-"+f.story.ID).messageShare(context.Background(), nil, in)
	if err != nil {
		t.Fatal(err)
	}
	if out.ID != "MS-0001" || len(out.Shares) != 1 || out.Shares[0].Holder != f.story.ID || out.Shares[0].Held != held.ID || out.Shares[0].By != "agent-"+f.story.ID || out.Shares[0].Split != in.Split {
		t.Errorf("message_share should answer the conversation with the share: %+v", out)
	}
	if last := out.Entries[len(out.Entries)-1]; last.Story != f.story.ID || !strings.Contains(last.Text, "shares `flai/internal/mcpserver/messages.go` with "+held.ID) {
		t.Errorf("the share's entry: %+v", last)
	}
	c, err := messages.Get(f.repo, "MS-0001")
	if err != nil {
		t.Fatal(err)
	}
	want, got := sortedJSON(t, messages.View(f.repo, c)), sortedJSON(t, out)
	if want != got {
		t.Errorf("message_share should answer what flai message show --json prints:\nwant %s\ngot  %s", want, got)
	}
}

// sortedJSON is v as JSON with the keys of every object sorted, so that a
// struct and a map holding the same fields compare equal.
func sortedJSON(t *testing.T, v any) string {
	t.Helper()
	raw, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	var fields any
	if err := json.Unmarshal(raw, &fields); err != nil {
		t.Fatal(err)
	}
	out, err := json.Marshal(fields)
	if err != nil {
		t.Fatal(err)
	}
	return string(out)
}

// readFile is a file's content.
func readFile(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}
