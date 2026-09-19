// Package adr records an architecture decision the way the conventions ask
// (S-0060): the number comes from the files present, the file is named and
// given its front matter, the index in design/adrs/README.md gets its row,
// an ADR it supersedes gets superseded_by set, and the author writes only
// the decision. Creation is one step that happens or does not, like an item
// made with a body (S-0059): flai check runs with everything in place, a
// finding it introduces undoes all of it, and what is kept is committed on
// its own.
package adr

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/check"
	"github.com/bytepunx/system-flow/flai/internal/docedit"
	"github.com/bytepunx/system-flow/flai/internal/execx"
	"github.com/bytepunx/system-flow/flai/internal/template"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// Statuses an author may give a new ADR.
var Statuses = []string{"proposed", "accepted"}

var (
	fileName = regexp.MustCompile(`^(\d{4})-.+\.md$`)
	statusRe = regexp.MustCompile(`(?m)^status: *\S+.*$`)
	dateRe   = regexp.MustCompile(`(?m)^date: *\S+.*$`)
	bySupRe  = regexp.MustCompile(`(?m)^superseded_by: *\[(.*)\] *$`)
)

// defaultBody is used when the project has no 0000-template.md.
const defaultBody = "## Context\n\n## Decision\n\n## Consequences\n\n## Alternatives considered\n"

// Options for New and Accept.
type Options struct {
	Title      string
	Status     string
	Supersedes []int
	Refines    []int
	Body       string // below the heading; the template's sections when empty
	Now        time.Time
	Autocommit bool
	Trailers   []string
}

// Result is what was recorded.
type Result struct {
	ID          string          `json:"id"`
	Number      int             `json:"number"`
	Title       string          `json:"title"`
	Status      string          `json:"status"`
	Path        string          `json:"path"`
	Changed     []string        `json:"changed"` // every path written, the new file first
	Warnings    []check.Finding `json:"warnings"`
	Committed   bool            `json:"committed"`
	Commit      string          `json:"commit,omitempty"`
	CommitError string          `json:"commit_error,omitempty"`
}

// Dir is the ADR folder of the checkout.
func Dir(repo *workitem.Repo) string {
	design := strings.TrimSuffix(repo.Manifest.Layout["design"], "/")
	if design == "" {
		design = "design"
	}
	return filepath.Join(repo.Root, design, "adrs")
}

// files maps ADR numbers to file names, the template (0000) left out.
func files(dir string) (map[int]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	out := map[int]string{}
	for _, e := range entries {
		if m := fileName.FindStringSubmatch(e.Name()); m != nil && !e.IsDir() {
			if n, _ := strconv.Atoi(m[1]); n > 0 {
				out[n] = e.Name()
			}
		}
	}
	return out, nil
}

// NextNumber is one more than the highest number among the files present.
// Gaps are not filled, and no counter kept in a document is consulted.
func NextNumber(dir string) (int, error) {
	fs, err := files(dir)
	if err != nil && !os.IsNotExist(err) {
		return 0, err
	}
	next := 1
	for n := range fs {
		if n >= next {
			next = n + 1
		}
	}
	return next, nil
}

// TemplateBody is what the project's 0000-template.md puts below its heading.
func TemplateBody(repo *workitem.Repo) string {
	data, err := os.ReadFile(filepath.Join(Dir(repo), "0000-template.md"))
	if err != nil {
		return defaultBody
	}
	_, body, err := workitem.SplitFrontMatter(string(data))
	if err != nil {
		return defaultBody
	}
	lines := strings.Split(body, "\n")
	for i, l := range lines {
		if strings.HasPrefix(l, "# ") {
			return strings.TrimLeft(strings.Join(lines[i+1:], "\n"), "\n")
		}
	}
	return strings.TrimLeft(body, "\n")
}

func adrID(n int) string { return fmt.Sprintf("ADR-%04d", n) }
func num4(n int) string  { return fmt.Sprintf("%04d", n) }

// yamlString quotes a title only when YAML would read it as something else.
func yamlString(s string) string {
	if s == "" || strings.ContainsAny(s, ":#[]{}&*!|>'\"%@`,") || strings.HasPrefix(s, "-") || strings.HasPrefix(s, "?") || s != strings.TrimSpace(s) {
		b, _ := json.Marshal(s)
		return string(b)
	}
	return s
}

func idList(ns []int) string {
	ids := make([]string, len(ns))
	for i, n := range ns {
		ids[i] = adrID(n)
	}
	return "[" + strings.Join(ids, ", ") + "]"
}

func numList(ns []int) string {
	parts := make([]string, len(ns))
	for i, n := range ns {
		parts[i] = num4(n)
	}
	if len(parts) <= 1 {
		return strings.Join(parts, "")
	}
	return strings.Join(parts[:len(parts)-1], ", ") + " and " + parts[len(parts)-1]
}

func uniqueSorted(ns []int) []int {
	seen := map[int]bool{}
	var out []int
	for _, n := range ns {
		if !seen[n] {
			seen[n] = true
			out = append(out, n)
		}
	}
	sort.Ints(out)
	return out
}

// snapshot remembers files so a refusal can put them back.
type snapshot struct {
	was     map[string][]byte
	created []string
}

func (s *snapshot) keep(path string) error {
	if _, ok := s.was[path]; ok {
		return nil
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	s.was[path] = data
	return nil
}

func (s *snapshot) undo() error {
	for _, p := range s.created {
		if err := os.Remove(p); err != nil && !os.IsNotExist(err) {
			return err
		}
	}
	for p, data := range s.was {
		if err := os.WriteFile(p, data, 0o644); err != nil {
			return err
		}
	}
	return nil
}

// New records a decision.
func New(repo *workitem.Repo, r execx.Runner, opt Options) (*Result, error) {
	title := strings.Join(strings.Fields(opt.Title), " ")
	if title == "" {
		return nil, fmt.Errorf("rule: an ADR needs a title: the decision, as a sentence")
	}
	status := opt.Status
	if status == "" {
		status = "proposed"
	}
	if status != "proposed" && status != "accepted" {
		return nil, fmt.Errorf("rule: status %q must be proposed or accepted", status)
	}
	if opt.Now.IsZero() {
		opt.Now = time.Now()
	}
	dir := Dir(repo)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, err
	}
	present, err := files(dir)
	if err != nil {
		return nil, err
	}
	supersedes, refines := uniqueSorted(opt.Supersedes), uniqueSorted(opt.Refines)
	for _, n := range append(append([]int{}, supersedes...), refines...) {
		if _, ok := present[n]; !ok {
			return nil, fmt.Errorf("rule: there is no %s to supersede or refine", adrID(n))
		}
	}
	for _, n := range supersedes {
		for _, m := range refines {
			if n == m {
				return nil, fmt.Errorf("rule: %s is both superseded and refined; choose one", adrID(n))
			}
		}
	}
	n, err := NextNumber(dir)
	if err != nil {
		return nil, err
	}
	slug := template.Slug(title)
	if slug == "" {
		slug = "decision"
	}
	if len(slug) > 80 {
		slug = strings.TrimRight(slug[:80], "-")
	}
	name := num4(n) + "-" + slug + ".md"
	path := filepath.Join(dir, name)

	before, err := check.Run(repo, opt.Now)
	if err != nil {
		return nil, err
	}
	snap := &snapshot{was: map[string][]byte{}}
	index := filepath.Join(dir, "README.md")

	body := strings.TrimSpace(opt.Body)
	if body == "" {
		body = strings.TrimSpace(TemplateBody(repo))
	}
	var fm strings.Builder
	fmt.Fprintf(&fm, "---\nid: %s\ntitle: %s\nstatus: %s\ndate: %s\nsupersedes: %s\nsuperseded_by: []\n", adrID(n), yamlString(title), status, opt.Now.UTC().Format("2006-01-02"), idList(supersedes))
	if len(refines) > 0 {
		fmt.Fprintf(&fm, "refines: %s\n", idList(refines))
	}
	fm.WriteString("---\n\n")
	content := fm.String() + "# " + adrID(n) + " " + title + "\n\n" + body + "\n"

	fail := func(err error) (*Result, error) {
		if uerr := snap.undo(); uerr != nil {
			return nil, errors.Join(err, fmt.Errorf("undoing the partial ADR failed: %w", uerr))
		}
		return nil, err
	}
	snap.created = append(snap.created, path)
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		return fail(err)
	}
	changed := []string{path}

	// The index row, with the note the index already uses.
	note := status
	if len(supersedes) > 0 {
		note += ", supersedes " + numList(supersedes)
	}
	if len(refines) > 0 {
		note += ", refines " + numList(refines)
	}
	row := fmt.Sprintf("| [%s](%s) | %s | %s |", num4(n), name, strings.ReplaceAll(title, "|", "\\|"), note)
	if _, err := os.Stat(index); err == nil {
		if err := snap.keep(index); err != nil {
			return fail(err)
		}
		if err := os.WriteFile(index, []byte(addRow(string(snap.was[index]), row)), 0o644); err != nil {
			return fail(err)
		}
		changed = append(changed, index)
	}

	// superseded_by is the one edit the conventions allow to an accepted ADR.
	for _, old := range supersedes {
		oldPath := filepath.Join(dir, present[old])
		if err := snap.keep(oldPath); err != nil {
			return fail(err)
		}
		text := string(snap.was[oldPath])
		m := bySupRe.FindStringSubmatch(text)
		switch {
		case m == nil:
			return fail(fmt.Errorf("rule: %s has no superseded_by line in its front matter to set", adrID(old)))
		case strings.Contains(m[1], adrID(n)):
		case strings.TrimSpace(m[1]) == "":
			text = bySupRe.ReplaceAllString(text, "superseded_by: ["+adrID(n)+"]")
		default:
			text = bySupRe.ReplaceAllString(text, "superseded_by: ["+strings.TrimSpace(m[1])+", "+adrID(n)+"]")
		}
		if err := os.WriteFile(oldPath, []byte(text), 0o644); err != nil {
			return fail(err)
		}
		changed = append(changed, oldPath)
		if data, err := os.ReadFile(index); err == nil {
			_ = os.WriteFile(index, []byte(noteRow(string(data), old, "superseded by "+num4(n))), 0o644)
		}
	}

	res := &Result{ID: adrID(n), Number: n, Title: title, Status: status, Path: relTo(repo.Root, path), Warnings: []check.Finding{}}
	for _, p := range changed {
		res.Changed = append(res.Changed, relTo(repo.Root, p))
	}
	if refused, err := introduced(repo, before, opt.Now); err != nil {
		return fail(err)
	} else if len(refused) > 0 {
		if err := snap.undo(); err != nil {
			return nil, fmt.Errorf("the check refused the ADR and undoing it failed: %w", err)
		}
		return nil, &docedit.RefusedError{Path: res.Path, Reason: fmt.Sprintf("flai check has %d finding(s) with this ADR; nothing was created", len(refused)), Findings: refused}
	}
	if opt.Autocommit && repo.Manifest.Autocommit() {
		res.commit(r, repo.Root, fmt.Sprintf("docs: %s %s", res.ID, title), opt.Trailers)
	}
	return res, nil
}

// Accept turns a proposed ADR into an accepted one, dated today.
func Accept(repo *workitem.Repo, r execx.Runner, n int, opt Options) (*Result, error) {
	if opt.Now.IsZero() {
		opt.Now = time.Now()
	}
	dir := Dir(repo)
	present, err := files(dir)
	if err != nil {
		return nil, err
	}
	name, ok := present[n]
	if !ok {
		return nil, fmt.Errorf("rule: there is no %s", adrID(n))
	}
	path, index := filepath.Join(dir, name), filepath.Join(dir, "README.md")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	fm, _, err := workitem.SplitFrontMatter(string(data))
	if err != nil {
		return nil, err
	}
	cur := strings.Fields(strings.TrimPrefix(statusRe.FindString(fm), "status:"))
	if len(cur) == 0 || cur[0] != "proposed" {
		was := "unknown"
		if len(cur) > 0 {
			was = cur[0]
		}
		return nil, fmt.Errorf("rule: %s is %s; only a proposed ADR can be accepted", adrID(n), was)
	}
	before, err := check.Run(repo, opt.Now)
	if err != nil {
		return nil, err
	}
	snap := &snapshot{was: map[string][]byte{}}
	if err := snap.keep(path); err != nil {
		return nil, err
	}
	newFM := statusRe.ReplaceAllString(fm, "status: accepted")
	newFM = dateRe.ReplaceAllString(newFM, "date: "+opt.Now.UTC().Format("2006-01-02"))
	text := strings.Replace(string(data), fm, newFM, 1)
	if err := os.WriteFile(path, []byte(text), 0o644); err != nil {
		return nil, err
	}
	changed := []string{path}
	if _, err := os.Stat(index); err == nil {
		if err := snap.keep(index); err == nil {
			_ = os.WriteFile(index, []byte(setRowStatus(string(snap.was[index]), n, "proposed", "accepted")), 0o644)
			changed = append(changed, index)
		}
	}
	title := ""
	for _, l := range strings.Split(newFM, "\n") {
		if strings.HasPrefix(l, "title:") {
			title = strings.Trim(strings.TrimSpace(strings.TrimPrefix(l, "title:")), `"`)
		}
	}
	res := &Result{ID: adrID(n), Number: n, Title: title, Status: "accepted", Path: relTo(repo.Root, path), Warnings: []check.Finding{}}
	for _, p := range changed {
		res.Changed = append(res.Changed, relTo(repo.Root, p))
	}
	if refused, err := introduced(repo, before, opt.Now); err != nil || len(refused) > 0 {
		_ = snap.undo()
		if err != nil {
			return nil, err
		}
		return nil, &docedit.RefusedError{Path: res.Path, Reason: fmt.Sprintf("flai check has %d finding(s) with %s accepted; nothing was changed", len(refused), res.ID), Findings: refused}
	}
	if opt.Autocommit && repo.Manifest.Autocommit() {
		res.commit(r, repo.Root, fmt.Sprintf("docs: %s accepted", res.ID), opt.Trailers)
	}
	return res, nil
}

// introduced returns the findings that were not there before.
func introduced(repo *workitem.Repo, before *check.Result, now time.Time) ([]check.Finding, error) {
	after, err := check.Run(repo, now)
	if err != nil {
		return nil, err
	}
	known := map[string]bool{}
	for _, f := range before.Findings {
		known[f.Rule+"\x00"+f.Path+"\x00"+f.Message] = true
	}
	var out []check.Finding
	for _, f := range after.Findings {
		if !known[f.Rule+"\x00"+f.Path+"\x00"+f.Message] {
			out = append(out, f)
		}
	}
	return out, nil
}

var rowRe = regexp.MustCompile(`^\| *\[(\d{4})\]\(`)

// addRow puts the row after the last row of the index table, or starts a
// table at the end when the index has none.
func addRow(index, row string) string {
	lines := strings.Split(strings.TrimRight(index, "\n"), "\n")
	last := -1
	for i, l := range lines {
		if rowRe.MatchString(l) || strings.HasPrefix(l, "|--") || strings.HasPrefix(l, "| ADR") {
			last = i
		}
	}
	if last < 0 {
		return strings.Join(lines, "\n") + "\n\n| ADR | Title | Status |\n|-----|-------|--------|\n" + row + "\n"
	}
	out := append(append(append([]string{}, lines[:last+1]...), row), lines[last+1:]...)
	return strings.Join(out, "\n") + "\n"
}

// noteRow appends a note to the status cell of one ADR's row.
func noteRow(index string, n int, note string) string {
	lines := strings.Split(index, "\n")
	for i, l := range lines {
		if m := rowRe.FindStringSubmatch(l); m != nil && m[1] == num4(n) && !strings.Contains(l, note) {
			trimmed := strings.TrimRight(l, " ")
			if strings.HasSuffix(trimmed, "|") {
				lines[i] = strings.TrimRight(strings.TrimSuffix(trimmed, "|"), " ") + ", " + note + " |"
			}
		}
	}
	return strings.Join(lines, "\n")
}

// setRowStatus replaces the leading status word of one ADR's status cell.
func setRowStatus(index string, n int, from, to string) string {
	lines := strings.Split(index, "\n")
	for i, l := range lines {
		if m := rowRe.FindStringSubmatch(l); m != nil && m[1] == num4(n) {
			cells := strings.Split(l, "|")
			if len(cells) >= 3 {
				k := len(cells) - 2
				cell := strings.TrimSpace(cells[k])
				if strings.HasPrefix(cell, from) {
					cells[k] = " " + to + strings.TrimPrefix(cell, from) + " "
					lines[i] = strings.Join(cells, "|")
				}
			}
		}
	}
	return strings.Join(lines, "\n")
}

func relTo(root, path string) string {
	if rel, err := filepath.Rel(root, path); err == nil && !strings.HasPrefix(rel, "..") {
		return filepath.ToSlash(rel)
	}
	return path
}

// commit records the changed paths, and only them. A failed commit is not a
// failed creation: the files are written and valid.
func (res *Result) commit(r execx.Runner, root, msg string, trailers []string) {
	if r == nil {
		return
	}
	if _, err := r.Run(root, "git", "rev-parse", "--is-inside-work-tree"); err != nil {
		return
	}
	if len(trailers) > 0 {
		msg += "\n\n" + strings.Join(trailers, "\n")
	}
	if _, err := r.Run(root, "git", append([]string{"add", "--"}, res.Changed...)...); err != nil {
		res.CommitError = err.Error()
		return
	}
	if _, err := r.Run(root, "git", append([]string{"commit", "--quiet", "-m", msg, "--"}, res.Changed...)...); err != nil {
		res.CommitError = err.Error()
		return
	}
	sha, err := r.Run(root, "git", "rev-parse", "--short", "HEAD")
	if err != nil {
		res.CommitError = err.Error()
		return
	}
	res.Committed, res.Commit = true, strings.TrimSpace(sha)
}
