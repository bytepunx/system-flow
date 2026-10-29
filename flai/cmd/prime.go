package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/spf13/cobra"

	ctxpack "github.com/bytepunx/system-flow/flai/internal/context"
	"github.com/bytepunx/system-flow/flai/internal/conventions"
	"github.com/bytepunx/system-flow/flai/internal/issues"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

func newPrimeCmd(a *app) *cobra.Command {
	var (
		cat     bool
		story   string
		budget  string
		role    string
		epic    string
		part    int
		reasons []string
	)
	c := &cobra.Command{
		Use:   "prime",
		Short: "Print the conventions an agent reads at session start, in order",
		Long: `List design/conventions in read order (README first, then by order) so an
agent, a hook, or a script can load them in one call. --cat prints the
content of each file with a header instead of the paths. Both, and --json,
give every convention whatever its roles, except design-system.md in a
project with no design system, that is, no guidance document at
design/system/design-system.md (ADR-0152).

--story S-nnnn prints the context pack for an agent working that story
(ADR-0047), fitted to a size budget (ADR-0049): a header naming the story,
its topics and where each came from, the pack's size against the budget, and
the size of each thing it prints; then each convention the story's agent
reads, those whose front matter roles are empty or list story (ADR-0068),
design-system.md only when the design system reaches the story (below),
as --cat prints it, with the sections whose topics include neither all nor
one of the story's left out (the front matter, the baseline marker, and the
Project additions heading always stay), never cut for the budget; the open
issues; everything the story, its epic, and its tasks link or name, whole
(a #fragment loads its section; a superseded ADR is replaced by what
supersedes it), except a document named only by its path written out and larger than an eighth of
the budget, which is briefed with the reason "named in <ID>" and a line
telling the agent to read it before relying on it or changing it
(ADR-0050); the design system, when the topics of its guidance include all
or one of the story's (ADR-0152): the guidance and the token file whole,
as if the story named them, with the reason "design system: topics <the
topics that matched>", each larger than an eighth of the budget briefed
with a line telling the agent to read it whole before building UI; a
brief of each design/system and design/tech file its topics
select, as its first paragraph and heading outline; the decision sentence of each ADR its topics
select or one link step reaches (linked from a named or briefed section,
refined by, or refining one); then the digest entries the story's topics,
other than all, or its claim select, never an absorbed one, ranked against
the story's title, goal, and criteria, loaded best first under a digest
heading, each with step digest and its reason (path: <claimed path> or
topics: <topics>, then rank n), a compacted one as a one-line brief, while
they stay within an eighth of the budget and, on a pack within its budget
before the step, the pack within it too (ADR-0194, TH-0624); then the sections ranked highest against the
story's title, goal, and criteria, each cut at its own heading (an ADR whole
when it fits), in rank order, while the budget has room. Then a catalog of
every document neither loaded nor briefed, with the outline of those loaded
in part, and one line per convention section left out, with its topics.
Nothing prints twice: the first reason wins and the others are listed on it.
When the conventions alone exceed the budget the pack is the conventions,
the digest entries within their eighth, and a catalog; when the conventions
and what is named exceed it, nothing is ranked, and the digest entries still
load within their eighth; the header says which. --budget sets the size, such as 80KB or 81920
bytes; the project's default is prime.budget in system-flow.yaml, and
flai's is 80KB. An archived story gets the pack it would get today.

--role explore or --role verify, with --story, prints the smaller pack for
a sub-agent the story's agent hands work to (ADR-0059): the conventions
whose front matter roles are empty or list the role (ADR-0068),
design-system.md only when the design system reaches the story, with the
sections the story's topics leave out taken out; the story's goal and
acceptance criteria; and briefs, never bodies, of what the story names,
the design system when it reaches the story, what its topics select, the ADRs one step reaches,
and the digest entries the story's pack selects, in its rank, each while the
budget has room, with a count of those left out. No open issues, nothing
ranked, no catalog. Its budget is half the story's agent's unless --budget
is given.

--role plan, orchestrate, analyze, or curate prints the pack for a
strategic agent, one that works above a story: the planner (plan), with
--epic E-nnnn or --story S-nnnn, and the orchestrator (orchestrate), the
analyzer (analyze), and the curator (curate), for the whole project, with
neither. Its topics are the role's, planning, orchestration, analysis, or
digest and planning, with the planner's item's as well; it
holds the conventions whose roles are empty or list the role, with the
sections its topics leave out taken out, and no README; the open issues;
for the planner, what the item names, whole, as a story's pack loads it (a
story's with its epic and tasks, an epic's alone), and the sections ranked highest against the item; briefs of the design, tech,
and ADRs its topics select and the ADRs one step reaches; for the
orchestrator, the digest's answers and recommendations ranked against the
title and first entry of each thread awaiting the operator, each thread's
best not yet taken in turn, with the reason thread TH-nnnn; rank n, within an
eighth of the budget (ADR-0194); and a catalog of the rest. It does not get the design system. Its budget is the story's agent's, and --budget sets it.

--role orchestrate --reason, repeatable, primes the orchestrator for its
run's reasons (ADR-0217) in place of the whole project. A reason is one of
` + ctxpack.ReasonForms() + `.
The pack loads the role's conventions once, with the sections those
reasons do not need left out: a section tagged all and other words stays
only when one of them is a reason's kind or the role's topic. Then, per
reason, what it needs, each item once however many reasons name it: for a
thread, the thread with its entries and the item it is on, as its title,
status, goal, and acceptance criteria (ADR-0218); for a review,
the story and its last flai verify result; for planned, the item, the
items under it, and the planner's open thread on it; for a draft, the
draft; for release, release_evaluate's figures; for ready, the board and
the promotion candidates. Then the digest's answers and recommendations
ranked against the reasons' threads, within a sixteenth of the budget. It
loads no briefs, no ranked sections, no open issues, and no catalog.
--reason is refused with another role and with --epic or --story, and an
unknown reason is refused naming the forms.

--role takes as well any role a project defines in a role file under
design/workflows/roles (ADR-0207), by its role: its pack is built as a
built-in role's is, from the topics its file names, the conventions whose
roles are empty or list it, and, when its file's item names epic or story,
the item given with --epic or --story; with item none, neither. A role
file for plan, orchestrate, analyze, or curate gives that role its topics
and item kinds in place of flai's own. A role no file defines, an item
kind its file does not take, and a role file that does not parse are
refused, naming the roles the project defines or the file's finding.

--part N prints part N of the pack as the MCP prime tool returns it, as
JSON: parts of at most 40,000 bytes each, every one carrying part and parts
(ADR-0104). --part implies --json; without it the whole pack prints.`,
		Example: `  flai prime
  flai prime --cat
  flai prime --json
  flai prime --story S-0136
  flai prime --story S-0136 --json
  flai prime --story S-0136 --budget 120KB
  flai prime --story S-0136 --role verify
  flai prime --story S-0136 --part 2
  flai prime --role plan --epic E-0016
  flai prime --role orchestrate --json
  flai prime --role orchestrate --reason thread:TH-0719 --reason review:S-0536`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			repo, err := a.project()
			if err != nil {
				return err
			}
			set, errs, err := conventions.Load(repo)
			if err != nil {
				return err
			}
			if set.Missing {
				return fmt.Errorf("no conventions folder at %s; render it from the template or run flai upgrade", relPath(repo.Root, set.Dir))
			}
			hasDesign, err := ctxpack.HasDesignSystem(repo)
			if err != nil {
				return err
			}
			if !hasDesign {
				set = set.WithoutDesignSystem()
			}
			readme := filepath.Join(set.Dir, "README.md")
			var paths []string
			if set.README != "" {
				paths = append(paths, readme)
			}
			for _, f := range set.Files {
				paths = append(paths, filepath.Join(repo.Root, f.Path))
			}
			for rel, e := range errs {
				a.logger().Warn("convention file unreadable", "component", "conventions", "path", rel, "err", e.Error())
			}
			strategicRoles, strategic := strategicRole(repo, role)
			withReasons := cmd.Flags().Changed("reason")
			switch {
			case cmd.Flags().Changed("part") && cat:
				return fmt.Errorf("--part prints a part of a context pack; --cat prints the conventions, give one of them")
			case cmd.Flags().Changed("part") && part < 1:
				return fmt.Errorf("--part %d: parts count from 1", part)
			case cmd.Flags().Changed("part") && story == "" && !strategic:
				return fmt.Errorf("--part prints a part of a context pack; give --story, or --role %s", strings.Join(strategicRoles, ", "))
			case withReasons && role != conventions.RoleOrchestrate:
				given := "no --role"
				if role != "" {
					given = "--role " + role
				}
				return fmt.Errorf("--reason primes the orchestrator for its run's reasons, and %s is not the orchestrator; give --role orchestrate and no --story or --epic", given)
			case withReasons && (story != "" || epic != ""):
				return fmt.Errorf("--reason primes the orchestrator for its run's reasons, which looks at the whole project, not one item; give --role orchestrate and no --story or --epic")
			case epic != "" && !strategic:
				return fmt.Errorf("--epic primes the planner for an epic; give --role plan too")
			case strategic:
				return a.primeStrategic(repo, role, epic, story, budget, reasons, part)
			case budget != "" && story == "":
				return fmt.Errorf("--budget sizes a story's context pack; give --story too")
			case role != "" && story == "":
				return fmt.Errorf("--role %s primes a sub-agent working a story; give --story too, or give a strategic role: %s", role, strings.Join(strategicRoles, ", "))
			}
			if story != "" {
				return a.primeStory(repo, story, role, budget, part)
			}
			list, _ := issues.List(repo)
			table := issues.SummaryTable(list)
			if a.jsonOut {
				return a.printJSON(map[string]any{"dir": relPath(repo.Root, set.Dir), "files": set.Files, "readme": set.README != "", "open_issues": len(strings.Split(strings.TrimSpace(table), "\n")) - 2})
			}
			if !cat {
				for _, p := range paths {
					fmt.Fprintln(a.out, relPath(repo.Root, p))
				}
				if table != "" {
					fmt.Fprintln(a.out, relPath(repo.Root, filepath.Join(issues.Dir(repo), issues.SummaryFile)))
				}
				return nil
			}
			for i, p := range paths {
				data, err := os.ReadFile(p)
				if err != nil {
					return err
				}
				if i > 0 {
					fmt.Fprintln(a.out)
				}
				rel := relPath(repo.Root, p)
				fmt.Fprintf(a.out, "%s\n%s\n\n", rel, strings.Repeat("=", len(rel)))
				fmt.Fprint(a.out, string(data))
			}
			if table != "" {
				fmt.Fprintf(a.out, "\nopen issues\n===========\n\n%s", table)
			}
			return nil
		},
	}
	c.Flags().BoolVar(&cat, "cat", false, "print file contents instead of paths")
	c.Flags().StringVar(&story, "story", "", "print the context pack for this story: the conventions its agent reads, the design, tech, and ADRs it selects, and a catalog of the rest")
	c.Flags().StringVar(&budget, "budget", "", "the size the context pack fits, such as 80KB (default: prime.budget in system-flow.yaml, else 80KB; half that with --role explore or verify)")
	c.Flags().StringVar(&role, "role", "", "print the pack for an agent in this role: explore or verify, a sub-agent of the story's agent (ADR-0059), with --story; or plan, orchestrate, analyze, curate, or a role a role file defines, a strategic agent")
	c.Flags().StringVar(&epic, "epic", "", "with --role plan, or a role whose role file takes an epic, print its pack for this epic")
	c.Flags().StringArrayVar(&reasons, "reason", nil, "with --role orchestrate, prime the orchestrator for this reason of its run (ADR-0217): "+ctxpack.ReasonForms()+"; repeatable")
	c.Flags().IntVar(&part, "part", 0, "print part N of the pack, from 1, as JSON, as the MCP prime tool returns it (ADR-0104)")
	return c
}

// strategicRole is the strategic roles the project defines (ADR-0207), as
// ctxpack.StrategicRoles lists them, and whether role is one of them. A role
// file that does not parse leaves the built-in roles listed, and takes every
// role but a sub-agent's as strategic, so that ctxpack.ForStrategic refuses
// it with the file's finding.
func strategicRole(repo *workitem.Repo, role string) ([]string, bool) {
	list, err := ctxpack.StrategicRoles(repo)
	if err != nil {
		list = conventions.StrategicRoles
	}
	if role == "" || slices.Contains(conventions.SubAgentRoles, role) {
		return list, false
	}
	return list, err != nil || slices.Contains(list, role)
}

// primeStory prints the context pack for a story (ADR-0047, ADR-0049), as
// ctxpack.ForStory builds it, or for a sub-agent in a role (ADR-0059), as
// ctxpack.ForRole does; with part, only that part of it (ADR-0104).
func (a *app) primeStory(repo *workitem.Repo, id, role, budget string, part int) error {
	var pack *ctxpack.Pack
	var err error
	if role != "" {
		pack, err = ctxpack.ForRole(repo, id, role, budget)
	} else {
		pack, err = ctxpack.ForStory(repo, id, budget)
	}
	if err != nil {
		return err
	}
	return a.printPack(pack, part)
}

// primeStrategic prints the pack for a strategic agent in role, for the
// planner's epic or story or, with neither, the whole project, as
// ctxpack.ForStrategic builds it, or, with reasons, the orchestrator's pack
// primed for them, as ctxpack.ForReasons builds it (ADR-0217); with part,
// only that part of it.
func (a *app) primeStrategic(repo *workitem.Repo, role, epic, story, budget string, reasons []string, part int) error {
	var pack *ctxpack.Pack
	var err error
	if len(reasons) > 0 {
		pack, err = ctxpack.ForReasons(repo, reasons, budget)
	} else {
		pack, err = ctxpack.ForStrategic(repo, role, epic, story, budget)
	}
	if err != nil {
		return err
	}
	return a.printPack(pack, part)
}

// printPack prints a context pack, as JSON with --json; given a part, it
// prints that part as JSON, as the MCP prime tool returns it.
func (a *app) printPack(pack *ctxpack.Pack, part int) error {
	if part != 0 {
		p, err := pack.Part(part, ctxpack.PartLimit)
		if err != nil {
			return err
		}
		data, err := toolJSON(p)
		if err != nil {
			return err
		}
		_, err = fmt.Fprintf(a.out, "%s\n", data)
		return err
	}
	if a.jsonOut {
		return a.printJSON(pack)
	}
	fmt.Fprint(a.out, pack.Header())
	fmt.Fprint(a.out, pack.Body())
	return nil
}

// toolJSON encodes v as the MCP go-sdk encodes a tool's structured output
// into the text of its result: json.Marshal, then, for an object, decoded
// and marshalled again after the output schema's defaults are applied, which
// sorts its keys. The pack's schema sets no defaults, so the second marshal
// changes the order of the keys only, never the size.
func toolJSON(v any) ([]byte, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	var decoded any
	if err := json.Unmarshal(data, &decoded); err != nil {
		return nil, err
	}
	return json.Marshal(decoded)
}
