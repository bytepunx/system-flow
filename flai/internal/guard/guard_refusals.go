package guard

import (
	"cmp"
	"fmt"
	"path"
	"path/filepath"
	"regexp"
	"slices"
	"strings"

	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// waitsOnSubAgents ends the refusal of a story's agent's wait_for_events: why
// it would not return, and how to wait instead.
const waitsOnSubAgents = "wait_for_events reports work items and threads, not sub-agents, so it would run to its timeout. Wait for a sub-agent by launching it with the Agent tool's run_in_background set to false, so that its result comes back as the tool's result; launch a layer's sub-agents in one message, each so. Use wait_for_events only for a thread awaiting the designer or the verified event of the story's close-out job (delegation.md)."

// wait says why a story's agent's own call is refused, or "" when it is not
// (S-0285): wait_for_events while a sub-agent of its session runs, unless a
// thread on its story or one of its tasks is open, so that it waits on the
// designer, whose answer the call returns on.
func (g Guard) wait(c Call) string {
	if c.Story == "" || c.Kind != FlaiTool || c.Tool != "wait_for_events" || g.Running == nil {
		return ""
	}
	running, err := g.Running()
	if err != nil || len(running) == 0 {
		return ""
	}
	if g.ThreadOpen != nil {
		if open, err := g.ThreadOpen(); err != nil || open {
			return ""
		}
	}
	names := make([]string, len(running))
	for i, s := range running {
		who := s.Type
		if who == "" {
			who = "unnamed"
		}
		names[i] = fmt.Sprintf("%s %s", who, s.ID)
	}
	return fmt.Sprintf("the story's agent cannot call wait_for_events while its sub-agents run (%s) and no thread on %s awaits the designer: %s", strings.Join(names, ", "), c.Story, waitsOnSubAgents)
}

// orchestrate words the refusal of the orchestrator's own call that its role
// file refuses as the orchestrator was refused under its permissions before
// the file held it, naming the call and the permission that would allow it,
// or is the zero Verdict when that rule allows the call.
func (g Guard) orchestrate(c Call) Verdict {
	switch c.Kind {
	case FileWrite:
		call := strings.TrimSpace(c.Name + " " + c.Path)
		return refusedOrchestrator(call, "use "+c.Name, nevers("it acts through flai and never edits a file"), "")
	case FlaiTool:
		tool, in := c.Tool, c.Args
		call, what := strings.Join(strings.Fields(tool+" "+in.ID+" "+in.To), " "), "call "+tool
		var needs, never string
		switch {
		case slices.Contains(MCPReads, tool), slices.Contains(MCPOrchestrates, tool):
			return Verdict{}
		case tool == "thread_reply", tool == "thread_resolve":
			tc := threadCall{verb: strings.TrimPrefix(tool, "thread_"), id: in.ID, recommendation: in.Recommendation, sourced: cites(in.Source)}
			why, permit := g.thread(tc)
			if why == "" {
				return Verdict{}
			}
			if tc.recommendation {
				call += " recommendation"
			}
			if tc.sourced {
				call += " source " + strings.TrimSpace(in.Source)
			}
			return refusedOrchestrator(call, threadDoes(tc), why, permit)
		case tool == "plan":
			what = "plan " + in.ID
			switch {
			case isEpic(in.ID):
				needs = manifest.PermitPlanBacklogEpics
			case isStory(in.ID):
				needs = manifest.PermitPlanBacklogStories
			default:
				never = plansCandidates
			}
		case tool == "item_edit":
			what = "edit " + in.ID
			switch {
			case finalizesDraft(in):
				what, needs = "finalize "+in.ID, manifest.PermitFinalizeDrafts
			case setsCostInputs(in):
				what, needs = "give "+in.ID+" cost of delay inputs", costNeeds(in.ID, g.Permissions)
			default:
				never = editsMCP
			}
		case tool == "item_move":
			what = fmt.Sprintf("move %s to %s", in.ID, in.To)
			switch in.To {
			case Ready:
				needs = manifest.PermitPromoteToReady
			case workitem.Done:
				never = acceptsWithFlai
			default:
				never = "it moves an item to ready and no further"
			}
		case tool == ReleasePublish:
			what, needs = "publish a release", manifest.PermitPublish
		case tool == QueueAdd:
			call = strings.TrimSpace(tool + " " + in.Story)
			what, needs = strings.TrimSpace("enqueue "+in.Story), manifest.PermitAcceptReviews
		case tool == SubprojectEdit:
			// it uses flai subproject --autocommit instead, which commits
			// the change on main; the tool is held to the same permission
			// all the same (ADR-0195)
			what, needs = "edit sub-projects", manifest.PermitEditSubprojects
		default:
			never = "of flai's tools that write, it calls only thread_open and activity_log, thread_reply and thread_resolve as the thread allows, and those its permissions allow"
		}
		switch {
		case never != "":
			return refusedOrchestrator(call, what, nevers(never), "")
		case !g.Permissions.Allows(needs):
			return refusedOrchestrator(call, what, off(needs), needs)
		}
	case Shell:
		r := g.orchestration()
		for _, words := range arguments(c.Command) {
			if why, needs := g.refuse(words, r); why != "" {
				call := strings.Join(words, " ")
				return refusedOrchestrator(call, fmt.Sprintf("run %q", call), why, needs)
			}
		}
	}
	return Verdict{}
}

// refusedOrchestrator is the orchestrator's refusal of call, said as what,
// for why; needs is the permission that would allow it, "" when none would.
func refusedOrchestrator(call, what, why, needs string) Verdict {
	v := Verdict{Outcome: Deny, Call: call, Reason: fmt.Sprintf("the orchestrator cannot %s: %s. %s", what, why, askOperator)}
	if needs != "" {
		v.Needs = "orchestration.permissions." + needs
	}
	return v
}

// off says the orchestrator is refused a call for want of the permission
// needs.
func off(needs string) string {
	return "it needs orchestration.permissions." + needs + ", which is off"
}

// nevers says the orchestrator never makes a call, whatever its
// permissions, for why.
func nevers(why string) string {
	return "the orchestrator never does it, whatever its permissions: " + why
}

// plan words the refusal of the planner's own call that its role file
// refuses as the planner was refused before the file held it, or is "" when
// that rule allows the call.
func (g Guard) plan(c Call) string {
	switch c.Kind {
	case FileWrite:
		return fmt.Sprintf("the planner cannot use %s: %s", c.Name, planner)
	case FlaiTool:
		switch tool, in := c.Tool, c.Args; {
		case tool == "item_new" && in.Type == "story" && !in.Draft:
			return fmt.Sprintf("the planner cannot create a story that is not a draft: %s; %s", drafts, planner)
		case slices.Contains(MCPReads, tool), slices.Contains(MCPPlans, tool):
			return ""
		case tool == "item_move" && in.To == Backlog:
			return ""
		case tool == "item_move":
			return fmt.Sprintf("the planner cannot move %s to %s: %s", in.ID, in.To, planner)
		}
		return fmt.Sprintf("the planner cannot call %s: %s", c.Tool, planner)
	case Shell:
		for _, words := range arguments(c.Command) {
			if why, _ := g.refuse(words, planning); why != "" {
				return fmt.Sprintf("the planner cannot run %q: %s; %s", strings.Join(words, " "), why, planner)
			}
		}
	}
	return ""
}

// analyze words the refusal of the analyzer's own call that its role file
// refuses as the analyzer was refused before the file held it (S-0223), or
// is "" when that rule allows the call.
func (g Guard) analyze(c Call) string {
	switch c.Kind {
	case FileWrite:
		if g.inReports(c.Path) {
			return ""
		}
		folder := "the design folder's analysis/"
		if g.Reports != "" {
			folder = strings.TrimSuffix(g.Reports, "/") + "/"
		}
		return fmt.Sprintf("the analyzer cannot use %s on %s: it edits only its report and the index, under %s; %s", c.Name, cmp.Or(c.Path, "no file"), folder, analyzer)
	case FlaiTool:
		switch {
		case slices.Contains(MCPReads, c.Tool), slices.Contains(MCPAnalyzes, c.Tool):
			return ""
		case slices.Contains(itemWrites, c.Tool):
			return fmt.Sprintf("the analyzer cannot call %s: %s; %s", c.Tool, authorsNoStories, analyzer)
		}
		return fmt.Sprintf("the analyzer cannot call %s: %s", c.Tool, analyzer)
	case Shell:
		for _, words := range arguments(c.Command) {
			if why, _ := g.refuse(words, analyzing); why != "" {
				return fmt.Sprintf("the analyzer cannot run %q: %s; %s", strings.Join(words, " "), why, analyzer)
			}
		}
	}
	return ""
}

// curate words the refusal of the curator's own call that its role file
// refuses as the curator was refused before the file held it (S-0431,
// ADR-0196), or is "" when that rule allows the call.
func (g Guard) curate(c Call) string {
	switch c.Kind {
	case FileWrite:
		return fmt.Sprintf("the curator cannot use %s on %s: %s", c.Name, cmp.Or(c.Path, "no file"), curator)
	case FlaiTool:
		if slices.Contains(MCPReads, c.Tool) || slices.Contains(MCPCurates, c.Tool) {
			return ""
		}
		return fmt.Sprintf("the curator cannot call %s: %s", c.Tool, curator)
	case Shell:
		for _, words := range arguments(c.Command) {
			if why, _ := g.refuse(words, curating); why != "" {
				return fmt.Sprintf("the curator cannot run %q: %s; %s", strings.Join(words, " "), why, curator)
			}
		}
	}
	return ""
}

// inReports says whether file, as a FileWrite call names it, is under the
// reports folder, as under says.
func (g Guard) inReports(file string) bool { return g.under(g.Reports, file) }

// under says whether file, as a FileWrite call names it, is under dir, a
// folder relative to Root and slash separated: taken from Root when it is
// relative, cleaned, and with the symbolic links along the part of it that
// exists resolved, as the folder's are, so that neither .. nor a link leads
// out of it. The folder itself is not under it, and with Root, dir, or file
// "" nothing is.
func (g Guard) under(dir, file string) bool {
	if g.Root == "" || dir == "" || file == "" {
		return false
	}
	if !filepath.IsAbs(file) {
		file = filepath.Join(g.Root, file)
	}
	folder := resolved(filepath.Join(g.Root, filepath.FromSlash(dir)))
	rel, err := filepath.Rel(folder, resolved(filepath.Clean(file)))
	return err == nil && rel != "." && rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator))
}

// resolved is p, clean and absolute, with the symbolic links of the longest
// part of it that exists resolved and the rest joined on as it stands.
func resolved(p string) string {
	rest := ""
	for {
		if r, err := filepath.EvalSymlinks(p); err == nil {
			return filepath.Join(r, rest)
		}
		parent := filepath.Dir(p)
		if parent == p {
			return filepath.Join(p, rest)
		}
		rest = filepath.Join(filepath.Base(p), rest)
		p = parent
	}
}

var (
	assignment = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*=`)
	subcommand = regexp.MustCompile(`^[a-z][a-z0-9-]*$`)
	shellFlags = regexp.MustCompile(`^-[a-zA-Z]*c[a-zA-Z]*$`)
)

// refuse says why r refuses one simple command, and the permission that
// would allow it, or "" when it does not. It does not trust the first word
// to be the program: wrappers such as env, sudo, timeout, xargs, and find
// -exec put it further along, so every word is looked at. A word is git run
// with a subcommand when a plain word follows it after git's own flags, and
// flai run with one when one of its commands follows; a shell given -c, or
// eval, has its script checked too.
func (g Guard) refuse(words []string, r rules) (why, needs string) {
	for i, w := range words {
		if assignment.MatchString(w) {
			continue
		}
		rest := words[i+1:]
		switch path.Base(w) {
		case "sh", "bash", "zsh", "dash", "ksh":
			for j, f := range rest {
				if shellFlags.MatchString(f) && j+1 < len(rest) {
					if why, needs := g.script(rest[j+1], r); why != "" {
						return why, needs
					}
					break
				}
			}
		case "eval":
			if why, needs := g.script(strings.Join(rest, " "), r); why != "" {
				return why, needs
			}
		case "flai", "flai.sh":
			cmd, sub := subcommands(rest, flaiValues)
			if !slices.Contains(g.Commands, cmd) || slices.Contains(rest, "--help") || slices.Contains(rest, "-h") {
				continue
			}
			if why, needs := r.flai(cmd, sub, rest); why != "" {
				return why, needs
			}
		case "git":
			cmd, args := gitSubcommand(rest)
			if r.git == nil || !subcommand.MatchString(cmd) || slices.Contains(gitReads, cmd) || gitListsOnly(cmd, args) {
				continue
			}
			return r.git(cmd), ""
		}
	}
	return "", ""
}

// script checks each simple command of a script a shell is given against r.
func (g Guard) script(text string, r rules) (why, needs string) {
	for _, c := range arguments(text) {
		if why, needs := g.refuse(c, r); why != "" {
			return why, needs
		}
	}
	return "", ""
}

// subcommands is the first two words after a program's own flags; a flag in
// takesValue consumes the word after it.
func subcommands(words []string, takesValue map[string]bool) (cmd, sub string) {
	var found []string
	for i := 0; i < len(words) && len(found) < 2; i++ {
		w := words[i]
		if strings.HasPrefix(w, "-") {
			if takesValue[w] && len(found) == 0 {
				i++
			}
			continue
		}
		found = append(found, w)
	}
	if len(found) > 0 {
		cmd = found[0]
	}
	if len(found) > 1 {
		sub = found[1]
	}
	return cmd, sub
}

// positionals are the words that are not flags; a flag in takesValue
// consumes the word after it, wherever it stands.
func positionals(words []string, takesValue map[string]bool) []string {
	var out []string
	for i := 0; i < len(words); i++ {
		if strings.HasPrefix(words[i], "-") {
			if takesValue[words[i]] {
				i++
			}
			continue
		}
		out = append(out, words[i])
	}
	return out
}
