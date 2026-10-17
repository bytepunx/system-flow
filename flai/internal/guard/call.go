package guard

import (
	"cmp"
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/ask"
	"github.com/bytepunx/system-flow/flai/internal/protected"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// Kind is what a call does, as far as the guard tells calls apart.
type Kind int

// The kinds of call: Other is one the guard holds no rule on, such as a read
// of a file; Shell runs a command line; FileWrite writes a file; FlaiTool
// calls one of flai's MCP tools.
const (
	Other Kind = iota
	Shell
	FileWrite
	FlaiTool
)

// Call is one tool call the guard decides on, in terms of no harness: who
// makes it, in which session, role, and story, and what it does. A harness's
// entry, such as Guard.Decide for Claude Code, reads its own hook input into
// a Call.
type Call struct {
	// Session is the harness session's ID.
	Session string
	// Subagent says a sub-agent makes the call, and Agent names its
	// definition, for refusals and the read-only lookup.
	Subagent bool
	Agent    string
	// Role and Story are the session's FLAI_ROLE and FLAI_STORY.
	Role  string
	Story string
	Kind  Kind
	// Name is the harness's own name for the tool, used only in refusals.
	Name string
	// Command is a Shell call's command line, Path the file a FileWrite
	// writes, and Tool a FlaiTool's name, without any prefix the harness
	// gives it, with Args its arguments.
	Command string
	Path    string
	Tool    string
	Args    Args
}

// Args are the arguments of a flai tool the guard reads: item_move's item
// and the state it moves the item to, plan's item, and the thread of
// thread_reply and thread_resolve; item_new's type and whether it makes a
// story a draft, and item_edit's draft; whether thread_reply is a
// recommendation and the source it cites. Fields are the names of every
// argument given, whatever the guard reads of it, so that an item_edit that
// finalizes a draft is told from one that changes more, and a draft false
// given from one left out. CostOfDelay is item_edit's cost_of_delay, each key
// it gives with its value as given, so that an edit that gives a story cost
// of delay inputs is told from one that sets a value or removes one (S-0328);
// it is nil when cost_of_delay is not an object. Commit is adr_new's commit,
// which a sub-agent may not give true (S-0287, ADR-0127).
type Args struct {
	ID             string                     `json:"id"`
	To             string                     `json:"to"`
	Type           string                     `json:"type"`
	Draft          bool                       `json:"draft"`
	Commit         bool                       `json:"commit"`
	Recommendation bool                       `json:"recommendation"`
	Source         string                     `json:"source"`
	Fields         []string                   `json:"-"`
	CostOfDelay    map[string]json.RawMessage `json:"-"`
}

// Outcome is what the guard decides on a call.
type Outcome int

// The outcomes: Allow lets the call through, Deny refuses it, and Ask puts
// it to a person: the harness's own permission tool asks, or Hold does.
const (
	Allow Outcome = iota
	Deny
	Ask
)

// Verdict is the guard's decision on a call: its Outcome, and for Deny and
// Ask the Reason. For the orchestrator's own refusals Call is the call
// refused, in one line, and Needs the permission that would allow it, as
// orchestration.permissions.<name>, or "" when none would. The zero Verdict
// allows the call.
type Verdict struct {
	Outcome Outcome
	Reason  string
	Call    string
	Needs   string
}

// deny is the Verdict that refuses a call for why, or allows it when why is
// "".
func deny(why string) Verdict {
	if why == "" {
		return Verdict{}
	}
	return Verdict{Outcome: Deny, Reason: why}
}

// Decide decides on c under g's policy: the shared paths and flai forge push
// for every caller, then the rules of a sub-agent, or of the session's own role, and the
// author a served agent writes as. A story's own agent's write to a path flai
// protects in its worktree, which they would allow, is put to the story's
// owner (Ask) while auto-approve is off.
func Decide(g Guard, c Call) Verdict {
	if v := g.shared(c); v.Outcome != Allow {
		return v
	}
	if v := g.forgePush(c); v.Outcome != Allow {
		return v
	}
	if c.Subagent {
		return g.asSubAgent(c)
	}
	if c.Kind == FlaiTool && c.Tool == ReleasePublish && c.Role != RoleOrchestrate && c.Role != RolePlan && c.Role != RoleAnalyze {
		return deny("only the orchestrator calls " + ReleasePublish + ": " + orchestratorPublishes)
	}
	var why string
	switch c.Role {
	case RolePlan:
		why = g.plan(c)
	case RoleAnalyze:
		why = g.analyze(c)
	case RoleOrchestrate:
		return g.orchestrate(c)
	case "":
		why = g.wait(c)
	}
	// a call its role refuses whoever it names is refused for the role's
	// reason, so that leaving out --by does not meet a second refusal
	if v := deny(cmp.Or(why, g.authors(c))); v.Outcome != Allow || c.Role != "" {
		return v
	}
	return g.ask(c)
}

// ask puts to the story's owner a write by a story's own agent to a path flai
// protects, other than one in .git, inside the story's worktree, while
// auto-approve is off (ADR-0086, ADR-0130); it allows every other call.
func (g Guard) ask(c Call) Verdict {
	if c.Kind != FileWrite || c.Story == "" || g.AutoApprove {
		return Verdict{}
	}
	rel, ok := inWorktree(c)
	if !ok || !protected.Path(rel) || protected.Git(rel) {
		return Verdict{}
	}
	story := workitem.CanonicalID(c.Story)
	return Verdict{Outcome: Ask, Reason: fmt.Sprintf("%s of %s waits on the story's owner: a path flai protects is written only with the owner's approval on a thread on %s (ADR-0086, ADR-0130)", c.Name, rel, story)}
}

// inWorktree is the path c writes relative to its story's worktree, slash
// separated, and whether it lies inside it: below the segment
// .flai-cache/worktrees/<story>/ of the cleaned path, where
// workitem.Repo.WorktreePath puts the worktree in any project.
func inWorktree(c Call) (string, bool) {
	if c.Story == "" || c.Path == "" {
		return "", false
	}
	segment := "/.flai-cache/worktrees/" + workitem.CanonicalID(c.Story) + "/"
	clean := filepath.ToSlash(filepath.Clean(c.Path))
	_, rel, ok := strings.Cut(clean, segment)
	return rel, ok && rel != ""
}

// askedWhy is why a person must answer a write the guard holds itself, in no
// harness's terms.
const askedWhy = "flai protects the path and does not let an agent write it without a person's approval (ADR-0086, ADR-0130)."

// askedChange is what the thread shows under the question: a Call carries
// the path it writes, not the content.
const askedChange = "The harness gave flai the path alone, not what the call would write."

// Hold decides on c under g's policy as Decide does and, when the verdict is
// Ask, holds the call: it asks the story's owner through a, within bound, the
// harness's, and allows the call or denies it with ask's reason. A reader
// whose harness calls no permission tool of flai's holds through it; any
// other verdict comes back as Decide gives it.
func Hold(ctx context.Context, g Guard, a ask.Asker, c Call, bound time.Duration) Verdict {
	v := Decide(g, c)
	if v.Outcome != Ask {
		return v
	}
	rel, _ := inWorktree(c)
	story, err := a.Repo.Get(c.Story)
	if err != nil {
		return deny(fmt.Sprintf("could not read %s for its owner, so %s was not written: %v", c.Story, rel, err))
	}
	if story.Status != workitem.InProgress {
		return deny(fmt.Sprintf("%s is %s, not in progress, so %s was not written: a path flai protects is asked for only in an in-progress story's worktree", story.ID, story.Status, rel))
	}
	req := ask.Request{Story: story.ID, Owner: story.Owner, Agent: g.selves(c)[0], Tool: c.Name, Path: rel, Why: askedWhy, Change: askedChange}
	if ok, why := a.Hold(ctx, req, bound); !ok {
		return deny(why)
	}
	return Verdict{}
}

// asSubAgent decides on a sub-agent's call: it reads, writes an ADR without
// committing it, writes a protected path only while auto-approve is on, and,
// when read-only, writes no file through the shell but under the temporary
// folder.
func (g Guard) asSubAgent(c Call) Verdict {
	who := cmp.Or(c.Agent, "unnamed")
	switch c.Kind {
	case FileWrite:
		if file, ok := protectedWrite(c); ok && !g.AutoApprove {
			return deny(fmt.Sprintf("a sub-agent (%s) cannot use %s on %s: a path flai protects, such as a file in a .claude/ folder or .mcp.json, is written only with the operator's approval on a thread, which would hold this call, and the layer with it, until they answer (ADR-0086, ADR-0106). Put the file's whole new content in your final message; the story's agent writes it.", who, c.Name, file))
		}
	case FlaiTool:
		switch {
		case c.Tool == ADRNew && c.Args.Commit:
			return deny(fmt.Sprintf("a sub-agent (%s) cannot call %s with commit true: %s; call it without commit %s. Put what you need done in your final message; the story's agent does it.", who, c.Tool, commitsADR, leavesADRToIt))
		case slices.Contains(MCPReads, c.Tool), slices.Contains(MCPStoryReads, c.Tool), c.Tool == ADRNew:
			return Verdict{}
		}
		return deny(fmt.Sprintf("a sub-agent (%s) cannot call %s: it reads, and only the story's agent changes work items, threads, and conversations with other stories or reads the inbox (ADR-0059). Put what you need done, or the question for the designer, in your final message.", who, c.Tool))
	case Shell:
		readOnly := g.ReadOnly != nil && g.ReadOnly(c.Agent)
		for _, s := range parse(c.Command) {
			run := strings.Join(s.words, " ")
			if why, _ := g.refuse(s.words, subAgent); why != "" {
				return deny(fmt.Sprintf("a sub-agent (%s) cannot run %q: %s. Put what you need done in your final message; the story's agent does it.", who, run, why))
			}
			if !readOnly {
				continue
			}
			if file := written(s); file != "" {
				return deny(fmt.Sprintf("a sub-agent (%s) cannot run %q: it reads, and it writes no file but under the temporary folder (ADR-0141); this writes %s. Put the correction in your final message; the story's agent makes it.", who, run, file))
			}
		}
	}
	return Verdict{}
}

// protectedWrite says whether c writes a protected path other than one in
// .git, and which file, as ProtectedWrite describes.
func protectedWrite(c Call) (string, bool) {
	if c.Kind != FileWrite || c.Path == "" {
		return "", false
	}
	clean := filepath.Clean(c.Path)
	return c.Path, protected.Path(clean) && !protected.Git(clean)
}
