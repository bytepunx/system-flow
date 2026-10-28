package guard

import (
	"cmp"
	"fmt"
	"path"
	"path/filepath"
	"slices"
	"strings"

	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/roles"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// A role a project defines in a role file under the design folder's
// workflows/roles is held to that file (S-0416, ADR-0207 §5): the reads every
// sub-agent has, and each MCP tool, flai command, and file edit its file
// lists, under the condition an entry names and, for an entry under a
// permissions entry, while the operator's setting is on. An entry flai check
// reports as roles.grant, such as one that reaches a call no setting allows,
// is never granted. A role whose file is missing or does not parse reads and
// is refused the rest.
//
// The planner, the analyzer, and the curator are held to their role files
// the same way, or, when the project has none, to flai's baseline copy of
// it; their rules in Go only word a refusal as before the files held them.

// DefinedRole says whether role is one a project's role file defines, which
// the guard holds to the file: neither "", nor a built-in role, nor the verb
// of the story's agent or its sub-agents, which no role file takes.
func DefinedRole(role string) bool {
	switch role {
	case "", RolePlan, RoleOrchestrate, RoleAnalyze, RoleCurate:
		return false
	}
	return !slices.Contains(roles.ReservedRoles, role)
}

// explains are the built-in roles held to their role file, each with the
// rule in Go that words a refusal of the file's as the role was refused
// before the file held it, "" when that rule would allow the call.
var explains = map[string]func(Guard, Call) string{
	RolePlan:    Guard.plan,
	RoleAnalyze: Guard.analyze,
	RoleCurate:  Guard.curate,
}

// BuiltInFromFile says whether role is a built-in role the guard holds to its
// role file: the planner's, the analyzer's, or the curator's.
func BuiltInFromFile(role string) bool {
	_, ok := explains[role]
	return ok
}

// HeldByFile says whether the guard holds role's own calls to a role file:
// a defined role's or a built-in one's.
func HeldByFile(role string) bool { return DefinedRole(role) || BuiltInFromFile(role) }

// costFlags are the flags of flai edit that set or clear a cost of delay
// input, which ungated-edit leaves out; a value it lets through.
var costFlags = append(slices.Clone(costInputFlags), "--clear-cost-of-delay")

// askOperatorOnRole ends a defined role's refusal; fixesRoleFile ends it
// while its file cannot hold it.
const (
	askOperatorOnRole = "Ask the operator with thread_open on the item if it needs doing"
	fixesRoleFile     = "Say so in your final summary, for the operator to fix the role file"
)

// held is a defined role as the guard holds it: its file, the lines of the
// entries flai check reports as roles.grant, and why the file cannot hold
// it, "" when it can.
type held struct {
	g       Guard
	file    roles.Role
	banned  map[int]bool
	problem string
}

// roleHeld is the role, by its verb, as g holds it: from RoleFile, else, for
// a built-in role, from flai's baseline file of it, or, while RoleProblem is
// set or there is no file, a file that lists nothing, named after the role.
func (g Guard) roleHeld(role string) held {
	g.Design = g.designFolder()
	file := g.RoleFile
	h := held{g: g, file: roles.Role{Name: role, Role: role}, problem: g.RoleProblem}
	if BuiltInFromFile(role) {
		base, _ := roles.BaselineFor(g.Design, role)
		h.file.Name = base.Name
		if file == nil {
			file = &base
		}
	}
	switch {
	case h.problem != "":
	case file == nil:
		h.problem = fmt.Sprintf("no role file defines the role %s; run flai check", role)
	default:
		h.file, h.banned = *file, map[int]bool{}
		for _, p := range roles.Validate([]roles.Role{h.file}, g.Design, roles.Vocabulary{}) {
			if p.Code == roles.CodeGrant {
				h.banned[p.Line] = true
			}
		}
	}
	return h
}

// designFolder is the design folder, for {design} in a role file's files:
// Design, else the parent of Reports, the analyzer's reports folder, the
// design folder's analysis/, when it is given; "" when neither is.
func (g Guard) designFolder() string {
	if g.Design != "" || g.Reports == "" {
		return g.Design
	}
	return path.Dir(strings.TrimSuffix(g.Reports, "/"))
}

// asRole decides on the own call of a session whose role a role file
// defines, or of a built-in role's, held to its file: a refusal of the
// file's is worded by the built-in role's rule in Go when that rule refuses
// the call too, as a plain refusal, and by the file otherwise.
func (g Guard) asRole(c Call) Verdict {
	h := g.roleHeld(c.Role)
	var call, what, why, needs string
	switch c.Kind {
	case FlaiTool:
		if slices.Contains(MCPReads, c.Tool) {
			return Verdict{}
		}
		call, what = toolCall(c), toolDoes(c)
		why, needs = h.tool(c.Tool, c.Args)
	case Shell:
		r := h.rules()
		for _, words := range commands(c.Command) {
			if why, needs = g.refuse(words, r); why != "" {
				call = strings.Join(words, " ")
				what = fmt.Sprintf("run %q", call)
				break
			}
		}
	case FileWrite:
		call, what = strings.TrimSpace(c.Name+" "+c.Path), fmt.Sprintf("use %s on %s", c.Name, cmp.Or(c.Path, "no file"))
		why = h.edit(c.Path)
	}
	if why == "" {
		return Verdict{}
	}
	if explain, ok := explains[c.Role]; ok && h.problem == "" {
		if reason := explain(g, c); reason != "" {
			return deny(reason)
		}
	}
	return h.refused(call, what, why, needs)
}

// refused is the role's refusal of call, said as what, for why; needs is the
// setting that would allow it, "" when none would. While the file cannot
// hold the role, why is the file's problem.
func (h held) refused(call, what, why, needs string) Verdict {
	end, source := askOperatorOnRole, "ADR-0207"
	if h.file.Path != "" {
		source = h.file.Path + ", " + source
	}
	if h.problem != "" {
		why, needs, end = h.problem, "", fixesRoleFile
	}
	return Verdict{Outcome: Deny, Call: call, Needs: needs, Reason: fmt.Sprintf("the role %s cannot %s: %s. %s (%s).", h.file.Name, what, why, end, source)}
}

// toolCall is a flai tool's call in one line, as the orchestrator's
// refusals write it.
func toolCall(c Call) string {
	in := c.Args
	call := strings.Join(strings.Fields(strings.Join([]string{c.Tool, in.ID, in.Story, in.To}, " ")), " ")
	if in.Recommendation {
		call += " recommendation"
	}
	if cites(in.Source) {
		call += " source " + strings.TrimSpace(in.Source)
	}
	return call
}

// toolDoes is what a flai tool's call does, said after "cannot".
func toolDoes(c Call) string {
	if c.Tool == "item_move" {
		return fmt.Sprintf("move %s to %s", c.Args.ID, c.Args.To)
	}
	return "call " + c.Tool
}

// entry is a call entry of a role file, with the setting of the permissions
// entry it stands under, "" for one under none.
type entry struct {
	call    roles.Call
	setting string
}

// entries are the file's call entries in tools, or in commands, its
// permissions' among them, but those it never grants: one flai check reports
// as roles.grant, or one that reaches a call no setting allows.
func (h held) entries(commands bool) []entry {
	var out []entry
	add := func(es []roles.Call, setting string) {
		for _, e := range es {
			if _, never := roles.Forbids(e.Call, e.When, h.file.WorksStory); !never && !h.banned[e.Line] {
				out = append(out, entry{e, setting})
			}
		}
	}
	pick := func(tools, cmds []roles.Call) []roles.Call {
		if commands {
			return cmds
		}
		return tools
	}
	add(pick(h.file.Tools, h.file.Commands), "")
	for _, p := range h.file.Permissions {
		if !h.banned[p.Line] {
			add(pick(p.Tools, p.Commands), p.Setting)
		}
	}
	return out
}

// setting is the permissions a setting names and whether it is on: "" is
// none, and on; one is on only while it is the role's own and the operator
// has it on, the orchestrator's in Permissions and another role's in
// RolePermissions.
func (h held) setting(s string) (manifest.Permissions, bool) {
	if s == "" {
		return manifest.Permissions{}, true
	}
	perm, owner, ok := roles.ParseSetting(s)
	switch {
	case !ok || owner == "" || owner != h.file.Name:
		return manifest.Permissions{}, false
	case owner == roles.Orchestrator:
		return h.g.Permissions, h.g.Permissions.Allows(perm)
	}
	return h.g.RolePermissions, h.g.RolePermissions.Allows(perm)
}

// judge says why the entries refuse a call, and the setting that would allow
// it; both "" when one allows it. An entry that names the call, as names
// says, allows it while its setting is on and its condition, if it names
// one, qualifies the entry and holds, as holds says under the setting's
// permissions.
func (h held) judge(es []entry, names func(string) bool, holds func(cond string, p manifest.Permissions) bool) (why, needs string) {
	var offs []string
	var unmet []roles.Cond
	for _, e := range es {
		if !names(e.call.Call) {
			continue
		}
		cond, known := roles.Condition(e.call.When)
		if e.call.When != "" && (!known || !cond.Qualifies(e.call.Call)) {
			continue
		}
		p, on := h.setting(e.setting)
		met := e.call.When == "" || holds(e.call.When, p)
		switch {
		case met && on:
			return "", ""
		case met:
			offs = append(offs, e.setting)
		case !slices.ContainsFunc(unmet, func(u roles.Cond) bool { return u.Name == cond.Name }):
			unmet = append(unmet, cond)
		}
	}
	switch {
	case len(offs) > 0:
		return "it needs " + offs[0] + ", which is off", offs[0]
	case len(unmet) > 0:
		under := make([]string, len(unmet))
		for i, u := range unmet {
			under[i] = fmt.Sprintf("%s (%s)", u.Name, u.Allows)
		}
		return "it may make it only under " + strings.Join(under, ", or "), ""
	}
	return "its role file lists no such call", ""
}

// tool says why the role's file refuses a flai tool with its arguments, and
// the setting that would allow it; both "" when it allows it.
func (h held) tool(tool string, in Args) (why, needs string) {
	return h.judge(h.entries(false),
		func(call string) bool { return call == tool },
		func(cond string, p manifest.Permissions) bool { return h.toolHolds(cond, p, tool, in) })
}

// toolHolds says whether the condition named cond holds for a flai tool with
// its arguments, under p, the permissions of the setting the entry stands
// under.
func (h held) toolHolds(cond string, p manifest.Permissions, tool string, in Args) bool {
	tc := threadCall{verb: strings.TrimPrefix(tool, "thread_"), id: in.ID, recommendation: in.Recommendation, sourced: cites(in.Source)}
	switch cond {
	case "draft-only":
		return in.Type != "story" || in.Draft
	case "to-backlog":
		return in.To == Backlog
	case "to-ready":
		return in.To == Ready
	case "ungated-edit":
		// draft false given finalizes a draft; a cost of delay value is
		// not an input
		finalizing := slices.Contains(in.Fields, "draft") && !in.Draft
		inputs := slices.ContainsFunc(costInputs, func(k string) bool { _, ok := in.CostOfDelay[k]; return ok })
		return !finalizing && !inputs && !slices.Contains(in.Fields, "clear_cost_of_delay")
	case "own-thread":
		return h.ownThread(tc)
	case "planner-thread":
		return storyPlanner(h.g.opener(tc.id))
	case "answer-threads":
		return tc.verb == "reply" && h.answers(tc, p)
	case "plan-epic":
		return isEpic(in.ID)
	case "plan-story":
		return isStory(in.ID)
	case "cost-inputs":
		return setsCostInputs(in)
	case "finalize-draft":
		return finalizesDraft(in)
	}
	return false
}

// rules are the role's rules on a command line: flai's reads and the
// commands its file lists, and git's reads.
func (h held) rules() rules {
	es := h.entries(true)
	return rules{
		flai: func(cmd, sub string, rest []string) (string, string) {
			if reads(cmd, sub, rest) {
				return "", ""
			}
			return h.judge(es,
				func(call string) bool { return h.matches(call, cmd, sub) },
				func(cond string, p manifest.Permissions) bool { return h.commandHolds(cond, p, cmd, sub, rest) })
		},
		git: func(string) string { return "it runs only git's reads" },
	}
}

// matches says whether a command entry names flai's command cmd run with
// sub, the word after it (ADR-0207 §5): flai X names X with none of its
// subcommands, or in any form when it has none; flai X Y names its
// subcommand Y; and flai X * names X in any form.
func (h held) matches(call, cmd, sub string) bool {
	rest, ok := strings.CutPrefix(call, "flai ")
	words := strings.Fields(rest)
	if !ok || len(words) == 0 || words[0] != cmd {
		return false
	}
	switch {
	case len(words) == 1:
		return !slices.Contains(h.subcommandsOf(cmd), sub)
	case len(words) == 2 && words[1] == "*":
		return true
	case len(words) == 2:
		return words[1] == sub
	}
	return false
}

// subcommandsOf are the subcommands of flai's command cmd: from Subcommands
// when it is given, else those cliReads, the role's file, the calls no
// setting allows, and the gated calls name.
func (h held) subcommandsOf(cmd string) []string {
	if h.g.Subcommands != nil {
		return h.g.Subcommands[cmd]
	}
	var out []string
	add := func(call string) {
		if w := strings.Fields(call); len(w) == 3 && w[0] == "flai" && w[1] == cmd && w[2] != "*" {
			out = append(out, w[2])
		}
	}
	for _, s := range cliReads[cmd] {
		if s != "" {
			out = append(out, s)
		}
	}
	for _, e := range h.file.Commands {
		add(e.Call)
	}
	for _, p := range h.file.Permissions {
		for _, e := range p.Commands {
			add(e.Call)
		}
	}
	for _, f := range roles.Never {
		add(f.Call)
	}
	for _, gt := range roles.Gated {
		add(gt.Call)
	}
	return out
}

// commandHolds says whether the condition named cond holds for flai's
// command cmd, with its subcommand sub and the words after flai, under p,
// the permissions of the setting the entry stands under.
func (h held) commandHolds(cond string, p manifest.Permissions, cmd, sub string, rest []string) bool {
	moves := positionals(rest, moveValues)
	moveTo := ""
	if cmd == "move" && len(moves) > 2 {
		moveTo = moves[2]
	}
	plans := positionals(rest, flaiValues)
	switch cond {
	case "draft-only":
		return drafted(rest)
	case "to-backlog":
		return moveTo == Backlog
	case "to-ready":
		return moveTo == Ready
	case "ungated-edit":
		return !slices.ContainsFunc(rest, finalizes) && !slices.ContainsFunc(costFlags, func(f string) bool { return named(rest, f) })
	case "own-thread":
		return h.ownThread(cliThread(sub, rest))
	case "planner-thread":
		return storyPlanner(h.g.opener(cliThread(sub, rest).id))
	case "answer-threads":
		return sub == "reply" && h.answers(cliThread(sub, rest), p)
	case "plan-epic":
		return len(plans) > 1 && isEpic(plans[1])
	case "plan-story":
		return len(plans) > 1 && isStory(plans[1])
	case "cost-inputs":
		return setsCostInputsOnly(rest)
	case "finalize-draft":
		return finalizesOnly(rest)
	case "by-policy":
		return ordersByPolicy(rest)
	case "accept-as-self":
		accepts := cmd == "accept" && !named(rest, "--batch") || moveTo == workitem.Done && isStory(moves[1])
		return accepts && h.self(value(rest, "--by"))
	case "enqueue-as-self":
		return cmd == "queue" && sub == "add" && h.self(value(rest, "--by"))
	}
	return false
}

// ownThread says whether a thread call is on a thread the role opened and,
// for a reply, is neither a recommendation nor sourced: a follow-up.
func (h held) ownThread(tc threadCall) bool {
	return h.self(h.g.opener(tc.id)) && !tc.recommendation && !tc.sourced
}

// answers says whether a reply on a thread the role did not open is one the
// setting's answer mode, in p, allows: a recommendation, or while it is
// autonomous an answer citing its source too.
func (h held) answers(tc threadCall, p manifest.Permissions) bool {
	return !h.self(h.g.opener(tc.id)) && (tc.recommendation || tc.sourced && p.AnswerMode() == manifest.AnswerAutonomous)
}

// self says whether name is one of the role's own: its name, the name flai
// serve gives its run on an item, <name>-<Item>, or the hook's FLAI_AGENT;
// "" never is.
func (h held) self(name string) bool {
	if name == "" {
		return false
	}
	names := []string{h.file.Name, h.g.Agent}
	if h.g.Item != "" {
		names = append(names, h.file.Name+"-"+workitem.CanonicalID(h.g.Item))
	}
	return slices.Contains(names, name)
}

// edit says why the role's file refuses an edit of file, or "" when a files
// entry holds it: a folder under which the file lies, as under resolves it,
// with {design} the design folder. An entry that is the root, leaves it, or
// reaches .git, .claude, or the role files' workflows folder holds nothing,
// so that no role widens itself.
func (h held) edit(file string) string {
	guarded := []string{".git", ".claude", path.Join(h.g.Design, path.Dir(roles.Folder))}
	var folders []string
	for _, f := range h.file.Files {
		if strings.Contains(f, "{design}") && h.g.Design == "" {
			continue
		}
		dir := path.Clean(strings.ReplaceAll(f, "{design}", h.g.Design))
		if dir == "." || !filepath.IsLocal(filepath.FromSlash(dir)) || slices.ContainsFunc(guarded, func(g string) bool {
			return dir == g || strings.HasPrefix(dir, g+"/") || strings.HasPrefix(g, dir+"/")
		}) {
			continue
		}
		if h.g.under(dir, file) {
			return ""
		}
		folders = append(folders, dir+"/")
	}
	if len(folders) == 0 {
		return "its role file names no folder it may edit"
	}
	return "it edits only a file under " + strings.Join(folders, " or ")
}
