package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	ctxpack "github.com/bytepunx/system-flow/flai/internal/context"
	"github.com/bytepunx/system-flow/flai/internal/conventions"
	"github.com/bytepunx/system-flow/flai/internal/issues"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

func newPrimeCmd(a *app) *cobra.Command {
	var (
		cat    bool
		story  string
		budget string
		role   string
	)
	c := &cobra.Command{
		Use:   "prime",
		Short: "Print the conventions an agent reads at session start, in order",
		Long: `List design/conventions in read order (README first, then by order) so an
agent, a hook, or a script can load them in one call. --cat prints the
content of each file with a header instead of the paths. Both give every
convention, whatever its roles.

--story S-nnnn prints the context pack for an agent working that story
(ADR-0047), fitted to a size budget (ADR-0049): a header naming the story,
its topics and where each came from, the pack's size against the budget, and
the size of each thing it prints; then each convention the story's agent
reads, those whose front matter roles are empty or list story (ADR-0068),
as --cat prints it, with the sections whose topics include neither all nor
one of the story's left out (the front matter, the baseline marker, and the
Project additions heading always stay), never cut for the budget; the open
issues; everything the story, its epic, and its tasks link or name, whole
(a #fragment loads its section; a superseded ADR is replaced by what
supersedes it), except a document named only by its path written out and larger than an eighth of
the budget, which is briefed with the reason "named in <ID>" and a line
telling the agent to read it before relying on it or changing it
(ADR-0050); a brief of each design/system and design/tech file its topics
select, as its first paragraph and heading outline; the decision sentence of each ADR its topics
select or one link step reaches (linked from a named or briefed section,
refined by, or refining one); then the sections ranked highest against the
story's title, goal, and criteria, each cut at its own heading (an ADR whole
when it fits), in rank order, while the budget has room. Then a catalog of
every document neither loaded nor briefed, with the outline of those loaded
in part, and one line per convention section left out, with its topics.
Nothing prints twice: the first reason wins and the others are listed on it.
When the conventions alone exceed the budget the pack is the conventions and
a catalog; when the conventions and what is named exceed it, nothing is
ranked; the header says which. --budget sets the size, such as 80KB or 81920
bytes; the project's default is prime.budget in system-flow.yaml, and
flai's is 80KB. An archived story gets the pack it would get today.

--role explore or --role verify, with --story, prints the smaller pack for
a sub-agent the story's agent hands work to (ADR-0059): the conventions
whose front matter roles are empty or list the role (ADR-0068), with the
sections the story's topics leave out taken out; the story's goal and
acceptance criteria; and briefs, never bodies, of what the story names,
what its topics select, and the ADRs one step reaches, each while the
budget has room, with a count of those left out. No open issues, nothing
ranked, no catalog. Its budget is half the story's agent's unless --budget
is given.`,
		Example: `  flai prime
  flai prime --cat
  flai prime --json
  flai prime --story S-0136
  flai prime --story S-0136 --json
  flai prime --story S-0136 --budget 120KB
  flai prime --story S-0136 --role verify`,
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
			if budget != "" && story == "" {
				return fmt.Errorf("--budget sizes a story's context pack; give --story too")
			}
			if role != "" && story == "" {
				return fmt.Errorf("--role primes a sub-agent working a story; give --story too")
			}
			if story != "" {
				return a.primeStory(repo, story, role, budget)
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
	c.Flags().StringVar(&budget, "budget", "", "the size the story's context pack fits, such as 80KB (default: prime.budget in system-flow.yaml, else 80KB; half that with --role)")
	c.Flags().StringVar(&role, "role", "", "print the pack for a sub-agent of the story's agent in this role: explore or verify (ADR-0059)")
	return c
}

// primeStory prints the context pack for a story (ADR-0047, ADR-0049), as
// ctxpack.ForStory builds it, or for a sub-agent in a role (ADR-0059), as
// ctxpack.ForRole does.
func (a *app) primeStory(repo *workitem.Repo, id, role, budget string) error {
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
	if a.jsonOut {
		return a.printJSON(pack)
	}
	fmt.Fprint(a.out, pack.Header())
	fmt.Fprint(a.out, pack.Body())
	return nil
}
