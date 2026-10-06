// Package issues manages design/issues: recurring friction recorded with
// a count and cost. See design/system/continuous-improvement.md and ADR-0014.
package issues

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/goccy/go-yaml"

	"github.com/bytepunx/system-flow/flai/internal/execx"
	"github.com/bytepunx/system-flow/flai/internal/storygit"
	"github.com/bytepunx/system-flow/flai/internal/template"
	"github.com/bytepunx/system-flow/flai/internal/usage"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// Folder is the fixed subfolder under layout.design.
const Folder = "issues"

// SummaryFile is the generated table of open issues.
const SummaryFile = "summary.md"

// Classes is the closed list.
var Classes = []string{"defect", "blocker", "efficiency", "impression"}

// Issue is one file.
type Issue struct {
	ID            string `yaml:"id" json:"id"`
	Title         string `yaml:"title" json:"title"`
	Class         string `yaml:"class" json:"class"`
	Status        string `yaml:"status" json:"status"`
	Count         int    `yaml:"count" json:"count"`
	Cost          string `yaml:"cost" json:"cost,omitempty"` // average per occurrence, Go duration
	FirstReported string `yaml:"first_reported" json:"first_reported"`
	LastReported  string `yaml:"last_reported" json:"last_reported"`
	Updated       string `yaml:"updated" json:"updated"`
	// Usage is what strategic agents spent on the issue: its Strategic
	// entries alone (S-0227).
	Usage *usage.Usage `yaml:"usage" json:"usage,omitempty"`

	// Unknown is the front matter this flai does not know, kept for writing
	// back (S-0181).
	Unknown []workitem.Field `yaml:"-" json:"-"`

	Path string `yaml:"-" json:"path"`
	Body string `yaml:"-" json:"-"`
}

// Dir is <layout.design>/issues.
func Dir(r *workitem.Repo) string {
	return filepath.Join(r.Manifest.Dir(r.Root, "design"), Folder)
}

// SummaryPath is summary.md's path relative to the repository root,
// slash-separated, as git names it.
func SummaryPath(r *workitem.Repo) string {
	return filepath.ToSlash(filepath.Join(r.Manifest.Layout["design"], Folder, SummaryFile))
}

var idPattern = regexp.MustCompile(`^I-\d{3,}$`)

// Parse decodes an issue document, keeping front-matter fields this flai
// does not know in Unknown (S-0181).
func Parse(doc string) (*Issue, error) {
	fm, body, err := workitem.SplitFrontMatter(doc)
	if err != nil {
		return nil, err
	}
	var is Issue
	if err := yaml.Unmarshal([]byte(fm), &is); err != nil {
		return nil, err
	}
	is.Unknown = workitem.UnknownFields(fm, is)
	is.Body = body
	return &is, nil
}

// Read loads one file.
func Read(path string) (*Issue, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	is, err := Parse(string(data))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", path, err)
	}
	is.Path = path
	return is, nil
}

// Validate checks the schema of one issue.
func (is *Issue) Validate() error {
	var errs []string
	if !idPattern.MatchString(is.ID) {
		errs = append(errs, fmt.Sprintf("id %q must look like I-0001", is.ID))
	}
	if is.Title == "" {
		errs = append(errs, "title is required")
	}
	if !contains(Classes, is.Class) {
		errs = append(errs, fmt.Sprintf("class %q must be one of %s", is.Class, strings.Join(Classes, ", ")))
	}
	if is.Status != "open" && is.Status != "closed" {
		errs = append(errs, fmt.Sprintf("status %q must be open or closed", is.Status))
	}
	if is.Count < 1 {
		errs = append(errs, "count must be at least 1")
	}
	if is.Cost != "" {
		if _, err := time.ParseDuration(is.Cost); err != nil {
			errs = append(errs, fmt.Sprintf("cost %q is not a duration like 20m", is.Cost))
		}
	}
	for name, v := range map[string]string{"first_reported": is.FirstReported, "last_reported": is.LastReported, "updated": is.Updated} {
		if _, err := time.Parse(workitem.TimeFormat, v); err != nil {
			errs = append(errs, fmt.Sprintf("%s %q is not a UTC timestamp", name, v))
		}
	}
	errs = append(errs, usageErrors(is.Usage)...)
	if len(errs) > 0 {
		sort.Strings(errs)
		return fmt.Errorf("%s", strings.Join(errs, "; "))
	}
	return nil
}

// Marshal renders the document with the same hand-written style as items.
func (is *Issue) Marshal() string {
	var b strings.Builder
	b.WriteString("---\n")
	fmt.Fprintf(&b, "id: %s\n", is.ID)
	fmt.Fprintf(&b, "title: %s\n", workitem.Scalar(is.Title))
	fmt.Fprintf(&b, "class: %s\n", is.Class)
	fmt.Fprintf(&b, "status: %s\n", is.Status)
	fmt.Fprintf(&b, "count: %d\n", is.Count)
	if is.Cost != "" {
		fmt.Fprintf(&b, "cost: %s\n", is.Cost)
	}
	fmt.Fprintf(&b, "first_reported: %s\n", is.FirstReported)
	fmt.Fprintf(&b, "last_reported: %s\n", is.LastReported)
	fmt.Fprintf(&b, "updated: %s\n", is.Updated)
	b.WriteString(usageBlock(is.Usage))
	workitem.WriteFields(&b, is.Unknown)
	b.WriteString("---\n")
	b.WriteString(is.Body)
	return b.String()
}

// Save writes the issue.
func (is *Issue) Save() error {
	if err := os.MkdirAll(filepath.Dir(is.Path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(is.Path, []byte(is.Marshal()), 0o644)
}

// List loads every issue, sorted by ID. A missing folder is an empty list.
func List(r *workitem.Repo) ([]*Issue, error) {
	matches, _ := filepath.Glob(filepath.Join(Dir(r), "*.md"))
	var out []*Issue
	for _, m := range matches {
		if base := filepath.Base(m); base == SummaryFile || base == "README.md" {
			continue
		}
		is, err := Read(m)
		if err != nil {
			return nil, err
		}
		workitem.WarnUnknown(m, is.Unknown)
		out = append(out, is)
	}
	sort.Slice(out, func(i, j int) bool { return num(out[i].ID) < num(out[j].ID) })
	return out, nil
}

// Get finds one issue by ID, given in any padding (I-12, I-012, I-0012).
func Get(r *workitem.Repo, id string) (*Issue, error) {
	canon := workitem.CanonicalID(id)
	for _, cand := range []string{id, canon, fmt.Sprintf("I-%03s", strings.TrimPrefix(canon, "I-"))} {
		matches, _ := filepath.Glob(filepath.Join(Dir(r), cand+"-*.md"))
		if len(matches) > 0 {
			return Read(matches[0])
		}
	}
	return nil, fmt.Errorf("%s not found", id)
}

// NextID allocates the next issue ID, zero-padded to workitem.IDWidth: one
// past the highest in this checkout, on main, in any worktree, and on any
// story branch (I-0065). A nil run is execx.System.
func NextID(r *workitem.Repo, run execx.Runner) string {
	if run == nil {
		run = execx.System{}
	}
	var names []string
	matches, _ := filepath.Glob(filepath.Join(Dir(r), "I-*.md"))
	for _, m := range matches {
		names = append(names, filepath.Base(m))
	}
	if folder, err := filepath.Rel(r.Root, Dir(r)); err == nil {
		names = append(names, storygit.FolderNames(run, r, filepath.ToSlash(folder))...)
	}
	max := 0
	for _, name := range names {
		if m := fileIDPattern.FindStringSubmatch(name); m != nil {
			if n, _ := strconv.Atoi(m[1]); n > max {
				max = n
			}
		}
	}
	return fmt.Sprintf("I-%0*d", workitem.IDWidth, max+1)
}

// fileIDPattern is an issue file's name, capturing its number.
var fileIDPattern = regexp.MustCompile(`^I-(\d+)-.*\.md$`)

// NewOptions describe an issue to record. Story is the story the first
// instance belongs to, or empty for none. Impact, when given, is written as
// its Impact section. Report is the analysis report under design/analysis
// that found it, or empty for none: the instance names it and the
// Remediation section links it (S-0224). Runner reads the other worktrees and
// branches for the next number; nil is execx.System.
type NewOptions struct {
	Title, Class, Cost, Note, Story string
	Impact                          Impact
	Report                          string
	Now                             time.Time
	Runner                          execx.Runner
}

// check refuses what New would refuse, before anything is written, and
// returns the story and the report, normalised, or empty for none.
func (opt NewOptions) check(r *workitem.Repo) (story, report string, err error) {
	if strings.TrimSpace(opt.Title) == "" {
		return "", "", fmt.Errorf("a title is required")
	}
	if !contains(Classes, opt.Class) {
		return "", "", fmt.Errorf("--class must be one of %s", strings.Join(Classes, ", "))
	}
	if opt.Cost != "" {
		if _, err := time.ParseDuration(opt.Cost); err != nil {
			return "", "", fmt.Errorf("--cost must be a duration like 20m: %w", err)
		}
	}
	if story, err = storyID(opt.Story); err != nil {
		return "", "", err
	}
	if err := opt.Impact.Validate(); err != nil {
		return "", "", err
	}
	if report, err = reportPath(r, opt.Report); err != nil {
		return "", "", err
	}
	return story, report, nil
}

// New creates an issue file with count 1. Anything it refuses leaves nothing
// written.
func New(r *workitem.Repo, opt NewOptions) (*Issue, error) {
	story, report, err := opt.check(r)
	if err != nil {
		return nil, err
	}
	now := opt.Now.UTC().Format(workitem.TimeFormat)
	id := NextID(r, opt.Runner)
	note := strings.TrimSpace(opt.Note)
	if note == "" {
		note = "First occurrence."
	}
	is := &Issue{ID: id, Title: opt.Title, Class: opt.Class, Status: "open", Count: 1, Cost: normalise(opt.Cost),
		FirstReported: now, LastReported: now, Updated: now,
		Path: filepath.Join(Dir(r), id+"-"+template.Slug(opt.Title)+".md"),
		Body: fmt.Sprintf("\n# %s %s\n\n## Description\n%s\n\n## Instances\n\n### %s\n%s%s%s\n\n## Remediation\n", id, opt.Title, opt.Title, now, storyLine(story), reportLine(report), note)}
	if !opt.Impact.IsZero() {
		is.Body = setImpact(is.Body, opt.Impact)
	}
	if report != "" {
		linkReport(is, r.Root, report)
	}
	if err := is.Validate(); err != nil {
		return nil, err
	}
	return is, is.Save()
}

// BumpOptions describe another occurrence of an issue. Story is the story it
// belongs to, or empty for none; Cost its wall-clock cost, or empty. Impact
// replaces the Impact lines it gives and keeps the others. Report is the
// analysis report under design/analysis that found it, or empty for none: the
// instance names it and the Remediation section links it, once (S-0224).
type BumpOptions struct {
	Story, Cost, Note string
	Impact            Impact
	Report            string
	Now               time.Time
}

// Bump records another occurrence: count, last_reported, averaged cost, and
// a new instance in the body naming the story it belongs to, if any.
func Bump(is *Issue, story, cost, note string, now time.Time) error {
	return BumpWith(nil, is, BumpOptions{Story: story, Cost: cost, Note: note, Now: now})
}

// BumpWith records another occurrence as Bump does, and with it the impact
// and the report opt gives. r is the project the issue is in; it may be nil
// when opt names no report. Anything it refuses leaves the issue as it was.
func BumpWith(r *workitem.Repo, is *Issue, opt BumpOptions) error {
	if is.Status != "open" {
		return fmt.Errorf("%s is closed; reopen it by editing status, or record a new issue", is.ID)
	}
	story, err := storyID(opt.Story)
	if err != nil {
		return err
	}
	var d time.Duration
	if opt.Cost != "" {
		if d, err = time.ParseDuration(opt.Cost); err != nil {
			return fmt.Errorf("--cost must be a duration like 20m: %w", err)
		}
	}
	if err := opt.Impact.Validate(); err != nil {
		return err
	}
	var report string
	if strings.TrimSpace(opt.Report) != "" {
		if r == nil {
			return fmt.Errorf("no project given to find report %q in", opt.Report)
		}
		if report, err = reportPath(r, opt.Report); err != nil {
			return err
		}
	}
	ts := opt.Now.UTC().Format(workitem.TimeFormat)
	if opt.Cost != "" {
		if is.Cost == "" {
			is.Cost = normalise(opt.Cost)
		} else {
			old, _ := time.ParseDuration(is.Cost)
			avg := (old*time.Duration(is.Count) + d) / time.Duration(is.Count+1)
			is.Cost = normalise(avg.Round(time.Minute).String())
		}
	}
	is.Count++
	is.LastReported, is.Updated = ts, ts
	note := strings.TrimSpace(opt.Note)
	if note == "" {
		note = "Occurred again."
	}
	is.Body = insertInstance(is.Body, ts, story, reportLine(report)+note)
	if !opt.Impact.IsZero() {
		is.Body = setImpact(is.Body, opt.Impact)
	}
	if report != "" {
		linkReport(is, r.Root, report)
	}
	return is.Save()
}

// Close marks the issue closed with a reason under Remediation.
func Close(is *Issue, reason string, now time.Time) error {
	if strings.TrimSpace(reason) == "" {
		return fmt.Errorf("--reason is required")
	}
	if is.Status == "closed" {
		return fmt.Errorf("%s is already closed", is.ID)
	}
	ts := now.UTC().Format(workitem.TimeFormat)
	is.Status, is.Updated = "closed", ts
	is.Body = strings.TrimRight(is.Body, "\n") + fmt.Sprintf("\nClosed %s: %s\n", ts, strings.TrimSpace(reason))
	return is.Save()
}

// insertInstance adds an occurrence at the end of the Instances section,
// before the section after it, Impact or Remediation, its story line first.
// Instances are headed by their timestamp, to the second. Two recorded in the
// same second share the heading, as narrative log entries do (I-0011): a
// second identical heading fails the duplicate-heading rule. The joined note
// names its story unless that instance already does.
func insertInstance(body, ts, story, note string) string {
	head, tail := strings.TrimRight(body, "\n"), ""
	if _, end, ok := section(body, "## Instances"); ok && end < len(body) {
		head, tail = strings.TrimRight(body[:end], "\n"), body[end:]
	} else if !ok {
		if idx := strings.Index(body, "\n## Remediation"); idx >= 0 {
			head, tail = strings.TrimRight(body[:idx], "\n"), body[idx+1:]
		}
	}
	entry := fmt.Sprintf("### %s\n%s%s\n\n", ts, storyLine(story), note)
	if last := strings.LastIndex(head, "\n### "); last >= 0 {
		heading, instance, _ := strings.Cut(head[last+1:], "\n")
		if strings.TrimSpace(heading) == "### "+ts {
			line := storyLine(story)
			if line != "" && contains(strings.Split(instance, "\n"), strings.TrimSuffix(line, "\n")) {
				line = ""
			}
			entry = line + note + "\n\n"
		}
	}
	return head + "\n\n" + entry + tail
}

// Summary renders summary.md: open issues, most expensive first.
func Summary(list []*Issue, now time.Time) string {
	type row struct {
		total time.Duration
		text  string
	}
	var rows []row
	for _, is := range list {
		if is.Status != "open" {
			continue
		}
		avg, total := "-", "-"
		var tot time.Duration
		if is.Cost != "" {
			d, _ := time.ParseDuration(is.Cost)
			tot = d * time.Duration(is.Count)
			avg, total = normalise(is.Cost), normalise(tot.String())
		}
		rows = append(rows, row{tot, fmt.Sprintf("| [%s](%s) | %s | %s | %d | %s | %s | %s |", is.ID, filepath.Base(is.Path), is.Class, is.Title, is.Count, avg, total, is.LastReported)})
	}
	sort.SliceStable(rows, func(i, j int) bool { return rows[i].total > rows[j].total })
	var b strings.Builder
	fmt.Fprintf(&b, "---\ntitle: Open issues\nupdated: %s\nstatus: active\n---\n\n# Open issues\n\nMost expensive first. Total is count times average cost. Generated by `flai issue summary`.\n\n", now.UTC().Format(workitem.TimeFormat))
	b.WriteString("| ID | Class | Title | Count | Avg cost | Total cost | Last reported |\n|----|-------|-------|-------|----------|------------|---------------|\n")
	for _, r := range rows {
		b.WriteString(r.text + "\n")
	}
	return b.String()
}

// WriteSummary regenerates summary.md.
func WriteSummary(r *workitem.Repo, now time.Time) ([]*Issue, error) {
	list, err := List(r)
	if err != nil {
		return nil, err
	}
	if err := os.MkdirAll(Dir(r), 0o755); err != nil {
		return nil, err
	}
	return list, os.WriteFile(filepath.Join(Dir(r), SummaryFile), []byte(Summary(list, now)), 0o644)
}

// SummaryTable returns the table part of the summary for embedding, or ""
// when there are no open issues.
func SummaryTable(list []*Issue) string {
	open := 0
	for _, is := range list {
		if is.Status == "open" {
			open++
		}
	}
	if open == 0 {
		return ""
	}
	s := Summary(list, time.Time{})
	i := strings.Index(s, "| ID |")
	if i < 0 {
		return ""
	}
	return s[i:]
}

// normalise trims zero components: 1h0m0s to 1h, 30m0s to 30m. Seconds that
// are not zero stay: 30s is 30s, not 3.
func normalise(d string) string {
	if strings.HasSuffix(d, "m0s") {
		d = strings.TrimSuffix(d, "0s")
	}
	if strings.HasSuffix(d, "h0m") {
		d = strings.TrimSuffix(d, "0m")
	}
	return d
}

func num(id string) int {
	n, _ := strconv.Atoi(strings.TrimPrefix(id, "I-"))
	return n
}

func contains(list []string, s string) bool {
	for _, x := range list {
		if x == s {
			return true
		}
	}
	return false
}
