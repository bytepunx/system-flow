package mcpserver

import (
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/bytepunx/system-flow/flai/internal/threads"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// Claude Code's permission prompt tool (S-0257). Run headless, Claude Code
// refuses a write under any .claude/ folder as sensitive unless a person or
// the tool named by --permission-prompt-tool approves it. flai serve names
// this one: it approves such a write in an in-progress story's worktree when
// the operator answers allow on a thread, or at once when the host's
// AutoApprove says so, and refuses everything else.

// AutoApprove reports whether the project at root lets permission_prompt
// allow a write under .claude/ without asking the operator.
type AutoApprove func(root string) bool

const permissionPromptDescription = "Claude Code calls this tool itself, through --permission-prompt-tool mcp__flai__permission_prompt, when a tool call would otherwise ask a person: the model need not call it. It approves only an Edit, Write, MultiEdit, or NotebookEdit of a file in a .claude/ folder inside an in-progress story's worktree, which Claude Code refuses as sensitive without a person's approval. Unless the host auto-approves such writes, it opens a thread on the story showing the change and waits until the operator replies: allow lets the write through, anything else refuses it with the operator's words as the reason. Everything else is refused at once. It answers {\"behavior\":\"allow\",\"updatedInput\":...} or {\"behavior\":\"deny\",\"message\":...} as Claude Code expects."

// permissionScope is what every refusal of a request outside the tool's remit says.
const permissionScope = "permission_prompt approves only an Edit, Write, MultiEdit, or NotebookEdit of a file in a .claude/ folder inside an in-progress story's worktree; anything else the session's permissions do not allow is refused, as before"

// permissionPoll is how often a held permission_prompt reads its thread for
// the operator's answer; tests shorten it.
var permissionPoll = time.Second

// permissionWords are the first words of an answer that allow the write.
var permissionWords = []string{"allow", "yes", "approve", "approved", "ok"}

// PermissionIn is Claude Code's permission request.
type PermissionIn struct {
	Project   string         `json:"project,omitempty"`
	ToolName  string         `json:"tool_name"`
	Input     map[string]any `json:"input,omitempty"`
	ToolUseID string         `json:"tool_use_id,omitempty"`
}

func (in PermissionIn) project() string { return in.Project }

// PermissionOut is the decision, sent to Claude Code as the JSON text of the
// result's single text block.
type PermissionOut struct {
	Behavior     string         `json:"behavior"`
	UpdatedInput map[string]any `json:"updatedInput,omitempty"`
	Message      string         `json:"message,omitempty"`
}

// permissionInputSchema is given rather than inferred: an inferred schema
// refuses properties it does not name, and Claude Code may send more.
var permissionInputSchema = map[string]any{
	"type": "object",
	"properties": map[string]any{
		"project":     map[string]any{"type": "string", "description": "the project, by key or folder: needed only when the server serves more than one and the path does not say which"},
		"tool_name":   map[string]any{"type": "string", "description": "the tool Claude Code asks permission for"},
		"input":       map[string]any{"type": "object", "description": "the tool call's input"},
		"tool_use_id": map[string]any{"type": "string"},
	},
	"required": []string{"tool_name"},
}

func allow(input map[string]any) PermissionOut {
	return PermissionOut{Behavior: "allow", UpdatedInput: input}
}

func deny(format string, args ...any) PermissionOut {
	return PermissionOut{Behavior: "deny", Message: fmt.Sprintf(format, args...)}
}

// permissionRoute serves permission_prompt for the project whose worktrees
// hold the path, when the request names none: Claude Code never does. It is a
// raw handler because Claude Code accepts only a single text block (I-0082):
// a typed one declares an output schema and answers structured content too.
func permissionRoute(p projects) mcp.ToolHandler {
	return func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		var in PermissionIn
		if err := json.Unmarshal(req.Params.Arguments, &in); err != nil {
			return permissionResult(deny("the arguments are not a permission request, an object with tool_name a string and input an object (%v): %s", err, permissionScope))
		}
		if in.Project == "" {
			if path, ok := permissionPath(in); ok {
				for _, s := range p.all() {
					if within(s.worktrees(), path) {
						return permissionResult(s.permissionPrompt(ctx, in))
					}
				}
			}
		}
		s, err := p.pick(in.project())
		if err != nil {
			return permissionResult(deny("%v", err))
		}
		return permissionResult(s.permissionPrompt(ctx, in))
	}
}

// permissionResult is the decision as Claude Code reads it: its JSON as the
// only content, a text block, with no structured content.
func permissionResult(out PermissionOut) (*mcp.CallToolResult, error) {
	text, err := json.Marshal(out)
	if err != nil {
		return nil, fmt.Errorf("encode the permission decision %+v: %w", out, err)
	}
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: string(text)}}}, nil
}

// permissionPath is the cleaned absolute path a request writes, if it is a
// write permission_prompt can approve.
func permissionPath(in PermissionIn) (string, bool) {
	key := "file_path"
	switch in.ToolName {
	case "Edit", "Write", "MultiEdit":
	case "NotebookEdit":
		key = "notebook_path"
	default:
		return "", false
	}
	p, _ := in.Input[key].(string)
	if p == "" || !filepath.IsAbs(p) {
		return "", false
	}
	return filepath.Clean(p), true
}

// within reports whether path lies strictly below dir.
func within(dir, path string) bool {
	rel, err := filepath.Rel(dir, path)
	return err == nil && rel != "." && rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator)) && !filepath.IsAbs(rel)
}

func (s *server) worktrees() string { return filepath.Join(s.repo.CacheDir(), "worktrees") }

// permissionPrompt decides a request for this project.
func (s *server) permissionPrompt(ctx context.Context, in PermissionIn) PermissionOut {
	path, ok := permissionPath(in)
	if !ok {
		why := fmt.Sprintf("%s is not an Edit, Write, MultiEdit, or NotebookEdit", orQuoted(in.ToolName))
		if slices.Contains([]string{"Edit", "Write", "MultiEdit", "NotebookEdit"}, in.ToolName) {
			why = fmt.Sprintf("the %s names no absolute path", in.ToolName)
		}
		return deny("%s: %s", why, permissionScope)
	}
	if !within(s.worktrees(), path) {
		return deny("%s is not inside a story's worktree under %s: %s", path, s.worktrees(), permissionScope)
	}
	rel, _ := filepath.Rel(s.worktrees(), path)
	parts := strings.Split(filepath.ToSlash(rel), "/")
	id := parts[0]
	notStory := deny("%s is not a story's worktree: %s", filepath.Join(s.worktrees(), id), permissionScope)
	it, err := s.repo.Get(id)
	if err != nil {
		return notStory // no such item is a refusal for Claude Code, not a tool error
	}
	if it.Type != workitem.Story || filepath.Clean(s.repo.WorktreePath(it.ID)) != filepath.Join(s.worktrees(), id) {
		return notStory
	}
	if it.Status != workitem.InProgress {
		return deny("%s is %s, not in progress: %s", it.ID, it.Status, permissionScope)
	}
	inTree := parts[1:]
	if len(inTree) < 2 || !slices.Contains(inTree[:len(inTree)-1], ".claude") {
		return deny("%s is not in a .claude/ folder of %s's worktree: %s", path, it.ID, permissionScope)
	}
	if escapes(s.repo.WorktreePath(it.ID), path) {
		return deny("%s leaves %s's worktree through a symbolic link: %s", path, it.ID, permissionScope)
	}
	shown := strings.Join(inTree, "/")
	if s.autoApprove != nil && s.autoApprove(projectRoot(s.repo)) {
		if s.logger != nil {
			s.logger.Info("permission allowed without asking", "component", "mcp", "agent", s.agent, "tool", in.ToolName, "path", path, "story", it.ID)
		}
		return allow(in.Input)
	}
	return s.askOperator(ctx, it, in, shown)
}

// escapes reports whether the deepest existing folder of path, symbolic
// links followed, lies outside the worktree.
func escapes(worktree, path string) bool {
	root, err := filepath.EvalSymlinks(worktree)
	if err != nil {
		return false // no worktree on disk: nothing to follow
	}
	dir := filepath.Dir(path)
	for {
		if real, err := filepath.EvalSymlinks(dir); err == nil {
			return real != root && !within(root, real)
		}
		up := filepath.Dir(dir)
		if up == dir {
			return false
		}
		dir = up
	}
}

func orQuoted(name string) string {
	if name == "" {
		return "a request without a tool name"
	}
	return name
}

// askOperator opens a thread on the story showing the change, waits for an
// answer from the story's owner, the operator (anyone but the agent when it
// has none), and resolves the thread with what was decided. Another agent's
// entry is not an answer.
func (s *server) askOperator(ctx context.Context, story *workitem.Item, in PermissionIn, rel string) PermissionOut {
	th, err := threads.New(s.repo, threads.NewOptions{
		Title:  fmt.Sprintf("Allow %s %s?", in.ToolName, rel),
		On:     story.ID,
		Author: s.agent,
		Text:   permissionRequest(s.agent, story.ID, story.Owner, in, rel),
		Now:    s.now(),
	})
	if err != nil {
		return deny("could not ask the operator on a thread: %v", err)
	}
	_ = s.mirror(th)
	// Only entries after the opening one can answer: the content shown may
	// itself hold lines that read as entries.
	asked := len(th.Entries())
	answer, ok := s.awaitAnswer(ctx, th.ID, asked, story.Owner)
	if !ok {
		s.settle(th.ID, fmt.Sprintf("refused: no answer before the session ended, so %s was not written", rel))
		return deny("no answer from the operator before the session ended")
	}
	if allowed(answer.Text) {
		s.settle(th.ID, fmt.Sprintf("allowed by %s: %s may %s %s", answer.Author, s.agent, in.ToolName, rel))
		return allow(in.Input)
	}
	s.settle(th.ID, fmt.Sprintf("refused by %s: %s was not written", answer.Author, rel))
	return deny("the operator refused: %s", answer.Text)
}

// awaitAnswer reads the thread until an entry after the first asked ones is
// by owner (by anyone but the agent when owner is empty), or the session
// ends.
func (s *server) awaitAnswer(ctx context.Context, id string, asked int, owner string) (threads.Entry, bool) {
	tick := time.NewTicker(permissionPoll)
	defer tick.Stop()
	for {
		select {
		case <-ctx.Done():
			return threads.Entry{}, false
		case <-s.closing: // nil, and so never ready, unless the server was given one
			return threads.Entry{}, false
		case <-tick.C:
			th, err := threads.Get(s.repo, id)
			if err != nil {
				continue // being written, or gone for a moment: read it again
			}
			entries := th.Entries()
			for i := asked; i < len(entries); i++ {
				if a := entries[i].Author; a != s.agent && (owner == "" || a == owner) {
					return entries[i], true
				}
			}
		}
	}
}

// settle resolves the thread as the agent, unless the operator resolved it.
func (s *server) settle(id, reason string) {
	th, err := threads.Get(s.repo, id)
	if err != nil || !th.Open() {
		return
	}
	if th, err = threads.Resolve(s.repo, id, s.agent, reason, s.now()); err == nil {
		_ = s.mirror(th)
	} else if s.logger != nil {
		s.logger.Warn("permission thread not resolved", "component", "mcp", "agent", s.agent, "thread", id, "err", err.Error())
	}
}

// allowed reports whether an answer's first word lets the write through.
func allowed(text string) bool {
	words := strings.Fields(strings.ToLower(text))
	if len(words) == 0 {
		return false
	}
	return slices.Contains(permissionWords, strings.TrimRight(words[0], ".,;:!?"))
}

var backticks = regexp.MustCompile("`+")

// fenced is text in a fenced code block longer than any backtick run in it,
// with the tab rule off around it when it holds a tab: what is shown is what
// would be written.
func fenced(text string) string {
	text = strings.TrimSuffix(text, "\n")
	n := 3
	for _, run := range backticks.FindAllString(text, -1) {
		if len(run) >= n {
			n = len(run) + 1
		}
	}
	fence := strings.Repeat("`", n)
	block := fence + "text\n" + text + "\n" + fence
	if text == "" {
		block = fence + "text\n" + fence
	}
	if strings.Contains(text, "\t") {
		block = "<!-- markdownlint-disable MD010 -->\n" + block + "\n<!-- markdownlint-enable MD010 -->"
	}
	return block
}

func str(m map[string]any, key string) string {
	v, _ := m[key].(string)
	return v
}

// permissionRequest is the thread's first entry: who asks, why a person must
// answer, how to, and the change itself.
func permissionRequest(agent, story, owner string, in PermissionIn, rel string) string {
	var b strings.Builder
	fmt.Fprintf(&b, "%s asks to %s `%s` in %s's worktree. Claude Code refuses writes under .claude/ without a person's approval.\n\n", agent, in.ToolName, rel, story)
	if owner != "" {
		fmt.Fprintf(&b, "Reply `allow`, as %s, the story's owner, to let it write. Anything else refuses it, and your words go back to the agent as the reason; a reply by anyone else is not an answer.\n\n", owner)
	} else {
		b.WriteString("Reply `allow` to let it write. Anything else refuses it, and your words go back to the agent as the reason.\n\n")
	}
	switch in.ToolName {
	case "Write":
		fmt.Fprintf(&b, "The whole content it would write:\n\n%s\n", fenced(str(in.Input, "content")))
	case "Edit":
		b.WriteString(editText(in.Input, "It would replace"))
	case "MultiEdit":
		edits, _ := in.Input["edits"].([]any)
		fmt.Fprintf(&b, "It would make %d edits, in order.\n", len(edits))
		for i, e := range edits {
			m, _ := e.(map[string]any)
			b.WriteString("\n")
			b.WriteString(editText(m, fmt.Sprintf("Edit %d would replace", i+1)))
		}
	case "NotebookEdit":
		mode := str(in.Input, "edit_mode")
		if mode == "" {
			mode = "replace"
		}
		cell := str(in.Input, "cell_id")
		if cell == "" {
			cell = "the first cell"
		} else {
			cell = "cell `" + cell + "`"
		}
		fmt.Fprintf(&b, "Edit mode %s, at %s, with this source:\n\n%s\n", mode, cell, fenced(str(in.Input, "new_source")))
	}
	return b.String()
}

// editText shows one edit's old and new text.
func editText(m map[string]any, lead string) string {
	var b strings.Builder
	fmt.Fprintf(&b, "%s this text:\n\n%s\n\nwith this text:\n\n%s\n", lead, fenced(str(m, "old_string")), fenced(str(m, "new_string")))
	if all, _ := m["replace_all"].(bool); all {
		b.WriteString("\nIt replaces every occurrence (replace_all).\n")
	}
	return b.String()
}
