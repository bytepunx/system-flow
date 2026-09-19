package cmd

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/bytepunx/system-flow/flai/internal/adr"
	"github.com/bytepunx/system-flow/flai/internal/docedit"
)

func newAdrCmd(a *app) *cobra.Command {
	c := &cobra.Command{
		Use:   "adr",
		Short: "Record architecture decisions",
	}
	c.AddCommand(newAdrNewCmd(a), newAdrAcceptCmd(a))
	return c
}

// adrNumber reads 27, 0027, or ADR-0027.
func adrNumber(s string) (int, error) {
	n, err := strconv.Atoi(strings.TrimPrefix(strings.ToUpper(strings.TrimSpace(s)), "ADR-"))
	if err != nil || n <= 0 {
		return 0, fmt.Errorf("%q is not an ADR number; write 27, 0027, or ADR-0027", s)
	}
	return n, nil
}

func adrNumbers(in []string) ([]int, error) {
	var out []int
	for _, s := range in {
		n, err := adrNumber(s)
		if err != nil {
			return nil, err
		}
		out = append(out, n)
	}
	return out, nil
}

// refusedADR prints a refusal the way flai doc save does and returns exit 4.
func (a *app) refusedADR(err error) (error, bool) {
	r, ok := docedit.IsRefused(err)
	if !ok {
		return nil, false
	}
	if a.jsonOut {
		_ = a.printJSON(map[string]any{"refused": r})
	} else {
		for _, f := range r.Findings {
			fmt.Fprintf(a.out, "%s:%d: %s: %s: %s\n", f.Path, f.Line, f.Level, f.Rule, f.Message)
		}
	}
	return &exitError{code: exitDocRefused, msg: r.Error()}, true
}

func (a *app) printADR(res *adr.Result, verb string) error {
	if a.jsonOut {
		return a.printJSON(res)
	}
	fmt.Fprintf(a.out, "%s %s (%s)\n  %s\n", res.ID, res.Title, res.Status, res.Path)
	switch {
	case res.Committed:
		fmt.Fprintf(a.out, "  %s and committed %s\n", verb, res.Commit)
	case res.CommitError != "":
		fmt.Fprintf(a.out, "  %s, NOT committed (%s)\n", verb, firstLine(res.CommitError))
	}
	return nil
}

func newAdrNewCmd(a *app) *cobra.Command {
	var opt adr.Options
	var supersedes, refines []string
	var bodyStdin, printBody bool
	c := &cobra.Command{
		Use:   "new \"<title>\"",
		Short: "Record a new ADR: number, file, front matter, and index row are supplied",
		Long: `Record an architecture decision. The number is one more than the highest
NNNN-*.md present in design/adrs (gaps are not filled). The file is
NNNN-slug.md with id, title, status, date, supersedes, superseded_by, and
refines when given; its row is added to design/adrs/README.md; each ADR it
supersedes gets superseded_by set, the one edit allowed to an accepted ADR.

The body below the heading is the project's 0000-template.md sections, or
standard input with --body-stdin. flai check runs with everything in place:
if it reports anything the ADR introduces, every file is put back, the
findings are printed, and the exit code is 4. --autocommit commits the new
file, the index, and any superseded ADRs on their own, unless the project
sets dashboard.autocommit: false. Nothing is pushed. --print-body prints the
template's sections and creates nothing.`,
		Example: `  flai adr new "Dashboards authenticate with a project token" --status accepted --refines 16
  flai adr new "Replace the SPA with server rendering" --supersedes 7 --body-stdin --autocommit < decision.md
  flai adr new --print-body`,
		Args: func(cmd *cobra.Command, args []string) error {
			if printBody {
				return cobra.NoArgs(cmd, args)
			}
			return cobra.ExactArgs(1)(cmd, args)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			repo, err := a.project()
			if err != nil {
				return err
			}
			if printBody {
				body := adr.TemplateBody(repo)
				if a.jsonOut {
					return a.printJSON(map[string]string{"body": body})
				}
				fmt.Fprint(a.out, body)
				return nil
			}
			opt.Title, opt.Now = args[0], a.now()
			if opt.Supersedes, err = adrNumbers(supersedes); err != nil {
				return err
			}
			if opt.Refines, err = adrNumbers(refines); err != nil {
				return err
			}
			if bodyStdin {
				data, err := io.ReadAll(cmd.InOrStdin())
				if err != nil {
					return err
				}
				if strings.TrimSpace(string(data)) == "" {
					return fmt.Errorf("--body-stdin was given and standard input is empty; write the decision, or leave the flag out for the template's sections")
				}
				opt.Body = string(data)
			}
			res, err := adr.New(repo, a.runner, opt)
			if rerr, ok := a.refusedADR(err); ok {
				return rerr
			}
			if err != nil {
				return err
			}
			return a.printADR(res, "recorded")
		},
	}
	f := c.Flags()
	f.StringVar(&opt.Status, "status", "proposed", "proposed or accepted")
	f.StringSliceVar(&supersedes, "supersedes", nil, "ADR this one supersedes (repeatable): 7, 0007, or ADR-0007")
	f.StringSliceVar(&refines, "refines", nil, "ADR this one refines (repeatable)")
	f.BoolVar(&bodyStdin, "body-stdin", false, "read the body below the heading from standard input")
	f.BoolVar(&opt.Autocommit, "autocommit", false, "commit what was written on its own, unless dashboard.autocommit is false")
	f.StringArrayVar(&opt.Trailers, "trailer", nil, "trailer line for the commit (repeatable)")
	f.BoolVar(&printBody, "print-body", false, "print the template's sections and create nothing")
	return c
}

func newAdrAcceptCmd(a *app) *cobra.Command {
	var opt adr.Options
	c := &cobra.Command{
		Use:   "accept <number>",
		Short: "Accept a proposed ADR: status accepted, dated today, index row updated",
		Long: `Only a proposed ADR can be accepted. Once accepted an ADR is immutable
except for superseded_by: to change the decision, record a new ADR that
supersedes it.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			repo, err := a.project()
			if err != nil {
				return err
			}
			n, err := adrNumber(args[0])
			if err != nil {
				return err
			}
			opt.Now = a.now()
			res, err := adr.Accept(repo, a.runner, n, opt)
			if rerr, ok := a.refusedADR(err); ok {
				return rerr
			}
			if err != nil {
				return err
			}
			return a.printADR(res, "accepted")
		},
	}
	c.Flags().BoolVar(&opt.Autocommit, "autocommit", false, "commit the change on its own, unless dashboard.autocommit is false")
	c.Flags().StringArrayVar(&opt.Trailers, "trailer", nil, "trailer line for the commit (repeatable)")
	return c
}
