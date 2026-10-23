package cmd

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/spf13/cobra"

	"github.com/bytepunx/system-flow/flai/internal/check"
	"github.com/bytepunx/system-flow/flai/internal/digest"
)

func newDigestCmd(a *app) *cobra.Command {
	c := &cobra.Command{
		Use:   "digest",
		Short: "Record and read the digest: the facts past stories learned (ADR-0191)",
		Long: `The digest is one markdown file per fact a story learned, D-nnnn-<slug>.md
in the design folder's digest folder (ADR-0191): a decision, the operator's
answer on a thread, a recommendation the operator confirmed or overruled, or
an issue's cause and fix. flai writes and checks each entry's front matter,
and nothing in the folder is ever deleted.`,
	}
	c.AddCommand(newDigestNewCmd(a), newDigestEditCmd(a), newDigestListCmd(a), newDigestShowCmd(a))
	return c
}

// digestBodyFlags are --body and --body-file, which both give the text below
// an entry's heading.
type digestBodyFlags struct {
	text, file string
}

func (b *digestBodyFlags) add(c *cobra.Command) {
	c.Flags().StringVar(&b.text, "body", "", "the text below the heading: what was asked or decided, the answer and its reason, and links to its sources")
	c.Flags().StringVar(&b.file, "body-file", "", "read the text below the heading from this file, or from standard input when it is -")
	c.MarkFlagsMutuallyExclusive("body", "body-file")
}

// given says whether either body flag was given.
func (b *digestBodyFlags) given(c *cobra.Command) bool {
	return c.Flags().Changed("body") || c.Flags().Changed("body-file")
}

// read is the body the flags give: --body's text, or --body-file's file, a
// relative path read from the working directory.
func (b *digestBodyFlags) read(a *app, c *cobra.Command) (string, error) {
	if !c.Flags().Changed("body-file") {
		return b.text, nil
	}
	if b.file == "-" {
		data, err := io.ReadAll(c.InOrStdin())
		if err != nil {
			return "", fmt.Errorf("reading the body from standard input: %w", err)
		}
		return string(data), nil
	}
	path := b.file
	if !filepath.IsAbs(path) && a.cwd != "" {
		path = filepath.Join(a.cwd, path)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("reading the body from --body-file %s: %w; give a file that exists, or - for standard input", b.file, err)
	}
	return string(data), nil
}

// refusedDigest prints a refused digest write the way a refused ADR is
// printed, each problem as a finding, and returns exit 4; ok is false when
// err is not a refusal.
func (a *app) refusedDigest(err error) (error, bool) {
	var r *digest.RefusedError
	if !errors.As(err, &r) {
		return nil, false
	}
	findings := make([]check.Finding, 0, len(r.Problems))
	for _, p := range r.Problems {
		findings = append(findings, check.Finding{Level: check.Error, Rule: p.Rule, Path: r.Path, Line: p.Line, Message: p.Message})
	}
	if a.jsonOut {
		_ = a.printJSON(map[string]any{"refused": digestRefusal{
			Path:     r.Path,
			Reason:   fmt.Sprintf("the digest entry has %d problem(s); nothing was written", len(r.Problems)),
			Findings: findings,
		}})
	} else {
		for _, f := range findings {
			fmt.Fprintf(a.out, "%s:%d: %s: %s: %s\n", f.Path, f.Line, f.Level, f.Rule, f.Message)
		}
	}
	return &exitError{code: exitDocRefused, msg: err.Error()}, true
}

// digestRefusal is a refused digest write as --json prints it, in the shape
// of a refused document save.
type digestRefusal struct {
	Path     string          `json:"path"`
	Reason   string          `json:"reason"`
	Findings []check.Finding `json:"findings"`
}

// printDigestResult prints the entry new or edit wrote.
func (a *app) printDigestResult(res digest.Result, verb string) error {
	if a.jsonOut {
		return a.printJSON(res)
	}
	fmt.Fprintf(a.out, "%s %s (%s, %s)\n  %s\n  %s\n", res.ID, res.Title, res.Entry.Kind, res.Entry.Status, res.Path, verb)
	return nil
}

func newDigestNewCmd(a *app) *cobra.Command {
	var kind string
	var sources, topics, paths []string
	var body digestBodyFlags
	c := &cobra.Command{
		Use:   "new \"<fact>\" --kind <kind> --source <id>...",
		Short: "Record a digest entry: its ID, file, and front matter are supplied",
		Long: `Record one fact a story learned as a new digest entry, in the checkout it
runs in. The ID is one more than the highest D-nnnn in this checkout, on
main, in any worktree, and on any story branch, so parallel writers do not
take the same one; gaps are not filled. The file is D-nnnn-<slug>.md in the
digest folder, dated today, with status active.

--kind is what the fact is: ` + strings.Join(digest.Kinds, ", ") + `. --source
names the work items (E-, S-, T-), threads (TH-), and issues (I-) it comes
from, at least one, each in the project or its archive. --source, --topic,
and --path are repeatable, or comma separated. The body is a few sentences
below the heading; it says where a secret lives, never its value.

The check runs on what was written: when it reports anything, such as a
source that names nothing or text with the shape of a key or token, the file
is removed, each problem is printed, and the exit code is 4.`,
		Example: `  flai digest new "Story worktrees install with the project's pnpm store" --kind cause --source I-0177 --source S-0481 --topic flaiover --path scripts/flaiover-install.sh
  flai digest new "The digest lives under design/digest" --kind answer --source TH-0476 --body-file answer.md --json`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if kind == "" {
				return fmt.Errorf("--kind is required: give what the fact is, one of %s", strings.Join(digest.Kinds, ", "))
			}
			if !slices.Contains(digest.Kinds, kind) {
				return fmt.Errorf("--kind %q is none of %s; give one of them", kind, strings.Join(digest.Kinds, ", "))
			}
			if len(sources) == 0 {
				return errors.New("--source is required: give at least one work item (E-, S-, T-), thread (TH-), or issue (I-) the fact comes from, such as --source S-0042")
			}
			text, err := body.read(a, cmd)
			if err != nil {
				return err
			}
			repo, err := a.project()
			if err != nil {
				return err
			}
			res, err := digest.New(repo, a.runner, digest.NewOptions{
				Title: args[0], Kind: kind, Sources: sources, Topics: topics, Paths: paths,
				Body: text, Now: a.now(), Env: digest.EnvFor(repo),
			})
			if rerr, ok := a.refusedDigest(err); ok {
				return rerr
			}
			if err != nil {
				return err
			}
			return a.printDigestResult(res, "recorded")
		},
	}
	f := c.Flags()
	f.StringVar(&kind, "kind", "", "what the fact is: "+strings.Join(digest.Kinds, ", "))
	f.StringSliceVar(&sources, "source", nil, "work item (E-, S-, T-), thread (TH-), or issue (I-) the fact comes from (repeatable, or comma separated; at least one)")
	f.StringSliceVar(&topics, "topic", nil, "topic the fact bears on, as a story's are (repeatable, or comma separated)")
	f.StringSliceVar(&paths, "path", nil, "repository path the fact bears on, from the project root (repeatable, or comma separated)")
	body.add(c)
	return c
}

func newDigestEditCmd(a *app) *cobra.Command {
	var topics, paths, sources, askedAgain []string
	var body digestBodyFlags
	c := &cobra.Command{
		Use:   "edit <D-nnnn>",
		Short: "Change a digest entry's body, topics, paths, sources, or asked_again",
		Long: `Change a digest entry in the checkout it runs in. Each flag given replaces
that part of the entry, and leaves the rest of the file as it was: --body or
--body-file the text below the heading, and --topic, --path, --source, and
--asked-again their lists, each repeatable or comma separated. Give a list
flag an empty value, such as --topic=, to empty it. The ID may be written in
any padding: D-7 is D-0007.

The check runs on what was written: when it reports a problem the change
introduced, the file is put back as it was, each problem is printed, and
the exit code is 4.`,
		Example: `  flai digest edit D-0007 --topic cli --topic go
  flai digest edit D-0007 --asked-again TH-0512 --asked-again TH-0530 --json
  flai digest edit D-0007 --body-file - < body.md`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			fl := cmd.Flags()
			opt := digest.EditOptions{}
			for _, l := range []struct {
				name string
				in   []string
				to   **[]string
			}{{"topic", topics, &opt.Topics}, {"path", paths, &opt.Paths}, {"source", sources, &opt.Sources}, {"asked-again", askedAgain, &opt.AskedAgain}} {
				if fl.Changed(l.name) {
					list := l.in
					*l.to = &list
				}
			}
			if body.given(cmd) {
				text, err := body.read(a, cmd)
				if err != nil {
					return err
				}
				opt.Body = &text
			}
			if opt.Body == nil && opt.Topics == nil && opt.Paths == nil && opt.Sources == nil && opt.AskedAgain == nil {
				return errors.New("nothing to change: give --body or --body-file, --topic, --path, --source, or --asked-again")
			}
			repo, err := a.project()
			if err != nil {
				return err
			}
			opt.Env = digest.EnvFor(repo)
			res, err := digest.Edit(repo, args[0], opt)
			if rerr, ok := a.refusedDigest(err); ok {
				return rerr
			}
			if err != nil {
				return err
			}
			return a.printDigestResult(res, "edited")
		},
	}
	f := c.Flags()
	f.StringSliceVar(&topics, "topic", nil, "topic the fact bears on; replaces the entry's topics (repeatable, or comma separated)")
	f.StringSliceVar(&paths, "path", nil, "repository path the fact bears on; replaces the entry's paths (repeatable, or comma separated)")
	f.StringSliceVar(&sources, "source", nil, "work item, thread, or issue the fact comes from; replaces the entry's sources (repeatable, or comma separated)")
	f.StringSliceVar(&askedAgain, "asked-again", nil, "thread that asked what the entry already answered; replaces the entry's asked_again (repeatable, or comma separated)")
	body.add(c)
	return c
}

func newDigestListCmd(a *app) *cobra.Command {
	var status, topic string
	c := &cobra.Command{
		Use:   "list",
		Short: "List the digest's entries with their sources, topics, paths, date, and status",
		Long: `List the entries in this checkout's digest folder, in ID order: each one's
ID, status, kind, date, and title, and below it its sources, topics, and
paths. --status keeps the entries in that status (` + strings.Join(digest.Statuses, ", ") + `);
--topic keeps those that bear on the topic. --json gives an array of the
entries, bodies included, [] when none.`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if status != "" && !slices.Contains(digest.Statuses, status) {
				return fmt.Errorf("--status %q is none of %s; give one of them", status, strings.Join(digest.Statuses, ", "))
			}
			repo, err := a.project()
			if err != nil {
				return err
			}
			all, err := digest.List(repo)
			if err != nil {
				return err
			}
			list := []digest.Entry{}
			for _, e := range all {
				if (status == "" || e.Status == status) && (topic == "" || slices.Contains(e.Topics, topic)) {
					list = append(list, e)
				}
			}
			if a.jsonOut {
				return a.printJSON(list)
			}
			if len(list) == 0 {
				fmt.Fprintln(a.out, "no digest entries")
				return nil
			}
			for _, e := range list {
				fmt.Fprintf(a.out, "%-6s %-9s %-14s %-10s %s\n", e.ID, e.Status, e.Kind, e.Date, e.Title)
				fmt.Fprintf(a.out, "       sources: %s  topics: %s  paths: %s\n", digestList(e.Sources), digestList(e.Topics), digestList(e.Paths))
			}
			return nil
		},
	}
	c.Flags().StringVar(&status, "status", "", "only the entries in this status: "+strings.Join(digest.Statuses, ", "))
	c.Flags().StringVar(&topic, "topic", "", "only the entries that bear on this topic")
	return c
}

func newDigestShowCmd(a *app) *cobra.Command {
	return &cobra.Command{
		Use:   "show <D-nnnn>",
		Short: "Print a digest entry: its front matter and its body",
		Long: `Print one digest entry: its title, kind, status, date, sources, topics,
paths, absorbed_by when it is absorbed, asked_again, cited, its file, and
the body below its heading. The ID may be written in any padding: D-7 is
D-0007. --json gives the entry.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			repo, err := a.project()
			if err != nil {
				return err
			}
			e, err := digest.Get(repo, args[0])
			if err != nil {
				return err
			}
			if a.jsonOut {
				return a.printJSON(e)
			}
			fmt.Fprintf(a.out, "%s %s\n", e.ID, e.Title)
			fields := [][2]string{
				{"kind", e.Kind}, {"status", e.Status}, {"date", e.Date},
				{"sources", digestList(e.Sources)}, {"topics", digestList(e.Topics)}, {"paths", digestList(e.Paths)},
			}
			if e.AbsorbedBy != "" {
				fields = append(fields, [2]string{"absorbed_by", e.AbsorbedBy})
			}
			fields = append(fields, [2]string{"asked_again", digestList(e.AskedAgain)}, [2]string{"cited", digestList(e.Cited)}, [2]string{"path", e.Path})
			for _, f := range fields {
				fmt.Fprintf(a.out, "  %-12s %s\n", f[0]+":", f[1])
			}
			if e.Body != "" {
				fmt.Fprintf(a.out, "\n%s\n", e.Body)
			}
			return nil
		},
	}
}

// digestList is a list as the text output shows it: comma separated, or -
// when empty.
func digestList(l []string) string {
	if len(l) == 0 {
		return "-"
	}
	return strings.Join(l, ", ")
}
