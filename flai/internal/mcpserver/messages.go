package mcpserver

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/bytepunx/system-flow/flai/internal/issues"
	"github.com/bytepunx/system-flow/flai/internal/messages"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// overlapInstructions tells a story's agent, in the server's instructions,
// what to do when a story in progress comes to claim the paths its story
// does: message that story's agent, not open a thread, and ask the operator
// only when the two do not agree (S-0331).
const overlapInstructions = "When its cause is a story in progress, not an accepted one, a write grew the two stories' claims to overlap on those paths: before you change them, message that story's agent with message_send, naming the paths in about, and narrow your touches if you can (I-0059). Answer a message to your story with message_reply before you go on, and read a conversation with message_get. Messages go between the agents of two open stories, apart from the operator's threads (ADR-0120): open a thread for the operator with thread_open only when the two of you do not agree."

// messagesWhat says what a message is, for each message tool's description
// (ADR-0120).
const messagesWhat = " A message goes between the agents of two open stories, each in progress or in review, and the messages between two stories make a conversation, one file in wip/messages, MS-nnnn-<slug>.md, in the main checkout (ADR-0120). Messages are kept apart from the operator's threads: none awaits the operator or appears among the threads; a question for the designer is still a thread, opened with thread_open."

// messagesFrom says whose story a message comes from: the calling session's,
// never one it names (S-0331).
const messagesFrom = " The story it comes from is this session's own: FLAI_STORY, else the story in this agent's name of the form agent-S-nnnn, else the story branch checked out where this server runs; there is no argument to name another, so an agent writes only as its own story. A session with none is refused, and nothing is written. The author of the entry is this agent's name."

// messagesReturns says what the message tools answer.
const messagesReturns = " Returns the conversation as flai message show --json prints it: id, title, from, to, about, status, closed and closed_reason, awaiting (the story that did not write the last message, while it is open), participants, created, updated, path (from the project root), and entries, each with its time, author, story, and text."

const messageSendDescription = "Start a conversation from this session's story to the story to with its first message, as flai message send does (S-0331): text is the message, whose first line is the conversation's title, and about the repository paths it is about, each a file or folder relative to the root that must exist, here or in this story's worktree. Message the other story's agent this way before you change paths both stories claim, naming them in about." + messagesWhat + messagesFrom + " A send to or from a story not in progress or in review, or archived, is refused with its state, as is a send to this story itself, and nothing is written." + messagesReturns

const messageReplyDescription = "Add a message to the conversation id from this session's story, one of its two, as flai message reply does (S-0331); the conversation then awaits the other story. Answer a message to your story before you go on with the paths it is about." + messagesWhat + messagesFrom + " A reply is refused from a story that is not one of the two, from a story no longer in progress or in review, and on a conversation that reads as closed (its status is closed, or either story is done, cancelled, or archived), which takes none: send a new message to start another; nothing is written." + messagesReturns

const messageGetDescription = "Read the conversation id (any zero padding, such as ms-4), as flai message show --json prints it (S-0331): its front matter, whether it reads as closed and why, the story it awaits while open, and every entry. Any session may read one; it needs no story of its own." + messagesWhat + messagesReturns

// MessageSendIn starts a conversation from the calling session's story. It
// has no from: an agent sends only as its own story (S-0331).
type MessageSendIn struct {
	Project string   `json:"project,omitempty" jsonschema:"the project, by key or folder: needed only when the server serves more than one"`
	To      string   `json:"to" jsonschema:"the story the message goes to, such as S-0331 (any zero padding): in progress or in review, and not this session's own"`
	Text    string   `json:"text" jsonschema:"the message; its first line is the conversation's title"`
	About   []string `json:"about,omitempty" jsonschema:"repository paths the conversation is about, each a file or folder relative to the root that must exist"`
}

func (in MessageSendIn) project() string { return in.Project }

// MessageReplyIn answers a conversation from the calling session's story.
type MessageReplyIn struct {
	Project string `json:"project,omitempty" jsonschema:"the project, by key or folder: needed only when the server serves more than one"`
	ID      string `json:"id" jsonschema:"the conversation, such as MS-0004 (any zero padding)"`
	Text    string `json:"text" jsonschema:"the reply"`
}

func (in MessageReplyIn) project() string { return in.Project }

// MessageGetIn names a conversation.
type MessageGetIn struct {
	Project string `json:"project,omitempty" jsonschema:"the project, by key or folder: needed only when the server serves more than one"`
	ID      string `json:"id" jsonschema:"the conversation, such as MS-0004 (any zero padding)"`
}

func (in MessageGetIn) project() string { return in.Project }

// ConversationOut is a conversation field for field as messages.View gives
// it, which flai message show --json prints.
type ConversationOut struct {
	ID           string           `json:"id"`
	Title        string           `json:"title"`
	From         string           `json:"from" jsonschema:"the story that started the conversation"`
	To           string           `json:"to" jsonschema:"the story it was started with"`
	About        []string         `json:"about" jsonschema:"the repository paths it is about"`
	Status       string           `json:"status" jsonschema:"open or closed, as stored"`
	Closed       bool             `json:"closed" jsonschema:"whether it reads as closed: its status is closed, or either story is done, cancelled, or archived; a closed conversation takes no reply"`
	ClosedReason string           `json:"closed_reason" jsonschema:"with closed, why"`
	Awaiting     string           `json:"awaiting" jsonschema:"while open, the story that did not write the last message, whose agent answers next"`
	Participants []string         `json:"participants" jsonschema:"the agents that wrote in it"`
	Created      string           `json:"created"`
	Updated      string           `json:"updated"`
	Path         string           `json:"path" jsonschema:"from the project root"`
	Entries      []messages.Entry `json:"entries" jsonschema:"every message in order, each with its time, author, story (empty on the entry that closed it), and text"`
}

// conversationOut is c as messages.View gives it, so that the tools answer
// what flai message show --json prints and cannot drift from it.
func conversationOut(r *workitem.Repo, c *messages.Conversation) (ConversationOut, error) {
	data, err := json.Marshal(messages.View(r, c))
	if err != nil {
		return ConversationOut{}, err
	}
	var out ConversationOut
	if err := json.Unmarshal(data, &out); err != nil {
		return ConversationOut{}, err
	}
	return out, nil
}

// messageStory is the calling session's own story, as recordingStory resolves
// it with nothing given; none refuses tool, which then writes nothing.
func (s *server) messageStory(tool string) (string, error) {
	if story := s.recordingStory(""); story != "" {
		return story, nil
	}
	return "", fmt.Errorf("%s writes as this session's own story, and it has none: FLAI_STORY is not set, the agent name %q is not of the form agent-S-nnnn, and no story branch is checked out where this server runs. Messages go between the agents of two open stories (ADR-0120): send from the session of a story in progress or in review, or ask the operator on a thread with thread_open; nothing was written", tool, s.agent)
}

func (s *server) messageSend(_ context.Context, _ *mcp.CallToolRequest, in MessageSendIn) (*mcp.CallToolResult, ConversationOut, error) {
	from, err := s.messageStory("message_send")
	if err != nil {
		return nil, ConversationOut{}, err
	}
	// about may name a path the story has made in its worktree and not yet
	// merged: look there first, as flai message send run in the worktree
	// does, then in the main checkout, where wip/messages always is
	r := issues.RepoFor(s.repo, from)
	if r != s.repo {
		r.MainRoot = projectRoot(s.repo)
	}
	c, err := messages.Send(r, messages.SendOptions{From: from, To: in.To, Author: s.agent, Text: in.Text, About: in.About, Now: s.now()})
	if err != nil {
		return nil, ConversationOut{}, err
	}
	out, err := conversationOut(s.repo, c)
	return nil, out, err
}

func (s *server) messageReply(_ context.Context, _ *mcp.CallToolRequest, in MessageReplyIn) (*mcp.CallToolResult, ConversationOut, error) {
	story, err := s.messageStory("message_reply")
	if err != nil {
		return nil, ConversationOut{}, err
	}
	c, err := messages.Reply(s.repo, in.ID, story, s.agent, in.Text, s.now())
	if err != nil {
		return nil, ConversationOut{}, err
	}
	out, err := conversationOut(s.repo, c)
	return nil, out, err
}

func (s *server) messageGet(_ context.Context, _ *mcp.CallToolRequest, in MessageGetIn) (*mcp.CallToolResult, ConversationOut, error) {
	c, err := messages.Get(s.repo, in.ID)
	if err != nil {
		return nil, ConversationOut{}, err
	}
	out, err := conversationOut(s.repo, c)
	return nil, out, err
}
