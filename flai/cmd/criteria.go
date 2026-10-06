package cmd

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/bytepunx/system-flow/flai/internal/docedit"
	"github.com/bytepunx/system-flow/flai/internal/itemedit"
)

// flai criteria lists, ticks, and unticks a story's acceptance criteria (S-0282).
func newCriteriaCmd(a *app) *cobra.Command {
	c := &cobra.Command{
		Use:   "criteria",
		Short: "List, tick, and untick an item's acceptance criteria by their numbers",
		Long: `An item's acceptance criteria are the checkboxes under its
## Acceptance criteria heading, numbered from 1 in the order they appear, as
flai criteria list prints them. An agent ticks a criterion once it has
verified it, and not before (work-management.md); flai criteria tick changes
that one box and every other byte of the body stays as it was.

tick and untick are an edit of the body, made as flai edit makes one: an
archived or closed item is refused (exit 4), as is a change flai check finds
fault with; with --hash, a change someone made meanwhile is a conflict
(exit 3); and agents connected over MCP are told the criteria changed. A
number there is not is refused, and nothing is written.

The MCP tool criteria_tick and the dashboard tick and untick criteria the
same way.`,
		Example: `  flai criteria list S-0282
  flai criteria tick S-0282 1 3
  flai criteria tick S-0282 1,3 --autocommit
  flai criteria untick S-0282 2
  flai criteria list S-0282 --json`,
	}
	c.AddCommand(newCriteriaListCmd(a), newCriteriaTickCmd(a, true), newCriteriaTickCmd(a, false))
	return c
}

// newCriteriaListCmd prints an item's acceptance criteria; it reads, so an
// archived item is listed too.
func newCriteriaListCmd(a *app) *cobra.Command {
	return &cobra.Command{
		Use:   "list <id>",
		Short: "Print an item's acceptance criteria, numbered, and which are ticked",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			repo, err := a.project()
			if err != nil {
				return err
			}
			v, err := itemedit.Show(repo, args[0])
			if err != nil {
				return err
			}
			list := itemedit.Criteria(v.Body)
			if a.jsonOut {
				return a.printJSON(map[string]any{"id": v.ID, "path": v.Path, "hash": v.Hash, "criteria": list})
			}
			printCriteria(a.out, v.ID, list)
			return nil
		},
	}
}

// newCriteriaTickCmd is flai criteria tick, or untick when tick is false.
func newCriteriaTickCmd(a *app, tick bool) *cobra.Command {
	verb := "tick"
	if !tick {
		verb = "untick"
	}
	var hash, byFlag string
	var trailers []string
	var autocommit bool
	c := &cobra.Command{
		Use:   verb + " <id> <n>...",
		Short: strings.ToUpper(verb[:1]) + verb[1:] + " an item's acceptance criteria by the numbers flai criteria list prints",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ns, err := criterionNumbers(args[1:])
			if err != nil {
				return err
			}
			repo, err := a.project()
			if err != nil {
				return err
			}
			v, err := itemedit.Show(repo, args[0])
			if err != nil {
				return err
			}
			if !v.Editable {
				r := &docedit.RefusedError{Path: v.Path, Reason: v.Reason}
				if a.jsonOut {
					_ = a.printJSON(map[string]any{"refused": r})
				}
				return &exitError{code: exitDocRefused, msg: r.Error()}
			}
			var body string
			if tick {
				body, err = itemedit.Tick(v.Body, ns, nil)
			} else {
				body, err = itemedit.Tick(v.Body, nil, ns)
			}
			var inv *itemedit.InvalidError
			if errors.As(err, &inv) {
				return fmt.Errorf("rule: %s", inv.Error())
			}
			if err != nil {
				return err
			}
			before := itemedit.Criteria(v.Body)
			after := itemedit.Criteria(body)
			var moved []string
			for i := range after {
				if after[i].Ticked != before[i].Ticked {
					moved = append(moved, strconv.Itoa(after[i].N))
				}
			}
			if hash == "" {
				hash = v.Hash
			}
			by := a.writer()
			if byFlag != "" {
				by = byFlag
			}
			message := fmt.Sprintf("%s criteria %s", verb, strings.Join(moved, ", "))
			res, err := itemedit.Apply(repo, a.runner, args[0], itemedit.Change{Body: &body}, itemedit.Options{Hash: hash, By: by, Message: message, Trailers: trailers, NoCommit: !autocommit, Now: a.now()})
			if c, ok := docedit.IsConflict(err); ok {
				if a.jsonOut {
					_ = a.printJSON(map[string]any{"conflict": c})
				}
				return &exitError{code: exitDocConflict, msg: c.Error()}
			}
			if r, ok := docedit.IsRefused(err); ok {
				if a.jsonOut {
					_ = a.printJSON(map[string]any{"refused": r})
				} else {
					for _, fd := range r.Findings {
						fmt.Fprintf(a.out, "%s:%d: %s: %s: %s\n", fd.Path, fd.Line, fd.Level, fd.Rule, fd.Message)
					}
				}
				return &exitError{code: exitDocRefused, msg: r.Error()}
			}
			if errors.As(err, &inv) {
				return fmt.Errorf("rule: %s", inv.Error())
			}
			if err != nil {
				return err
			}
			now, err := itemedit.Show(repo, res.ID)
			if err != nil {
				return err
			}
			list := itemedit.Criteria(now.Body)
			if a.jsonOut {
				return a.printJSON(struct {
					*itemedit.Result
					Criteria []itemedit.Criterion `json:"criteria"`
				}{res, list})
			}
			if res.Unchanged {
				fmt.Fprintf(a.out, "%s: unchanged\n", res.ID)
			} else {
				fmt.Fprintf(a.out, "%s: %sed %s\n", res.ID, verb, strings.Join(moved, ", "))
			}
			switch {
			case res.Committed:
				fmt.Fprintf(a.out, "  committed %s (%d file(s))\n", res.Commit, len(res.Files))
			case res.CommitError != "":
				fmt.Fprintf(a.out, "  changed, not committed: %s\n", res.CommitError)
			}
			printCriteria(a.out, res.ID, list)
			return nil
		},
	}
	f := c.Flags()
	f.StringVar(&hash, "hash", "", "the hash flai criteria list --json or flai edit --show printed; a change made meanwhile is then a conflict")
	f.StringVar(&byFlag, "by", "", "who ticks, as agents are told (default: FLAI_AGENT, then the config author)")
	f.BoolVar(&autocommit, "autocommit", false, "commit the item, unless dashboard.autocommit is false")
	f.StringArrayVar(&trailers, "trailer", nil, "trailer line for the commit (repeatable)")
	return c
}

// criterionNumbers reads the numbers of criteria from words, each one number
// or several separated by commas.
func criterionNumbers(words []string) ([]int, error) {
	var out []int
	seen := map[int]bool{}
	for _, w := range words {
		for _, part := range strings.Split(w, ",") {
			part = strings.TrimSpace(part)
			if part == "" {
				continue
			}
			n, err := strconv.Atoi(part)
			if err != nil || n < 1 {
				return nil, fmt.Errorf("rule: %q is not the number of an acceptance criterion: give the numbers flai criteria list prints, such as 1 3 or 1,3", part)
			}
			if !seen[n] {
				seen[n] = true
				out = append(out, n)
			}
		}
	}
	return out, nil
}

// printCriteria prints id's criteria under a line that counts the ticked.
func printCriteria(w io.Writer, id string, list []itemedit.Criterion) {
	ticked := 0
	for _, c := range list {
		if c.Ticked {
			ticked++
		}
	}
	fmt.Fprintf(w, "%s: %d of %d ticked\n", id, ticked, len(list))
	for _, c := range list {
		mark := " "
		if c.Ticked {
			mark = "x"
		}
		fmt.Fprintf(w, "  %d [%s] %s\n", c.N, mark, c.Text)
	}
}
