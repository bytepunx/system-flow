package guard

import (
	"encoding/json"
	"fmt"
	"io"
	"slices"
	"strings"
)

// OpenCode is the harness name of OpenCode, whose plugin flai-guard.js, in
// .opencode/plugins, runs flai guard --harness opencode on each call (S-0407).
const OpenCode = "opencode"

// OpenCodeRefuses is the exit status with which flai guard refuses an
// OpenCode call: the plugin throws with the guard's standard error, which
// OpenCode hands back to the model as the tool's error.
const OpenCodeRefuses = 2

// The events the OpenCode plugin hands the guard: a tool call before it runs,
// a session created, and a tool call after it ran. An input that names no
// event is a tool call.
const (
	OpenCodeToolBefore     = "tool.execute.before"
	OpenCodeSessionCreated = "session.created"
	OpenCodeToolAfter      = "tool.execute.after"
)

// OpenCodeMCPPrefix starts the name of each of flai's MCP tools as OpenCode
// knows them, from a server named flai: <server>_<tool>, so flai_inbox.
const OpenCodeMCPPrefix = "flai_"

// openCodeTask is OpenCode's Task tool, which runs a sub-agent in a child
// session and returns when it ends.
const openCodeTask = "task"

// openCodeFileWrites are OpenCode's tools that write the file their filePath
// names; apply_patch writes each file its patchText names.
var openCodeFileWrites = []string{"edit", "write", "multiedit"}

// openCodePatches are OpenCode's tools that write the files a patch names.
var openCodePatches = []string{"apply_patch", "patch"}

// OpenCodeInput is what the plugin hands flai guard on standard input, one
// JSON object per call: the Event, the Tool and its Args, the SessionID the
// call is made in, its CallID, and Session, the session the event is about,
// as OpenCode's client reads it: the caller's for a tool call, the session
// made for session.created, and the sub-agent's child session for a task's
// tool.execute.after. A session whose ParentID is set is a sub-agent's: the
// Task tool runs a sub-agent in a child session of its caller's
// (agent-adapters.md § Read from the source at v1.19.5).
type OpenCodeInput struct {
	Event     string          `json:"event"`
	Tool      string          `json:"tool"`
	SessionID string          `json:"sessionID"`
	CallID    string          `json:"callID"`
	Args      json.RawMessage `json:"args"`
	Session   OpenCodeSession `json:"session"`
}

// OpenCodeSession is the part of an OpenCode session the guard reads: its
// ID, its parent's, and the name of the agent it runs.
type OpenCodeSession struct {
	ID       string `json:"id"`
	ParentID string `json:"parentID"`
	Agent    string `json:"agent"`
}

// openCodeArgs is the part of an OpenCode tool's arguments the guard reads
// beside a flai tool's Args: bash's command, the file edit and write change,
// the content write gives its file, and apply_patch's patch.
type openCodeArgs struct {
	Command   string `json:"command"`
	FilePath  string `json:"filePath"`
	Content   string `json:"content"`
	PatchText string `json:"patchText"`
}

// openCodeWhole is OpenCode's tool that replaces its file whole with the
// content it gives.
const openCodeWhole = "write"

// openCode is OpenCode's Reader.
type openCode struct{}

// Read reads the OpenCode plugin's input: a session.created whose session has
// a parent is a sub-agent starting in the parent's session, a task's
// tool.execute.after the sub-agent of its child session stopping, and a
// tool.execute.before, or an input naming no event, a tool call in g's role
// and story. Any other input is a call the guard holds no rule on.
func (openCode) Read(data []byte, g Guard) (Hook, error) {
	var in OpenCodeInput
	if err := json.Unmarshal(data, &in); err != nil {
		return Hook{}, err
	}
	switch in.Event {
	case "", OpenCodeToolBefore:
		return g.openCodeCall(in)
	case OpenCodeSessionCreated:
		if in.Session.ParentID != "" {
			return Hook{SubAgent: &SubAgentChange{Session: in.Session.ParentID, ID: in.Session.ID, Type: in.Session.Agent}}, nil
		}
	case OpenCodeToolAfter:
		if in.Tool == openCodeTask && in.Session.ID != "" {
			return Hook{SubAgent: &SubAgentChange{Session: in.Session.ParentID, ID: in.Session.ID, Type: in.Session.Agent, Stopped: true}}, nil
		}
	}
	return Hook{Call: Call{Session: in.SessionID, Role: g.Role, Story: g.Story, Name: in.Tool}}, nil
}

// openCodeCall is the Call in makes, in g's role and story: bash runs a
// command line, edit, write, and multiedit write their filePath, apply_patch
// writes each file its patch names, the first as Call and the rest in Also,
// and a tool of OpenCodeMCPPrefix is a flai tool; a call from a session with
// a parent is a sub-agent's.
func (g Guard) openCodeCall(in OpenCodeInput) (Hook, error) {
	c := Call{Session: in.SessionID, Role: g.Role, Story: g.Story, Name: in.Tool}
	if in.Session.ParentID != "" {
		c.Subagent, c.Agent = true, in.Session.Agent
	}
	args := in.Args
	if len(args) == 0 || string(args) == "null" {
		args = json.RawMessage("{}")
	}
	tool, flai := strings.CutPrefix(in.Tool, OpenCodeMCPPrefix)
	if flai {
		var input Input
		if err := json.Unmarshal(args, &input); err != nil {
			return Hook{}, fmt.Errorf("reading the arguments of %s: %w", in.Tool, err)
		}
		c.Kind, c.Tool, c.Args = FlaiTool, tool, input.Args
		return Hook{Call: c}, nil
	}
	var a openCodeArgs
	if err := json.Unmarshal(args, &a); err != nil {
		return Hook{}, fmt.Errorf("reading the arguments of %s: %w", in.Tool, err)
	}
	switch {
	case in.Tool == "bash":
		c.Kind, c.Command = Shell, a.Command
	case slices.Contains(openCodeFileWrites, in.Tool):
		c.Kind, c.Path = FileWrite, a.FilePath
		if in.Tool == openCodeWhole {
			c.Whole, c.Content = true, a.Content
		}
	case slices.Contains(openCodePatches, in.Tool):
		c.Kind = FileWrite
		paths := PatchPaths(a.PatchText)
		if len(paths) == 0 {
			return Hook{Call: c}, nil
		}
		c.Path = paths[0]
		h := Hook{Call: c}
		for _, p := range paths[1:] {
			more := c
			more.Path = p
			h.Also = append(h.Also, more)
		}
		return h, nil
	}
	return Hook{Call: c}, nil
}

// patchMarkers start the lines of an OpenCode patch that name a file it
// adds, changes, deletes, or moves a file to.
var patchMarkers = []string{"*** Add File: ", "*** Update File: ", "*** Delete File: ", "*** Move to: "}

// PatchPaths are the files an OpenCode apply_patch patch names, in the order
// it names them, each once.
func PatchPaths(patch string) []string {
	var paths []string
	for line := range strings.Lines(patch) {
		line = strings.TrimRight(line, "\r\n")
		for _, m := range patchMarkers {
			if p, ok := strings.CutPrefix(line, m); ok {
				if p = strings.TrimSpace(p); p != "" && !slices.Contains(paths, p) {
					paths = append(paths, p)
				}
			}
		}
	}
	return paths
}

// Answer refuses a denied call with OpenCodeRefuses and its reason on
// stderr, which the plugin throws with, and passes any other with exit 0 and
// nothing on stderr: an Ask goes on to OpenCode's own permission rules, the
// deny rules flai gives it (ADR-0260), since the guard does not yet hold a
// call to ask on OpenCode.
func (openCode) Answer(v Verdict, stderr io.Writer) int {
	if v.Outcome != Deny {
		return 0
	}
	fmt.Fprintln(stderr, v.Reason)
	return OpenCodeRefuses
}
