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
	if len(seen) != 3 {
		t.Errorf("message_send, message_reply, and message_get should be listed: %v", seen)
	}
	root := t.TempDir()
	makeProject(t, filepath.Join(root, "alpha"), "alpha")
	for name, in := range map[string]string{"project": f.cs.InitializeResult().Instructions, "folder": folderSetup(t, root).cs.InitializeResult().Instructions} {
		for _, want := range []string{"message that story's agent with message_send, naming the paths in about", "Answer a message to your story with message_reply before you go on", "message_get", "open a thread for the operator with thread_open only when the two of you do not agree"} {
			if !strings.Contains(in, want) {
				t.Errorf("the %s instructions do not say %q: %s", name, want, in)
			}
		}
		if strings.Contains(in, "on a thread before you change them") {
			t.Errorf("the %s instructions still send an overlap to a thread: %s", name, in)
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
