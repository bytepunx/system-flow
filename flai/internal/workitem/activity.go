package workitem

import (
	"errors"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/goccy/go-yaml"

	"github.com/bytepunx/system-flow/flai/internal/atomicfile"
)

// The strategic agents (ADR-0075) that each keep one activity document under
// wip/agents (ADR-0079).
const (
	ActivityPlanner      = "planner"
	ActivityOrchestrator = "orchestrator"
	ActivityAnalyzer     = "analyzer"
)

// ActivityKinds are the strategic agents' kinds, in the order flai lists them.
var ActivityKinds = []string{ActivityPlanner, ActivityOrchestrator, ActivityAnalyzer}

// IsActivityKind reports whether kind names a strategic agent.
func IsActivityKind(kind string) bool { return contains(ActivityKinds, kind) }

// IsActivityDocument reports whether a file under wip/agents, by its name or
// path, is a strategic agent's activity document rather than a narrative.
func IsActivityDocument(name string) bool {
	base := filepath.Base(name)
	return strings.HasSuffix(base, ".md") && IsActivityKind(strings.TrimSuffix(base, ".md"))
}

// Activity is a strategic agent's activity document, wip/agents/<kind>.md:
// its totals in front matter and its log entries, oldest first.
type Activity struct {
	Kind           string  `yaml:"kind" json:"kind"`
	AccruedCost    float64 `yaml:"accrued_cost" json:"accrued_cost"`
	AccruedSeconds int64   `yaml:"accrued_seconds" json:"accrued_seconds"`
	TasksCompleted int     `yaml:"tasks_completed" json:"tasks_completed"`
	// LastRun is when the newest activity ended, in TimeFormat; empty before
	// the first.
	LastRun string `yaml:"last_run" json:"last_run"`

	Path    string          `yaml:"-" json:"path"`
	Entries []ActivityEntry `yaml:"-" json:"entries"`
	// Refusals are the calls flai guard refused the agent, oldest first,
	// under ## Refusals (S-0218).
	Refusals []ActivityRefusal `yaml:"-" json:"refusals,omitempty"`
}

// ActivityEntry is one activity in the log.
type ActivityEntry struct {
	// At is when the activity ended.
	At time.Time `json:"at"`
	// Summary is one line saying what the agent did.
	Summary string `json:"summary"`
	// Trigger is one line saying what started the activity, written for a
	// planner run flai serve started (ADR-0084): asked, or the replanner's
	// triggers, separated by semicolons. It is empty on older entries and on
	// those an agent logs itself.
	Trigger string `json:"trigger,omitempty"`
	// Items are the IDs of the items the activity touched.
	Items   []string `json:"items"`
	Seconds int64    `json:"seconds"`
	// Cost is in US dollars; Estimated when it was apportioned or priced
	// rather than reported.
	Cost      float64 `json:"cost"`
	Estimated bool    `json:"estimated"`
}

// ActivityRefusal is a call flai guard refused a strategic agent (S-0218).
// It is no activity: it has no seconds or cost, and the totals do not count
// it.
type ActivityRefusal struct {
	// At is when the call was refused.
	At time.Time `json:"at"`
	// Call is the call refused, in one line.
	Call string `json:"call"`
	// Needs is the permission that would allow the call, such as
	// orchestration.permissions.publish; empty when none would.
	Needs string `json:"needs,omitempty"`
}

// ActivityPath is wip/agents/<kind>.md.
func (r *Repo) ActivityPath(kind string) string {
	return filepath.Join(r.AgentsDir(), kind+".md")
}

// Activity loads a strategic agent's activity document. When the document
// does not exist yet it returns an empty one, with no entries and zero
// totals, for that kind and path.
func (r *Repo) Activity(kind string) (*Activity, error) {
	if !IsActivityKind(kind) {
		return nil, unknownActivityKind(kind)
	}
	path := r.ActivityPath(kind)
	a, err := ReadActivity(path)
	if errors.Is(err, os.ErrNotExist) {
		return &Activity{Kind: kind, Path: path}, nil
	}
	return a, err
}

// Activities loads the activity documents that exist, in ActivityKinds order.
func (r *Repo) Activities() ([]*Activity, error) {
	var out []*Activity
	for _, kind := range ActivityKinds {
		a, err := ReadActivity(r.ActivityPath(kind))
		if errors.Is(err, os.ErrNotExist) {
			continue
		}
		if err != nil {
			return nil, err
		}
		out = append(out, a)
	}
	return out, nil
}

// ReadActivity loads and validates an activity document: its front matter is
// parsed strictly, its kind must be its file's name, and its log is parsed
// back into entries.
func ReadActivity(path string) (*Activity, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	a, err := parseActivity(string(data))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", path, err)
	}
	if want := strings.TrimSuffix(filepath.Base(path), ".md"); a.Kind != want {
		return nil, fmt.Errorf("%s: kind %q must be the file's name, %q", path, a.Kind, want)
	}
	a.Path = path
	return a, nil
}

func parseActivity(doc string) (*Activity, error) {
	fm, body, err := SplitFrontMatter(doc)
	if err != nil {
		return nil, err
	}
	var a Activity
	if err := yaml.UnmarshalWithOptions([]byte(fm), &a, yaml.Strict()); err != nil {
		return nil, fmt.Errorf("front matter must hold only kind, accrued_cost, accrued_seconds, tasks_completed, and last_run: %w", err)
	}
	if !IsActivityKind(a.Kind) {
		return nil, unknownActivityKind(a.Kind)
	}
	if a.LastRun != "" {
		if _, err := time.Parse(TimeFormat, a.LastRun); err != nil {
			return nil, fmt.Errorf("last_run %q must be a time like %s", a.LastRun, TimeFormat)
		}
	}
	if a.Entries, a.Refusals, err = parseActivityLog(body); err != nil {
		return nil, err
	}
	return &a, nil
}

var activityHeading = regexp.MustCompile(`^### (\d{4}-\d\d-\d\dT\d\d:\d\d:\d\dZ)(?: \(\d+\))?$`)

const (
	logSection      = "## Log"
	refusalsSection = "## Refusals"
	summaryField    = "- Summary: "
	triggerField    = "- Trigger: "
	itemsField      = "- Items: "
	secondsField    = "- Seconds: "
	costField       = "- Cost: "
	callField       = "- Call: "
	needsField      = "- Needs: "
	noItems         = "none"
	noPermission    = "none"
	usd             = " USD"
	estimated       = ", estimated"
)

// parseActivityLog reads the entries under ## Log and the refusals under
// ## Refusals, which follows it when there are any, as Marshal writes them.
func parseActivityLog(body string) ([]ActivityEntry, []ActivityRefusal, error) {
	lines := strings.Split(body, "\n")
	start, refusals := slices.Index(lines, logSection), -1
	if start < 0 {
		return nil, nil, errors.New("no ## Log section; flai writes it with the first activity")
	}
	end := len(lines)
	if i := slices.Index(lines[start:], refusalsSection); i >= 0 {
		refusals, end = start+i, start+i
	}
	var entries []ActivityEntry
	err := parseActivityBlocks(logSection, "log entry", lines[start+1:end], []string{summaryField, itemsField, secondsField, costField}, func(at time.Time) func(string) (string, error) {
		entries = append(entries, ActivityEntry{At: at})
		e := &entries[len(entries)-1]
		return func(l string) (string, error) { return parseActivityField(e, l) }
	})
	if err != nil || refusals < 0 {
		return entries, nil, err
	}
	var out []ActivityRefusal
	err = parseActivityBlocks(refusalsSection, "refusal", lines[refusals+1:], []string{callField, needsField}, func(at time.Time) func(string) (string, error) {
		out = append(out, ActivityRefusal{At: at})
		f := &out[len(out)-1]
		return func(l string) (string, error) { return parseRefusalField(f, l) }
	})
	return entries, out, err
}

// parseActivityBlocks reads the blocks of a section, each a heading with the
// time and its field lines: begin starts a block at its time and returns
// what sets a field line on it and returns the line's prefix. Each block must
// have the lines required, and has no line twice and no other line.
func parseActivityBlocks(section, what string, lines, required []string, begin func(time.Time) func(string) (string, error)) error {
	var set func(string) (string, error)
	var at string
	var seen map[string]bool
	finish := func() error {
		if set == nil {
			return nil
		}
		for _, f := range required {
			if !seen[f] {
				return fmt.Errorf("%s %s has no %q line", what, at, strings.TrimSpace(f))
			}
		}
		return nil
	}
	for _, l := range lines {
		if strings.TrimSpace(l) == "" {
			continue
		}
		if m := activityHeading.FindStringSubmatch(l); m != nil {
			if err := finish(); err != nil {
				return err
			}
			t, err := time.Parse(TimeFormat, m[1])
			if err != nil {
				return fmt.Errorf("%s heading %q: %w", what, l, err)
			}
			at, seen, set = m[1], map[string]bool{}, begin(t)
			continue
		}
		if set == nil {
			return fmt.Errorf("line %q under %s is not an entry heading like ### %s", l, section, TimeFormat)
		}
		field, err := set(l)
		if err != nil {
			return fmt.Errorf("%s %s: %w", what, at, err)
		}
		if seen[field] {
			return fmt.Errorf("%s %s has two %q lines", what, at, strings.TrimSpace(field))
		}
		seen[field] = true
	}
	return finish()
}

// parseRefusalField sets the field line l names on f and returns its prefix.
func parseRefusalField(f *ActivityRefusal, l string) (string, error) {
	switch {
	case strings.HasPrefix(l, callField):
		v := strings.TrimPrefix(l, callField)
		call, opened := strings.CutPrefix(v, "`")
		call, closed := strings.CutSuffix(call, "`")
		if !opened || !closed || call == "" || strings.Contains(call, "`") {
			return "", fmt.Errorf("call %q must be one code span, like `flai accept S-0001`", v)
		}
		f.Call = call
		return callField, nil
	case strings.HasPrefix(l, needsField):
		if v := strings.TrimPrefix(l, needsField); v != noPermission {
			f.Needs = v
		}
		return needsField, nil
	}
	return "", fmt.Errorf("line %q is not one of Call or Needs", l)
}

// parseActivityField sets the field line l names on e and returns its prefix.
func parseActivityField(e *ActivityEntry, l string) (string, error) {
	switch {
	case strings.HasPrefix(l, summaryField):
		e.Summary = strings.TrimPrefix(l, summaryField)
		return summaryField, nil
	case strings.HasPrefix(l, triggerField):
		e.Trigger = strings.TrimPrefix(l, triggerField)
		return triggerField, nil
	case strings.HasPrefix(l, itemsField):
		if v := strings.TrimPrefix(l, itemsField); v != noItems {
			e.Items = strings.Split(v, ", ")
		}
		return itemsField, nil
	case strings.HasPrefix(l, secondsField):
		n, err := strconv.ParseInt(strings.TrimPrefix(l, secondsField), 10, 64)
		if err != nil {
			return "", fmt.Errorf("seconds %q must be a whole number", strings.TrimPrefix(l, secondsField))
		}
		e.Seconds = n
		return secondsField, nil
	case strings.HasPrefix(l, costField):
		v := strings.TrimPrefix(l, costField)
		v, e.Estimated = strings.CutSuffix(v, estimated)
		amount, ok := strings.CutSuffix(v, usd)
		f, err := strconv.ParseFloat(amount, 64)
		if !ok || err != nil {
			return "", fmt.Errorf("cost %q must be dollars like 0.4213 USD, optionally followed by %q", strings.TrimPrefix(l, costField), estimated)
		}
		e.Cost = f
		return costField, nil
	}
	return "", fmt.Errorf("line %q is not one of Summary, Trigger, Items, Seconds, or Cost", l)
}

// Marshal renders the activity document: its log, then its refusals under
// ## Refusals when it has any. Entries whose end times share a second get an
// ordinal after the time, and refusals likewise, so no two headings of a
// section are the same (markdownlint MD024).
func (a *Activity) Marshal() string {
	var b strings.Builder
	b.WriteString("---\n")
	fmt.Fprintf(&b, "kind: %s\n", a.Kind)
	fmt.Fprintf(&b, "accrued_cost: %s\n", formatCost(a.AccruedCost))
	fmt.Fprintf(&b, "accrued_seconds: %d\n", a.AccruedSeconds)
	fmt.Fprintf(&b, "tasks_completed: %d\n", a.TasksCompleted)
	if a.LastRun == "" {
		b.WriteString("last_run: \"\"\n")
	} else {
		fmt.Fprintf(&b, "last_run: %s\n", a.LastRun)
	}
	b.WriteString("---\n\n")
	fmt.Fprintf(&b, "# %s activity\n\n", strings.ToUpper(a.Kind[:1])+a.Kind[1:])
	b.WriteString("flai writes this document, one entry per activity, newest last; do not edit it by hand.\n\n")
	b.WriteString("## Log\n")
	count := map[string]int{}
	for _, e := range a.Entries {
		items := noItems
		if len(e.Items) > 0 {
			items = strings.Join(e.Items, ", ")
		}
		cost := formatCost(e.Cost) + usd
		if e.Estimated {
			cost += estimated
		}
		fmt.Fprintf(&b, "\n### %s\n\n%s%s\n", heading(count, e.At), summaryField, e.Summary)
		if e.Trigger != "" {
			fmt.Fprintf(&b, "%s%s\n", triggerField, e.Trigger)
		}
		fmt.Fprintf(&b, "%s%s\n%s%d\n%s%s\n", itemsField, items, secondsField, e.Seconds, costField, cost)
	}
	if len(a.Refusals) > 0 {
		fmt.Fprintf(&b, "\n%s\n", refusalsSection)
	}
	count = map[string]int{}
	for _, f := range a.Refusals {
		needs := f.Needs
		if needs == "" {
			needs = noPermission
		}
		fmt.Fprintf(&b, "\n### %s\n\n%s`%s`\n%s%s\n", heading(count, f.At), callField, f.Call, needsField, needs)
	}
	return b.String()
}

// heading is the heading of a block that ended at at: the time, and an
// ordinal after it when count has seen the second before.
func heading(count map[string]int, at time.Time) string {
	s := at.UTC().Format(TimeFormat)
	count[s]++
	if n := count[s]; n > 1 {
		return fmt.Sprintf("%s (%d)", s, n)
	}
	return s
}

// AppendActivity logs an activity of the strategic agent kind: it creates
// the document when missing, appends the entry last, adds its cost and
// seconds to the totals, counts it, moves last_run to its end when that is
// later, and writes the document atomically.
func (r *Repo) AppendActivity(kind string, e ActivityEntry) (*Activity, error) {
	if !IsActivityKind(kind) {
		return nil, unknownActivityKind(kind)
	}
	e, err := cleanActivityEntry(kind, e)
	if err != nil {
		return nil, err
	}
	return r.updateActivity(kind, "log the "+kind+"'s activity", func(a *Activity) {
		a.Entries = append(a.Entries, e)
		a.AccruedCost = roundCost(a.AccruedCost + e.Cost)
		a.AccruedSeconds += e.Seconds
		a.TasksCompleted++
		if at := e.At.Format(TimeFormat); at > a.LastRun {
			a.LastRun = at
		}
	})
}

// AppendRefusal logs a call flai guard refused the strategic agent kind
// under ## Refusals (S-0218): it creates the document when missing, appends
// the refusal last, leaves the totals alone, and writes the document
// atomically.
func (r *Repo) AppendRefusal(kind string, f ActivityRefusal) (*Activity, error) {
	if !IsActivityKind(kind) {
		return nil, unknownActivityKind(kind)
	}
	f.Call = strings.ReplaceAll(strings.Join(strings.Fields(f.Call), " "), "`", "'")
	f.Needs = strings.TrimSpace(f.Needs)
	switch {
	case f.Call == "":
		return nil, errors.New("a refusal needs the call refused")
	case f.At.IsZero():
		return nil, errors.New("a refusal needs the time of the call")
	case strings.ContainsAny(f.Needs, " \t\n") || f.Needs == noPermission:
		return nil, fmt.Errorf("permission %q must be one word, like orchestration.permissions.publish", f.Needs)
	}
	f.At = f.At.UTC().Truncate(time.Second)
	return r.updateActivity(kind, "log the "+kind+"'s refusal", func(a *Activity) {
		a.Refusals = append(a.Refusals, f)
	})
}

// updateActivity changes kind's activity document, created when missing,
// under its lock: it reads it, applies change, and writes it atomically once
// the markdown lint passes it. doing says what the change is for errors.
func (r *Repo) updateActivity(kind, doing string, change func(*Activity)) (*Activity, error) {
	unlock, err := r.lockActivity(kind)
	if err != nil {
		return nil, fmt.Errorf("cannot %s: %w", doing, err)
	}
	defer unlock()
	path := r.ActivityPath(kind)
	a, before := &Activity{Kind: kind, Path: path}, ""
	data, err := os.ReadFile(path)
	switch {
	case err == nil:
		before = string(data)
		if a, err = ReadActivity(path); err != nil {
			return nil, fmt.Errorf("cannot %s: %w", doing, err)
		}
	case !errors.Is(err, os.ErrNotExist):
		return nil, err
	}
	change(a)
	after := a.Marshal()
	if err := r.LintGuard(path, before, after); err != nil {
		return nil, err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, err
	}
	if err := atomicfile.WriteFile(path, []byte(after), 0o644); err != nil {
		return nil, err
	}
	return a, nil
}

// cleanActivityEntry refuses an entry flai cannot log and puts the rest in
// the form the log keeps: one-line summary and trigger, trimmed items, whole
// seconds of UTC, and the cost to the cent's hundredth.
func cleanActivityEntry(kind string, e ActivityEntry) (ActivityEntry, error) {
	e.Summary = strings.Join(strings.Fields(e.Summary), " ")
	e.Trigger = strings.Join(strings.Fields(e.Trigger), " ")
	if e.Summary == "" {
		return e, fmt.Errorf("an activity needs a summary: say in one line what the %s did", kind)
	}
	if e.At.IsZero() {
		return e, fmt.Errorf("an activity needs the time it ended")
	}
	if e.Seconds < 0 || e.Cost < 0 || math.IsNaN(e.Cost) || math.IsInf(e.Cost, 0) {
		return e, fmt.Errorf("an activity's seconds (%d) and cost (%v) must be zero or more", e.Seconds, e.Cost)
	}
	var items []string
	for _, it := range e.Items {
		it = strings.TrimSpace(it)
		if it == "" {
			continue
		}
		if strings.ContainsAny(it, ", \t\n") || it == noItems {
			return e, fmt.Errorf("item %q must be an ID, like S-0230, with no commas or spaces", it)
		}
		items = append(items, it)
	}
	e.Items = items
	e.At = e.At.UTC().Truncate(time.Second)
	e.Cost = roundCost(e.Cost)
	return e, nil
}

// activityLockAge is how old a lock on an activity document is when it is
// taken as left behind by a flai that died holding it.
const activityLockAge = 10 * time.Second

// lockActivity takes the lock on kind's activity document, a file in
// .flai-cache that one writer at a time creates, so that two flai processes
// writing the document at once, such as flai serve logging an activity and
// flai guard a refusal, do not lose each other's entry. It waits while
// another holds the lock, and removes a lock older than activityLockAge.
func (r *Repo) lockActivity(kind string) (unlock func(), err error) {
	path := filepath.Join(r.CacheDir(), "activity-"+kind+".lock")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, err
	}
	for {
		f, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o600)
		if err == nil {
			if err := f.Close(); err != nil {
				_ = os.Remove(path)
				return nil, err
			}
			return func() { _ = os.Remove(path) }, nil
		}
		if !errors.Is(err, os.ErrExist) {
			return nil, err
		}
		if st, err := os.Stat(path); err == nil && time.Since(st.ModTime()) > activityLockAge {
			_ = os.Remove(path)
			continue
		}
		time.Sleep(10 * time.Millisecond)
	}
}

func unknownActivityKind(kind string) error {
	return fmt.Errorf("unknown strategic agent %q: use one of %s", kind, strings.Join(ActivityKinds, ", "))
}

func roundCost(c float64) float64 { return math.Round(c*1e4) / 1e4 }

func formatCost(c float64) string { return strconv.FormatFloat(c, 'f', 4, 64) }
