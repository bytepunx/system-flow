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

	"github.com/bytepunx/system-flow/flai/internal/ask"
	"github.com/bytepunx/system-flow/flai/internal/protected"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// Claude Code's permission prompt tool (S-0257, ADR-0106). Run headless,
// Claude Code refuses a write to a path it protects, such as a .claude/
// folder or .mcp.json, unless a person or the tool named by
// --permission-prompt-tool approves it. flai serve names this one: it
// approves such a write in an in-progress story's worktree when the operator
// answers allow on a thread, or at once when the host's AutoApprove says so,
// and refuses everything else, a path in .git always.

// AutoApprove reports whether the project at root lets permission_prompt
// allow a write to a protected path without asking the operator.
type AutoApprove func(root string) bool

const permissionPromptDescription = "Claude Code calls this tool itself, through --permission-prompt-tool mcp__flai__permission_prompt, when a tool call would otherwise ask a person: the model need not call it. It approves only an Edit, Write, MultiEdit, or NotebookEdit of a path Claude Code protects inside an in-progress story's worktree, such as a file in a .claude/ folder or .mcp.json, which Claude Code refuses without a person's approval; it never approves a path in .git. Unless the host auto-approves such writes, it opens a thread on the story showing the change and waits at most four minutes for the operator's reply: allow lets the write through, anything else refuses it with the operator's words as the reason. Unanswered by then, it refuses the write naming the thread, which stays open; the same write made again takes the answer given on it since, or waits again. Everything else is refused at once. It answers {\"behavior\":\"allow\",\"updatedInput\":...} or {\"behavior\":\"deny\",\"message\":...} as Claude Code expects."

// permissionScope is what every refusal of a request outside the tool's remit says.
const permissionScope = "permission_prompt approves only an Edit, Write, MultiEdit, or NotebookEdit of a path Claude Code protects, other than .git, inside an in-progress story's worktree; anything else the session's permissions do not allow is refused, as before"

// permissionPoll is how often a held permission_prompt reads its thread for
// the operator's answer; tests shorten it.
var permissionPoll = time.Second

// permissionWait is the longest a permission_prompt call waits for an answer
// (ADR-0124): below Claude Code's idle timeout for an MCP call, 30 minutes
// over stdio and 5 over HTTP (I-0103). Tests shorten it.
var permissionWait = 4 * time.Minute

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
	shown := strings.Join(parts[1:], "/")
	if protected.Git(shown) {
		return deny("%s is in .git, which permission_prompt never approves: git's link to the repository and its history are not an agent's to write; %s", path, permissionScope)
	}
	if !protected.Path(shown) {
		return deny("%s is not a path Claude Code protects in %s's worktree: %s", path, it.ID, permissionScope)
	}
	if escapes(s.repo.WorktreePath(it.ID), path) {
		return deny("%s leaves %s's worktree through a symbolic link: %s", path, it.ID, permissionScope)
	}
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

// askOperator maps Claude Code's request onto the hold-and-ask (ADR-0130),
// bounded by permissionWait, and its answer onto Claude Code's decision.
func (s *server) askOperator(ctx context.Context, story *workitem.Item, in PermissionIn, rel string) PermissionOut {
	asker := ask.Asker{Repo: s.repo, Now: s.now, Poll: permissionPoll, Closing: s.closing, Logger: s.logger}
	if ok, why := asker.Hold(ctx, permissionRequest(s.agent, story, in, rel), permissionWait); !ok {
		return deny("%s", why)
	}
	return allow(in.Input)
}

// claudeProtects is why a person must answer, as Claude Code has it.
const claudeProtects = "Claude Code protects the path and does not write it without a person's approval."

// permissionRequest is Claude Code's request for the agent to write rel in
// the story's worktree, with the change its input would make.
func permissionRequest(agent string, story *workitem.Item, in PermissionIn, rel string) ask.Request {
	return ask.Request{Story: story.ID, Owner: story.Owner, Agent: agent, Tool: in.ToolName, Path: rel, Why: claudeProtects, Change: changeText(in)}
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

// changeText shows what the tool call would write: the whole content, each
// edit's old and new text, or the notebook cell's source.
func changeText(in PermissionIn) string {
	var b strings.Builder
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
