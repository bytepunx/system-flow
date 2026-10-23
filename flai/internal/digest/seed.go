package digest

import (
	"errors"
	"fmt"
	"regexp"
	"slices"
	"strings"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/execx"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// seedKinds are the kinds Seed writes. Decisions are left to a model, which
// can judge which of them a design document already records.
var seedKinds = []string{"answer", "recommendation", "cause"}

// questionLimit is the most characters of a thread's question an entry
// quotes before it is cut on a word and ends in an ellipsis.
const questionLimit = 600

// SeedOptions say what Seed writes and against what it checks it.
type SeedOptions struct {
	Since  time.Time // candidates learned before it are left out; zero is every one
	DryRun bool      // report what would be written and write nothing
	Env    Env       // what each entry is checked against; zero is EnvFor(repo)
}

// SeedEntry is an entry Seed wrote, or would write on a dry run, which gives
// it no ID or path.
type SeedEntry struct {
	ID     string `json:"id"`
	Path   string `json:"path"`
	Kind   string `json:"kind"`
	Source string `json:"source"`
	Ref    string `json:"ref"`
	Title  string `json:"title"`
}

// SeedSkip is a candidate Seed did not write and why: "cited by D-nnnn",
// "no cause", "no fix", or "refused: <rules>".
type SeedSkip struct {
	Kind   string `json:"kind"`
	Source string `json:"source"`
	Ref    string `json:"ref"`
	Title  string `json:"title"`
	Reason string `json:"reason"`
}

// SeedResult is what Seed listed, wrote, and skipped.
type SeedResult struct {
	Listed  map[string]int `json:"listed"` // candidates harvested, by kind
	Written []SeedEntry    `json:"written"`
	Skipped []SeedSkip     `json:"skipped"`
	DryRun  bool           `json:"dry_run"`
}

// SeedTally is the candidates of one kind Seed listed, wrote, and skipped.
type SeedTally struct {
	Listed  int `json:"listed"`
	Written int `json:"written"`
	Skipped int `json:"skipped"`
}

// ByKind counts the candidates listed, written, and skipped, by kind.
func (r SeedResult) ByKind() map[string]SeedTally {
	out := map[string]SeedTally{}
	for k, n := range r.Listed {
		t := out[k]
		t.Listed = n
		out[k] = t
	}
	for _, w := range r.Written {
		t := out[w.Kind]
		t.Written++
		out[w.Kind] = t
	}
	for _, s := range r.Skipped {
		t := out[s.Kind]
		t.Skipped++
		out[s.Kind] = t
	}
	return out
}

// Seed writes an entry, through New, for each answer, overruled
// recommendation, and closed issue's cause and fix that Harvest lists and no
// entry records yet, dated the day it was learned. A thread's candidate is
// recorded by an entry citing the thread whose body holds the candidate's
// time; an issue's by any entry citing it. So a second run writes nothing.
// An entry the check refuses is skipped and named by the rules it broke;
// any other error stops the seed and is returned with what was written so
// far.
func Seed(repo *workitem.Repo, r execx.Runner, opts SeedOptions) (SeedResult, error) {
	res := SeedResult{Listed: map[string]int{}, Written: []SeedEntry{}, Skipped: []SeedSkip{}, DryRun: opts.DryRun}
	for _, k := range seedKinds {
		res.Listed[k] = 0
	}
	env := opts.Env
	if env.Known == nil && env.Exists == nil {
		env = EnvFor(repo)
	}
	candidates, err := Harvest(repo, opts.Since, seedKinds)
	if err != nil {
		return res, err
	}
	entries, err := List(repo)
	if err != nil {
		return res, err
	}
	for _, c := range candidates {
		res.Listed[c.Kind]++
		skip := func(title, reason string) {
			res.Skipped = append(res.Skipped, SeedSkip{Kind: c.Kind, Source: c.Source, Ref: c.Ref, Title: title, Reason: reason})
		}
		if reason := skipReason(c, entries); reason != "" {
			skip(c.Title, reason)
			continue
		}
		body := seedBody(c, repo.Manifest.Owner)
		title := seedTitle(c.Title)
		if opts.DryRun {
			res.Written = append(res.Written, SeedEntry{Kind: c.Kind, Source: c.Source, Ref: c.Ref, Title: title})
			entries = append(entries, Entry{Sources: []string{c.Source}, Body: body})
			continue
		}
		wrote, err := New(repo, r, NewOptions{Title: title, Kind: c.Kind, Sources: []string{c.Source}, Topics: c.Topics, Paths: c.Paths, Body: body, Now: learned(c), Env: env})
		var refused *RefusedError
		if errors.As(err, &refused) {
			if len(Secrets(c.Title)) > 0 {
				title = "" // a skip names the rules it broke, never the text
			}
			skip(title, "refused: "+brokenRules(refused.Problems))
			continue
		}
		if err != nil {
			return res, fmt.Errorf("seeding the digest with the %s from %s %s: %w", c.Kind, c.Source, c.Ref, err)
		}
		res.Written = append(res.Written, SeedEntry{ID: wrote.ID, Path: wrote.Path, Kind: c.Kind, Source: c.Source, Ref: c.Ref, Title: wrote.Title})
		entries = append(entries, wrote.Entry)
	}
	return res, nil
}

// skipReason is why Seed does not write a candidate, or empty when it does:
// a cause with no text or no fix, an issue any entry cites, or a thread's
// candidate an entry citing the thread records by holding its time.
func skipReason(c Candidate, entries []Entry) string {
	thread := c.Kind == "answer" || c.Kind == "recommendation"
	switch {
	case c.Kind == "cause" && strings.TrimSpace(c.Text) == "":
		return "no cause"
	case c.Kind == "cause" && strings.TrimSpace(c.Fix) == "":
		return "no fix"
	}
	source := canonicalSource(c.Source)
	var by []string
	for _, e := range entries {
		if !slices.Contains(clean(e.Sources, canonicalSource), source) || thread && !strings.Contains(e.Body, c.Ref) {
			continue
		}
		id := e.ID
		if id == "" {
			id = "an entry this run writes" // a dry run's, which has no ID
		}
		by = append(by, id)
	}
	if len(by) == 0 {
		return ""
	}
	return "cited by " + strings.Join(by, ", ")
}

// brokenRules names the rules of the problems, each once, in their order.
func brokenRules(problems []Problem) string {
	var rules []string
	for _, p := range problems {
		if !slices.Contains(rules, p.Rule) {
			rules = append(rules, p.Rule)
		}
	}
	return strings.Join(rules, ", ")
}

// seedTitle is a candidate's title as New writes it, on one line, and as an
// entry's heading can end: without the trailing punctuation markdownlint's
// MD026 reports.
func seedTitle(title string) string {
	title = strings.Join(strings.Fields(title), " ")
	if t := strings.TrimRight(title, ".,;:!"); t != "" {
		return t
	}
	return title
}

// seedBody is the body Seed writes below an entry's heading: a sentence
// naming who answered or when the issue closed, which holds the candidate's
// Ref, each text it quotes in a fenced block of its own, and the source.
func seedBody(c Candidate, owner string) string {
	var blocks []string
	if c.Kind == "cause" {
		blocks = []string{c.Source + " was closed on " + c.Date + ".", "The cause:", fenced(c.Text), "The fix:", fenced(c.Fix)}
	} else {
		blocks = []string{"The operator, " + owner + ", answered " + c.Source + " at " + c.Ref + "."}
		if q := firstParagraph(c.Question); q != "" {
			blocks = append(blocks, "What was asked:", fenced(q))
		}
		if c.Kind == "recommendation" {
			blocks = append(blocks, "The recommendation it overruled:", fenced(c.Overruled), "The operator's answer:", fenced(c.Text))
		} else {
			blocks = append(blocks, "The answer:", fenced(c.Text))
		}
	}
	blocks = append(blocks, "Source: "+c.Source+".")
	return strings.Join(blocks, "\n\n") + "\n"
}

var (
	// paragraphBreak is a blank line, or one of spaces and tabs only.
	paragraphBreak = regexp.MustCompile(`\n[ \t]*\n`)
	// backticks is a run of backticks.
	backticks = regexp.MustCompile("`+")
)

// firstParagraph is a question's first paragraph, cut on a word near
// questionLimit.
func firstParagraph(text string) string {
	text = strings.TrimSpace(strings.NewReplacer("\r\n", "\n", "\r", "\n").Replace(text))
	if at := paragraphBreak.FindStringIndex(text); at != nil {
		text = text[:at[0]]
	}
	return clipOnWord(text, questionLimit)
}

// fenced quotes text in a fenced block of language text, so no markdown in
// it is read as markdown: its fence is a backtick longer than the longest
// run of backticks in it and at least three, its tabs are spaces, its lines
// have no trailing spaces, and it has no run of blank lines.
func fenced(text string) string {
	text = strings.NewReplacer("\r\n", "\n", "\r", "\n", "\t", "    ").Replace(text)
	var lines []string
	for _, l := range strings.Split(text, "\n") {
		l = strings.TrimRight(l, " ")
		if l == "" && (len(lines) == 0 || lines[len(lines)-1] == "") {
			continue
		}
		lines = append(lines, l)
	}
	for len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}
	longest := 0
	for _, run := range backticks.FindAllString(text, -1) {
		longest = max(longest, len(run))
	}
	fence := strings.Repeat("`", max(3, longest+1))
	return fence + "text\n" + strings.Join(lines, "\n") + "\n" + fence
}
