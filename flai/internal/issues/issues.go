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

	"github.com/bytepunx/system-flow/flai/internal/template"
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

	Path string `yaml:"-" json:"path"`
	Body string `yaml:"-" json:"-"`
}

// Dir is <layout.design>/issues.
func Dir(r *workitem.Repo) string {
	return filepath.Join(r.Manifest.Dir(r.Root, "design"), Folder)
}

var idPattern = regexp.MustCompile(`^I-\d{3,}$`)

// Parse decodes an issue document.
func Parse(doc string) (*Issue, error) {
	fm, body, err := workitem.SplitFrontMatter(doc)
	if err != nil {
		return nil, err
	}
	var is Issue
	if err := yaml.UnmarshalWithOptions([]byte(fm), &is, yaml.Strict()); err != nil {
		return nil, err
	}
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

// NextID allocates the next issue ID, zero-padded to workitem.IDWidth.
func NextID(r *workitem.Repo) string {
	max := 0
	matches, _ := filepath.Glob(filepath.Join(Dir(r), "I-*.md"))
	for _, m := range matches {
		if n := num(strings.SplitN(filepath.Base(m), "-", 3)[0] + "-" + strings.SplitN(filepath.Base(m), "-", 3)[1]); n > max {
			max = n
		}
	}
	return fmt.Sprintf("I-%0*d", workitem.IDWidth, max+1)
}

// NewOptions describe an issue to record.
type NewOptions struct {
	Title, Class, Cost, Note string
	Now                      time.Time
}

// New creates an issue file with count 1.
func New(r *workitem.Repo, opt NewOptions) (*Issue, error) {
	if strings.TrimSpace(opt.Title) == "" {
		return nil, fmt.Errorf("a title is required")
	}
	if !contains(Classes, opt.Class) {
		return nil, fmt.Errorf("--class must be one of %s", strings.Join(Classes, ", "))
	}
	if opt.Cost != "" {
		if _, err := time.ParseDuration(opt.Cost); err != nil {
			return nil, fmt.Errorf("--cost must be a duration like 20m: %w", err)
		}
	}
	now := opt.Now.UTC().Format(workitem.TimeFormat)
	id := NextID(r)
	note := strings.TrimSpace(opt.Note)
	if note == "" {
		note = "First occurrence."
	}
	is := &Issue{ID: id, Title: opt.Title, Class: opt.Class, Status: "open", Count: 1, Cost: normalise(opt.Cost),
		FirstReported: now, LastReported: now, Updated: now,
		Path: filepath.Join(Dir(r), id+"-"+template.Slug(opt.Title)+".md"),
		Body: fmt.Sprintf("\n# %s %s\n\n## Description\n%s\n\n## Instances\n\n### %s\n%s\n\n## Remediation\n", id, opt.Title, opt.Title, now, note)}
	if err := is.Validate(); err != nil {
		return nil, err
	}
	return is, is.Save()
}

// Bump records another occurrence: count, last_reported, averaged cost, and
// a new instance in the body.
func Bump(is *Issue, cost, note string, now time.Time) error {
	if is.Status != "open" {
		return fmt.Errorf("%s is closed; reopen it by editing status, or record a new issue", is.ID)
	}
	ts := now.UTC().Format(workitem.TimeFormat)
	if cost != "" {
		d, err := time.ParseDuration(cost)
		if err != nil {
			return fmt.Errorf("--cost must be a duration like 20m: %w", err)
		}
		if is.Cost == "" {
			is.Cost = normalise(cost)
		} else {
			old, _ := time.ParseDuration(is.Cost)
			avg := (old*time.Duration(is.Count) + d) / time.Duration(is.Count+1)
			is.Cost = normalise(avg.Round(time.Minute).String())
		}
	}
	is.Count++
	is.LastReported, is.Updated = ts, ts
	if strings.TrimSpace(note) == "" {
		note = "Occurred again."
	}
	is.Body = insertInstance(is.Body, ts, strings.TrimSpace(note))
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

// insertInstance adds an occurrence at the end of the Instances section.
// Instances are headed by their timestamp, to the second. Two recorded in the
// same second share the heading, as narrative log entries do (I-0011): a
// second identical heading fails the duplicate-heading rule.
func insertInstance(body, ts, note string) string {
	head, tail := strings.TrimRight(body, "\n"), ""
	if idx := strings.Index(body, "\n## Remediation"); idx >= 0 {
		head, tail = strings.TrimRight(body[:idx], "\n"), body[idx+1:]
	}
	entry := fmt.Sprintf("### %s\n%s\n\n", ts, note)
	if last := strings.LastIndex(head, "\n### "); last >= 0 {
		heading, _, _ := strings.Cut(head[last+1:], "\n")
		if strings.TrimSpace(heading) == "### "+ts {
			entry = note + "\n\n"
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

// normalise trims zero components: 1h0m0s to 1h, 30m0s to 30m.
func normalise(d string) string {
	if d == "" {
		return ""
	}
	d = strings.TrimSuffix(d, "0s")
	if strings.HasSuffix(d, "h0m") {
		d = strings.TrimSuffix(d, "0m")
	}
	if d == "" {
		return "0s"
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
