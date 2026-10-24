package cmd

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/bytepunx/system-flow/flai/internal/check"
	ctxpack "github.com/bytepunx/system-flow/flai/internal/context"
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
	c.AddCommand(newDigestNewCmd(a), newDigestEditCmd(a), newDigestListCmd(a), newDigestShowCmd(a), newDigestHarvestCmd(a), newDigestSeedCmd(a),
		newDigestAbsorbCmd(a), newDigestCitedCmd(a), newDigestCompactCmd(a))
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
paths, and the document that absorbed it when it is absorbed. --status keeps the entries in that status (` + strings.Join(digest.Statuses, ", ") + `);
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
				absorbed := ""
				if e.AbsorbedBy != "" {
					absorbed = "  absorbed_by: " + e.AbsorbedBy
				}
				fmt.Fprintf(a.out, "       sources: %s  topics: %s  paths: %s%s\n", digestList(e.Sources), digestList(e.Topics), digestList(e.Paths), absorbed)
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

func newDigestHarvestCmd(a *app) *cobra.Command {
	var since string
	var kinds []string
	c := &cobra.Command{
		Use:   "harvest",
		Short: "List the facts the digest could record: answers, recommendations, decisions, and causes",
		Long: `List every candidate for the digest in this checkout, learned on or after
--since, or every one when it is not given: the operator's answers and the
recommendations the operator overruled on every thread, resolved or not, the
decisions in done stories' narratives, and the causes and fixes of closed
issues. It is deterministic and writes nothing, so the curator reads the same
list every run.

Each line gives the candidate's kind, its date, its source with where in it
(the thread entry's time, or the narrative's path), its title, and the
entries that already cite the source; the last line counts them by kind.
--since is a day, YYYY-MM-DD in UTC, or an RFC 3339 time. --kind keeps the
candidates of that kind (` + strings.Join(digest.Kinds, ", ") + `), repeatable or
comma separated. --json gives an array of the candidates with every field,
their topics, paths, and text included, [] when none.`,
		Example: `  flai digest harvest --since 2026-10-01
  flai digest harvest --kind answer,recommendation --json`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			from, err := digestSince(since)
			if err != nil {
				return err
			}
			for _, k := range kinds {
				if !slices.Contains(digest.Kinds, strings.TrimSpace(k)) {
					return fmt.Errorf("--kind %q is none of %s; give one of them", k, strings.Join(digest.Kinds, ", "))
				}
			}
			repo, err := a.project()
			if err != nil {
				return err
			}
			list, err := digest.Harvest(repo, from, kinds)
			if err != nil {
				return err
			}
			if a.jsonOut {
				return a.printJSON(list)
			}
			if len(list) == 0 {
				fmt.Fprintln(a.out, "no candidates")
				return nil
			}
			count := map[string]int{}
			uncited := 0
			for _, cand := range list {
				count[cand.Kind]++
				cited := ""
				if len(cand.CitedBy) > 0 {
					cited = "  (cited by " + strings.Join(cand.CitedBy, ", ") + ")"
				} else {
					uncited++
				}
				title, _, _ := strings.Cut(cand.Title, "\n")
				fmt.Fprintf(a.out, "%-14s %-10s %s: %s%s\n", cand.Kind, cand.Date, digestWhere(cand.Source, cand.Ref), title, cited)
			}
			var counts []string
			for _, k := range digest.Kinds {
				if count[k] > 0 {
					counts = append(counts, fmt.Sprintf("%d %s", count[k], k))
				}
			}
			fmt.Fprintf(a.out, "%d candidates, %d cited by no entry: %s\n", len(list), uncited, strings.Join(counts, ", "))
			return nil
		},
	}
	c.Flags().StringVar(&since, "since", "", "only the candidates learned on or after this day, YYYY-MM-DD in UTC, or this RFC 3339 time")
	c.Flags().StringSliceVar(&kinds, "kind", nil, "only the candidates of this kind: "+strings.Join(digest.Kinds, ", ")+" (repeatable, or comma separated)")
	return c
}

func newDigestSeedCmd(a *app) *cobra.Command {
	var since string
	var dryRun bool
	c := &cobra.Command{
		Use:   "seed",
		Short: "Write a digest entry for each operator answer and closed issue's cause no entry records yet",
		Long: `Seed the digest in the checkout it runs in from what harvest lists, learned
on or after --since, or all of it when it is not given. It writes an answer
or recommendation entry for each answer the operator gave on a thread, and a
cause entry for each closed issue with a cause and a fix, each dated the day
it was learned and quoting its source. It never seeds decisions: a model
judges which of them a design document already records.

A source an entry already records is skipped: an issue any entry cites, or a
thread's answer when an entry citing the thread holds the answer's time. So
a second run writes nothing. Each entry is checked as flai digest new checks
it; one the check refuses, such as an answer with the shape of a key or
token, is skipped and named by the rules it broke, never its text, and the
rest are written.

Each entry written is printed with its ID, kind, source, and title, and each
skipped one with its reason, then the candidates listed, written, and
skipped by kind. --dry-run prints what would be written and writes nothing.
--json gives the result: listed by kind, written, skipped, and dry_run. A
skipped entry is not an error: the exit code is 0 unless the seed fails.`,
		Example: `  flai digest seed --dry-run
  flai digest seed --since 2026-10-01 --json`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			from, err := digestSince(since)
			if err != nil {
				return err
			}
			repo, err := a.project()
			if err != nil {
				return err
			}
			res, err := digest.Seed(repo, a.runner, digest.SeedOptions{Since: from, DryRun: dryRun})
			if err != nil {
				if len(res.Written) > 0 && !a.jsonOut {
					a.printSeedResult(res)
				}
				return err
			}
			if a.jsonOut {
				return a.printJSON(res)
			}
			a.printSeedResult(res)
			return nil
		},
	}
	c.Flags().StringVar(&since, "since", "", "only the candidates learned on or after this day, YYYY-MM-DD in UTC, or this RFC 3339 time")
	c.Flags().BoolVar(&dryRun, "dry-run", false, "print what would be written and write nothing")
	return c
}

func newDigestAbsorbCmd(a *app) *cobra.Command {
	var by, heading string
	c := &cobra.Command{
		Use:   "absorb <D-nnnn> --by <path> [--heading <heading>]",
		Short: "Mark a digest entry absorbed by the document that now records its fact",
		Long: `Mark a digest entry absorbed, in the checkout it runs in (ADR-0199): its
status becomes absorbed, and its absorbed_by the path, from the project root,
of the design document, convention, or ADR that now records the fact, or
<path>#<anchor> with --heading, which must be one of the document's
headings. A heading is its text, a heading path such as "Rules › Go" (>
also joins), or its anchor, as flai doc show --heading takes it. The body
is kept. A story's pack never loads an absorbed entry. An entry absorbed
again names the newer document, and a compacted entry may be absorbed. The
ID may be written in any padding: D-7 is D-0007.

A path that is not from the project root or names no file, and a heading
the document does not have, are refused: nothing is written, the reason is
printed, and the exit code is 4. --json gives the entry, as flai digest
edit does.`,
		Example: `  flai digest absorb D-0007 --by design/system/flai-cli.md
  flai digest absorb D-0007 --by design/conventions/git.md --heading "Rules" --json`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if strings.TrimSpace(by) == "" {
				return errors.New("--by is required: give the path, from the project root, of the document that now records the fact, such as --by design/system/overview.md")
			}
			repo, err := a.project()
			if err != nil {
				return err
			}
			res, err := digest.Absorb(repo, args[0], digest.AbsorbOptions{By: by, Heading: heading, Anchor: digestAnchor, Env: digest.EnvFor(repo)})
			if rerr, ok := a.refusedDigest(err); ok {
				return rerr
			}
			if err != nil {
				return err
			}
			return a.printDigestResult(res, "absorbed by "+res.Entry.AbsorbedBy)
		},
	}
	c.Flags().StringVar(&by, "by", "", "the path, from the project root, of the document that now records the fact (required)")
	c.Flags().StringVar(&heading, "heading", "", "the heading of the document that records it: its text, a heading path, or its anchor")
	return c
}

// digestAnchor is the anchor of the section a heading names in a document,
// found as flai doc show --heading finds it.
func digestAnchor(path, content, heading string) (string, error) {
	part, err := ctxpack.Section(path, content, heading)
	if err != nil {
		return "", err
	}
	doc, err := ctxpack.ParseDoc(path, "", content)
	if err != nil {
		return "", err
	}
	for _, sec := range doc.Sections {
		if sec.Line == part.Line {
			return ctxpack.Slug(sec.Heading), nil
		}
	}
	return "", fmt.Errorf("the heading %q of %s was found at line %d but no section starts there; this is a flai bug: report it", heading, path, part.Line)
}

func newDigestCitedCmd(a *app) *cobra.Command {
	return &cobra.Command{
		Use:   "cited",
		Short: "List each digest entry's citers and latest citation, and write its cited list",
		Long: `List every entry in this checkout's digest with the work items, stories,
and threads that cite it and its latest citation (ADR-0199), and write each
entry's cited front matter as its citers, in ID order, where that changed.

A citation is the entry's ID, D-nnnn in any zero padding, written in a work
item's file (an epic, story, or task, on the board or in the archive), in a
narrative (in the agents folder or the archive), which counts as its story,
or in a thread. Its date is the citing file's updated time. Digest entries
citing one another do not count. A citing file whose updated time cannot be
read is an error naming it.

Each line gives the entry's ID, status, latest citation, or - when none,
and its citers, with "written" where cited changed; the last line counts
the entries and those written. Each write is checked as flai digest edit's
is. --json gives an array of the entries, each with its id, title, status,
date, path, cited, latest (empty when none), and changed, [] when none.`,
		Example: `  flai digest cited
  flai digest cited --json`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			repo, err := a.project()
			if err != nil {
				return err
			}
			list, err := digest.Cited(repo, digest.EnvFor(repo))
			if rerr, ok := a.refusedDigest(err); ok {
				return rerr
			}
			if err != nil {
				return err
			}
			if a.jsonOut {
				return a.printJSON(list)
			}
			if len(list) == 0 {
				fmt.Fprintln(a.out, "no digest entries")
				return nil
			}
			written := 0
			for _, c := range list {
				mark := ""
				if c.Changed {
					mark = "  (written)"
					written++
				}
				fmt.Fprintf(a.out, "%-6s %-9s %-20s %s%s\n", c.ID, c.Status, digestLatest(c.Latest), digestList(c.Cited), mark)
			}
			fmt.Fprintf(a.out, "%d entries, cited written on %d\n", len(list), written)
			return nil
		},
	}
}

func newDigestCompactCmd(a *app) *cobra.Command {
	var dryRun bool
	c := &cobra.Command{
		Use:   "compact",
		Short: "Compact the active digest entries no one cited within digest.compact_after to their title",
		Long: `Compact the digest in the checkout it runs in (ADR-0199). It runs flai
digest cited first, then turns each active entry whose latest citation, or
its own date when that is later or it has none, is older than
digest.compact_after (720h, 30 days, when it is not set) into compacted: its
status becomes compacted and its body is emptied below its # D-nnnn <title>
heading, so a story's pack loads only its title. Absorbed and compacted
entries are left as they are, so a second run compacts nothing more.
Nothing is deleted: git keeps every earlier text. With digest.compact_after
0, compaction is off and nothing is done; a digest.compact_after that is
not a duration is refused, naming it.

Each entry compacted is printed with its ID, title, and its latest citation
or its date, then the count, the time an entry had to be cited or dated
since to be kept, and the entries whose cited was written. --dry-run prints
what would be compacted and written, and writes nothing. --json gives the
result: window, off, dry_run, before, cited (as flai digest cited --json
gives it), and compacted, each with its id, title, path, latest (empty when
none), and date.`,
		Example: `  flai digest compact --dry-run
  flai digest compact --json`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			repo, err := a.project()
			if err != nil {
				return err
			}
			window, err := repo.Manifest.Digest.CompactAfterDuration()
			if err != nil {
				return err
			}
			res, err := digest.Compact(repo, digest.CompactOptions{Window: window, Now: a.now(), DryRun: dryRun, Env: digest.EnvFor(repo)})
			if rerr, ok := a.refusedDigest(err); ok {
				return rerr
			}
			if err != nil {
				return err
			}
			if a.jsonOut {
				return a.printJSON(res)
			}
			a.printCompactResult(res)
			return nil
		},
	}
	c.Flags().BoolVar(&dryRun, "dry-run", false, "print what would be compacted and written, and write nothing")
	return c
}

// printCompactResult prints each entry compact compacted, or would, and the
// counts.
func (a *app) printCompactResult(res digest.CompactResult) {
	if res.Off {
		fmt.Fprintln(a.out, "compaction is off: digest.compact_after is 0; nothing compacted")
		return
	}
	verb, count, written := "compacted", "compacted", "written"
	if res.DryRun {
		verb, count, written = "would compact", "to compact", "would be written"
	}
	for _, c := range res.Compacted {
		judged := "no citation, dated " + c.Date
		if c.Latest != "" {
			judged = "last cited " + c.Latest
		}
		fmt.Fprintf(a.out, "%s %s %s (%s)\n", verb, c.ID, c.Title, judged)
	}
	changed := 0
	for _, c := range res.Cited {
		if c.Changed {
			changed++
		}
	}
	fmt.Fprintf(a.out, "%d %s: not cited or dated since %s (digest.compact_after %s); cited %s on %d\n", len(res.Compacted), count, res.Before, res.Window, written, changed)
}

// digestLatest is a latest citation as the text output shows it, or - when
// there is none.
func digestLatest(s string) string {
	if s == "" {
		return "-"
	}
	return s
}

// printSeedResult prints each entry seed wrote, or would write, each it
// skipped with its reason, and the counts by kind.
func (a *app) printSeedResult(res digest.SeedResult) {
	written := "written"
	if res.DryRun {
		written = "would write"
	}
	for _, w := range res.Written {
		lead := w.ID
		if res.DryRun {
			lead = "would write"
		}
		fmt.Fprintf(a.out, "%s %s %s %s\n", lead, w.Kind, w.Source, w.Title)
	}
	for _, s := range res.Skipped {
		title := ""
		if s.Title != "" {
			title = " " + s.Title
		}
		fmt.Fprintf(a.out, "skipped %s %s%s: %s\n", s.Kind, digestWhere(s.Source, s.Ref), title, s.Reason)
	}
	tally := res.ByKind()
	for _, k := range digest.Kinds {
		t, ok := tally[k]
		if !ok {
			continue
		}
		fmt.Fprintf(a.out, "%-15s listed %d, %s %d, skipped %d\n", k+":", t.Listed, written, t.Written, t.Skipped)
	}
}

// digestSince is the time --since gives, a day in UTC or an RFC 3339 time,
// or the zero time, every candidate, when it is empty.
func digestSince(s string) (time.Time, error) {
	if s == "" {
		return time.Time{}, nil
	}
	if t, err := time.Parse(time.DateOnly, s); err == nil {
		return t, nil
	}
	if t, err := time.Parse(time.RFC3339, s); err == nil {
		return t, nil
	}
	return time.Time{}, fmt.Errorf("--since %q is not a date: give a day, YYYY-MM-DD in UTC, such as 2026-10-01, or an RFC 3339 time, such as 2026-10-01T12:00:00Z", s)
}

// digestWhere is a candidate's source, followed by where in it when it says.
func digestWhere(source, ref string) string {
	if ref == "" {
		return source
	}
	return source + " " + ref
}

// digestList is a list as the text output shows it: comma separated, or -
// when empty.
func digestList(l []string) string {
	if len(l) == 0 {
		return "-"
	}
	return strings.Join(l, ", ")
}
